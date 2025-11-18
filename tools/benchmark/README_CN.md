# Benchmark 结果

* grpc:

```
grpc protocol performance test
      Name                     | N   | Min   | Median | Mean  | StdDev | Max  
      ========================================================================
      GetUser [duration]       | 500 | 100µs | 300µs  | 400µs | 600µs  | 5.2ms
      ------------------------------------------------------------------------
      GetUsers [duration]      | 500 | 100µs | 300µs  | 300µs | 200µs  | 1.5ms
      ------------------------------------------------------------------------
      GetUserByName [duration] | 496 | 100µs | 200µs  | 300µs | 100µs  | 1ms  

pixiu to grpc protocol performance test
      Name                     | N   | Min   | Median | Mean  | StdDev | Max   
      =========================================================================
      GetUser [duration]       | 500 | 600µs | 2ms    | 3.4ms | 4.3ms  | 25.9ms
      -------------------------------------------------------------------------
      GetUsers [duration]      | 500 | 600µs | 1.3ms  | 2.1ms | 3.3ms  | 25.1ms
      -------------------------------------------------------------------------
      GetUserByName [duration] | 500 | 800µs | 1.9ms  | 2.7ms | 3.8ms  | 33.4ms
```

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