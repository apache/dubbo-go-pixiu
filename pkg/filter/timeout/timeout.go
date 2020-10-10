package timeout

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"time"
)

func init() {
	extension.SetFilterFunc(constant.TimeoutFilter, NewTimeoutFilter().Run())
}

// timeoutFilter is a filter for control request time out.
type timeoutFilter struct {
	waitTime time.Duration
}

// NewTimeoutFilter create filter.
func NewTimeoutFilter() *timeoutFilter {
	return &timeoutFilter{
		waitTime: time.Minute,
	}
}

// DoFilter do filter.
func (f *timeoutFilter) Run() context.FilterFunc {
	return func(c context.Context) {

	}
}
