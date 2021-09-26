package server

import (
	"fmt"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

type (
	ApiConfigListener interface {
		OnAddAPI(r router.API) error
		OnDeleteRouter(r config.Resource) error
	}
	// ApiConfigManager similar to RouterManager
	ApiConfigManager struct {
		als map[string]ApiConfigListener
	}
)

func CreateDefaultApiConfigManager(server *Server, bs *model.Bootstrap) *ApiConfigManager {
	acm := &ApiConfigManager{als: make(map[string]ApiConfigListener)}
	return acm
}

func (acm *ApiConfigManager) AddApiConfigListener(adapterID string, l ApiConfigListener) {
	existedListener, existed := acm.als[adapterID]
	if existed {
		panic(fmt.Errorf("AddApiConfigListener %T and %T listen same adapter: %s", l, existedListener, adapterID))
	}
	acm.als[adapterID] = l
}

func (acm *ApiConfigManager) AddAPI(adapterID string, r router.API) error {
	l, existed := acm.als[adapterID]
	if !existed {
		return errors.Errorf("no listener found")
	}
	return l.OnAddAPI(r)
}

func (acm *ApiConfigManager) DeleteAPI(adapterID string, r config.Resource) error {
	l, existed := acm.als[adapterID]
	if !existed {
		return errors.Errorf("no listener found")
	}
	return l.OnDeleteRouter(r)
}
