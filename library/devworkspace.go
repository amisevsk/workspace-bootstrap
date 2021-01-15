package library

import (
	"context"
	"errors"
	"fmt"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// GetDevWorkspace reads the current DevWorkspace from the cluster; depends on workspaceNameEnvVar
func GetDevWorkspace(client client.Client) (*v1alpha2.DevWorkspace, error) {
	namespace, err := GetCurrentNamespace()
	if err != nil {
		return nil, fmt.Errorf("failed to get current Kubernetes namespace: %w", err)
	}
	dwName := os.Getenv(workspaceNameEnvVar)
	if dwName == "" {
		return nil, fmt.Errorf("could not get current DevWorkspace name")
	}
	namespacedName := types.NamespacedName{
		Name:      dwName,
		Namespace: namespace,
	}
	dw := &v1alpha2.DevWorkspace{}
	err = client.Get(context.Background(), namespacedName, dw)
	if err != nil {
		return nil, fmt.Errorf("failed to read DevWorkspace from cluster: %w", err.Error())
	}
	return dw, nil
}

func GetDevWorkspaceTemplate(client client.Client, name, namespace string) (*v1alpha2.DevWorkspaceTemplate, error) {
	namespacedName := types.NamespacedName{
		Name:name,
		Namespace: namespace,
	}
	dwt := &v1alpha2.DevWorkspaceTemplate{}
	err := client.Get(context.Background(), namespacedName, dwt)
	if err != nil {
		return nil, fmt.Errorf("failed to read DevWorkspaceTemplate from cluster: %w", err)
	}
	return dwt, nil
}

func ReadDevWorkspaceTemplateFromFile(path string) (*v1alpha2.DevWorkspaceTemplate, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read DevWorkspaceTemplate at path %s: %w", path, err)
	}
	dwt := &v1alpha2.DevWorkspaceTemplate{}
	err = yaml.Unmarshal(bytes, dwt)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DevWorkspaceTemplate: %w", err)
	}
	return dwt, nil
}

func FindDevWorkspaceTemplate() (string, error) {
	projectsRoot := os.Getenv(projectsRootEnvVar)
	if projectsRoot == "" {
		projectsRoot = defaultProjectsRoot
	}

	var dwtFile string
	found := false
	multipleDWTsMarker := errors.New("found multiple DevWorkspaceTemplate")
	err := filepath.Walk(projectsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == devworkspaceTemplateName {
			if found {
				return multipleDWTsMarker
			}
			dwtFile = path
			found = true
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, multipleDWTsMarker) {
			return "", fmt.Errorf("found multiple files named %s in projects", devworkspaceTemplateName)
		}
		return "", fmt.Errorf("failed to find devworkspace template in %s: %w", projectsRoot, err)
	}
	if !found {
		return "", fmt.Errorf("failed to find devworkspace template in %s", projectsRoot)
	}
	return dwtFile, nil
}