#!/bin/bash

# Copyright 2018 Google, Inc. All rights reserved.
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

source "$(dirname $(readlink -f ${BASH_SOURCE}))/../test/library.sh"

: ${PROJECT_ID:="knative-environments"}
readonly PROJECT_ID
readonly K8S_CLUSTER_NAME=${1:?"First argument must be the kubernetes cluster name."}
readonly K8S_CLUSTER_ZONE=us-central1-a
readonly K8S_CLUSTER_MACHINE=n1-standard-8
readonly K8S_CLUSTER_NODES=5
readonly ISTIO_VERSION=0.8.0
readonly SERVING_RELEASE=https://storage.googleapis.com/knative-releases/latest/release.yaml
readonly ISTIO_YAML=https://storage.googleapis.com/knative-releases/latest/istio.yaml
export ISTIO_VERSION
readonly PROJECT_USER=$(gcloud config get-value core/account)
readonly CURRENT_PROJECT=$(gcloud config get-value project)

function cleanup() {
  gcloud config set project ${CURRENT_PROJECT}
}

cd ${SERVING_ROOT_DIR}
trap cleanup EXIT

echo "Using project ${PROJECT_ID} and user ${PROJECT_USER}"
gcloud config set project ${PROJECT_ID}

existing_cluster=$(gcloud container clusters list \
    --zone=${K8S_CLUSTER_ZONE} --filter="name=${K8S_CLUSTER_NAME}")
if [[ -n "${existing_cluster}" ]]; then
  header "Deleting previous cluster ${K8S_CLUSTER_NAME} in ${PROJECT_ID}"
  gcloud -q container clusters delete \
    --zone=${K8S_CLUSTER_ZONE} ${K8S_CLUSTER_NAME}
fi

header "Creating cluster ${K8S_CLUSTER_NAME} in ${PROJECT_ID}"
gcloud --project=${PROJECT_ID} container clusters create \
  --cluster-version=${SERVING_GKE_VERSION} \
  --image-type "${SERVING_GKE_IMAGE^^}" \
  --zone=${K8S_CLUSTER_ZONE} \
  --scopes=cloud-platform \
  --machine-type=${K8S_CLUSTER_MACHINE} \
  --enable-autoscaling --min-nodes=1 --max-nodes=${K8S_CLUSTER_NODES} \
  ${K8S_CLUSTER_NAME}

header "Setting cluster admin"
acquire_cluster_admin_role ${PROJECT_USER} ${K8S_CLUSTER_NAME} ${K8S_CLUSTER_ZONE}

header "Installing istio"
kubectl apply -f ${ISTIO_YAML}
wait_until_pods_running istio-system

kubectl label namespace default istio-injection=enabled

header "Installing Knative Serving"
# Install might fail before succeding, so we retry a few times.
# For details, see https://github.com/knative/install/issues/13
installed=0
for i in {1..10}; do
  kubectl apply -f ${SERVING_RELEASE} && installed=1 && break
  sleep 30
done
if (( ! installed )); then
  echo "ERROR: could not install Knative Serving"
  exit 1
fi

wait_until_pods_running knative-serving-system
wait_until_pods_running build-system

header "Knative Serving deployed successfully to ${K8S_CLUSTER_NAME}"
