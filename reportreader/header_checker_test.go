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
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestCheckHeaders(t *testing.T) {
	var (
		aHeader = sarama.RecordHeader{
			Key:   []byte("a"),
			Value: []byte("header"),
		}
		newStopHeader = sarama.RecordHeader{
			Key:   []byte("stop"),
			Value: []byte("true"),
		}
		normalMsg = sarama.ConsumerMessage{
			Headers: []*sarama.RecordHeader{
				&aHeader,
			},
		}
		stopMsg = sarama.ConsumerMessage{
			Headers: []*sarama.RecordHeader{
				&newStopHeader,
			},
		}
	)
	t.Run("a normal header", func(t *testing.T) {
		ans := checkHeaders(&normalMsg)
		assert.False(t, ans)
	})

	t.Run("a stop header", func(t *testing.T) {
		ans := checkHeaders(&stopMsg)
		assert.True(t, ans)
	})
}

func TestHeadersToString(t *testing.T) {
	var headers = []*sarama.RecordHeader{
		{
			Key:   []byte("a"),
			Value: []byte("header"),
		},
	}

	got := headersToStrings(headers)
	want := map[string]string{
		"a": "header",
	}
	assert.Equal(t, want, got)
}
