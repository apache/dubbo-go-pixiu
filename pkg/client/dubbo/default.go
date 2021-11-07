package dubbo

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
)

var defaultMappingParams []config.MappingParam

// getMappingParams If configuration items are set in the api_config configuration file, then the configuration in api_config will be used first,
// if not, the default will be used, and a total rule for obtaining the original data will be generated at the end
func getMappingParams(mps []config.MappingParam) []config.MappingParam {
	newMappingParam := make([]config.MappingParam, len(mps))
	copy(newMappingParam, mps)
	for _, defaultMappingParam := range defaultMappingParams {
		equal := false
		for _, customMappingParam := range mps {
			if customMappingParam.Name == defaultMappingParam.Name {
				equal = true
			}
		}
		if equal {
			continue
		}
		newMappingParam = append(newMappingParam, defaultMappingParam)
	}
	return newMappingParam
}
