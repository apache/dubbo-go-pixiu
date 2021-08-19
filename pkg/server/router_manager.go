package server

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	RouterListener interface {
		OnAddRouter(r *model.Router)
		OnDeleteRouter(r *model.Router)
	}

	RouterManager struct {
		rls []RouterListener
	}
)

func CreateDefaultRouterManager(server *Server, bs *model.Bootstrap) *RouterManager {
	rm := &RouterManager{}
	return rm
}

func (rm *RouterManager) AddRouterListener(l RouterListener) {
	rm.rls = append(rm.rls, l)
}

func (rm *RouterManager) AddRouter(r *model.Router) {
	for _, l := range rm.rls {
		l.OnAddRouter(r)
	}
}

func (rm *RouterManager) DeleteRouter(r *model.Router) {
	for _, l := range rm.rls {
		l.OnDeleteRouter(r)
	}
}
