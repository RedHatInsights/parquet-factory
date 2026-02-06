#!/usr/bin/env python3
# Copyright 2021 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
This script takes a file with a set of lines in the form `timestamp partition-offset message`,
orders the lines by the timestamp and publish the records with the same intervals, taking first
message as base timestamp.
"""
import argparse
import asyncio
import datetime
import time

import aiokafka


async def main():
    """Handle the execution of the script."""
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "records_file", help="Path to a file containing the records to be injected",
    )
    parser.add_argument(
        "kafka_server", help="Kafka server to inject messages",
    )
    parser.add_argument(
        "topic", help="Kafka topic to inject the messages",
    )
    parser.add_argument(
        "--update_timestamp", action="store_true",
        help="Moves the timestamp forward to avoid the retention policy"
    )
    parser.add_argument(
        "--ignore-partition", action="store_true",
        help="Ignore the partition field in the input file and publish to any available partition"
    )

    args = parser.parse_args()

    print(f"Opening {args.records_file} for reading")
    with open(args.records_file) as records_file:
        records = records_file.readlines()

    print("Done. Splitting by columns...")
    records = [record.split(" ", 3) for record in records]
    print("Done. Sorting by timestamp...")
    records.sort(key=lambda x: x[0])  # sorting by first column, timestamp

    shift = 0
    if args.update_timestamp:
        shift = (time.time() * 1000) - 601200000. - int(records[0][0])
        print(f"Original ts: {records[0][0]}\nNew ts: {int(records[0][0]) + shift}\nShift: {shift}")

    producer = aiokafka.AIOKafkaProducer(bootstrap_servers=[args.kafka_server])
    await producer.start()
    print(f"Start sending {len(records)} messages")

    try:
        for record in records:
            timestamp = int(record[0]) + shift
            offset = int(record[1])
            partition= int(record[2]) if not args.ignore_partition else None
            payload = record[3][:-1].encode()  # Remove trailing \n
            await producer.send(args.topic, value=payload, timestamp_ms=timestamp, partition=partition)

    finally:
        await producer.stop()



if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
