package plugins

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitPluginsGroup(t *testing.T) {

	config, err := config.LoadAPIConfigFromFile("/Users/zhengxianle/dubbo-go-proxy/configs/api_config.yaml")
	assert.Empty(t, err)

	InitPluginsGroup(config.PluginsGroup, config.PluginFilePath)
}

func TestInitApiUrlWithFilterChain(t *testing.T) {
	config, err := config.LoadAPIConfigFromFile("/Users/zhengxianle/dubbo-go-proxy/configs/api_config.yaml")
	assert.Empty(t, err)

	InitPluginsGroup(config.PluginsGroup, config.PluginFilePath)
	InitApiUrlWithFilterChain(config.Resources)
}

func TestGetApiFilterFuncsWithApiUrl(t *testing.T) {
	config, err := config.LoadAPIConfigFromFile("/Users/zhengxianle/dubbo-go-proxy/configs/api_config.yaml")
	assert.Empty(t, err)

	InitPluginsGroup(config.PluginsGroup, config.PluginFilePath)
	InitApiUrlWithFilterChain(config.Resources)

	fltc := GetApiFilterFuncsWithApiUrl("/")

	assert.Equal(t, len(fltc), 3)
}
