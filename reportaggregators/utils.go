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

package reportaggregators

import (
	"fmt"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/RedHatInsights/parquet-factory/metrics"
)

var timestampRe = regexp.MustCompile(`archives/compressed/[0-9a-f]+/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/(20[0-9]{2})([0-1][0-9])/(0[1-9]|[12]\d|3[01])/([0-2][0-9])([0-6][0-9])([0-6][0-9])\.tar\.gz$`)

// Metadata represents the key "metadata" of the received reports
type Metadata struct {
	ClusterID string `json:"cluster_id"`
}

// LogInsertedRow is a helper to print an appropriate log when a row is inserted
func LogInsertedRow(archivePath, table string) {
	metrics.InsertedRows.With(metrics.WithTableLabel(table)).Inc()
	log.Trace().Str("archive_path", archivePath).Str("table", table).Msg("inserted row")
}

// ExtractCollectedDate allow to extract the collected timestamp from the path to an archive
func ExtractCollectedDate(path string) (time.Time, error) {
	dateMatch := timestampRe.FindAllSubmatch([]byte(path), -1)
	if dateMatch == nil {
		return time.Time{}, fmt.Errorf("unable to parse archive path")
	}

	parsedDate, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", dateMatch[0][1], dateMatch[0][2], dateMatch[0][3], dateMatch[0][4], dateMatch[0][5], dateMatch[0][6]))
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}
