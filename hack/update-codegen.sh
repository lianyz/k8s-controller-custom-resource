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

if [ "$#" -lt 2 ] || [ "${1}" == "--help" ]; then
  cat <<EOF
Usage: $(basename "$0") <group> <version>

  <group> group of controller
  <version> version of controller

Examples:
  $(basename "$0") samplecontroller v1
EOF
  exit 0
fi

# 注意事项：该脚本生成的zz_generated_deepcopy.go路径不对
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
# CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
CODEGEN_PKG=${GOPATH}/src/github.com/kubernetes/code-generator
# 以下命令等价于: MODULE=$(cat go.mod | head -n 1 | awk '{print $2}')
# 读取go.mod的第一行，用awk分割第一行，读取第二个词
MODULE=$(head -n 1 < go.mod | awk '{print $2}')
OUTPUT_PKG=pkg/client
APIS_PKG=pkg/apis
GROUP=$1
VERSION=$2
bash "${CODEGEN_PKG}"/generate-groups.sh all \
"${MODULE}"/${OUTPUT_PKG} "${MODULE}"/${APIS_PKG} \
"${GROUP}":"${VERSION}" \
--output-base "${SCRIPT_ROOT}"/../../.. \
--go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt \
