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

package configInfo

import (
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// @Tags Config
// @Summary batch Release PluginGroup Config
// @Description publish the PluginGroup from the unpublished space to the published space.
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/plugin_group/publish [put]
// BatchReleasePluginGroup publishes the staged PluginGroup configuration.
func BatchReleasePluginGroup(c *gin.Context) {
	if err := logic.BizPublishPluginGroupConfig(); err != nil {
		logger.Warnf("Batch Release PluginGroup err, %v", err)
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("success"))
}
