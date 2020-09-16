package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

var ApiLoadTypeMap = make(map[ApiLoadType]ApiLoad, 8)

type ApiLoadType string

const (
	File  ApiLoadType = "file"
	Nacos ApiLoadType = "nacos"
)

func AddApiLoad(fileApiConfPath string, config model.ApiConfig) {
	if len(fileApiConfPath) > 0 {
		ApiLoadTypeMap[File] = NewFileApiLoader(WithFilePath(fileApiConfPath))
	}
	if config.Nacos != nil {
		ApiLoadTypeMap[Nacos] = NewNacosApiLoader(WithNacosAddress(config.Nacos.Address))
	}
}

func MergeApi() {

}
