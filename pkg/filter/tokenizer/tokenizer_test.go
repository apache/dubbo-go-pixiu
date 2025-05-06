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
	c.TargetResp = &client.UnaryResponse{
		Data: []byte(`{
		"usage": {
			"prompt_tokens": 7,
			"completion_tokens": 32,
			"total_tokens": 39,
			"prompt_tokens_details": {
				"cached_tokens": 0
			},
			"prompt_cache_hit_tokens": 0,
			"prompt_cache_miss_tokens": 7
		}
	}`)}
	filter.Encode(c)
}

func TestStreamResponse(t *testing.T) {
	filter := &Filter{}

	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	s := io.NopCloser(strings.NewReader(`data: {
		"usage": {
			"prompt_tokens": 7,
			"completion_tokens": 32,
			"total_tokens": 39,
			"prompt_tokens_details": {
				"cached_tokens": 0
			},
			"prompt_cache_hit_tokens": 0,
			"prompt_cache_miss_tokens": 7
		}
	}

`))
	c.TargetResp = &client.StreamResponse{Stream: s}
	filter.Encode(c)
	buf := make([]byte, 1024)
	c.TargetResp.(*client.StreamResponse).Stream.Read(buf)
	time.Sleep(3 * time.Millisecond)
	c.TargetResp.(*client.StreamResponse).Stream.Close()
}
