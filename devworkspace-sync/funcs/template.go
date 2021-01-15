package funcs

import (
	"context"
	"fmt"
	"github.com/amisevsk/workspace-bootstrap/library"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

// SyncDevWorkspaceTemplate finds a file named devworkspace.template.yaml and syncs it to the cluster. If template
// already exists on the cluster and is not controlled by the current workspace, an error is returned.
func SyncDevWorkspaceTemplate() error {
	client, err := library.GetK8sClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %w", err)
	}

	dwtPath, err := library.FindDevWorkspaceTemplate()
	if err != nil {
		return fmt.Errorf("failed to find DevWorkspaceTemplate: %w", err)
	}
	log.Printf("Found DevWorkspaceTemplate at %s", dwtPath)

	dwt, err := library.ReadDevWorkspaceTemplateFromFile(dwtPath)
	if err != nil {
		return fmt.Errorf("failed to read DevWorkspaceTemplate: %w", err)
	}
	log.Printf("Successfully read DevWorkspaceTemplate")

	if dwt.Namespace == "" {
		ns, err := library.GetCurrentNamespace()
		if err != nil {
			return fmt.Errorf("failed to get the current namespace and DevWorkspaceTemplate does not supply one")
		}
		dwt.Namespace = ns
	}

	log.Printf("Reading current DevWorkspace to ensure Template is compatible")
	// Get current DevWorkspace to set ownerref correctly
	currDW, err := library.GetDevWorkspace(client)
	if err != nil {
		return fmt.Errorf("failed to get the current DevWorkspace: %w", err)
	}
	err = library.SetControllerRef(currDW, dwt)
	if err != nil {
		return fmt.Errorf("failed to set controllerref on DevWorkspaceTemplate")
	}

	err = client.Create(context.Background(), dwt)
	if err == nil {
		log.Printf("Successfully created DevWorkspaceTemplate on cluster: name %q, namespace: %q", dwt.Name, dwt.Namespace)
		return nil
	}
	if !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("encountered unexpected error when trying to sync template to cluster (retry): %w", err)
	}

	// template exists, need to update (if allowed)
	log.Printf("DevWorkspaceTemplate with name %s already exists in current namespace; updating", dwt.Name)
	clusterDWT := &v1alpha2.DevWorkspaceTemplate{}
	err = client.Get(context.Background(), types.NamespacedName{Name: dwt.Name, Namespace: dwt.Namespace}, clusterDWT)
	if err != nil {
		return fmt.Errorf("unexpected error while trying to get devworkspace template from cluster (retry): %w", err)
	}
	if !library.OwnerRefsMatch(metav1.GetControllerOf(dwt), metav1.GetControllerOf(clusterDWT)) {
		return fmt.Errorf("template already exists on cluster and is not controlled by the current DevWorkspace")
	}
	clusterDWT.Spec = dwt.Spec
	err = client.Update(context.Background(), clusterDWT)
	if err != nil {
		return fmt.Errorf("failed to update DevWorkspaceTemplate on cluster: %w", err)
	}
	log.Printf("Successfully updated DevWorkspaceTemplate on cluster")
	return nil
}
