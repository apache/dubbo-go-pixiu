# 如何运行 dubbo 简单的测试例子

## 启动 pixiu

- 构建 pixiu

```bash
cd /dubbogo-go-pixiu/cmd/pixiu

go build
```

- 把 pixiu 移动到运行目录

```bash
mv -b pixiu ../../samples/dubbogo/simple
```

- 启动某个例子

通过 `./start.sh [name]` 启动 pixiu，可以选择如下：

```bash
./start.sh body
./start.sh mix
./start.sh query
./start.sh uri
./start.sh pixiu
```

- 接口访问

先到指定测试例子的路径下，比如 `/samples/dubbogo/simple/body`

执行命令：

```bash
./request.sh
```

