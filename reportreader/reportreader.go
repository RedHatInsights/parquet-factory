// Copyright 2021 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reportreader

import (
	"context"
	"crypto/sha512"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/RedHatInsights/parquet-factory/dataaggregator"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/metrics"
	"github.com/RedHatInsights/parquet-factory/utils"

	tlsutils "github.com/RedHatInsights/insights-operator-utils/tls"
)

const (
	offsetTag    = "offset"
	partitionTag = "partition"
	topicTag     = "topic"
	archiveTag   = "archive_path"
)

// KafkaConsumer represents the implementation of a Consumer
type KafkaConsumer struct {
	Topic             string
	GroupID           string
	consumer          sarama.Consumer      // Defer close it
	offsetManager     sarama.OffsetManager // Defer close it
	wg                sync.WaitGroup
	Aggregator        dataaggregator.DataAggregator
	partitionTracker  *PartitionTracker
	limits            limitChecker
	consumerTimeout   time.Duration
	processedMessages *utils.ArchivePathSet
}

// New constructs a new implementation of a KafkaConsumer
func New(config conf.KafkaConfig, aggregator dataaggregator.DataAggregator) (*KafkaConsumer, error) {
	// Create sarama.Config object
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = "parquet-factory"
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = false
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	if strings.Contains(config.SecurityProtocol, "SSL") {
		log.Info().Msgf("Security protocol uses TLS: %s", config.SecurityProtocol)
		saramaConfig.Net.TLS.Enable = true
	}

	if strings.EqualFold(config.SecurityProtocol, "SSL") {
		// configuring TLS if needed
		tlsConfig, err := tlsutils.NewTLSConfig(config.CertPath)
		if err != nil {
			return nil, err
		}
		saramaConfig.Net.TLS.Config = tlsConfig
	} else if strings.HasPrefix(config.SecurityProtocol, "SASL_") {
		log.Info().Msg("Configuring SASL authentication")
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = config.ClientID
		saramaConfig.Net.SASL.Password = config.ClientSecret

		if strings.EqualFold(config.SaslMechanism, sarama.SASLTypeSCRAMSHA512) {
			log.Info().Msg("Configuring SCRAM")
			saramaConfig.Net.SASL.Handshake = true
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &SCRAMClient{HashGeneratorFcn: sha512.New}
			}
		}
	}

	// sarama client used by consumer group
	client, err := sarama.NewClient(
		config.Addresses,
		saramaConfig,
	) // Defer close it
	if err != nil {
		log.Error().Err(err).Msg("Unable to create a new Kafka client")
		return nil, err
	}

	c, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create a new Kafka consumer")
		return nil, err
	}

	offsetManager, err := sarama.NewOffsetManagerFromClient(config.GroupID, client)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create an OffsetManager")
		return nil, err
	}

	log.Debug().Int("Shift", conf.GetConfiguration().TimeShift).Msg("Kafka consumer with timeshift")

	consumer := &KafkaConsumer{
		Topic:            config.Topic,
		GroupID:          config.GroupID,
		consumer:         c,
		offsetManager:    offsetManager,
		wg:               sync.WaitGroup{},
		Aggregator:       aggregator,
		partitionTracker: NewPartitionTracker(),
		limits: limitChecker{
			limitTimestamp: utils.GetHourOnly(
				time.Now().Add(
					time.Duration(conf.GetConfiguration().TimeShift) * time.Minute)),
			maxRecords:       config.MaxRecords,
			consumedMessages: 0,
			mutex:            sync.RWMutex{},
		},
		consumerTimeout:   time.Duration(config.ConsumerTimeout) * time.Second,
		processedMessages: utils.NewArchivePathSet(),
	}

	err = consumer.getInitialOffsetTracker()
	return consumer, err
}

// Start init consuming Kafka records
func (c *KafkaConsumer) Start() context.Context {
	log.Info().Msg("Starting consumer...")
	context, cancel := context.WithCancel(context.Background())

	partitions, err := c.consumer.Partitions(c.Topic)
	if err != nil {
		log.Error().Err(err).Msgf(`Cannot retrieve partitions list for topic "%s"`, c.Topic)
		cancel()
		return context
	}

	go func() {
		for _, partition := range partitions {
			c.wg.Add(1)
			go func(partition int32) {
				if err := c.consumePartition(partition); err != nil {
					log.Error().Err(err).Msgf("error consuming partition %d", partition)
					cancel()
				}
			}(partition)
		}

		log.Info().Msg("Waiting for the partition consumers to finish")
		if waitTimeout(&c.wg, c.consumerTimeout) {
			log.Info().Msg("Timed out waiting for consumers.")
		} else {
			log.Info().Msg("All the partition consumers finished.")
		}
		cancel()
	}()

	return context
}

//gocyclo:ignore
func (c *KafkaConsumer) consumePartition(p int32) error {
	defer c.wg.Done()
	for {
		// Check if we've already reached the limit before starting to consume
		if !c.limits.CanConsumeMore() {
			log.Info().Msg("Limit reached before consuming any messages. Exiting consumer.")
			return nil
		}

		// Get the offset from the offset tracker for the current partition
		offset, err := c.partitionTracker.GetOffset(c.Topic, p)
		if err != nil {
			log.Error().Err(err).Msg("Stopping consumer goroutine. Impossible to fetch the commit")
			return err
		}

		pConsumer, err := c.initPartitionConsumer(p, offset)
		if err != nil {
			log.Warn().Err(err).Str(topicTag, c.Topic).Int32(partitionTag, p).Msg("Cannot consume topic/partition")
			return err
		}

		for m := range pConsumer.Messages() {
			// check it's offset is lower than the marked offset
			if c.checkOffset(m) {
				consumerLog(log.Warn(), m, "This offset is lower than the stored one")
				continue
			}

			// check limits
			if !c.limits.CheckMessage(m) {
				consumerLog(log.Info(), m, "FINISH")
				return nil
			}

			if checkHeaders(m) {
				consumerLog(log.Info(), m, "I've been asked to aggregate all messages.FINISH")
				return c.markMessage(m)
			}

			// check if message has been already processed in current run
			path, err := utils.GetPathFromRawMsg(m.Value)
			if err != nil {
				log.Error().Err(err).Msg("can't retrieve path from kafka message, skipping")
				continue
			}
			if !c.processedMessages.Add(path) {
				log.Warn().Msg("factory was about to duplicate a row, skipping")
				continue
			}

			// Process message
			consumerLog(log.Info(), m, "message processed")
			if err := c.Aggregator.Handle(m.Value); err != nil {
				log.Error().Err(err).Msg("Unable to dispatch event")
				continue
			}
			if err = c.markMessage(m); err != nil {
				return err
			}
		}

		pConsumer.AsyncClose()
	}
}

func (c *KafkaConsumer) checkOffset(m *sarama.ConsumerMessage) bool {
	lastOffsetStored, err := c.partitionTracker.GetOffset(m.Topic, m.Partition)
	if err != nil {
		consumerLog(log.Error().Err(err), m, "Can't retrieve the stored offset")
		return false
	}
	log.Debug().
		Int64("message offset", m.Offset).
		Int64("stored offset", lastOffsetStored).
		Msg("checkOffset")
	return m.Offset <= lastOffsetStored
}

func (c *KafkaConsumer) markMessage(m *sarama.ConsumerMessage) error {
	// Increase the number of messages processed on this topic
	c.limits.MessageProcessed()
	consumerLog(log.Debug(), m, "marked processed")
	metrics.OffsetProcessed.Inc()
	// Mark offset
	if err := c.partitionTracker.RecordOffset(m); err != nil {
		consumerLog(log.Error().Err(err), m, "error marking consumed")
		return err
	}
	consumerLog(log.Debug(), m, "marked consumed")

	metrics.OffsetConsummed.Inc()

	return nil
}

func (c *KafkaConsumer) getInitialOffsetTracker() error {
	log.Info().Msg("Getting initial offsets")
	partitions, err := c.consumer.Partitions(c.Topic)
	if err != nil {
		return err
	}

	defer func() {
		if err := c.offsetManager.Close(); err != nil {
			log.Error().Err(err).Msg("Unable to Close the OffsetManager")
		}
	}()

	for _, partition := range partitions {
		partManager, err := c.offsetManager.ManagePartition(c.Topic, partition)
		if err != nil {
			return err
		}
		defer partManager.AsyncClose()

		if err := c.partitionTracker.TrackPartition(c.Topic, partition); err != nil {
			return err
		}

		offset, _ := partManager.NextOffset()
		msg := &sarama.ConsumerMessage{
			Topic:     c.Topic,
			Partition: partition,
			Offset:    offset,
		}
		if err = c.partitionTracker.RecordOffset(msg); err != nil {
			return err
		}
	}

	return nil
}

// initPartitionConsumer facilitates the creation of a PartitionConsumer, handling common errors
func (c *KafkaConsumer) initPartitionConsumer(partition int32, offset int64) (sarama.PartitionConsumer, error) {
	var pConsumer sarama.PartitionConsumer

	offset = changeOffsetIfNotNewestOrOldest(offset)
	log.Debug().
		Int64(offsetTag, offset).
		Int32(partitionTag, partition).
		Msg("Start consuming partition")
	pConsumer, err := c.consumer.ConsumePartition(c.Topic, partition, offset)
	if err != nil {
		serr, ok := err.(sarama.KError)

		if !ok || serr != sarama.ErrOffsetOutOfRange {
			// The error is not related to sarama lib or it is not the expected
			return nil, err
		}

		// Stored offset is out of retention policy bounds
		log.Info().Str(topicTag, c.Topic).Int32(partitionTag, partition).Msg("Stored offset is out of bounds. Skip until a valid one")
		pConsumer, err = c.consumer.ConsumePartition(c.Topic, partition, sarama.OffsetOldest)

		if err != nil {
			return nil, err
		}
	}
	if pConsumer == nil {
		return nil, errors.New("pConsumer is nil")
	}
	return pConsumer, nil
}

func changeOffsetIfNotNewestOrOldest(offset int64) (newOffset int64) {
	// Keep the value if it is a logic one (OffsetNewest or OffsetOldest)
	if offset == sarama.OffsetOldest || offset == sarama.OffsetNewest {
		newOffset = offset
	} else {
		newOffset = offset + 1
	}
	return
}

// Close closes the Kafka consumer and releases all resources
func (c *KafkaConsumer) Close() error {
	if c.consumer != nil {
		if err := c.consumer.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Kafka consumer")
			return err
		}
	}
	return nil
}
