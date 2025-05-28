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

# Prepare environment for testing and automation.
#
# source automation/setup.sh
# cd ${TMP_PROJECT_PATH}

echo 'Copy the project under a temporary'
TMP_PROJECT_PATH=/tmp/src/github.com/k8snetworkplumbingwg/ovs-cni
rm -rf ${TMP_PROJECT_PATH}
mkdir -p ${TMP_PROJECT_PATH}
cp -rf $(pwd)/. ${TMP_PROJECT_PATH}

echo 'Exporting temporary project path'
export TMP_PROJECT_PATH
