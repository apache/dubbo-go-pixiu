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
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestPost1(t *testing.T) {
	url := "http://localhost:8883/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=GetUserByName"
	data := "{\"types\":\"string\",\"values\":\"tc\"}"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPost2(t *testing.T) {
	url := "http://localhost:8883/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=UpdateUserByName"
	data := "{\"types\":\"string,object\",\"values\":[\"tc\",{\"id\":\"0001\",\"code\":1,\"name\":\"tc\",\"age\":15}]}"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPost3(t *testing.T) {
	url := "http://localhost:8883/api/v1/test-dubbo/UserService/com.dubbogo.pixiu.UserService?group=test&version=1.0.0&method=GetUserByCode"
	data := "{\"types\":\"int\",\"values\":1}"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}
