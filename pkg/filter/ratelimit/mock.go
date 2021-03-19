package ratelimit

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
)

func GetMockedRateLimitConfig() (*Config, error) {
	c := &Config{}
	if err := yaml.UnmarshalYMLConfig("./mock/config.yml", c); err != nil {
		return nil, err
	}
	return c, nil
}
