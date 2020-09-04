package filter

import (
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func init() {
	extension.SetFilterFunc(constant.LoggerFilter, Logger())
}

func Logger() context.FilterFunc {
	return func(c context.Context) {
		start := time.Now()

		c.Next()

		latency := time.Now().Sub(start)

		logger.Infof("[dubboproxy go] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.StatusCode(), latency, c.GetMethod(), c.GetUrl())
	}
}
