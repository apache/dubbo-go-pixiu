package zookeeper

import (
	"encoding/json"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/remoting/zookeeper"
	ex "github.com/apache/dubbo-go/common/extension"
	"github.com/apache/dubbo-go/metadata/definition"
	dr "github.com/apache/dubbo-go/registry"
	"github.com/apache/dubbo-go/remoting/zookeeper/curator_discovery"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/go-zookeeper/zk"
	"strings"
	"sync"
	"time"
)

var _ registry.Listener = new(applicationServiceListener)

// applicationServiceListener normally monitors the /services/[:application]/[:id]
type applicationServiceListener struct {
	urls   []interface{}
	path   string
	client *zookeeper.ZooKeeperClient

	exit chan struct{}
	wg   sync.WaitGroup
}

// newApplicationServiceListener creates a new zk service listener
func newApplicationServiceListener(path string, client *zookeeper.ZooKeeperClient) *applicationServiceListener {
	return &applicationServiceListener{
		path:   path,
		client: client,
		exit:   make(chan struct{}),
	}
}

func (asl *applicationServiceListener) WatchAndHandle() {
	defer asl.wg.Done()

	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	for {
		e, err := asl.client.ExistW(asl.path)
		// error handling
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", asl.path, err)
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching (path{%s}) got errNilNode,so exit listen", asl.path)
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on (path{%s}) exceed max fail times: %s,so exit listen",
					asl.path, MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		tickerTTL := defaultTTL
		ticker := time.NewTicker(tickerTTL)
		asl.handleEvent(asl.path)
	WATCH:
		for {
			select {
			case <-ticker.C:
				asl.handleEvent(asl.path)
			case zkEvent := <-e:
				logger.Warnf("get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
					zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, zookeeper.StateToString(zkEvent.State), zkEvent.Err)
				ticker.Stop()
				if zkEvent.Type != zk.EventNodeDeleted && zkEvent.Type != zk.EventNodeDataChanged {
					break WATCH
				}
				asl.handleEvent(zkEvent.Path)
				break WATCH
			case <-asl.exit:
				logger.Warnf("listen(path{%s}) goroutine exit now...", asl.path)
				ticker.Stop()
				return
			}
		}

	}
}

func (asl *applicationServiceListener) handleEvent(path string) {
	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	data, err := asl.client.GetContent(path)
	if err != nil {
		logger.Warnf("Get service instance %s failed due to %s", path, err.Error())
		// disable the API
		for _, url := range asl.urls {
			bkConf, _, _ := registry.ParseDubboString(url.(string))
			apiPattern := registry.GetAPIPattern(bkConf)
			if err := localAPIDiscSrv.RemoveAPIByPath(config.Resource{Path: apiPattern}); err != nil {
				logger.Errorf("Error={%s} when try to remove API by path: %s", err.Error(), apiPattern)
			}
		}
		return
	}

	asl.urls = asl.getUrls(data, path)
	for _, url := range asl.urls {
		bkConfig, _, err := registry.ParseDubboString(url.(string))
		if err != nil {
			logger.Warnf("Parse dubbo interface provider %s failed; due to %s", url.(string), err.Error())
			continue
		}
		if len(bkConfig.ApplicationName) == 0 || len(bkConfig.Interface) == 0 {
			continue
		}
		methods, err := asl.getMethods(bkConfig.Interface)
		if err != nil {
			logger.Warnf("Get methods of interface %s failed; due to %s", bkConfig.Interface, err.Error())
			continue
		}

		apiPattern := registry.GetAPIPattern(bkConfig)
		mappingParams := []config.MappingParam{
			{
				Name:  "requestBody.values",
				MapTo: "opt.values",
			},
			{
				Name:  "requestBody.types",
				MapTo: "opt.types",
			},
		}
		for i := range methods {
			api := registry.CreateAPIConfig(apiPattern, bkConfig, methods[i], mappingParams)
			if err := localAPIDiscSrv.AddAPI(api); err != nil {
				logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
			}
		}
	}
}

// getUrls get exported urls from instance
func (asl *applicationServiceListener) getUrls(data []byte, path string) []interface{} {
	iss := &curator_discovery.ServiceInstance{}
	err := json.Unmarshal(data, iss)
	if err != nil {
		logger.Warnf("Parse service instance %s failed due to %s", path, err.Error())
		return nil
	}
	instance := toZookeeperInstance(iss)

	metaData := instance.GetMetadata()
	metadataStorageType, ok := metaData[constant.MetadataStorageTypeKey]
	if !ok {
		metadataStorageType = constant.DefaultMetadataStorageType
	}
	// get metadata service proxy factory according to the metadataStorageType
	proxyFactory := ex.GetMetadataServiceProxyFactory(metadataStorageType)
	if proxyFactory == nil {
		return nil
	}
	metadataService := proxyFactory.GetProxy(instance)
	if metadataService == nil {
		logger.Warnf("Get metadataService of instance %s failed", instance)
		return nil
	}
	// call GetExportedURLs to get the exported urls
	urls, err := metadataService.GetExportedURLs(constant.AnyValue, constant.AnyValue, constant.AnyValue, constant.AnyValue)
	if err != nil {
		logger.Errorf("Get exported urls of instance %s failed; due to %s", instance, err.Error())
		return nil
	}
	return urls
}

// toZookeeperInstance convert to registry's service instance
func toZookeeperInstance(cris *curator_discovery.ServiceInstance) dr.ServiceInstance {
	pl, ok := cris.Payload.(map[string]interface{})
	if !ok {
		logger.Errorf("toZookeeperInstance{%s} payload is not map[string]interface{}", cris.Id)
		return nil
	}
	mdi, ok := pl["metadata"].(map[string]interface{})
	if !ok {
		logger.Errorf("toZookeeperInstance{%s} metadata is not map[string]interface{}", cris.Id)
		return nil
	}
	md := make(map[string]string, len(mdi))
	for k, v := range mdi {
		md[k] = fmt.Sprint(v)
	}
	return &dr.DefaultServiceInstance{
		Id:          cris.Id,
		ServiceName: cris.Name,
		Host:        cris.Address,
		Port:        cris.Port,
		Enable:      true,
		Healthy:     true,
		Metadata:    md,
	}
}

// getMethods will return the methods of a service
func (asl *applicationServiceListener) getMethods(in string) ([]string, error) {
	methods := []string{}

	path := strings.Join([]string{methodsRootPath, in}, constant.PathSlash)
	data, err := asl.client.GetContent(path)
	if err != nil {
		return nil, err
	}

	sd := &definition.ServiceDefinition{}
	err = json.Unmarshal(data, sd)
	if err != nil {
		return nil, err
	}

	for _, m := range sd.Methods {
		methods = append(methods, m.Name)
	}
	return methods, nil
}

// Close closes this listener
func (asl *applicationServiceListener) Close() {
	close(asl.exit)
	asl.wg.Wait()
}
