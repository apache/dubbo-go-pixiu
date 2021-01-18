package urlPath

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
