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

package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RedHatInsights/insights-operator-utils/logger"
	"github.com/RedHatInsights/insights-operator-utils/types"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
)

const (
	noSaslConfig              = "warning: SASL configuration is missing"
	noBrokerConfig            = "warning: no broker configurations found in clowder config"
	configFileEnvVariableName = "PARQUET_FACTORY_CONFIG_FILE"
	envPrefix                 = "PARQUET_FACTORY_"
)

// KafkaConfig represents the configuration for the Kafka consumer
type KafkaConfig struct {
	Addresses        []string `mapstructure:"address" toml:"address"`
	SecurityProtocol string   `mapstructure:"security_protocol" toml:"security_protocol"`
	CertPath         string   `mapstructure:"cert_path" toml:"cert_path"`
	SaslMechanism    string   `mapstructure:"sasl_mechanism" toml:"sasl_mechanism"`
	ClientID         string   `mapstructure:"client_id" toml:"client_id"`
	ClientSecret     string   `mapstructure:"client_secret" toml:"client_secret"` // #nosec G117 -- Configuration field, not a hardcoded secret
	Topic            string   `mapstructure:"topic" toml:"topic"`
	GroupID          string   `mapstructure:"group_id" toml:"group_id"`
	MaxRecords       int      `mapstructure:"max_consumed_records" toml:"max_consumed_records"`
	MaxRetries       int      `mapstructure:"max_retries" toml:"max_retries"`
	ConsumerTimeout  int      `mapstructure:"consumer_timeout" toml:"consumer_timeout"` // Seconds
}

// S3Config represents the configuration for the S3 client
type S3Config struct {
	Endpoint       string `mapstructure:"endpoint" toml:"endpoint"`
	Bucket         string `mapstructure:"bucket" toml:"bucket"`
	FilePathPrefix string `mapstructure:"prefix" toml:"prefix"`
	Region         string `mapstructure:"region" toml:"region"`
	AccessKey      string `mapstructure:"access_key" toml:"access_key"` // #nosec G117 -- Configuration field, not a hardcoded secret
	SecretKey      string `mapstructure:"secret_key" toml:"secret_key"` // #nosec G117 -- Configuration field, not a hardcoded secret
	UseSSL         bool   `mapstructure:"use_ssl" toml:"use_ssl"`
}

// Config represents the configuration for the parquet-factory
type Config struct {
	RulesKafkaConsumer KafkaConfig                       `mapstructure:"kafka_rules" toml:"kafka_rules"`
	S3                 S3Config                          `mapstructure:"s3" toml:"s3"`
	Logging            logger.LoggingConfiguration       `mapstructure:"logging" toml:"logging"`
	CloudWatch         logger.CloudWatchConfiguration    `mapstructure:"cloudwatch" toml:"cloudwatch"`
	Sentry             logger.SentryLoggingConfiguration `mapstructure:"sentry" toml:"sentry"`
	TimeShift          int                               `mapstructure:"time_shift" toml:"time_shift"` // Minutes
	Metrics            types.MetricsConfiguration        `mapstructure:"metrics" toml:"metrics"`
}

// config holds the loaded configuration
var config Config

// LoadConfiguration loads the configuration from the file pointed for the environment
// variable configFileEndVariableName. If not, it uses configFilePath
func LoadConfiguration(configFilePath string) error {
	configFile, specified := os.LookupEnv(configFileEnvVariableName)

	if !specified {
		configFile = configFilePath
	}

	directory, basename := filepath.Split(configFile)
	file := strings.TrimSuffix(basename, filepath.Ext(basename))
	if directory == "" {
		directory = "."
	}

	viper.SetConfigName(file)
	viper.AddConfigPath(directory)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Something wrong happened parsing configuration")
		fmt.Println(err)
		return err
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Configuration could not be unmarshaled")
		fmt.Println(err)
		return err
	}

	if err := updateConfigFromClowder(&config); err != nil {
		fmt.Println("Error loading clowder configuration")
		return err
	}

	return nil
}

func updateConfigFromClowder(c *Config) error {
	if !clowder.IsClowderEnabled() {
		fmt.Println("Clowder is disabled")
		return nil
	}
	fmt.Println("Clowder is enabled")
	if clowder.LoadedConfig.Kafka == nil {
		fmt.Println("No Kafka configuration available in Clowder, using default one")
	} else {
		updateBrokerCfgFromClowder(c)
	}

	if clowder.LoadedConfig.ObjectStore == nil {
		fmt.Println("No S3 configuration available in Clowder, using default one")
	} else {
		updateBucketCfgFromClowder(c)
	}
	return nil
}

func updateBrokerCfgFromClowder(c *Config) {
	updateTopicMapping(c)
	if len(clowder.LoadedConfig.Kafka.Brokers) == 0 {
		fmt.Println(noBrokerConfig)
		return
	}

	addresses := []string{}
	for _, broker := range clowder.LoadedConfig.Kafka.Brokers {
		if broker.Port != nil {
			addresses = append(addresses, fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port))
		} else {
			addresses = append(addresses, broker.Hostname)
		}
	}
	c.RulesKafkaConsumer.Addresses = addresses

	// SSL config
	clowderBrokerCfg := clowder.LoadedConfig.Kafka.Brokers[0]
	if clowderBrokerCfg.Authtype != nil {
		fmt.Println("kafka is configured to use authentication")
		if clowderBrokerCfg.Sasl != nil {
			c.RulesKafkaConsumer.ClientID = *clowderBrokerCfg.Sasl.Username
			c.RulesKafkaConsumer.ClientSecret = *clowderBrokerCfg.Sasl.Password
			c.RulesKafkaConsumer.SaslMechanism = *clowderBrokerCfg.Sasl.SaslMechanism
			c.RulesKafkaConsumer.SecurityProtocol = *clowderBrokerCfg.SecurityProtocol
			if caPath, err := clowder.LoadedConfig.KafkaCa(clowderBrokerCfg); err == nil {
				c.RulesKafkaConsumer.CertPath = caPath
			}
		} else {
			fmt.Println(noSaslConfig)
		}
	}
}

func updateTopicMapping(c *Config) {
	// Updating topic from clowder mapping if available
	if topicCfg, ok := clowder.KafkaTopics[c.RulesKafkaConsumer.Topic]; ok {
		c.RulesKafkaConsumer.Topic = topicCfg.Name
	} else {
		fmt.Printf("warning: no kafka mapping found for topic %s", c.RulesKafkaConsumer.Topic)
	}
}

func updateBucketCfgFromClowder(c *Config) {
	if bucketCfg, ok := clowder.ObjectBuckets[c.S3.Bucket]; ok {
		c.S3.AccessKey = *bucketCfg.AccessKey
		c.S3.SecretKey = *bucketCfg.SecretKey
		c.S3.Endpoint = fmt.Sprintf("%s:%d", *bucketCfg.Endpoint, clowder.LoadedConfig.ObjectStore.Port)
		if bucketCfg.Tls != nil {
			c.S3.UseSSL = *bucketCfg.Tls
		}
		if bucketCfg.Region != nil {
			c.S3.Region = *bucketCfg.Region
		}

		fmt.Println("[DEBUG] Bucket configuration from Clowder")
		fmt.Printf("    Bucket Access Key: %s\n", c.S3.AccessKey)
		// for the secret key, show first 3 characters and then dots
		fmt.Printf("    Bucket Secret Key: %s...\n", c.S3.SecretKey[:3])
		fmt.Printf("    Bucket Endpoint: %s\n", c.S3.Endpoint)
		fmt.Printf("    Bucket Use SSL: %t\n", c.S3.UseSSL)
		fmt.Printf("    Bucket Region: %s\n", c.S3.Region)
	} else {
		fmt.Printf("warning: no bucket mapping found for bucket %s", c.S3.Bucket)
	}
}

// GetConfiguration returns the loaded configuration
func GetConfiguration() Config {
	return config
}

// GetCloudWatchConfiguration returns CloudWatch configuration
func GetCloudWatchConfiguration() logger.CloudWatchConfiguration {
	return config.CloudWatch
}

// GetSentryConfiguration returns CloudWatch configuration
func GetSentryConfiguration() logger.SentryLoggingConfiguration {
	return config.Sentry
}

// GetLoggingConfiguration returns logging configuration
func GetLoggingConfiguration() logger.LoggingConfiguration {
	return config.Logging
}

// GetS3Configuration returns logging configuration
func GetS3Configuration() S3Config {
	return config.S3
}

// GetMetricsConfiguration returns metrics configuration
func GetMetricsConfiguration() types.MetricsConfiguration {
	return config.Metrics
}
