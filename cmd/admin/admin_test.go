package main

import (
	"go.etcd.io/etcd/embed"
	"os"
	"testing"
	"time"
)

import (
	getty "github.com/apache/dubbo-getty"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
)

const defaultEtcdV3WorkDir = "/tmp/default-pixiu-admin.etcd"

type AdminTestSuite struct {
	suite.Suite
	etcd *embed.Etcd
}

// start etcd server
func (suite *AdminTestSuite) SetupSuite() {

	t := suite.T()

	cfg := embed.NewConfig()
	// avoid conflict with default etcd work-dir
	cfg.Dir = defaultEtcdV3WorkDir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-e.Server.ReadyNotify():
		t.Log("Server is ready!")
	case <-getty.GetTimeWheel().After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		t.Logf("Server took too long to start!")
	}

	suite.etcd = e

	gin.SetMode(gin.TestMode)
}

// stop etcd server
func (suite *AdminTestSuite) TearDownSuite() {
	suite.etcd.Close()
	// clean the etcd workdir
	if err := os.RemoveAll(defaultEtcdV3WorkDir); err != nil {
		suite.FailNow(err.Error())
	}
}

func TestAdminSuite(t *testing.T) {
	suite.Run(t, &AdminTestSuite{})
}

func (suite *AdminTestSuite) TestSetBaseInfo() {

}

func (suite *AdminTestSuite) TestGetBaseInfo() {

}

func TestGetBaseInfo(t *testing.T) {
	str := "httpVerb: GET\nonAir: true\ntimeout: 1000ms\ninboundRequest:\n    requestType: http\n    queryStrings:\n    - name: name\n        required: true\nintegrationRequest:\n    requestType: http\n    host: 127.0.0.1:8889\n    path: /UserProvider/GetUserByName\n    mappingParams:\n    - name: queryStrings.name\n        mapTo: queryStrings.name\n    group: \"test\"\n    version: 1.0.0"
	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(str), res)
	if err != nil {

	}
}
