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

package main_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/RedHatInsights/parquet-factory/reportaggregators/rulereportaggregator"
	"github.com/RedHatInsights/parquet-factory/testhelpers"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/tisnik/go-capture"

	main "github.com/RedHatInsights/parquet-factory/cmd/parquet-factory"
	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/dataaggregator/mock"
	"github.com/RedHatInsights/parquet-factory/reportreader"
)

var timeout = 1 * time.Second

func TestPrintVersionInfo(t *testing.T) {
	output, err := capture.StandardOutput(func() {
		log.Logger = log.Output(zerolog.New(os.Stdout))
		main.PrintVersionInfo()
	})

	assert.NoError(t, err)
	assert.Contains(t, output, "Version: ")
	assert.Contains(t, output, "Build time: ")
	assert.Contains(t, output, "Branch: ")
	assert.Contains(t, output, "Commit: ")
}

func TestWaitForConsumers(t *testing.T) {
	done := make(chan bool)

	go func() {
		main.WaitForConsumers([]*reportreader.KafkaConsumer{})
		done <- true
	}()

	select {
	case <-time.After(timeout):
		t.Fatal("test run out of time")
	case <-done:
	}
}

func TestCreateConsumerError(t *testing.T) {
	kafkaCfg := &conf.KafkaConfig{}
	consumer, err := main.CreateKafkaConsumer(kafkaCfg, nil)
	assert.Nil(t, consumer)
	assert.Error(t, err)
}

func TestCreateConsumer(t *testing.T) {
	topic := "topic"
	mockBroker := testhelpers.NewBrokerWithTopic(t, topic)
	defer mockBroker.Close()

	kafkaCfg := &conf.KafkaConfig{
		Addresses: []string{mockBroker.Addr()},
		Topic:     topic,
	}

	consumer, err := main.CreateKafkaConsumer(kafkaCfg, rulereportaggregator.NewRulesReportAggregator())
	assert.NotNil(t, consumer)
	assert.NoError(t, err)
	if consumer != nil {
		defer func() {
			assert.NoError(t, consumer.Close())
		}()
	}
	fmt.Println(err)
}

func TestStartMetricsMetrics(t *testing.T) {
	err := main.StartMetrics()
	assert.NoError(t, err)
}

func TestStartConsumers(t *testing.T) {
	topic1 := "t1"
	topic2 := "t2"
	mockBroker := testhelpers.NewBrokerWith2Topics(t, topic1, topic2)
	defer mockBroker.Close()

	cfg := &conf.Config{
		RulesKafkaConsumer: conf.KafkaConfig{
			Addresses:  []string{mockBroker.Addr()},
			Topic:      topic1,
			MaxRecords: 0, // Stop immediately without consuming any messages
		},
	}

	err := main.StartKafkaCollection(*cfg, nil)
	assert.NoError(t, err)
}

func TestCommitOffset(t *testing.T) {
	topic1 := "t1"
	topic2 := "t2"
	mockBroker := testhelpers.NewBrokerWith2Topics(t, topic1, topic2)
	defer mockBroker.Close()

	mockWorkingAggregator := mock.Aggregator{}
	mockFaultyAggregator := mock.FaultyAggregator{}
	cfg := &conf.Config{
		RulesKafkaConsumer: conf.KafkaConfig{
			Addresses: []string{mockBroker.Addr()},
			Topic:     topic1,
		},
	}

	commitConsumer, err := main.CreateKafkaConsumer(&cfg.RulesKafkaConsumer, &mockWorkingAggregator)
	if err != nil {
		assert.NoError(t, err)
	}
	if commitConsumer != nil {
		defer func() {
			assert.NoError(t, commitConsumer.Close())
		}()
	}

	notCommitConsumer, err := main.CreateKafkaConsumer(&cfg.RulesKafkaConsumer, &mockFaultyAggregator) //MBY BAD TODO
	if err != nil {
		assert.NoError(t, err)
	}
	if notCommitConsumer != nil {
		defer func() {
			assert.NoError(t, notCommitConsumer.Close())
		}()
	}

	var consumers = make([]*reportreader.KafkaConsumer, 2)
	consumers[0] = commitConsumer
	consumers[1] = notCommitConsumer
	errorCount := 0
	commitCount := 0
	for _, consumer := range consumers {
		_, err := consumer.Aggregator.WriteResults(nil)
		if err != nil {
			errorCount++
		} else {
			// commit offset only if no errors occurred
			if err := consumer.OffsetCommit(); err != nil {
				assert.NoError(t, err)
			}
			commitCount++
		}
	}

	assert.Equal(t, errorCount, 1)
	assert.Equal(t, commitCount, 1)
}
