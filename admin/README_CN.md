# Pixiu-Admin 部署与配置文档

[English](README.md) | **中文**

**Pixiu-Admin** 是基于 **Pixiu** 生态的管理平台，主要用于配置、监控和管理 Pixiu 网关资源。通过 Web 用户界面和 RESTful API
提供集中式管理功能。本文档介绍了如何在 Linux 环境下部署和配置 Pixiu-Admin。

后端API接口文档请参考 [API.md](../admin/API_CN.md)。

## 部署文档

### 使用 Docker 启动

首先，确保您在项目的根目录（包含 Dockerfile 的目录）下，使用以下命令启动 Pixiu-Admin：

```bash
docker-compose up -d
```

### 使用源码部署

如果不使用 Docker，可以通过以下命令下载源代码：

```bash
git clone https://github.com/apache/dubbo-go-pixiu
```

### 部署 etcd

手动部署 etcd 服务，使用以下命令：

```bash
docker run -d -p2379:2379 --env ALLOW_NONE_AUTHENTICATION=yes --name etcd bitnami/etcd
```

对于 M1/M1 Pro 用户，使用以下命令：

```bash
docker run -d -p2379:2379 --platform linux/amd64 --env ALLOW_NONE_AUTHENTICATION=yes --name etcd bitnami/etcd:3.5.1
```

### 运行 Admin

#### 源代码运行

进入项目目录并运行：

```bash
cd dubbo-go-pixiu
# 直接运行
go run ./cmd/admin/admin.go -c /your/local/path/conf.yaml
# 后台运行
nohup go run ./cmd/admin/admin.go -c /your/local/path/conf.yaml &
```

#### 运行 Pixiu

默认配置见 [pixiu_with_admin_config.yaml](../configs/pixiu_with_admin_config.yaml)

```bash
go run ./cmd/pixiu/pixiu.go gateway start -c ./configs/pixiu_with_admin_config.yaml
```

### 测试运行 admin-web

进入 `web` 目录并安装依赖：

```bash
cd ./admin/web/
yarn install  # 安装依赖
yarn run serve  # 测试运行
```

#### admin-web 配置

编辑 `web` 目录下的 `vue.config.js`，配置后端服务地址：

```
devServer: {
    host: '0.0.0.0',
        port: 8080,  // Web app address
        hot: true,
        https: false,
        open: false,
        disableHostCheck: true,
        proxy: {
        "/config": {
            target: "http://127.0.0.1:8081",  // Backend service address
                ws: true,  // Enable websockets
                changeOrigin: true,  // Enable proxy
        }
    }
}
```

运行成功后，可以在浏览器访问 [http://127.0.0.1:8081/login.html#/Overview](http://127.0.0.1:8081/login.html#/Overview)。

## 二、相关操作

### 管理映射（Resource）

#### 创建映射配置

1. 点击 "映射配置"，进入映射配置列表界面。
2. 点击右上角的新增按钮，创建新的映射配置。

![1.png](../docs/images/admin/1.png)

在代码编辑器中键入映射配置，点击确认创建。

![2.png](../docs/images/admin/2.png)

#### 映射配置示例

```yaml
path: '/api/v1/test-dubbo/user'
type: restful
description: user
filters:
  - filter0
methods:
  - httpVerb: GET
    onAir: true
    timeout: 1000ms
    inboundRequest:
      requestType: http
      queryStrings:
        - name: name
          required: true
    integrationRequest:
      requestType: dubbo
      mappingParams:
        - name: queryStrings.name
          mapTo: 0
          mapType: "java.lang.String"
      applicationName: "UserProvider"
      interface: "com.ic.user.UserProvider"
      method: "GetUserByName"
      group: "test"
      version: 1.0.0
      clusterName: "test_dubbo"
  - httpVerb: POST
    onAir: true
    timeout: 10s
    inboundRequest:
      requestType: http
    integrationRequest:
      requestType: dubbo
      mappingParams:
        - name: requestBody._all
          mapTo: 0
          mapType: "object"
      applicationName: "UserProvider"
      interface: "com.ic.user.UserProvider"
      method: "CreateUser"
      group: "test"
      version: 1.0.0
      clusterName: "test_dubbo"
```

#### 查看和删除映射

1. 映射列表刷新后，可以查看和删除配置。

![3.png](../docs/images/admin/3.png)

2. 点击删除会删除该映射配置，点击查看跳转至配置详情界面。

![4.png](../docs/images/admin/4.png)

#### 编辑映射

可以在编辑区修改映射配置，点击 "修改" 保存修改。

![5.png](../docs/images/admin/5.png)

#### 方法映射

点击新增按钮，添加新的方法映射。编辑示例如下：

```yaml
httpVerb: PUT
onAir: true
timeout: 10s
inboundRequest:
  requestType: http
integrationRequest:
  requestType: dubbo
  mappingParams:
    - name: requestBody._all
      mapTo: 0
      mapType: "object"
  applicationName: "UserProvider"
  interface: "com.ic.user.UserProvider"
  method: "CreateUser"
  group: "test"
  version: 1.0.0
  clusterName: "test_dubbo"
```

点击确认后，方法映射会出现在列表中。

![6.png](../docs/images/admin/6.png)

#### 查看和删除方法映射

可以查看方法映射的详细信息或删除该映射。

比如，将第二个方法映射的 httpVerb 从 POST 修改为 DELETE。

![7.png](../docs/images/admin/7.png)

注意：id不能进行修改，即使修改保存后也会改变为旧值。点击确定后，会更新该方法映射。

![8.png](../docs/images/admin/8.png)

### 管理插件组

#### 创建插件组

点击左侧插件配置菜单项，可以查看插件相关配置。

![9.png](../docs/images/admin/9.png)

点击右上方新增，可以创建新的插件组。

![10.png](../docs/images/admin/10.png)

插件组配置示例：

```yaml
groupName: "group2"
plugins:
  - name: "rate limit"
    version: "0.1.0"
    priority: 1000
    externalLookupName: "ExternalPluginRateLimit"
  - name: "log"
    version: "0.2.0"
    priority: 2000
    externalLookupName: "ExternalPluginLog"
```

保存后，列表刷新展示新创建的插件组。

![11.png](../docs/images/admin/11.png)

和映射配置类似，点击查看会弹出编辑框，可以对插件组配置进行修改；点击删除会删除整个插件组配置。

#### 查看和删除插件组

点击查看可以编辑插件组配置，点击删除删除插件组。

### 管理限流配置

#### 配置限流

点击 "限流配置" 菜单，进行限流组件的配置。具体配置示例：

![12.png](../docs/images/admin/12.png)

```yaml
resources:
  - name: test-http
    items:
      - pattern: /api/v1/test-dubbo/user
      - matchStrategy: 1
        pattern: /api/*/test-dubbo/user
rules:
  - flowRule:
      resource: ""
      tokencalculatestrategy: 0
      threshold: 100
      enable: true
```

点击保存后，限流配置生效。

![13.png](../docs/images/admin/13.png)

## 三、Pixiu 远程配置

### 启动和配置

启动 Pixiu 并指定配置文件，在配置文件中，定义 etcd 地址和配置路径：

```yaml
api_meta_config:
  address: "127.0.0.1:2379"
  api_config_path: "/pixiu/config/api"
```

### 测试

在 admin 中创建资源配置，并使用 `curl` 测试 Pixiu 转发功能：

```bash
curl "http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc"
curl -X POST "http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc"
```

如果请求未找到对应服务，返回错误信息；若配置正确，则返回服务响应。

## 许可证

本项目采用 Apache License 2.0 开源许可。

