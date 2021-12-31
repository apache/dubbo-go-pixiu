package test

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	hessian "github.com/apache/dubbo-go-hessian2"
	tpconst "github.com/dubbogo/triple/pkg/common/constant"
	"testing"
)

func TestTriple2Dubbo(t *testing.T) {
	tripleRefConf := newTripleRefConf("com.dubbogo.pixiu.UserService", tpconst.TRIPLE)
	resp, err := tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser1",
		[]string{"java.lang.String"},
		[]hessian.Object{"A003"},
	)

	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser1(userId string) res: %+v", resp)
}

func newTripleRefConf(iface, protocol string) config.ReferenceConfig {

	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		Cluster:       "failover",
		RegistryIDs:   []string{"zk"},
		Protocol:      protocol,
		Generic:       "true",
		URL:           "tri://127.0.0.1:9999/" + iface + "?" + constant.SerializationKey + "=hessian2",
	}

	rootConfig := config.NewRootConfigBuilder().
		Build()
	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad(appName)

	return refConf
}
