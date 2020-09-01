package proxy

import "github.com/dubbogo/dubbo-go-proxy/pkg/logger"

func Recover() FilterFunc {
	return func(c Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Infof("[dubboproxy go] error:%+v", err)
				c.WriteErr(err)
			}
		}()
		c.Next()
	}
}
