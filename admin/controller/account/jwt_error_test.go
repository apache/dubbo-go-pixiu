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
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

func TestJWTConfigurationErrorsDoNotExposeDetails(t *testing.T) {
	t.Setenv("DUBBOGO_PIXIU_JWT_SIGN_KEY", "short-key")
	for _, test := range []struct {
		name    string
		handler gin.HandlerFunc
	}{
		{name: "login", handler: func(c *gin.Context) { generateToken(c, "admin") }},
		{name: "logout", handler: Logout},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			test.handler(context)
			require.Equal(t, http.StatusInternalServerError, response.Code)
			require.Contains(t, response.Body.String(), "authentication is not configured")
			require.NotContains(t, response.Body.String(), "DUBBOGO_PIXIU_JWT_SIGN_KEY")
		})
	}
}
