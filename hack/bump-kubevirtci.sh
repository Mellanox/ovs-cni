#!/bin/bash -e
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


KUBEVIRTCI_TAG=$(curl -L -Ss https://storage.googleapis.com/kubevirt-prow/release/kubevirt/kubevirtci/latest)
[[ ${#KUBEVIRTCI_TAG} != "18" ]] && echo "error getting KUBEVIRTCI_TAG" && exit 1

sed -i "s/\(KUBEVIRTCI_TAG:-\)[^}]*/\1${KUBEVIRTCI_TAG}/" cluster/cluster.sh

git --no-pager diff cluster/cluster.sh | grep KUBEVIRTCI_TAG || true
