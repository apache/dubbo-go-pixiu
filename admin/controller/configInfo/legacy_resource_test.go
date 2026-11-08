/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package configInfo

import (
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
)

func TestDeleteResourceInfoRejectsUnsafeIDsBeforeStorageAccess(t *testing.T) {
	for _, query := range []string{"", "?resourceId=0", "?resourceId=-1", "?resourceId=invalid"} {
		t.Run(query, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Request = httptest.NewRequest(http.MethodDelete, "/"+query, nil)

			DeleteResourceInfo(context)

			var result struct {
				Code string `json:"code"`
				Data string `json:"data"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
				t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
			}
			if result.Code != adminconfig.ERR || result.Data != "resourceId must be a positive integer" {
				t.Fatalf("unsafe resource id response: %+v", result)
			}
		})
	}
}

func TestLegacyMethodIDKeepsBodyIDWhenQueryIsOmitted(t *testing.T) {
	methodID, err := legacyMethodID("", 17)
	if err != nil {
		t.Fatalf("legacyMethodID: %v", err)
	}
	if methodID != 17 {
		t.Fatalf("method ID = %d, want body ID 17", methodID)
	}
}

func TestLegacyMethodIDQueryOverridesBodyID(t *testing.T) {
	methodID, err := legacyMethodID(" 23 ", 17)
	if err != nil {
		t.Fatalf("legacyMethodID: %v", err)
	}
	if methodID != 23 {
		t.Fatalf("method ID = %d, want query ID 23", methodID)
	}
}

func TestLegacyMethodIDRejectsMalformedQuery(t *testing.T) {
	if _, err := legacyMethodID("abc", 17); err == nil {
		t.Fatal("legacyMethodID accepted malformed query ID")
	}
}
