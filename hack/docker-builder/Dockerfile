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

FROM quay.io/fedora/fedora:40-x86_64

RUN dnf -y install make git sudo gcc rsync-daemon rsync openvswitch hostname && \
    dnf -y clean all

ENV GOPATH="/go"
RUN \
    DESTINATION=/opt && \
    VERSION=1.22.9 && \
    TARBALL=go${VERSION}.linux-amd64.tar.gz && \
    URL=https://dl.google.com/go && \
    mkdir -p ${DESTINATION} && \
    curl -L ${URL}/${TARBALL} -o ${DESTINATION}/${TARBALL} && \
    tar -xf ${DESTINATION}/${TARBALL} -C ${DESTINATION}
ENV PATH="/opt/go/bin:$PATH"

ADD rsyncd.conf /etc/rsyncd.conf

RUN mkdir -p /opt/cni/bin

ENV PATH="/opt/cni/bin:$PATH"

ADD https://github.com/containernetworking/plugins/releases/download/v1.0.1/cni-plugins-linux-amd64-v1.0.1.tgz /opt/cni/bin

RUN \
    cd  /opt/cni/bin && \
    tar -xzf cni-plugins-linux-amd64-v1.0.1.tgz && \
    rm -f cni-plugins-linux-amd64-v1.0.1.tgz

RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest

ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
