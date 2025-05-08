package controller

import (
	"fmt"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	AnnotationIdleAfter    = "kubenap/idleAfter"
	AnnotationReplicaCount = "kubenap/replicaCount"
	AnnotationService      = "kubenap/service"
	AnnotationIngress      = "kubenap/ingress"
)

type ParsedConfig struct {
	IdleAfter    time.Duration
	ReplicaCount int32
	ServiceName  string
	IngressName  string
}

// ParseAnnotations extracts kubenap config from Deployment annotations.
func ParseAnnotations(dep *appsv1.Deployment) (*ParsedConfig, error) {
	ann := dep.Annotations
	if ann == nil {
		return nil, fmt.Errorf("missing annotations")
	}

	idleStr, ok := ann[AnnotationIdleAfter]
	if !ok {
		return nil, fmt.Errorf("missing annotation: %s", AnnotationIdleAfter)
	}

	idleAfter, err := time.ParseDuration(idleStr)
	if err != nil {
		return nil, fmt.Errorf("invalid idleAfter duration: %w", err)
	}

	repStr, ok := ann[AnnotationReplicaCount]
	if !ok {
		return nil, fmt.Errorf("missing annotation: %s", AnnotationReplicaCount)
	}

	replicas, err := strconv.Atoi(repStr)
	if err != nil {
		return nil, fmt.Errorf("invalid replicaCount: %w", err)
	}

	return &ParsedConfig{
		IdleAfter:    idleAfter,
		ReplicaCount: int32(replicas),
		ServiceName:  ann[AnnotationService],
		IngressName:  ann[AnnotationIngress],
	}, nil
}
