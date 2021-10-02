<a name="qorlm"></a>
# Introduction
This document describes how to run and debug Pixiu samples  on Windows. Except for testing, there is no need to enter a command line. <br />The run of each sample is divided into four steps :
> 1. Start the zookeeper
> 1. Start dubbo-go provider
> 1. Start pixiu
> 1. Test

The following will introduce how to start `body sample` under samples/dubbogo/simple.
<a name="b10EM"></a>
# Prerequisite

1. Use an IDE, such as Goland
1. Install the zookeeper
1. Install the curl for windows
   <a name="gq18V"></a>
# Quick Start
<a name="s2qQB"></a>
## 1 Start zookeeper
Go to the `zookeeper/bin` directory  and double-click zkServer.cmd to start zookeeper.
<a name="db9k8"></a>
## 2 Start dubbo-go provider

- Firstly, modify run configuration for samples/dubbogo/simple/body/server/app/server.go, and add two environment variables.
> - CONF_PROVIDER_FILE_PATH：samples\dubbogo\simple\body\server\profiles\dev\server.yml
> - APP_LOG_CONF_FILE：samples\dubbogo\simple\body\server\profiles\dev\log.yml

- Secondly, execute main function of server.go.
  <a name="YzNS4"></a>
## 3 Start pixiu
Pixiu starts the application through command parameters, so we only need to pass the command and parameters to the main function when starting.<br />

- Firstly, modify run configuration for cmd/pixiu/pixiu.go, and add program startup arguments `gateway start -c samples\dubbogo\simple\body\pixiu\conf.yaml`.
- Secondly, modify the file `samples\dubbogo\simple\body\pixiu\conf.yaml`, set the path field as `samples\dubbogo\simple\body\pixiu\api_config.yaml`.
```yaml
http_filters:
  - name: dgp.filter.http.apiconfig
    config:
      path: samples\dubbogo\simple\body\pixiu\api_config.yaml
```

- Finally, execute main function of pixiu.go.
  <a name="W1LMZ"></a>
## 4 Test
Open cmd terminal，and enter the following command lines :
```bash
curl -s -X POST "localhost:8881/api/v1/test-dubbo/user" -d "{\"id\":\"0003\",\"code\":3,\"name\":\"dubbogo\",\"age\":99}" --header "Content-Type: application/json"
curl -s -X PUT "localhost:8881/api/v1/test-dubbo/user" -d "{\"id\":\"0003\",\"code\":3,\"name\":\"dubbogo\",\"age\":77}" --header "Content-Type: application/json"
curl -s -X PUT "localhost:8881/api/v1/test-dubbo/user2" -d "{\"name\":\"dubbogo\",\"user\":{\"id\":\"0003\",\"code\":3,\"name\":\"dubbogo\",\"age\":88}}" --header "Content-Type: application/json"
```
output：
```
{"age":99,"code":3,"iD":"0003","name":"dubbogo"}
true
true
```
Congratulations, successful startup.
