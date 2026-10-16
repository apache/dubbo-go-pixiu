# Benchmark 结果

测试环境：macOS、Apple Silicon，每个方法 500 次采样。Dubbo、gRPC、Triple 三套测试并行运行。

## Dubbo

| 方法 | Min | Median | Mean | StdDev | Max |
| --- | --- | --- | --- | --- | --- |
| GetUser（经 Pixiu） | 200µs | **400µs** | 600µs | 1.5ms | 14.2ms |
| GetGender（经 Pixiu） | 200µs | **300µs** | 400µs | 500µs | 4.3ms |
| GetUser0（经 Pixiu） | 200µs | 400µs | 400µs | 200µs | 1.6ms |
| GetUsers（经 Pixiu） | 200µs | 400µs | 400µs | 200µs | 1.7ms |
| GetUser（直连） | 100µs | 200µs | 200µs | 100µs | 500µs |

对比 issue #820 中的 Dubbo 数据：GetUser `600µs → 400µs`（约 1.5x），GetGender `500µs → 300µs`（约 1.7x）。

## gRPC

| 方法 | Min | Median | Mean | StdDev | Max |
| --- | --- | --- | --- | --- | --- |
| GetUser（经 Pixiu） | 200µs | **600µs** | 800µs | 1.1ms | 8.3ms |
| GetUsers（经 Pixiu） | 200µs | **600µs** | 700µs | 400µs | 3.6ms |
| GetUserByName（经 Pixiu） | 200µs | 400µs | 500µs | 300µs | 2ms |
| GetUser（直连） | 100µs | 200µs | 200µs | 400µs | 2.8ms |

## Triple

| 方法 | Min | Median | Mean | StdDev | Max |
| --- | --- | --- | --- | --- | --- |
| GetUser（经 Pixiu） | 200µs | **1ms** | 1.5ms | 1.9ms | 10.7ms |
| GetUsers（经 Pixiu） | 300µs | **1.3ms** | 1.4ms | 700µs | 4.1ms |
| SayHello（经 Pixiu） | 300µs | **800µs** | 900µs | 500µs | 7ms |
| GetUser（直连） | 200µs | 500µs | 600µs | 300µs | 2.5ms |

相较 issue #820 中 Triple `SayHello` 的 `12.7ms`，最新结果为 `800µs`，约降低 **15.9x**。汇总数据按同一方法对比：Dubbo/gRPC 使用 `GetUser`，Triple 使用 `SayHello`。

# 运行方法

1. 构建 pixiu 可执行文件

将工作路径修改至 **dubbo-go-pixiu** 根目录，编译 pixiu 可执行文件

```
go build -o tools/benchmark/dist/pixiu cmd/pixiu/pixiu.go
```

最终可执行文件路径为 `tools/benchmark/dist/pixiu`，后续测试会依赖此路径。

2. 运行 zookeeper 服务

```
docker run -d --name zk -p 2181:2181 zookeeper:latest 
```

3. 运行测试代码

将工作路径修改至 `dubbo-go-pixiu/tools/benchmark/test`

```
# 运行所有测试
go test -v ./...

# 运行 dubbo 测试
go test -v dubbo_suite/dubbo_test.go 

# 运行 gRPC 测试
go test -v grpc_suite/grpc_test.go

# 运行 triple 测试
go test -v triple_suite/proto_suite/proto_test.go 
```
