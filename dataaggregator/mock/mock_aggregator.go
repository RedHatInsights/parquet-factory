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

package mock

import (
	"errors"

	"github.com/RedHatInsights/parquet-factory/s3writer"
)

// Aggregator is a DataAggregator to be used when no checks are needed
type Aggregator struct{}

// Handle does nothing, it is just a fake
func (a *Aggregator) Handle(interface{}) error {
	return nil
}

// WriteResults does nothing, it is just a fake
func (a *Aggregator) WriteResults(s3writer.S3ParquetWriter) (int, error) {
	return 0, nil
}

// FaultyAggregator return an error on WriteResults
type FaultyAggregator struct{}

// Handle does nothing, it is just a fake
func (a *FaultyAggregator) Handle(interface{}) error {
	return nil
}

// WriteResults simulates a writing failure
func (a *FaultyAggregator) WriteResults(s3writer.S3ParquetWriter) (int, error) {
	return 0, errors.New("test error")
}
