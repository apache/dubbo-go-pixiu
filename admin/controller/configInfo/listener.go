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
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// @Tags Config
// @Summary get all Listener list
// @Description get all Listener list
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/listener/list [get]
// GetListenerList get all Listener list
func GetListenerList(c *gin.Context) {
	rst, err := logic.BizGetListeners()
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(rst))
}

// @Tags Config
// @Summary create a Listener
// @Description pass the YAML/JSON for the Listener through the form's content field to create the Listener.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param content formData string true "Listener content"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /config/api/listener [put]
// CreateListener create a Listener
func CreateListener(c *gin.Context) {
	body := c.PostForm("content")
	res := &config.Listener{}
	err := yaml.UnmarshalYML([]byte(body), res)
	logger.Debug(body)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	err = logic.BizCreateListener(res)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("create Listener success!"))
}

// @Tags Config
// @Summary delete Listener
// @Description delete the Listener based on the name
// @Produce application/json
// @Param name query string true "Listener id"
// @Success 200 {object} string
// @Router /config/api/listener [delete]
// DeleteListener delete Listener
func DeleteListener(c *gin.Context) {
	id := c.Query(logic.Listener)
	err := logic.BizDeleteListener(id)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}

	c.JSON(http.StatusOK, adminconfig.WithRet("delete Listener success!"))
}

// @Tags Config
// @Summary get Listener detail
// @Description get Listener details based on name
// @Produce application/json
// @Param name query string true "Listener name"
// @Success 200 {object} string
// @Router /config/api/listener/detail [get]
// DetailListener get Listener detail
func DetailListener(c *gin.Context) {
	name := c.Query(logic.Listener)
	res, err := logic.BizGetListener(name)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(res))
}

// @Tags Config
// @Summary update Listener
// @Description pass the YAML/JSON for the Listener through the form's content field to update the Listener.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param content formData string true "Listener content"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /config/api/listener [post]
// UpdateListener update Listener
func UpdateListener(c *gin.Context) {
	body := c.PostForm("content")
	res := &config.Listener{}
	err := yaml.UnmarshalYML([]byte(body), res)
	logger.Debug(body)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	err = logic.BizUpdateListener(res)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("update Listener success!"))
}
