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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/samples/http/grpc/proto"
)

const (
	url  = "http://localhost:8881/api/v1/provider.UserProvider/GetUser"
	data = "{\"userId\":1}"
)

func TestGet(t *testing.T) {
	c := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte{}))
	assert.NoError(t, err)
	resp, err := c.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)
	var r proto.GetUserResponse
	err = json.Unmarshal(body, &r)
	assert.NoError(t, err)
	assert.Equal(t, "user(s) query successfully", r.Message)
	assert.Equal(t, 2, len(r.Users))
	assert.Equal(t, int32(1), r.Users[0].UserId)
	assert.Equal(t, "Kenway", r.Users[0].Name)
	assert.Equal(t, int32(2), r.Users[1].UserId)
	assert.Equal(t, "Ken", r.Users[1].Name)
}

func TestPost(t *testing.T) {
	c := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	assert.NoError(t, err)
	resp, err := c.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)
	var r proto.GetUserResponse
	err = json.Unmarshal(body, &r)
	assert.NoError(t, err)
	assert.Equal(t, "user(s) query successfully", r.Message)
	assert.Equal(t, 1, len(r.Users))
	assert.Equal(t, int32(1), r.Users[0].UserId)
	assert.Equal(t, "Kenway", r.Users[0].Name)
}
