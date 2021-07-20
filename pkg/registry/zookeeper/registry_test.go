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
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go/common"
	ex "github.com/apache/dubbo-go/common/extension"
	"github.com/apache/dubbo-go/metadata/service"
	"github.com/apache/dubbo-go/protocol"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	"github.com/apache/dubbo-go/registry"
	"github.com/dubbogo/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/service/api"
)

func TestNewZKRegistry(t *testing.T) {
	regConfig := model.Registry{
		Protocol: "zookeeper",
		Timeout:  "2s",
		Address:  "127.0.0.1:9100",
	}
	reg, err := newZKRegistry(regConfig)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	regConfig = model.Registry{
		Protocol: "zookeeper",
		Timeout:  "2xxxxxx",
		Address:  "127.0.0.1:9100",
	}
	reg, err = newZKRegistry(regConfig)
	assert.Nil(t, reg)
	assert.NotNil(t, err)

	regConfig = model.Registry{
		Protocol: "zookeeper",
		Timeout:  "2s",
		Address:  "",
	}
	reg, err = newZKRegistry(regConfig)
	assert.Nil(t, reg)
	assert.NotNil(t, err)
}

func TestLoadInterfaces(t *testing.T) {
	api.Init()
	tc, err := zk.StartTestCluster(1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Stop()

	conn, _, _ := tc.ConnectAll()
	conn.Create("/dubbo", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/com.ikurento.user.UserProvider1", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/com.ikurento.user.UserProvider1/providers", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/com.ikurento.user.UserProvider1/providers/"+
		url.QueryEscape("dubbo://192.168.3.46:20002/org.apache.dubbo.UserProvider1?anyhost=true&app.version=0.0.1&application=UserInfoServer1&bean.name=UserProvider&cluster=failover&interface=org.apache.dubbo.UserProvider1&loadbalance=random&methods=GetUser&name=UserInfoServer&organization=dubbo.io&registry.role=3&service.filter=echo,token,accesslog,tps,generic_service,execute,pshutdown&side=provider&ssl-enabled=false&timestamp=1624716984&warmup=100"),
		nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/com.ikurento.user.UserProvider2", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/com.ikurento.user.UserProvider2/providers", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/com.ikurento.user.UserProvider2/providers/"+
		url.QueryEscape("dubbo://192.168.3.46:20001/org.apache.dubbo.UserProvider2?anyhost=true&app.version=0.0.1&application=UserInfoServer2&bean.name=UserProvider&cluster=failover&interface=org.apache.dubbo.UserProvider2&loadbalance=random&methods=GetUser,SetUser,AllUsers&name=UserInfoServer&organization=dubbo.io&registry.role=3&service.filter=echo,token,accesslog,tps,generic_service,execute,pshutdown&side=provider&ssl-enabled=false&timestamp=1624716984&warmup=100"),
		nil, 0, zk.WorldACL(zk.PermAll))

	hosts := make([]string, len(tc.Servers))
	for i, srv := range tc.Servers {
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	}
	regConfig := model.Registry{
		Protocol: "zookeeper",
		Timeout:  "20s",
		Address:  hosts[0],
	}
	reg, _ := newZKRegistry(regConfig)
	r := reg.(*ZKRegistry)
	c := r.GetClient()
	connState := c.GetConnState()
	for connState != zk.StateConnected && connState != zk.StateHasSession {
		<-time.After(200 * time.Millisecond)
		connState = c.GetConnState()
	}

	apis, errs := r.LoadInterfaces()
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, apis[0].URLPattern, "/UserInfoServer2/org.apache.dubbo.UserProvider2/0.0.1/GetUser")
	assert.Nil(t, err)
	assert.Equal(t, apis[1].URLPattern, "/UserInfoServer2/org.apache.dubbo.UserProvider2/0.0.1/SetUser")
	assert.Nil(t, err)
	assert.Equal(t, apis[2].URLPattern, "/UserInfoServer2/org.apache.dubbo.UserProvider2/0.0.1/AllUsers")
	assert.Nil(t, err)
	assert.Equal(t, apis[3].URLPattern, "/UserInfoServer1/org.apache.dubbo.UserProvider1/0.0.1/GetUser")
	assert.Nil(t, err)
	c.Destroy()
}

func TestLoadApplications(t *testing.T) {
	api.Init()
	tc, err := zk.StartTestCluster(1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Stop()

	conn, _, _ := tc.ConnectAll()
	conn.Create("/services", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/services/UserInfoServer", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/services/UserInfoServer/127.0.0.1:20001", []byte("{\"Name\":\"UserInfoServer\",\"Id\":\"127.0.0.1:20001\",\"Address\":\"127.0.0.1\",\"Port\":20001,\"Payload\":{\"id\":\"172.23.240.1:20001\",\"metadata\":{\"dubbo.endpoints\":\"[{\\\"port\\\":20001,\\\"protocol\\\":\\\"dubbo\\\"}]\",\"dubbo.exported-services.revision\":\"480513435\",\"dubbo.metadata-service.url-params\":\"{\\\"dubbo\\\":{\\\"app.version\\\":\\\"0.0.1\\\",\\\"bean.name\\\":\\\"MetadataService\\\",\\\"environment\\\":\\\"dev\\\",\\\"export\\\":\\\"true\\\",\\\"interface\\\":\\\"org.apache.dubbo.metadata.MetadataService\\\",\\\"message_size\\\":\\\"0\\\",\\\"module\\\":\\\"dubbo-go user-info server\\\",\\\"name\\\":\\\"UserInfoServer\\\",\\\"organization\\\":\\\"dubbo.io\\\",\\\"port\\\":\\\"20001\\\",\\\"registry.role\\\":\\\"3\\\",\\\"release\\\":\\\"dubbo-golang-1.5.7\\\",\\\"service.filter\\\":\\\"echo,token,accesslog,tps,generic_service,execute,pshutdown\\\",\\\"side\\\":\\\"provider\\\",\\\"ssl-enabled\\\":\\\"false\\\",\\\"version\\\":\\\"1.0.0\\\"}}\",\"dubbo.metadata.storage-type\":\"local\",\"dubbo.subscribed-services.revision\":\"N/A\"},\"name\":\"UserInfoServer\"},\"RegistrationTimeUTC\":0}"), 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo", nil, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/metadata", nil, 0, zk.WorldACL(zk.PermAll))
	methodsData := []byte("{\"CanonicalName\":\"org.apache.dubbo.UserProvider\",\"CodeSource\":\"\",\"Methods\":[{\"Name\":\"GetUser\",\"ParameterTypes\":[\"slice\"],\"ReturnType\":\"ptr\",\"Parameters\":null}, {\"Name\":\"SetUser\",\"ParameterTypes\":[\"slice\"],\"ReturnType\":\"ptr\",\"Parameters\":null}],\"Types\":null}")
	conn.Create("/dubbo/metadata/org.apache.dubbo.UserProvider", methodsData, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/metadata/org.apache.dubbo.UserProvider/dubbo", methodsData, 0, zk.WorldACL(zk.PermAll))
	conn.Create("/dubbo/metadata/org.apache.dubbo.UserProvider/dubbo/provider", methodsData, 0, zk.WorldACL(zk.PermAll))

	createPxy()
	hosts := make([]string, len(tc.Servers))
	for i, srv := range tc.Servers {
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	}
	regConfig := model.Registry{
		Protocol: "zookeeper",
		Timeout:  "20s",
		Address:  hosts[0],
	}
	reg, _ := newZKRegistry(regConfig)
	r := reg.(*ZKRegistry)
	c := r.GetClient()
	connState := c.GetConnState()
	for connState != zk.StateConnected && connState != zk.StateHasSession {
		<-time.After(200 * time.Millisecond)
		connState = c.GetConnState()
	}

	apis, errs := r.LoadApplications()
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, apis[0].URLPattern, "/UserInfoServer/org.apache.dubbo.UserProvider/0.0.1/GetUser")
	assert.Nil(t, err)
	assert.Equal(t, apis[1].URLPattern, "/UserInfoServer/org.apache.dubbo.UserProvider/0.0.1/SetUser")
	assert.Nil(t, err)

	c.Destroy()
}

func createPxy() service.MetadataService {
	ex.SetProtocol("dubbo", func() protocol.Protocol {
		return &mockProtocol{}
	})

	ins := &registry.DefaultServiceInstance{
		Id:          "127.0.0.1:20001",
		ServiceName: "UserInfoServer",
		Host:        "127.0.0.1",
		Port:        20001,
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"dubbo.endpoints":                    "[{\"port\":20001,\"protocol\":\"dubbo\"}]",
			"dubbo.exported-services.revision":   "480513435",
			"dubbo.metadata-service.url-params":  "{\"dubbo\":{\"app.version\":\"0.0.1\",\"bean.name\":\"MetadataService\",\"environment\":\"dev\",\"export\":\"true\",\"interface\":\"org.apache.dubbo.metadata.MetadataService\",\"message_size\":\"0\",\"module\":\"dubbo-go user-info server\",\"name\":\"UserInfoServer\",\"organization\":\"dubbo.io\",\"port\":\"20001\",\"registry.role\":\"3\",\"release\":\"dubbo-golang-1.5.7\",\"service.filter\":\"echo,token,accesslog,tps,generic_service,execute,pshutdown\",\"side\":\"provider\",\"ssl-enabled\":\"false\",\"version\":\"1.0.0\"}}",
			"dubbo.metadata.storage-type":        "local",
			"dubbo.subscribed-services.revision": "N/A",
		},
	}
	return ex.GetMetadataServiceProxyFactory(constant.DefaultMetadataStorageType).GetProxy(ins)
}

type mockProtocol struct {
}

func (m mockProtocol) Export(protocol.Invoker) protocol.Exporter {
	panic("implement me")
}

func (m mockProtocol) Refer(*common.URL) protocol.Invoker {
	return &mockInvoker{}
}

func (m mockProtocol) Destroy() {
	panic("implement me")
}

type mockInvoker struct {
}

func (m *mockInvoker) GetURL() *common.URL {
	panic("implement me")
}

func (m *mockInvoker) IsAvailable() bool {
	panic("implement me")
}

func (m *mockInvoker) Destroy() {
	panic("implement me")
}

func (m *mockInvoker) Invoke(context.Context, protocol.Invocation) protocol.Result {
	return &protocol.RPCResult{
		Rest: &[]interface{}{"dubbo://:20001/org.apache.dubbo.UserProvider?accesslog=&app.version=0.0.1&application=UserInfoServer&auth=&bean.name=UserProvider&cluster=failover&environment=dev&execute.limit=&execute.limit.rejected.handler=&export=true&group=&interface=org.apache.dubbo.UserProvider&loadbalance=random&message_size=4&methods.GetUser.loadbalance=random&methods.GetUser.retries=1&methods.GetUser.tps.limit.interval=&methods.GetUser.tps.limit.rate=&methods.GetUser.tps.limit.strategy=&methods.GetUser.weight=0&module=dubbo-go+user-info+server&name=UserInfoServer&organization=dubbo.io&owner=&param.sign=&registry.role=3&release=dubbo-golang-1.5.7&retries=&serialization=&service.filter=echo%2Ctoken%2Caccesslog%2Ctps%2Cgeneric_service%2Cexecute%2Cpshutdown&side=provider&ssl-enabled=false&timestamp=1626573430&tps.limit.interval=&tps.limit.rate=&tps.limit.rejected.handler=&tps.limit.strategy=&tps.limiter=&version=&warmup=100"},
	}
}
