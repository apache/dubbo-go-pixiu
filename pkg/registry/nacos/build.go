package nacos

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	perrors "github.com/pkg/errors"
	"net"
	"strconv"
	"strings"
	"time"
)

func GetNacosConfig(ad *model.Adapter) ([]nacosConstant.ServerConfig, nacosConstant.ClientConfig, error) {
	if ad == nil {
		return []nacosConstant.ServerConfig{}, nacosConstant.ClientConfig{}, perrors.New("url is empty!")
	}
	var serverConfigs = []nacosConstant.ServerConfig{}
	registryConfig := ad.Config["registries"].(map[string]interface{})
	for k := range registryConfig {
		if k == "Nacos" {
			nacosConfig, ok := registryConfig[k]
			if !ok {
				logger.Error(registryConfig)
			}
			m := nacosConfig.(map[string]interface{})
			address := m["address"].(string)
			addresses := strings.Split(address, ",")
			serverConfigs = make([]nacosConstant.ServerConfig, 0, len(address))
			for _, ad := range addresses {
				ip, portStr, err := net.SplitHostPort(ad)
				if err != nil {
					return []nacosConstant.ServerConfig{}, nacosConstant.ClientConfig{},
						perrors.WithMessagef(err, "split [%s] ", address)
				}
				port, _ := strconv.Atoi(portStr)
				serverConfigs = append(serverConfigs, nacosConstant.ServerConfig{IpAddr: ip, Port: uint64(port)})
			}
		}
	}

	clientConfig := nacosConstant.ClientConfig{
		TimeoutMs: uint64(int32(3000 / time.Millisecond)),
		//BeatInterval:        url.GetParamInt(constant.NACOS_BEAT_INTERVAL_KEY, 5000),
		BeatInterval: int64(int32(5000 / time.Millisecond)),
		//NamespaceId:         url.GetParam(constant.NACOS_NAMESPACE_ID, ""),
		NamespaceId: "",
		//AppName:             url.GetParam(constant.NACOS_APP_NAME_KEY, ""),
		//Endpoint:            url.GetParam(constant.NACOS_ENDPOINT, ""),
		//RegionId:            url.GetParam(constant.NACOS_REGION_ID_KEY, ""),
		//AccessKey:           url.GetParam(constant.NACOS_ACCESS_KEY, ""),
		//SecretKey:           url.GetParam(constant.NACOS_SECRET_KEY, ""),
		//OpenKMS:             url.GetParamBool(constant.NACOS_OPEN_KMS_KEY, false),
		OpenKMS: false,
		//UpdateThreadNum:     url.GetParamByIntValue(constant.NACOS_UPDATE_THREAD_NUM_KEY, 20),
		//NotLoadCacheAtStart: url.GetParamBool(constant.NACOS_NOT_LOAD_LOCAL_CACHE, true),
		Username:   "",
		Password:   "",
		RotateTime: "1h",
		MaxAge:     3,
		LogLevel:   "debug",
	}
	//configmap := make(map[string]interface{},2)
	//configmap["serverConfigs"] = serverConfigs
	//configmap["clientConfig"] = clientConfig
	return serverConfigs, clientConfig, nil
}

//func NewNacosClient(adapt *model.Adapter) (naming_client.INamingClient, error) {
//	address := adapt.Config["address"].(string)
//	if len(address) == 0 {
//		return nil, perrors.New("nacos address is empty!")
//	}
//	configMap := make(map[string]interface{}, 2)
//
//	//addresses := strings.Split(url.Path, ",")
//	serverConfigs := make([]nacosConstant.ServerConfig, 0, len(address))
//	ip, portStr, err := net.SplitHostPort(address)
//	if err != nil {
//		return nil, perrors.WithMessagef(err, "split [%s] ", address)
//	}
//	port, _ := strconv.Atoi(portStr)
//	serverConfigs = append(serverConfigs, nacosConstant.ServerConfig{
//		IpAddr: ip,
//		Port:   uint64(port),
//	})
//
//	configMap["serverConfigs"] = serverConfigs
//
//	var clientConfig nacosConstant.ClientConfig
//	if err != nil{
//		logger.Info("fail to parse timeout")
//	}
//	//clientConfig.TimeoutMs = uint64()
//	//clientConfig.ListenInterval = 2 * clientConfig.TimeoutMs
//	//clientConfig.CacheDir = url.GetParam(constant.NACOS_CACHE_DIR_KEY, "")
//	//clientConfig.LogDir = url.GetParam(constant.NACOS_LOG_DIR_KEY, "")
//	//clientConfig.Endpoint = rc.Address
//	clientConfig.Username = adapt.Config["Username"].(string)
//	clientConfig.Password = adapt.Config["Password"].(string)
//	clientConfig.NotLoadCacheAtStart = true
//	//clientConfig.NamespaceId = rc.GetParam(constant.NACOS_NAMESPACE_ID, "")
//
//	configMap["clientConfig"] = clientConfig
//
//	return clients.CreateNamingClient(configMap)
//}
