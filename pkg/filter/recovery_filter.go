package filter

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func init() {
	extension.SetFilterFunc(constant.RecoveryFilter, Recover())
}

// Recover
func Recover() context.FilterFunc {
	return func(c context.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Infof("[dubboproxy go] error:%+v", err)
				c.WriteErr(err)
			}
		}()
		c.Next()
	}
}
