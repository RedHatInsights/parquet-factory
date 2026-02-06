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
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/parquet-factory/dataaggregator/mock"
	"github.com/RedHatInsights/parquet-factory/metrics"
	"github.com/RedHatInsights/parquet-factory/reportreader"
)

func TestConsumePartition(t *testing.T) {
	var (
		timeout        = 2 * time.Second
		limitTimestamp = time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	)

	type test struct {
		name              string
		partitions        []int32
		offsets           map[int32]int64
		limitReached      map[int32]bool
		consumerErr       sarama.KError
		expectTimeout     bool
		lastCall          string
		messageOffsets    []int64
		messageTimestamps time.Time
		headers           map[string]string
		maxRecords        *int // optional, defaults to 10 if nil
	}

	tests := []test{
		{
			name:       "tracked offset is out of bounds",
			partitions: []int32{0, 1},
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			consumerErr:   sarama.ErrOffsetOutOfRange,
			expectTimeout: false,
			lastCall:      "ConsumePartition",
		},
		{
			name:       "tracked offset is inside of bounds",
			partitions: []int32{0, 1},
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			expectTimeout: false,
			lastCall:      "ConsumePartition",
		},
		{
			name:       "messages in chan and messageTimestamps < limitTimestamp",
			partitions: []int32{0, 1},
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			expectTimeout:     true,
			lastCall:          "ConsumePartition",
			messageOffsets:    []int64{1},
			messageTimestamps: limitTimestamp.Add(-1 * time.Hour),
		},
		{
			name:       "messages in chan and messageTimestamps > limitTimestamp",
			partitions: []int32{0, 1},
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			expectTimeout:     false,
			lastCall:          "ConsumePartition",
			messageOffsets:    []int64{1},
			messageTimestamps: limitTimestamp.Add(1 * time.Hour),
		},
		{
			name:       "messages in chan and messageTimestamps < limitTimestamp but they have stop header",
			partitions: []int32{0, 1},
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			expectTimeout:     false,
			lastCall:          "ConsumePartition",
			messageOffsets:    []int64{1},
			messageTimestamps: limitTimestamp.Add(-1 * time.Hour),
			headers:           map[string]string{"stop": "true"},
		},
		{
			name:       "limit already reached before consuming any messages (maxRecords = 0)",
			partitions: []int32{0, 1},
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			expectTimeout: false,
			lastCall:      "Partitions",
			maxRecords:    func() *int { i := 0; return &i }(),
		},
	}

	err := metrics.InitMetrics("testEnv")
	assert.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var testMockConsumer = NewMockConsumer()

			maxRecords := 10
			if tc.maxRecords != nil {
				maxRecords = *tc.maxRecords
			}

			config := reportreader.MockConfiguration{
				Topic:          "test_topic",
				GroupID:        "test_group",
				MaxRecords:     maxRecords,
				LimitTimestamp: limitTimestamp,
				Aggregator:     &mock.Aggregator{},
				OffsetManager: &mockOffsetManager{
					partitionOffsetManager: &mockPartitionOffsetManager{},
				},
				Consumer: &testMockConsumer,
			}

			config.PartitionTracker = reportreader.NewMockPartitionTracker(
				map[string]map[int32]int64{
					"test_topic": tc.offsets,
				},
				map[string]map[int32]bool{
					"test_topic": tc.limitReached,
				},
			)

			sut, err := reportreader.NewMockKafkaConsumer(config)
			assert.NoError(t, err)

			testMockConsumer.partitions = tc.partitions
			if tc.consumerErr != sarama.ErrNoError {
				testMockConsumer.errConsumePartition = tc.consumerErr
			}

			if len(tc.messageOffsets) > 0 {
				messageChan := make(chan *sarama.ConsumerMessage, 1)
				testMockConsumer.partitionConsumer = mockPartitionConsumer{
					messageChan: messageChan,
				}

				go fillMessageChan(
					messageChan,
					config.Topic,
					tc.partitions,
					tc.messageOffsets,
					tc.headers,
					tc.messageTimestamps,
				)

				defer func() {
					for len(messageChan) > 0 {
						<-messageChan // empty the channel
					}
					close(messageChan)
				}()
			}

			ctx := sut.Start()

			if tc.expectTimeout {
				assert.Error(t, waitForContext(ctx, timeout), "expected the context not to be canceled")
			} else {
				assert.NoError(t, waitForContext(ctx, timeout), "expected the context to be canceled")
				assert.Equal(t, "context canceled", ctx.Err().Error())
			}

			assert.Equal(t, tc.lastCall, testMockConsumer.lastCall)
		})
	}
}

func waitForContext(ctx context.Context, timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		fmt.Println("timeout waiting for context to cancel", ctx.Err())
		return errors.New("timeout waiting for context to cancel")
	case <-ctx.Done():
		return nil
	}
}

func fillMessageChan(messageChan chan *sarama.ConsumerMessage, topic string, partitions []int32, offsets []int64, headers map[string]string, timestamp time.Time) {
	recordsHeaders := make([]*sarama.RecordHeader, 0, len(headers))
	for k, v := range headers {
		recordsHeaders = append(recordsHeaders, &sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	once := 0
	for _, partition := range partitions {
		for _, offset := range offsets {
			fmt.Printf("Sending offset %d to partition %d\n", offset, partition)
			msg := sarama.ConsumerMessage{
				Headers:   recordsHeaders,
				Timestamp: timestamp,
				Key:       []byte(`key`),
				Value:     []byte(fmt.Sprintf(`{"path": "test/path%d.gz"}`, once)),
				Topic:     topic,
				Partition: partition,
				Offset:    offset,
			}

			once++
			messageChan <- &msg
			fmt.Printf("Sent offset %d to partition %d\n", offset, partition)
		}
	}
}
