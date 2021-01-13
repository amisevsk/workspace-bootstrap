package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	devworkspace "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	repoDevfileEnvVar    = "DEVFILE"
	defaultDevfileEnvVar = "DEFAULT_DEVFILE"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(devworkspace.AddToScheme(scheme))
}

func stop(err error) {
	if err != nil {
		log.Printf("Failed to bootstrap workspace: %s", err.Error())
	}
	time.Sleep(60 * time.Minute)
}

func main() {
	log.Println("Beginning devfile bootstrap")
	client, err := getK8sClient()
	if err != nil {
		stop(err)
	}
	dw, err := readDevWorkspace(client)
	if err != nil {
		stop(err)
	}
	log.Println("Read DevWorkspace on cluster complete")
	devfile, err := getActualDevfile()
	if err != nil {
		stop(err)
	}
	log.Println("Read devfile complete")

	dw.Spec.Template = devfile.Spec.Template
	err = client.Patch(context.Background(), dw, k8sclient.Merge)
	if err != nil {
		stop(fmt.Errorf("failed to update DevWorkspace with devfile from repository: %w", err))
	}
	log.Println("Updated DevWorkspace with spec from repository")
	stop(nil)
}

func getActualDevfile() (*devworkspace.DevWorkspace, error) {
	repoDevfilePath := os.Getenv(repoDevfileEnvVar)
	defaultDevfilePath := os.Getenv(defaultDevfileEnvVar)
	log.Printf("Reading devfile.yaml from repo cloned to %s", strings.TrimSuffix(repoDevfilePath, "devfile.yaml"))
	if repoDevfilePath == "" && defaultDevfilePath == "" {
		return nil, fmt.Errorf("could not find devfile and no default is set")
	}
	if repoDevfilePath != "" {
		devfile, err := readDevfile(repoDevfilePath)
		if err != nil {
			if errors.Is(err, errInvalidSchemaVersion) {
				log.Printf("Devfile found in repository is unsupported; using default DevWorkspace")
				return readDevfile(defaultDevfilePath)
			}
			return nil, fmt.Errorf("failed to read repo devfile: %w", err)
		}
		return devfile, nil
	}
	log.Printf("Cloned repository does not contain devfile.yaml; using default DevWorkspace")
	return readDevfile(defaultDevfilePath)
}

func getK8sClient() (k8sclient.Client, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := k8sclient.New(cfg, k8sclient.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}
