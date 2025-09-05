package registry

import (
	"context"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Controller is a provider-agnostic control-plane interface for MCP registries.
// It discovers, watches and forwards changes via an onChange callback.
type Controller interface {
	Run(ctx context.Context, interval time.Duration) error
	Close() error
}

// BuildFunc creates a Controller for a given registry configuration.
// Implemented by each provider and registered via RegisterProvider.
type BuildFunc func(reg model.Registry, onChange func(*model.McpServerConfig)) (Controller, error)
