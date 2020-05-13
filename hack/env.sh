#!/usr/bin/env bash

# Copyright 2020 Apache Bench Operator Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# General vars
export AB_OPERATOR_NAME=${AB_OPERATOR_NAME:-"apache-bench-operator"}
export AB_OPERATOR_NAMESPACE=${AB_OPERATOR_NAMESPACE:-"benchmark"}
export AB_OPERATOR_VERSION=${AB_OPERATOR_VERSION:-`awk '$1 == "Version" {gsub(/"/, "", $3); print $3}' version/version.go`}
export AB_OPERATOR_PREVIOUS_VERSION=${AB_OPERATOR_PREVIOUS_VERSION:-`awk '$1 == "Version" {gsub(/"/, "", $3); print $3}' version/version.go`}
export AB_OPERATOR_BUILD_DIR=${AB_OPERATOR_BUILD_DIR:-"build"}
export AB_OPERATOR_DEPLOY_DIR=${AB_OPERATOR_DEPLOY_DIR:-"deploy"}
export AB_OPERATOR_DOCS_DIR=${AB_OPERATOR_DOCS_DIR:-"docs"}
export AB_OPERATOR_BRANCH_NAME=${AB_OPERATOR_BRANCH_NAME:-`git status -b -uno | awk 'NR==1{print $3}'`}

# Container image vars
export AB_OPERATOR_IMAGE_BUILDER=${AB_OPERATOR_IMAGE_BUILDER:-"podman"}
export AB_OPERATOR_IMAGE_REPO=${AB_OPERATOR_IMAGE_REPO:-"quay.io/jmckind/${AB_OPERATOR_NAME}"}
export AB_OPERATOR_IMAGE_TAG=${AB_OPERATOR_IMAGE_TAG:-${AB_OPERATOR_BRANCH_NAME}}
export AB_OPERATOR_IMAGE=${AB_OPERATOR_IMAGE:-"${AB_OPERATOR_IMAGE_REPO}:${AB_OPERATOR_IMAGE_TAG}"}

# Ensure go module support is enabled
export GO111MODULE=on
