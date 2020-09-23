package api_load

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service/api"
	"testing"
	"time"
)

func TestApiLoad_StartLoadApi(t *testing.T) {
}

func TestApiLoad_AddApiLoad(t *testing.T) {
	ads := api.NewLocalMemoryApiDiscoveryService()
	apiManager := NewApiManager(time.Second*5, ads)
	apiManager.AddApiLoad(model.ApiConfig{
		File: &model.File{
			FileApiConfPath: filePath,
		},
	})

}
