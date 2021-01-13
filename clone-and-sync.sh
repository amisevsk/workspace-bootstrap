#!/bin/bash

TMPDIR=/tmp/devworkspace-bootstrap
mkdir -p "$TMPDIR"
cd "$TMPDIR"

# Get DevWorkspace json
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
kubectl get devworkspace ${CHE_WORKSPACE_NAME} -n ${CHE_WORKSPACE_NAMESPACE} --token ${TOKEN} -o json > "${TMPDIR}/devworkspace.json"

cat "${TMPDIR}/devworkspace.json"

# Get project info from devworkspace.json
PROJECT_NAME=$(jq -r '.spec.template.projects[0].name' devworkspace.json)
PROJECT=$(jq -r '.spec.template.projects[0].git.remotes.origin' devworkspace.json)

# Clone project
cd "${PROJECTS_ROOT}"
git clone "${PROJECT}" "${PROJECT_NAME}"
cd "${PROJECT_NAME}"

# Read devfile from project and update devworkspace
if [ ! -f devfile.yaml ]; then
  echo "Could not find devfile.yaml in project"
  # sleep? exit?
fi
yq '{parent, projects, components, commands, events}' devfile.yaml > "${TMPDIR}/git-devworkspace.json"
cat "${TMPDIR}/git-devworkspace.json"

cd "$TMPDIR"
jq -sc '.[0] + {"spec": {"template": .[1]}}' devworkspace.json git-devworkspace.json > merged-devworkspace.json

cat merged-devworkspace.json
kubectl patch devworkspace ${CHE_WORKSPACE_NAME} -n ${CHE_WORKSPACE_NAMESPACE} --token ${TOKEN} \
  --type merge --patch "$(cat merged-devworkspace.json)"
