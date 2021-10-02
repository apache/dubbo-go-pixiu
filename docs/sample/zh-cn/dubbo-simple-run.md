# 如何运行 dubbo 简单的测试例子

## 启动 Pixiu

- 打开 samples 路径

```bash
cd /dubbo-go-pixiu/samples/dubbogo/simple
```

- 使用 start.sh 启动案例

使用 start.sh help 来获取帮助信息

```
./start.sh help 

dubbo-go-pixiu start helper
./start.sh action project
hint:
./start.sh prepare body for prepare config file and up docker in body project
./start.sh startPixiu body for start dubbo or http server in body project
./start.sh startServer body for start pixiu in body project
./start.sh startTest body for start unit test in body project
./start.sh clean body for clean

```

使用  `./start.sh [action] [project]` 来启动案例，比如说：

```bash
./start.sh prepare body
./start.sh startPixiu body
./start.sh startServer body
./start.sh startTest body
./start.sh clean body
```

- 接口访问

先到指定测试例子的路径下，比如 `/samples/dubbogo/simple/body`

执行命令：

```bash
./request.sh
```

