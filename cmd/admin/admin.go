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
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	perrors "github.com/pkg/errors"
	"github.com/urfave/cli"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	etcdv3 "github.com/apache/dubbo-go-pixiu/pkg/remoting/etcd3"
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

var (
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start apache pixiu admin",
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
		Usage: "stop apache pixiu admin",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)

// Version admin version
const Version = "0.1.0"

func newAdminApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "apache pixiu admin"
	app.Version = Version
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo"
	app.Usage = "Dubbogo pixiu admin"
	app.Flags = cmdStart.Flags

	//commands
	app.Commands = []cli.Command{
		cmdStart,
		cmdStop,
	}

	//action
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
	client = etcdv3.NewConfigClient(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split(bootstrap.EtcdConfig.Address, ",")...),
	)
	defer client.Close()

	http.HandleFunc("/config/api", GetAPIConfig)
	http.HandleFunc("/config/api/set", SetAPIConfig)

	http.ListenAndServe(bootstrap.Server.Address, nil)
}

// GetAPIConfig handle get api config http request
func GetAPIConfig(w http.ResponseWriter, req *http.Request) {
	config, err := client.Get(bootstrap.EtcdConfig.Path)
	if err != nil {
		logger.Errorf("GetAPIConfig err, %v\n", err)
		w.Write([]byte("Error"))
	}
	w.Write([]byte(config))
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
