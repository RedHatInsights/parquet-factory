---
layout: page
title: Configuration
nav_order: 2
---

# Configuration
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

Parquet Factory is configured using a **toml** file. By default, it will load
the `config.toml` file in the working directory, but any other file can be used
by defining the environment variable `PARQUET_FACTORY__CONFIG_FILE`.

Also each key in config can be overwritten by corresponding env var. For example
if you have config

```toml
[kafka_rules]
address = "kafka:9092"
security_protocol = "PLAINTEXT"
sasl_mechanism = "PLAIN"
client_id = "username"
client_secret = "password"
topic = "incoming_rules_topic"
group_id = "parquet-factory-group"
cert_path = "somewhere"
max_consumed_records = 5000
```

You can define environment variables like
`PARQUET_FACTORY__KAFKA_RULES__ADDRESS` to overwrite the `address` inside the
`kafka_rules` section.

It's very powerful to avoid storing sensitive information (like passwords)
inside the configuration file.

## Rule hits consumer configuration

The Kafka consumer for rule hits topic is configured in the section
`[kafka_rules]`:

```toml
[kafka_rule]
address = "kafka:9092"
security_protocol = "PLAINTEXT"
sasl_mechanism = "PLAIN"
client_id = "username"
client_secret = "password"
topic = "incoming_rules_topic"
group_id = "parquet-factory-group"
cert_path = "somewhere"
max_consumed_records = 5000
consumer_timeout = 240
```

* `address` is the host and port to the Kafka broker to be used.
* `security_protocol` is the `security.protocol` configuration property used by
  the Kafka consumer. Currently, `PLAINTEXT`, `SSL` and `SASL_SSL` are supported
  and tested.
* `sasl_mechanism`: only used when `security_protocol` is set to `SASL_SSL`. It
  corresponds with `sasl.mechanisms` Kafka property. Currently, only `PLAIN` is
  tested.
* `client_id`: used with `SASL_SSL` authentication, it corresponds with `sasl.username`
  Kafka property.
* `client_secret`: used with `SASL_SSL` authentication, it corresponds with `sasl.password`
  Kafka property.
* `topic` is the topic name to consume messages from.
* `group_id` is the consumer group identifier to be used in this topic.
* `cert_path` is a path in the file system to a certificate to be used to
  connect to the Kafka broker.
* `max_consumed_records` is an integer representing the maximum number of Kafka
  records that `parquet-factory` is able to read from the rule hits topic in a
  single execution.
* `max_retries` is an integer indicating the maximum number of retries that the
  consumer will try before exiting.
* `consumer_timeout` timeout in seconds that PF rules consumer will wait for all partitions to finish.

## Features extraction consumer configuration

It is very similar to the previous one, the only change is the section name:
`[kafka_features]`:

```toml
[kafka_features]
address = "kafka:9092"
topic = "incoming_features_topic"
group_id = "parquet-factory-group"
cert_path = "somewhere"
max_consumed_records = 5000
max_retries = 3
consumer_timeout = 240
```

All the properties inside the section have the same meaning as the ones
described in the previous section of this document.

## S3 configuration

The generated parquet files will be stored directly into a S3 instance. To
configure the access to this network storage, the `[s3]` section is used:

```toml
[s3]
endpoint = "localhost:9000"
bucket = "ceph"
prefix = "pipeline_data"
region = "us-east-1"
access_key = "minio"
secret_key = "minio123"
use_ssl = false
```

* `endpoint` is the address used to access the S3 storage, in the form of a pair
  of hostname and port.
* `bucket` indicates the storage bucket where the Parquet files will be
  uploaded.
* `prefix` indicates the first part of the generated file names in the S3
  storage. All the generated Parquet files will share the same prefix.
* `region` where the storage is located. Some endpoints will ignore this.
* `access_key` and `secret_key` are a pair of string credentials used to be
  authenticated by the S3 server.
* `use_ssl` indicates whether use SSL to connect to the S3 instance or not.

## Logging configuration

The logging configuration is made according to the
[Insights Operator Utils](https://pkg.go.dev/github.com/RedHatInsights/insights-operator-utils)
package. Parquet Factory is able to configure `zerolog` package using this framework that
allows to easily configure the logging to CloudWatch, Sentry or Kafka topic.

### General logging configuration

```toml
[logging]
debug = false
use_stderr = false
log_level = "info"
logging_to_cloud_watch_enabled = false
logging_to_sentry_enabled = false
logging_to_kafka_enabled = false
```

In the configuration above you can see all the possible configuration parameters and
their default values.

* `debug` when enabled, this parameter enables the usage of [`zerolog.ConsoleWriter`](https://pkg.go.dev/github.com/rs/zerolog?utm_source=godoc#readme-pretty-logging)
  which perform a pretty and coloured printing of the logs to the console.
* `use_stderr`: by default, the standard output is used to print the logs. If you need to print them to
  the standard error output, you can set this to `true`.
* `log_level` indicated the minimun level that will be printed
* `logging_to_cloud_watch_enabled` send the logs to a configured CloudWatch instance (read bellow).
* `logging_to_sentry_enabled` as the previous one, but sending the error logs to Sentry.
* `logging_to_kafka_enabled` send the logs in JSON format to a Kafka topic.

### Logging to different cloud services

As stated in the previous section, the logs can be sent to different cloud services that can handle
or store the logs for monitoring or debug purposes. The currently supported services are CloudWatch,
Sentry and Kafka.

Each one has its own configuration parameters. For further info and insights about it, take a look to
[logger package documentation](https://pkg.go.dev/github.com/RedHatInsights/insights-operator-utils/logger).
