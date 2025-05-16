package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SuspendController struct {
	Clientset     kubernetes.Interface
	ActivityCache *ActivityCache
}

const AnnotationOriginalReplicas = "kubenap/originalReplicas"

// SuspendIfIdle checks whether a deployment should be suspended and scales it down
func (sc *SuspendController) SuspendIfIdle(ctx context.Context, dep *appsv1.Deployment, cfg *ParsedConfig) {
	key := dep.Namespace + "/" + dep.Name

	lastSeen, ok := sc.ActivityCache.Get(dep.Namespace, dep.Name)

	if !ok {
		log.Printf("No activity recorded for %s, skipping suspended check", key)
		return
	}

	if time.Since(lastSeen) < cfg.IdleAfter {
		log.Printf("%s is still active (last seen %s ago)", key, time.Since(lastSeen))
		return
	}

	if dep.Spec.Replicas != nil && *dep.Spec.Replicas == 0 {
		log.Printf("%s is already scaled down", key)
		return
	}

	if dep.Annotations == nil || dep.Annotations[AnnotationOriginalReplicas] == "" {
		currentReplicas := int32(1)
		if dep.Spec.Replicas != nil {
			currentReplicas = *dep.Spec.Replicas
		}

		patch := map[string]interface{}{
			"metadata": map[string]interface{}{
				"annotations": map[string]string{
					AnnotationOriginalReplicas: fmt.Sprintf("%d", currentReplicas),
				},
			},
		}

		data, err := json.Marshal(patch)
		if err == nil {
			_, err = sc.Clientset.AppsV1().Deployments(dep.Namespace).Patch(
				ctx,
				dep.Name,
				types.MergePatchType,
				data,
				metav1.PatchOptions{},
			)
			if err != nil {
				log.Printf("Failed to annotate original replica count for %s/%s: %v", dep.Namespace, dep.Name, err)
			}
		} else {
			log.Printf("Failed to serialize patch for originalReplicas: %v", err)
		}
	}

	// Patch replicas to 0
	zero := int32(0)
	newDep := dep.DeepCopy()
	newDep.Spec.Replicas = &zero

	_, err := sc.Clientset.AppsV1().Deployments(dep.Namespace).Update(ctx, newDep, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to scale down %s: %v", key, err)
		return
	}
	log.Printf("Suspended %s due to inactivity (last seen %s ago)", key, time.Since(lastSeen))
}
