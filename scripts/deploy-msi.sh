#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

###############################################################################

set -e
set -u
set -o pipefail

ROOT="${DIR}/.."

cd "${ROOT}"
make

source "${ROOT}/test/common.sh"

export INSTANCE_NAME_DEFAULT="${INSTANCE_NAME_PREFIX}-$(printf "%x" $(date '+%s'))-${LOCATION}"
export INSTANCE_NAME="${INSTANCE_NAME:-${INSTANCE_NAME_DEFAULT}}"
export CLUSTER_DEFINITION="${ROOT}/examples/kubernetes.json"

#export CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID="msi"
#export CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET="msi"
export CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID="${SERVICE_PRINCIPAL_CLIENT_ID}"
export CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET="${SERVICE_PRINCIPAL_CLIENT_SECRET}"

export CUSTOM_HYPERKUBE_SPEC="docker.io/colemickens/hyperkube-amd64:v1.5.2-lazyinit-msi"

deploy