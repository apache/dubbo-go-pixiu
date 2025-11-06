# 指标上报过滤器 (dgp.filter.http.metric)

[English](metric.md) | 中文

---

## 概述

`dgp.filter.http.metric` 过滤器为 Pixiu 网关提供统一的指标上报功能。它整合了之前的两个过滤器（`dgp.filter.http.metric` 和 `dgp.filter.http.prometheusmetric`）的功能，并支持 **Pull** 和 **Push** 两种模式，集成了 OpenTelemetry。

> **注意**：此过滤器默认使用 **Push** 模式。如需使用 Pull 模式，请在配置中显式指定 `mode: "pull"`。

### 核心特性

- **统一入口**：单个过滤器支持 Pull 和 Push 两种模式
- **OpenTelemetry 集成**：Pull 模式使用 OpenTelemetry（与 Pixiu Tracing 保持一致）
- **Context 扩展**：其他过滤器可通过 `HttpContext.RecordMetric()` 记录自定义指标
- **向后兼容**：复用原有过滤器的逻辑

---

## 模式说明

### Pull 模式（推荐）

指标通过 HTTP 端点暴露，供 Prometheus 抓取。

**特点：**
- 使用 OpenTelemetry SDK
- 指标通过全局 HTTP 端点暴露
- Prometheus 标准拉取模型
- 支持来自 Context 的动态自定义指标

**适用场景：**
- 长期运行的服务
- Kubernetes 环境
- 开发和测试环境

### Push 模式

指标主动推送到 Prometheus Push Gateway。

**特点：**
- 使用 Prometheus 原生 SDK
- 每 N 个请求推送一次指标
- 批量推送减少网络开销
- 支持来自 Context 的动态自定义指标（与 Pull 模式相同）

**适用场景：**
- 防火墙后的服务
- 短生命周期批处理任务
- 无法暴露入站端口的环境

---

## 配置说明

### Pull 模式

```yaml
static_resources:
  listeners:
    - name: "net/http"
      protocol_type: "HTTP"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8888
      filter_chains:
        filters:
          - name: dgp.filter.httpconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: /
                    route:
                      cluster: backend
              http_filters:
                # 其他过滤器...
                - name: dgp.filter.http.httpproxy
                  config: {}
                
                # MetricReporter 必须放在最后
                - name: dgp.filter.http.metric
                  config:
                    mode: "pull"

# 全局指标配置（控制 HTTP 端点）
metric:
  enable: true
  prometheus_port: 2222  # 访问地址：http://localhost:2222/metrics
```

### Push 模式

```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "push"
      push_config:
        gateway_url: "http://push-gateway:9091"  # Push Gateway 地址（默认：http://localhost:9091）
        job_name: "pixiu"                        # 任务名称（默认：pixiu）
        push_interval: 100                       # 每 100 个请求推送一次（默认：100）
        metric_path: "/metrics"                  # 推送路径（默认：/metrics）
```

**注意**：`push_config` 中的所有字段都有默认值。如果省略，将自动应用默认值。

**最小化配置（使用所有默认值）**：
```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "push"
      # push_config 可以省略或为空，将使用所有默认值
```

**极简配置（默认使用 Push 模式）**：
```yaml
http_filters:
  - name: dgp.filter.http.metric
    config: {}
    # 默认使用 push 模式和所有默认配置
```

**默认值**：
- `gateway_url`: `http://localhost:9091`
- `job_name`: `pixiu`
- `push_interval`: `100`
- `metric_path`: `/metrics`

---

## 内置指标

**Pull 和 Push 模式现在使用统一的指标名称和类型：**

| 指标名称 | 类型 | 描述 | 标签 | 状态 |
|---------|------|------|------|------|
| `pixiu_requests_total` | Counter | 请求总数 | code, method, host, url | ⚠️ 已弃用 |
| `pixiu_request_count` | Counter | 请求总数 | code, method, host, url | ✅ 推荐 |
| `pixiu_request_elapsed` | Counter | 请求总耗时（毫秒）| code, method, host, url | ✅ |
| `pixiu_request_error_count` | Counter | 错误总数 | code, method, host, url | ✅ |
| `pixiu_request_content_length` | Counter | 请求大小（字节）| code, method, url | ✅ |
| `pixiu_response_content_length` | Counter | 响应大小（字节）| code, method, url | ✅ |
| `pixiu_process_time_millisec` | Histogram | 请求处理时长分布（毫秒）| code, method, url | ✅ |

**向后兼容说明**：
- 为保持向后兼容，目前**同时导出**两个指标：`pixiu_requests_total`（旧）和 `pixiu_request_count`（新）
- 推荐使用新的指标名称 `pixiu_request_count`
- `pixiu_requests_total` 将在未来版本中移除
- 两种模式使用相同的指标名称，便于在不同模式间切换

---

## 自定义指标（扩展功能）

其他过滤器可以记录自定义指标，由 MetricReporter 自动收集并上报。

**Pull 和 Push 模式都支持**：通过 `HttpContext.RecordMetric()` 记录的自定义指标在两种模式下都能完整导出。

### 在自定义过滤器中使用

```go
package myfilter

import (
    "fmt"
    "time"
    "github.com/apache/dubbo-go-pixiu/pkg/context/http"
    "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
)

type Filter struct {
    startTime time.Time
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
    f.startTime = time.Now()
    return filter.Continue
}

func (f *Filter) Encode(ctx *http.HttpContext) filter.FilterStatus {
    // 记录 Counter 指标
    ctx.RecordMetric("my_requests_total", "counter", 1.0, map[string]string{
        "method": ctx.GetMethod(),
        "status": fmt.Sprintf("%d", ctx.GetStatusCode()),
    })
    
    // 记录 Histogram 指标
    latency := time.Since(f.startTime).Milliseconds()
    ctx.RecordMetric("my_request_duration_ms", "histogram", float64(latency), map[string]string{
        "endpoint": ctx.GetUrl(),
    })
    
    // 记录 Gauge 指标
    ctx.RecordMetric("my_active_connections", "gauge", float64(42), nil)
    
    return filter.Continue
}
```

### 支持的指标类型

| 类型 | 描述 | 示例 |
|------|------|------|
| `counter` | 单调递增的值 | 请求数、错误数 |
| `histogram` | 值的分布统计 | 延迟、请求大小 |
| `gauge` | 可增可减的值 | 活跃连接数、内存使用 |

---

## 使用指南

### Pull 模式配置

**步骤 1：配置全局指标端点**

```yaml
# conf.yaml
metric:
  enable: true
  prometheus_port: 2222
```

**步骤 2：添加 MetricReporter 过滤器**

```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "pull"
```

**步骤 3：配置 Prometheus**

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'pixiu'
    static_configs:
      - targets: ['localhost:2222']
    scrape_interval: 15s
```

**步骤 4：访问指标**

```bash
curl http://localhost:2222/metrics
```

### Push 模式配置

**步骤 1：启动 Push Gateway**

```bash
docker run -d -p 9091:9091 prom/pushgateway
```

**步骤 2：配置 MetricReporter**

```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "push"
      push_config:
        gateway_url: "http://localhost:9091"
        job_name: "pixiu"
        push_interval: 100
        metric_path: "/metrics"
```

**步骤 3：验证指标**

```bash
# 查看 Push Gateway 中的指标
curl http://localhost:9091/metrics
```

---

## 与旧版过滤器的区别

### 替代旧版 dgp.filter.http.prometheusmetric (Push)

> **重要说明**：`dgp.filter.http.metric` 过滤器现在统一支持 Pull 和 Push 两种模式，默认为 Push 模式。旧版的 `dgp.filter.http.prometheusmetric` 过滤器已被标记为废弃。

**旧配置（已废弃）：**
```yaml
- name: dgp.filter.http.prometheusmetric
  config:
    metric_collect_rules:
      push_gateway_url: "http://localhost:9091"
      counter_push: true
      push_interval_threshold: 100
      push_job_name: "pixiu"
```

**新配置（推荐）：**
```yaml
- name: dgp.filter.http.metric
  config:
    mode: "push"
    push_config:
      gateway_url: "http://localhost:9091"
      job_name: "pixiu"
      push_interval: 100
      metric_path: "/metrics"
```

**使用默认配置（更简单）：**
```yaml
- name: dgp.filter.http.metric
  config: {}
  # 默认使用 push 模式，gateway_url=http://localhost:9091, job_name=pixiu, push_interval=100
```

---

## 故障排查

### Pull 模式

**问题**：无法访问指标端点

**解决**：确保全局指标配置已启用：
```yaml
metric:
  enable: true
  prometheus_port: 2222
```

**问题**：指标未显示

**解决**：检查过滤器是否放在 http_filters 列表的最后

### Push 模式

**问题**：指标未推送到 Gateway

**解决**：
1. 验证 Push Gateway 正在运行：`curl http://push-gateway:9091`
2. 检查推送间隔阈值是否达到（发送 N 个请求）
3. 查看日志中的推送错误

**问题**：重复指标警告

**解决**：这在运行多个测试时是正常的；指标会被正确聚合

---

## 注意事项

- **过滤器顺序**：MetricReporter 应放在 http_filters 列表的**最后**
- **Pull 端点**：由全局 `metric.prometheus_port` 控制，而非过滤器配置
- **指标统一**：Pull 和 Push 模式使用相同的指标名称和类型
- **配置默认值**：Push 模式的所有配置字段都可以省略，会自动应用默认值
- **线程安全**：过滤器是线程安全的，支持高并发场景

