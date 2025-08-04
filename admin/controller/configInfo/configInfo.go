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
	"net/http"
	"strconv"
	"strings"
)

import (
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
)

import (
	config2 "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// GetBaseInfo get pixiu base info such as name,desc
func GetBaseInfo(c *gin.Context) {
	conf, err := logic.BizGetBaseInfo()
	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	data, _ := yaml.MarshalYML(conf)
	c.JSON(http.StatusOK, config2.WithRet(string(data)))
}

// SetBaseInfo modify pixiu base info such as name,desc
func SetBaseInfo(c *gin.Context) {
	body := c.PostForm("content")

	baseInfo := &config2.BaseInfo{}
	err := yaml.UnmarshalYML([]byte(body), baseInfo)
	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}

	setErr := logic.BizSetBaseInfo(baseInfo, true)
	if setErr != nil {
		c.JSON(http.StatusOK, config2.WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet("success"))
}

// GetResourceList get all resource list
func GetResourceList(c *gin.Context) {
	unpublished := getUnpublishedVal(c)

	res, err := logic.BizGetResourceList(unpublished)
	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, config2.WithRet(string(data)))
}

// GetResourceDetail get resource detail with yml
func GetResourceDetail(c *gin.Context) {
	unpublished := getUnpublishedVal(c)
	id := c.Query(logic.ResourceID)
	res, err := logic.BizGetResourceDetail(id, unpublished)
	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet(res))
}

// CreateResourceInfo create resource
func CreateResourceInfo(c *gin.Context) {
	body := c.PostForm("content")
	unpublished := getUnpublishedVal(c)

	res := &fc.Resource{}
	err := yaml.UnmarshalYML([]byte(body), res)
	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}

	var setErr1, setErr2 error // err1 represent write publish space, err2 represent write unpublished space
	if unpublished {
		setErr2 = logic.BizSetResourceInfo(res, true, true)
	} else {
		setErr1 = logic.BizSetResourceInfo(res, true, false)
		setErr2 = logic.BizSetResourceInfo(res, true, true)
	}

	//setErr := logic.BizSetResourceInfo(res, true, unpublished)
	if setErr1 != nil {
		c.JSON(http.StatusOK, config2.WithError(setErr1))
		return
	}
	if setErr2 != nil {
		c.JSON(http.StatusOK, config2.WithError(setErr2))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet("Success"))
}

// ModifyResourceInfo modify resource
func ModifyResourceInfo(c *gin.Context) {
	id := c.Query(logic.ResourceID)
	body := c.PostForm("content")
	unpublished := getUnpublishedVal(c)

	res := &fc.Resource{}
	err := yaml.UnmarshalYML([]byte(body), res)
	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}

	if id != "" {
		res.ID, err = strconv.Atoi(id)
		if err != nil {
			logger.Warnf("resourceID not number err, %v\n", err)
			c.JSON(http.StatusOK, config2.WithError(err))
			return
		}
	}

	old, err := getResourceDetail(id, unpublished)
	if err == nil && old != nil {
		// when resource path change, should modify all method below it
		if old.Path != res.Path {
			afterResourcePathChange(id, res.Path, unpublished)
		}
	}

	setErr := logic.BizSetResourceInfo(res, false, unpublished)
	if setErr != nil {
		c.JSON(http.StatusOK, config2.WithError(setErr))
		return
	}

	c.JSON(http.StatusOK, config2.WithRet("Success"))
}

func afterResourcePathChange(resourceId, path string, unpublished bool) {
	mList, err := logic.BizGetMethodList(resourceId, unpublished)
	if err != nil {
		return
	}
	for i := range mList {
		m := &mList[i]
		m.ResourcePath = path
		setErr := logic.BizSetResourceMethod(resourceId, m, false, unpublished)
		if setErr != nil {
			logger.Warnf("afterResourcePathChange err, %v\n", err)
			continue
		}
	}
}

// DeleteResourceInfo delete resource
func DeleteResourceInfo(c *gin.Context) {
	id := c.Query(logic.ResourceID)
	unpublished := getUnpublishedVal(c)
	if unpublished {
		// Check whether the configuration has been released when deleting the configuration
		old, err := getResourceDetail(id, false)
		if err != nil {
			c.JSON(http.StatusOK, config2.WithError(err))
			return
		}
		if old != nil {
			c.JSON(http.StatusOK, config2.WithError(errors.New("The configuration has been published and cannot be deleted")))
			return
		}
	}
	err := logic.BizDeleteResourceInfo(id, unpublished)
	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}

	c.JSON(http.StatusOK, config2.WithRet("Success"))
}

// GetMethodList get all method list below one resource
func GetMethodList(c *gin.Context) {
	resourceId := c.Query(logic.ResourceID) // unique id
	unpublished := getUnpublishedVal(c)

	res, err := logic.BizGetMethodList(resourceId, unpublished)
	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, config2.WithRet(string(data)))
}

// GetMethodDetail get method detail with yml
func GetMethodDetail(c *gin.Context) {
	resourceId := c.Query(logic.ResourceID)
	methodId := c.Query(logic.MethodID) // unique id
	unpublished := getUnpublishedVal(c)
	res, err := logic.BizGetMethodDetail(resourceId, methodId, unpublished)
	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet(res))
}

// DeleteResourceInfo delete method
func DeleteMethodInfo(c *gin.Context) {
	resourceId := c.Query(logic.ResourceID)
	methodId := c.Query(logic.MethodID)
	unpublished := getUnpublishedVal(c)
	if unpublished {
		old, err := logic.BizGetMethodDetail(resourceId, methodId, false)
		if err != nil {
			c.JSON(http.StatusOK, config2.WithError(err))
			return
		}
		if old != "" {
			c.JSON(http.StatusOK, config2.WithError(errors.New("The configuration has been published and cannot be deleted")))
			return
		}
	}
	err := logic.BizDeleteMethodInfo(resourceId, methodId, unpublished)

	if err != nil {
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet("Success"))
}

// CreateMethodInfo create method
func CreateMethodInfo(c *gin.Context) {
	body := c.PostForm("content")
	resourceId := c.Query(logic.ResourceID)
	unpublished := getUnpublishedVal(c)

	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}

	resource, err := getResourceDetail(resourceId, unpublished)
	if err != nil {
		logger.Warnf("CreateMethodInfo can't query resource  err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	res.ResourcePath = resource.Path

	setErr := logic.BizSetResourceMethod(resourceId, res, true, unpublished)
	if setErr != nil {
		c.JSON(http.StatusOK, config2.WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet("Success"))
}

func getResourceDetail(id string, unpublished bool) (*fc.Resource, error) {
	res, err := logic.BizGetResourceDetail(id, unpublished)
	if err != nil {
		return nil, err
	}

	resource := &fc.Resource{}
	err = yaml.UnmarshalYML([]byte(res), resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// ModifyMethodInfo modify method
func ModifyMethodInfo(c *gin.Context) {
	body := c.PostForm("content")
	resourceId := c.Query(logic.ResourceID)
	methodId := c.Query(logic.MethodID)
	unpublished := getUnpublishedVal(c)

	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(body), res)
	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}

	if methodId != "" {
		res.ID, err = strconv.Atoi(methodId)
		if err != nil {
			logger.Warnf("methodID not number err, %v\n", err)
			c.JSON(http.StatusOK, config2.WithError(err))
			return
		}
	}

	resource, err := getResourceDetail(resourceId, unpublished)
	if err != nil {
		logger.Warnf("CreateMethodInfo can't query resource  err, %v\n", err)
		c.JSON(http.StatusOK, config2.WithError(err))
		return
	}
	res.ResourcePath = resource.Path

	setErr := logic.BizSetResourceMethod(resourceId, res, false, unpublished)
	if setErr != nil {
		c.JSON(http.StatusOK, config2.WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, config2.WithRet("Success"))
}

// getUnpublishedVal Determine the configuration type of the current operation
func getUnpublishedVal(c *gin.Context) bool {
	// The front-end request carries the unpublished field to determine which configuration is currently operating
	// 1 represent true (unpublished, delay publish), 0 represent false (published, direct publish)
	unpublishedVal := c.PostForm("unpublished")
	if strings.EqualFold(unpublishedVal, "1") {
		return true
	} else {
		return false
	}
}

// BatchReleaseResource Publish all configuration information
func BatchReleaseResource(c *gin.Context) {
	fromKList, fromVList, fromErr := logic.BRGetResourceList(true) // from represent unpublished space
	toKList, toVList, _ := logic.BRGetResourceList(false)          // to represent published space
	// Do not handle toList errors
	if fromErr != nil {
		logger.Warnf("Batch Release Resource err, %v\n", fromErr)
		c.JSON(http.StatusOK, config2.WithError(fromErr))
		return
	}
	// todo Optimize comparison method to reduce time complexity
	for i, fromK := range fromKList {
		fromV := fromVList[i]
		fromKTmp := strings.Split(fromK, "/")
		flag := false
		for j, toK := range toKList {
			toV := toVList[j]
			toKTmp := strings.Split(toK, "/")
			flag = strings.EqualFold(fromKTmp[len(fromKTmp)-1], toKTmp[len(toKTmp)-1])
			if flag {
				if !strings.EqualFold(fromV, toV) {
					err := logic.BRUpdate(toK, fromV)
					if err != nil {
						logger.Warnf("Batch Release Resource err, %v\n", err)
						c.JSON(http.StatusOK, config2.WithError(err))
						return
					}
				}
				break
			}
		}
		if !flag {
			err := logic.BRCreate(fromKTmp[len(fromKTmp)-1], fromV, logic.Resources)
			if err != nil {
				logger.Warnf("Batch Release Resource err, %v\n", err)
				c.JSON(http.StatusOK, config2.WithError(err))
				return
			}
		}
	}
}

// BatchReleaseMethod Batch Release Method Config
func BatchReleaseMethod(c *gin.Context) {
	// todo
}

// BatchReleasePluginGroup Batch Release PluginGroup Config
func BatchReleasePluginGroup(c *gin.Context) {
	fromKList, fromVList, fromErr := logic.BRGetPluginGroupList(true) // from represent unpublished space
	toKList, toVList, _ := logic.BRGetPluginGroupList(false)          // to represent published space
	if fromErr != nil {
		logger.Warnf("Batch Release PluginGroup err, %v\n", fromErr)
		c.JSON(http.StatusOK, config2.WithError(fromErr))
		return
	}
	fromKTmp := strings.Split(fromKList[0], "/")
	if toKList == nil {
		err := logic.BRCreate(fromKTmp[len(fromKTmp)-1], fromVList[0], logic.PluginGroup)
		if err != nil {
			logger.Warnf("Batch Release PluginGroup err, %v\n", err)
			c.JSON(http.StatusOK, config2.WithError(err))
		}
		return
	}
	if !strings.EqualFold(fromVList[0], toVList[0]) {
		err := logic.BRUpdate(toKList[0], fromVList[0])
		if err != nil {
			logger.Warnf("Batch Release PluginGroup err, %v\n", err)
			c.JSON(http.StatusOK, config2.WithError(err))
		}
	}
}
