## LLM 网关 endpoint 配置

[English](./endpoint.md) | 中文

本文档解释了如何在您的网关路由配置中为大语言模型 (LLM) 配置上游 endpoint。

### endpoint 结构

集群中的每个 endpoint 都由一个 `id` 定义，并且可以包含一个 `llm_meta` 块用于自定义行为。

```yaml
clusters:
  - name: "my_llm_cluster"
    endpoints:
      - id: "provider-1-main"
        socket_address:
           domains:
              - api.deepseek.com
        llm_meta:
          # ... 其他特定于 LLM 的配置在此处 ...
      - id: "provider-2-fallback"
        socket_address:
           domains:
              - api.openai.com/v1
        llm_meta:
          # ... 其他特定于 LLM 的配置在此处 ...
```

### `llm_meta` 配置字段

`llm_meta` 块包含所有与网关应如何处理此 LLM endpoint 相关的特定配置。

`fallback`

- **类型**: `boolean`
- **描述**: 决定如果在此endpoint上的所有重试尝试都失败后，网关是否应继续处理集群中的下一个 endpoint。
    - `true`: 如果此endpoint失败，网关将尝试下一个可用的 endpoint。
    - `false`: 如果此endpoint失败，则处理停止，并将最后一个错误返回给客户端。对于 fallback 链中的最后一个 endpoint ，此值应设置为 `false`。

`api_key`
- **类型**: `string`
- **描述**: 此 endpoint 使用的 API 密钥。当将请求转发到 LLM 服务时，此密钥将包含在请求头中。

`retry_policy`

- **类型**: `object`
- **描述**: 一个定义了当对此endpoint的请求失败时要使用的重试策略的对象。

   `retry_policy` 对象包含以下字段：

   `name`

  - **类型**: `string`
  - **描述**: 要使用的已注册重试策略的名称。名称不区分大小写。

   `config`

  - **类型**: `object`
  - **描述**: 一个键值对的映射，特定于所选的重试策略名称。

### 可用的重试策略

以下是您可以按名称指定的内置重试策略。

1. `CountBased` (基于次数)

    此策略在没有延迟的情况下重试固定的次数。

    #### 配置参数:

    - `times` (integer): 初始尝试失败后重试的次数。值为 `3` 意味着将有 1 次初始尝试和最多 3 次重试，总共 4 次尝试。

    示例:

    ```yaml
    retry_policy:
      name: "CountBased"
      config:
        times: 2
    ```

2. `ExponentialBackoff` (指数退避)

   此策略重试指定的次数，并随每次后续重试增加延迟时间。这是处理速率限制和临时网络问题的推荐策略。

    #### 配置参数:

    - `times` (integer): 初始尝试失败后重试的次数。
    - `initialInterval` (string): 第一次重试前初始等待时间的持续时长。
    - `maxInterval` (string): 重试之间的最大可能延迟。计算出的退避延迟将以此值为上限。
    - `multiplier` (float): 每次尝试后延迟时间乘以的因子。值为 `2.0` 将使延迟加倍。

    示例:

    ```yaml
    retry_policy:
      name: "ExponentialBackoff"
      config:
        times: 3
        initialInterval: 200ms
        maxInterval: 5s
        multiplier: 2.0
    ```

3. `NoRetry` (不重试)

   这是默认的重试策略。此策略仅执行初始尝试，如果失败则不执行任何重试。

    #### 配置参数:

    此策略不需要 `config` 块。

   示例:

   ```yaml
   retry_policy:
     name: "NoRetry"
   ```

### 完整配置示例

此示例展示了一个包含两个endpoint的集群。

第一个 (`deepseek-primary`) 使用 `ExponentialBackoff` 重试策略，并在失败时 fallback 到下一个 endpoint。

第二个 (`deepseek-fallback`) 是最终 endpoint，使用简单的 `CountBased` 重试策略，并禁用了fallback。

```yaml
clusters:
  - name: deepseek_cluster
    lb_policy: lb # 不需要为此集群启用负载均衡。
    endpoints:
      # --- 主endpoint ---
      - id: deepseek-primary
        socket_address:
          domains:
            - api.deepseek.com
        llm_meta:
          # 如果所有重试都失败，则移至下一个 endpoint。
          fallback: true
          # 此 endpoint 使用的 API 密钥。
          api_key: "your_deepseek_api_key_here"
          # 为主 endpoint 使用稳健的重试策略。
          retry_policy:
            name: ExponentialBackoff
            config:
              times: 3
              initialInterval: 200ms
              maxInterval: 8s
              multiplier: 2.5

      # --- fallback endpoint ---
      - id: openai-fallback
        socket_address:
          domains:
            - api.openai.com/v1
        llm_meta:
          # 这是最后的选择；不要再进一步 fallback。
          fallback: false
          # 此 endpoint 使用的 API 密钥。
          api_key: "your_openai_api_key_here"
          # 为 fallback endpoint 使用更简单、更快速的重试。
          retry_policy:
            name: CountBased
            config:
              times: 1
```