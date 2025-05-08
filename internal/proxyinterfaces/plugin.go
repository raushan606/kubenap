package proxyinterfaces

import (
	"context"
	"net/http"
)

// ProxyPlugin allows proxy-specific logic for interacting with the kubenap controller.
type ProxyPlugin interface {

	// ShouldWake determines whether the given HTTP request should trigger a wake action.
	ShouldWake(r *http.Request) (bool, error)

	// ExtractOriginalPath extracts the original request path to resume the correct application.
	ExtractOriginalPath(r *http.Request) (string, error)

	// OnResumeCompleted is called after the app has been resumed and ready to serve.
	// The plugin can optionally perform logging or metric forwarding.
	OnResumeCompleted(ctx context.Context, serviceName string) error
}
