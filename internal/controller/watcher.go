package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"log"
)

type DeploymentWatcher struct {
	clientset *kubernetes.Clientset
	stopCh    chan struct{}
}

// NewDeploymentWatcher returns a new deployment watcher instance.
func NewDeploymentWatcher(clientset *kubernetes.Clientset) *DeploymentWatcher {
	return &DeploymentWatcher{
		clientset: clientset,
		stopCh:    make(chan struct{}),
	}
}

// Start begins watching deployments with kubenap/enabled: "true" label.
func (dw *DeploymentWatcher) Start(ctx context.Context) error {
	factory := informers.NewSharedInformerFactory(dw.clientset, 0)
	informer := factory.Apps().V1().Deployments().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dw.handleDeployment,
		UpdateFunc: func(oldObj, newObj interface{}) { dw.handleDeployment(newObj) },
	})

	go informer.Run(dw.stopCh)

	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		return fmt.Errorf("timed out waiting for cache sync")
	}

	log.Println("Deployment Watcher Running..")
	<-ctx.Done()
	return nil
}
func (dw *DeploymentWatcher) handleDeployment(obj interface{}) {
	dep, ok := obj.(*appsv1.Deployment)
	if !ok {
		return
	}

	if val, ok := dep.Labels["kubenap/enabled"]; ok && val == "true" {
		log.Printf("Watched deployment: %s/%s", dep.Namespace, dep.Name)
		// TODO: track or enqueue for suspend handling
	}
}
