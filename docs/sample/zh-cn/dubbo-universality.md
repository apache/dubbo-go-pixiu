# 使用 dubbo 通用性请求

> POST 请求 [samples](https://github.com/dubbogo/dubbo-go-proxy-samples/tree/master/dubbo/apiconfig/proxy)

## 建议

> 使用此方式，你能够给一个集群定义一个接口来请求对应 dubbo 提供的服务

### 接口配置

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/:application/:interface'
    type: restful
    description: common
    methods:
      - httpVerb: POST
        onAir: true
        timeout: 100s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody.values
              mapTo: 0
              opt:
                open: true
                usable: true
                name: values
            - name: requestBody.types
              mapTo: 1
              opt:
                open: true
                name: types
            - name: uri.application
              mapTo: 2
              opt:
                open: true
                name: application
            - name: uri.interface
              mapTo: 3
              opt:
                open: true
                name: interface
            - name: queryStrings.method
              mapTo: 4
              opt:
                open: true
                name: method
            - name: queryStrings.group
              mapTo: 5
              opt:
                open: true
                name: group
            - name: queryStrings.version
              mapTo: 6
              opt:
                open: true
                name: version
          paramTypes: ["object", "object", "string", "string", "string", "string", "string"]
          # 这个是必须注意的，实际传给 dubbo 的 paramTypes，在最终调用的时候优先级高于 paramTypes。
          toParamTypes: ["string"]
          clusterName: "test_dubbo"
```

### 测试例子

- 单个 string 参数

```bash
curl host:port/api/v1/test-dubbo/UserService/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByName -X POST -d '{"types":["string"],"values":"tc"}' --header "Content-Type: application/json"
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
curl host:port/api/v1/test-dubbo/UserService/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByCode -X POST -d '{"types":["int"],"values":1}' --header "Content-Type: application/json"
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
curl host:port/api/v1/test-dubbo/UserService/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=UpdateUserByName -X POST -d '{"types":["string","body"],"values":["tc",{"id":"0001","code":1,"name":"tc","age":15}]}' --header "Content-Type: application/json"
```

result

```bash
true
```

### 特殊配置

#### 可配码

```go
const (
	optionKeyTypes       = "types"
	optionKeyGroup       = "group"
	optionKeyVersion     = "version"
	optionKeyInterface   = "interface"
	optionKeyApplication = "application"
	optionKeyMethod      = "method"
	optionKeyValues      = "values"
)
```

#### 选择项

组装泛化调用的参数

```go
// GenericService uses for generic invoke for service call
type GenericService struct {
	Invoke       func(ctx context.Context, req []interface{}) (interface{}, error) `dubbo:"$invoke"`
	referenceStr string
}
```

- types

> dubbo 泛化类型

用于 dubbogo `GenericService#Invoke` 函数的第二个参数。

- method

用于 dubbogo `GenericService#Invoke` 函数的第一个参数。

- group

Dubbo 组配置 `ReferenceConfig#Group`。

- version

Dubbo 版本配置 `ReferenceConfig#Version`。

- interface

Dubbo 接口配置 `ReferenceConfig#InterfaceName`。

- application

目前暂时用于缓存，索引的一部分查找对应的缓存对象。

- values

值的处理，用于 `GenericService#Invoke` 函数的第三个参数。

#### 解释

##### 单个参数

请求体

```json
{
    "types": [
        "string"
    ],
    "values": "tc"
}
```

```yaml
            - name: requestBody.types
              mapTo: 1
              opt:
                open: true
                name: types
```

- `requestBody.types` 意味着对 body 的 json 内容取 key 的值为 types。
- `opt.name` 意味着扩展名称，和 proxy 提供的默认实现匹配。
- `opt.open` 打开，目前是 `true` 才会创建对应的扩展，如果后续配置缩减的话考虑删除。
- `opt.usable` 表示上游服务是否需要这个参数，对应代码里面 `setTarget` 这个行为，默认加了扩展是 `false`，即这个字段只有行为，不会成为 RPC 的参数。

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

请注意这种特殊情况的配置目前自由度不是很高，如果有不能满足的场景请及时反馈到[问题](https://github.com/dubbogo/dubbo-go-proxy/issues)

[上一页](./dubbo.md)