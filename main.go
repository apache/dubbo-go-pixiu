package main

import (
	_ "github.com/dubbogo/dubbo-go-proxy/common/config"
	_ "github.com/dubbogo/dubbo-go-proxy/service/filter"
	"github.com/dubbogo/dubbo-go-proxy/service/proxy/http"
)

// for example
func main() {
	http.Run()
}
