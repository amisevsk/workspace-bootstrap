package library

import (
	"errors"
	"fmt"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/pkg/devfile"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sigs.k8s.io/yaml"
)

var (
	ErrInvalidSchemaVersion    = errors.New("unsupported schemaVersion found in devfile")
	devfileSchemaVersionRegexp = regexp.MustCompile(`2\.[0-9]+\..*`)
)

// devfileScaffold can be used to unmarshal a devfile.yaml to a Go struct
type devfileScaffold struct {
	devfile.DevfileHeader
	v1alpha2.DevWorkspaceTemplateSpec
}

// ReadDevfile reads a devfile from a path on disk
func ReadDevfile(path string) (*v1alpha2.DevWorkspace, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read devfile at path %s", path)
	}
	devfileYaml := devfileScaffold{}
	err = yaml.Unmarshal(bytes, &devfileYaml)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal devfile at path %s: %w", path, err)
	}
	if !devfileSchemaVersionRegexp.MatchString(devfileYaml.SchemaVersion) {
		return nil, ErrInvalidSchemaVersion
	}

	dw := &v1alpha2.DevWorkspace{}
	dw.Name = devfileYaml.Metadata.Name
	dw.Spec.Template.Parent = devfileYaml.Parent
	dw.Spec.Template.Projects = devfileYaml.Projects
	dw.Spec.Template.Components = devfileYaml.Components
	dw.Spec.Template.Commands = devfileYaml.Commands
	dw.Spec.Template.Events = devfileYaml.Events

	return dw, nil
}

func FindRepoDevfile() (string, error) {
	projectsRoot := os.Getenv(projectsRootEnvVar)
	if projectsRoot == "" {
		projectsRoot = defaultProjectsRoot
	}
	var devfilePath string
	devfiles, err := filepath.Glob(filepath.Join(projectsRoot, "*", devfileName))
	if err != nil {
		return "", err
	}
	switch len(devfiles) {
	case 0:
		return "", fmt.Errorf("no devfile found")
	case 1:
		devfilePath = devfiles[0]
	default:
		return "", fmt.Errorf("multiple devfiles found in repo")
	}
	return devfilePath, nil
}
