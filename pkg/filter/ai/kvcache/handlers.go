package kvcache

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func (f *Filter) manageCache(ctx context.Context, model string, prompt string) {
	tokens, err := f.tokenManager.GetTokens(ctx, model, prompt)
	if err != nil {
		logger.Warnf("[KVCache] tokenize failed: %v", err)
		return
	}
	cacheStatus, err := f.lmcacheClient.Lookup(ctx, &LookupRequest{Tokens: tokens})
	if err != nil {
		logger.Warnf("[KVCache] lookup failed: %v", err)
		return
	}
	decision := f.cacheStrategy.MakeDecision(ctx, cacheStatus)
	if err := f.cacheStrategy.ExecuteDecision(ctx, decision, tokens); err != nil {
		logger.Warnf("[KVCache] execute strategy failed: %v", err)
	}
}

func readRequestBody(req *http.Request) ([]byte, error) {
	if req == nil || req.Body == nil {
		return nil, nil
	}
	if req.GetBody != nil {
		reader, err := req.GetBody()
		if err == nil {
			defer reader.Close()
			return io.ReadAll(reader)
		}
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodyBytes)), nil
	}
	return bodyBytes, nil
}

func extractPromptAndModel(body []byte) (string, string, error) {
	if len(body) == 0 {
		return "", "", nil
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", err
	}
	model, _ := payload["model"].(string)
	prompt := coercePrompt(payload["prompt"])
	if prompt == "" {
		prompt = extractPromptFromMessages(payload["messages"])
	}
	return strings.TrimSpace(prompt), model, nil
}

func coercePrompt(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				parts = append(parts, str)
			}
		}
		return strings.Join(parts, "\n")
	default:
		return ""
	}
}

func extractPromptFromMessages(value any) string {
	msgs, ok := value.([]any)
	if !ok {
		return ""
	}
	parts := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			continue
		}
		if content, ok := msgMap["content"].(string); ok {
			parts = append(parts, content)
		}
	}
	return strings.Join(parts, "\n")
}

func effectiveTimeout(hc *contexthttp.HttpContext, cfg *Config) time.Duration {
	if cfg == nil {
		return 0
	}
	timeout := cfg.RequestTimeout
	if hc != nil && hc.Timeout > 0 && (timeout <= 0 || hc.Timeout < timeout) {
		timeout = hc.Timeout
	}
	if timeout <= 0 {
		return 2 * time.Second
	}
	return timeout
}
