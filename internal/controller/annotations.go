package controller

import (
	"log"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	AnnotationIdleAfter    = "kubenap/idleAfter"
	AnnotationReplicaCount = "kubenap/replicaCount"
	AnnotationService      = "kubenap/service"
	AnnotationIngress      = "kubenap/ingress"
	DefaultIdleAfter       = 10 * time.Minute
	DefaultReplicaCount    = 1
)

type ParsedConfig struct {
	IdleAfter    time.Duration
	ReplicaCount int32
	ServiceName  string
	IngressName  string
}

// ParseAnnotations extracts kubenap config from Deployment annotations,
// falling back to sensible defaults where necessary.
func ParseAnnotations(dep *appsv1.Deployment) (*ParsedConfig, error) {
	ann := dep.Annotations
	if ann == nil {
		log.Printf("No annotations found on deployment %s/%s, using defaults", dep.Namespace, dep.Name)
		return &ParsedConfig{
			IdleAfter:    DefaultIdleAfter,
			ReplicaCount: DefaultReplicaCount,
		}, nil
	}

	// idleAfter
	idleAfter := DefaultIdleAfter
	if idleStr, ok := ann[AnnotationIdleAfter]; ok {
		parsed, err := time.ParseDuration(idleStr)
		if err != nil {
			log.Printf("Invalid idleAfter format on %s/%s: %v, using default %s",
				dep.Namespace, dep.Name, err, DefaultIdleAfter)
		} else {
			idleAfter = parsed
		}
	} else {
		log.Printf("idleAfter not set on %s/%s, using default %s", dep.Namespace, dep.Name, DefaultIdleAfter)
	}

	// replicaCount
	replicaCount := DefaultReplicaCount
	if repStr, ok := ann[AnnotationReplicaCount]; ok {
		replicas, err := strconv.Atoi(repStr)
		if err != nil {
			log.Printf("Invalid replicaCount on %s/%s: %v, using default %d",
				dep.Namespace, dep.Name, err, DefaultReplicaCount)
		} else {
			replicaCount = replicas
		}
	} else {
		log.Printf("replicaCount not set on %s/%s, using default %d", dep.Namespace, dep.Name, DefaultReplicaCount)
	}

	return &ParsedConfig{
		IdleAfter:    idleAfter,
		ReplicaCount: int32(replicaCount),
		ServiceName:  ann[AnnotationService],
		IngressName:  ann[AnnotationIngress],
	}, nil
}
