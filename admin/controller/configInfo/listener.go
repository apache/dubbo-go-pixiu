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
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"

	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// GetListenerList get all Listener list
func GetListenerList(c *gin.Context) {
	rst, err := logic.BizGetListeners()
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet(rst))
}

// CreateListener create a Listener
func CreateListener(c *gin.Context) {
	body := c.PostForm("content")
	res := &fc.Listener{}
	err := yaml.UnmarshalYML([]byte(body), res)
	logger.Debug(body)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	err = logic.BizCreateListener(res)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("create Listener success!"))
}

// DeleteListener delete resource
func DeleteListener(c *gin.Context) {
	id := c.Query(logic.Listener)
	err := logic.BizDeleteListener(id)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}

	c.JSON(http.StatusOK, config.WithRet("delete Listener success!"))
}

// DetailListener get Listener detail
func DetailListener(c *gin.Context) {
	name := c.Query(logic.Listener)
	res, err := logic.BizGetListener(name)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet(res))
}

// UpdateListener update Listener
func UpdateListener(c *gin.Context) {
	body := c.PostForm("content")
	res := &fc.Listener{}
	err := yaml.UnmarshalYML([]byte(body), res)
	logger.Debug(body)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	err = logic.BizUpdateListener(res)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("update Listener success!"))
}
