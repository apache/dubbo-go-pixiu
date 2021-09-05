# 概述

Pixiu 作为一个出色的网关，实现与 SpringCloud 服务对接能力。



## 文档框架

- 项目说明：背景，目的，预期能力，现状未来
- 方案设计：系统架构，核心技术方案设计，执行流程设计，启动加载流程设计
- 其他

# 项目说明

TODO:
- 流程图
    - 启动流程（UML）
    - 运行时流程（UML）
    - 停机流程（UML）
- 代码视角技术方案分层图



## 问题

-



## 功能说明



- [ ] 启动阶段：pixiu 启动阶段加载 sc 环境（可配置化）
    - [ ] 实现 CloudAdapter 启动逻辑
        - [ ] 加载 cluster，获取 sc 注册中心所有服务实例信息
        - [ ] 加载 router
            - [ ] 根据 配置的 service api 信息，结合 sc 服务实例信息创建 sc router 信息（router match，router action）
    - [ ] 新增 拉取指定注册中心服务信息
        - [ ] **nacos** 实现
            - [ ] 实现 load.go 方法：LoadAllServices，GetCluster
                - [ ] 获取所有的 Service 服务信息
                - [ ] 获取所有的 ServiceId
            - [ ] 本地缓存所有服务信息
            - [ ] 提供定时任务刷新本地服务信息
            - [ ] 监听nacos服务变更事件，更新本地服务
        - [ ] [delay]zk
        - [ ] [delay]consul
    - [ ] 新增 本地缓存，缓存服务实例信息，区分不同注册中心
        - [ ] 提供查询所有服务实例信息接口
        - [ ] 提供查询所有 serviceId 接口
        - [ ] 提供查询指定服务实例根据服务ID 接口
    - [ ] 新增 定时刷新远程与本地服务信息，统一触发
    - [ ] 新增 监听注册中心远程服务变更事件（回调），刷新本地服务实例缓存，统一提供刷新接口
- [ ] 运行时：pixiu 运行时 sc 场景兼容
    - [ ] [improve]新增 sc 请求处理代理 fitler
- [ ] 关闭：
    - [ ] 优雅停机实现



# 方案设计

## 相关代码

- Pixiu 集成 SpringCloud 插件 ：[cloud.go](../cloud.go)
- Nacos 注册中心相关：[nacos](../../../registry/nacos)
    - listener : 服务变更事件监听
    - nacos: 服务发现能力
    - refresh: 服务定时刷新？？？
- Java 实现的SpringCloud 示例，方便调试 ：[samples/springcloud](../../../../samples/springcloud)
- ...

## 启动流程
![alt PixiuStart](pixiu_start_sc.svg "Pixiu Start Process")

## 配置文件

Pixiu SpringCloud 配置示例：

```yaml
static_resources:
  listeners:
    - name: "net/http"
  clusters:
    - name: "test_dubbo"
  adapters:
    - name: "dgp.cluster_adapters.spring_cloud"
      config:
        registries:
          "zookeeper":
            protocol: "zookeeper"
            timeout: "3s"
            address: "127.0.0.1:2181"
            username: ""
            password: ""
```







