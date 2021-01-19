package main

import (
	"fmt"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	etcdv3 "github.com/dubbogo/dubbo-go-proxy/pkg/remoting/etcd3"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client *etcdv3.Client


const (
	// EtcdAddress etcd address
	EtcdAddress   string = "127.0.0.1:2379"
	// AdminAddress admin server host and port
	AdminAddress  string = "127.0.0.1:8080"
	// APIConfigPath api config path in etcd
	APIConfigPath string = "/proxy/config/api"
)

func main() {

	client = etcdv3.NewConfigClient(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split(EtcdAddress, ",")...),
	)
	defer client.Close()

	http.HandleFunc("/config/api", GetAPIConfig)
	http.HandleFunc("/config/api/set", SetAPIConfig)

	http.ListenAndServe(AdminAddress, nil)

}

// GetAPIConfig handle get api config http request
func GetAPIConfig(w http.ResponseWriter, req *http.Request) {
	config, err := client.Get(APIConfigPath)
	if err != nil {
		fmt.Printf(err.Error())
		w.Write([]byte("Error"))
	}
	w.Write([]byte(config))
}

// SetAPIConfig handle modify api config http request
func SetAPIConfig(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return
	}
	// validate the api config
	apiConf := &config.APIConfig{}
	err = yaml.UnmarshalYML([]byte(body), apiConf)

	if err != nil {
		fmt.Printf("UnmarshalYML body error, %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	setErr := client.Update(APIConfigPath, string(body))
	if setErr != nil {
		fmt.Printf("update etcd error, %v", setErr)
		w.Write([]byte(setErr.Error()))
	}
	fmt.Printf("success")
	w.Write([]byte("Success"))
}
