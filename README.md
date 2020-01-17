[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-proxy.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-proxy)

### instructions

The gateway can support calling Java Dubbo cluster and golang Dubbo cluster

You can use 'go build' and './start.sh' to start gateway.
then,you can do http request to you dubbo interface.


#### HTTP request format when do not have metadata center:

Group and version is the mapping data in Dubbo service. 

```
{application Name}/​{Interface name}?version={version}&group={group}&method={method}
```

http POST body: 

```json
{
    "paramTypes" : ["org.apache.dubbo.demo.model.User"],
    "paramValues": [
        {
            "id": 23,
            "username": "testUser"
        }
    ]
}
```

#### HTTP request format when  have metadata center:
```
{application Name}/​{Interface name}?version={version}&group={group}&method={method}
```
```json
{
    "paramValues": [
        {
            "id": 23,
            "username": "testUser"
        }
    ]
}

```






### 服务端配置元数据中心：

添加依赖
```
<dependency>
    <groupId>org.apache.dubbo</groupId>
    <artifactId>dubbo-metadata-report-redis</artifactId>
</dependency>
```

在provider和consumer应用中配置MetadataReportConfig
```
@Bean
public MetadataReportConfig metadataReportConfig() {
    MetadataReportConfig metadataReportConfig = new MetadataReportConfig();
    metadataReportConfig.setAddress("redis://localhost:6379");
    return metadataReportConfig;
}

```

 启动provider，可以在redis中发现多了一个key

```
org.apache.dubbo.demo.DemoService:provider:dubbo-demo-annotation-provider.metaData
{
    "parameters": {
        "side": "provider",
        "release": "",
        "methods": "sayHello",
        "deprecated": "false",
        "dubbo": "2.0.2",
        "default.dynamic": "false",
        "interface": "org.apache.dubbo.demo.DemoService",
        "generic": "false",
        "default.deprecated": "false",
        "application": "dubbo-demo-annotation-provider",
        "default.register": "true",
        "dynamic": "false",
        "bean.name": "providers:dubbo:org.apache.dubbo.demo.DemoService",
        "register": "true",
        "anyhost": "true"
    },
    "canonicalName": "org.apache.dubbo.demo.DemoService",
    "codeSource": "file:/D:/Workspace/Dubbo/dubbo-demo/dubbo-demo-interface/target/classes/",
    "methods": [
        {
            "name": "sayHello",
            "parameterTypes": [
                "java.lang.String"
            ],
            "returnType": "java.lang.String"
        }
    ],
    "types": [
        {
            "type": "byte"
        },
        {
            "type": "java.lang.String",
            "properties": {
                "coder": {
                    "type": "byte"
                },
                "value": {
                    "type": "byte[]"
                },
                "hash": {
                    "type": "int"
                }
            }
        },
        {
            "type": "int"
        }
    ]
}
```

