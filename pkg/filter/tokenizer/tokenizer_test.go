/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tokenizer

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
)

func TestUnaryResponse(t *testing.T) {
	filter := &Filter{}

	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	c.TargetResp = &client.UnaryResponse{Data: []byte("{\"object\":\"chat.completion\",\"created\":1744871476,\"model\":\"deepseek-chat\",\"choices\":[{\"index\":0,\"message\":{\"role\":\"assistant\",\"content\":\"The sum of 3 and 5 is calculated as follows:\\n\\n\\\\[ 3 + 5 = 8 \\\\]\\n\\nSo, the answer is **8**.\"},\"logprobs\":null,\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":7,\"completion_tokens\":32,\"total_tokens\":39,\"prompt_tokens_details\":{\"cached_tokens\":0},\"prompt_cache_hit_tokens\":0,\"prompt_cache_miss_tokens\":7},\"system_fingerprint\":\"fp_3d5141a69a_prod0225\"}")}
	filter.Encode(c)
}

func TestStreamResponse(t *testing.T) {
	filter := &Filter{}

	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	c.TargetResp = &client.StreamResponse{Stream: io.NopCloser(strings.NewReader("data: {\"object\":\"chat.completion\",\"created\":1744871476,\"model\":\"deepseek-chat\",\"choices\":[{\"index\":0,\"message\":{\"role\":\"assistant\",\"content\":\"The sum of 3 and 5 is calculated as follows:\\n\\n\\\\[ 3 + 5 = 8 \\\\]\\n\\nSo, the answer is **8**.\"},\"logprobs\":null,\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":7,\"completion_tokens\":32,\"total_tokens\":39,\"prompt_tokens_details\":{\"cached_tokens\":0},\"prompt_cache_hit_tokens\":0,\"prompt_cache_miss_tokens\":7},\"system_fingerprint\":\"fp_3d5141a69a_prod0225\"}"))}
	filter.Encode(c)
	time.Sleep(1 * time.Millisecond)
}
