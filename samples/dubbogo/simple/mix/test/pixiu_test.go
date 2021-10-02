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

package test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	url := "http://localhost:8882/api/v1/test-dubbo/user/tc?age=99"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(s), "0001"))
}

func TestPut1(t *testing.T) {
	url := "http://localhost:8882/api/v1/test-dubbo/user/tc"
	data := "{\"id\":\"0001\",\"code\":1,\"name\":\"tc\",\"age\":66}"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "true", string(s))
}

func TestPut2(t *testing.T) {
	url := "http://localhost:8882/api/v1/test-dubbo/user?name=tc"
	data := "{\"id\":\"0001\",\"code\":1,\"name\":\"tc\",\"age\":55}"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "true", string(s))
}
