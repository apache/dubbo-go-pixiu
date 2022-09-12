[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-pixiu.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-pixiu)

[English](./README.md) | 中文

# 简介

**Dubbo-Go-Pixiu**(官网: https://dubbo.apache.org/zh/docs3-v2/dubbo-go-pixiu/) 是一款 Dubbo 生态下的高性能 API 网关和多语言解决方案 Sidecar

![](https://dubbo-go-pixiu.github.io/img/pixiu-dubbo-ecosystem.png)
Pixiu 是一款开源的 Dubbo 生态的 API 网关和 接入 dubbo 集群的语言解决方案。作为 API 网关形态， Pixiu 能接收外界的网络请求，将其转换为 dubbo 等协议请求，转发给背后集群；作为 Sidecar，Pixiu 期望可以代替代理服务注册到 Dubbo 集群，让多语言服务接入 Dubbo 集群提供更快捷的解决方案


## 快速开始

你可以在 https://github.com/apache/dubbo-go-pixiu-samples 中找到所有有关 pixiu 功能的案例，可以按照如下的步骤进行操作。

#### 进入示例代码目录

```
cd dubbogo/simple
```

可以使用 start.sh 脚本快速启动案例项目，可以执行如下命令来获得更多信息

```
./start.sh [action] [project]
./start.sh help
```

下列步骤中，我们将启动 body 案例项目

#### 准备配置文件和外部依赖 docker

使用 start.sh 的 prepare 命令来准备配置文件和外部docker依赖

```
./start.sh prepare body
```

如果想要手动准备文件，需要注意：
- 将 conf.yaml 中的 $PROJECT_DIR 修改为本地绝对路径

#### 启动 dubbo 服务或者 http 服务

```
./start.sh startServer body
```

#### 启动 pixiu

```
./start.sh startPixiu body
```

可以使用下列命令来手动启动 pixiu

```
 go run cmd/pixiu/*.go gateway start -c /[absolute-path]/dubbo-go-pixiu/samples/dubbogo/simple/body/pixiu/conf.yaml
```


#### 尝试请求

可以使用 curl 或者执行单元测试来验证一下

```
curl -X POST 'localhost:8881/api/v1/test-dubbo/user' -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header 'Content-Type: application/json' 
./start.sh startTest body
```

#### 清除

```
./start.sh clean body
```


## docker启动示例

#### 
```shell
docker pull dubbogopixiu/dubbo-go-pixiu:latest
```
```
docker run --name pixiuname -p 8888:8888 \
    -v /yourpath/conf.yaml:/etc/pixiu/conf.yaml \
    -v /yourpath/log.yml:/etc/pixiu/log.yml \
    dubbogopixiu/dubbo-go-pixiu:latest
```


## 特性

- 多协议支持：目前已支持 Http、Dubbo2、Triple、gRPC 协议代理和转换，其他协议持续集成中
- 安全认证：支持 HTTPS、JWT Token 校验等安全认证措施
- 注册中心集成：支持从 Dubbo 或 Spring Cloud 集群中获取服务元数据，支持 ZK、Nacos 注册中心
- 流量治理：集成 sentinel，支持多种协议限流
- 可观测性：集成 opentelemetry 和 jaeger，便于进行分布式链路追踪
- 自持 admin 和可视化界面：拥有 pixiu-admin 进行远程管理和可视化

### 控制面

Pixiu 控制面是 frok 自 [istio](https://github.com/istio/istio) v1.14.3 版本。提供包括服务发现、流量管理、安全等多种能力。

## 联系我们

项目在快速迭代中，欢迎使用， 欢迎给出建议或者提交pr。


### 社区

**官方钉钉群(31203920)**:

[![flowchart](./docs/images/group-pixiu-dingding.jpg)](docs/images/group-pixiu-dingding.jpg)
## License

Apache License, Version 2.0
