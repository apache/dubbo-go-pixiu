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

package xds

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	xdspb "github.com/apache/dubbo-go-pixiu/pkg/model/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"strings"
)

type CdsManager struct {
	DiscoverApi
}

// Fetch overwrite DiscoverApi.Fetch.
func (c *CdsManager) Fetch() error {
	r, err := c.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		logger.Error("can not fetch lds", err)
		return err
	}
	clusters := make([]*xdspb.Cluster, 0, len(r))
	for _, one := range r {
		_cluster := &xdspb.Cluster{}
		if err := one.To(_cluster); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("clusters from xds server %v", _cluster)
		clusters = append(clusters, _cluster)
	}

	return c.setupCluster(clusters)
}

func (c *CdsManager) Delta() error {
	readCh, err := c.DiscoverApi.Delta()
	if err != nil {
		return err
	}
	go c.asyncHandler(readCh)
	return nil
}

func (c *CdsManager) asyncHandler(read chan *apiclient.DeltaResources) {
	for one := range read {
		clusters := make([]*xdspb.Cluster, 0, len(one.NewResources))
		for _, one := range one.NewResources {
			_cluster := &xdspb.Cluster{}
			if err := one.To(_cluster); err != nil {
				logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
				continue
			}
			logger.Infof("clusters from xds server %v", _cluster)
			clusters = append(clusters, _cluster)

		}
		c.removeCluster(one.RemovedResource)
		if len(clusters) > 0 {
			if err := c.setupCluster(clusters); err != nil {
				logger.Errorf("can not setup cluster.", err)
			}
		}
	}
}

func (c *CdsManager) removeCluster(clusterNames []string) {
	server.GetClusterManager().RemoveCluster(clusterNames)
}

func (c *CdsManager) setupCluster(clusters []*xdspb.Cluster) error {
	clusterMgt := server.GetClusterManager()
	for _, cluster := range clusters {
		switch {
		case clusterMgt.HasCluster(cluster.Name):
			clusterMgt.UpdateCluster(c.makeCluster(cluster))
		default:
			clusterMgt.AddCluster(c.makeCluster(cluster))
		}
	}
	return nil
}

func (c *CdsManager) makeCluster(cluster *xdspb.Cluster) *model.Cluster {
	return &model.Cluster{
		Name:             cluster.Name,
		TypeStr:          cluster.TypeStr,
		Type:             c.makeClusterType(cluster),
		EdsClusterConfig: c.makeEdsClusterConfig(cluster.EdsClusterConfig),
		LbStr:            cluster.LbStr,
		Lb:               c.makeLoadBalancePolicy(cluster),
		HealthChecks:     c.makeHealthChecks(cluster.HealthChecks),
		Endpoints:        c.makeEndpoints(cluster.Endpoints),
	}
}

func (c *CdsManager) makeLoadBalancePolicy(cluster *xdspb.Cluster) model.LbPolicy {
	return model.LbPolicy(model.LbPolicyValue[cluster.LbStr])
}

func (c *CdsManager) makeClusterType(cluster *xdspb.Cluster) model.DiscoveryType {
	return model.DiscoveryType(model.DiscoveryTypeValue[cluster.TypeStr])
}

func (c *CdsManager) makeEndpoints(endpoint *xdspb.Endpoint) []*model.Endpoint {
	r := make([]*model.Endpoint, 0, 1)
	r = append(r, &model.Endpoint{
		ID:       endpoint.Id,
		Name:     endpoint.Name,
		Address:  c.makeAddress(endpoint),
		Metadata: endpoint.Metadata,
	})
	return r
}

func (c *CdsManager) makeAddress(endpoint *xdspb.Endpoint) model.SocketAddress {
	if endpoint == nil {
		return model.SocketAddress{} // todo make default socket address
	}
	return model.SocketAddress{ // todo check nil
		ProtocolStr:  endpoint.Address.ProtocolStr,
		Protocol:     c.makeProtocol(endpoint),
		Address:      endpoint.Address.Address,
		Port:         int(endpoint.Address.Port),
		ResolverName: endpoint.Address.ResolverName,
		//Domains:      endpoint.Address.do, todo generate domains from proto
		//CertsDir: endpoint.Address., //todo generate certsDir
	}
}

func (c *CdsManager) makeProtocol(endpoint *xdspb.Endpoint) model.ProtocolType {
	return model.ProtocolType(model.ProtocolTypeValue[strings.ToUpper(endpoint.Address.ProtocolStr)])
}

func (c *CdsManager) makeHealthChecks(checks []*xdspb.HealthCheck) (result []model.HealthCheck) {
	//todo implement me after fix model.HealthCheck type define
	//result = make([]model.HealthCheck, 0, len(checks))
	//for _, check := range checks {
	//	switch one := check.GetChecker().(type) {
	//	case *xdspb.HealthCheck_HttpChecker:
	//		result = append(result, model.HttpHealthCheck{
	//			Host:             one.HttpChecker.Host,
	//			Path:             one.HttpChecker.Path,
	//			UseHttp2:         one.HttpChecker.UseHttp2,
	//			ExpectedStatuses: one.HttpChecker.ExpectedStatuses,
	//		})
	//	case *xdspb.HealthCheck_GrpcChecker:
	//		result = append(result, model.GrpcHealthCheck{
	//			ServiceName: one.GrpcChecker.ServiceName,
	//			Authority:   one.GrpcChecker.Authority,
	//		})
	//	case *xdspb.HealthCheck_CustomChecker:
	//		result = append(result, model.CustomHealthCheck{
	//			Name: one.CustomChecker.Name,
	//			Config: func() interface{} {
	//				if one.CustomChecker.Config == nil {
	//					return nil
	//				}
	//				return one.CustomChecker.Config.AsMap()
	//			}(),
	//		})
	//	}
	//}
	return
}

func (c *CdsManager) makeEdsClusterConfig(edsConfig *xdspb.EdsClusterConfig) model.EdsClusterConfig {
	//todo implement me
	if edsConfig == nil {
		return model.EdsClusterConfig{} //todo make default EdsClusterConfig
	}
	return model.EdsClusterConfig{
		EdsConfig: model.ConfigSource{
			Path: edsConfig.EdsConfig.Path,
		},
		ServiceName: edsConfig.ServiceName,
	}
}
