package healthchecker

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/pkg/errors"
)

type (
	// HealthCheckerFactory describe health checker factory
	HealthCheckerFactory interface {
		CreateHealthChecker(cfg map[string]interface{}, endpoint model.Endpoint) HealthChecker
	}

	// HealthChecker upstream cluster health checker
	HealthChecker interface {
		// CheckHealth check health
		CheckHealth() bool
	}
)

var (
	healthCheckerFactoryRegistry = map[string]HealthCheckerFactory{}
)

// Register registers health checker factory.
func RegisterHealthCheckerFactory(name string, f HealthCheckerFactory) {
	if name == "" {
		panic(fmt.Errorf("%T: empty name", f))
	}

	existedFactory, existed := healthCheckerFactoryRegistry[name]
	if existed {
		panic(fmt.Errorf("%T and %T got same factory: %s", f, existedFactory, name))
	}

	healthCheckerFactoryRegistry[name] = f
}

// GetHttpHealthCheckerFactory get factory by kind
func GetHttpHealthCheckerFactory(name string) (HealthCheckerFactory, error) {
	existedFilter, existed := healthCheckerFactoryRegistry[name]
	if existed {
		return existedFilter, nil
	}
	return nil, errors.Errorf("factory not found %s", name)
}

// CreateHealthCheck is a extendable function that can create different health checker
// by different health check session.
// The Default session is TCPDial session
func CreateHealthCheck(cfg model.HealthCheckConfig) HealthChecker {
	f, ok := GetHttpHealthCheckerFactory(cfg.Protocol)
	if !ok {
		// not registered, use default session factory
		f = &TCPDialSessionFactory{}
	}
	return newHealthChecker(cfg, f)
}
