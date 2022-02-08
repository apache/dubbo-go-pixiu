
# Pixiu - SpringCloud

## Pixiu SpringCloud Discovery With Nacos

there are two spring cloud java project in serve dir which offer same restful api:
- user-service
- auth-service

you can use those to mock below situation
- two service with each one instance
- one service with two instances, just change spring.application.name in  auth-service/src/resource/application.properties from auth-service to user-service notice ! change the port if registry center is nacos
- start a new instance or kill the old to mock instance up or down

## Pixiu SpringCloud practice : Service Discovery With Zookeeper

### Start SpringCloud 

> p.s. the SpringCloud Server name is `pixiu-springcloud-server` , setting with [application.yml](src/main/resources/application.yml)

first, start Zookeeper use for Register Center, maybe you can use [docker-compose](docker-compose.yml)

```shell
services:
  zookeeper:
    image: zookeeper
    ports:
      - "2181:2181"
```

second, start the SpringCloud Server, use [start.sh](start.sh) or run with IDEA 

```shell
mvn spring-boot:run
```

finally, check Server
```shell
curl http://localhost:9127/hi
```
result
```shell
Hello Pixiu World! from org.springframework.cloud.zookeeper.serviceregistry.ServiceInstanceRegistration@63a45760
```
or

find service instance info on zookeeper:
```shell
curl http://localhost:9127/service/pixiu-springcloud-server
```

result
```shell
[
  {
    "serviceId": "pixiu-springcloud-server",
    "host": "127.0.0.1",
    "port": 9127,
    "secure": false,
    "uri": "http://127.0.0.1:9127",
    "metadata": {
      "instance_status": "UP"
    },
    "serviceInstance": {
      "name": "pixiu-springcloud-server",
      "id": "1a7b120a-e77a-4a8d-ad3b-ca75473288a4",
      "address": "127.0.0.1",
      "port": 9127,
      "sslPort": null,
      "payload": {
        "@class": "org.springframework.cloud.zookeeper.discovery.ZookeeperInstance",
        "id": "pixiu-springcloud-server",
        "name": "pixiu-springcloud-server",
        "metadata": {
          "instance_status": "UP"
        }
      },
      "registrationTimeUTC": 1640852269145,
      "serviceType": "DYNAMIC",
      "uriSpec": {
        "parts": [
          {
            "value": "scheme",
            "variable": true
          },
          {
            "value": "://",
            "variable": false
          },
          {
            "value": "address",
            "variable": true
          },
          {
            "value": ":",
            "variable": false
          },
          {
            "value": "port",
            "variable": true
          }
        ]
      },
      "enabled": true
    },
    "instanceId": "1a7b120a-e77a-4a8d-ad3b-ca75473288a4",
    "scheme": null
  }
]
```


### Start Pixiu

start Pixiu 
```shell
gateway start -c samples/springcloud/pixiu/conf.yaml
```

practice  proxy SpringCloud service base on Pixiu
```shell
curl http://127.0.0.1:8888/pixiu-springcloud-server/hi
```
