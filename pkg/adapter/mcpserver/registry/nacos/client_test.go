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
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/stretchr/testify/assert"
)

// Test constants
const (
	// Network and Service constants
	testIP                   = "127.0.0.1"
	testPort                 = 8080
	testNamespace            = "public"
	testGroupName            = "DEFAULT_GROUP"
	testGroupNameMcpVersions = "mc-server-versions"

	// Test data constants
	testServerCount      = 151
	testFilteredCount    = 150
	testRetryMaxAttempts = 50
	testRetryInterval    = 10 * time.Millisecond

	// Test IDs and names
	testMcpServerID    = "a4768d16-8263-48ea-8994-e003a2c80271"
	testLocalServerID  = "52df06fe-5433-4154-b8e2-3fbb33ca5a33"
	testServerName     = "explore"
	testServiceName    = "explore"
	testServiceNameNew = "explore-new"

	// Version constants
	testVersion112    = "1.0.12"
	testVersion113    = "1.0.13"
	testVersionLatest = "1.0.2"

	// Configuration keys
	testConfigKey      = "test"
	testConfigKey1     = "test1"
	testConfigKey3     = "test3"
	testCredentialKey  = "test_test"
	testCredentialKey1 = "test1_test1"
	testCredentialKey3 = "test3_test3"

	// Configuration values
	testSecretKey  = "secret_key"
	testSecretKey1 = "secret_key_1"
	testSecretKey3 = "secret_key_3"

	// Test data separator
	configKeySeparator = "$$"
)

// Test data generation functions

// createMcpServerVersionConfig creates a basic MCP server version configuration JSON
func createMcpServerVersionConfig(id, name, protocol, frontProtocol, version string) string {
	return fmt.Sprintf(`{"id":"%s","name":"%s","protocol":"%s","frontProtocol":"%s","description":"%s","enabled":true,"capabilities":["TOOL"],"latestPublishedVersion":"%s","versionDetails":[{"version":"1.0.0","release_date":"2025-06-09T05:41:16Z","is_latest":false},{"version":"1.0.1","release_date":"2025-06-09T05:41:37Z","is_latest":false},{"version":"%s","release_date":"2025-06-09T05:42:46Z","is_latest":true}]}`, id, name, protocol, frontProtocol, name, version, version)
}

// createExploreServerVersionConfig creates the explore server version configuration
func createExploreServerVersionConfig(latestVersion string) string {
	baseVersions := `[{"version":"1.0.0","release_date":"2025-06-05T10:11:40Z","is_latest":false},{"version":"1.0.1","release_date":"2025-06-05T10:12:59Z","is_latest":false},{"version":"1.0.2","release_date":"2025-06-05T10:21:28Z","is_latest":false},{"version":"1.0.3","release_date":"2025-06-05T10:21:39Z","is_latest":false},{"version":"1.0.4","release_date":"2025-06-05T10:25:04Z","is_latest":false},{"version":"1.0.6","release_date":"2025-06-05T10:25:24Z","is_latest":false},{"version":"1.0.8","release_date":"2025-06-05T10:27:38Z","is_latest":false},{"version":"1.0.9","release_date":"2025-06-05T10:32:13Z","is_latest":false},{"version":"1.0.10","release_date":"2025-06-05T10:32:28Z","is_latest":false},{"version":"1.0.11","release_date":"2025-06-05T11:04:09Z","is_latest":true},{"version":"1.0.12"}]`

	return fmt.Sprintf(`{"id":"%s","name":"%s","protocol":"https","frontProtocol":"mcp-streamable","description":"%s","enabled":true,"capabilities":["TOOL"],"latestPublishedVersion":"%s","versionDetails":%s}`, testMcpServerID, testServerName, testServerName, latestVersion, baseVersions)
}

// createMcpServerConfig creates MCP server configuration JSON
func createMcpServerConfig(id, version, serviceName string) string {
	return fmt.Sprintf(`{"id":"%s","name":"%s","protocol":"https","frontProtocol":"mcp-streamable","description":"%s","versionDetail":{"version":"%s"},"remoteServerConfig":{"serviceRef":{"namespaceId":"%s","groupName":"%s","serviceName":"%s"},"exportPath":""},"enabled":true,"capabilities":["TOOL"],"toolsDescriptionRef":"%s-%s-mcp-tools.json"}`, id, testServerName, testServerName, version, testNamespace, testGroupName, serviceName, id, version)
}

// createMcpToolsConfig creates MCP tools configuration JSON
func createMcpToolsConfig(configRef string) string {
	return fmt.Sprintf(`{"tools":[{"name":"%s","description":"%s","inputSchema":{"type":"object","properties":{"tags":{"description":"tags","type":"string"}}}}],"toolsMeta":{"%s":{"enabled":true,"templates":{"json-go-template":{"requestTemplate":{"method":"GET","url":"/v0/explore?key={{ ${nacos.%s}.key }}","argsToUrlParam":true}}}}}}`, testServerName, testServerName, testServerName, configRef)
}

// createCredentialConfig creates credential configuration JSON
func createCredentialConfig(keyValue string) string {
	return fmt.Sprintf(`{
    "key": "%s"
}`, keyValue)
}

// createBrokenJSON creates intentionally broken JSON for testing error handling
func createBrokenJSON() string {
	return "{"
}

type MockedNacosConfigClient struct {
	configs           map[string]any
	configListenerMap map[string][]func(string, string, string, string)
}

func (m MockedNacosConfigClient) GetConfig(param vo.ConfigParam) (string, error) {
	if result, exist := m.configs[param.DataId+"$$"+param.Group]; exist {
		config, ok := result.(string)
		if ok {
			return config, nil
		}

		err, ok := result.(error)
		if ok {
			return "", err
		}

		return "", fmt.Errorf("unknown config type")
	}
	return "", nil
}

func (m MockedNacosConfigClient) PublishConfig(_ vo.ConfigParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosConfigClient) DeleteConfig(_ vo.ConfigParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosConfigClient) ListenConfig(params vo.ConfigParam) (err error) {
	if _, ok := m.configListenerMap[params.Group]; !ok {
		m.configListenerMap[params.Group] = []func(string, string, string, string){}
	}
	m.configListenerMap[params.DataId+"$$"+params.Group] = append(m.configListenerMap[params.DataId+"$$"+params.Group], params.OnChange)
	return nil
}

func (m MockedNacosConfigClient) CancelListenConfig(params vo.ConfigParam) (err error) {
	delete(m.configListenerMap, params.DataId+"$$"+params.Group)
	return nil
}

func (m MockedNacosConfigClient) SearchConfig(param vo.SearchConfigParam) (*model.ConfigPage, error) {
	dataIdRegex := strings.ReplaceAll(param.DataId, "*", ".*")
	groupRegex := strings.ReplaceAll(param.Group, "*", ".*")
	result := []model.ConfigItem{}

	for key, value := range m.configs {
		dataIdAndGroup := strings.Split(key, "$$")
		dataId := dataIdAndGroup[0]
		group := dataIdAndGroup[1]
		if regexp.MustCompile(dataIdRegex).MatchString(dataId) && regexp.MustCompile(groupRegex).MatchString(group) {
			result = append(result, model.ConfigItem{
				DataId:  dataId,
				Group:   group,
				Content: value.(string),
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DataId < result[j].DataId
	})

	offset := param.PageSize * (param.PageNo - 1)
	size := param.PageSize
	if offset+param.PageSize > len(result) {
		size = len(result) - offset
	}
	finalResult := result[offset : offset+size]
	return &model.ConfigPage{
		TotalCount:     len(result),
		PageNumber:     param.PageNo,
		PagesAvailable: len(result)/param.PageSize + 1,
		PageItems:      finalResult,
	}, nil
}

func (m MockedNacosConfigClient) CloseClient() {
	//TODO implement me
	panic("implement me")
}

type MockedNacosNamingClient struct {
	listenerMap map[string][]func(services []model.Instance, err error)
}

func (m MockedNacosNamingClient) RegisterInstance(_ vo.RegisterInstanceParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) BatchRegisterInstance(_ vo.BatchRegisterInstanceParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) DeregisterInstance(_ vo.DeregisterInstanceParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) UpdateInstance(_ vo.UpdateInstanceParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) GetService(param vo.GetServiceParam) (model.Service, error) {
	return model.Service{
		Name:      param.ServiceName,
		GroupName: param.GroupName,
		Hosts: []model.Instance{
			{
				Ip:   testIP,
				Port: testPort,
			},
		},
	}, nil
}

func (m MockedNacosNamingClient) SelectAllInstances(_ vo.SelectAllInstancesParam) ([]model.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) SelectInstances(_ vo.SelectInstancesParam) ([]model.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) SelectOneHealthyInstance(_ vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) Subscribe(param *vo.SubscribeParam) error {
	if m.listenerMap[param.ServiceName+"$$"+param.GroupName] == nil {
		m.listenerMap[param.ServiceName+"$$"+param.GroupName] = []func([]model.Instance, error){}
	}
	m.listenerMap[param.ServiceName+"$$"+param.GroupName] = append(m.listenerMap[param.ServiceName+"$$"+param.GroupName], param.SubscribeCallback)
	return nil
}

func (m MockedNacosNamingClient) Unsubscribe(_ *vo.SubscribeParam) error {
	return nil
}

func (m MockedNacosNamingClient) GetAllServicesInfo(_ vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) ServerHealthy() bool {
	//TODO implement me
	panic("implement me")
}

func (m MockedNacosNamingClient) CloseClient() {
	//TODO implement me
	panic("implement me")
}

func TestNacosRegistryClient_ListMcpServer(t *testing.T) {
	// Test case 1: List multiple pages with valid MCP servers
	mockedConfigs := map[string]any{}
	for i := 0; i < testServerCount; i++ {
		configKey := fmt.Sprintf("%d-mcp-versions.json%smcp-server-versions", i, configKeySeparator)
		configValue := createMcpServerVersionConfig(fmt.Sprintf("%d", i), "test", "http", "mcp-streamable", testVersionLatest)
		mockedConfigs[configKey] = configValue
	}

	client := NacosRegistryClient{
		configClient: MockedNacosConfigClient{configs: mockedConfigs},
	}

	servers, err := client.ListMcpServer()
	if err != nil {
		t.Fatalf("Failed to list MCP servers: %v", err)
	}
	assert.Equal(t, testServerCount, len(servers))

	// Verify no duplicate server IDs
	serverMap := map[string]string{}
	for _, info := range servers {
		if _, exists := serverMap[info.Id]; exists {
			t.Fatalf("Duplicate server ID found: %s", info.Id)
		}
		serverMap[info.Id] = info.Id
	}

	// Test case 2: Local server should not be listed (frontProtocol is "stdio")
	localServerKey := fmt.Sprintf("65-mcp-versions.json%smcp-server-versions", configKeySeparator)
	localServerConfig := createMcpServerVersionConfig(testLocalServerID, "test", "http", "stdio", testVersionLatest)
	mockedConfigs[localServerKey] = localServerConfig

	servers, err = client.ListMcpServer()
	if err != nil {
		t.Fatalf("Failed to list MCP servers after adding local server: %v", err)
	}
	assert.Equal(t, testFilteredCount, len(servers))

	// Test case 3: Broken config should not cause failure, just be skipped
	mockedConfigs[localServerKey] = createBrokenJSON()
	servers, err = client.ListMcpServer()
	if err != nil {
		t.Fatalf("Failed to list MCP servers with broken config: %v", err)
	}
	assert.Equal(t, testFilteredCount, len(servers))
}

func TestNacosRegistryClient_ListenToMcpServer(t *testing.T) {
	// Create test configuration data using helper functions
	versionConfigKey := fmt.Sprintf("%s-mcp-versions.json%smcp-server-versions", testMcpServerID, configKeySeparator)
	serverConfigKey112 := fmt.Sprintf("%s-%s-mcp-server.json%smcp-server", testMcpServerID, testVersion112, configKeySeparator)
	toolsConfigKey112 := fmt.Sprintf("%s-%s-mcp-tools.json%smcp-tools", testMcpServerID, testVersion112, configKeySeparator)
	serverConfigKey113 := fmt.Sprintf("%s-%s-mcp-server.json%smcp-server", testMcpServerID, testVersion113, configKeySeparator)
	toolsConfigKey113 := fmt.Sprintf("%s-%s-mcp-tools.json%smcp-tools", testMcpServerID, testVersion113, configKeySeparator)

	configClient := MockedNacosConfigClient{
		configs: map[string]any{
			versionConfigKey:   createExploreServerVersionConfig(testVersion112),
			serverConfigKey112: createMcpServerConfig(testMcpServerID, testVersion112, testServiceName),
			toolsConfigKey112:  createMcpToolsConfig(fmt.Sprintf("%s/%s", testConfigKey, testConfigKey)),
			serverConfigKey113: createMcpServerConfig(testMcpServerID, testVersion113, testServiceName),
			toolsConfigKey113:  createMcpToolsConfig(fmt.Sprintf("%s/%s", testConfigKey3, testConfigKey3)),
			fmt.Sprintf("%s%s%s", testConfigKey, configKeySeparator, testConfigKey):   createCredentialConfig(testSecretKey),
			fmt.Sprintf("%s%s%s", testConfigKey1, configKeySeparator, testConfigKey1): createCredentialConfig(testSecretKey1),
			fmt.Sprintf("%s%s%s", testConfigKey3, configKeySeparator, testConfigKey3): createCredentialConfig(testSecretKey3),
		},
		configListenerMap: map[string][]func(string, string, string, string){},
	}

	namingClient := MockedNacosNamingClient{
		listenerMap: map[string][]func(services []model.Instance, err error){},
	}
	client := NacosRegistryClient{
		configClient: configClient,
		namingClient: namingClient,
		servers:      map[string]*ServerContext{},
	}

	// Verify initial server listing
	servers, err := client.ListMcpServer()
	if err != nil {
		t.Fatalf("Failed to list MCP servers: %v", err)
	}
	assert.Equal(t, 1, len(servers))

	// Set up listener for configuration changes
	var newConfig *McpServerConfig
	err = client.ListenToMcpServer(testMcpServerID, func(info *McpServerConfig) {
		newConfig = info
	})
	if err != nil {
		t.Fatalf("Failed to start listening to MCP server: %v", err)
	}

	// Wait for initial configuration to be loaded
	for i := 0; i < testRetryMaxAttempts && newConfig == nil; i++ {
		time.Sleep(testRetryInterval)
	}

	// Verify initial configuration
	expectedServerConfig := createMcpServerConfig(testMcpServerID, testVersion112, testServiceName)
	expectedToolsConfig := createMcpToolsConfig(fmt.Sprintf("%s/%s", testConfigKey, testConfigKey))
	// Replace nacos template with processed version
	expectedToolsConfig = strings.ReplaceAll(expectedToolsConfig, fmt.Sprintf("${nacos.%s/%s}", testConfigKey, testConfigKey), fmt.Sprintf(".config.credentials.%s", testCredentialKey))

	assert.Equal(t, expectedServerConfig, newConfig.ServerSpecConfig)
	assert.Equal(t, expectedToolsConfig, newConfig.ToolsSpecConfig)
	assert.Equal(t, 1, len(newConfig.Credentials))
	assert.Equal(t, map[string]any{"key": testSecretKey}, newConfig.Credentials[testCredentialKey])

	// Test case 1: Change tool nacos template reference
	listener := configClient.configListenerMap[toolsConfigKey112][0]
	updatedToolsConfig := createMcpToolsConfig(fmt.Sprintf("%s/%s", testConfigKey1, testConfigKey1))
	listener(testNamespace, "mcp-tools", toolsConfigKey112, updatedToolsConfig)

	// Wait for tools update to propagate
	for i := 0; i < testRetryMaxAttempts; i++ {
		if newConfig != nil && strings.Contains(newConfig.ToolsSpecConfig, testCredentialKey1) {
			break
		}
		time.Sleep(testRetryInterval)
	}

	// Verify updated tools configuration
	expectedUpdatedToolsConfig := strings.ReplaceAll(updatedToolsConfig, fmt.Sprintf("${nacos.%s/%s}", testConfigKey1, testConfigKey1), fmt.Sprintf(".config.credentials.%s", testCredentialKey1))
	assert.Equal(t, expectedUpdatedToolsConfig, newConfig.ToolsSpecConfig)
	assert.Equal(t, 1, len(newConfig.Credentials))
	assert.Equal(t, map[string]any{"key": testSecretKey1}, newConfig.Credentials[testCredentialKey1])

	// Test case 2: Change backend service name
	serviceListener := configClient.configListenerMap[serverConfigKey112][0]
	updatedServerConfig := createMcpServerConfig(testMcpServerID, testVersion112, testServiceNameNew)
	serviceListener(testNamespace, "mcp-server", serverConfigKey112, updatedServerConfig)

	for i := 0; i < testRetryMaxAttempts; i++ {
		if newConfig != nil && strings.Contains(newConfig.ServerSpecConfig, testServiceNameNew) {
			break
		}
		time.Sleep(testRetryInterval)
	}

	// Test case 3: Publish new version of MCP server
	versionListener := configClient.configListenerMap[versionConfigKey][0]
	updatedVersionConfig := createExploreServerVersionConfig(testVersion113)
	versionListener(testNamespace, testGroupNameMcpVersions, versionConfigKey, updatedVersionConfig)

	// Wait for version update to trigger server config change
	for i := 0; i < testRetryMaxAttempts; i++ {
		if newConfig != nil && strings.Contains(newConfig.ServerSpecConfig, fmt.Sprintf("\"version\":\"%s\"", testVersion113)) {
			break
		}
		time.Sleep(testRetryInterval)
	}

	// Wait for tools config to update to new version reference
	for i := 0; i < testRetryMaxAttempts; i++ {
		if newConfig != nil && strings.Contains(newConfig.ToolsSpecConfig, testCredentialKey3) {
			break
		}
		time.Sleep(testRetryInterval)
	}

	// Verify final configuration state
	expectedFinalServerConfig := createMcpServerConfig(testMcpServerID, testVersion113, testServiceName)
	expectedFinalToolsConfig := createMcpToolsConfig(fmt.Sprintf("%s/%s", testConfigKey3, testConfigKey3))
	expectedFinalToolsConfig = strings.ReplaceAll(expectedFinalToolsConfig, fmt.Sprintf("${nacos.%s/%s}", testConfigKey3, testConfigKey3), fmt.Sprintf(".config.credentials.%s", testCredentialKey3))

	assert.Equal(t, expectedFinalServerConfig, newConfig.ServerSpecConfig)
	assert.Equal(t, expectedFinalToolsConfig, newConfig.ToolsSpecConfig)
	assert.Equal(t, 1, len(newConfig.Credentials))
	assert.Equal(t, map[string]any{"key": testSecretKey3}, newConfig.Credentials[testCredentialKey3])
}
