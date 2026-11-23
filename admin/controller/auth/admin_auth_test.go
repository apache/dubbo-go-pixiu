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

package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/gin-gonic/gin"
)

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		username   string
		isAdmin    bool
		checkErr   error
		wantStatus int
	}{
		{name: "administrator allowed", username: "admin", isAdmin: true, wantStatus: http.StatusNoContent},
		{name: "ordinary user denied", username: "user", wantStatus: http.StatusForbidden},
		{name: "role lookup failure denied", username: "admin", checkErr: errors.New("database unavailable"), wantStatus: http.StatusForbidden},
		{name: "missing username denied", wantStatus: http.StatusForbidden},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			checkerCalled := false
			router.GET("/", func(c *gin.Context) {
				if test.username != "" {
					c.Set("claims", &CustomClaims{Username: test.username})
				}
				c.Next()
			}, adminAuth(func(username string) (bool, error) {
				checkerCalled = true
				if username != test.username {
					t.Errorf("admin checker username = %q, want %q", username, test.username)
				}
				return test.isAdmin, test.checkErr
			}), func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			router.ServeHTTP(response, request)
			if response.Code != test.wantStatus {
				t.Fatalf("response status = %d, want %d", response.Code, test.wantStatus)
			}
			if checkerCalled != (test.username != "") {
				t.Fatalf("admin checker called = %t, want %t", checkerCalled, test.username != "")
			}
		})
	}
}
