
# Quick Start

Start SpringCloud :

```bash
cd samples/springcloud/server/spring-cloud && bash start.sh
```

Start Pixiu : 

```bash
go run cmd/pixiu/*.go gateway start -c  ${absolute-path}/dubbo-go-pixiu/samples/springcloud/pixiu/conf.yaml
```

Call the server of SpringCloud by Pixiu :

```bash

# the serviceId is `user-provider`
curl http://localhost:8888/user-service/echo/Pixiu

# the serviceId is `auth-provider`
curl http://localhost:8888/auth-service/echo/Pixiu
```
result on console  : 
```log
Hello Nacos Discovery Pixiu
```
