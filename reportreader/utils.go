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
	"encoding/json"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/reportaggregators/rulereportaggregator"
)

const (
	cannotParseMsg = "could not be parsed"
	contentTag     = "content"
)

// consumerLog facilitates logging
func consumerLog(event *zerolog.Event, m *sarama.ConsumerMessage, logMsg string) {
	event = event.
		Int32(partitionTag, m.Partition).
		Int64(offsetTag, m.Offset).
		Str(topicTag, m.Topic).
		Time("timestamp", m.Timestamp).
		Interface("headers", headersToStrings(m.Headers))

	config := conf.GetConfiguration()

	switch m.Topic {
	case config.RulesKafkaConsumer.Topic:
		var parsed rulereportaggregator.RulesResultsReport
		if err := json.Unmarshal(m.Value, &parsed); err != nil {
			event = event.
				Str(archiveTag, cannotParseMsg).
				Str(contentTag, string(m.Value))
		} else {
			event = event.Str(archiveTag, parsed.Path)
		}
	}
	event.Msg(logMsg)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	if timeout == 0*time.Second {
		wg.Wait()
		return false
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
