#!/bin/bash
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


BROKER=${BROKER:-kafka:9092}
OPTIONS=${OPTIONS-""}

TOPICS="ccx-prod-insights-operator-archive-rules-results ccx-prod-insights-operator-archive-features"
PARTITIONS="0 1"

OUTPUT="kafka_input.csv"

KAFKACAT_CMD="kafkacat -b $BROKER $OPTIONS -C "

# Print the CSV file
echo "msg,topic,partition,offset,timestamp" > ${OUTPUT}

for t in ${TOPICS}; do
    for p in ${PARTITIONS}; do
        ${KAFKACAT_CMD} -t ${t} -p ${p} -o beginning -c 1 -f "beginning,%t,%p,%o,%T\n" >> ${OUTPUT}
        ${KAFKACAT_CMD} -t ${t} -p ${p} -o -1 -c 1 -f "end,%t,%p,%o,%T\n" >> ${OUTPUT}
    done
done

exit 0
