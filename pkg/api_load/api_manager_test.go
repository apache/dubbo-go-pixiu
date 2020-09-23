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
package api_load

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service/api"
)

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestApiLoad_StartLoadApi(t *testing.T) {
}

func TestApiLoad_AddApiLoad(t *testing.T) {
	ads := api.NewLocalMemoryApiDiscoveryService()
	apiManager := NewApiManager(time.Second*5, ads)
	apiManager.AddApiLoader(model.ApiConfig{
		File: &model.File{
			FileApiConfPath: filePath,
		},
	})
	apiLoader, err := apiManager.GetApiLoad(File)
	assert.Nil(t, err)
	assert.NotNil(t, apiLoader)

}

func TestApiManager_StartLoadApi(t *testing.T) {
	ads := api.NewLocalMemoryApiDiscoveryService()
	apiManager := NewApiManager(time.Second, ads)
	apiManager.AddApiLoader(model.ApiConfig{
		File: &model.File{
			FileApiConfPath: filePath,
		},
	})
	err := apiManager.StartLoadApi()
	assert.Nil(t, err)
	go apiManager.SelectMergeApiTask()
	appendApiName := "test_append_api"
	appendApi := model.Api{
		Name:     appendApiName,
		ITypeStr: "HTTP",
		OTypeStr: "REST",
		Method:   "GET",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "BDTService",
				Group:           "test1",
				Version:         "2.0.0",
				Interface:       "com.ikurento.user.UserProvider",
				Method:          "GetUser",
				Types: []string{
					"java.lang.String",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	bytes, err := yaml.Marshal([]model.Api{appendApi})
	assert.Nil(t, err)
	originContent, _ := ioutil.ReadFile(filePath)
	err = appendContent(filePath, bytes)
	ok := assert.Nil(t, err)
	if ok {
		defer ioutil.WriteFile(filePath, originContent, os.ModePerm)
	}
	j, _ := json.Marshal(appendApi)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	select {
	case <-ctx.Done():
		_, err = apiManager.ads.GetApi(*service.NewDiscoveryRequest(j))
		assert.Nil(t, err)
		return
	}
}
