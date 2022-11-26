#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

# 注意事项：该脚本生成的zz_generated_deepcopy.go路径不对
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
# CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
CODEGEN_PKG=${GOPATH}/src/k8s.io/code-generator
MODULE=github.com/lianyz/k8s-controller-custom-resource
OUTPUT_PKG=pkg/client
APIS_PKG=pkg/apis
GROUP=samplecrd
VERSION=v1
bash "${CODEGEN_PKG}"/generate-groups.sh all \
${OUTPUT_PKG} ${MODULE}/${APIS_PKG} \
${GROUP}:${VERSION} \
--output-base "${SCRIPT_ROOT}" \
--go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt \
#####################样例 start##################################
#注意事项：
#MODULE需和go.mod文件内容一致
#"${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
#  sample-controller/pkg/generated sample-controller/pkg/apis \
#  samplecontroller:v1 \
#  --output-base "$(dirname "${BASH_SOURCE[0]}")/../.." \
#  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt
#####################样例 end##################################