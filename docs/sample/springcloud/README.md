# Quick Start

Start Nacos [Docker environment]:

```bash
cd samples/springcloud/docker
run docker-compose.yml/services
```

Start SpringCloud [Java environment]:

```bash
cd samples/springcloud/server

# the port is 8074
run auth-service

# the port is 8071
run user-service
```

Start Pixiu:

```bash
go run cmd/pixiu/*.go gateway start -c samples/springcloud/pixiu/conf.yaml
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
