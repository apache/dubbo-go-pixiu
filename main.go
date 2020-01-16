package main

import (
	_ "github.com/pantianying/dubbo-go-proxy/common/config"
	_ "github.com/pantianying/dubbo-go-proxy/service/filter"
	"github.com/pantianying/dubbo-go-proxy/service/proxy/http"
)

// for example
func main() {
	http.Run()
}
