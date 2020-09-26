package client

import (
	"errors"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/httpclient"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	_ "github.com/gin-gonic/gin"
	"sync"
)

type Client interface {
	Init() error
	Close() error
	Call(req *Request) (resp Response, err error)
}

type ClientPool struct {
	poolMap map[model.ApiType]*sync.Pool
}

func (pool *ClientPool) getClient(t model.ApiType) (Client, error) {
	if pool.poolMap[t] != nil {
		return pool.poolMap[t].Get().(Client), nil
	}
	return nil, errors.New("protocol not supported yet")
}

func NewClientPool() *ClientPool {
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
