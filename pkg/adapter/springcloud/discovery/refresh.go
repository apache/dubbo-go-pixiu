package discovery

import (
	"errors"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"sync"
	"time"
)
const REGISTRY_REFRESH_INTERVAL = 30

type periodicalRefreshClient struct {
	client Client
}

func newPeriodicalRefresh(c Client) *periodicalRefreshClient {
	return &periodicalRefreshClient{c}
}

// 向eureka查询注册列表, 刷新本地列表
func (r *periodicalRefreshClient) StartPeriodicalRefresh() error {
	logger.Infof("refresh registry every %d sec", REGISTRY_REFRESH_INTERVAL)

	refreshRegistryChan := make(chan error)

	isBootstrap := true
	go func() {
		ticker := time.NewTicker(REGISTRY_REFRESH_INTERVAL * time.Second)

		for {
			logger.Info("registry refresh started")
			err := r.doRefresh()
			if nil != err {
				// 如果是第一次查询失败, 退出程序
				if isBootstrap {
					logger.Errorf("failed to refresh registry", err)
					refreshRegistryChan <- err
					return

				} else {
					logger.Error(err)
				}

			}
			logger.Info("done refreshing registry")

			if isBootstrap {
				isBootstrap = false
				close(refreshRegistryChan)
			}

			<-ticker.C
		}
	}()

	return <- refreshRegistryChan
}

func (r *periodicalRefreshClient) doRefresh() error {
	instances, err := r.client.QueryServices()

	if nil != err {
		logger.Errorf("failed to query all services", err)
		return err
	}

	if nil == instances {
		logger.Info("no instance found")
		return nil
	}

	logger.Infof("total app count: %d", len(instances))

	newRegistryMap := r.groupByService(instances)

	r.refreshRegistryMap(newRegistryMap)

	return nil

}


// 将所有实例按服务名进行分组
func (r *periodicalRefreshClient) groupByService(instances []*InstanceInfo) *sync.Map {
	servMap := new(sync.Map)
	for _, ins := range instances {
		infosGeneric, exist := servMap.Load(ins.ServiceName)
		if !exist {
			infosGeneric = make([]*InstanceInfo, 0, 5)
			infosGeneric = append(infosGeneric.([]*InstanceInfo), ins)

		} else {
			infosGeneric = append(infosGeneric.([]*InstanceInfo), ins)
		}
		servMap.Store(ins.ServiceName, infosGeneric)
	}
	return servMap
}


// 更新本地注册列表
// s: gogate server对象
// newRegistry: 刚从eureka查出的最新服务列表
func (r *periodicalRefreshClient) refreshRegistryMap(newRegistry *sync.Map) {
	if nil == r.client.GetInternalRegistryStore() {
		r.client.SetInternalRegistryStore(NewInsInfoArrSyncMap())
	}

	// 找出本地列表存在, 但新列表中不存在的服务
	exclusiveKeys, _ := FindExclusiveKey(r.client.GetInternalRegistryStore().GetMap(), newRegistry)
	// 删除本地多余的服务
	DelKeys(r.client.GetInternalRegistryStore().GetMap(), exclusiveKeys)
	// 将新列表中的服务合并到本地列表中
	MergeSyncMap(newRegistry, r.client.GetInternalRegistryStore().GetMap())
}

/*
* 从map中删除指定的key

* PARAMS:
*	- baseMap: 要删除key的map
*	- keys: 要删除的key数组
 */
func DelKeys(baseMap *sync.Map, keys []interface{}) error {
	if nil == baseMap {
		return errors.New("baseMap cannot be null")
	}

	for _, key := range keys {
		baseMap.Delete(key)
	}

	return nil
}

/*
* 两个map取并集
*
* PARAMS:
*	- fromMap: 源map
*	- toMap: 合并后的map
*
 */
func MergeSyncMap(fromMap, toMap *sync.Map) error {
	if nil == fromMap || nil == toMap {
		return errors.New("fromMap or toMap cannot be null")
	}

	fromMap.Range(func(key, value interface{}) bool {
		toMap.Store(key, value)
		return true
	})

	return nil
}

/*
* 找出在baseMap中存在但yMap中不存在的元素
*
* PARAMS:
*	- baseMap: 独有元素所在的map
*	- yMap: 对比map
*
* RETURNS:
*	baseMap中独有元素的key的数组
 */
func FindExclusiveKey(baseMap, yMap *sync.Map) ([]interface{}, error) {
	if nil == baseMap || nil == yMap {
		return nil, errors.New("fromMap or toMap cannot be null")
	}

	var keys []interface{}
	baseMap.Range(func(key, value interface{}) bool {
		_, exist := yMap.Load(key)
		if !exist {
			keys = append(keys, key)
		}

		return true
	})

	return keys, nil
}