package model

import (
	"sync"
)

var (
	CacheApi = sync.Map{}
	EmptyApi = &Api{}
)

// Api is api gateway concept, control request from browser、Mobile APP、third party people
type Api struct {
	Name     string      `json:"name" yaml:"name"`
	ITypeStr string      `json:"itype" yaml:"itype"`
	IType    ApiType     `json:"-" yaml:"-"`
	OTypeStr string      `json:"otype" yaml:"otype"`
	OType    ApiType     `json:"-" yaml:"-"`
	Status   Status      `json:"status" yaml:"status"`
	Metadata interface{} `json:"metadata" yaml:"metadata"`
	Method   string      `json:"method" yaml:"method"`
	RequestMethod
}

// NewApi
func NewApi() *Api {
	return &Api{}
}

// FindApi find a api, if not exist, return false
func (a *Api) FindApi(name string) (*Api, bool) {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api), true
	}

	return nil, false
}

// MatchMethod
func (a *Api) MatchMethod(method string) bool {
	i := RequestMethodValue[method]
	if a.RequestMethod == RequestMethod(i) {
		return true
	}

	return false
}

// IsOk api status equals Up
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
