package proxyrewrite

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/go-playground/assert/v2"
	"net/http"
	"testing"
)

func TestDecode(t *testing.T) {
	url := "/user-service/query"

	request, _ := http.NewRequest("Get", url, nil)
	ctx := mock.GetMockHTTPContext(request)

	factory := &FilterFactory{cfg: &Config{UriRegex: []string{"^/([^/]*)/(.*)$", "/$2"}}}
	chain := filter.NewDefaultFilterChain()
	_ = factory.PrepareFilterChain(ctx, chain)
	chain.OnDecode(ctx)

	assert.Equal(t, ctx.GetUrl(), "/query")
}

func TestConfigUpdate(t *testing.T) {
	url := "/user-service/query"

	request, _ := http.NewRequest("Get", url, nil)
	ctx := mock.GetMockHTTPContext(request)

	cfg := &Config{UriRegex: []string{"^/([^/]*)/(.*)$", "/$2"}}

	factory := &FilterFactory{cfg: cfg}
	chain := filter.NewDefaultFilterChain()
	_ = factory.PrepareFilterChain(ctx, chain)

	cfg.UriRegex = []string{}

	chain.OnDecode(ctx)

	assert.Equal(t, ctx.GetUrl(), "/query")
}
