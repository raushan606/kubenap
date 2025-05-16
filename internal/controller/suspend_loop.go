package controller

import (
	"context"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

type SuspendLoop struct {
	Clientset      kubernetes.Interface
	ActivityCache  *ActivityCache
	DeploymentList func() []*appsv1.Deployment // injects watched deployments
	ConfigLoader   func(*appsv1.Deployment) (*ParsedConfig, error)
	Suspender      *SuspendController
	Interval       time.Duration
}

// Start runs the suspend loop in the background
func (sl *SuspendLoop) Start(ctx context.Context) {
	ticker := time.NewTicker(sl.Interval)
	defer ticker.Stop()

	log.Printf("Suspend loop started (interval: %s)", sl.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Suspend loop exiting")
			return
		case <-ticker.C:
			sl.runOnce(ctx)
		}
	}
}

// runOnce checks all deployments for idleness
func (sl *SuspendLoop) runOnce(ctx context.Context) {
	deployments := sl.DeploymentList()
	for _, dep := range deployments {
		cfg, err := sl.ConfigLoader(dep)
		if err != nil {
			log.Printf("Failed to load config for %s/%s: %v", dep.Namespace, dep.Name, err)
			continue
		}
		sl.Suspender.SuspendIfIdle(ctx, dep, cfg)
	}
}