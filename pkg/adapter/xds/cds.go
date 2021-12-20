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
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	model2 "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds/model"
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
	clusters := make([]*model2.Cluster, 0, len(r))
	for _, one := range r {
		extClusters := &model2.PixiuExtensionClusters{}
		if err := one.To(extClusters); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("clusters from xds server %v", extClusters)
		clusters = append(clusters, extClusters.Clusters...)
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
		clusters := make([]*model2.Cluster, 0, len(one.NewResources))
		for _, one := range one.NewResources {
			_cluster := &model2.PixiuExtensionClusters{}
			if err := one.To(_cluster); err != nil {
				logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
				continue
			}
			logger.Infof("clusters from xds server %v", _cluster)
			clusters = append(clusters, _cluster.Clusters...)

		}
		if err := c.setupCluster(clusters); err != nil {
			logger.Errorf("can not setup cluster.", err)
		}
	}
}

func (c *CdsManager) removeCluster(clusterNames []string) {
	server.GetClusterManager().RemoveCluster(clusterNames)
}

func (c *CdsManager) setupCluster(clusters []*model2.Cluster) error {
	clusterMgt := server.GetClusterManager()
	laterApplies := make([]func() error, 0, len(clusters))
	toRemoveHash := make(map[string]struct{}, len(clusters))

	for _, _cluster := range clusters {
		toRemoveHash[_cluster.Name] = struct{}{}
	}
	for _, cluster := range clusters {
		delete(toRemoveHash, cluster.Name)
		switch {
		case clusterMgt.HasCluster(cluster.Name):
			laterApplies = append(laterApplies, func() error {
				clusterMgt.UpdateCluster(c.makeCluster(cluster))
				return nil
			})
		default:
			laterApplies = append(laterApplies, func() error {
				clusterMgt.AddCluster(c.makeCluster(cluster))
				return nil
			})
		}
	}

	c.removeClusters(toRemoveHash)
	for _, fn := range laterApplies { //do update and add new cluster.
		if err := fn(); err != nil {
			logger.Errorf("can not modify cluster", err)
		}
	}
	return nil
}

func (c *CdsManager) removeClusters(toRemoveList map[string]struct{}) {
	removeClusters := make([]string, 0, len(toRemoveList))
	for clusterName, _ := range toRemoveList {
		removeClusters = append(removeClusters, clusterName)
	}
}

func (c *CdsManager) makeCluster(cluster *model2.Cluster) *model.Cluster {
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

func (c *CdsManager) makeLoadBalancePolicy(cluster *model2.Cluster) model.LbPolicy {
	return model.LbPolicy(model.LbPolicyValue[cluster.LbStr])
}

func (c *CdsManager) makeClusterType(cluster *model2.Cluster) model.DiscoveryType {
	return model.DiscoveryType(model.DiscoveryTypeValue[cluster.TypeStr])
}

func (c *CdsManager) makeEndpoints(endpoint *model2.Endpoint) []*model.Endpoint {
	r := make([]*model.Endpoint, 0, 1)
	r = append(r, &model.Endpoint{
		ID:       endpoint.Id,
		Name:     endpoint.Name,
		Address:  c.makeAddress(endpoint),
		Metadata: endpoint.Metadata,
	})
	return r
}

func (c *CdsManager) makeAddress(endpoint *model2.Endpoint) model.SocketAddress {
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

func (c *CdsManager) makeProtocol(endpoint *model2.Endpoint) model.ProtocolType {
	return model.ProtocolType(model.ProtocolTypeValue[strings.ToUpper(endpoint.Address.ProtocolStr)])
}

func (c *CdsManager) makeHealthChecks(checks []*model2.HealthCheck) (result []model.HealthCheck) {
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

func (c *CdsManager) makeEdsClusterConfig(edsConfig *model2.EdsClusterConfig) model.EdsClusterConfig {
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
