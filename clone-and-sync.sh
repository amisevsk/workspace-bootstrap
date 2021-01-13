#!/bin/bash

set -e

TMPDIR=/tmp/devworkspace-bootstrap
mkdir -p "$TMPDIR"
cd "$TMPDIR"

# Get DevWorkspace json
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
kubectl get devworkspace ${CHE_WORKSPACE_NAME} \
  --namespace $(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace) \
  --token ${TOKEN} \
  -o json > "${TMPDIR}/devworkspace.json"

echo "Read DevWorkspace:"
cat "${TMPDIR}/devworkspace.json"

# Get project info from devworkspace.json
PROJECT_NAME=$(jq -r '.spec.template.projects[0].name' devworkspace.json)
PROJECT=$(jq -r '.spec.template.projects[0].git.remotes.origin' devworkspace.json)

if [ ! -d "${PROJECTS_ROOT}/${PROJECT_NAME}" ]; then
  echo "Cloning Devworkspace project from ${PROJECT} to ${PROJECTS_ROOT}/${PROJECT_NAME}"
  # Clone project
  cd "${PROJECTS_ROOT}"
  git clone "${PROJECT}" "${PROJECT_NAME}"
  cd "${PROJECT_NAME}"
else
  echo "Project already cloned to ${PROJECTS_ROOT}/${PROJECT_NAME}"
fi

if [ -f devfile.yaml ]; then
  export DEVFILE="${PROJECTS_ROOT}/${PROJECT_NAME}/devfile.yaml"
  echo "Found devfile in repository: ${DEVFILE}"
fi

exec "$@"
