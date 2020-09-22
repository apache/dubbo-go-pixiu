package api

import (
	"encoding/json"
	"errors"
)

import (
	"github.com/goinggo/mapstructure"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
)

func init() {
	extension.SetApiDiscoveryService(constant.LocalMemoryApiDiscoveryService, NewLocalMemoryApiDiscoveryService())
}

type LocalMemoryApiDiscoveryService struct {
}

func NewLocalMemoryApiDiscoveryService() *LocalMemoryApiDiscoveryService {
	return &LocalMemoryApiDiscoveryService{}
}

func (ads *LocalMemoryApiDiscoveryService) RemoveAllApi() error {
	model.CacheApi.Range(func(name, _ interface{}) bool {
		model.CacheApi.Delete(name)
		return true
	})
	return nil
}

func (ads *LocalMemoryApiDiscoveryService) RemoveApi(name string) error {
	model.CacheApi.Delete(name)
	return nil
}

func (ads *LocalMemoryApiDiscoveryService) AddApi(request service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	aj := model.NewApi()
	if err := json.Unmarshal(request.Body, aj); err != nil {
		return *service.EmptyDiscoveryResponse, err
	}

	if _, loaded := model.CacheApi.LoadOrStore(aj.Name, aj); loaded {
		// loaded
		logger.Warnf("api:%s is exist", aj)
	} else {
		// store
		if aj.Metadata == nil {

		} else {
			if v, ok := aj.Metadata.(map[string]interface{}); ok {
				if d, ok := v["dubbo"]; ok {
					dm := &dubbo.DubboMetadata{}
					if err := mapstructure.Decode(d, dm); err != nil {
						return *service.EmptyDiscoveryResponse, err
					}
					aj.Metadata = dm
				}
			}

			aj.RequestMethod = model.RequestMethod(model.RequestMethodValue[aj.Method])
		}
	}

	return *service.NewDiscoveryResponseWithSuccess(true), nil
}

func (ads *LocalMemoryApiDiscoveryService) GetApi(request service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	n := string(request.Body)

	if a, ok := model.CacheApi.Load(n); ok {
		return *service.NewDiscoveryResponse(a), nil
	}

	return *service.EmptyDiscoveryResponse, errors.New("not found")
}
