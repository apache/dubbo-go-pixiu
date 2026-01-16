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
	"strconv"
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/model"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/i18n"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	pkgconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	pkgmodel "github.com/apache/dubbo-go-pixiu/pkg/model"
)

type ContentRequest struct {
	Content string `json:"content"`
}

func getContentFromRequest(c *gin.Context) string {
	if strings.Contains(c.ContentType(), "application/json") {
		var req ContentRequest
		if err := c.ShouldBindJSON(&req); err == nil && req.Content != "" {
			return req.Content
		}
	}
	return c.PostForm("content")
}

func (h *Handler) ListClusters(c *gin.Context) {
	clusters, err := h.etcd.ListClusters()
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, clusters)
}

func (h *Handler) GetCluster(c *gin.Context) {
	name := c.Param("name")
	content, err := h.etcd.GetCluster(name)
	if err != nil {
		config.Error(c, err)
		return
	}
	if c.Query("format") == "yaml" {
		config.Success(c, content)
		return
	}
	var cluster pkgmodel.ClusterConfig
	if err := yaml.UnmarshalYML([]byte(content), &cluster); err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, cluster)
}

func (h *Handler) CreateCluster(c *gin.Context) {
	body := getContentFromRequest(c)
	var cluster pkgmodel.ClusterConfig
	if err := yaml.UnmarshalYML([]byte(body), &cluster); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.CreateCluster(&cluster); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgClusterCreated)
}

func (h *Handler) UpdateCluster(c *gin.Context) {
	body := getContentFromRequest(c)
	var cluster pkgmodel.ClusterConfig
	if err := yaml.UnmarshalYML([]byte(body), &cluster); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.UpdateCluster(&cluster); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgClusterUpdated)
}

func (h *Handler) DeleteCluster(c *gin.Context) {
	name := c.Param("name")
	if err := h.etcd.DeleteCluster(name); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgClusterDeleted)
}

func (h *Handler) ListListeners(c *gin.Context) {
	listeners, err := h.etcd.ListListeners()
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, listeners)
}

func (h *Handler) GetListener(c *gin.Context) {
	name := c.Param("name")
	content, err := h.etcd.GetListener(name)
	if err != nil {
		config.Error(c, err)
		return
	}
	if c.Query("format") == "yaml" {
		config.Success(c, content)
		return
	}
	var listener pkgmodel.Listener
	if err := yaml.UnmarshalYML([]byte(content), &listener); err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, listener)
}

func (h *Handler) CreateListener(c *gin.Context) {
	body := getContentFromRequest(c)
	var listener pkgmodel.Listener
	if err := yaml.UnmarshalYML([]byte(body), &listener); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.CreateListener(&listener); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgListenerCreated)
}

func (h *Handler) UpdateListener(c *gin.Context) {
	body := getContentFromRequest(c)
	var listener pkgmodel.Listener
	if err := yaml.UnmarshalYML([]byte(body), &listener); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.UpdateListener(&listener); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgListenerUpdated)
}

func (h *Handler) DeleteListener(c *gin.Context) {
	name := c.Param("name")
	if err := h.etcd.DeleteListener(name); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgListenerDeleted)
}

func (h *Handler) ListResources(c *gin.Context) {
	resources, err := h.etcd.ListResources(false)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, resources)
}

func (h *Handler) GetResource(c *gin.Context) {
	id := c.Param("id")
	content, err := h.etcd.GetResource(id, false)
	if err != nil {
		config.Error(c, err)
		return
	}
	if c.Query("format") == "yaml" {
		config.Success(c, content)
		return
	}
	var resource pkgconfig.Resource
	if err := yaml.UnmarshalYML([]byte(content), &resource); err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, resource)
}

func (h *Handler) CreateResource(c *gin.Context) {
	body := getContentFromRequest(c)
	var resource pkgconfig.Resource
	if err := yaml.UnmarshalYML([]byte(body), &resource); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.CreateResource(&resource, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgResourceCreated)
}

func (h *Handler) UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		config.ErrorMessage(c, "invalid resource ID")
		return
	}
	body := getContentFromRequest(c)
	var resource pkgconfig.Resource
	if err := yaml.UnmarshalYML([]byte(body), &resource); err != nil {
		config.Error(c, err)
		return
	}
	resource.ID = id
	if err := h.etcd.UpdateResource(&resource, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgResourceUpdated)
}

func (h *Handler) DeleteResource(c *gin.Context) {
	id := c.Param("id")
	if err := h.etcd.DeleteResource(id, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgResourceDeleted)
}

func (h *Handler) ListMethods(c *gin.Context) {
	resourceID := c.Query("resourceId")
	methods, err := h.etcd.ListMethods(resourceID, false)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, methods)
}

func (h *Handler) GetMethod(c *gin.Context) {
	resourceID := c.Query("resourceId")
	methodID := c.Param("id")
	method, err := h.etcd.GetMethod(resourceID, methodID, false)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, method)
}

func (h *Handler) CreateMethod(c *gin.Context) {
	resourceID := c.Query("resourceId")
	body := getContentFromRequest(c)
	var method pkgconfig.Method
	if err := yaml.UnmarshalYML([]byte(body), &method); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.CreateMethod(resourceID, &method, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgMethodCreated)
}

func (h *Handler) UpdateMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		config.ErrorMessage(c, "invalid method ID")
		return
	}
	resourceID := c.Query("resourceId")
	body := getContentFromRequest(c)
	var method pkgconfig.Method
	if err := yaml.UnmarshalYML([]byte(body), &method); err != nil {
		config.Error(c, err)
		return
	}
	method.ID = id
	if err := h.etcd.UpdateMethod(resourceID, &method, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgMethodUpdated)
}

func (h *Handler) DeleteMethod(c *gin.Context) {
	resourceID := c.Query("resourceId")
	methodID := c.Param("id")
	if err := h.etcd.DeleteMethod(resourceID, methodID, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgMethodDeleted)
}

func (h *Handler) ListPluginGroups(c *gin.Context) {
	groups, err := h.etcd.ListPluginGroups(false)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, groups)
}

func (h *Handler) GetPluginGroup(c *gin.Context) {
	name := c.Param("name")
	content, err := h.etcd.GetPluginGroup(name, false)
	if err != nil {
		config.Error(c, err)
		return
	}
	if c.Query("format") == "yaml" {
		config.Success(c, content)
		return
	}
	var group model.PluginGroup
	if err := yaml.UnmarshalYML([]byte(content), &group); err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, group)
}

func (h *Handler) CreatePluginGroup(c *gin.Context) {
	body := getContentFromRequest(c)
	var group model.PluginGroup
	if err := yaml.UnmarshalYML([]byte(body), &group); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.CreatePluginGroup(&group, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgPluginCreated)
}

func (h *Handler) UpdatePluginGroup(c *gin.Context) {
	body := getContentFromRequest(c)
	var group model.PluginGroup
	if err := yaml.UnmarshalYML([]byte(body), &group); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.etcd.UpdatePluginGroup(&group, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgPluginUpdated)
}

func (h *Handler) DeletePluginGroup(c *gin.Context) {
	name := c.Param("name")
	if err := h.etcd.DeletePluginGroup(name, false); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgPluginDeleted)
}
