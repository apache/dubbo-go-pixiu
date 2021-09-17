package nacos

import (
	"bytes"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	nacosClient "github.com/dubbogo/gost/database/kv/nacos"
)

func init() {
	var _ registry.Loader = new(NacosRegistry)
}

// nacos registry
type NacosRegistry struct {
	url          []*common.URL
	namingClient *nacosClient.NacosNamingClient
	registryUrls []*common.URL
}

const (
	NacosDefaultContext = "/nacos"
)

func (receiver *NacosRegistry) NewRegistryLoader(ad *model.Adapter) (registry.Loader, error) {
	return newNacosRegistryLoad(ad), nil
}

func newNacosRegistryLoad(ad *model.Adapter) registry.Loader {
	if ad == nil {
		return nil
	}
	sc, cc, err := GetNacosConfig(ad)
	if err != nil {
		logger.Errorf("nameClient failed")
	}
	var tmpClient *NacosRegistry
	//for _, adapter := range sr.Adapters {
	application := ad.Name
	Client, err := nacosClient.NewNacosNamingClient(application, false, sc, cc)
	if err != nil {
		logger.Info("fail to init Nacos NamingClient")
	}
	tmpClient = &NacosRegistry{
		namingClient: Client,
		registryUrls: []*common.URL{},
	}
	//}
	return tmpClient
}

func (receiver *NacosRegistry) NewResistry(ad *model.Adapter) registry.Loader {
	return newNacosRegistryLoad(ad)
}

func (receiver *NacosRegistry) LoadAllServices() ([]*common.URL, error) {
	namingClient := receiver.namingClient.Client()
	serviceInfo, err := namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "",
		GroupName: "DEFAULT_GROUP",
		PageNo:    1,
		PageSize:  20,
	})
	dom := serviceInfo.Doms
	for d := range dom {
		fmt.Println(d)
	}
	var urlList = []*common.URL{}
	if err != nil {
		return nil, err
	}
	for i := 0; int64(i) < serviceInfo.Count; i++ {
		serviceInstance, err := namingClient.SelectInstances(vo.SelectInstancesParam{
			ServiceName: serviceInfo.Doms[i],
			GroupName:   "DEFAULT_GROUP",
			HealthyOnly: true,
		})
		if err != nil {
			logger.Errorf("fail to get instance", err)
		}
		url := generateUrl(serviceInstance[0])
		urlList = append(urlList, url)
	}
	receiver.url = urlList
	for i, url := range urlList {
		fmt.Println("The %d's url is %s", i, url)
	}
	return urlList, err
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
