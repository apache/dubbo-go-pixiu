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
	etcdv3 "github.com/dubbogo/dubbo-go-proxy/pkg/remoting/etcd3"
	perrors "github.com/pkg/errors"
	"github.com/urfave/cli"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)


type AdminBootstrap struct {
	Address  string  `yaml:"address" json:"address" mapstructure:"address"`
	EtcdConfig EtcdConfig `yaml:"etcd" json:"etcd" mapstructure:"etcd"`
}

type EtcdConfig struct {
	Address  string  `yaml:"admin" json:"admin" mapstructure:"admin"`
	Path string  `yaml:"path" json:"path" mapstructure:"path"`
}



var (
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start dubbogo proxy admin",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: "DUBBOGO_PROXY_CONFIG",
				Value:  "configs/conf.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			_ , err := LoadAPIConfigFromFile(configPath)
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
		Usage: "stop dubbogo proxy",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)


func newAdminApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "dubbogo proxy admin"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo"
	app.Usage = "Dubbogo proxy is a lightweight gateway."
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
	bootstrap = adminBootstrap;
	return adminBootstrap, nil
}


func main() {
	app := newAdminApp(&cmdStart)
	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Run(os.Args)
}


var client *etcdv3.Client
var bootstrap *AdminBootstrap

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

	http.ListenAndServe(bootstrap.Address, nil)
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
	apiConf := &config.APIConfig{}
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
