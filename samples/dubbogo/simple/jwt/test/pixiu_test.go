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

package jwt

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

func TestLocalRemote(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:1314/remote", nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp)
}

func TestLocal(t *testing.T) {

	verify(t, "http://localhost:8888/health", http.StatusUnauthorized)

	s := verify(t, "http://localhost:8888/user", http.StatusOK)
	assert.True(t, strings.Contains(s, "user"))

	s = verify(t, "http://localhost:8888/user/pixiu", http.StatusOK)
	assert.True(t, strings.Contains(s, "pixiu"))

	s = verify(t, "http://localhost:8888/prefix", http.StatusOK)
	assert.True(t, strings.Contains(s, "prefix"))
}

func TestRemote(t *testing.T) {
	s := verify(t, "http://localhost:8888/prefix", http.StatusOK)
	assert.True(t, strings.Contains(s, "prefix"))
}

func TestSpringCloud(t *testing.T) {
	s := verify(t, "http://localhost:8888/user-service/echo/Pixiu", http.StatusOK)
	assert.True(t, strings.Contains(s, "User"))
	s = verify(t, "http://localhost:8888/auth-service/echo/Pixiu", http.StatusOK)
	assert.True(t, strings.Contains(s, "Auth"))
}

func verify(t *testing.T, url string, status int) string {
	var token = "eyJraWQiOiJlZThkNjI2ZCIsInR5cCI6IkpXVCIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJXZWlkb25nIiwiYXVkIjoiVGFzaHVhbiIsImlzcyI6Imp3a3Mtc2VydmljZS5hcHBzcG90LmNvbSIsImlhdCI6MTYzMTM2OTk1NSwianRpIjoiNDY2M2E5MTAtZWU2MC00NzcwLTgxNjktY2I3NDdiMDljZjU0In0.LwD65d5h6U_2Xco81EClMa_1WIW4xXZl8o4b7WzY_7OgPD2tNlByxvGDzP7bKYA9Gj--1mi4Q4li4CAnKJkaHRYB17baC0H5P9lKMPuA6AnChTzLafY6yf-YadA7DmakCtIl7FNcFQQL2DXmh6gS9J6TluFoCIXj83MqETbDWpL28o3XAD_05UP8VLQzH2XzyqWKi97mOuvz-GsDp9mhBYQUgN3csNXt2v2l-bUPWe19SftNej0cxddyGu06tXUtaS6K0oe0TTbaqc3hmfEiu5G0J8U6ztTUMwXkBvaknE640NPgMQJqBaey0E4u0txYgyvMvvxfwtcOrDRYqYPBnA"
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, status, resp.StatusCode)
	assert.NotNil(t, resp)
	s, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(s))
	return string(s)
}
