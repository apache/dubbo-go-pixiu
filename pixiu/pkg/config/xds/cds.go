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
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api"
	xdspb "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
)

type CdsManager struct {
	DiscoverApi
	clusterMg controls.ClusterManager
}

// Fetch overwrite DiscoverApi.Fetch.
func (c *CdsManager) Fetch() error {
	r, err := c.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		return err
	}
	clusters := make([]*xdspb.Cluster, 0, len(r))
	for _, one := range r {
		extClusters := &xdspb.PixiuExtensionClusters{}
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
		clusters := make([]*xdspb.Cluster, 0, len(one.NewResources))
		for _, one := range one.NewResources {
			cluster := &xdspb.PixiuExtensionClusters{}
			if err := one.To(cluster); err != nil {
				logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
				continue
			}
			logger.Infof("clusters from xds server %v", cluster)
			clusters = append(clusters, cluster.Clusters...)

		}
		if err := c.setupCluster(clusters); err != nil {
			logger.Errorf("can not setup cluster.", err)
		}
	}
}

func (c *CdsManager) removeCluster(clusterNames []string) {
	c.clusterMg.RemoveCluster(clusterNames)
}

func (c *CdsManager) setupCluster(clusters []*xdspb.Cluster) error {

	laterApplies := make([]func() error, 0, len(clusters))
	toRemoveHash := make(map[string]struct{}, len(clusters))

	store, err := c.clusterMg.CloneXdsControlStore()
	if err != nil {
		return errors.WithMessagef(err, "can not clone cluster store when update cluster")
	}
	//todo this will remove the cluster which defined locally.
	for _, cluster := range store.Config() {
		toRemoveHash[cluster.Name] = struct{}{}
	}
	for _, cluster := range clusters {
		delete(toRemoveHash, cluster.Name)

		makeCluster := c.makeCluster(cluster)
		switch {
		case c.clusterMg.HasCluster(cluster.Name):
			laterApplies = append(laterApplies, func() error {
				c.clusterMg.UpdateCluster(makeCluster)
				return nil
			})
		default:
			laterApplies = append(laterApplies, func() error {
				c.clusterMg.AddCluster(makeCluster)
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
	for clusterName := range toRemoveList {
		removeClusters = append(removeClusters, clusterName)
	}
	if len(toRemoveList) == 0 {
		return
	}
	c.removeCluster(removeClusters)
}

func (c *CdsManager) makeCluster(cluster *xdspb.Cluster) *model.ClusterConfig {
	return &model.ClusterConfig{
		Name:             cluster.Name,
		TypeStr:          cluster.TypeStr,
		Type:             c.makeClusterType(cluster),
		EdsClusterConfig: c.makeEdsClusterConfig(cluster.EdsClusterConfig),
		LbStr:            c.makeLoadBalancePolicy(cluster.LbStr),
		HealthChecks:     c.makeHealthChecks(cluster.HealthChecks),
		Endpoints:        c.makeEndpoints(cluster.Endpoints),
	}
}

func (c *CdsManager) makeLoadBalancePolicy(lb string) model.LbPolicyType {
	return model.LbPolicyTypeValue[lb]
}

func (c *CdsManager) makeClusterType(cluster *xdspb.Cluster) model.DiscoveryType {
	return model.DiscoveryTypeValue[cluster.TypeStr]
}

func (c *CdsManager) makeEndpoints(endpoints []*xdspb.Endpoint) []*model.Endpoint {
	r := make([]*model.Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		r[i] = &model.Endpoint{
			ID:       endpoint.Id,
			Name:     endpoint.Name,
			Address:  c.makeAddress(endpoint),
			Metadata: endpoint.Metadata,
		}
	}
	return r
}

func (c *CdsManager) makeAddress(endpoint *xdspb.Endpoint) model.SocketAddress {
	if endpoint == nil || endpoint.Address == nil {
		return model.SocketAddress{}
	}
	return model.SocketAddress{
		Address:      endpoint.Address.Address,
		Port:         int(endpoint.Address.Port),
		ResolverName: endpoint.Address.ResolverName,
		Domains:      endpoint.Address.Domains,
		CertsDir:     endpoint.Address.CertsDir,
	}
}

func (c *CdsManager) makeHealthChecks(checks []*xdspb.HealthCheck) (result []model.HealthCheckConfig) {
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
	if edsConfig == nil {
		return model.EdsClusterConfig{}
	}
	return model.EdsClusterConfig{
		EdsConfig: model.ConfigSource{
			Path:            edsConfig.EdsConfig.Path,
			ApiConfigSource: c.makeApiConfigSource(edsConfig.EdsConfig.ApiConfigSource),
		},
		ServiceName: edsConfig.ServiceName,
	}
}

func (c *CdsManager) makeApiConfigSource(apiConfig *xdspb.ApiConfigSource) (result model.ApiConfigSource) {
	apiType, ok := model.ApiTypeValue[apiConfig.APITypeStr]
	if !ok {
		logger.Errorf("unknown apiType %s", apiConfig.APITypeStr)
		return
	}

	return model.ApiConfigSource{
		APIType:        api.ApiType(apiType),
		APITypeStr:     apiConfig.APITypeStr,
		ClusterName:    apiConfig.ClusterName,
		RefreshDelay:   apiConfig.RefreshDelay,
		RequestTimeout: apiConfig.RequestTimeout,
		GrpcServices:   nil, //todo create node of pb
	}
}
