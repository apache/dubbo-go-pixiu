## LLM Service Discovery and Registration

English | [中文](registry_CN.md)

This document aims to guide LLM service providers on how to dynamically register their service instances with the LLM Gateway via a Nacos registry. By following these guidelines, the gateway will be able to automatically discover your service and apply appropriate routing, retry, and fallback strategies based on the metadata you provide.

### Registration Mechanism Overview

The core mechanism of service discovery is that your LLM service registers as a **Nacos instance** and provides a specific set of **metadata** upon registration. The LLM Gateway listens for service changes in Nacos, reads this metadata, and dynamically converts it into a fully functional gateway `endpoint` configuration.

A basic Nacos registration request includes the following key information:

- **`ServiceName`**: The name of your service collection (e.g., `deepseek-service`).
- **`Ip` & `Port`**: The network address where your service instance listens for traffic.
- **`Metadata`**: A collection of key-value pairs, **which is crucial for configuring all gateway behaviors**.

### `metadata` Configuration Fields

All gateway-specific configurations are passed through the `metadata` field of the Nacos instance. Below are all the supported `metadata` keys and their descriptions.

`cluster`

- **Type**: `string`
- **Required**: Yes
- **Description**: Defines which cluster (`cluster`) in the gateway this endpoint should belong to. The gateway aggregates service instances with the same `cluster` name based on this value.
- **Example**: `"deepseek_cluster"`

`id`

- **Type**: `string`
- **Required**: Yes
- **Description**: The unique identifier for the endpoint. It must be unique within its `cluster`.
- **Example**: `"deepseek-primary-instance-1"`

`ip`

- **Type**: `string`
- **Required**: No
- **Description**: An optional parameter, useful when your service needs to register an address with the gateway that is different from its internal IP (e.g., a publicly accessible IP). If not provided, the gateway will use the default value `0.0.0.0`.
- **Example**: `"203.0.113.55"`

`port`

- **Type**: `string` (representing an integer)
- **Required**: No
- **Description**: An optional parameter to override the `Port` registered by the Nacos instance itself. Its usage is similar to the `ip` field.
- **Example**: `"9000"`

`name`

- **Type**: `string`
- **Required**: No
- **Description**: A human-readable name to identify this endpoint. Primarily used for logging and monitoring.
- **Example**: `"DeepSeek V2 Chat (Primary)"`

`address`

- **Type**: `string`
- **Required**: No
- **Description**: A string split by comma, each string stands for a address 
- **Example**: `"api.deepseek.com"`

`llm-meta.fallback`

- **Type**: `string` ("true" or "false")
- **Required**: No, defaults to `"false"`
- **Description**: Determines whether the gateway should proceed to the next endpoint in the cluster if all retry attempts on this endpoint fail.

`llm-meta.api_key`

- **Type**: `string` 
- **Required**: No
- **Description**: The API key to be used by this endpoint.

`llm-meta.retry_policy.name`

- **Type**: `string`
- **Required**: No, defaults to `"NoRetry"`
- **Description**: Specifies the name of the retry policy to use when a request to this endpoint fails. The name is case-insensitive.

`llm-meta.retry_policy.config`

- **Type**: `string` (JSON format)
- **Required**: No
- **Description**: A JSON string containing configuration parameters specific to the selected retry policy.

### Retry Policy Metadata

You can configure retry behavior by using a combination of `llm-meta.retry_policy.name` and `llm-meta.retry_policy.config`.

1.  `CountBased`

    - **name**: `"CountBased"`
    - **config**: A JSON string containing a `times` field.
        - `times` (integer): Number of retry attempts.
    - **Example**:
        - `llm-meta.retry_policy.name`: `"CountBased"`
        - `llm-meta.retry_policy.config`: `{"times": 2}`

2.  `ExponentialBackoff`

    - **name**: `"ExponentialBackoff"`
    - **config**: A JSON string containing `times`, `initialInterval`, `maxInterval`, and `multiplier` fields.
        - `times` (integer): Number of retries.
        - `initialInterval` (string): Initial wait duration (e.g., "200ms").
        - `maxInterval` (string): Maximum wait duration (e.g., "5s").
        - `multiplier` (float): The multiplier factor for the delay.
    - **Example**:
        - `llm-meta.retry_policy.name`: `"ExponentialBackoff"`
        - `llm-meta.retry_policy.config`: `{"times": 3, "initialInterval": "200ms", "maxInterval": "5s", "multiplier": 2.0}`

3.  `NoRetry`

    - **name**: `"NoRetry"`
    - **config**: This policy does not require a `config` field.

### Complete Registration Example (Go)

Detailed usage please refer to our [official samples](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/llm/nacos)。

The configuration file of pixiu, need to enable the adapter of llmregistrycenter, as follows:

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

This example demonstrates how to register a fully-featured LLM service instance using the Nacos Go SDK. The instance will be configured as the primary endpoint in the `deepseek_cluster`, using an exponential backoff retry policy, and will fall back to the next service in the cluster upon failure.

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// ... (Code for creating the Nacos client is omitted here)
	// client, err := createNacosClient()

	// 1. Prepare the JSON configuration for the retry policy
	retryConfig := map[string]interface{}{
		"times":           3,
		"initialInterval": "200ms",
		"maxInterval":     "8s",
		"multiplier":      2.5,
	}
	retryConfigJSON, _ := json.Marshal(retryConfig)

	// 2. Construct the metadata containing all gateway configurations
	metadata := map[string]string{
		// --- Core Endpoint Configuration ---
		"cluster": "deepseek_cluster",
		"id":      "deepseek-primary",
		"name":    "DeepSeek V2 Chat (Primary)",

		// Optional (use ip+port or address): The instance's IP and Port
		"ip":   "203.0.113.55",
		"port": "9000",

		// Optional (use ip+port or address): address field
		"address": "api.deepseek.com",

		// --- LLM-Specific Metadata ---
		"llm-meta.fallback": "true",

		// API Keys in JSON string format
		"llm-meta.api_keys": "key-xxxxxxxx",

		// --- Retry Policy Configuration ---
		"llm-meta.retry_policy.name":   "ExponentialBackoff",
		"llm-meta.retry_policy.config": string(retryConfigJSON),
	}

	// 3. Register the Nacos instance
	// Note: The Ip and Port here are the actual listening addresses of the service instance,
	// while the ip and port in the metadata are the addresses you want the gateway to access.
	_, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.1.10", // The service's internal IP
		Port:        8001,           // The service's internal port
		ServiceName: "deepseek-service",
		GroupName:   "DEFAULT_GROUP",
		Ephemeral:   true,
		Healthy:     true,
		Weight:      10,
		Metadata:    metadata,
	})

	if err != nil {
		log.Fatalf("Failed to register service instance: %v", err)
	}

	log.Println("Service instance registered successfully!")
	// ...
}

```