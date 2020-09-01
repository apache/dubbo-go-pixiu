package proxy

import (
	"io/ioutil"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

func init() {
	AddFilterFunc(pkg.HttpTransferDubboFilter, HttpDubbo())
}

func HttpDubbo() FilterFunc {
	return func(c Context) {
		doDubbo(c.(*HttpContext))
	}
}

func doDubbo(c *HttpContext) {
	api := c.GetApi()

	if bytes, err := ioutil.ReadAll(c.r.Body); err != nil {
		logger.Errorf("[dubboproxy go] read body err:%v!", err)
		c.WriteFail()
		c.Abort()
	} else {
		if api.Client == nil {
			api.Client = SingleDubboClient()
		}

		if resp, err := api.Client.Call(NewRequest(bytes, api)); err != nil {
			logger.Errorf("[dubboproxy go] client do err:%v!", err)
			c.WriteFail()
			c.Abort()
		} else {
			c.WriteResponse(resp)
			c.Next()
		}
	}
}
