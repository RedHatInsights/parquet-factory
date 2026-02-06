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
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/parquet-factory/metrics"
)

// OffsetCommit commits the offset to every partition using a consumer group handler
func (c *KafkaConsumer) OffsetCommit() error {
	// Iterate over all topic and partitions
	for _, topic := range c.partitionTracker.GetTopics() {
		for _, partition := range c.partitionTracker.GetPartitionsForTopic(topic) {
			// Create the object to handle the partition offset
			partManager, err := c.offsetManager.ManagePartition(topic, partition)
			if err != nil {
				log.Error().Err(err).Msg("Unable to manage partition")
				return err
			}

			// get the offset to be committed
			offset, err := c.partitionTracker.GetOffset(topic, partition)
			if err != nil {
				log.Error().Err(err).Msg("Unable to commit messages")
				return err
			}

			partManager.MarkOffset(offset, "")
			metrics.OffsetMarked.Inc()
			log.Debug().
				Str(topicTag, topic).
				Int32(partitionTag, partition).
				Int64(offsetTag, offset).
				Msg("This stored offset will be committed now")
			partManager.AsyncClose()
		}
	}
	c.offsetManager.Commit()
	if err := c.offsetManager.Close(); err != nil {
		log.Error().Err(err).Msg("Unable to close OffsetManager")
	}
	log.Debug().Msg("All marked offsets have been committed")

	return nil
}
