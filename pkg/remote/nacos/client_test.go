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

package nacos

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestNacosV2Integration(t *testing.T) {
	address := os.Getenv("NACOS_TEST_ADDRESS")
	if address == "" {
		t.Skip("set NACOS_TEST_ADDRESS to a Nacos 2.x main port; main port + 1000 must also be reachable")
	}
	client, err := NewNacosClient(&model.RemoteConfig{Address: address, Timeout: "5s"})
	require.NoError(t, err)
	t.Cleanup(client.Close)
	service := fmt.Sprintf("pixiu-v2-smoke-%d", time.Now().UnixNano())
	registered, err := client.namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip: "127.0.0.1", Port: 18080, Weight: 1, Enable: true, Healthy: true,
		ServiceName: service, GroupName: "DEFAULT_GROUP", ClusterName: "DEFAULT", Ephemeral: true,
	})
	require.NoError(t, err)
	require.True(t, registered)
	t.Cleanup(func() {
		removed, err := client.namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip: "127.0.0.1", Port: 18080, ServiceName: service,
			GroupName: "DEFAULT_GROUP", Cluster: "DEFAULT", Ephemeral: true,
		})
		assert.NoError(t, err)
		assert.True(t, removed)
	})
	require.Eventually(t, func() bool {
		instances, err := client.SelectInstances(vo.SelectInstancesParam{
			ServiceName: service, GroupName: "DEFAULT_GROUP", HealthyOnly: true,
		})
		return err == nil && len(instances) == 1 && instances[0].Ip == "127.0.0.1" && instances[0].Port == 18080
	}, 20*time.Second, 200*time.Millisecond, "registered instance must be discoverable over the v2 gRPC connection")
}

// closeRecorder embeds the SDK naming interface and records CloseClient calls
// so NacosClient.Close can be asserted without a real gRPC connection.
type closeRecorder struct {
	naming_client.INamingClient

	mu         sync.Mutex
	closeCount int
}

func (c *closeRecorder) CloseClient() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeCount++
}

func (c *closeRecorder) closeCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closeCount
}

// TestNacosClient_Close delegates to the underlying v2 naming client's
// CloseClient so the gRPC connection is released on shutdown.
func TestNacosClient_Close(t *testing.T) {
	recorder := &closeRecorder{}
	client := &NacosClient{namingClient: recorder}

	assert.Equal(t, 0, recorder.closeCalls(), "underlying client must not be closed before Close")

	client.Close()

	assert.Equal(t, 1, recorder.closeCalls(),
		"Close must delegate to the underlying naming client's CloseClient")
}
