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

package config

import (
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/pkg/i18n"
)

// Response represents a standard API response.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// Success sends a successful response with data.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
	})
}

// SuccessMessage sends a successful response with a message.
func SuccessMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
	})
}

// Error sends an error response.
func Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: err.Error(),
	})
}

// ErrorMessage sends an error response with a custom message.
func ErrorMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: message,
	})
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden sends a 403 response.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
	})
}

// getLang extracts the language from gin context, defaults to English.
func getLang(c *gin.Context) string {
	lang := c.GetString("lang")
	if lang == "" {
		return i18n.LangEnglish
	}
	return lang
}

// SuccessMessageI18n sends a successful response with a translated message.
func SuccessMessageI18n(c *gin.Context, msgKey string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: i18n.T(getLang(c), msgKey),
	})
}

// ErrorMessageI18n sends an error response with a translated message.
func ErrorMessageI18n(c *gin.Context, msgKey string) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: i18n.T(getLang(c), msgKey),
	})
}

// UnauthorizedI18n sends a 401 response with a translated message.
func UnauthorizedI18n(c *gin.Context, msgKey string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: i18n.T(getLang(c), msgKey),
	})
}

// ForbiddenI18n sends a 403 response with a translated message.
func ForbiddenI18n(c *gin.Context, msgKey string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: i18n.T(getLang(c), msgKey),
	})
}
