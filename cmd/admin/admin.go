package main

import (
	"encoding/json"
	etcdv3 "github.com/dubbogo/dubbo-go-proxy/pkg/remoting/etcd3"
	"net/http"
	"strings"
	"time"
)

var client *etcdv3.Client

type SetConfig struct {
	Config string
}

func main() {

	client = etcdv3.NewServiceDiscoveryClient(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split("127.0.0.1:2379", ",")...),
	)

	http.HandleFunc("/config/api", GetApiConfig)
	http.HandleFunc("/config/api/set", SetApiConfig)

	http.ListenAndServe(":8080", nil)

}

func GetApiConfig(w http.ResponseWriter, req *http.Request) {
	config, err := client.Get("/proxy/config/api")
	if err != nil {
		w.Write([]byte("Error"))
	}
	w.Write([]byte(config))
}

func SetApiConfig(w http.ResponseWriter, req *http.Request) {

	var p SetConfig
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	setErr := client.Create("/proxy/config/api", p.Config)
	if setErr != nil {
		w.Write([]byte(setErr.Error()))
	}
	w.Write([]byte("Success"))
}
