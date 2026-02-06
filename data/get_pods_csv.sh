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


PODNAMES=`oc get pods | awk '$1 ~ /parquet-factory/ { print $1 }'`

if [[ $? != 0 ]]; then
    echo "Error retrieving the parquet-factory pods from cluster. Did you performed `oc login`?"
    exit 1
else
    echo "Parquet Factory pods list retrieved"
fi
OUTPUT="pods_timing.csv"

# Print header of the CSV
echo "pod_name;exit_status;start_time;end_time;duration;num_messages" > ${OUTPUT}

for pod in ${PODNAMES}; do
    echo -n "$pod;" >> ${OUTPUT}
    oc describe pod ${pod} | awk -v ORS=";" '/Started:|Finished:/ { $1=""; print } /Exit\sCode:/ { print $NF}' >> ${OUTPUT}
    echo -n ";" >> ${OUTPUT}
    oc logs ${pod} | grep "Received message" | grep -v "count=" | wc -l >> ${OUTPUT}
    # echo "" >> ${OUTPUT}
done

exit 0
