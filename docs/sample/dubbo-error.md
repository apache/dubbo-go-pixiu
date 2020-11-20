# dubbo error

## proxy bug condition

Server log.

```bash
2020-11-19T20:30:26.070+0800    ERROR   filter_impl/generic_service_filter.go:98        [Generic Service Filter] method:GetUserTimeout invocation arguments number was wrong
github.com/apache/dubbo-go/common/logger.Errorf
```

In case, if it happens, proxy will return nil because dubbo server return nil. You can see dubbo log follow when log format is debug. 

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

## err log in proxy

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
