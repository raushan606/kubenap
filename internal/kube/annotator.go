package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	
)

const LastActivityKey = "kubenap/lastActivityAt"

// PatchLastActivityAt sets the lastActivityAt annotation on a deployment.
func PatchLastActivityAt(ctx context.Context, clientset kubernetes.Interface, namespace, name string, ts time.Time) error {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]string{
				LastActivityKey: ts.UTC().Format(time.RFC3339),
			},
		},
	}

	data, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("failed to marshal patch: %w", err)
	}

	_, err = clientset.AppsV1().Deployments(namespace).Patch(
		ctx,
		name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to patch deployment: %w", err)
	}

	return nil
}