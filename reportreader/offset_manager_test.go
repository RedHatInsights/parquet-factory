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

	"github.com/stretchr/testify/assert"
	"github.com/RedHatInsights/parquet-factory/metrics"
	"github.com/RedHatInsights/parquet-factory/reportreader"
)

func TestOffsetCommit(t *testing.T) {
	type test struct {
		name                             string
		offsets                          map[int32]int64
		limitReached                     map[int32]bool
		offsetManagerManagePartitionrErr error
		offsetManagerCloseErr            error
		expectAnError                    bool
	}
	tests := []test{
		{
			name: "offset is tracked",
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
		},
		{
			name: "error managing partition",
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			offsetManagerManagePartitionrErr: errors.New("error managing partition"),
			expectAnError:                    true,
		},
		{
			name: "error closing offset manager",
			offsets: map[int32]int64{
				0: 0,
				1: 0,
			},
			limitReached: map[int32]bool{
				0: false,
				1: false,
			},
			offsetManagerCloseErr: errors.New("error closing offset manager"),
			expectAnError:         false, // Although it errors, the OffsetCommit doesn't return an error
		},
	}
	err := metrics.InitMetrics("testEnv")
	assert.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var testMockConsumer = mockConsumer{}

			config := reportreader.MockConfiguration{
				Topic:      "test_topic",
				GroupID:    "test_group",
				MaxRecords: 10,
				Aggregator: nil,
				OffsetManager: &mockOffsetManager{
					partitionOffsetManager: &mockPartitionOffsetManager{},
					managePartitionErr:     tc.offsetManagerManagePartitionrErr,
					closeErr:               tc.offsetManagerCloseErr,
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

			err = sut.OffsetCommit()

			if tc.expectAnError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
