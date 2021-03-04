# FAQ

In case you did not find any answer here and in [closed issues](https://github.com/dubbogo/dubbo-go-pixiu/issues?q=is%3Aissue+is%3Aclosed), [create new issue](https://github.com/dubbogo/dubbo-go-pixiu/issues/new/choose).

### Why first dubbogo request timeout?

First request for dubbogo need create a GenericService, will exec ReferenceConfig.GenericLoad(), it takes some account of time. If you worry about this case, you need preheating apis for yourself.

We plan an auto preheating in the future.



 