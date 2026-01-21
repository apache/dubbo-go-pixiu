# gRPC 代理过滤器 (dgp.filter.grpc.proxy)

[English](grpcproxy.md) | 中文

---

## 概述

`dgp.filter.grpc.proxy` 过滤器为 Pixiu 网关提供 gRPC 代理功能，并支持 **gRPC Server Reflection**。该特性使得网关可以在不需要预编译 proto 文件的情况下，动态解析和检查消息内容。

### 核心特性

- **三种反射模式**: 根据需求在性能和功能之间选择
- **动态消息解码**: 在运行时解析和检查 gRPC 消息，无需 proto 文件
- **基于 TTL 的缓存**: 高效的方法描述符缓存，自动清理过期条目
- **协议检测**: 同时支持 gRPC 和 Triple 协议
- **优雅降级**: 混合模式提供自动回退到透传模式

---

## 反射模式

### 透传模式（Passthrough，默认）

执行透明的二进制代理，不解码消息。

**特点:**
- 最高性能（无消息解析开销）
- 无法检查消息内容
- 纯二进制转发

**使用场景:**
- 不需要消息检查的高吞吐场景
- 仅基于服务/方法名称的简单路由
- 对向后兼容性要求极高的场景

**配置:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "passthrough"  # 或省略（默认值）
```

### 反射模式（Reflection）

使用 gRPC Server Reflection API 动态解码和检查消息内容。

**特点:**
- 运行时完整消息解码
- 支持基于内容的路由和过滤
- 需要后端服务器支持反射
- 由于反射调用，有轻微性能开销

**使用场景:**
- 基于内容的路由（根据消息字段路由）
- 字段级别的过滤或转换
- 完整消息检查的日志记录和调试
- 需要消息验证的 API 网关场景

**配置:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
      descriptor_cache_ttl: 300  # 5分钟缓存
```

### 混合模式（Hybrid）

先尝试反射，失败时回退到透传模式。

**特点:**
- 两全其美：可用时进行内容检查
- 对不支持反射的服务优雅降级
- 反射超时防止阻塞

**使用场景:**
- 反射支持程度不一的混合环境
- 迁移场景（逐步启用反射）
- 需要高可用的生产环境

**配置:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "hybrid"
      reflection_timeout: 5s  # 等待反射的最大时间
```

---

## 配置

### 完整配置示例

```yaml
static_resources:
  listeners:
    - name: "grpc-gateway"
      protocol_type: "GRPC"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8882
      filter_chains:
        filters:
          - name: dgp.filter.network.grpcconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: "/echo.EchoService/"
                    route:
                      cluster: "echo-grpc"
              grpc_filters:
                - name: dgp.filter.grpc.proxy
                  config:
                    # 反射模式（默认: "passthrough"）
                    reflection_mode: "reflection"

                    # 方法描述符缓存 TTL（默认: 300秒）
                    descriptor_cache_ttl: 300

                    # 启用 Triple 协议检测（默认: false）
                    enable_protocol_detection: true

                    # 混合模式的反射超时（默认: 5秒）
                    reflection_timeout: 5s

                    # TLS 配置（可选）
                    enable_tls: false
                    tls_cert_file: ""
                    tls_key_file: ""

                    # 连接设置（可选）
                    keepalive_time: 300s
                    keepalive_timeout: 5s
                    connect_timeout: 5s
                    max_concurrent_streams: 0

  clusters:
    - name: "echo-grpc"
      lb_policy: "RoundRobin"
      endpoints:
        - socket_address:
            address: 127.0.0.1
            port: 50051
            protocol_type: "GRPC"
```

### 配置字段

| 字段 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `reflection_mode` | string | `"passthrough"` | 反射模式: `"passthrough"`、`"reflection"` 或 `"hybrid"` |
| `descriptor_cache_ttl` | int | `300` | 方法描述符缓存 TTL（秒） |
| `enable_protocol_detection` | bool | `false` | 启用 Triple 协议检测 |
| `reflection_timeout` | string | `"5s"` | 混合模式下的反射超时时间 |
| `enable_tls` | bool | `false` | 启用后端连接的 TLS |
| `tls_cert_file` | string | `""` | TLS 证书文件路径 |
| `tls_key_file` | string | `""` | TLS 密钥文件路径 |
| `keepalive_time` | string | `"300s"` | 后端连接的保活时间 |
| `keepalive_timeout` | string | `"5s"` | 保活超时时间 |
| `connect_timeout` | string | `"5s"` | 连接超时时间 |
| `max_concurrent_streams` | uint32 | `0`（无限制） | 最大并发流数 |

---

## 启用服务器反射

要使用 `reflection` 或 `hybrid` 模式，你的后端 gRPC 服务器必须启用 Server Reflection。

### Go (gRPC-Go)

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

func main() {
    server := grpc.NewServer()

    // 注册你的服务
    echo.RegisterEchoServiceServer(server, &echoServer{})

    // 启用服务器反射（重要！）
    reflection.Register(server)

    // 启动服务器
    lis, _ := net.Listen("tcp", ":50051")
    server.Serve(lis)
}
```

### Java (gRPC-Java)

```java
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.reflection.v1alpha.ServerReflectionGrpc;

public class EchoServer {
    public static void main(String[] args) throws Exception {
        Server server = ServerBuilder
            .forPort(50051)
            .addService(new EchoServiceImpl())
            // 启用服务器反射
            .addService(ServerReflectionGrpc.newInstance())
            .build();

        server.start();
        server.awaitTermination();
    }
}
```

### Python (gRPC-Python)

```python
import grpc
from grpc_reflection.v1alpha import reflection
from concurrent import futures

class EchoService(echo_pb2_grpc.EchoServiceServicer):
    # ... 实现 ...

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    echo_pb2_grpc.add_EchoServiceServicer_to_server(EchoService(), server)

    # 启用服务器反射
    reflection.enable_server_reflection(
        service_names=[echo.DESCRIPTOR.services_by_name['EchoService'].full_name,
                       reflection.SERVICE_NAME],
        server=server
    )

    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()
```

---

## 模式对比

| 特性 | 透传模式 | 反射模式 | 混合模式 |
|------|---------|---------|---------|
| **性能** | ⭐⭐⭐⭐⭐（最佳） | ⭐⭐⭐（良好） | ⭐⭐⭐⭐（更好） |
| **消息检查** | ❌ 否 | ✅ 是 | ✅ 是（可用时） |
| **需要反射** | ❌ 否 | ✅ 是 | ❌ 可选 |
| **回退支持** | N/A | ❌ 否 | ✅ 是 |
| **使用场景** | 高性能代理 | 基于内容的路由 | 混合环境 |

---

## 描述符缓存

反射模式使用基于 TTL 的缓存来存储从后端服务器获取的方法描述符。这减少了反射请求的数量并提高了性能。

### 缓存行为

| 操作 | 描述 |
|------|------|
| **缓存命中** | 立即返回缓存的描述符 |
| **缓存未命中** | 通过反射 API 获取描述符并缓存 |
| **过期** | 描述符在 TTL 后过期，下次访问时重新获取 |
| **驱逐** | 缓存满时进行 LRU 驱逐 |

### 缓存配置

```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
      # 缓存 TTL（秒），默认: 300
      descriptor_cache_ttl: 300
```

**推荐的 TTL 值:**
- 开发环境: `60`（1分钟）- 频繁更新
- 测试环境: `300`（5分钟）- 平衡
- 生产环境: `1800`（30分钟）- 稳定服务

---

## Triple 协议检测

Pixiu 支持 **Dubbo Triple 协议**，这是 Apache Dubbo 社区开发的 gRPC 兼容协议。

### 什么是 Triple 协议？

Triple 是一个与 gRPC 兼容的协议，具有以下特点：
- 像 gRPC 一样使用 HTTP/2
- 支持 Protobuf 序列化
- 添加 Triple 特定的元数据头（`tri-*`）
- 使 gRPC 服务能与 Dubbo 生态系统协作

### 启用协议检测

```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
      # 启用 Triple 协议检测
      enable_protocol_detection: true
```

**启用后:**
- Pixiu 自动检测 Triple 或 gRPC 协议
- 提取 Triple 特定的元数据（`tri-service-version`、`tri-group` 等）
- 确保与 Dubbo Triple 服务的兼容性

---

## 使用示例

### 示例 1: 简单透传代理

```yaml
static_resources:
  listeners:
    - name: "grpc-proxy"
      protocol_type: "GRPC"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8882
      filter_chains:
        filters:
          - name: dgp.filter.network.grpcconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: "/"
                    route:
                      cluster: "backend-grpc"
              grpc_filters:
                - name: dgp.filter.grpc.proxy
                  config: {}  # 默认透传模式
  clusters:
    - name: "backend-grpc"
      endpoints:
        - socket_address:
            address: backend.example.com
            port: 50051
```

### 示例 2: 使用反射进行基于内容的路由

```yaml
static_resources:
  listeners:
    - name: "grpc-gateway"
      protocol_type: "GRPC"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8882
      filter_chains:
        filters:
          - name: dgp.filter.network.grpcconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: "/myapi/"
                    route:
                      cluster: "api-v1"
              grpc_filters:
                - name: dgp.filter.grpc.proxy
                  config:
                    reflection_mode: "reflection"
                    descriptor_cache_ttl: 600
```

### 示例 3: 迁移场景的混合模式

```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      # 先尝试反射，失败时回退到透传
      reflection_mode: "hybrid"
      reflection_timeout: 3s
      enable_protocol_detection: true
```

---

## 性能考虑

### 模式性能

| 模式 | 延迟影响 | 吞吐量 | 内存使用 |
|------|---------|--------|---------|
| 透传模式 | 最小 | 最高 | 最低 |
| 反射模式 | +10-20% | 高 | 中等（缓存） |
| 混合模式 | +5-10% | 较高 | 中等 |

### 优化建议

1. **使用透传模式** 当不需要消息检查时
2. **启用缓存** 并为反射模式设置合适的 TTL
3. **设置合理的超时** 防止混合模式阻塞
4. **监控缓存命中率** 以调整 TTL 值

---

## 故障排除

### 问题: 反射模式返回 "service not found"

**原因**: 后端服务器未启用 Server Reflection。

**解决方案**: 在服务器上启用反射：
```go
reflection.Register(grpcServer)
```

### 问题: 混合模式回退到透传模式

**原因**: 反射超时或反射服务不可用。

**解决方案**:
1. 检查后端是否启用了反射
2. 增加 `reflection_timeout` 值
3. 检查到反射服务的网络连接

### 问题: 内存使用过高

**原因**: 描述符缓存过大或 TTL 过长。

**解决方案**:
```yaml
# 减少缓存 TTL
descriptor_cache_ttl: 60  # 1分钟而不是5分钟
```

### 问题: Triple 服务不工作

**原因**: 未启用协议检测。

**解决方案**:
```yaml
enable_protocol_detection: true
```

---

## 迁移指南

### 从透传模式迁移到反射模式

**步骤 1**: 在后端服务器上启用反射
**步骤 2**: 更新 Pixiu 配置：
```yaml
# 之前
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config: {}

# 之后
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
```

### 从旧的 ForceCodec 迁移到反射模式

**旧配置（已弃用）:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      force_codec: "passthrough"
```

**新配置（推荐）:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "passthrough"  # 行为相同，名称更清晰
```

---

## 相关资源

- [gRPC Server Reflection Protocol](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) - 官方规范
- [Issue #821](https://github.com/apache/dubbo-go-pixiu/issues/821) - 原始功能请求
- [PR #849](https://github.com/apache/dubbo-go-pixiu/pull/849) - 实现详情
- [示例代码](https://github.com/apache/dubbo-go-pixiu-samples/tree/master/grpc/reflection) - 完整的示例代码

---

## 注意事项

- **过滤器顺序**: gRPC Proxy 过滤器应该是 `grpc_filters` 列表中的**唯一**过滤器
- **服务器反射**: `reflection` 模式必需，`hybrid` 模式可选
- **线程安全**: 过滤器是线程安全的，支持高并发场景
- **向后兼容**: `force_codec` 字段已弃用；请使用 `reflection_mode` 替代
