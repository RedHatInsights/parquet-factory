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
	"testing"

	"github.com/IBM/sarama"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/parquet-factory/reportreader"
)

const topicName string = "topicName"

func TestPartitionTrackerTrackPartition(t *testing.T) {
	partition := int32(0)

	t.Run("valid topic and partition", func(t *testing.T) {
		sut := reportreader.NewPartitionTracker()
		err := sut.TrackPartition(topicName, partition)
		assert.NoError(t, err)
	})

	t.Run("partition and topic already being tracked", func(t *testing.T) {
		sut := reportreader.NewPartitionTracker()
		assert.NoError(t, sut.TrackPartition(topicName, partition))
		// Second call to track partition for the same topic/partition should fail
		err := sut.TrackPartition(topicName, partition)
		assert.Error(t, err)
	})

	t.Run("this partition doesn't exist but the partitions map has already been initialized", func(t *testing.T) {
		var newPartition int32 = 1
		sut := reportreader.NewPartitionTracker()
		assert.NoError(t, sut.TrackPartition(topicName, partition))
		// A call to track a new partition for the same topic shouldn't fail
		err := sut.TrackPartition(topicName, newPartition)
		assert.NoError(t, err)
	})
}

type testCase struct {
	name             string
	offset           int64
	expectedOffset   int64
	expectingAnError bool
}

func TestPartitionTrackerGetOffset(t *testing.T) {
	partition := int32(0)

	sut := reportreader.NewPartitionTracker()
	err := sut.TrackPartition(topicName, partition)
	assert.NoError(t, err)

	// get non tracked offset
	_, err = sut.GetOffset(topicName, partition)
	assert.Error(t, err)

	// Initialize the partition with offset 0
	initMsg := sarama.ConsumerMessage{
		Topic:     topicName,
		Partition: partition,
		Offset:    0,
	}

	err = sut.RecordOffset(&initMsg)
	assert.NoError(t, err)

	testCases := []testCase{
		{
			name:           "sarama OffsetOldest",
			offset:         sarama.OffsetOldest,
			expectedOffset: sarama.OffsetOldest,
		},
		{
			name:           "sarama OffsetNewest",
			offset:         sarama.OffsetNewest,
			expectedOffset: sarama.OffsetNewest,
		},
		{
			name:           "a valid offset",
			offset:         1000,
			expectedOffset: 1000, // Should return stored offset, not the next one
		},
		{
			name:             "smaller offset (negative lag)",
			offset:           -17,
			expectedOffset:   1000, // the previous one
			expectingAnError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeMsg := sarama.ConsumerMessage{
				Topic:     topicName,
				Partition: partition,
				Offset:    tc.offset,
			}

			err := sut.RecordOffset(&fakeMsg)
			if tc.expectingAnError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "a negative lag should generate an error")
			}

			// get the offset tracked from the fakeMsg
			offset, err := sut.GetOffset(topicName, partition)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedOffset, offset)
		})
	}
	t.Run("get offset for a non tracked partition", func(t *testing.T) {
		// get offset for a non tracked partition
		_, err = sut.GetOffset(topicName, int32(1000))
		assert.Error(t, err)
	})
}

func TestPartitionTrackerGetTopics(t *testing.T) {
	partition := int32(0)

	sut := reportreader.NewPartitionTracker()
	err := sut.TrackPartition(topicName, partition)
	assert.NoError(t, err)

	// If no offset has been recorded yet, GetTopics return should be empty
	topics := sut.GetTopics()
	assert.Equal(t, []string{}, topics)

	// Record a fake offset
	initMsg := sarama.ConsumerMessage{
		Topic:     topicName,
		Partition: partition,
		Offset:    0,
	}

	err = sut.RecordOffset(&initMsg)
	assert.NoError(t, err)

	topics = sut.GetTopics()
	assert.Equal(t, []string{topicName}, topics)
}

func TestPartitionTrackerGetPartitionsForTopic(t *testing.T) {
	partition := int32(0)

	sut := reportreader.NewPartitionTracker()
	err := sut.TrackPartition(topicName, partition)
	assert.NoError(t, err)

	// If no offset has been recorded yet, GetPartitionsForTopic return should be empty
	partitions := sut.GetPartitionsForTopic(topicName)
	assert.Equal(t, []int32{}, partitions)

	// Record a fake offset
	initMsg := sarama.ConsumerMessage{
		Topic:     topicName,
		Partition: partition,
		Offset:    0,
	}

	err = sut.RecordOffset(&initMsg)
	assert.NoError(t, err)

	partitions = sut.GetPartitionsForTopic(topicName)
	assert.Equal(t, []int32{0}, partitions)
}
