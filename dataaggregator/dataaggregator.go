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

package dataaggregator

import "github.com/RedHatInsights/parquet-factory/s3writer"

// DataAggregator defines the interface for every instance able to aggregate
// some specific data
type DataAggregator interface {
	Handle(interface{}) error
	WriteResults(s3writer.S3ParquetWriter) (int, error)
}
