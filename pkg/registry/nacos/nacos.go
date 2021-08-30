package nacos

import (
	"bytes"
	"strings"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	nacosClient "github.com/dubbogo/gost/database/kv/nacos")

// nacos registry
type NacosRegistry struct {

	*common.URL
	namingClient *nacosClient.NacosNamingClient
	registryUrls []*common.URL

}

func (receiver *NacosRegistry) LoadAllServices() ([]*common.URL, error) {

	return nil, nil
}

func (receiver NacosRegistry) GetCluster() (string, error) {

	return "", nil
}

func appendParam(target *bytes.Buffer, url *common.URL, key string) {
	value := url.GetParam(key, "")
	target.Write([]byte(constant.NACOS_SERVICE_NAME_SEPARATOR))
	if strings.TrimSpace(value) != "" {
		target.Write([]byte(value))
	}
}






