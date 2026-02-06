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

set -e # Exit if any command fails

BLUE=$(tput setaf 4)
NC=$(tput sgr0) # No Color

log_debug() {
    echo "${BLUE}[DEBUG]${NC} $1"
}

run_bdd_tests() {
    LOCAL_PARQUET_FACTORY=$(pwd)
    cd /tmp
    log_debug "Using https://github.com/RedHatInsights/ccx-bdd-tests for running BDD tests."
    if [ -d "/tmp/ccx-bdd-tests" ]; then
        log_debug "The repository is already cloned. Updating to last version."
        cd ccx-bdd-tests
        git pull origin main
    else
        log_debug "Clonning CCX BDD tests repository in the temporary path /tmp/ccx-bdd-tests."
        git clone https://github.com/RedHatInsights/ccx-bdd-tests.git
        cd ccx-bdd-tests
    fi

    log_debug "Running BDD tests for $LOCAL_PARQUET_FACTORY"
    make LOCAL_SOURCE_CODE_FOLDER="$LOCAL_PARQUET_FACTORY" debug-parquet-factory

    cd "$LOCAL_PARQUET_FACTORY"
}

if [ "$CI" = true ] ; then
    # Don't ask if its being run in CI
    run_bdd_tests
else
    while true; do
        read -r -p "Do you wish to run BDD tests? This could take a while. [y/n] " yn
        case $yn in
            [Yy]* ) run_bdd_tests; break;;
            [Nn]* ) echo "No other integration tests implemented yet"; exit 0;;
            * ) echo "Please answer 'y' (yes) or 'n' (no).";;
        esac
    done
fi

exit 0