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

package account

import (
	"errors"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

func TestRegisterRequiresAdministrator(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newRegisterContext(t, map[string]string{"username": "new-admin", "password": "password123"})
	ctx.Set("claims", &auth.CustomClaims{Username: "operator"})
	created := false

	registerAccount(ctx, func(username string) (bool, error) {
		if username != "operator" {
			t.Fatalf("checked username = %q, want operator", username)
		}
		return false, errors.New("not an administrator")
	}, func(string, string) error {
		created = true
		return nil
	})

	if created {
		t.Fatal("created an account for a non-administrator")
	}
	if !strings.Contains(recorder.Body.String(), "not an administrator") {
		t.Fatalf("response should explain the authorization failure: %s", recorder.Body.String())
	}
}

func TestRegisterTrimsUsernameAndHashesPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newRegisterContext(t, map[string]string{"username": "  new-admin  ", "password": "password123"})
	ctx.Set("claims", &auth.CustomClaims{Username: "admin"})
	var gotUsername, gotPassword string

	registerAccount(ctx, func(username string) (bool, error) {
		if username != "admin" {
			t.Fatalf("checked username = %q, want admin", username)
		}
		return true, nil
	}, func(username, password string) error {
		gotUsername, gotPassword = username, password
		return nil
	})

	if gotUsername != "new-admin" {
		t.Fatalf("created username = %q, want new-admin", gotUsername)
	}
	if gotPassword != utils.Md5("password123") {
		t.Fatal("registration did not hash the submitted password")
	}
	if !strings.Contains(recorder.Body.String(), "Register successfully") {
		t.Fatalf("expected successful registration response, got %s", recorder.Body.String())
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, recorder := newRegisterContext(t, map[string]string{"username": "new-admin", "password": "short"})
	ctx.Set("claims", &auth.CustomClaims{Username: "admin"})
	created := false

	registerAccount(ctx, func(string) (bool, error) { return true, nil }, func(string, string) error {
		created = true
		return nil
	})

	if created {
		t.Fatal("created an account with a short password")
	}
	if !strings.Contains(recorder.Body.String(), "8 and 128") {
		t.Fatalf("expected password validation error, got %s", recorder.Body.String())
	}
}

func newRegisterContext(t *testing.T, fields map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	values := url.Values{}
	for key, value := range fields {
		values.Set(key, value)
	}
	request := httptest.NewRequest("POST", "/register", strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	return ctx, recorder
}
