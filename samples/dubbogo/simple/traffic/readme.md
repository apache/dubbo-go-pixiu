#Traffic Filter quick start 
### Start Zookeeper:
```shell
cd samples/dubbogo/simple/traffic
run docker-compose.yml/services
```
### Start Http Server
```shell
cd server
go run server.go
```
```shell
cd server/v1
go run server.go
```
```shell
cd server/v2
go run server.go
```

### Start Pixiu
```shell
go run cmd/pixiu/*.go gateway start -c samples/dubbogo/simple/traffic/pixiu/conf.yaml
```
### Start test
```shell
curl http://localhost:8888/user
curl -H "canary-by-header: 10" http://localhost:8888/user
curl -H "canary-weight: 90" http://localhost:8888/user
```
