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

package handler

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/i18n"
)

// Language returns a middleware that parses Accept-Language header and sets the language in context.
func (h *Handler) Language() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader("Accept-Language")
		lang := i18n.ParseAcceptLanguage(acceptLang)
		c.Set("lang", lang)
		c.Next()
	}
}

// RequirePermission returns a middleware that checks if the user has the required permission.
func (h *Handler) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			config.UnauthorizedI18n(c, i18n.MsgNotAuthenticated)
			c.Abort()
			return
		}

		hasPermission, err := h.mysql.HasPermission(username, resource, action)
		if err != nil {
			h.logger.Error("failed to check permission: " + err.Error())
			config.ErrorMessageI18n(c, i18n.MsgCheckPermFailed)
			c.Abort()
			return
		}

		if !hasPermission {
			config.ForbiddenI18n(c, i18n.MsgPermissionDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
