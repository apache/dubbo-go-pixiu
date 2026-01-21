package kvcache

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/hashicorp/go-uuid"
)

func stripIncomingHeaders(headers http.Header) {
	headers.Del(HeaderSessionID)
	headers.Del(HeaderCacheMode)
	headers.Del(HeaderModel)
	headers.Del(HeaderCacheHit)
	headers.Del(HeaderCacheDetail)
	headers.Del(HeaderPolicy)
}

func generateSessionID(existing string) (string, bool) {
	if existing != "" {
		return existing, false
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return "", false
	}
	return id, true
}

func getModels(u *url.URL, targets []string) string {
	if len(targets) == 1 {
		return targets[0]
	}
	if u == nil {
		return ""
	}
	model := strings.TrimSpace(u.Query().Get("model"))
	return model
}

func readResponseHeader(hc *contexthttp.HttpContext, key string) string {
	if hc == nil {
		return ""
	}
	if resp, ok := hc.SourceResp.(*http.Response); ok && resp != nil {
		return resp.Header.Get(key)
	}
	return hc.Writer.Header().Get(key)
}

func storeContextValue(hc *contexthttp.HttpContext, key string, value any) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[key] = value
}

func getContextValue(hc *contexthttp.HttpContext, key string) any {
	if hc.Params == nil {
		return nil
	}
	return hc.Params[key]
}

func getContext(hc *contexthttp.HttpContext) context.Context {
	if hc != nil && hc.Ctx != nil {
		return hc.Ctx
	}
	return context.Background()
}
