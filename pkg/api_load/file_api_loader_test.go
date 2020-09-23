package api_load

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

import (
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

const filePath = "./mock/api_config.yaml"

func newFileLoader() ApiLoader {
	return &FileApiLoader{
		filePath: filePath,
	}
}

func TestFileApiLoader_InitLoad(t *testing.T) {
	fileApiLoader := newFileLoader()
	err := fileApiLoader.InitLoad()
	assert.Nil(t, err)
	apis, err := fileApiLoader.GetLoadedApiConfigs()
	assert.Nil(t, err)
	assert.NotNil(t, apis)
	bytes, err := yaml.Marshal(apis)
	assert.Nil(t, err)
	content, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	assert.YAMLEq(t, string(bytes), string(content))
}

func TestFileApiLoader_GetPrior(t *testing.T) {
	fileApiLoader := newFileLoader()
	prior := fileApiLoader.GetPrior()
	assert.Equal(t, 0, prior)
}

func TestFileApiLoader_HotLoad(t *testing.T) {
	fileApiLoader := newFileLoader()
	fileChangeNotifier, err := fileApiLoader.HotLoad()
	assert.Nil(t, err)
	assert.NotNil(t, fileChangeNotifier)
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
	t.Log(string(bytes))
	err = appendContent(filePath, bytes)
	ok := assert.Nil(t, err)
	if ok {
		defer ioutil.WriteFile(filePath, originContent, os.ModePerm)
	}
	select {
	case <-fileChangeNotifier:
		apis, err := fileApiLoader.GetLoadedApiConfigs()
		assert.Nil(t, err)
		assert.Len(t, apis, 2)
		for _, api := range apis {
			if api.Name == appendApiName {
				assert.True(t, true)
				return
			}
		}
		assert.True(t, false)
		break
	case <-time.After(time.Second * 10000):
		assert.Fail(t, "did't get change notice!")
	}
}

func appendContent(filename string, content []byte) error {
	fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer fl.Close()
	n, err := fl.Write(content)
	if err == nil && n < len(content) {
		err = io.ErrShortWrite
	}
	return err
}
