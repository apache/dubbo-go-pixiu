# Start

How to start the dubbo-go-pixiu


#### 1 cd samples dir

```
cd samples/dubbo/simple
```

we can use start.sh to run samples quickly. for more info, execute command as below for more help

```
./start.sh [action] [project]
./start.sh help
```

we run body samples below step

#### 2 prepare config file and docker

prepare command will prepare dubbo-server and pixiu config file and start docker container needed

```
./start.sh prepare body
```

if prepare config file manually, notice:
- modify $PROJECT_DIR in conf.yaml to absolute path in your compute

#### 3 start dubbo or http server

```
./start.sh startServer body
```

#### 4 start pixiu

```
./start.sh startPixiu body
```

if run pixiu manually, use command as below

```
 go run cmd/pixiu/*.go gateway start -c /[absolute-path]/dubbo-go-pixiu/samples/dubbo/simple/body/pixiu/conf.yaml
```


#### 5. Try a request

use curl to try or use unit test

```bash
curl -X POST 'localhost:8881/api/v1/test-dubbo/user' -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header 'Content-Type: application/json' 
./start.sh startTest body
```

#### 6. Clean

```
./start.sh clean body
```
