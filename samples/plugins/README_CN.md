# 插件式过滤器

通过插件式过滤器，你可以更加灵活地扩展Proxy的功能（无需重新打包Proxy），可以动态插拔和滚动升级（规划中）。

插件依赖于GO提供的[Plugin](https://golang.org/pkg/plugin/)， 有如下几个限制：

- 使用相同的Go版本

  > 插件实现和主应用程序都必须使用**完全相同**的 Go 工具链**版本**构建

- 使用相同的GOPATH

  > 插件实现和主应用程序都必须使用**完全相同**的 GOPATH 构建

- 使用相同版本的公共依赖项

  > 插件实现和主应用程序都必须使用**完全相同**的公共依赖项，在本例中唯一的公共依赖项是[github.com/dubbogo/dubbo-go-proxy-filter](https://github.com/dubbogo/dubbo-go-proxy-filter)。所以，请务必保证其与Proxy中的版本一致。

- 不支持Windows
  
  > 建议使用Docker或虚拟机



### 创建你的Plugin Filter
这里有一些[例子](./plugin)可供参考。



### 打包:

> 调试时可本地打包，生产环境或Windows建议使用Docker

###### 本地

```
go build -o out/plugin.so -buildmode=plugin ./plugin/*.go
```

###### Docker

> 将localFilePath替换为本地的目录，打包好的plugin.so文件用于后续使用

```
docker build -t plugin-test . && docker run --name plugin -v <localFilePath>:/go/src/out plugin-test
```

### 使用:

使用[配置文件](config/api_config.yaml)启动dubbogo proxy

> 在使用之前，请先更新配置文件中的插件加载目录“ pluginFilePath ”为“localFilePath/plugin.so”。	

访问并查看日志



