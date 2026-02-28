# Benchmark Results

Test environment: macOS, Apple Silicon, N=500 samples per method

## gRPC Protocol

### gRPC Direct

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 100µs | 200µs | 400µs | 1.1ms | 8.2ms |
| GetUsers | 100µs | 200µs | 200µs | 100µs | 400µs |
| GetUserByName | 100µs | 200µs | 200µs | 100µs | 600µs |
| SayHello | 100µs | 200µs | 200µs | 100µs | 400µs |

### gRPC via Pixiu

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 600µs | 1.4ms | 1.8ms | 1.8ms | 12.7ms |
| GetUsers | 500µs | 1.4ms | 1.7ms | 1.2ms | 14.6ms |
| GetUserByName | 500µs | 1.3ms | 2ms | 2ms | 17.1ms |
| SayHello | 600µs | 1.3ms | 1.8ms | 2ms | 32.8ms |

## Triple Protocol

### Triple Direct

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 200µs | 1.8ms | 2.5ms | 2.4ms | 14.5ms |
| GetUsers | 200µs | 1.5ms | 2.0ms | 2.0ms | 12.0ms |
| GetUserByName | 200µs | 1.2ms | 1.8ms | 1.8ms | 10.5ms |
| SayHello | 200µs | 900µs | 1.4ms | 1.3ms | 12.5ms |

### Triple via Pixiu

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 1ms | 2.1ms | 2.8ms | 3.1ms | 24.4ms |
| GetUsers | 800µs | 2.0ms | 2.5ms | 2.5ms | 20.0ms |
| GetUserByName | 700µs | 1.8ms | 2.3ms | 2.0ms | 18.0ms |
| SayHello | 700µs | 2ms | 2.5ms | 1.4ms | 12ms |

## Dubbo Protocol

### Dubbo Direct

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 100µs | 200µs | 600µs | 2.8ms | 20.7ms |
| GetGender | 0s | 100µs | 200µs | 100µs | 600µs |
| GetUser0 | 100µs | 100µs | 200µs | 0s | 400µs |
| GetUsers | 100µs | 200µs | 200µs | 100µs | 600µs |
| GetUser2 | 0s | 200µs | 200µs | 100µs | 500µs |
| GetErr | 300µs | 500µs | 500µs | 200µs | 1.4ms |

### Dubbo via Pixiu

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 700µs | 2ms | 3ms | 5.9ms | 44.7ms |
| GetGender | 300µs | 800µs | 1ms | 500µs | 3.1ms |
| GetUser0 | 300µs | 700µs | 800µs | 500µs | 4.2ms |
| GetUsers | 100µs | 500µs | 600µs | 400µs | 3.2ms |

## Performance Summary

Pixiu proxy adds approximately 0.3-1.1ms overhead compared to direct protocol calls.

| Protocol | Direct (Median) | via Pixiu (Median) | Overhead |
|----------|-----------------|--------------------|-----------|
| gRPC | ~200µs | ~1.3ms | ~1.1ms |
| Triple | ~1.4ms | ~2ms | ~0.6ms |
| Dubbo | ~200µs | ~700µs | ~0.5ms |

# How to Run

1. Build the Pixiu executable

Change the working directory to the root of **dubbo-go-pixiu** and build the Pixiu executable.

```bash
go build -o tools/benchmark/dist/pixiu cmd/pixiu/pixiu.go
```

The final executable will be located at `tools/benchmark/dist/pixiu`, which will be used in subsequent tests.

2. Start the Zookeeper service

```bash
docker run -d --name zk -p 2181:2181 zookeeper:latest
```

3. Run the test code

Change the working directory to `dubbo-go-pixiu/tools/benchmark/test`.

```bash
# Run All Tests
go test -v ./...

# Run Dubbo Tests
go test -v dubbo_suite/dubbo_test.go

# Run gRPC Tests
go test -v grpc_suite/grpc_test.go

# Run Triple Tests
go test -v triple_suite/proto_suite/proto_test.go
```
