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

package reportreader

import (
	"sync"
	"time"

	"github.com/RedHatInsights/parquet-factory/dataaggregator"
	"github.com/RedHatInsights/parquet-factory/utils"

	"github.com/IBM/sarama"
)

type MockConfiguration struct {
	Topic            string
	GroupID          string
	MaxRecords       int
	LimitTimestamp   time.Time
	Aggregator       dataaggregator.DataAggregator
	OffsetManager    sarama.OffsetManager
	Consumer         sarama.Consumer
	PartitionTracker *PartitionTracker
}

// NewMockKafkaConsumer returns a KafkaConsumer with a stub sarama.OffsetManager,
// sarama.Consumer and PartitionTracker.
func NewMockKafkaConsumer(config MockConfiguration) (*KafkaConsumer, error) {
	consumer := &KafkaConsumer{
		Topic:            config.Topic,
		GroupID:          config.GroupID,
		consumer:         config.Consumer,
		offsetManager:    config.OffsetManager,
		wg:               sync.WaitGroup{},
		Aggregator:       config.Aggregator,
		partitionTracker: config.PartitionTracker,
		limits: limitChecker{
			limitTimestamp:   config.LimitTimestamp,
			maxRecords:       config.MaxRecords,
			consumedMessages: 0,
			mutex:            sync.RWMutex{},
		},
		processedMessages: utils.NewArchivePathSet(),
	}

	err := consumer.getInitialOffsetTracker()
	return consumer, err
}

// NewMockPartitionTracker returns a PartitionTracker given the specified offsets.
func NewMockPartitionTracker(offsets map[string]map[int32]int64, limitReached map[string]map[int32]bool) *PartitionTracker {
	return &PartitionTracker{
		offsets:      offsets,
		limitReached: limitReached,
	}
}
