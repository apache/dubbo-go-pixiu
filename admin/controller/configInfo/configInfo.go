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
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// @Tags Config
// @Summary get basic configuration of Pixiu
// @Description get pixiu base info such as name,desc
// @Produce application/json
// @Success 200 {string} string "YAML content"
// @Router /config/api/base [get]
// GetBaseInfo get pixiu base info such as name,desc
func GetBaseInfo(c *gin.Context) {
	conf, err := logic.BizGetBaseInfo()
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	data, _ := yaml.MarshalYML(conf)
	c.JSON(http.StatusOK, adminconfig.WithRet(string(data)))
}

// @Tags Config
// @Summary modify pixiu base info such as name,desc
// @Description Pass YAML content through the form's content field to set basic information.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param content formData string true "YAML content"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /config/api/base/ [post]
// @Router /config/api/base/ [put]
// SetBaseInfo modify pixiu base info such as name,desc
func SetBaseInfo(c *gin.Context) {
	body := c.PostForm("content")

	baseInfo := &adminconfig.BaseInfo{}
	err := yaml.UnmarshalYML([]byte(body), baseInfo)
	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}

	setErr := logic.BizSetBaseInfo(baseInfo, true)
	if setErr != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("success"))
}
