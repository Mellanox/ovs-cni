#!/bin/bash
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


allowed_hashicorp_modules=(
    "github.com/hashicorp/errwrap"
    "github.com/hashicorp/go-multierror"
    "github.com/hashicorp/hcl"
)

error_found=false
while read -r line; do
    if ! [[ " ${allowed_hashicorp_modules[*]} " == *" $line "* ]]; then
        echo "found non allowlisted hashicorp module: $line"
        error_found=true
    fi
done < <(grep -i hashicorp go.mod | grep -o 'github.com/[^ ]*')

if [[ $error_found == true ]]; then
    echo "Non allowlisted hashicorp modules found, exiting with an error."
    echo "HashiCorp adapted BSL, which we cant use on our projects."
    echo "Please review the licensing, and either add it to the list if it isn't BSL,"
    echo "or use a different library."
    exit 1
fi
echo "All included hashicorp modules are allowlisted"
