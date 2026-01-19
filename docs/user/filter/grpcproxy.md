# gRPC Proxy Filter (dgp.filter.grpc.proxy)

English | [中文](grpcproxy_CN.md)

---

## Overview

The `dgp.filter.grpc.proxy` filter provides gRPC proxy functionality for Pixiu gateway with support for **gRPC Server Reflection**. This feature enables dynamic message parsing and inspection at the gateway level without requiring pre-compiled proto files.

### Key Features

- **Three Reflection Modes**: Choose between performance and functionality based on your needs
- **Dynamic Message Decoding**: Parse and inspect gRPC messages at runtime without proto files
- **TTL-based Caching**: Efficient caching of method descriptors with automatic cleanup
- **Protocol Detection**: Support for both gRPC and Triple protocol compatibility
- **Graceful Fallback**: Hybrid mode provides automatic fallback to passthrough

---

## Reflection Modes

### Passthrough Mode (Default)

Performs transparent binary proxying without decoding messages.

**Characteristics:**
- Highest performance (no message parsing overhead)
- No content inspection capabilities
- Pure binary forwarding

**Use Cases:**
- High-throughput scenarios where message inspection is not needed
- Simple routing based on service/method names only
- When backward compatibility is critical

**Configuration:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "passthrough"  # or omit (default)
```

### Reflection Mode

Uses gRPC Server Reflection API to dynamically decode and inspect message contents.

**Characteristics:**
- Full message decoding at runtime
- Enables content-based routing and filtering
- Requires backend server to support reflection
- Slight performance overhead due to reflection calls

**Use Cases:**
- Content-aware routing (route based on message fields)
- Field-level filtering or transformation
- Logging and debugging with full message inspection
- API gateway scenarios requiring message validation

**Configuration:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
      descriptor_cache_ttl: 300  # 5 minutes cache
```

### Hybrid Mode

Tries reflection first, falls back to passthrough on failure.

**Characteristics:**
- Best of both worlds: content inspection when available
- Graceful degradation for services without reflection
- Reflection timeout prevents blocking

**Use Cases:**
- Mixed environments with varying reflection support
- Migration scenarios (gradually enabling reflection)
- Production environments requiring high availability

**Configuration:**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "hybrid"
      reflection_timeout: 5s  # Max time to wait for reflection
```

---

## Configuration

### Complete Configuration Example

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
                    # Reflection mode (default: "passthrough")
                    reflection_mode: "reflection"

                    # Cache TTL for method descriptors (default: 300s)
                    descriptor_cache_ttl: 300

                    # Enable Triple protocol detection (default: false)
                    enable_protocol_detection: true

                    # Reflection timeout for hybrid mode (default: 5s)
                    reflection_timeout: 5s

                    # TLS configuration (optional)
                    enable_tls: false
                    tls_cert_file: ""
                    tls_key_file: ""

                    # Connection settings (optional)
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

### Configuration Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `reflection_mode` | string | `"passthrough"` | Reflection mode: `"passthrough"`, `"reflection"`, or `"hybrid"` |
| `descriptor_cache_ttl` | int | `300` | Cache TTL for method descriptors in seconds |
| `enable_protocol_detection` | bool | `false` | Enable Triple protocol detection |
| `reflection_timeout` | string | `"5s"` | Max time to wait for reflection in hybrid mode |
| `enable_tls` | bool | `false` | Enable TLS for backend connections |
| `tls_cert_file` | string | `""` | Path to TLS certificate file |
| `tls_key_file` | string | `""` | Path to TLS key file |
| `keepalive_time` | string | `"300s"` | Keepalive time for backend connections |
| `keepalive_timeout` | string | `"5s"` | Keepalive timeout |
| `connect_timeout` | string | `"5s"` | Connection timeout |
| `max_concurrent_streams` | uint32 | `0` (unlimited) | Max concurrent streams |

---

## Enabling Server Reflection

To use `reflection` or `hybrid` modes, your backend gRPC server must have Server Reflection enabled.

### Go (gRPC-Go)

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

func main() {
    server := grpc.NewServer()

    // Register your service
    echo.RegisterEchoServiceServer(server, &echoServer{})

    // Enable server reflection (IMPORTANT!)
    reflection.Register(server)

    // Start server
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
            // Enable server reflection
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
    # ... implementation ...

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    echo_pb2_grpc.add_EchoServiceServicer_to_server(EchoService(), server)

    # Enable server reflection
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

## Mode Comparison

| Feature | Passthrough | Reflection | Hybrid |
|---------|-------------|------------|--------|
| **Performance** | ⭐⭐⭐⭐⭐ (Best) | ⭐⭐⭐ (Good) | ⭐⭐⭐⭐ (Better) |
| **Message Inspection** | ❌ No | ✅ Yes | ✅ Yes (when available) |
| **Requires Reflection** | ❌ No | ✅ Yes | ❌ Optional |
| **Fallback** | N/A | ❌ No | ✅ Yes |
| **Use Case** | High-performance proxy | Content-aware routing | Mixed environments |

---

## Descriptor Cache

The reflection mode uses a TTL-based cache to store method descriptors retrieved from the backend server. This reduces the number of reflection requests and improves performance.

### Cache Behavior

| Operation | Description |
|-----------|-------------|
| **Cache Hit** | Return cached descriptor immediately |
| **Cache Miss** | Fetch descriptor via reflection API, cache it |
| **Expiration** | Descriptors expire after TTL, re-fetched on next access |
| **Eviction** | LRU eviction when cache is full |

### Cache Configuration

```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
      # Cache TTL in seconds (default: 300)
      descriptor_cache_ttl: 300
```

**Recommended TTL Values:**
- Development: `60` (1 minute) - frequent updates
- Testing: `300` (5 minutes) - balanced
- Production: `1800` (30 minutes) - stable services

---

## Triple Protocol Detection

Pixiu supports the **Dubbo Triple protocol**, a gRPC-compatible protocol developed by the Apache Dubbo community.

### What is Triple Protocol?

Triple is a gRPC-compatible protocol that:
- Uses HTTP/2 like gRPC
- Supports Protobuf serialization
- Adds Triple-specific metadata headers (`tri-*`)
- Enables gRPC services to work with Dubbo ecosystem

### Enabling Protocol Detection

```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
      # Enable Triple protocol detection
      enable_protocol_detection: true
```

**When enabled:**
- Pixiu automatically detects Triple vs gRPC protocol
- Extracts Triple-specific metadata (`tri-service-version`, `tri-group`, etc.)
- Ensures compatibility with Dubbo Triple services

---

## Usage Examples

### Example 1: Simple Passthrough Proxy

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
                  config: {}  # Default passthrough mode
  clusters:
    - name: "backend-grpc"
      endpoints:
        - socket_address:
            address: backend.example.com
            port: 50051
```

### Example 2: Content-Aware Routing with Reflection

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

### Example 3: Hybrid Mode for Migration

```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      # Try reflection first, fallback to passthrough
      reflection_mode: "hybrid"
      reflection_timeout: 3s
      enable_protocol_detection: true
```

---

## Performance Considerations

### Mode Performance

| Mode | Latency Impact | Throughput | Memory Usage |
|------|----------------|------------|--------------|
| Passthrough | Minimal | Highest | Lowest |
| Reflection | +10-20% | High | Moderate (cache) |
| Hybrid | +5-10% | Higher | Moderate |

### Optimization Tips

1. **Use Passthrough** when message inspection is not needed
2. **Enable Caching** with appropriate TTL for reflection mode
3. **Set Reasonable Timeout** for hybrid mode to prevent blocking
4. **Monitor Cache Hit Ratio** to tune TTL values

---

## Troubleshooting

### Problem: Reflection mode returns "service not found"

**Cause**: Backend server does not have Server Reflection enabled.

**Solution**: Enable reflection on your server:
```go
reflection.Register(grpcServer)
```

### Problem: Hybrid mode falls back to passthrough

**Cause**: Reflection timeout exceeded or reflection service unavailable.

**Solution**:
1. Check if reflection is enabled on the backend
2. Increase `reflection_timeout` value
3. Check network connectivity to reflection service

### Problem: High memory usage

**Cause**: Descriptor cache too large or TTL too long.

**Solution**:
```yaml
# Reduce cache TTL
descriptor_cache_ttl: 60  # 1 minute instead of 5 minutes
```

### Problem: Triple services not working

**Cause**: Protocol detection not enabled.

**Solution**:
```yaml
enable_protocol_detection: true
```

---

## Migration Guide

### From Passthrough to Reflection

**Step 1**: Enable reflection on backend servers
**Step 2**: Update Pixiu configuration:
```yaml
# Before
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config: {}

# After
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "reflection"
```

### From Legacy ForceCodec to Reflection

**Old Configuration (Deprecated):**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      force_codec: "passthrough"
```

**New Configuration (Recommended):**
```yaml
grpc_filters:
  - name: dgp.filter.grpc.proxy
    config:
      reflection_mode: "passthrough"  # Same behavior, clearer name
```

---

## Related Resources

- [gRPC Server Reflection Protocol](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) - Official specification
- [Issue #821](https://github.com/apache/dubbo-go-pixiu/issues/821) - Original feature request
- [PR #849](https://github.com/apache/dubbo-go-pixiu/pull/849) - Implementation details
- [Sample Code](https://github.com/apache/dubbo-go-pixiu-samples/tree/master/grpc/reflection) - Complete working examples

---

## Notes

- **Filter Order**: gRPC Proxy filter should be the **only** filter in `grpc_filters` list
- **Server Reflection**: Required for `reflection` mode, optional for `hybrid` mode
- **Thread Safety**: Filter is thread-safe and supports high-concurrency scenarios
- **Backward Compatibility**: `force_codec` field is deprecated; use `reflection_mode` instead
