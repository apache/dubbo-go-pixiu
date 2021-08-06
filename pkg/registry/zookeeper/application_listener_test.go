package zookeeper

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/service/api"
	"github.com/dubbogo/go-zookeeper/zk"
	"testing"
	"time"
)

func TestAppli(t *testing.T) {
	api.Init()
	//tc, err := zk.StartTestCluster(1, nil, nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer tc.Stop()

	regConfig := model.Registry{
		Protocol: "zookeeper",
		Timeout:  "20s",
		Address:  "127.0.0.1:2181",
	}
	config.Load("D:\\codes\\dubbo-go-pixiu\\samples\\dubbogo\\simple\\query\\pixiu\\conf.yaml")
	dubbo.SingletonDubboClient().Init()

	reg, _ := newZKRegistry(regConfig)
	r := reg.(*ZKRegistry)
	c := r.GetClient()
	connState := c.GetConnState()
	for connState != zk.StateConnected && connState != zk.StateHasSession {
		<-time.After(200 * time.Millisecond)
		connState = c.GetConnState()
	}

	time.Sleep(30 * time.Second)
}
