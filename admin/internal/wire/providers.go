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

package wire

import (
	gxetcd "github.com/dubbogo/gost/database/kv/etcd/v3"

	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/handler"
	"github.com/apache/dubbo-go-pixiu/admin/internal/store"
	"github.com/apache/dubbo-go-pixiu/admin/internal/xds"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
)

// ConfigPath is a distinct type for the config file path.
type ConfigPath string

// EtcdBasePath is a distinct type for the etcd base path.
type EtcdBasePath string

// XDSPort is a distinct type for the xDS server port.
type XDSPort uint

// ProvideConfig loads the configuration from the given path.
func ProvideConfig(path ConfigPath) (*config.Config, error) {
	return config.Load(string(path))
}

// ProvideLogger creates a new zap logger from config.
func ProvideLogger(cfg *config.Config) (*zap.Logger, error) {
	return config.NewLogger(cfg.Zap)
}

// ProvideMySQLConfig extracts MySQL config from the main config.
func ProvideMySQLConfig(cfg *config.Config) config.MySQLConfig {
	return cfg.MySQL
}

// ProvideEtcdConfig extracts Etcd config from the main config.
func ProvideEtcdConfig(cfg *config.Config) config.EtcdConfig {
	return cfg.Etcd
}

// ProvideJWTConfig extracts JWT config from the main config.
func ProvideJWTConfig(cfg *config.Config) config.JWTConfig {
	return cfg.JWT
}

// ProvideEtcdClient creates a new etcd client from config.
func ProvideEtcdClient(cfg config.EtcdConfig) (*gxetcd.Client, func(), error) {
	client, err := config.NewEtcdClient(cfg)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		client.Close()
	}
	return client, cleanup, nil
}

// ProvideEtcdBasePath extracts the etcd base path from config.
func ProvideEtcdBasePath(cfg config.EtcdConfig) EtcdBasePath {
	return EtcdBasePath(cfg.Path)
}

// ProvideXDSPort returns the xDS server port.
func ProvideXDSPort() XDSPort {
	return 18000
}

// ProvideEtcdStore creates the Etcd store.
func ProvideEtcdStore(client *gxetcd.Client, basePath EtcdBasePath) *store.Etcd {
	return store.NewEtcd(client, string(basePath))
}

// ProvideXDSServer creates the xDS server.
func ProvideXDSServer(etcd *store.Etcd, logger *zap.Logger, port XDSPort) *xds.Server {
	return xds.NewServer(etcd, logger, uint(port))
}

// App holds all application dependencies.
type App struct {
	Config    *config.Config
	Logger    *zap.Logger
	MySQL     *store.MySQL
	Etcd      *store.Etcd
	Handler   *handler.Handler
	XDSServer *xds.Server
}

// NewApp creates a new App with all dependencies.
func NewApp(
	cfg *config.Config,
	logger *zap.Logger,
	mysql *store.MySQL,
	etcd *store.Etcd,
	h *handler.Handler,
	xdsServer *xds.Server,
) *App {
	return &App{
		Config:    cfg,
		Logger:    logger,
		MySQL:     mysql,
		Etcd:      etcd,
		Handler:   h,
		XDSServer: xdsServer,
	}
}
