---
layout: page
title: Deployment
nav_order: 4
---

## Online deployment

`parquet-factory` is deployed as a
[Cronjob](https://docs.openshift.com/container-platform/4.6/nodes/jobs/nodes-nodes-jobs.html)
in the PSI cluster.
It is configured to be run every hour.

## Local deployment

If you intend to work on `parquet-factory` locally, you can use the `docker-compose.yaml` configuration
in order to create a deployment that can be easily used by your local version.

### Usage

#. Start the pods with `podman-compose up` or `docker-compose up`.
#. Adapt the `config.toml` to your local deployment, if needed. The provided one
   contains the needed configuration in order to use the services deployed locally by
   docker|podman compose.

If you intend to test something related to Kafka and you want to use more than 1 partition per topic, you will need to configure the topics using the next command:

```
kafka-topics --bootstrap-server kafka:9092 --create --topic incoming_features_topic --partitions 4
```

The command `kafka-topic` can be found in the default `cp-kafka` container image used in the previous steps, so you can run it:

```
podman run -it --rm --network=host -v /etc/resolv.conf:/etc/resolv.conf:Z cp-kafka kafka-topics --bootstrap-server kafka:9092 --create --topic incoming_features_topic --partitions 4
```

### Generating Kafka messages

In order to help to reproduce the normal behaviour of the `parquet-factory`, an utility named `load_kafka.py` was written. This program read a set of records from a file and send them to a Kafka topic.

The file should have an special format:

```
TIMESTAMP OFFSET PARTITION MESSAGE
```

Using this tool, timestamp and partition will be preserved for every message sent to the topic, which will help to ensure that the behaviour is closer to the real one.

In order to generate the input file needed by this script, you can use `kafkacat`:

`kafkacat -b kafka_broker:443 -C -o beginning -e -f "%T %o %p %s\n" -t TOPIC > output_file.txt`
