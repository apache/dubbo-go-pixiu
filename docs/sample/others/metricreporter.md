# Metric Reporter Filter (dgp.filter.http.metricreporter)

English | [中文](metricreporter_CN.md)

---

## Overview

The `dgp.filter.http.metricreporter` filter provides unified metric reporting for Pixiu gateway. It consolidates the functionality of two previous filters (`dgp.filter.http.metric` and `dgp.filter.http.prometheusmetric`) and supports both **Pull** and **Push** modes with OpenTelemetry integration.

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
                - name: dgp.filter.http.metricreporter
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
  - name: dgp.filter.http.metricreporter
    config:
      mode: "push"
      push_config:
        gateway_url: "http://push-gateway:9091"  # Push Gateway URL
        job_name: "pixiu"                        # Job name
        push_interval: 100                       # Push every 100 requests
        metric_path: "/metrics"                  # Push path
```

---

## Built-in Metrics

### Pull Mode Metrics

| Metric Name | Type | Description |
|-------------|------|-------------|
| `pixiu_request_count` | Counter | Total number of requests |
| `pixiu_request_elapsed` | Counter | Total request elapsed time (ms) |
| `pixiu_request_error_count` | Counter | Total error count |
| `pixiu_request_content_length` | Counter | Request size (bytes) |
| `pixiu_response_content_length` | Counter | Response size (bytes) |
| `pixiu_process_time_millicec` | Histogram | Request processing time distribution (ms) |

**Labels**: `code`, `method`, `url`, `host`

### Push Mode Metrics

| Metric Name | Type | Description |
|-------------|------|-------------|
| `requests_total` | Counter | Total HTTP requests |
| `request_duration_seconds` | Histogram | Request latency |
| `request_size_bytes` | Histogram | Request size |
| `response_size_bytes` | Histogram | Response size |

**Labels**: `code`, `method`, `url`, `host`

---

## Custom Metrics (Extension Feature)

Other filters can record custom metrics that will be reported by MetricReporter.

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
  - name: dgp.filter.http.metricreporter
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
  - name: dgp.filter.http.metricreporter
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

### Replaces dgp.filter.http.metric (Pull)

**Old Configuration:**
```yaml
- name: dgp.filter.http.metric
  config: {}
```

**New Configuration:**
```yaml
- name: dgp.filter.http.metricreporter
  config:
    mode: "pull"
```

### Replaces dgp.filter.http.prometheusmetric (Push)

**Old Configuration:**
```yaml
- name: dgp.filter.http.prometheusmetric
  config:
    metric_collect_rules:
      push_gateway_url: "http://localhost:9091"
      counter_push: true
      push_interval_threshold: 100
      push_job_name: "pixiu"
```

**New Configuration:**
```yaml
- name: dgp.filter.http.metricreporter
  config:
    mode: "push"
    push_config:
      gateway_url: "http://localhost:9091"
      job_name: "pixiu"
      push_interval: 100
      metric_path: "/metrics"
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
- **Custom Metrics**: Pull mode fully supports dynamic metrics; Push mode logs them (can be extended)
- **OpenTelemetry**: Pull mode is consistent with Pixiu Tracing technology stack

