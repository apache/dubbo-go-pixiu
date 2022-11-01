package dubbogen

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	v3 "github.com/apache/dubbo-go-pixiu/pilot/pkg/xds/v3"
	istiolog "istio.io/pkg/log"
)

var log = istiolog.RegisterScope("dubbogen", "xDS Generator for Proxyless dubbo", 0)

type DubboConfigGenerator struct{}

func (d *DubboConfigGenerator) Generate(proxy *model.Proxy, w *model.WatchedResource, req *model.PushRequest) (model.Resources, model.XdsLogDetails, error) {
	switch w.TypeUrl {
	case v3.DubboServiceNameMappingType:
		return d.buildServiceNameMappings(proxy, req, w.ResourceNames), model.DefaultXdsLogDetails, nil
	default:
		log.Warnf("not support now %s", w.TypeUrl)
		return nil, model.DefaultXdsLogDetails, nil
	}
}
