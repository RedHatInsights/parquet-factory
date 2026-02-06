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

	"github.com/IBM/sarama"
)

var stopHeader = sarama.RecordHeader{
	Key:   []byte("stop"),
	Value: []byte("true"),
}

// checkHeaders returns true if stopHeader is in the message
func checkHeaders(m *sarama.ConsumerMessage) bool {
	return isHeaderInSlice(m.Headers, stopHeader)
}

func isHeaderInSlice(slice []*sarama.RecordHeader, val sarama.RecordHeader) bool {
	for _, item := range slice {
		if bytes.Equal(item.Key, val.Key) && bytes.Equal(item.Value, val.Value) {
			return true
		}
	}
	return false
}

func headersToStrings(slice []*sarama.RecordHeader) map[string]string {
	output := make(map[string]string)
	for _, item := range slice {
		output[string(item.Key)] = string(item.Value)
	}
	return output
}
