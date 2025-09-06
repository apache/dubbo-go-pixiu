# How run dubbo simple test samples

## Start Pixiu

- cd sample dir

```bash
cd /dubbo-go-pixiu/samples/dubbogo/simple
```

- use start.sh to start pixiu and server quickly

start.sh help for help info

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

Use `./start.sh [action] [project]` to start sample，for exmaple：

```bash
./start.sh prepare body
./start.sh startPixiu body
./start.sh startServer body
./start.sh startTest body
./start.sh clean body
```

- interface access

Go to the path of the specified test example，such as `/samples/dubbogo/simple/body`

Executive Command：

```bash
./request.sh
```

