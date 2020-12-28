# Get the parameter from the part of uri,query

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

### Test

- from uri and query

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