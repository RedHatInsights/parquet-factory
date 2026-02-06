// Copyright 2022 Red Hat, Inc
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

package testhelpers

import (
	"testing"

	"github.com/IBM/sarama"
)

// NewBrokerWithTopic returns a sarama.MockBroker that can handle minimal requests
func NewBrokerWithTopic(t testing.TB, topic string) *sarama.MockBroker {
	broker := sarama.NewMockBroker(t, 0)
	broker.SetHandlerByMap(getHandlersMapForMockConsumer(t, broker, topic))
	return broker
}

// NewBrokerWith2Topics returns a sarama.MockBroker that can handle the initialization of 2 brokers
func NewBrokerWith2Topics(t testing.TB, topic1, topic2 string) *sarama.MockBroker {
	broker := sarama.NewMockBroker(t, 0)
	broker.SetHandlerByMap(getHandlersMapForMockConsumer2Topics(t, broker, topic1, topic2))
	return broker
}

// NewBrokerWithoutTopic returns a sarama.MockBroker that will fail
func NewBrokerWithoutTopic(t testing.TB, topic string) *sarama.MockBroker {
	broker := sarama.NewMockBroker(t, 0)
	broker.SetHandlerByMap(getHandlersMapTopicDoesNotExists(t, broker, topic))
	return broker
}

// getHandlersMapForMockConsumer returns handlers for mock broker to successfully create a new consumer
func getHandlersMapForMockConsumer(t testing.TB, mockBroker *sarama.MockBroker, topicName string) map[string]sarama.MockResponse {
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

// getHandlersMapForMockConsumer2Topics returns handlers for mock broker to successfully create a new consumer
func getHandlersMapForMockConsumer2Topics(t testing.TB, mockBroker *sarama.MockBroker, topic1Name, topic2Name string) map[string]sarama.MockResponse {
	return map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()).
			SetLeader(topic1Name, 0, mockBroker.BrokerID()).
			SetLeader(topic2Name, 0, mockBroker.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(topic1Name, 0, -1, 0).
			SetOffset(topic1Name, 0, -2, 0).
			SetOffset(topic2Name, 0, -1, 0).
			SetOffset(topic2Name, 0, -2, 0),
		"FetchRequest": sarama.NewMockFetchResponse(t, 1),
		"FindCoordinatorRequest": sarama.NewMockFindCoordinatorResponse(t).
			SetCoordinator(sarama.CoordinatorGroup, "", mockBroker),
		"OffsetFetchRequest": sarama.NewMockOffsetFetchResponse(t).
			SetOffset("", topic1Name, 0, 0, "", sarama.ErrNoError).
			SetOffset("", topic2Name, 0, 0, "", sarama.ErrNoError),
	}
}

// getHandlersMapTopicDoesNotExists return a metadata response with a non existent topic.
func getHandlersMapTopicDoesNotExists(t testing.TB, mockBroker *sarama.MockBroker, topicName string) map[string]sarama.MockResponse {
	return map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()).
			SetLeader(topicName+"_non_existent", 0, mockBroker.BrokerID()),
	}
}
