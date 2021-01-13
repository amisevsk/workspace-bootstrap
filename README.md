# Devfile 2.0 DevWorkspace Bootstrapper (demo)

This repo contains the definition for a simple Devfile 2.0 component that can be used to bootstrap a workspace from a publically-hosted Git repo, without needing to fetch the repo's devfile first.

Depends on the [DevWorkspace operator](https://github.com/devfile/devworkspace-operator) to be deployed on the cluster where the DevWorkspace is applied.

The default project (https://github.com/amisevsk/devfile-demo-workspace.git) requires changes from https://github.com/devfile/devworkspace-operator/pull/240 to be deployed to the cluster.

The bash script `bootstrap-workspace.sh` is a quick demo describing how this bootstrapper can be deployed to a cluster that supports it; see `./bootstrap-workspace.sh --help` for details.

## Design
The bootstrapper is implemented as a simple container that
1. Reads the current DevWorkspace spec for the workspace it's deployed in
2. Reads the first project in that DevWorkspace spec
3. Clones the repo pointed at by the first project to `${PROJECT_ROOT}` (set by default as a part of the devfile API)
4. Reads the devfile.yaml stored in the cloned repo (if any)
5. Merges the spec of the devfile.yaml with the current DevWorkspace, which results in the deployment being restarted with the spec from the repo's devfile.

## Quick testing instructions:
Requires: `minikube`, `kubectl`, and `make`
```bash
minikube start && minikube addons enable ingress
pushd $(mktemp -d)
git clone https://github.com/devfile/devworkspace-operator.git
pushd devworkspace-operator
git fetch origin pull/240/head:plugin-flattening
git checkout plugin-flattening
# optional: export IMG=my-img; make docker
export IMG="quay.io/amisevsk/devworkspace-controller:plugin-flattening"
make install_cert_manager install install_plugin_templates
popd; popd;
./bootstrap-workspace.sh # -r <REPO> -n <WORKSPACE_NAME>
```

## Additional Info
- devfile 2.0 API repo: https://github.com/devfile/api
- devworkspace-operator repo: https://github.com/devfile/devworkspace-operator