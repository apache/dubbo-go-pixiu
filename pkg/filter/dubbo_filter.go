package filter

import (
	"io/ioutil"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
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
	extension.SetFilterFunc(constant.HttpTransferDubboFilter, HttpDubbo())
}

// HttpDubbo http 2 dubbo
func HttpDubbo() context.FilterFunc {
	return func(c context.Context) {
		doDubbo(c.(*http.HttpContext))
	}
}

func doDubbo(c *http.HttpContext) {
	api := c.GetApi()

	if bytes, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.Errorf("[dubboproxy go] read body err:%v!", err)
		c.WriteFail()
		c.Abort()
	} else {
		if resp, err := dubbo.SingleDubboClient().Call(client.NewRequest(bytes, api)); err != nil {
			logger.Errorf("[dubboproxy go] client do err:%v!", err)
			c.WriteFail()
			c.Abort()
		} else {
			c.WriteResponse(resp)
			c.Next()
		}
	}
}
