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

package reportreader_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RedHatInsights/parquet-factory/testhelpers"

	"github.com/RedHatInsights/parquet-factory/reportaggregators/rulereportaggregator"

	"github.com/IBM/sarama"
	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/reportreader"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	topic := "topic"

	testCases := []struct {
		name          string
		setupBroker   func(t *testing.T) *sarama.MockBroker
		getConfig     func(brokerAddr string) conf.KafkaConfig
		expectError   bool
		errorContains string
	}{
		{
			name: "broker is OK",
			setupBroker: func(t *testing.T) *sarama.MockBroker {
				return testhelpers.NewBrokerWithTopic(t, topic)
			},
			getConfig: func(brokerAddr string) conf.KafkaConfig {
				return conf.KafkaConfig{
					Addresses: []string{brokerAddr},
					Topic:     topic,
				}
			},
			expectError: false,
		},
		{
			name: "topic doesn't exist",
			setupBroker: func(t *testing.T) *sarama.MockBroker {
				return testhelpers.NewBrokerWithoutTopic(t, topic)
			},
			getConfig: func(brokerAddr string) conf.KafkaConfig {
				return conf.KafkaConfig{
					Addresses: []string{brokerAddr},
					Topic:     topic,
				}
			},
			expectError:   true,
			errorContains: "kafka server: Request was for a topic or partition that does not exist on this broker",
		},
		{
			name:        "can't create the client",
			setupBroker: nil,
			getConfig: func(brokerAddr string) conf.KafkaConfig {
				return conf.KafkaConfig{
					Addresses: []string{"not an address"},
					Topic:     topic,
				}
			},
			expectError:   true,
			errorContains: "kafka: client has run out of available brokers to talk to",
		},
		{
			name:        "can't create the client and use TLS",
			setupBroker: nil,
			getConfig: func(brokerAddr string) conf.KafkaConfig {
				return conf.KafkaConfig{
					Addresses: []string{"not an address"},
					Topic:     topic,
					CertPath:  "../testdata/cert.pem",
				}
			},
			expectError:   true,
			errorContains: "kafka: client has run out of available brokers to talk to",
		},
		{
			name: "use SSL without cert",
			setupBroker: func(t *testing.T) *sarama.MockBroker {
				return testhelpers.NewBrokerWithTopic(t, topic)
			},
			getConfig: func(brokerAddr string) conf.KafkaConfig {
				return conf.KafkaConfig{
					Addresses:        []string{brokerAddr},
					Topic:            topic,
					SecurityProtocol: "SSL",
				}
			},
			expectError:   true,
			errorContains: "no cert path provided. Skip",
		},
		{
			name: "use SASL",
			setupBroker: func(t *testing.T) *sarama.MockBroker {
				return testhelpers.NewBrokerWithTopic(t, topic)
			},
			getConfig: func(brokerAddr string) conf.KafkaConfig {
				return conf.KafkaConfig{
					Addresses:        []string{brokerAddr},
					Topic:            topic,
					SecurityProtocol: "SASL_SSL",
					SaslMechanism:    "SCRAM-SHA512",
					ClientID:         "a-client-id",
				}
			},
			expectError:   true,
			errorContains: "invalid configuration (Net.SASL.Password must not be empty when SASL is enabled)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var mockBroker *sarama.MockBroker
			var brokerAddr string

			if tc.setupBroker != nil {
				mockBroker = tc.setupBroker(t)
				defer mockBroker.Close()
				brokerAddr = mockBroker.Addr()
			}

			config := tc.getConfig(brokerAddr)
			consumer, err := reportreader.New(config, rulereportaggregator.NewRulesReportAggregator())

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				assert.NoError(t, err)
				if consumer != nil {
					defer func() {
						assert.NoError(t, consumer.Close())
					}()
				}
			}
		})
	}
}

func TestStart(t *testing.T) {
	var (
		timeout                        = 2 * time.Second
		testMockConsumer               = mockConsumer{}
		testMockPartitionOffsetManager = mockPartitionOffsetManager{}
		testMockOffsetManager          = mockOffsetManager{
			partitionOffsetManager: &testMockPartitionOffsetManager,
		}
		testPartitionTracker = &reportreader.PartitionTracker{}
	)
	config := reportreader.MockConfiguration{
		Topic:            "test_topic",
		GroupID:          "test_group",
		MaxRecords:       10,
		Aggregator:       nil,
		OffsetManager:    &testMockOffsetManager,
		Consumer:         &testMockConsumer,
		PartitionTracker: testPartitionTracker,
	}

	t.Run("if partitions cannot be listed, the context should be cancelled", func(t *testing.T) {
		sut, err := reportreader.NewMockKafkaConsumer(config)
		assert.NoError(t, err)

		testMockConsumer.errPartitions = errors.New("cannot list partitions")
		defer func() {
			testMockConsumer.errPartitions = nil
		}()

		ctx := sut.Start()

		assert.NoError(t, waitForContext(ctx, timeout), "expected the context to be canceled")

		assert.Equal(t, "context canceled", ctx.Err().Error())
		assert.Equal(t, "Partitions", testMockConsumer.lastCall)
	})

	t.Run("if partitions can be listed but not consumed, the context should be cancelled", func(t *testing.T) {
		sut, err := reportreader.NewMockKafkaConsumer(config)
		assert.NoError(t, err)
		assert.Equal(t, "Partitions", testMockConsumer.lastCall)

		testMockConsumer.partitions = []int32{0, 1}
		testMockConsumer.errConsumePartition = errors.New("error consuming partitions")
		testPartitionTracker = reportreader.NewMockPartitionTracker(
			map[string]map[int32]int64{},
			map[string]map[int32]bool{},
		)
		defer func() {
			testMockConsumer.partitions = []int32{}
			testMockConsumer.errConsumePartition = nil
			testPartitionTracker = &reportreader.PartitionTracker{}
		}()

		ctx := sut.Start()

		assert.NoError(t, waitForContext(ctx, timeout), "expected the context to be canceled")

		assert.Equal(t, "context canceled", ctx.Err().Error())
		assert.Equal(t, "Partitions", testMockConsumer.lastCall)
	})
}

// GetHandlersMapForMockConsumer returns handlers for mock broker to successfully create a new consumer
func GetHandlersMapForMockConsumer(t testing.TB, mockBroker *sarama.MockBroker, topicName string) map[string]sarama.MockResponse {
	return map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()).
			SetLeader(topicName, 0, mockBroker.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(topicName, 0, -1, 0).
			SetOffset(topicName, 0, -2, 0),
		"FetchRequest": sarama.NewMockFetchResponse(t, 1),
		"FindCoordinatorRequest": sarama.NewMockFindCoordinatorResponse(t).
			SetCoordinator(sarama.CoordinatorGroup, "", mockBroker),
		"OffsetFetchRequest": sarama.NewMockOffsetFetchResponse(t).
			SetOffset("", topicName, 0, 0, "", sarama.ErrNoError),
	}
}

// GetHandlersMapTopicDoesNotExists return a metadata response with a non existent topic.
func GetHandlersMapTopicDoesNotExists(t testing.TB, mockBroker *sarama.MockBroker, topicName string) map[string]sarama.MockResponse {
	return map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()).
			SetLeader(topicName+"_non_existent", 0, mockBroker.BrokerID()),
	}
}
