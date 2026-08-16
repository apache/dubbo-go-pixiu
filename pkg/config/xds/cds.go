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
	"strconv"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
)

type CdsManager struct {
	DiscoverApi
	clusterMg controls.ClusterManager
}

const endpointHealthMetadataKey = "pixiu.io/unhealthy"

// Fetch overwrite DiscoverApi.Fetch.
func (c *CdsManager) Fetch() error {
	r, err := c.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		return err
	}
	clusters := make([]*xdsmodel.Cluster, 0, len(r))
	for _, one := range r {
		extClusters := &xdsmodel.PixiuExtensionClusters{}
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
	for delta := range read {
		err := c.applyDelta(delta)
		delta.Complete(err)
		if err != nil {
			logger.Errorf("can not setup cluster: %v", err)
		}
	}
}

func (c *CdsManager) applyDelta(delta *apiclient.DeltaResources) error {
	if delta == nil {
		return nil
	}

	clusters := make([]*xdsmodel.Cluster, 0, len(delta.NewResources))
	for _, resource := range delta.NewResources {
		cluster := &xdsmodel.PixiuExtensionClusters{}
		if err := resource.To(cluster); err != nil {
			return errors.Wrapf(err, "unknown resource %q, expect Cluster", resource.GetName())
		}
		logger.Infof("clusters from xds server %v", cluster)
		clusters = append(clusters, cluster.Clusters...)
	}

	for _, name := range delta.RemovedResources {
		if name == constant.ClusterType {
			c.clusterMg.RemoveXDSClusters(c.clusterMg.XDSClusterNames())
		}
	}
	if len(delta.NewResources) == 0 {
		return nil
	}
	return c.setupCluster(clusters)
}

func (c *CdsManager) removeCluster(clusterNames []string) {
	c.clusterMg.RemoveXDSClusters(clusterNames)
}

func (c *CdsManager) setupCluster(clusters []*xdsmodel.Cluster) error {

	laterApplies := make([]func() error, 0, len(clusters))
	toRemoveHash := make(map[string]struct{}, len(clusters))

	for _, name := range c.clusterMg.XDSClusterNames() {
		toRemoveHash[name] = struct{}{}
	}
	for _, cluster := range clusters {
		delete(toRemoveHash, cluster.Name)

		makeCluster := c.makeCluster(cluster)
		laterApplies = append(laterApplies, func() error {
			return c.clusterMg.UpsertXDSCluster(makeCluster)
		})
	}

	c.removeClusters(toRemoveHash)
	for _, fn := range laterApplies { //do update and add new cluster.
		if err := fn(); err != nil {
			return errors.Wrap(err, "can not modify cluster")
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

func (c *CdsManager) makeCluster(cluster *xdsmodel.Cluster) *model.ClusterConfig {
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

func (c *CdsManager) makeClusterType(cluster *xdsmodel.Cluster) model.DiscoveryType {
	return model.DiscoveryTypeValue[cluster.TypeStr]
}

func (c *CdsManager) makeEndpoints(endpoints []*xdsmodel.Endpoint) []*model.Endpoint {
	r := make([]*model.Endpoint, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}
		unhealthy, _ := strconv.ParseBool(endpoint.Metadata[endpointHealthMetadataKey])
		r = append(r, &model.Endpoint{
			ID:        endpoint.Id,
			Name:      endpoint.Name,
			Address:   c.makeAddress(endpoint),
			Metadata:  endpoint.Metadata,
			UnHealthy: unhealthy,
		})
	}
	return r
}

func (c *CdsManager) makeAddress(endpoint *xdsmodel.Endpoint) model.SocketAddress {
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

func (c *CdsManager) makeHealthChecks(checks []*xdsmodel.HealthCheck) (result []model.HealthCheckConfig) {
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

func (c *CdsManager) makeEdsClusterConfig(edsConfig *xdsmodel.EdsClusterConfig) model.EdsClusterConfig {
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

func (c *CdsManager) makeApiConfigSource(apiConfig *xdsmodel.ApiConfigSource) (result model.ApiConfigSource) {
	apiType, ok := model.ApiTypeValue[apiConfig.APITypeStr]
	if !ok {
		logger.Errorf("unknown apiType %s", apiConfig.APITypeStr)
		return
	}

	return model.ApiConfigSource{
		APIType:        model.ApiType(apiType),
		APITypeStr:     apiConfig.APITypeStr,
		ClusterName:    apiConfig.ClusterName,
		RefreshDelay:   apiConfig.RefreshDelay,
		RequestTimeout: apiConfig.RequestTimeout,
		GrpcServices:   nil, //todo create node of pb
	}
}
