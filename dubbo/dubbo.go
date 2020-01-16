package dubbo

import (
	"context"
	"github.com/dubbogo/dubbo-go-proxy/common/config"
	"strings"
	"sync"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/logger"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)
import (
	"github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
)

var (
	Client GenericClientPool
	dgCfg  dg.ConsumerConfig
)

type GenericClientPool struct {
	mLock              sync.RWMutex
	GenericServicePool map[string]*dg.GenericService
}

func (d *GenericClientPool) Init() {
	dgCfg = dg.GetConsumerConfig()
	//can change config here
	dg.SetConsumerConfig(dgCfg)
	dg.Load()
	d.GenericServicePool = make(map[string]*dg.GenericService)
}

func (d *GenericClientPool) get(appName string) *dg.GenericService {
	d.mLock.RLock()
	defer d.mLock.RUnlock()
	return d.GenericServicePool[appName]
}
func (d *GenericClientPool) check(appName string) bool {
	d.mLock.RLock()
	defer d.mLock.RUnlock()
	if _, ok := d.GenericServicePool[appName]; ok {
		return true
	} else {
		return false
	}
}
func (d *GenericClientPool) create(interfaceName, version, group string) *dg.GenericService {
	key := strings.Join([]string{interfaceName, version, group}, "_")
	referenceConfig := dg.NewReferenceConfig(interfaceName, context.TODO())
	referenceConfig.InterfaceName = interfaceName
	referenceConfig.Cluster = constant.DEFAULT_CLUSTER
	var registers []string
	for k := range dgCfg.Registries {
		registers = append(registers, k)
	}
	referenceConfig.Registry = strings.Join(registers, ",")
	referenceConfig.Protocol = dubbo.DUBBO
	referenceConfig.Version = version
	referenceConfig.Generic = true
	referenceConfig.Retries = config.Config.Retries
	d.mLock.Lock()
	defer d.mLock.Unlock()
	referenceConfig.GenericLoad(interfaceName)
	time.Sleep(200 * time.Millisecond) //sleep to wait invoker create
	clientService := referenceConfig.GetRPCService().(*dg.GenericService)

	d.GenericServicePool[key] = clientService
	return clientService
}

/* todo 预注册加载配置，能够提前和zk，避免接口首次请求时时间慢*/

func (d *GenericClientPool) Get(interfaceName, version, group string) *dg.GenericService {
	key := strings.Join([]string{interfaceName, version, group}, "_")
	if d.check(key) {
		return d.get(key)
	} else {
		return d.create(interfaceName, version, group)
	}
}
