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
