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
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

const (
	recordingOffsetError = "Error recording offset"
)

// OffsetTracker defines the methods to retrieve the cached offset for a given topic and partition
type OffsetTracker interface {
	GetTopics() []string
	GetPartitionsForTopic(topic string) []int32
	GetOffset(topic string, partition int32) (int64, error)
}

// PartitionTracker allow to store the read offset for given topic/partition pairs while it can be committed
type PartitionTracker struct {
	offsetMutex  sync.RWMutex
	limitsMutex  sync.RWMutex
	offsets      map[string]map[int32]int64
	limitReached map[string]map[int32]bool
}

// NewPartitionTracker creates an PartitionTracker ready to be used
func NewPartitionTracker() *PartitionTracker {
	return &PartitionTracker{
		offsets:      map[string]map[int32]int64{},
		limitReached: map[string]map[int32]bool{},
	}
}

// RecordOffset store the offset that should be committed when needed
func (pt *PartitionTracker) RecordOffset(msg *sarama.ConsumerMessage) error {
	log.Debug().
		Int64(offsetTag, msg.Offset).
		Int32(partitionTag, msg.Partition).
		Str(topicTag, msg.Topic).
		Msg("Storing the offset to be committed in the future")

	pt.offsetMutex.Lock()
	defer pt.offsetMutex.Unlock()

	if partitionOffsets, ok := pt.offsets[msg.Topic]; ok {
		switch {
		case msg.Offset == -2:
			log.Warn().Msg("never consummed from this topic before")
		case partitionOffsets[msg.Partition] > msg.Offset:
			err := errors.New("fetched a smaller offset than the cached one")
			log.Info().
				Err(err).
				Int64(offsetTag, msg.Offset).
				Int64("cached offset", partitionOffsets[msg.Partition]).
				Int32(partitionTag, msg.Partition).
				Str(topicTag, msg.Topic).
				Msg(recordingOffsetError)
			log.Error().Err(err).Msg(recordingOffsetError)
			return err
		}
		partitionOffsets[msg.Partition] = msg.Offset
	} else {
		log.Warn().Msg("No offsets stored for this topic")
		pt.offsets[msg.Topic] = map[int32]int64{
			msg.Partition: msg.Offset,
		}
	}
	log.Debug().
		Int64(offsetTag, msg.Offset).
		Int32(partitionTag, msg.Partition).
		Str(topicTag, msg.Topic).
		Msg("Stored the offset to be committed in the future")
	return nil
}

// TrackPartition register the topic:partition in order to be aware of it
// This function will return an error if the topic:partition is already tracked, else nil
func (pt *PartitionTracker) TrackPartition(topic string, partition int32) error {
	pt.limitsMutex.Lock()
	defer pt.limitsMutex.Unlock()

	if partitions, ok := pt.limitReached[topic]; ok {
		if _, ok := partitions[partition]; ok {
			// if topic & partition already exist in the map, error
			return fmt.Errorf("this partition (%s:%d) is already tracked", topic, partition)
		}
		// topic exist, but new partition
		partitions[partition] = false
	} else {
		// topic didn't exist, create map for it
		pt.limitReached[topic] = map[int32]bool{
			partition: false,
		}
	}
	return nil
}

// GetOffset retrieve the current cached offset for a given topic and partition
func (pt *PartitionTracker) GetOffset(topic string, partition int32) (int64, error) {
	pt.offsetMutex.RLock()
	defer pt.offsetMutex.RUnlock()

	if partitionOffsets, ok := pt.offsets[topic]; ok {
		if offset, ok := partitionOffsets[partition]; ok {
			return offset, nil
		}
		return -1, fmt.Errorf("no offset cached for topic %s in the partition %d", topic, partition)
	}
	return -1, fmt.Errorf("no offset cached for any partition in the topic %s", topic)
}

// GetTopics retrieve an array of topics that has some offset cached
func (pt *PartitionTracker) GetTopics() []string {
	pt.offsetMutex.RLock()
	defer pt.offsetMutex.RUnlock()

	topics := make([]string, 0, len(pt.offsets))
	for topic := range pt.offsets {
		topics = append(topics, topic)
	}

	return topics
}

// GetPartitionsForTopic retrieve an array of partition numbers that have some offset cached in a given topic
func (pt *PartitionTracker) GetPartitionsForTopic(topic string) []int32 {
	pt.offsetMutex.RLock()
	defer pt.offsetMutex.RUnlock()

	if topicPartitions, ok := pt.offsets[topic]; ok {
		partitions := make([]int32, 0, len(topicPartitions))
		for part := range topicPartitions {
			partitions = append(partitions, part)
		}
		return partitions
	}

	return []int32{}
}
