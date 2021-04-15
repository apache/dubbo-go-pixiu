package main

import (
	getty "github.com/apache/dubbo-getty"
	"github.com/coreos/etcd/embed"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
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
