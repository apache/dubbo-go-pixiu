# 使用 dubbo 通用性请求

> POST 请求 [samples](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/dubbogo/simple/proxy)

## 直连泛化调用契约（破坏性变更）

Pixiu 直连泛化调用不再使用请求体里的 `types` 作为方法签名来源。直连模式现在要求：

- `integrationRequest.url`
- `integrationRequest.interface`
- `integrationRequest.method`
- `integrationRequest.parameterTypes`
- `integrationRequest.serialization`

`mappingParams` 仍然只负责传值，不再定义方法签名。

对于 Triple 直连泛化调用，provider 必须暴露 generic `$invoke` 入口。当前
dubbo-go generic client 在该路径下使用 non-IDL 模式。IDL 生成的 Triple handler
通常只注册具体 RPC 方法，不会暴露 `$invoke`；因此对这类 provider 发起
`$invoke` 泛化调用时没有匹配的 Triple handler，可能返回 `404 Not Found`。

## 建议

> 使用此方式，你能够给一个集群定义一个接口来请求对应 dubbo 提供的服务
> 下面的示例使用 registry 模式。在该模式下，`opt.types` 仍可提供泛化调用签名。
> 直连模式必须改用 `integrationRequest.parameterTypes`。

### 接口配置

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/:interface'
    type: restful
    description: common
    methods:
      - httpVerb: POST
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody.values
              mapTo: opt.values
            - name: requestBody.types
              mapTo: opt.types
            - name: uri.interface
              mapTo: opt.interface
            - name: queryStrings.method
              mapTo: opt.method
            - name: queryStrings.group
              mapTo: opt.group
            - name: queryStrings.version
              mapTo: opt.version
          clusterName: "test_dubbo"
```

### 测试例子

- 单个 string 参数

```bash
curl host:port/api/v1/test-dubbo/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByName -X POST -d '{"types":["string"],"values":"tc"}' --header "Content-Type: application/json"
```

result

```json
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T20:54:38.746+08:00"
}
```

- 单个 int 参数

```bash
curl host:port/api/v1/test-dubbo/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByCode -X POST -d '{"types":["int"],"values":1}' --header "Content-Type: application/json"
```

result

```json
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T20:54:38.746+08:00"
}
```

- 多个参数

```bash
curl host:port/api/v1/test-dubbo/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=UpdateUserByName -X POST -d '{"types":["string","body"],"values":["tc",{"id":"0001","code":1,"name":"tc","age":15}]}' --header "Content-Type: application/json"
```

result

```bash
true
```

### 特殊配置

#### 可配码

支持的 `mapTo` 选项：

```yaml
- opt.types
- opt.group
- opt.version
- opt.interface
- opt.method
- opt.values
```

#### generic 模式

`integrationRequest.generic` 选择泛化调用的方式，取值：

```yaml
- "true"        # 默认，map 泛化：values 作为有类型的参数发送，结果以 map 返回
- gson          # 请求以 JSON 文本发送，结果以 JSON 文本返回
- protobuf-json # 请求以 protobuf JSON 文本发送，结果以 protobuf JSON 文本返回
- bean          # Java bean 兼容方式
```

map 方式沿用上面的 `opt.values` 与 `opt.types` 约定：`values` 是参数列表，`types` 是每个参数的 Java 类型名。

`gson` 与 `protobuf-json` 方式把整个请求消息作为 JSON 文本发送，因此 `values` 只能有一个元素，`parameterTypes` 也只能声明一个类型，该类型描述请求消息本身：

```yaml
integrationRequest:
  requestType: dubbo
  interface: com.example.Greeter
  method: SayHello
  protocol: tri
  serialization: hessian2
  generic: protobuf-json
  parameterTypes:
    - com.example.HelloRequest
  mappingParams:
    - name: requestBody.values
      mapTo: opt.values
```

请求体 `{"values":{"name":"test"}}` 会以 `{"name":"test"}` 作为唯一的泛化参数发送，provider 返回的 JSON 文本直接作为响应体写出，不会被包成带引号的字符串。

provider 侧要求：Triple 场景下 provider 需要声明 `serialization: hessian2`，并提供泛化入口 `$invoke`。dubbo-go 的非 IDL 服务默认提供该入口，IDL 导出的服务在 [apache/dubbo-go#3752](https://github.com/apache/dubbo-go/pull/3752) 之后同样提供。

#### 选择项

在mapTo 里面使用特定的关键字(列表如下)，貔貅可以自动组装泛化调用的参数

```go
// GenericService uses for generic invoke for service call
type GenericService struct {
	Invoke func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) `dubbo:"$invoke"`
}
```

- opt.types

> dubbo 泛化类型

当未配置 `integrationRequest.parameterTypes` 时，作为 dubbogo `GenericService#Invoke` 的 `types` 参数。
直连泛化调用请改用 `integrationRequest.parameterTypes` 声明方法签名。

- opt.method

作为 dubbogo `GenericService#Invoke` 的 `methodName` 参数。

- opt.group

Dubbo reference group。

- opt.version

Dubbo reference version。

- opt.interface

Dubbo service interface。

- opt.values

作为 dubbogo `GenericService#Invoke` 的 `args` 参数。

#### 解释

##### 单个参数

请求体

```json
{
    "types": ["string"],
    "values": "tc"
}
```

```yaml
  - name: requestBody.types
    mapTo: opt.types
```

- `requestBody.types` 表示读取请求体里的 `types` 字段。
- `opt.types` 表示在当前 registry 模式示例中，将该值作为泛化调用的 `types` 参数。

##### 多个参数

请求体

```json
{
  "types": [
    "java.lang.String",
    "object"
  ],
  "values": [
    "tc",
    {
      "id": "0001",
      "code": 1,
      "name": "tc",
      "age": 99
    }
  ]
}
```

请注意这种特殊情况的配置目前自由度不是很高，如果有不能满足的场景请及时反馈到[问题](https://github.com/apache/dubbo-go-pixiu/issues)

[上一页](dubbo.md)
