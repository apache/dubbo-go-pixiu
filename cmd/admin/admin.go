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
	ETCD_ADDRESS string = "127.0.0.1:2379"
	ADMIN_ADDRESS string = "127.0.0.1:8080"
	API_CONFIG_PATH string  = "/proxy/config/api"
)

func main() {

	client = etcdv3.NewConfigClient(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split(ETCD_ADDRESS, ",")...),
	)
	defer client.Close()

	http.HandleFunc("/config/api", GetApiConfig)
	http.HandleFunc("/config/api/set", SetApiConfig)

	http.ListenAndServe(ADMIN_ADDRESS, nil)

}

func GetApiConfig(w http.ResponseWriter, req *http.Request) {
	config, err := client.Get(API_CONFIG_PATH)
	if err != nil {
		fmt.Printf(err.Error())
		w.Write([]byte("Error"))
	}
	w.Write([]byte(config))
}

func SetApiConfig(w http.ResponseWriter, req *http.Request) {

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

	setErr := client.Update(API_CONFIG_PATH, string(body))
	if setErr != nil {
		fmt.Printf("update etcd error, %v", setErr)
		w.Write([]byte(setErr.Error()))
	}
	fmt.Printf("success")
	w.Write([]byte("Success"))
}
