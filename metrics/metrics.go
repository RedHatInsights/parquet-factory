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

package metrics

import (
	"github.com/RedHatInsights/insights-operator-utils/metrics/push"
	"github.com/prometheus/client_golang/prometheus"
)

type envInitializer struct {
	environment string
}

var (
	err         error
	tableLabels = []string{
		"table",
	}

	// OffsetMarked number of messages which offset has been marked.
	OffsetMarked prometheus.Gauge
	// OffsetConsummed number of messages consumed.
	OffsetConsummed prometheus.Gauge
	// OffsetProcessed number of messages processed.
	OffsetProcessed prometheus.Gauge
	// FilesGenerated number of files generated, partitioned by table.
	FilesGenerated *prometheus.CounterVec
	// InsertedRows number of rows written, partitioned by table.
	InsertedRows *prometheus.CounterVec
	// ErrorCount is a metric that saves the number of errors
	ErrorCount prometheus.Counter
	// State stores the state of the cronjob job
	State prometheus.Gauge
	// Idle while the job is not running
	Idle float64
	// Init while the job is starting
	Init float64 = 1
	// ConnectToKafka before start consumming
	ConnectToKafka float64 = 2
	// Consume while reading Kafka messages
	Consume float64 = 3
	// GenerateTables while generating Kafka tables
	GenerateTables float64 = 4
)

func (envInit envInitializer) getOffsetMarked() (prometheus.Collector, error) {
	OffsetMarked, err = push.NewGaugeWithError(prometheus.GaugeOpts{
		Name:        "offset_marked",
		Help:        "number of messages which offset has been marked",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	})

	return OffsetMarked, err
}

func (envInit envInitializer) getOffsetConsummed() (prometheus.Collector, error) {
	OffsetConsummed, err = push.NewGaugeWithError(prometheus.GaugeOpts{
		Name:        "offset_consummed",
		Help:        "number of messages consumed",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	})

	return OffsetConsummed, err
}

func (envInit envInitializer) getOffsetProcessed() (prometheus.Collector, error) {
	OffsetProcessed, err = push.NewGaugeWithError(prometheus.GaugeOpts{
		Name:        "offset_processed",
		Help:        "number of messages processed",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	})

	return OffsetProcessed, err
}

func (envInit envInitializer) getFilesGenerated() (prometheus.Collector, error) {
	FilesGenerated, err = push.NewCounterVecWithError(prometheus.CounterOpts{
		Name:        "files_generated",
		Help:        "number of files generated",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	}, tableLabels)

	return FilesGenerated, err
}

func (envInit envInitializer) getInsertedRows() (prometheus.Collector, error) {
	InsertedRows, err = push.NewCounterVecWithError(prometheus.CounterOpts{
		Name:        "inserted_rows",
		Help:        "number of rows written",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	}, tableLabels)

	return InsertedRows, err
}

func (envInit envInitializer) getErrorCount() (prometheus.Collector, error) {
	ErrorCount, err = push.NewCounterWithError(prometheus.CounterOpts{
		Name:        "error_count",
		Help:        "saves the number of errors",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	})

	return ErrorCount, err
}

func (envInit envInitializer) getState() (prometheus.Collector, error) {
	State, err = push.NewGaugeWithError(prometheus.GaugeOpts{
		Name:        "state",
		Help:        "state of the cronjob job",
		ConstLabels: prometheus.Labels{"environment": envInit.environment},
	})

	return State, err
}

// WithTableLabel returns the prometheus label for that table metric
func WithTableLabel(table string) prometheus.Labels {
	return prometheus.Labels{"table": table}
}

// InitMetrics fills the collector variables with some Prometheus metrics and automatically registers them.
func InitMetrics(environment string) error {
	// set the environment
	envInit := envInitializer{
		environment: environment,
	}
	initFunctions := []func() (prometheus.Collector, error){
		envInit.getOffsetMarked,
		envInit.getOffsetConsummed,
		envInit.getOffsetProcessed,
		envInit.getFilesGenerated,
		envInit.getInsertedRows,
		envInit.getErrorCount,
		envInit.getState,
	}
	return push.InitMetrics(initFunctions)
}
