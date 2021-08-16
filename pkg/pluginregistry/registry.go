package pluginregistry

import (
	// http filters
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/accesslog"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/api"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/authority"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/header"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/host"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/metric"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/response"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/timeout"
	// adapter
)
