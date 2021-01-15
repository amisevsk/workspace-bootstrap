package funcs

import (
	"context"
	"fmt"
	"github.com/amisevsk/workspace-bootstrap/library"
	"log"
)

// SyncDevfile syncs the devfile defined in the repository to the cluster
// The DevWorkspace may be restarted after this operation.
func SyncDevfile() error {
	client, err := library.GetK8sClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	dwPath, err := library.FindRepoDevfile()
	if err != nil {
		return fmt.Errorf("failed to find devfile.yaml in repository: %w", err)
	}
	log.Printf("Found devfile in repo: %s", dwPath)

	dw, err := library.ReadDevfile(dwPath)
	if err != nil {
		return fmt.Errorf("failed to read devfile.yaml from repository: %w", err)
	}
	log.Printf("Successfully processed devfile from %s", dwPath)

	clusterDW, err := library.GetDevWorkspace(client)
	if err != nil {
		return fmt.Errorf("failed to get current DevWorkspace from cluster: %w", err)
	}
	clusterDW.Spec.Template = dw.Spec.Template

	err = client.Update(context.Background(), clusterDW)
	if err != nil {
		return fmt.Errorf("failed to update DevWorkspace: %w", err)
	}
	log.Printf("Devfile successfully applied. Workspace may restart.")
	return nil
}