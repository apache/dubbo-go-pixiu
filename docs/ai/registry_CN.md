## LLM 服务发现与注册

[English](registry.md) | 中文

本文档旨在指导 LLM 服务提供商如何通过 Nacos 注册中心，将其服务实例动态地注册到 LLM 网关。通过遵循这些准则，网关将能够自动发现您的服务，并根据您提供的元数据（metadata）应用相应的路由、重试和 fallback 策略。

### 注册机制概述

服务发现的核心机制是：您的 LLM 服务作为一个**Nacos 实例**进行注册，并在注册时提供一组特定的**元数据**。LLM 网关会监听 Nacos 中的服务变更，读取这些元数据，并将其动态地转换为一个功能齐全的网关 `endpoint` 配置。

一个基本的 Nacos 注册请求包含以下关键信息：

- **`ServiceName`**: 您的服务集合的名称 (例如, `deepseek-service`)。
- **`Ip` & `Port`**: 您的服务实例监听流量的网络地址。
- **`Metadata`**: 一个键值对集合，**这是配置所有网关行为的关键**。

### `metadata` 配置字段

所有网关的特定配置都通过 Nacos 实例的 `metadata` 字段进行传递。以下是所有支持的 `metadata` 键及其说明。

`cluster`

- **类型**: `string`
- **必需**: 是
- **描述**: 定义此 endpoint 应归属到网关中的哪个集群 (`cluster`)。网关会根据此值将具有相同 `cluster` 名称的服务实例聚合在一起。
- **示例**: `"deepseek_cluster"`

`id`

- **类型**: `string`
- **必需**: 是
- **描述**: endpoint 的唯一标识符。它在所属的 `cluster` 内必须是唯一的。
- **示例**: `"deepseek-primary-instance-1"`

`ip`

- **类型**: `string`
- **必需**: 否
- **描述**: 可选参数，当您的服务需要向网关注册一个不同于其内部 IP 的地址时（例如，一个可公开访问的 IP），此字段非常有用。如果未提供，网关将使用默认值 `0.0.0.0`。
- **示例**: `"203.0.113.55"`

`port`

- **类型**: `string` (代表一个整数)
- **必需**: 否
- **描述**: 可选参数，用于覆盖 Nacos 实例本身注册的 `Port`。用法与 `ip` 字段类似。
- **示例**: `"9000"`

`name`

- **类型**: `string`
- **必需**: 否
- **描述**: 一个人类可读的名称，用于标识此 endpoint。主要用于日志和监控。
- **示例**: `"DeepSeek V2 Chat (Primary)"`

`address`

- **类型**: `string`
- **必需**: 否
- **描述**: 以逗号分隔的字符串，每一个字符串代表了一个 address
- **示例**: `"api.deepseek.com"`

`llm-meta.fallback`

- **类型**: `string` ("true" 或 "false")
- **必需**: 否, 默认为 `"false"`
- **描述**: 决定如果在此 endpoint 上的所有重试尝试都失败后，网关是否应继续处理集群中的下一个 endpoint。

`llm-meta.api_key`

- **类型**: `string`
- **必需**: 否
- **描述**: 此 endpoint 使用的 API 密钥。

`llm-meta.retry_policy.name`

- **类型**: `string`
- **必需**: 否, 默认为 `"NoRetry"`
- **描述**: 指定当对此 endpoint 的请求失败时要使用的重试策略的名称。名称不区分大小写。

`llm-meta.retry_policy.config`

- **类型**: `string` (JSON 格式)
- **必需**: 否
- **描述**: 一个 JSON 字符串，包含特定于所选重试策略的配置参数。

### 重试策略元数据

您可以通过组合使用 `llm-meta.retry_policy.name` 和 `llm-meta.retry_policy.config` 来配置重试行为。

1.  `CountBased` (基于次数)

    - **name**: `"CountBased"`
    - **config**: 一个包含 `times` 字段的 JSON 字符串。
        - `times` (integer): 重试的次数。
    - **示例**:
        - `llm-meta.retry_policy.name`: `"CountBased"`
        - `llm-meta.retry_policy.config`: `{"times": 2}`

2.  `ExponentialBackoff` (指数退避)

    - **name**: `"ExponentialBackoff"`
    - **config**: 一个包含 `times`, `initialInterval`, `maxInterval`, 和 `multiplier` 字段的 JSON 字符串。
        - `times` (integer): 重试次数。
        - `initialInterval` (string): 初始等待时长 (例如 "200ms")。
        - `maxInterval` (string): 最大等待时长 (例如 "5s")。
        - `multiplier` (float): 延迟时间的乘数因子。
    - **示例**:
        - `llm-meta.retry_policy.name`: `"ExponentialBackoff"`
        - `llm-meta.retry_policy.config`: `{"times": 3, "initialInterval": "200ms", "maxInterval": "5s", "multiplier": 2.0}`

3.  `NoRetry` (不重试)

    - **name**: `"NoRetry"`
    - **config**: 此策略不需要 `config` 字段。

### 完整注册示例 (Go)

更详细的使用方法可以参考[官方示例](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/llm/nacos)。

pixiu 配置文件，需要启用llmregistrycenter这个适配器，示例如下：

```yaml
  adapters:
    - id: test
      name: dgp.adapter.llmregistrycenter
      config:
        registries:
          nacos:
            protocol: nacos
            address: "127.0.0.1:8848"
            timeout: "5s"
            group: test_llm_registry_group
            namespace: public
```

此示例展示了如何使用 Nacos Go SDK 注册一个功能完整的 LLM 服务实例。该实例将被配置为 `deepseek_cluster` 中的主 endpoint，使用指数退避重试策略，并在失败时 fallback 到集群中的下一个服务。

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// ... (此处省略了创建 Nacos 客户端的代码)
	// client, err := createNacosClient()

	// 1. 准备重试策略的 JSON 配置
	retryConfig := map[string]interface{}{
		"times":           3,
		"initialInterval": "200ms",
		"maxInterval":     "8s",
		"multiplier":      2.5,
	}
	retryConfigJSON, _ := json.Marshal(retryConfig)

	// 2. 构造包含所有网关配置的 metadata
	metadata := map[string]string{
		// --- 核心 Endpoint 配置 ---
		"cluster": "deepseek_cluster",
		"id":      "deepseek-primary",
		"name":    "DeepSeek V2 Chat (Primary)",

		// 可选: (使用ip+port或者address): 实例的 IP 和 Port
		"ip":   "203.0.113.55",
		"port": "9000",

		// 可选: (使用ip+port或者address): address 列表
		"address": "api.deepseek.com",

		// --- LLM 特定元数据 ---
		"llm-meta.fallback": "true",

		// 使用 JSON 字符串格式的 API Keys
		"llm-meta.api_key": "key-xxxxxxxx",

		// --- 重试策略配置 ---
		"llm-meta.retry_policy.name":   "ExponentialBackoff",
		"llm-meta.retry_policy.config": string(retryConfigJSON),
	}

	// 3. 注册 Nacos 实例
	// 注意：这里的 Ip 和 Port 是服务实例的实际监听地址，
	// 而 metadata 中的 ip 和 port 是希望网关访问的地址。
	_, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.1.10", // 服务的内部 IP
		Port:        8001,           // 服务的内部端口
		ServiceName: "deepseek-service",
		GroupName:   "DEFAULT_GROUP",
		Ephemeral:   true,
		Healthy:     true,
		Weight:      10,
		Metadata:    metadata,
	})

	if err != nil {
		log.Fatalf("注册服务实例失败: %v", err)
	}

	log.Println("服务实例注册成功！")
	// ...
}

```