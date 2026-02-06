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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tisnik/go-capture"

	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/reportaggregators/rulereportaggregator"
)

var (
	testTopic  = "my_topic"
	testPath   = "my_path"
	ruleReport = rulereportaggregator.RulesResultsReport{
		Path: testPath,
		Report: rulereportaggregator.RuleReport{
			Info: []rulereportaggregator.InfoReport{
				{
					InfoID: "an id",
				},
			},
		},
	}
	logMessage = "a simple log"
)

func TestConsumerLog(t *testing.T) {
	err := conf.LoadConfiguration("../testdata/config1.toml")
	require.NoError(t, err, "cannot load configuration from file")

	config := conf.GetConfiguration()

	fmt.Println()

	type testCase struct {
		name           string
		value          []byte
		topic          string
		expectedString string
	}

	testCases := []testCase{
		{
			name:  "unexpected topic",
			topic: testTopic,
		},
		{
			name:           "parsable rules message in rules topic",
			topic:          config.RulesKafkaConsumer.Topic,
			value:          structToBytes(ruleReport),
			expectedString: ruleReport.Path,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := sarama.ConsumerMessage{
				Key:       []byte("a key"),
				Value:     tc.value,
				Topic:     tc.topic,
				Partition: 1,
				Offset:    1,
			}
			output, err := capture.ErrorOutput(func() {
				log.Logger = log.Output(zerolog.New(os.Stderr))
				consumerLog(log.Debug(), &msg, logMessage)
			})
			require.NoError(t, err, "unable to capture standard output")

			assert.Contains(t, output, logMessage)
			assert.Contains(t, output, tc.topic)
			assert.Contains(t, output, tc.expectedString)
		})
	}
}

func TestWaitTimeout(t *testing.T) {
	t.Run("timed out", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)

		timeout := time.Second
		assert.True(t, waitTimeout(&wg, timeout))
	})
	t.Run("no timed out", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Done()

		timeout := time.Second
		assert.False(t, waitTimeout(&wg, timeout))
	})
}

func structToBytes(testStruct interface{}) []byte {
	bytesBuffer := new(bytes.Buffer)
	_ = json.NewEncoder(bytesBuffer).Encode(testStruct)

	return bytesBuffer.Bytes()
}
