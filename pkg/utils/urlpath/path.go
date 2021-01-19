package urlpath

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"strings"
)

func Split(path string) []string {
	return strings.Split(strings.TrimLeft(path, constant.PathSlash), constant.PathSlash)
}

func VariableName(key string) string {
	return strings.TrimPrefix(key, constant.PathParamIdentifier)
}

func IsPathVariable(key string) bool {
	if key == "" {
		return false
	}

	if strings.HasPrefix(key, constant.PathParamIdentifier) {
		return true
	}
	//return key[0] == '{' && key[len(key)-1] == '}'
	return false
}
