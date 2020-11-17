# dubbo error

err log in proxy.

- dubbo server return error:

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