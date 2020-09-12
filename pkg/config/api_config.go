package config

import (
	"log"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	perrors "github.com/pkg/errors"
)

// import (
// 	"sync"
// )

// var (
// 	CacheApi = sync.Map{}
// )

// // Api is api gateway concept, control request from browser、Mobile APP、third party people
// type Api struct {
// 	Name     string      `json:"name"`
// 	ITypeStr string      `json:"itype"`
// 	IType    ApiType     `json:"-"`
// 	OTypeStr string      `json:"otype"`
// 	OType    ApiType     `json:"-"`
// 	Status   Status      `json:"status"`
// 	Metadata interface{} `json:"metadata"`
// 	Method   string      `json:"method"`
// 	RequestMethod
// }

// var EmptyApi = &Api{}

// func NewApi() *Api {
// 	return &Api{}
// }

// func (a *Api) FindApi(name string) (*Api, bool) {
// 	if v, ok := CacheApi.Load(name); ok {
// 		return v.(*Api), true
// 	}

// 	return nil, false
// }

// func (a *Api) MatchMethod(method string) bool {
// 	i := RequestMethodValue[method]
// 	if a.RequestMethod == RequestMethod(i) {
// 		return true
// 	}

// 	return false
// }

// func (a *Api) IsOk(name string) bool {
// 	if v, ok := CacheApi.Load(name); ok {
// 		return v.(*Api).Status == Up
// 	}

// 	return false
// }

// // Offline api offline
// func (a *Api) Offline(name string) {
// 	if v, ok := CacheApi.Load(name); ok {
// 		v.(*Api).Status = Down
// 	}
// }

// // Online api online
// func (a *Api) Online(name string) {
// 	if v, ok := CacheApi.Load(name); ok {
// 		v.(*Api).Status = Up
// 	}
// }

// HTTPVerb defines the restful api http verb
type HTTPVerb string

const (
	// MethodAny any method
	MethodAny HTTPVerb = "ANY"
	// MethodGet get
	MethodGet HTTPVerb = "GET"
	// MethodHead head
	MethodHead HTTPVerb = "HEAD"
	// MethodPost post
	MethodPost HTTPVerb = "POST"
	// MethodPut put
	MethodPut HTTPVerb = "PUT"
	// MethodPatch patch
	MethodPatch HTTPVerb = "PATCH" // RFC 5789
	// MethodDelete delete
	MethodDelete HTTPVerb = "DELETE"
	// MethodOptions options
	MethodOptions HTTPVerb = "OPTIONS"
)

// APIConfig defines the data structure of the api gateway configuration
type APIConfig struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Resources   []Resource   `yaml:"resources"`
	Definitions []Definition `yaml:"definitions"`
}

// Resource defines the API path
type Resource struct {
	Type        string     `yaml:"type"` // Restful, Dubbo
	Path        string     `yaml:"path"`
	Description string     `yaml:"description"`
	Filters     []string   `yaml:"filters"`
	Methods     []Method   `yaml:"methods"`
	Resources   []Resource `yaml:"resources,omitempty"`
}

// Method defines the method of the api
type Method struct {
	OnAir          bool `yaml:"onAir"`
	HTTPVerb       `yaml:"httpVerb"`
	InboundRequest `yaml:"inboundRequest"`
}

// InboundRequest defines the details of the inbound
type InboundRequest struct {
	Type         string           `yaml:"type"` //http, TO-DO: dubbo
	Headers      []Params         `yaml:"headers"`
	QueryStrings []Params         `yaml:"queryStrings"`
	RequestBody  []BodyDefinition `yaml:"requestBody"`
}

// Params defines the simple parameter definition
type Params struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
}

// BodyDefinition connects the request body to the definitions
type BodyDefinition struct {
	DefinitionName string `yaml:"definitionName"`
	ContentType    string `yaml:"contentType"`
}

// IntegrationRequest defines the backend request format and target
type IntegrationRequest struct {
	Type                string `yaml:"type"` // dubbo, TO-DO: http
	dubbo.DubboMetadata `yaml:"dubboMetaData,inline"`
}

// MappingParams defines the mapping rules of headers and queryStrings
type MappingParams struct {
	Name  string `yaml:"name"`
	MapTo string `yaml:"mapTo"`
}

// Definition defines the complex json request body
type Definition struct {
	Name   string `yaml:"name"`
	Schema string `yaml:"schema"` // use json schema
}

// LoadAPIConfigFromFile load the api config from file
func LoadAPIConfigFromFile(path string) (*APIConfig, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	log.Printf("Load API configuration file form %s", path)
	apiConf := &APIConfig{}
	_, err := yaml.UnmarshalYMLConfig(path, apiConf)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	return apiConf, nil
}
