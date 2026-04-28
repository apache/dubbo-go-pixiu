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

package opa

import (
	"net/http"
	"strings"
)

import (
	"github.com/gin-gonic/gin"

	perrors "github.com/pkg/errors"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
)

// @Tags Config
// @Summary upload OPA policy (server mode)
// @Router /config/api/opa/policy [put]
func PutOPAPolicy(c *gin.Context) {
	// For PUT, we typically expect Form Data
	serverURL := resolveOPAServerURL(c.PostForm("server_url"))
	policyID := resolveOPAPolicyID(c.PostForm("policy_id"))
	bearerToken := c.PostForm("bearer_token")
	content := c.PostForm("content")

	if content == "" {
		c.JSON(http.StatusOK, adminconfig.WithError(perrors.New("rego content is required")))
		return
	}

	if err := logic.BizPutOPAPolicy(serverURL, policyID, bearerToken, content); err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Update Success"))
}

// @Tags Config
// @Summary get OPA policy (server mode)
// @Router /config/api/opa/policy [get]
func GetOPAPolicy(c *gin.Context) {
	var query adminconfig.OPAQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}

	serverURL := resolveOPAServerURL(query.ServerURL)
	policyID := resolveOPAPolicyID(query.PolicyID)

	result, err := logic.BizGetOPAPolicy(serverURL, policyID, query.BearerToken)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(result))
}

// @Tags Config
// @Summary delete OPA policy (server mode)
// @Router /config/api/opa/policy [delete]
func DeleteOPAPolicy(c *gin.Context) {
	var query adminconfig.OPAQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}

	serverURL := resolveOPAServerURL(query.ServerURL)
	policyID := resolveOPAPolicyID(query.PolicyID)

	if err := logic.BizDeleteOPAPolicy(serverURL, policyID, query.BearerToken); err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Delete Success"))
}

func resolveOPAServerURL(serverURL string) string {
	serverURL = strings.TrimSpace(serverURL)
	if serverURL != "" {
		return serverURL
	}
	if adminconfig.Bootstrap != nil {
		if trimmed := strings.TrimSpace(adminconfig.Bootstrap.OPA.ServerURL); trimmed != "" {
			return trimmed
		}
	}
	return adminconfig.DefaultOPAServerURL
}

func resolveOPAPolicyID(policyID string) string {
	policyID = strings.TrimSpace(policyID)
	if policyID != "" {
		return policyID
	}
	if adminconfig.Bootstrap != nil {
		if trimmed := strings.TrimSpace(adminconfig.Bootstrap.OPA.PolicyID); trimmed != "" {
			return trimmed
		}
	}
	return adminconfig.DefaultOPAPolicyID
}
