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
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestLimitCheckerCheckMessage(t *testing.T) {
	var (
		limitDate = time.Now()
		message   = sarama.ConsumerMessage{}
	)

	t.Run("when consumedMessages < maxRecords, if a message is previous to current hour it should pass the limit check", func(t *testing.T) {
		sut := limitChecker{
			limitTimestamp:   limitDate,
			maxRecords:       1,
			consumedMessages: 0,
			mutex:            sync.RWMutex{},
		}
		message.Timestamp = limitDate.Add(-1 * time.Hour)
		assert.True(t, sut.CheckMessage(&message), "the message should pass")
	})

	t.Run("when consumedMessages < maxRecords, if a message is posterior to current hour it shouldn't pass the limit check", func(t *testing.T) {
		sut := limitChecker{
			limitTimestamp:   limitDate,
			maxRecords:       1,
			consumedMessages: 0,
			mutex:            sync.RWMutex{},
		}
		message.Timestamp = limitDate.Add(1 * time.Hour)
		assert.False(t, sut.CheckMessage(&message), "the message shouldn't pass")
	})

	t.Run("if consumedMessages > maxRecords, no message should pass", func(t *testing.T) {
		sut := limitChecker{
			limitTimestamp:   limitDate,
			maxRecords:       0,
			consumedMessages: 1,
			mutex:            sync.RWMutex{},
		}
		message.Timestamp = limitDate.Add(-1 * time.Hour)
		assert.False(t, sut.CheckMessage(&message), "the message shouldn't pass")
	})
}

func TestLimitCheckerMessageProcessed(t *testing.T) {
	sut := limitChecker{
		consumedMessages: 0,
		mutex:            sync.RWMutex{},
	}
	sut.MessageProcessed()
	assert.Equal(t, 1, sut.consumedMessages)
}

func TestLimitCheckerCanConsumeMore(t *testing.T) {
	t.Run("when consumedMessages < maxRecords, should return true", func(t *testing.T) {
		sut := limitChecker{
			maxRecords:       10,
			consumedMessages: 5,
			mutex:            sync.RWMutex{},
		}
		assert.True(t, sut.CanConsumeMore(), "should be able to consume more messages")
	})

	t.Run("when consumedMessages == maxRecords, should return false", func(t *testing.T) {
		sut := limitChecker{
			maxRecords:       10,
			consumedMessages: 10,
			mutex:            sync.RWMutex{},
		}
		assert.False(t, sut.CanConsumeMore(), "should not be able to consume more messages")
	})

	t.Run("when consumedMessages > maxRecords, should return false", func(t *testing.T) {
		sut := limitChecker{
			maxRecords:       10,
			consumedMessages: 15,
			mutex:            sync.RWMutex{},
		}
		assert.False(t, sut.CanConsumeMore(), "should not be able to consume more messages")
	})
}
