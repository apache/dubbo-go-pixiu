# 从 URI获取参数

> GET 请求 [sample](https://github.com/apache/dubbo-go-pixiu/tree/develop/sample/dubbo/multi)

## 简单示例

### 接口配置

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/student/:name'
    type: restful
    description: student
    methods:
      - httpVerb: GET
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
              mapType: "string"
          applicationName: "StudentService"
          interface: "com.dubbogo.pixiu.StudentService"
          method: "GetStudentByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
              mapType: "string"
            - name: requestBody._all
              mapTo: 1
              mapType: "object"
          applicationName: "StudentService"
          interface: "com.dubbogo.pixiu.StudentService"
          method: "UpdateStudentByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/teacher/:name'
    type: restful
    description: teacher
    methods:
      - httpVerb: GET
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
              mapType: "string"
          applicationName: "TeacherService"
          interface: "com.dubbogo.pixiu.TeacherService"
          method: "GetTeacherByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
              mapType: "string"
            - name: requestBody._all
              mapTo: 1
              mapType: "object"
          applicationName: "TeacherService"
          interface: "com.dubbogo.pixiu.TeacherService"
          method: "UpdateTeacherByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### server config
```yaml
# dubbo server yaml configure file
# application config
application:
  organization: "dubbogopixiu.com"
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
    # Equivalent to the interface in the  dubbo.xml file
    interface: "com.dubbogo.pixiu.StudentService"
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
    # Equivalent to the interface in the  dubbo.xml file
    interface: "com.dubbogo.pixiu.TeacherService"
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

### 测试

- 来自 uri

```bash
curl localhost:8888/api/v1/test-dubbo/student/tc-student -X GET 
and
curl localhost:8888/api/v1/test-dubbo/teacher/tc-teacher -X GET 
```

结果

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

[Previous](dubbo.md)
