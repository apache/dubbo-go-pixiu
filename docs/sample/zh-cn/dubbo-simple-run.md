# 如何运行 dubbo 简单的测试例子

## 启动 proxy

- 构建 proxy

```bash
cd /dubbogo-go-proxy/cmd/proxy

go build
```

- 把 proxy 移动到运行目录

```bash
mv -b proxy ../../samples/dubbogo/simple
```

- 启动某个例子

通过 `./start.sh [name]` 启动 proxy，可以选择如下：

```bash
./start.sh body
./start.sh mix
./start.sh query
./start.sh uri
./start.sh proxy
```

- 接口访问

先到指定测试例子的路径下，比如 `/samples/dubbogo/simple/body`

执行命令：

```bash
./request.sh
```

