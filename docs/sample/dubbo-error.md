# How to check dubbo error

## Proxy bug condition

Server log.

```bash
2020-11-19T20:30:26.070+0800    ERROR   filter_impl/generic_service_filter.go:98        [Generic Service Filter] method:GetUserTimeout invocation arguments number was wrong
github.com/apache/dubbo-go/common/logger.Errorf
```

If this case happens, proxy will return nil because dubbo server return nil. You can see dubbo log follow when log format is debug. 

```bash
2020-11-19T20:30:26.072+0800    DEBUG   proxy/proxy.go:172      [makeDubboCallProxy] result: 0xc0001aeeb0, err: <nil>
2020-11-19T20:30:26.072+0800    DEBUG   dubbo/dubbo.go:152      [dubbo-go-proxy] dubbo client resp:<nil>
2020-11-19T20:30:26.073+0800    DEBUG   remote/call.go:117      [dubbo-go-proxy] client call resp:<nil>
```

When generic invoke, err will return nil, because wrote the code as below.

```go
	if len(oldParams) != len(argsType) {
		logger.Errorf("[Generic Service Filter] method:%s invocation arguments number was wrong", methodName)
		return &protocol.RPCResult{}
	}
```

## Err log in proxy

Dubbo server return error:

```bash
2020-11-17T11:19:18.019+0800    ERROR   remote/call.go:87       [dubbo-go-proxy] client call err:data is exist!
github.com/dubbogo/dubbo-go-proxy/pkg/logger.Errorf
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/logger/logging.go:52
github.com/dubbogo/dubbo-go-proxy/pkg/filter/remote.(*clientFilter).doRemoteCall
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/filter/remote/call.go:87
github.com/dubbogo/dubbo-go-proxy/pkg/filter/remote.(*clientFilter).Do.func1
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/filter/remote/call.go:65
github.com/dubbogo/dubbo-go-proxy/pkg/context/http.(*HttpContext).Next
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/context/http/context.go:54
github.com/dubbogo/dubbo-go-proxy/pkg/filter/timeout.(*timeoutFilter).Do.func1.1
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/filter/timeout/timeout.go:70
```

Return value

```json
{
    "message": "data is exist"
}
```

## No provider info

```bash
2020-11-20T15:56:59.011+0800    ERROR   remote/call.go:112      [dubbo-go-proxy] client call err:Failed to invoke the method $invoke. No provider available for the service dubbo://:@:/?interface=com.ic.user.UserProvider&group=test&version=1.0.0 from registry zookeeper://127.0.0.1:2181?group=&registry=zookeeper&registry.label=true&registry.preferred=false&registry.role=0&registry.timeout=3s&registry.ttl=&registry.weight=0&registry.zone=&simplified=false on the consumer 30.11.176.51 using the dubbo version 1.3.0 .Please check if the providers have been started and registered.!
```

[Previous](./dubbo.md)