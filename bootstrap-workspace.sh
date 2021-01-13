#!/bin/bash
#
# This script is a demo of how workspaces can be bootstrapped, and what info is
# needed for the process
#

set -e

NAME="workspace-bootstrap"
REPO="https://github.com/amisevsk/devfile-demo-workspace.git"

USAGE="
Usage: ./bootstrap-workspace.sh [OPTIONS]
Options:
    --help
        Print this message.
    --workspace-name, -n <NAME>
        Name used for the DevWorkspace (default: workspace-bootstrap)
    --repo, -r <REPO>
        Repository to create DevWorkspace from (default: https://github.com/amisevsk/devfile-demo-workspace.git)
    --project-name <PROJECT_NAME>
        Name of project in DevWorkspace (default: from repo URL)
"

function print_usage() {
    echo -e "$USAGE"
}

function parse_arguments() {
    while [[ $# -gt 0 ]]; do
        key="$1"
        case $key in
            -n|--workspace-name)
            WORKSPACE_NAME="$2"
            shift; shift;
            ;;
            -r|--repo)
            REPO="$2"
            shift; shift;
            ;;
            --project-name)
            PROJECT_NAME="$2"
            shift; shift;
            ;;
            *)
            print_usage
            exit 0
        esac
    done
}

parse_arguments "$@"

if [ -z PROJECT_NAME ]; then
  PROJECT_NAME="${REPO##*/}"
  PROJECT_NAME="${PROJECT_NAME%.git}"
  echo "Using '$PROJECT_NAME' as project name"
fi

cat <<EOF | cat
kind: DevWorkspace
apiVersion: workspace.devfile.io/v1alpha2
metadata:
  name: ${WORKSPACE_NAME}
spec:
  started: true
  routingClass: 'basic'
  template:
    projects:
      - name: ${PROJECT_NAME}
        git:
          remotes:
            origin: ${REPO}
    components:
      - name: workspace-root
        container:
          image: docker.io/amisevsk/workspace-root:dev
EOF