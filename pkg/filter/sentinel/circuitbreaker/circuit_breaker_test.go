package circuitbreaker

import (
	stdHttp "net/http"
	"testing"
)

import (
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	pkgs "github.com/apache/dubbo-go-pixiu/pkg/filter/sentinel"
)

func TestFilter(t *testing.T) {
	f := FilterFactory{cfg: &Config{}}

	mockYaml, err := yaml.MarshalYML(mockConfig())
	assert.Nil(t, err)

	assert.Nil(t, yaml.UnmarshalYML(mockYaml, f.Config()))

	assert.Nil(t, f.Apply())

	decoder := &Filter{cfg: f.cfg, matcher: f.matcher}
	request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
	c := mock.GetMockHTTPContext(request)

	assert.Equal(t, decoder.Decode(c), filter.Continue)
}

func mockConfig() *Config {
	c := Config{
		Rules: []*pkgs.Resource{
			{
				Name: "test-dubbo",
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/test-dubbo/user"},
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/test-dubbo/user/*"},
				},
			},
		},
		Resources: []Resources{{Resource: "test-dubbo",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   1000,
			Threshold:        1.0,
		}},
	}
	return &c
}
