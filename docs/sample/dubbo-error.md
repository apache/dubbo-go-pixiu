# dubbo error

err log in proxy

```bash
2020-11-14T17:56:58.652+0800    ERROR   remote/call.go:89       [dubboproxy go] client do err:data is exist!
github.com/dubbogo/dubbo-go-proxy/pkg/logger.Errorf
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/logger/logging.go:52
github.com/dubbogo/dubbo-go-proxy/pkg/filter/remote.(*clientFilter).doRemoteCall
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/filter/remote/call.go:89
github.com/dubbogo/dubbo-go-proxy/pkg/filter/remote.(*clientFilter).Do.func1
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/filter/remote/call.go:66
github.com/dubbogo/dubbo-go-proxy/pkg/context/http.(*HttpContext).Next
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/context/http/context.go:54
github.com/dubbogo/dubbo-go-proxy/pkg/filter/timeout.(*timeoutFilter).Do.func1.1
        /Users/tc/Documents/workspace_2020/dubbo-go-proxy/pkg/filter/timeout/timeout.go:74
```