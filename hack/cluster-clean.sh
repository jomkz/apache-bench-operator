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

HACK_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source ${HACK_DIR}/env.sh

# Clean up operator resources in a Kubernetes cluster

# Deployment (operator)
kubectl delete deployment -n ${AB_OPERATOR_NAMESPACE} ${AB_OPERATOR_NAME}

# Roles/Bindings
kubectl delete rolebinding -n ${AB_OPERATOR_NAMESPACE} ${AB_OPERATOR_NAME}
kubectl delete role -n ${AB_OPERATOR_NAMESPACE} ${AB_OPERATOR_NAME}

# ServiceAccount
kubectl delete sa -n ${AB_OPERATOR_NAMESPACE} ${AB_OPERATOR_NAME}

# CustomResourceDefinition
kubectl delete crd apachebenches.httpd.apache.org

# Scorecard artifacts
kubectl delete secret -n ${AB_OPERATOR_NAMESPACE} scorecard-kubeconfig
