package proxy

import (
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func Logger() FilterFunc {
	return func(c Context) {
		start := time.Now()

		c.Next()

		latency := time.Now().Sub(start)

		logger.Infof("[dubboproxy go] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.StatusCode(), latency, c.GetMethod(), c.GetUrl())
	}
}
