/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/gin-gonic/gin"
	perrors "github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	etcdv3 "github.com/dubbogo/gost/database/kv/etcd/v3"
)

// AdminBootstrap admin bootstrap config
type AdminBootstrap struct {
	Server     ServerConfig `yaml:"server" json:"server" mapstructure:"server"`
	EtcdConfig EtcdConfig   `yaml:"etcd" json:"etcd" mapstructure:"etcd"`
}

// ServerConfig admin http server config
type ServerConfig struct {
	Address string `yaml:"address" json:"address" mapstructure:"address"`
}

// EtcdConfig admin etcd client config
type EtcdConfig struct {
	Address string `yaml:"address" json:"admin" mapstructure:"admin"`
	Path    string `yaml:"path" json:"path" mapstructure:"path"`
}

type BaseInfo struct {
	Name           string `json:"name" yaml:"name"`
	Description    string `json:"description" yaml:"description"`
	PluginFilePath string `json:"pluginFilePath" yaml:"pluginFilePath"`
}

type RetData struct {
	Code string      `json:"code" yaml:"code"`
	Data interface{} `json:"data" yaml:"data"`
}

func WithError(err error) RetData {
	return RetData{ERR, err.Error()}
}

func WithRet(data interface{}) RetData {
	return RetData{OK, data}
}

var (
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start dubbogo pixiu admin",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: "PROXY_ADMIN_CONFIG",
				Value:  "configs/admin_config.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			_, err := LoadAPIConfigFromFile(configPath)
			if err != nil {
				logger.Errorf("load admin config  error:%+v", err)
				return err
			}
			Start()
			return nil
		},
	}

	cmdStop = cli.Command{
		Name:  "stop",
		Usage: "stop dubbogo pixiu admin",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)

// Version admin version
const Version = "0.1.0"
const Base = "base"
const Resources = "Resources"
const Plugin = "plugin"
const OK = "10001"
const ERR = "10002"

func newAdminApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "dubbogo pixiu admin"
	app.Version = Version
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo"
	app.Usage = "Dubbogo pixiu admin"
	app.Flags = cmdStart.Flags

	// commands
	app.Commands = []cli.Command{
		cmdStart,
		cmdStop,
	}

	// action
	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}
		return startCmd.Action.(func(c *cli.Context) error)(c)
	}

	return app
}

// LoadAPIConfigFromFile load config from file
func LoadAPIConfigFromFile(path string) (*AdminBootstrap, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	adminBootstrap := &AdminBootstrap{}
	err := yaml.UnmarshalYMLConfig(path, adminBootstrap)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	bootstrap = adminBootstrap
	return adminBootstrap, nil
}

func main() {
	app := newAdminApp(&cmdStart)
	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Run(os.Args)
}

var (
	client    *etcdv3.Client
	bootstrap *AdminBootstrap
)

// Start start init etcd client and start admin http server
func Start() {
	newClient, err := etcdv3.NewConfigClient(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(20*time.Second),
		etcdv3.WithEndpoints(strings.Split(bootstrap.EtcdConfig.Address, ",")...),
	)

	if err != nil {
		logger.Errorf("update etcd error, %v\n", err)
		return
	}

	client = newClient
	defer client.Close()

	r := SetupRouter()
	r.Run(bootstrap.Server.Address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/config/api/base", GetBaseInfo)
	r.POST("/config/api/base/set", SetBaseInfo)
	r.GET("/config/api/resource/list", GetResourceList)
	r.POST("/config/api/resource/create", CreateResourceInfo)
	//r.POST("/config/api/resource/modify", UpdateResourceInfo)

	return r
}

// SetBaseInfo handle modify base info http request
func SetBaseInfo(c *gin.Context) {
	body := c.PostForm("content")

	baseInfo := &BaseInfo{}
	err := yaml.UnmarshalYML([]byte(body), baseInfo)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	setErr := BizSetBaseInfo(baseInfo)

	if setErr != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.String(http.StatusOK, "success")
}

func BizSetBaseInfo(info *BaseInfo) error {
	// validate the api config

	data, _ := yaml.MarshalYML(info)
	setErr := client.Update(getRootPath(Base), string(data))

	if setErr != nil {
		logger.Warnf("update etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetBaseInfo error")
	}
	return nil
}

// GetBaseInfo get base info
func GetBaseInfo(c *gin.Context) {
	config, err := BizGetBaseInfo()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := yaml.MarshalYML(config)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

func BizGetBaseInfo() (*BaseInfo, error) {
	config, err := client.Get(getRootPath(Base))
	if err != nil {
		logger.Errorf("GetBaseInfo err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetBaseInfo error")
	}
	data := &BaseInfo{}
	_ = yaml.UnmarshalYML([]byte(config), data)

	return data, nil
}

// GetResourceList get all resource list
func GetResourceList(c *gin.Context) {
	res, err := BizGetResourceList()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

func BizGetResourceList() ([]fc.Resource, error) {
	_, vList, err := client.GetChildrenKVList(getRootPath(Resources))
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.Resource
	for _, v := range vList {
		res := fc.Resource{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, res)
	}

	return ret, nil
}

// SetResourceInfo modify resource info exclude
func SetResourceInfo(w http.ResponseWriter, req *http.Request) {

}

func BizCreateResourceInfo(res *fc.Resource) error {

	methods := res.Methods
	// 创建 methods
	BizCreateResourceMethod(res.Path, methods)
	// 创建resource
	res.Methods = nil
	data, _ := yaml.MarshalYML(res)
	setErr := client.Create(getRootPath(Resources)+"/"+res.Path, string(data))

	if setErr != nil {
		logger.Warnf("update etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizCreateResourceInfo error")
	}
	return nil
}

func BizModifyResourceInfo(res fc.Resource) {

}

// CreateResourceInfo create resource
func CreateResourceInfo(c *gin.Context) {
	body := c.PostForm("content")

	res := &fc.Resource{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := BizCreateResourceInfo(res)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// DeleteResourceInfo delete resource
func DeleteResourceInfo(w http.ResponseWriter, req *http.Request) {

}

// GetResourceMethodList get method info list below resource
func GetResourceMethodList(w http.ResponseWriter, req *http.Request) {

}

// ModifyResourceMethod get method info below resource
func ModifyResourceMethod(w http.ResponseWriter, req *http.Request) {

}

// BizCreateResourceMethod batch create method below specific path
func BizCreateResourceMethod(root string, methods []fc.Method) error {
	var kList, vList []string

	for _, method := range methods {
		kList = append(kList, root+"/"+string(method.HTTPVerb))
		data, _ := yaml.MarshalYML(method)
		vList = append(vList, string(data))
	}

	err := client.BatchCreate(kList, vList)
	if err != nil {
		logger.Warnf("update etcd error, %v\n", err)
		return perrors.WithMessage(err, "BizCreateResourceMethod error")
	}
	return nil
}

// DeleteResourceMethod delete method below resource
func DeleteResourceMethod(w http.ResponseWriter, req *http.Request) {

}

// GetPluginGroupList get plugin group list
func GetPluginGroupList(w http.ResponseWriter, req *http.Request) {

}

// ModifyPluginGroup modify plugin group
func ModifyPluginGroup(w http.ResponseWriter, req *http.Request) {

}

// DeletePluginGroup delete plugin group
func DeletePluginGroup(w http.ResponseWriter, req *http.Request) {

}

// SetAPIConfig handle modify api config http request
func SetAPIConfig(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("read body err, %v\n", err)
		return
	}
	// validate the api config
	apiConf := &fc.APIConfig{}
	err = yaml.UnmarshalYML([]byte(body), apiConf)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	setErr := client.Update(bootstrap.EtcdConfig.Path, string(body))
	if setErr != nil {
		logger.Warnf("update etcd error, %v\n", err)
		w.Write([]byte(setErr.Error()))
	}
	w.Write([]byte("Success"))
}

func getRootPath(key string) string {
	return bootstrap.EtcdConfig.Path + "/" + key
}
