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

package zookeeper

import (
	"testing"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type DemoListener struct {
}

func (d *DemoListener) OnAddAPI(r router.API) error {
	return nil
}

func (d *DemoListener) OnRemoveAPI(r router.API) error {
	return nil
}

func (d *DemoListener) OnDeleteRouter(r config.Resource) error {
	return nil
}

func TestNewZKRegistry(t *testing.T) {

	regConfig := model.Registry{
		Protocol: "zookeeper",
		Timeout:  "2s",
		Address:  "127.0.0.1:9100",
	}
	reg, err := newZKRegistry(regConfig, &DemoListener{})
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	regConfig = model.Registry{
		Protocol: "zookeeper",
		Timeout:  "2xxxxxx",
		Address:  "127.0.0.1:9100",
	}
	reg, err = newZKRegistry(regConfig, &DemoListener{})
	assert.Nil(t, reg)
	assert.NotNil(t, err)

	regConfig = model.Registry{
		Protocol: "zookeeper",
		Timeout:  "2s",
		Address:  "",
	}
	reg, err = newZKRegistry(regConfig, &DemoListener{})
	assert.Nil(t, reg)
	assert.NotNil(t, err)
}

//
//func TestLoadInterfaces(t *testing.T) {
//	api.Init()
//	tc, err := zk.StartTestCluster(1, nil, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer tc.Stop()
//
//	conn, _, _ := tc.ConnectAll()
//	conn.Create("/dubbo", nil, 0, zk.WorldACL(zk.PermAll))
//	conn.Create("/dubbo/com.ikurento.user.UserProvider1", nil, 0, zk.WorldACL(zk.PermAll))
//	conn.Create("/dubbo/com.ikurento.user.UserProvider1/providers", nil, 0, zk.WorldACL(zk.PermAll))
//	conn.Create("/dubbo/com.ikurento.user.UserProvider1/providers/"+
//		url.QueryEscape("dubbo://192.168.3.46:20002/org.apache.dubbo.UserProvider1?anyhost=true&app.version=0.0.1&application=UserInfoServer1&bean.name=UserProvider&cluster=failover&interface=org.apache.dubbo.UserProvider1&loadbalance=random&methods=GetUser&name=UserInfoServer&organization=dubbo.io&registry.role=3&service.filter=echo,token,accesslog,tps,generic_service,execute,pshutdown&side=provider&ssl-enabled=false&timestamp=1624716984&warmup=100"),
//		nil, 0, zk.WorldACL(zk.PermAll))
//	conn.Create("/dubbo/com.ikurento.user.UserProvider2", nil, 0, zk.WorldACL(zk.PermAll))
//	conn.Create("/dubbo/com.ikurento.user.UserProvider2/providers", nil, 0, zk.WorldACL(zk.PermAll))
//	conn.Create("/dubbo/com.ikurento.user.UserProvider2/providers/"+
//		url.QueryEscape("dubbo://192.168.3.46:20001/org.apache.dubbo.UserProvider2?anyhost=true&app.version=0.0.1&application=UserInfoServer2&bean.name=UserProvider&cluster=failover&interface=org.apache.dubbo.UserProvider2&loadbalance=random&methods=GetUser,SetUser,AllUsers&name=UserInfoServer&organization=dubbo.io&registry.role=3&service.filter=echo,token,accesslog,tps,generic_service,execute,pshutdown&side=provider&ssl-enabled=false&timestamp=1624716984&warmup=100"),
//		nil, 0, zk.WorldACL(zk.PermAll))
//
//	hosts := make([]string, len(tc.Servers))
//	for i, srv := range tc.Servers {
//		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
//	}
//	regConfig := model.Registry{
//		Protocol: "zookeeper",
//		Timeout:  "20s",
//		Address:  hosts[0],
//	}
//	reg, _ := newZKRegistry(regConfig)
//	r := reg.(*ZKRegistry)
//	c := r.GetClient()
//	connState := c.GetConnState()
//	for connState != zk.StateConnected && connState != zk.StateHasSession {
//		<-time.After(200 * time.Millisecond)
//		connState = c.GetConnState()
//	}
//	r.LoadInterfaces()
//	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
//	api, err := localAPIDiscSrv.GetAPI("/UserInfoServer1/org.apache.dubbo.UserProvider1/0.0.1/GetUser", config.MethodPost)
//	assert.Equal(t, api.URLPattern, "/userinfoserver1/org.apache.dubbo.userprovider1/0.0.1/getuser")
//	assert.Nil(t, err)
//	api, err = localAPIDiscSrv.GetAPI("/UserInfoServer2/org.apache.dubbo.UserProvider2/0.0.1/GetUser", config.MethodPost)
//	assert.Equal(t, api.URLPattern, "/userinfoserver2/org.apache.dubbo.userprovider2/0.0.1/getuser")
//	assert.Nil(t, err)
//	api, err = localAPIDiscSrv.GetAPI("/UserInfoServer2/org.apache.dubbo.UserProvider2/0.0.1/SetUser", config.MethodPost)
//	assert.Equal(t, api.URLPattern, "/userinfoserver2/org.apache.dubbo.userprovider2/0.0.1/setuser")
//	assert.Nil(t, err)
//	c.Destroy()
//}
