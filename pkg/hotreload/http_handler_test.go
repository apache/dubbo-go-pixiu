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

package hotreload

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestReloadHandlerRejectsOversizedBody(t *testing.T) {
	withReloadSecret(t, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/-/reload", strings.NewReader(strings.Repeat("a", maxReloadBodyBytes+1)))
	req.RemoteAddr = "192.0.2.1:12345"
	req.Header.Set("X-Reload-Token", "test-secret")

	rr := httptest.NewRecorder()
	(&ReloadHandler{}).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, rr.Code)
}

func TestReloadHandlerAllowsBodyAtSizeLimit(t *testing.T) {
	withReloadSecret(t, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/-/reload", strings.NewReader(strings.Repeat("a", maxReloadBodyBytes)))
	req.RemoteAddr = "192.0.2.1:12345"
	req.Header.Set("X-Reload-Token", "test-secret")

	rr := httptest.NewRecorder()
	(&ReloadHandler{}).ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusRequestEntityTooLarge, rr.Code)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestReloadHandlerEmptyBodyFallsBackToFileReload(t *testing.T) {
	withReloadSecret(t, "test-secret")
	oldConfigPath := configPath
	configPath = ""
	t.Cleanup(func() {
		configPath = oldConfigPath
	})

	req := httptest.NewRequest(http.MethodPost, "/-/reload", strings.NewReader(""))
	req.RemoteAddr = "192.0.2.1:12345"
	req.Header.Set("X-Reload-Token", "test-secret")

	rr := httptest.NewRecorder()
	(&ReloadHandler{}).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "config path not set")
}

func withReloadSecret(t *testing.T, secret string) {
	t.Helper()

	oldSecret := reloadSecret
	reloadSecret = secret
	t.Cleanup(func() {
		reloadSecret = oldSecret
	})
}
