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

package reportreader_test

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/RedHatInsights/parquet-factory/utils"
)

type mockPartitionOffsetManager struct {
	nextOffset int64
	metadata   string
	err        error
	errChan    chan *sarama.ConsumerError
}

func (p *mockPartitionOffsetManager) NextOffset() (int64, string) {
	fmt.Println("Calling NextOffset")
	return p.nextOffset, p.metadata
}

func (p *mockPartitionOffsetManager) MarkOffset(offset int64, metadata string) {
	fmt.Println("Calling MarkOffset")
}

func (p *mockPartitionOffsetManager) ResetOffset(offset int64, metadata string) {
	fmt.Println("Calling ResetOffset")
}

func (p *mockPartitionOffsetManager) Errors() <-chan *sarama.ConsumerError {
	fmt.Println("Calling Errors")
	return p.errChan
}

func (p *mockPartitionOffsetManager) AsyncClose() {
	fmt.Println("Calling AsyncClose")
}

func (p *mockPartitionOffsetManager) Close() error {
	fmt.Println("Calling Close")
	return p.err
}

type mockOffsetManager struct {
	partitionOffsetManager sarama.PartitionOffsetManager
	managePartitionErr     error
	closeErr               error
}

func (om *mockOffsetManager) ManagePartition(topic string, partition int32) (sarama.PartitionOffsetManager, error) {
	return om.partitionOffsetManager, om.managePartitionErr
}

func (om *mockOffsetManager) Close() error {
	return om.closeErr
}

func (om *mockOffsetManager) Commit() {}

type mockPartitionConsumer struct {
	errClose            error
	messageChan         chan *sarama.ConsumerMessage
	errChan             chan *sarama.ConsumerError
	highWaterMarkOffset int64
}

func (pc mockPartitionConsumer) AsyncClose() {
	fmt.Println("Calling AsyncClose")
}

func (pc mockPartitionConsumer) Close() error {
	fmt.Println("Calling Close")
	return pc.errClose
}

// the broker.
func (pc mockPartitionConsumer) Messages() <-chan *sarama.ConsumerMessage {
	fmt.Println("Calling Messages")
	return pc.messageChan
}

func (pc mockPartitionConsumer) Errors() <-chan *sarama.ConsumerError {
	fmt.Println("Calling Errors")
	return pc.errChan
}

func (pc mockPartitionConsumer) HighWaterMarkOffset() int64 {
	fmt.Println("Calling HighWaterMarkOffset")
	return pc.highWaterMarkOffset
}

func (pc mockPartitionConsumer) IsPaused() bool {
	fmt.Println("Calling IsPaused")
	return false
}

func (pc mockPartitionConsumer) Pause() {
	fmt.Println("Calling Pause on PartitionConsumer")
}

func (pc mockPartitionConsumer) Resume() {
	fmt.Println("Calling Resume on PartitionConsumer")
}

type mockConsumer struct {
	topics              []string
	partitions          []int32
	partitionConsumer   sarama.PartitionConsumer
	errTopics           error
	errPartitions       error
	errConsumePartition error
	errClose            error
	highWaterMarks      map[string]map[int32]int64
	lastCall            string
	processedMessages   *utils.ArchivePathSet
}

func (c *mockConsumer) Topics() ([]string, error) {
	fmt.Println("Calling Topics")
	c.lastCall = "Topics"
	return c.topics, c.errTopics
}

func (c *mockConsumer) Partitions(topic string) ([]int32, error) {
	fmt.Println("Calling Partitions")
	c.lastCall = "Partitions"
	return c.partitions, c.errPartitions
}

func (c *mockConsumer) ConsumePartition(topic string, partition int32, offset int64) (sarama.
	PartitionConsumer, error) {
	fmt.Println("Calling ConsumePartition")
	c.lastCall = "ConsumePartition"
	return c.partitionConsumer, c.errConsumePartition
}

func (c *mockConsumer) HighWaterMarks() map[string]map[int32]int64 {
	fmt.Println("Calling HighWaterMarks")
	c.lastCall = "HighWaterMarks"
	return c.highWaterMarks
}

func (c *mockConsumer) Close() error {
	fmt.Println("Calling Close")
	c.lastCall = "Close"
	return c.errClose
}

func (c *mockConsumer) Pause(topicPartitions map[string][]int32) {
	fmt.Println("Calling Pause")
	c.lastCall = "Pause"
}

func (c *mockConsumer) Resume(topicPartitions map[string][]int32) {
	fmt.Println("Calling Resume")
	c.lastCall = "Resume"
}

func (c *mockConsumer) PauseAll() {
	fmt.Println("Calling PauseAll")
	c.lastCall = "PauseAll"
}

func (c *mockConsumer) ResumeAll() {
	fmt.Println("Calling ResumeAll")
	c.lastCall = "ResumeAll"
}

func NewMockConsumer() mockConsumer {
	var testMockConsumer = mockConsumer{}
	testMockConsumer.processedMessages = utils.NewArchivePathSet()
	return testMockConsumer
}
