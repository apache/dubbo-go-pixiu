package main

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/ghodss/yaml"

	//"github.com/knadh/koanf/parsers/yaml"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

type Load interface {
	load(param vo.ConfigParam) (string, error)
}

func main() {
	//配置连接信息
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "localhost",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig: constant.NewClientConfig(
				constant.WithNamespaceId("dubbo-go-pixiu"),
				constant.WithTimeoutMs(5000),
				constant.WithNotLoadCacheAtStart(true),
				constant.WithLogLevel("info"),
			),
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	//读取文件
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "pixiu.config",
		Group:  "DEFAULT_GROUP",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)

	cfg := &model.Bootstrap{}
	err = yaml.Unmarshal([]byte(content), cfg)
	if err != nil {
		log.Fatalln("[config] [yaml load] convert YAML to JSON failed, ", err)
	}
}
