package sidecarprovider

import (
	"reflect"
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	dg "dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/registry"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func registerServiceInstance() {
	url := selectMetadataServiceExportedURL()
	if url == nil {
		return
	}
	instance, err := createInstance(url)
	if err != nil {
		panic(err)
	}
	p := extension.GetProtocol(constant.RegistryKey)
	var rp registry.RegistryFactory
	var ok bool
	if rp, ok = p.(registry.RegistryFactory); !ok {
		panic("dubbo registry protocol{" + reflect.TypeOf(p).String() + "} is invalid")
	}
	rs := rp.GetRegistries()
	for _, r := range rs {
		var sdr registry.ServiceDiscoveryHolder
		if sdr, ok = r.(registry.ServiceDiscoveryHolder); !ok {
			continue
		}
		// publish app level data to registry
		err := sdr.GetServiceDiscovery().Register(instance)
		if err != nil {
			panic(err)
		}
	}
	// publish metadata to remote
	if dg.GetApplicationConfig().MetadataType == constant.RemoteMetadataStorageType {
		if remoteMetadataService, err := extension.GetRemoteMetadataService(); err == nil && remoteMetadataService != nil {
			remoteMetadataService.PublishMetadata(dg.GetApplicationConfig().Name)
		}
	}
}

func createInstance(url *common.URL) (registry.ServiceInstance, error) {
	appConfig := dg.GetApplicationConfig()
	port, err := strconv.ParseInt(url.Port, 10, 32)
	if err != nil {
		return nil, perrors.WithMessage(err, "invalid port: "+url.Port)
	}

	host := url.Ip
	if len(host) == 0 {
		host = common.GetLocalIp()
	}

	// usually we will add more metadata
	metadata := make(map[string]string, 8)
	metadata[constant.MetadataStorageTypePropertyName] = appConfig.MetadataType

	instance := &registry.DefaultServiceInstance{
		ServiceName: appConfig.Name,
		Host:        host,
		Port:        int(port),
		ID:          host + constant.KeySeparator + url.Port,
		Enable:      true,
		Healthy:     true,
		Metadata:    metadata,
	}

	for _, cus := range extension.GetCustomizers() {
		cus.Customize(instance)
	}

	return instance, nil
}

// selectMetadataServiceExportedURL get already be exported url
func selectMetadataServiceExportedURL() *common.URL {
	var selectedUrl *common.URL
	metaDataService, err := extension.GetLocalMetadataService(constant.DefaultKey)
	if err != nil {
		logger.Warnf("get metadata service exporter failed, pls check if you import _ \"dubbo.apache.org/dubbo-go/v3/metadata/service/local\"")
		return nil
	}
	urlList, err := metaDataService.GetExportedURLs(constant.AnyValue, constant.AnyValue, constant.AnyValue, constant.AnyValue)
	if err != nil {
		panic(err)
	}
	if len(urlList) == 0 {
		return nil
	}
	for _, url := range urlList {
		selectedUrl = url
		// rest first
		if url.Protocol == "rest" {
			break
		}
	}
	return selectedUrl
}