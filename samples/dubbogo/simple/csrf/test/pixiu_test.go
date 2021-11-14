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

package csrf

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

var token string

func GetToken(t *testing.T) bool {
	return t.Run("token", func(t *testing.T) {
		urlStr := "http://localhost:1314/login?key=pixiu&secret=pixiu888"
		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest("GET", urlStr, nil)
		assert.NoError(t, err)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp)
		s, _ := ioutil.ReadAll(resp.Body)
		token = string(s)
	})
}

func TestCsrfHeader(t *testing.T) {
	if GetToken(t) && token != "" {
		urlStr := "http://localhost:8888/user/"
		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest("GET", urlStr, nil)
		assert.NoError(t, err)
		req.Header.Set("csrfSalt", "pixiu")
		req.Header.Set("pixiu", token)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp)
		s, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(s))
		assert.True(t, strings.Contains(string(s), "success"))
	}
}

func TestCsrfQuery(t *testing.T) {
	if GetToken(t) && token != "" {
		urlStr := "http://localhost:8888/user?pixiu=" + token
		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest("GET", urlStr, nil)
		assert.NoError(t, err)
		req.Header.Set("csrfSalt", "pixiu")
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp)
		s, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(s))
		assert.True(t, strings.Contains(string(s), "success"))
	}
}
