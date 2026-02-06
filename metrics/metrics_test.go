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

package metrics_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/RedHatInsights/parquet-factory/metrics"
)

var testEnv = "testEnv"

func TestInitMetrics(t *testing.T) {
	t.Run("optimal setup", func(t *testing.T) {
		err := metrics.InitMetrics(testEnv)
		assert.NoError(t, err)
	})
	t.Run("already registered metrics", func(t *testing.T) {
		err := metrics.InitMetrics(testEnv)
		assert.NoError(t, err)
	})
}

func TestWithTableLabel(t *testing.T) {
	var testTable = "my_table"
	assert.Equal(t,
		metrics.WithTableLabel(testTable),
		prometheus.Labels{"table": testTable})
}
