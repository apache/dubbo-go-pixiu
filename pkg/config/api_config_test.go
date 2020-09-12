package config_test

import (
	"log"
	"testing"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func TestLoadAPIConfigFromFile(t *testing.T) {
	apiC, err := config.LoadAPIConfigFromFile("")
	assert.Empty(t, apiC)
	assert.EqualError(t, err, "Config file not specified")

	apiC, err = config.LoadAPIConfigFromFile("./mock/api_config.yml")
	assert.Empty(t, err)
	assert.Equal(t, apiC.Name, "api name")
	bytes, _ := yaml.Marshal(apiC)
	log.Printf("%s", bytes)
}
