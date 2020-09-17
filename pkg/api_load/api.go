package api_load

import (
	"fmt"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
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
	mergeLock      sync.Locker
	limiter        *time.Ticker
	sleepTime      int64
	mergeTask      chan struct{}
	ApiLoadTypeMap map[ApiLoadType]ApiLoader
}

func NewApiLoad() *ApiLoad {
	return &ApiLoad{
		ApiLoadTypeMap: make(map[ApiLoadType]ApiLoader, 8),
		mergeTask:      make(chan struct{}, 1),
	}
}

func (al *ApiLoad) AddApiLoad(fileApiConfPath string, config model.ApiConfig) {
	if len(fileApiConfPath) > 0 {
		al.ApiLoadTypeMap[File] = NewFileApiLoader(WithFilePath(fileApiConfPath))
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
		logger.Warnf("proxy won't hot load api since limiter is null.")
		return nil
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
				case <-changeNotifier:
					al.SelectMergeApiTask()
					break
				}
			}
		}()
	}
	return nil
}

func (al *ApiLoad) ClearMergeTask() {
	wait := time.After(time.Millisecond * 50)
	for {
		select {
		case <-al.mergeTask:
			logger.Debug("drop merge task")
			break
		case <-wait:
			return
		}
	}
}

func (al *ApiLoad) SelectMergeApiTask() (err error) {
	for {
		select {
		case <-al.limiter.C:
			if len(al.mergeTask) > 0 {
				go al.DoMergeApiTask()
			}
			break
		default:
			al.ClearMergeTask()
			time.Sleep(time.Millisecond * time.Duration(al.sleepTime))
			break
		}
	}
	return
}

func (al *ApiLoad) DoMergeApiTask() {
	al.mergeLock.Lock()
	defer al.mergeLock.Unlock()
	wait := time.After(time.Millisecond * 50)
	select {
	case <-wait:
		break
	case <-al.mergeTask:
		// If apiLoadType is File,then try cover it's apis using other's apis from registry center
		var totalApis map[string]model.Api
		sortedApiLoader := []int{}
		sortedApiLoaderMap := make(map[int]ApiLoadType, len(al.ApiLoadTypeMap))
		for apiLoadType, loader := range al.ApiLoadTypeMap {
			sortedApiLoader = append(sortedApiLoader, loader.GetPrior())
			sortedApiLoaderMap[loader.GetPrior()] = apiLoadType
		}

		sort.Ints(sortedApiLoader)
		for _, sortNo := range sortedApiLoader {
			loadType := sortedApiLoaderMap[sortNo]
			apiLoader := al.ApiLoadTypeMap[loadType]
			fileApiConfigs, err := apiLoader.GetLoadedApiConfigs()
			if err != nil {
				logger.Error("get file apis error:%v", err)
			} else {
				for _, fleApiConfig := range fileApiConfigs {
					if fleApiConfig.Status != model.Up {
						continue
					}
					totalApis[al.buildApiID(fleApiConfig)] = fleApiConfig
				}
			}
		}
		// todo添加api
		break
	}
}

func (al *ApiLoad) buildApiID(api model.Api) string {
	return fmt.Sprintf("name:%s,ITypeStr:%s,OTypeStr:%s,Method:%s",
		api.Name, api.ITypeStr, api.OTypeStr)
}
