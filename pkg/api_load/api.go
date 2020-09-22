package api_load

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
	"sort"
	"sync"
	"time"
)

//var ApiLoadTypeMap = make(map[ApiLoadType]ApiLoader, 8)

type ApiLoadType string

const (
	File  ApiLoadType = "file"
	Nacos ApiLoadType = "nacos"
)

type ApiLoad struct {
	mergeLock       sync.Locker
	limiter         *time.Timer
	rateLimiterTime time.Duration
	mergeTask       chan struct{}
	ApiLoadTypeMap  map[ApiLoadType]ApiLoader
	ads             service.ApiDiscoveryService
}

func NewApiLoad(rateLimiterTime time.Duration, ads service.ApiDiscoveryService) *ApiLoad {
	if rateLimiterTime > time.Millisecond*50 {
		rateLimiterTime = time.Millisecond * 50
	}
	return &ApiLoad{
		ApiLoadTypeMap:  make(map[ApiLoadType]ApiLoader, 8),
		mergeTask:       make(chan struct{}, 1),
		limiter:         time.NewTimer(rateLimiterTime),
		rateLimiterTime: rateLimiterTime,
		ads:             ads,
	}
}

func (al *ApiLoad) AddApiLoad(config model.ApiConfig) {
	if config.File != nil {
		al.ApiLoadTypeMap[File] = NewFileApiLoader(WithFilePath(config.File.FileApiConfPath))
	}
	if config.Nacos != nil {
		al.ApiLoadTypeMap[Nacos] = NewNacosApiLoader(WithNacosAddress(config.Nacos.Address))
	}
}

func (al *ApiLoad) StartLoadApi() error {
	for _, loader := range al.ApiLoadTypeMap {
		err := loader.InitLoad()
		if err != nil {
			logger.Warn("proxy init api error:%v", err)
			break
		}
	}

	//al.MergeApi()

	if al.limiter == nil {
		return errors.New("proxy won't hot load api since limiter is null.")
	}

	for _, loader := range al.ApiLoadTypeMap {
		changeNotifier, err := loader.HotLoad()
		if err != nil {
			logger.Warn("proxy hot load api error:%v", err)
			break
		}

		go func() {
			for {
				select {
				case _, ok := <-changeNotifier:
					if !ok {
						logger.Debug("changeNotifier of apiloader was closed!")
						return
					}
					al.AddMergeTask()
					break
				}
			}
		}()
	}
	return nil
}

func (al *ApiLoad) AddMergeTask() error {
	select {
	case al.mergeTask <- struct{}{}:
		logger.Debug("added a merge task, waiting to merge api.")
		break
	case <-time.After(5 * time.Second):
		logger.Errorf("add merge task fail:wait timeout.")
		break
	}
	return nil
}

func (al *ApiLoad) SelectMergeApiTask() (err error) {
	for {
		select {
		case <-al.limiter.C:
			if len(al.mergeTask) > 0 {
				_, err = al.DoMergeApiTask()
				if err != nil {
					logger.Warnf("error merge api task:%v", err)
				}
			}
			al.limiter.Reset(time.Second)
			break
		default:
			time.Sleep(time.Millisecond * al.rateLimiterTime / 10)
			break
		}
	}
	return
}

func (al *ApiLoad) DoMergeApiTask() (skip bool, err error) {
	al.mergeLock.Lock()
	defer al.mergeLock.Unlock()
	wait := time.After(time.Millisecond * 50)
	select {
	case <-wait:
		logger.Debug("merge api task is too frequent.")
		skip = true
		return
	case <-al.mergeTask:
		// If apiLoadType is File,then try cover it's apis using other's apis from registry center
		var multiApisMerged map[string]model.Api
		var sortedApiLoader []int
		sortedApiLoaderMap := make(map[int]ApiLoadType, len(al.ApiLoadTypeMap))
		for apiLoadType, loader := range al.ApiLoadTypeMap {
			sortedApiLoader = append(sortedApiLoader, loader.GetPrior())
			sortedApiLoaderMap[loader.GetPrior()] = apiLoadType
		}

		sort.Ints(sortedApiLoader)
		for _, sortNo := range sortedApiLoader {
			loadType := sortedApiLoaderMap[sortNo]
			apiLoader := al.ApiLoadTypeMap[loadType]
			var apiConfigs []model.Api
			apiConfigs, err = apiLoader.GetLoadedApiConfigs()
			if err != nil {
				logger.Error("get file apis error:%v", err)
				return
			} else {
				for _, fleApiConfig := range apiConfigs {
					if fleApiConfig.Status != model.Up {
						continue
					}
					multiApisMerged[al.buildApiID(fleApiConfig)] = fleApiConfig
				}
			}
		}

		var totalApis []model.Api
		for _, api := range multiApisMerged {
			totalApis = append(totalApis, api)
		}
		err = al.ads.RemoveAllApi()
		if err != nil {
			logger.Errorf("remove all older apis error:%v", err)
			return
		}
		err = al.add2ApiDiscoveryService(totalApis)
		if err != nil {
			logger.Errorf("add newer apis error:%v", err)
			return
		}
		return
	}
}

func (al *ApiLoad) add2ApiDiscoveryService(apis []model.Api) error {
	for _, api := range apis {
		j, _ := json.Marshal(api)
		_, err := al.ads.AddApi(*service.NewDiscoveryRequest(j))
		if err != nil {
			logger.Errorf("error add api:%s", j)
			return err
		}
	}
	return nil
}

func (al *ApiLoad) buildApiID(api model.Api) string {
	return fmt.Sprintf("name:%s,ITypeStr:%s,OTypeStr:%s,Method:%s",
		api.Name, api.ITypeStr, api.OTypeStr)
}
