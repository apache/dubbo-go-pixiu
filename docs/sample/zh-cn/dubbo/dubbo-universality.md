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
