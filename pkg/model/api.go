package model

import (
	"sync"
)

var (
	CacheApi = sync.Map{}
)

// Api is api gateway concept, control request from browser、Mobile APP、third party people
type Api struct {
	Name     string      `json:"name"`
	ITypeStr string      `json:"itype"`
	IType    ApiType     `json:"-"`
	OTypeStr string      `json:"otype"`
	OType    ApiType     `json:"-"`
	Status   Status      `json:"status"`
	Metadata interface{} `json:"metadata"`
	Method   string      `json:"method"`
	RequestMethod
}

var EmptyApi = &Api{}

func NewApi() *Api {
	return &Api{}
}

func (a *Api) FindApi(name string) (*Api, bool) {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api), true
	}

	return nil, false
}

func (a *Api) MatchMethod(method string) bool {
	i := RequestMethodValue[method]
	if a.RequestMethod == RequestMethod(i) {
		return true
	}

	return false
}

func (a *Api) IsOk(name string) bool {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api).Status == Up
	}

	return false
}

// Offline api offline
func (a *Api) Offline(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*Api).Status = Down
	}
}

// Online api online
func (a *Api) Online(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*Api).Status = Up
	}
}
