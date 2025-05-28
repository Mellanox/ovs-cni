#!/bin/bash -xe
# Copyright 2025 ovs-cni authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0


# This script should be able to execute functional tests against Kubernetes
# cluster on any environment with basic dependencies listed in
# check-patch.packages installed and podman / docker running.
#
# yum -y install automation/check-patch.packages
# automation/check-patch.e2e.sh

teardown() {
    make cluster-down
}

main() {
    export KUBEVIRT_PROVIDER='k8s-1.26-centos9'

    source automation/setup.sh
    cd ${TMP_PROJECT_PATH}

    echo 'Run golint'
    make lint

    echo 'Run functional tests'
    make docker-test

    echo 'Run e2e tests'
    make cluster-down
    make cluster-up
    trap teardown EXIT SIGINT SIGTERM SIGSTOP
    make cluster-sync
    make E2E_TEST_ARGS="-ginkgo.v -test.v -ginkgo.noColor -test.timeout 20m --junit-output=$ARTIFACTS/junit.functest.xml" functest
}

[[ "${BASH_SOURCE[0]}" == "$0" ]] && main "$@"
