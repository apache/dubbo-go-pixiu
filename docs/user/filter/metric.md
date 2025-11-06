# Metric Reporter Filter (dgp.filter.http.metric)

English | [中文](metric_CN.md)

---

## Overview

The `dgp.filter.http.metric` filter provides unified metric reporting for Pixiu gateway. It consolidates the functionality of two previous filters (`dgp.filter.http.metric` and `dgp.filter.http.prometheusmetric`) and supports both **Pull** and **Push** modes with OpenTelemetry integration.

> **Note**: This filter defaults to **Push** mode. To use Pull mode, explicitly specify `mode: "pull"` in the configuration.

### Key Features

- **Unified Entry Point**: Single filter for both Pull and Push modes
- **OpenTelemetry Integration**: Pull mode uses OpenTelemetry for metrics (consistent with Pixiu Tracing)
- **Context-based Extension**: Other filters can record custom metrics via `HttpContext.RecordMetric()`
- **Backward Compatible**: Reuses logic from original filters

---

## Modes

### Pull Mode (Recommended)

Metrics are exposed via HTTP endpoint for Prometheus to scrape.

**Characteristics:**
- Uses OpenTelemetry SDK
- Metrics exposed via global HTTP endpoint
- Standard Prometheus pull model
- Supports dynamic custom metrics from context

**Use Cases:**
- Long-running services
- Kubernetes environments
- Development and testing

### Push Mode

Metrics are actively pushed to Prometheus Push Gateway.

**Characteristics:**
- Uses Prometheus native SDK
- Metrics pushed every N requests
- Batch push reduces network overhead
- Supports dynamic custom metrics from context (same as Pull mode)

**Use Cases:**
- Services behind firewalls
- Short-lived batch jobs
- Cannot expose inbound ports

---

## Configuration

### Pull Mode

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
                # Other filters...
                - name: dgp.filter.http.httpproxy
                  config: {}
                
                # MetricReporter must be placed last
                - name: dgp.filter.http.metric
                  config:
                    mode: "pull"

# Global metric configuration (controls the HTTP endpoint)
metric:
  enable: true
  prometheus_port: 2222  # Access metrics at http://localhost:2222/metrics
```

### Push Mode

```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "push"
      push_config:
        gateway_url: "http://push-gateway:9091"  # Push Gateway URL (default: http://localhost:9091)
        job_name: "pixiu"                        # Job name (default: pixiu)
        push_interval: 100                       # Push every 100 requests (default: 100)
        metric_path: "/metrics"                  # Push path (default: /metrics)
```

**Note**: All fields in `push_config` have default values. If omitted, defaults will be applied automatically.

**Minimal Configuration (uses all defaults)**:
```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "push"
      # push_config can be omitted or empty to use all defaults
```

**Ultra-minimal Configuration (defaults to Push mode)**:
```yaml
http_filters:
  - name: dgp.filter.http.metric
    config: {}
    # Defaults to push mode with all default settings
```

**Default Values**:
- `gateway_url`: `http://localhost:9091`
- `job_name`: `pixiu`
- `push_interval`: `100`
- `metric_path`: `/metrics`

---

## Built-in Metrics

**Pull and Push modes now use unified metric names and types:**

| Metric Name | Type | Description | Labels | Status |
|-------------|------|-------------|--------|--------|
| `pixiu_requests_total` | Counter | Total number of requests | code, method, host, url | ⚠️ Deprecated |
| `pixiu_request_count` | Counter | Total number of requests | code, method, host, url | ✅ Recommended |
| `pixiu_request_elapsed` | Counter | Total request elapsed time (milliseconds) | code, method, host, url | ✅ |
| `pixiu_request_error_count` | Counter | Total error count | code, method, host, url | ✅ |
| `pixiu_request_content_length` | Counter | Request size (bytes) | code, method, url | ✅ |
| `pixiu_response_content_length` | Counter | Response size (bytes) | code, method, url | ✅ |
| `pixiu_process_time_millisec` | Histogram | Request processing time distribution (milliseconds) | code, method, url | ✅ |

**Backward Compatibility**:
- For backward compatibility, both metrics are currently exported: `pixiu_requests_total` (old) and `pixiu_request_count` (new)
- Use the new metric name `pixiu_request_count` for new deployments
- `pixiu_requests_total` will be removed in future versions
- Both modes use identical metric names for easy switching

---

## Custom Metrics (Extension Feature)

Other filters can record custom metrics that will be automatically collected and reported by MetricReporter.

**Supported in both Pull and Push modes**: Custom metrics recorded via `HttpContext.RecordMetric()` are fully exported in both modes.

### Usage in Custom Filters

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
    // Record counter metric
    ctx.RecordMetric("my_requests_total", "counter", 1.0, map[string]string{
        "method": ctx.GetMethod(),
        "status": fmt.Sprintf("%d", ctx.GetStatusCode()),
    })
    
    // Record histogram metric
    latency := time.Since(f.startTime).Milliseconds()
    ctx.RecordMetric("my_request_duration_ms", "histogram", float64(latency), map[string]string{
        "endpoint": ctx.GetUrl(),
    })
    
    // Record gauge metric
    ctx.RecordMetric("my_active_connections", "gauge", float64(42), nil)
    
    return filter.Continue
}
```

### Supported Metric Types

| Type | Description | Example |
|------|-------------|---------|
| `counter` | Monotonically increasing value | Request count, error count |
| `histogram` | Value distribution | Latency, request size |
| `gauge` | Value that can go up or down | Active connections, memory usage |

---

## Setup Guide

### Pull Mode Setup

**Step 1: Configure Global Metric Endpoint**

```yaml
# conf.yaml
metric:
  enable: true
  prometheus_port: 2222
```

**Step 2: Add MetricReporter Filter**

```yaml
http_filters:
  - name: dgp.filter.http.metric
    config:
      mode: "pull"
```

**Step 3: Configure Prometheus**

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'pixiu'
    static_configs:
      - targets: ['localhost:2222']
    scrape_interval: 15s
```

**Step 4: Access Metrics**

```bash
curl http://localhost:2222/metrics
```

### Push Mode Setup

**Step 1: Start Push Gateway**

```bash
docker run -d -p 9091:9091 prom/pushgateway
```

**Step 2: Configure MetricReporter**

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

**Step 3: Verify Metrics**

```bash
# View metrics in Push Gateway
curl http://localhost:9091/metrics
```

---

## Differences from Previous Filters

### Replaces Legacy dgp.filter.http.prometheusmetric (Push)

> **Important**: The `dgp.filter.http.metric` filter now supports both Pull and Push modes, with Push as the default. The legacy `dgp.filter.http.prometheusmetric` filter has been marked as deprecated.

**Old Configuration (Deprecated):**
```yaml
- name: dgp.filter.http.prometheusmetric
  config:
    metric_collect_rules:
      push_gateway_url: "http://localhost:9091"
      counter_push: true
      push_interval_threshold: 100
      push_job_name: "pixiu"
```

**New Configuration (Recommended):**
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

**With Default Configuration (Simpler):**
```yaml
- name: dgp.filter.http.metric
  config: {}
  # Defaults to push mode with gateway_url=http://localhost:9091, job_name=pixiu, push_interval=100
```

---

## Troubleshooting

### Pull Mode

**Problem**: Cannot access metrics endpoint

**Solution**: Ensure global metric configuration is enabled:
```yaml
metric:
  enable: true
  prometheus_port: 2222
```

**Problem**: Metrics not showing up

**Solution**: Check if filter is placed at the end of http_filters list

### Push Mode

**Problem**: Metrics not pushed to Gateway

**Solution**: 
1. Verify Push Gateway is running: `curl http://push-gateway:9091`
2. Check push interval threshold is met (send N requests)
3. Review logs for push errors

**Problem**: Duplicate metrics warnings

**Solution**: Normal when running multiple tests; metrics will be properly aggregated

---

## Notes

- **Filter Order**: MetricReporter should be placed **last** in the http_filters list
- **Pull Endpoint**: Controlled by global `metric.prometheus_port`, not filter config
- **Unified Metrics**: Pull and Push modes use identical metric names and types
- **Configuration Defaults**: All Push mode configuration fields can be omitted with automatic defaults
- **Thread Safety**: Filter is thread-safe and supports high-concurrency scenarios

