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

package rulereportaggregator

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/RedHatInsights/parquet-factory/metrics"
	"github.com/RedHatInsights/parquet-factory/reportaggregators"
	"github.com/RedHatInsights/parquet-factory/s3writer"
)

// RuleHit represents each RuleHit report received
type RuleHit struct {
	RuleID string `json:"rule_id"`
}

// InfoReport represents each Info report received
type InfoReport struct {
	InfoID    string                 `json:"info_id"`
	Component string                 `json:"component"`
	Details   map[string]interface{} `json:"details"`
	Key       string                 `json:"key"`
}

// RuleReport represents the "report" key of the received reports"
type RuleReport struct {
	Info    []InfoReport `json:"info"`
	Reports []RuleHit    `json:"reports"`
}

// RulesResultsReport represents the whole received report with the needed keys
type RulesResultsReport struct {
	Path     string                     `json:"path"`
	Metadata reportaggregators.Metadata `json:"metadata"`
	Report   RuleReport                 `json:"report"`
}

// RulesResultsReportAggregator stores an array of RulesResultsReport
type RulesResultsReportAggregator struct {
	ReceivedReports []RulesResultsReport
	mutex           sync.RWMutex
}

// NewRulesReportAggregator initialize a RulesResultsReportAggregator variable
func NewRulesReportAggregator() *RulesResultsReportAggregator {
	return &RulesResultsReportAggregator{
		ReceivedReports: []RulesResultsReport{},
	}
}

// Handle parses an incoming message from Kafka and store it in the aggregation
func (aggregator *RulesResultsReportAggregator) Handle(data interface{}) error {
	message, ok := data.([]byte)
	if !ok {
		return errors.New("the argument doesn't match the expected type")
	}

	var parsed RulesResultsReport
	if err := json.Unmarshal(message, &parsed); err != nil {
		log.Info().Err(err).Msgf("Unable to parse message: %v", message)
		log.Error().Err(err).Msg("Unable to parse message")
		return err
	}

	aggregator.mutex.Lock()
	defer aggregator.mutex.Unlock()
	aggregator.ReceivedReports = append(aggregator.ReceivedReports, parsed)
	return nil
}

// WriteResults writes the aggregated results into the  provided S3ParquetWriter
func (aggregator *RulesResultsReportAggregator) WriteResults(writer s3writer.S3ParquetWriter) (int, error) {
	metrics.State.Set(metrics.GenerateTables)
	var writtenResults int

	ruleHitsFiles, err := aggregator.createRuleHitTable(writer)
	if err != nil {
		log.Error().Err(err).Msg("error saving rule hit tables")
		deleteErr := writer.DeleteFiles(ruleHitsFiles)
		if deleteErr != nil {
			log.Error().Err(deleteErr).Msg("error deleting incomplete rule report!")
		}
		return 0, err
	}

	writtenResults += len(ruleHitsFiles)

	archivesFiles, err := aggregator.createArchivesTable(writer)
	if err != nil {
		log.Error().Err(err).Msg("error saving archives tables")
		deleteErr := writer.DeleteFiles(archivesFiles)
		if deleteErr != nil {
			log.Error().Err(deleteErr).Msg("error deleting incomplete rule report!")
		}
		return 0, err
	}

	writtenResults += len(archivesFiles)

	return writtenResults, nil
}
