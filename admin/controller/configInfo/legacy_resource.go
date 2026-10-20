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
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	commonyaml "github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

// These handlers preserve the Resource/Method Admin API for clients deployed
// before AdminRouteBinding. They continue to read/write the legacy runtime
// etcd keys, so old YAML and published records remain usable.
func GetResourceList(c *gin.Context) {
	resources, err := logic.BizGetResourceList(legacyUnpublished(c))
	if err != nil {
		legacyError(c, err)
		return
	}
	data, _ := json.Marshal(resources)
	c.JSON(http.StatusOK, adminconfig.WithRet(string(data)))
}

func GetResourceDetail(c *gin.Context) {
	detail, err := logic.BizGetResourceDetail(c.Query(logic.ResourceID), legacyUnpublished(c))
	if err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(detail))
}

func CreateResourceInfo(c *gin.Context) {
	resource, err := decodeLegacyResource(c)
	if err != nil {
		legacyError(c, err)
		return
	}
	// Creation always assigns server-owned identities; positive IDs are only
	// accepted by the logic layer when mirroring the same records to draft.
	resource.ID = 0
	for index := range resource.Methods {
		resource.Methods[index].ID = 0
	}
	if legacyUnpublished(c) {
		err = logic.BizSetResourceInfo(resource, true, true)
	} else {
		err = logic.BizSetResourceInfo(resource, true, false)
		if err == nil {
			// Mirror the published records into the draft namespace with the
			// same Resource/Method IDs. Allocating a second set here would make
			// subsequent edits and publishes address a different route.
			err = logic.BizSetResourceInfo(resource, false, true)
		}
		if err == nil {
			resourceID := strconv.Itoa(resource.ID)
			for index := range resource.Methods {
				if err = logic.BizSetResourceMethod(resourceID, &resource.Methods[index], false, true); err != nil {
					break
				}
			}
		}
	}
	if err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

func ModifyResourceInfo(c *gin.Context) {
	resource, err := decodeLegacyResource(c)
	if err != nil {
		legacyError(c, err)
		return
	}
	id := strings.TrimSpace(c.Query(logic.ResourceID))
	if id != "" {
		resource.ID, err = strconv.Atoi(id)
		if err != nil {
			legacyError(c, err)
			return
		}
	}
	unpublished := legacyUnpublished(c)
	oldContent, oldErr := logic.BizGetResourceDetail(strconv.Itoa(resource.ID), unpublished)
	var old config.Resource
	if oldErr == nil {
		_ = commonyaml.UnmarshalYML([]byte(oldContent), &old)
	}
	if err := logic.BizSetResourceInfo(resource, false, unpublished); err != nil {
		legacyError(c, err)
		return
	}
	if oldErr == nil && old.Path != resource.Path {
		methods, err := logic.BizGetMethodList(strconv.Itoa(resource.ID), unpublished)
		if err == nil {
			for index := range methods {
				methods[index].ResourcePath = resource.Path
				if err := logic.BizSetResourceMethod(strconv.Itoa(resource.ID), &methods[index], false, unpublished); err != nil {
					legacyError(c, err)
					return
				}
			}
		}
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

func DeleteResourceInfo(c *gin.Context) {
	id := c.Query(logic.ResourceID)
	resourceID, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil || resourceID <= 0 {
		legacyError(c, errors.New("resourceId must be a positive integer"))
		return
	}
	id = strconv.Itoa(resourceID)
	unpublished := legacyUnpublished(c)
	if unpublished {
		published, err := logic.BizResourceExists(id, false)
		if err != nil {
			legacyError(c, err)
			return
		}
		if published {
			legacyError(c, errors.New("The configuration has been published and cannot be deleted"))
			return
		}
	}
	if err := logic.BizDeleteResourceInfo(id, unpublished); err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

func GetMethodList(c *gin.Context) {
	methods, err := logic.BizGetMethodList(c.Query(logic.ResourceID), legacyUnpublished(c))
	if err != nil {
		legacyError(c, err)
		return
	}
	data, _ := json.Marshal(methods)
	c.JSON(http.StatusOK, adminconfig.WithRet(string(data)))
}

func GetMethodDetail(c *gin.Context) {
	detail, err := logic.BizGetMethodDetail(c.Query(logic.ResourceID), c.Query(logic.MethodID), legacyUnpublished(c))
	if err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(detail))
}

func CreateMethodInfo(c *gin.Context) { saveLegacyMethod(c, true) }

func ModifyMethodInfo(c *gin.Context) { saveLegacyMethod(c, false) }

func saveLegacyMethod(c *gin.Context, create bool) {
	var method config.Method
	if err := commonyaml.UnmarshalYML([]byte(c.PostForm("content")), &method); err != nil {
		legacyError(c, err)
		return
	}
	resourceID := c.Query(logic.ResourceID)
	resourceContent, err := logic.BizGetResourceDetail(resourceID, legacyUnpublished(c))
	if err != nil {
		legacyError(c, err)
		return
	}
	var resource config.Resource
	if err := commonyaml.UnmarshalYML([]byte(resourceContent), &resource); err != nil {
		legacyError(c, err)
		return
	}
	method.ResourcePath = resource.Path
	if !create {
		methodID, err := strconv.Atoi(c.Query(logic.MethodID))
		if err != nil {
			legacyError(c, err)
			return
		}
		method.ID = methodID
	} else {
		method.ID = 0
	}
	if err := logic.BizSetResourceMethod(resourceID, &method, create, legacyUnpublished(c)); err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

func DeleteMethodInfo(c *gin.Context) {
	resourceID, methodID := c.Query(logic.ResourceID), c.Query(logic.MethodID)
	resourceNumber, resourceErr := strconv.Atoi(strings.TrimSpace(resourceID))
	methodNumber, methodErr := strconv.Atoi(strings.TrimSpace(methodID))
	if resourceErr != nil || resourceNumber <= 0 || methodErr != nil || methodNumber <= 0 {
		legacyError(c, errors.New("resourceId and methodId must be positive integers"))
		return
	}
	resourceID, methodID = strconv.Itoa(resourceNumber), strconv.Itoa(methodNumber)
	unpublished := legacyUnpublished(c)
	if unpublished {
		published, err := logic.BizMethodExists(resourceID, methodID, false)
		if err != nil {
			legacyError(c, err)
			return
		}
		if published {
			legacyError(c, errors.New("The configuration has been published and cannot be deleted"))
			return
		}
	}
	if err := logic.BizDeleteMethodInfo(resourceID, methodID, unpublished); err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

func BatchReleaseResource(c *gin.Context) { releaseLegacyResources(c) }

func BatchReleaseMethod(c *gin.Context) { releaseLegacyResources(c) }

func releaseLegacyResources(c *gin.Context) {
	if err := logic.BizBatchReleaseResources(); err != nil {
		legacyError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

func decodeLegacyResource(c *gin.Context) (*config.Resource, error) {
	var resource config.Resource
	if err := commonyaml.UnmarshalYML([]byte(c.PostForm("content")), &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func legacyUnpublished(c *gin.Context) bool {
	value := c.PostForm("unpublished")
	if value == "" {
		value = c.Query("unpublished")
	}
	return strings.EqualFold(strings.TrimSpace(value), "1")
}

func legacyError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, adminconfig.WithError(err))
}
