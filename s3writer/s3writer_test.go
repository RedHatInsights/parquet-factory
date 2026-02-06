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

package s3writer_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/s3writer"
)

var s3TestConf = conf.S3Config{
	Endpoint:       "testEndpoint",
	Bucket:         "testBucket",
	FilePathPrefix: "fleet_data",
	Region:         "EU",
	AccessKey:      "testAccessKey",
	SecretKey:      "testSecretKey",
	UseSSL:         true,
}

func TestNew(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		value, err := s3writer.New(s3TestConf)
		assert.NoError(t, err)
		assert.NotNil(t, value)
	})

	t.Run("invalid configuration", func(t *testing.T) {
		err := os.Setenv("AWS_STS_REGIONAL_ENDPOINTS", "invalid-value")
		assert.NoError(t, err)
		defer func() {
			err := os.Setenv("AWS_STS_REGIONAL_ENDPOINTS", "")
			assert.NoError(t, err)
		}()

		value, err := s3writer.New(s3TestConf)
		assert.NoError(t, err) // AWS SDK v2 is more lenient than v1 and doesn't fail on AWS_STS_REGIONAL_ENDPOINTS invalid values when using explicit credentials
		assert.NotNil(t, value)
	})
}
