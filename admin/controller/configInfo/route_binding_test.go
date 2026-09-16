/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package configInfo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

func TestGetRouteBindingSchema(t *testing.T) {
	response := invokeRouteBindingHandler(t, http.MethodGet, nil, GetRouteBindingSchema)
	var result struct {
		Code string                `json:"code"`
		Data []schema.ObjectSchema `json:"data"`
	}
	decodeRouteBindingResponse(t, response, &result)
	if result.Code != adminconfig.OK {
		t.Fatalf("response code: want %s, got %s", adminconfig.OK, result.Code)
	}
	if len(result.Data) != 1 || result.Data[0].Kind != schema.KindAdminRouteBinding {
		t.Fatalf("schema response: %+v", result.Data)
	}
}

func TestValidateRouteBindingReturnsDefaults(t *testing.T) {
	body := []byte(`{
  "kind": "AdminRouteBinding",
  "metadata": {"name": "user-get"},
  "spec": {
    "entry": {"path": "/api/users", "method": "GET"},
    "target": {
      "application": "UserProvider",
      "interface": "com.example.UserService",
      "method": "GetUser",
      "cluster": "user-dubbo"
    }
  }
}`)
	response := invokeRouteBindingHandler(t, http.MethodPost, body, ValidateRouteBinding)
	var result struct {
		Code string                         `json:"code"`
		Data routeBindingValidationResponse `json:"data"`
	}
	decodeRouteBindingResponse(t, response, &result)
	if result.Code != adminconfig.OK || !result.Data.Valid {
		t.Fatalf("validation response: %+v", result)
	}
	if result.Data.Object.Spec["entry"].(map[string]any)["protocol"] != "http" {
		t.Fatalf("entry defaults: %+v", result.Data.Object.Spec["entry"])
	}
	if result.Data.Object.Spec["target"].(map[string]any)["protocol"] != "dubbo" {
		t.Fatalf("target defaults: %+v", result.Data.Object.Spec["target"])
	}
	if result.Data.Object.Spec["enabled"] != true {
		t.Fatalf("enabled default: %+v", result.Data.Object.Spec["enabled"])
	}
	if result.Data.Object.Spec["publish"].(map[string]any)["mode"] != "draft" {
		t.Fatalf("publish defaults: %+v", result.Data.Object.Spec["publish"])
	}
}

func TestValidateRouteBindingReturnsFieldIssues(t *testing.T) {
	body := []byte(`{
  "kind": "AdminRouteBinding",
  "metadata": {"name": "bad-route"},
  "spec": {
    "entry": {"path": "api/users", "method": "TRACE"},
    "target": {"application": "UserProvider"}
  }
}`)
	response := invokeRouteBindingHandler(t, http.MethodPost, body, ValidateRouteBinding)
	var result struct {
		Code string                    `json:"code"`
		Data routeBindingErrorResponse `json:"data"`
	}
	decodeRouteBindingResponse(t, response, &result)
	if result.Code != adminconfig.ERR {
		t.Fatalf("response code: want %s, got %s", adminconfig.ERR, result.Code)
	}
	if len(result.Data.Issues) == 0 {
		t.Fatalf("expected validation issues, got %+v", result.Data)
	}
	foundPathIssue := false
	for _, issue := range result.Data.Issues {
		if issue.Path == "spec.entry.path" {
			foundPathIssue = true
			break
		}
	}
	if !foundPathIssue {
		t.Fatalf("path validation issue missing: %+v", result.Data.Issues)
	}
}

func TestPreviewRouteBindingReturnsLegacyYAML(t *testing.T) {
	body := []byte(`{
  "kind": "AdminRouteBinding",
  "metadata": {"name": "user-get"},
  "spec": {
    "entry": {"path": "/api/users", "method": "GET"},
    "target": {
      "application": "UserProvider",
      "interface": "com.example.UserService",
      "method": "GetUser",
      "cluster": "user-dubbo"
    }
  }
}`)
	response := invokeRouteBindingHandler(t, http.MethodPost, body, PreviewRouteBinding)
	var result struct {
		Code string                      `json:"code"`
		Data routeBindingPreviewResponse `json:"data"`
	}
	decodeRouteBindingResponse(t, response, &result)
	if result.Code != adminconfig.OK {
		t.Fatalf("response code: want %s, got %s", adminconfig.OK, result.Code)
	}
	if result.Data.YAML == "" || !bytes.Contains([]byte(result.Data.YAML), []byte("resources:")) {
		t.Fatalf("legacy preview yaml: %q", result.Data.YAML)
	}
	if bytes.Contains([]byte(result.Data.YAML), []byte("AdminRouteBinding")) || bytes.Contains([]byte(result.Data.YAML), []byte("publish:")) {
		t.Fatalf("control-plane fields leaked into legacy preview: %q", result.Data.YAML)
	}
}

func invokeRouteBindingHandler(t *testing.T, method string, body []byte, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest(method, "/", bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request
	handler(context)
	return response
}

func decodeRouteBindingResponse(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if response.Code != http.StatusOK {
		t.Fatalf("http status: want %d, got %d", http.StatusOK, response.Code)
	}
	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response %q: %v", response.Body.String(), err)
	}
}
