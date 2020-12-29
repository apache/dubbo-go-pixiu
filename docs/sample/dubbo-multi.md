# Get the parameter from the part of uri

> GET request [sample](https://github.com/dubbogo/dubbo-go-proxy/tree/develop/sample/dubbo/multi)

## Simple Demo

### Api Config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/student/:name'
    type: restful
    description: student
    filters:
      - filter0
    methods:
      - httpVerb: GET
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
          applicationName: "StudentService"
          interface: "com.dubbogo.proxy.StudentService"
          method: "GetStudentByName"
          paramTypes: [ "string" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
            - name: requestBody._all
              mapTo: 1
          applicationName: "StudentService"
          interface: "com.dubbogo.proxy.StudentService"
          method: "UpdateStudentByName"
          paramTypes: [ "string", "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/teacher/:name'
    type: restful
    description: teacher
    filters:
      - filter0
    methods:
      - httpVerb: GET
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
          applicationName: "TeacherService"
          interface: "com.dubbogo.proxy.TeacherService"
          method: "GetTeacherByName"
          paramTypes: [ "string" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
            - name: requestBody._all
              mapTo: 1
          applicationName: "TeacherService"
          interface: "com.dubbogo.proxy.TeacherService"
          method: "UpdateTeacherByName"
          paramTypes: [ "string", "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### server config
```yaml
# dubbo server yaml configure file
# application config
application:
  organization: "dubbogoproxy.com"
  name: "BDTService"
  module: "dubbogo user-info server"
  version: "0.0.1"
  owner: "ZX"
  environment: "dev"

registries:
  "demoZk1":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2181"
  "demoZk2":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2182"

services:
  "StudentProvider":
    registry: "demoZk1, demoZk2"
    protocol: "dubbo"
    # 相当于dubbo.xml中的interface
    interface: "com.dubbogo.proxy.StudentService"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    group: test
    version: 1.0.0
    methods:
      - name: "GetStudentByName"
        retries: 1
        loadbalance: "random"
  "TeacherProvider":
    registry: "demoZk1, demoZk2"
    protocol: "dubbo"
    # 相当于dubbo.xml中的interface
    interface: "com.dubbogo.proxy.TeacherService"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    group: test
    version: 1.0.0
    methods:
      - name: "GetTeacherByName"
        retries: 1
        loadbalance: "random"

protocols:
  "dubbo":
    name: "dubbo"
    port: 20000

protocol_conf:
  dubbo:
    session_number: 700
    session_timeout: "20s"
    getty_session_param:
      compress_encoding: false
      tcp_no_delay: true
      tcp_keep_alive: true
      keep_alive_period: "120s"
      tcp_r_buf_size: 262144
      tcp_w_buf_size: 65536
      pkg_rq_size: 1024
      pkg_wq_size: 512
      tcp_read_timeout: "1s"
      tcp_write_timeout: "5s"
      wait_timeout: "1s"
      max_msg_len: 1024
      session_name: "server"
```

### Test

- from uri

```bash
curl localhost:8888/api/v1/test-dubbo/student/tc-student -X GET 
and
curl localhost:8888/api/v1/test-dubbo/teacher/tc-teacher -X GET 
```

result

```bash
{
    "age":18,
    "code":1,
    "iD":"0001",
    "name":"tc-student",
    "time":"2020-12-26T21:50:04.148+08:00"
}
and 
{
    "age":18,
    "code":1,
    "iD":"0001",
    "name":"tc-teacher",
    "time":"2020-12-26T21:50:04.148+08:00"
}
```

[Previous](./dubbo.md)