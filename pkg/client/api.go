package client

import (
	"sync"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

var (
	CacheApi = sync.Map{}
)

// Api is api gateway concept, control request from browser、Mobile APP、third party people
type Api struct {
	Name     string        `json:"name"`
	ITypeStr string        `json:"itype"`
	IType    model.ApiType `json:"-"`
	OTypeStr string        `json:"otype"`
	OType    model.ApiType `json:"-"`
	Status   model.Status  `json:"status"`
	Metadata interface{}   `json:"metadata"`
	Method   string        `json:"method"`
	model.RequestMethod
	Client Client
	lock   sync.Mutex
}

var EmptyApi = &Api{}

func NewApi() *Api {
	return &Api{lock: sync.Mutex{}}
}

func (a *Api) FindApi(name string) (*Api, bool) {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api), true
	}

	return nil, false
}

func (a *Api) MatchMethod(method string) bool {
	i := model.RequestMethodValue[method]
	if a.RequestMethod == model.RequestMethod(i) {
		return true
	}

	return false
}

func (a *Api) IsOk(name string) bool {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api).Status == model.Up
	}

	return false
}

// Offline api offline
func (a *Api) Offline(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*Api).Status = model.Down
	}
}

// Online api online
func (a *Api) Online(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*Api).Status = model.Up
	}
}
