# 本地例子

> 自行准备好 zookeeper 环境

## 启动服务提供方

配置好下面两个环境变量：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/sample/server/config/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/sample/server/config/log.yml

运行 sample/server/app/server.go

## 启动 proxy

配置好下面的环境变量：
- -c /Users/tc/Documents/workspace_2020/dubbo-go-proxy/sample/proxy/conf.yaml 
- -a /Users/tc/Documents/workspace_2020/dubbo-go-proxy/sample/proxy/api_config.yaml

运行 cmd/proxy/proxy.go

因为没有配置 dubbogo 的配置文件，所以启动的时候出现日志如下：
```bash
2020/11/06 10:29:09 [InitLog] warn: ioutil.ReadFile(file:./conf/log.yml) = error:open ./conf/log.yml: no such file or directory
2020/11/06 10:29:09 [InitLog] warn: log configure file name is nil
2020/11/06 10:29:28 [consumerInit] application configure(consumer) file name is nil
2020/11/06 10:29:28 [providerInit] application configure(provider) file name is nil
```

