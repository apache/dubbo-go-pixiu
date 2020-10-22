package pool

import (
	"errors"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/httpclient"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"sync"
)

//ClientPool  a pool of client.
type ClientPool struct {
	poolMap map[model.ApiType]*sync.Pool
}

var (
	_clinetPool *ClientPool
	once        = sync.Once{}
)

func newClientPool() *ClientPool {
	clientPool := &ClientPool{
		poolMap: make(map[model.ApiType]*sync.Pool),
	}
	clientPool.poolMap[model.DUBBO] = &sync.Pool{
		New: func() interface{} {
			return dubbo.NewDubboClient()
		},
	}
	clientPool.poolMap[model.REST] = &sync.Pool{
		New: func() interface{} {
			return httpclient.NewHttpClient()
		},
	}
	return clientPool
}

// SingletonPool singleton pool
func SingletonPool() *ClientPool {
	if _clinetPool == nil {
		once.Do(func() {
			_clinetPool = newClientPool()
		})
	}

	return _clinetPool
}

// GetClient  a factory method to get a client according to apiType .
func (pool *ClientPool) GetClient(t model.ApiType) (client.Client, error) {
	if pool.poolMap[t] != nil {
		return pool.poolMap[t].Get().(client.Client), nil
	}
	return nil, errors.New("protocol not supported yet")
}
