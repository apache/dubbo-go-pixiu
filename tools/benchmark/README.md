# Benchmark Results

Test environment: macOS, Apple Silicon, N=500 samples per method. The three suites were run in parallel with `go test -count=1 ./...`.

## gRPC Protocol

### gRPC Direct

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 100µs | 200µs | 200µs | 400µs | 2.8ms |
| GetUsers | 100µs | 200µs | 200µs | 100µs | 400µs |
| GetUserByName | 100µs | 200µs | 200µs | 100µs | 400µs |
| SayHello | 100µs | 200µs | 200µs | 0s | 500µs |

### gRPC via Pixiu

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 200µs | 600µs | 800µs | 1.1ms | 8.3ms |
| GetUsers | 200µs | 600µs | 700µs | 400µs | 3.6ms |
| GetUserByName | 200µs | 400µs | 500µs | 300µs | 2ms |
| SayHello | 200µs | 500µs | 600µs | 300µs | 1.8ms |

## Triple Protocol

### Triple Direct

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 200µs | 500µs | 600µs | 300µs | 2.5ms |
| GetUsers | 200µs | 500µs | 500µs | 200µs | 1.7ms |
| GetUserByName | 200µs | 600µs | 800µs | 700µs | 7ms |
| SayHello | 200µs | 500µs | 500µs | 200µs | 1.2ms |

### Triple via Pixiu

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 200µs | 1ms | 1.5ms | 1.9ms | 10.7ms |
| GetUsers | 300µs | 1.3ms | 1.4ms | 700µs | 4.1ms |
| GetUserByName | 300µs | 900µs | 1.2ms | 1ms | 7ms |
| SayHello | 300µs | 800µs | 900µs | 500µs | 7ms |

## Dubbo Protocol

### Dubbo Direct

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 100µs | 200µs | 200µs | 100µs | 500µs |
| GetGender | 100µs | 200µs | 200µs | 100µs | 500µs |
| GetUser0 | 100µs | 200µs | 200µs | 100µs | 900µs |
| GetUsers | 100µs | 200µs | 200µs | 100µs | 500µs |

### Dubbo via Pixiu

| Method | Min | Median | Mean | StdDev | Max |
|--------|-----|--------|------|--------|-----|
| GetUser | 200µs | 400µs | 600µs | 1.5ms | 14.2ms |
| GetGender | 200µs | 300µs | 400µs | 500µs | 4.3ms |
| GetUser0 | 200µs | 400µs | 400µs | 200µs | 1.6ms |
| GetUsers | 200µs | 400µs | 400µs | 200µs | 1.7ms |

## Performance Summary

Pixiu proxy adds approximately 0.2-0.4ms median overhead in this run. Results vary with local load and scheduling. The summary compares the same method (`GetUser` for Dubbo and gRPC, `SayHello` for Triple).

| Protocol | Direct (Median) | via Pixiu (Median) | Overhead |
|----------|-----------------|--------------------|-----------|
| gRPC | ~200µs | ~600µs | ~0.4ms |
| Triple | ~500µs | ~800µs | ~0.3ms |
| Dubbo | ~200µs | ~400µs | ~0.2ms |

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
