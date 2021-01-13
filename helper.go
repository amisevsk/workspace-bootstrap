package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	devworkspace "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/pkg/devfile"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

const (
	workspaceNameEnvVar = "CHE_WORKSPACE_NAME" // Fragile; should be updated.
)

var (
	errInvalidSchemaVersion    = errors.New("unsupported schemaVersion found in devfile")
	devfileSchemaVersionRegexp = regexp.MustCompile(`2\.[0-9]+\..*`)
)

// devfileScaffold can be used to unmarshal a devfile.yaml to a Go struct
type devfileScaffold struct {
	devfile.DevfileHeader
	devworkspace.DevWorkspaceTemplateSpec
}

// readDevWorkspace reads the current DevWorkspace from the cluster; depends on workspaceNameEnvVar
func readDevWorkspace(client client.Client) (*devworkspace.DevWorkspace, error) {
	namespace, err := getCurrentNamespace()
	if err != nil {
		return nil, fmt.Errorf("failed to get current Kubernetes namespace: %w", err)
	}
	dwName := os.Getenv(workspaceNameEnvVar)
	if dwName == "" {
		return nil, fmt.Errorf("could not get current DevWorkspace name")
	}
	fmt.Printf("Reading DevWorkspace with name %s in namespace %s", dwName, namespace)
	namespacedName := types.NamespacedName{
		Name:      dwName,
		Namespace: namespace,
	}
	dw := &devworkspace.DevWorkspace{}
	err = client.Get(context.Background(), namespacedName, dw)
	if err != nil {
		return nil, fmt.Errorf("failed to read DevWorkspace from cluster: %s", err.Error())
	}
	return dw, nil
}

// readDevfile reads a devfile from a path on disk
func readDevfile(path string) (*devworkspace.DevWorkspace, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read devfile at path %s", path)
	}
	log.Printf("Read devfile: %s", bytes)
	devfileYaml := devfileScaffold{}
	err = yaml.Unmarshal(bytes, &devfileYaml)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal devfile at path %s: %w", path, err)
	}
	if !devfileSchemaVersionRegexp.MatchString(devfileYaml.SchemaVersion) {
		log.Printf("Unsupported schemaVersion found in devfile: %q", devfileYaml.SchemaVersion)
		return nil, errInvalidSchemaVersion
	}

	dw := &devworkspace.DevWorkspace{}
	dw.Name = devfileYaml.Metadata.Name
	dw.Spec.Template.Parent = devfileYaml.Parent
	dw.Spec.Template.Projects = devfileYaml.Projects
	dw.Spec.Template.Components = devfileYaml.Components
	dw.Spec.Template.Commands = devfileYaml.Commands
	dw.Spec.Template.Events = devfileYaml.Events

	return dw, nil
}

// getCurrentNamespace get the current namespace this container is running in
func getCurrentNamespace() (string, error) {
	nsBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("could not read namespace from mounted serviceaccount info")
		}
		return "", err
	}
	ns := strings.TrimSpace(string(nsBytes))
	return ns, nil
}
