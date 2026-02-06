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

package reportaggregators_test

import (
	"testing"
	"time"

	"github.com/RedHatInsights/parquet-factory/reportaggregators"
	"github.com/stretchr/testify/assert"
)

// TestExtractCollectedDate checks several cases of ExtractCollectDate
func TestExtractCollectedDate(t *testing.T) {
	var mapTests = []struct {
		name  string
		path  string
		date  string
		valid bool
	}{
		{
			"Valid 1",
			"archives/compressed/60/60d2f3b5-7dee-4044-929a-62ebc6bf6188/202105/05/200331.tar.gz",
			"2021-05-05T20:03:31+00:00",
			true,
		}, {
			"Valid 2",
			"archives/compressed/60/60d2f3b5-7dee-4044-929a-62ebc6bf6188/202104/12/123456.tar.gz",
			"2021-04-12T12:34:56+00:00",
			true,
		}, {
			"Invalid path",
			"archives/compressed/60/60d2f3b5-7dee-4044-929a-62ebc6bf6188/123456.tar.gz",
			"",
			false,
		}, {
			"Invalid hour",
			"archives/compressed/60/60d2f3b5-7dee-4044-929a-62ebc6bf6188/202104/30/293456.tar.gz",
			"",
			false,
		},
	}

	for _, tt := range mapTests {
		t.Run(tt.name, func(t *testing.T) {
			extracted, err := reportaggregators.ExtractCollectedDate(tt.path)
			if tt.valid {
				expected, _ := time.Parse(time.RFC3339, tt.date)

				assert.Equal(t, expected.Unix(), extracted.Unix())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
