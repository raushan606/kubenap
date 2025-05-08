package controller

import (
	"context"
	"time"
)

// DeploymentMonitor defines the interface for tracking and managing autosuspend-enabled deployments.
// It provides methods to check the status of deployments, suspend them, and resume them when needed.
type DeploymentMonitor interface {

	// ResumeDeployment ensures the target deployment is scaled up to the desired number of replicas and ready.
	ResumeDeployment(ctx context.Context, serviceName string) error

	// SuspendIdleDeployments scans for idle deployments and scales them to zero.
	SuspendIdleDeployments(ctx context.Context) error

	// UpdateActivity registers the last activity time for a given service.
	UpdateActivity(ctx context.Context, serviceName string, timestamp time.Time) error
}
