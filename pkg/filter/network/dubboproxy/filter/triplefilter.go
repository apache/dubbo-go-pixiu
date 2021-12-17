package filter

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	tpconst "github.com/dubbogo/triple/pkg/common/constant"
)

func Invoke( invocation protocol.Invocation) (protocol.Result, error) {
	tripleRefConf := newRefConf("org.apache.dubbo.samples.UserProviderTriple", tpconst.TRIPLE)

}

func newRefConf(iface, protocol string) config.ReferenceConfig {

	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		Cluster:       "failover",
		RegistryIDs:   []string{"zk"},
		Protocol:      protocol,
		URL: "127.0.0.1:20001",
		Generic:       "true",
	}

	rootConfig := config.NewRootConfigBuilder().
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad("pixiu")

	return refConf
}
