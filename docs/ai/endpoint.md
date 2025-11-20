
## LLM Gateway Endpoint Configuration

English | [中文](./endpoint_CN.md)

This document explains how to configure upstream endpoints for Large Language Models (LLMs) within your gateway's routing configuration.

### Endpoint Structure

Each endpoint within a cluster is defined by an id and can contain an llm_meta block for custom behavior.clusters:

```yaml
clusters:
  - name: "my_llm_cluster"
    endpoints:
      - id: "provider-1-main"
        socket_address:
           domains:
              - api.deepseek.com
        llm_meta:
          # ... other LLM-specific configuration goes here ...
      - id: "provider-2-fallback"
        socket_address:
           domains:
              - api.openai.com/v1
        llm_meta:
          # ... other LLM-specific configuration goes here ...
```

### `llm_meta` Configuration Fields

The llm_meta block holds all the configuration specific to how the gateway should treat this LLM endpoint.

`fallback`

- Type: `boolean`
- Description: Determines if the gateway should proceed to the next endpoint in the cluster if all retry attempts on this endpoint fail. When the value is `true`, and if this endpoint fails, the gateway will attempt the next available endpoint. When the value is `false`, and if this endpoint fails, the process stops, and the last error is returned to the client.

`api_key`
- Type: `string`
- Description: The API key to be used by this endpoint. This key will be included in the request headers when forwarding requests to the LLM service.

`retry_policy`

- Type: `object`
- Description: An object that defines the retry strategy to use if a request to this endpoint fails.

   The `retry_policy` object contains the following fields:

   `name`

  - Type: `string`
  - Description: The name of the registered retry policy to use. The name is case-insensitive.

   `config`

  - Type: `object`
  - Description: A map of key-value pairs specific to the chosen retry policy name.

### Available Retry Policies

Here are the built-in retry policies you can specify by name.

1. `CountBased`

   This policy retries a fixed number of times with no delay between attempts.

   #### Config Parameters:

    - times (integer): The number of times to retry after the initial attempt fails. A value of 3 means there will be 1 initial attempt and up to 3 retries, for a total of 4 attempts.

   Example:

    ```yaml
    retry_policy:
      name: "CountBased"
      config:
        times: 2
    ```

2. `ExponentialBackoff`

   This policy retries a specified number of times, increasing the delay between each subsequent retry. This is the recommended strategy for handling rate limits and transient network issues.

   #### Config Parameters:

    - times (integer): The number of times to retry after the initial attempt fails.
    - initialInterval (string): The duration of the initial wait time before the first retry (e.g., "100ms", "1s").
    - maxInterval (string): The maximum possible delay between retries. The calculated backoff delay will be capped at this value.
    - multiplier (float): The factor by which the delay is multiplied after each attempt. A value of 2.0 will double the delay.

   Example:
    ```yaml
    retry_policy:
      name: "ExponentialBackoff"
      config:
        times: 3
        initialInterval: 200ms
        maxInterval: 5s
        multiplier: 2.0
    ```

3. `NoRetry`

   This is the default retry policy. This policy only performs the initial attempt and does not perform any retries if it fails.

   #### Config Parameters:

   This policy does not require a `config` block.

   Example:
    ```yaml
    retry_policy:
      name: "NoRetry"
    ```

### Complete Configuration Example

This example shows a cluster with two endpoints.

The first (deepseek-primary) uses an `ExponentialBackoff` retry policy and will fall back to the next endpoint on failure.

The second (deepseek-fallback) is the final endpoint, using a simple `CountBased` retry policy and with fallback disabled.

```yaml
clusters:
  - name: deepseek_cluster
    lb_policy: lb # No need to enable load balancing for this cluster.
    endpoints:
      # --- Primary Endpoint ---
      - id: deepseek-primary
        socket_address:
          domains:
            - api.deepseek.com
        llm_meta:
          # If all retries fail, move to the next endpoint.
          fallback: true
          # Your API key for this endpoint.
          api_key: "your_deepseek_api_key_here"
          # Use a robust retry strategy for the primary endpoint.
          retry_policy:
            name: ExponentialBackoff
            config:
              times: 3
              initialInterval: 200ms
              maxInterval: 8s
              multiplier: 2.5

      # --- Fallback Endpoint ---
      - id: openai-fallback
        socket_address:
          domains:
            - api.openai.com/v1
        llm_meta:
          # This is the last resort; do not fall back further.
          fallback: false
          # Your API key for this endpoint.
          api_key: "your_openai_api_key_here"
          # Use a simpler, faster retry for the fallback.
          retry_policy:
            name: CountBased
            config:
              times: 1
```