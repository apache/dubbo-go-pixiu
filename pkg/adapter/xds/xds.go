/*
 * Licensed to the Apache Software Foundation (ASF) under one or more * contributor license agreements.  See the NOTICE file distributed with
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

package xds

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/model/xds"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

func init() {
	adapter.RegisterAdapterPlugin(&DiscoveryPlugin{})
}

type (
	DiscoverApi interface {
		Fetch(localVersion string) ([]*apiclient.ProtoAny, error)
		Delta() (chan *apiclient.DeltaResources, error)
	}

	AdapterConfig struct {
	}

	Adapter struct {
		ID     string
		Name   string
		cfg    *AdapterConfig
		ads    DiscoverApi //aggregate discover service manager todo to implement
		cds    *CdsManager //cluster discover service manager
		lds    *LdsManager //listener discover service manager
		exitCh chan struct{}
	}

	GrpcClusterManager struct {
		clusters *sync.Map // map[clusterName]*GrpcCluster
	}

	GrpcCluster struct {
		name   string //cluster name
		config *model.Cluster
		once   sync.Once
		conn   *grpc.ClientConn
	}
)

// GetGrpcCluster get the cluster or create it first time.
func (g *GrpcClusterManager) GetGrpcCluster(name string) (apiclient.GrpcCluster, error) {
	store, err := server.GetClusterManager().CloneStore()
	if err != nil {
		return nil, errors.WithMessagef(err, "clone cluster store failed")
	}
	if !store.HasCluster(name) {
		return nil, errors.Errorf("can not find cluster of %s", name)
	}

	if load, ok := g.clusters.Load(name); ok {
		grpcCluster := load.(*GrpcCluster) // grpcClusterManager only
		return grpcCluster, nil
	}
	var clusterCfg *model.Cluster
	for _, cfg := range store.Config {
		if cfg.Name == name {
			clusterCfg = cfg
			break
		}
	}
	if clusterCfg == nil {
		return nil, errors.Errorf("can not find cluster of %s", name)
	}
	newCluster := &GrpcCluster{
		name:   name,
		config: clusterCfg,
	}
	g.clusters.Store(name, newCluster)
	return newCluster, nil
}

func (g *GrpcClusterManager) Close() (err error) {
	//todo enhance the close process when concurrent
	g.clusters.Range(func(_, value interface{}) bool {
		if conn := value.(*grpc.ClientConn); conn != nil {
			if err = conn.Close(); err != nil {
				logger.Errorf("can not close grpc connection.", err)
			}
		}
		return true
	})
	return nil
}

func (g *GrpcCluster) GetConnect() *grpc.ClientConn {
	g.once.Do(func() {
		creds := insecure.NewCredentials()
		//if *xdsCreds { // todo
		//	log.Println("Using xDS credentials...")
		//	var err error
		//	if creds, err = xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()}); err != nil {
		//		log.Fatalf("failed to create client-side xDS credentials: %v", err)
		//	}
		//}
		if len(g.config.Endpoints) == 0 {
			panic("expect endpoint.")
		}
		address := g.config.Endpoints[0].Address.Address
		port := g.config.PickOneEndpoint().Address.Port
		target := fmt.Sprintf("%s:%d", address, port)
		logger.Infof("to connect xds server %s ...", target)
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //todo fix timeout cancel warning
		conn, err := grpc.DialContext(ctx, target,
			grpc.WithTransportCredentials(creds),
			grpc.WithBlock(),
		)
		logger.Infof("connected xds server (%s)", target)
		if err != nil {
			panic(fmt.Sprintf("grpc.Dial(%s) failed: %v", target, err))
		}
		g.conn = conn
	})
	return g.conn
}

func (g *GrpcCluster) IsAlive() bool {
	return true
}

func (g *GrpcCluster) Close() error {
	if err := g.conn.Close(); err != nil {
		return errors.Wrapf(err, "can not close. %v", g.config)
	}
	return nil
}

var (
	grpcClusterManager = GrpcClusterManager{
		clusters: &sync.Map{},
	}
)

func (a *Adapter) createApiManager(config *model.ApiConfigSource,
	node *model.Node,
	resourceTypes ...apiclient.ResourceTypeName) DiscoverApi {
	if config == nil {
		return nil
	}

	switch config.APIType {
	case model.ApiTypeGRPC:
		return apiclient.CreateGrpcApiClient(config, node, &grpcClusterManager, a.exitCh, resourceTypes...)
	default:
		logger.Fatalf("un-support the api type %s", config.APITypeStr)
		return nil
	}
}

func (a *Adapter) Start() {
	dm := server.GetDynamicResourceManager()
	// lds fetch just run on init phase.
	if dm.GetLds() != nil {
		a.lds = &LdsManager{DiscoverApi: a.createApiManager(dm.GetLds(), dm.GetNode(), xds.ListenerType)}
		if err := a.lds.Fetch(); err != nil {
			logger.Errorf("can not fetch lds")
		}
	}
	// catch the ongoing cds config change.
	if dm.GetCds() != nil {
		a.cds = &CdsManager{DiscoverApi: a.createApiManager(dm.GetCds(), dm.GetNode(), xds.ClusterType)}
		if err := a.cds.Delta(); err != nil {
			logger.Errorf("can not fetch lds")
		}
	}

}

func (a *Adapter) Stop() {
	//todo close the grpc connection
	if err := grpcClusterManager.Close(); err != nil {
		logger.Errorf("grpcClusterManager close failed. %v", err)
	}
	close(a.exitCh)
}

func (a *Adapter) Apply() error {
	// do nothing
	return nil
}

func (a *Adapter) Config() interface{} {
	return a.cfg
}
