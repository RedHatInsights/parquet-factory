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

package conf_test

import (
	"os"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/logger"
	"github.com/RedHatInsights/insights-operator-utils/types"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/parquet-factory/conf"
)

func TestLoadConfiguration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
}

func TestLoadConfigurationNotExist(t *testing.T) {
	os.Clearenv()
	err := conf.LoadConfiguration("idontbelivethisexist")

	if err == nil {
		t.Fail()
	}
}

func TestGetConfiguration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
	cfg := conf.GetConfiguration()

	assert.Equal(t, []string{"kafka:9092"}, cfg.RulesKafkaConsumer.Addresses)
	assert.Equal(t, 240, cfg.RulesKafkaConsumer.ConsumerTimeout)
	assert.Equal(t, 30, cfg.TimeShift)
}

func TestCloudWatchConfiguration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
	cfg := conf.GetCloudWatchConfiguration()

	assert.Equal(t, logger.CloudWatchConfiguration{}, cfg)
}

func TestGetSentryConfiguration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
	cfg := conf.GetSentryConfiguration()

	assert.Equal(
		t,
		logger.SentryLoggingConfiguration{
			SentryDSN:         "put_your_url_here",
			SentryEnvironment: "put_your_env_here",
		},
		cfg,
	)
}

func TestLoggingConfiguration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
	cfg := conf.GetLoggingConfiguration()

	assert.Equal(
		t,
		logger.LoggingConfiguration{
			Debug:                      true,
			LoggingToCloudWatchEnabled: false,
			LoggingToSentryEnabled:     false,
		},
		cfg,
	)
}

func TestGetS3Configuration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
	cfg := conf.GetS3Configuration()

	assert.Equal(
		t,
		conf.S3Config{
			Endpoint:       "localhost:9000",
			Bucket:         "ceph",
			FilePathPrefix: "fleet_aggregations",
			Region:         "us-east-1",
			AccessKey:      "minio",
			SecretKey:      "minio123",
		},
		cfg,
	)
}

func TestGetMetricsConfiguration(t *testing.T) {
	os.Clearenv()
	mustLoadConfiguration(t, "../testdata/config1")
	cfg := conf.GetMetricsConfiguration()

	assert.Equal(
		t,
		types.MetricsConfiguration{
			Job:              "job_name",
			GatewayURL:       "gateway_url",
			GatewayAuthToken: "gateway_auth_token",
			TimeBetweenPush:  60,
		},
		cfg,
	)
}

type testCase struct {
	name           string
	clowderFile    string
	expectedRegion string
	expectedUseSSL bool
}

func TestCheckUpdateClowder(t *testing.T) {
	testCases := []testCase{
		{
			name:           "clowder config complete",
			clowderFile:    "../testdata/clowderconfig.json",
			expectedRegion: "testRegion",
			expectedUseSSL: true,
		},
		{
			name:           "clowder config with bucket as in ephemeral",
			clowderFile:    "../testdata/clowderconfig-ephemeralbucket.json",
			expectedRegion: "us-east-1",
			expectedUseSSL: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()
			err := os.Setenv("ACG_CONFIG", tc.clowderFile)
			assert.Nil(t, err)
			clowder.LoadedConfig, err = clowder.LoadConfig(tc.clowderFile)

			// FIXME: this is a workaround because clowder doesn't load the config on start up
			// it's a copy of the clowder.init() function
			clowder.ObjectBuckets = make(map[string]clowder.ObjectStoreBucket)
			if clowder.LoadedConfig.ObjectStore != nil {
				for _, bucket := range clowder.LoadedConfig.ObjectStore.Buckets {
					clowder.ObjectBuckets[bucket.RequestedName] = bucket
				}
			}

			assert.Nil(t, err)
			mustLoadConfiguration(t, "../testdata/config1")
			cfg := conf.GetConfiguration()

			assert.Equal(
				t,
				[]string{"clowderkafkabroker:9092"},
				cfg.RulesKafkaConsumer.Addresses,
			)

			// TODO: cannot be tested because the ACG_CONFIG envvar is not loaded when
			// clowder package is loaded and init is executed
			// assert.Equal(
			// 	t,
			// 	"translated_rules_topic",
			// 	cfg.RulesKafkaConsumer.Topic,
			// )

			assert.Equal(t, "http://clowders3server:9000", cfg.S3.Endpoint)
			assert.Equal(t, tc.expectedRegion, cfg.S3.Region)
			assert.Equal(t, "testAccessKey", cfg.S3.AccessKey)
			assert.Equal(t, "testSecretKey", cfg.S3.SecretKey)
			assert.Equal(t, tc.expectedUseSSL, cfg.S3.UseSSL)
		})
	}
}
