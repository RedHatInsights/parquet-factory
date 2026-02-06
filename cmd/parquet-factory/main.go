/*
Copyright © 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Entry point to the insights results aggregator service.
//
// The service contains consumer (usually Kafka consumer) that consumes
// messages from given source, processes those messages and stores them
// in configured data store. It also starts REST API servers with
// endpoints that expose several types of information: list of organizations,
// list of clusters for given organization, and cluster health.
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/RedHatInsights/insights-operator-utils/logger"
	"github.com/RedHatInsights/insights-operator-utils/metrics/push"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/dataaggregator"
	"github.com/RedHatInsights/parquet-factory/metrics"
	"github.com/RedHatInsights/parquet-factory/reportaggregators/rulereportaggregator"
	"github.com/RedHatInsights/parquet-factory/reportreader"
	"github.com/RedHatInsights/parquet-factory/s3writer"
)

const (
	// SUCCESS indicates that the program finishes without any problem
	SUCCESS = 0
	// BADCONFIG indicates that the program finished because a misconfiguration
	BADCONFIG = 1
	// CONSUMERERROR indicates that the program finishes because a problem in Kafka consumers
	CONSUMERERROR = 2
	// METRICSERROR indicates that the program finishes because it cannot push the metrics
	METRICSERROR = 3
	// S3ERROR indicates that the program finishes because there was an error related to S3
	S3ERROR = 4
	// defaultConfigFilename specifies the name of the default config file
	defaultConfigFilename = "config"

	topicTag = "topic"
)

var (
	// buildVersion contains the major.minor version of the CLI client
	buildVersion = "*not set*"

	// buildTime contains timestamp when the CLI client has been built
	buildTime = "*not set*"

	// buildBranch contains Git branch used to build this application
	buildBranch = "*not set*"

	// buildCommit contains Git commit used to build this application
	buildCommit = "*not set*"
)

func createKafkaConsumer(cfg *conf.KafkaConfig, aggregator dataaggregator.DataAggregator) (*reportreader.KafkaConsumer, error) {
	consumer, err := reportreader.New(*cfg, aggregator)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to create the Kafka consumer for topic %s", cfg.Topic)
		return nil, err
	}

	return consumer, nil
}

func startKafkaCollection(config conf.Config, s3Writer *s3writer.S3Writer) error {
	metrics.State.Set(metrics.ConnectToKafka)
	ruleHitsAggregator := rulereportaggregator.NewRulesReportAggregator()

	ruleConsumer, err := createKafkaConsumer(&config.RulesKafkaConsumer, ruleHitsAggregator)
	if err != nil {
		log.Error().Err(err).Msg("cannot create consumer")
		return err
	}
	metrics.State.Set(metrics.Consume)
	var consumers [1]*reportreader.KafkaConsumer
	consumers[0] = ruleConsumer

	waitForConsumers(consumers[:])

	log.Info().Msg("Consumers ready for writing")
	for _, consumer := range consumers {
		log.Info().Str(topicTag, consumer.Topic).Msg("running aggregator")
		numFilesWritten, err := consumer.Aggregator.WriteResults(s3Writer)
		if err != nil {
			log.Error().Str(topicTag, consumer.Topic).
				Err(err).Msg("aggregator failure, no results were stored")
		} else {
			// commit offset only if no errors occurred
			if numFilesWritten == 0 {
				log.Info().Msg("No files needed to be written")
			}
			log.Info().Str(topicTag, consumer.Topic).Msg("committing offset")
			if err := consumer.OffsetCommit(); err != nil {
				log.Error().Err(err).Str(topicTag, consumer.Topic).
					Msg("Some problem happened when commiting offsets")
			}
		}
	}

	return nil
}

func waitForConsumers(consumers []*reportreader.KafkaConsumer) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// consumers coordination channel
	done := make(chan struct{}, len(consumers))

	for _, consumer := range consumers {
		go func(c *reportreader.KafkaConsumer) {
			<-consumer.Start().Done()
			done <- struct{}{}
		}(consumer)
	}

	remaining := len(consumers)
	for remaining > 0 {
		select {
		case sig := <-sigterm:
			log.Info().Msgf("Signal %v received. Finishing without storing or commiting offset", sig)
			endProgram(SUCCESS)

		case <-done:
			remaining--
			log.Info().Msgf("Topic context cancelled. Remaining topics %d", remaining)
		}
	}
}

func main() {
	if err := conf.LoadConfiguration(defaultConfigFilename); err != nil {
		log.Error().Msgf("Configuration cannot be loaded: %s", err)
		endProgram(BADCONFIG)
	}
	config := conf.GetConfiguration()

	if err := startMetrics(); err != nil {
		endProgram(METRICSERROR)
	}

	if err := logger.InitZerolog(
		conf.GetLoggingConfiguration(),
		conf.GetCloudWatchConfiguration(),
		conf.GetSentryConfiguration()); err != nil {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		log.Warn().Err(err).Msg(`Logger configuration cannot be loaded. Using "debug=true" by default`)
	}

	log.Info().Msg("Parquet service")
	printVersionInfo()
	s3Writer, err := s3writer.New(conf.GetS3Configuration())
	if err != nil {
		endProgram(S3ERROR)
	}

	metrics.State.Set(metrics.Consume)

	if err = startKafkaCollection(config, s3Writer); err != nil {
		endProgram(CONSUMERERROR)
	}

	log.Info().Msg("See you")
	endProgram(SUCCESS)
}

func endProgram(status int) {
	if status != SUCCESS && status != BADCONFIG {
		metrics.ErrorCount.Inc()
	}
	metrics.State.Set(metrics.Idle)

	metricsConf := conf.GetMetricsConfiguration()
	err := push.SendMetrics(metricsConf.Job, metricsConf.GatewayURL, metricsConf.GatewayAuthToken)
	if err != nil {
		log.Error().Err(err).Msg("Cannot push metrics")
	}

	logger.CloseZerolog()
	os.Exit(status)
}
