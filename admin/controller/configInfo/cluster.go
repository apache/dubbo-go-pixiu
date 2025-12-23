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
// @Summary get cluster list
// @Description get all clusters' info
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/cluster/list [get]
// GetClusterList get all cluster list
func GetClusterList(c *gin.Context) {
	rst, err := logic.BizGetClusters()
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(rst))
}

// @Tags Config
// @Summary create cluster
// @Description Create a cluster by passing the YAML/JSON configuration for the cluster through the form's content field.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param content formData string true "Cluster content"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /config/api/cluster [put]
// CreateCluster create a cluster
func CreateCluster(c *gin.Context) {
	body := c.PostForm("content")
	res := &config.Cluster{}
	err := yaml.UnmarshalYML([]byte(body), res)
	logger.Debug(body)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	err = logic.BizCreateCluster(res)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("create cluster success!"))
}

// @Tags Config
// @Summary delete cluster
// @Description delete cluster according to cluster ID
// @Produce application/json
// @Param id query string true "Cluster ID"
// @Success 200 {object} string
// @Router /config/api/cluster [delete]
// DeleteCluster delete resource
func DeleteCluster(c *gin.Context) {
	id := c.Query(logic.ClusterID)
	err := logic.BizDeleteCluster(id)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}

	c.JSON(http.StatusOK, adminconfig.WithRet("delete cluster success!"))
}

// @Tags Config
// @Summary get cluster detail
// @Description get cluster detail according to cluster ID
// @Produce application/json
// @Param id query string true "Cluster ID"
// @Success 200 {object} string
// @Router /config/api/cluster/detail [get]
// DetailCluster get cluster detail
func DetailCluster(c *gin.Context) {
	id := c.Query(logic.ClusterID)
	res, err := logic.BizGetCluster(id)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(res))
}

// @Tags Config
// @Summary update cluster
// @Description pass the Cluster's YAML/JSON via the form's content field to update the cluster.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param content formData string true "Cluster content"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /config/api/cluster [post]
// UpdateCluster update cluster
func UpdateCluster(c *gin.Context) {
	body := c.PostForm("content")
	res := &config.Cluster{}
	err := yaml.UnmarshalYML([]byte(body), res)
	logger.Debug(body)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	err = logic.BizUpdateCluster(res)
	if err != nil {
		c.JSON(http.StatusOK, adminconfig.WithError(err))
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("update cluster success!"))
}
