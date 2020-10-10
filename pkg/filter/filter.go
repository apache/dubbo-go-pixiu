package filter

import "github.com/dubbogo/dubbo-go-proxy/pkg/context"

// Filter filter interface.
type Filter interface {
	// Run do filter, next() to next filter, before is pre logic, after is post logic.
	Run() context.FilterFunc
}
