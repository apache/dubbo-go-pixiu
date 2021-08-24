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
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestGet1(t *testing.T) {
	url := "http://localhost:8885/api/v1/test-dubbo/user/name/tc"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"age\":18,\"code\":1,\"iD\":\"0001\",\"name\":\"tc\",\"time\":\"2021-08-01T18:08:41+08:00\"}", string(s))
}

func TestGet2(t *testing.T) {
	url := "http://localhost:8885/api/v1/test-dubbo/user/code/1"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"age\":18,\"code\":1,\"iD\":\"0001\",\"name\":\"tc\",\"time\":\"2021-08-01T18:08:41+08:00\"}", string(s))
}

func TestGet3(t *testing.T) {
	url := "http://localhost:8885/api/v1/test-dubbo/user/name/tc/age/99"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"age\":18,\"code\":1,\"iD\":\"0001\",\"name\":\"tc\",\"time\":\"2021-08-01T18:08:41+08:00\"}", string(s))
}
