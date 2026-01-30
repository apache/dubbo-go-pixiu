package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type tokenizeRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type tokenizeResponse struct {
	Count  int   `json:"count"`
	Tokens []int `json:"tokens"`
	MaxLen int   `json:"max_model_len"`
}

type tokensRequest struct {
	Tokens []int `json:"tokens"`
}

type lmcacheResponse struct {
	EventID   string `json:"event_id"`
	NumTokens int    `json:"num_tokens"`
}

type lookupResponse struct {
	EventID    string                 `json:"event_id"`
	LayoutInfo map[string]cacheLayout `json:"layout_info"`
}

type cacheLayout struct {
	Location   string `json:"0"`
	TokenCount int    `json:"1"`
}

func main() {
	var (
		vllmAddr    = flag.String("vllm", ":8000", "vLLM mock listen address")
		lmcacheAddr = flag.String("lmcache", ":8080", "LMCache mock listen address")
	)
	flag.Parse()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/tokenize", handleTokenize)
		mux.HandleFunc("/v1/chat/completions", handleChatCompletions)
		fmt.Printf("[mock] vLLM listening on %s\n", *vllmAddr)
		if err := http.ListenAndServe(*vllmAddr, mux); err != nil {
			panic(err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/lookup", handleLookup)
	mux.HandleFunc("/pin", handleSimpleLMCache)
	mux.HandleFunc("/compress", handleSimpleLMCache)
	mux.HandleFunc("/evict", handleSimpleLMCache)
	fmt.Printf("[mock] LMCache listening on %s\n", *lmcacheAddr)
	if err := http.ListenAndServe(*lmcacheAddr, mux); err != nil {
		panic(err)
	}
}

func handleTokenize(w http.ResponseWriter, r *http.Request) {
	var req tokenizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	prompt := strings.TrimSpace(req.Prompt)
	tokens := make([]int, 0, len(prompt))
	for _, r := range prompt {
		if r == ' ' || r == '\n' || r == '\t' {
			continue
		}
		tokens = append(tokens, int(r))
	}
	resp := tokenizeResponse{
		Count:  len(tokens),
		Tokens: tokens,
		MaxLen: 4096,
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"id":      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   "mock-model",
		"choices": []map[string]any{
			{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": "mock response",
				},
				"finish_reason": "stop",
			},
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleLookup(w http.ResponseWriter, r *http.Request) {
	var req tokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	layout := map[string]cacheLayout{
		"vllm-instance-1": {
			Location:   "LocalCPUBackend",
			TokenCount: len(req.Tokens),
		},
	}
	resp := lookupResponse{
		EventID:    fmt.Sprintf("lookup-%d", time.Now().UnixNano()),
		LayoutInfo: layout,
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleSimpleLMCache(w http.ResponseWriter, r *http.Request) {
	var req tokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	resp := lmcacheResponse{
		EventID:   fmt.Sprintf("event-%d", time.Now().UnixNano()),
		NumTokens: len(req.Tokens),
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
