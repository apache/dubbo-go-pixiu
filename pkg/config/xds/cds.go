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
	"fmt"
	"strconv"
)

import (
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

type xdsClusterReplacer interface {
	ReplaceXDSClusters(clusters []*model.ClusterConfig) error
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

	if len(delta.NewResources) == 0 && !containsResource(delta.RemovedResources, constant.ClusterType) {
		return nil
	}
	return c.setupCluster(clusters)
}

func (c *CdsManager) setupCluster(clusters []*xdsmodel.Cluster) error {
	converted := make([]*model.ClusterConfig, 0, len(clusters))
	names := make(map[string]struct{}, len(clusters))
	for _, cluster := range clusters {
		if cluster == nil || cluster.Name == "" {
			return errors.New("xDS cluster must have a name")
		}
		if _, duplicate := names[cluster.Name]; duplicate {
			return errors.Errorf("duplicate xDS cluster %q", cluster.Name)
		}
		names[cluster.Name] = struct{}{}
		convertedCluster, err := c.makeCluster(cluster)
		if err != nil {
			return errors.Wrapf(err, "xDS cluster %q", cluster.Name)
		}
		converted = append(converted, convertedCluster)
	}
	replacer, ok := c.clusterMg.(xdsClusterReplacer)
	if !ok {
		return errors.New("cluster manager does not support transactional xDS replacement")
	}
	return errors.Wrap(replacer.ReplaceXDSClusters(converted), "can not replace xDS clusters")
}

func (c *CdsManager) makeCluster(cluster *xdsmodel.Cluster) (*model.ClusterConfig, error) {
	clusterType, err := c.makeClusterType(cluster)
	if err != nil {
		return nil, err
	}
	lb, err := c.makeLoadBalancePolicy(cluster.LbStr)
	if err != nil {
		return nil, err
	}
	endpoints, err := c.makeEndpoints(cluster.Endpoints)
	if err != nil {
		return nil, err
	}
	edsConfig, err := c.makeEdsClusterConfig(cluster.EdsClusterConfig)
	if err != nil {
		return nil, err
	}
	return &model.ClusterConfig{
		Name:             cluster.Name,
		TypeStr:          cluster.TypeStr,
		Type:             clusterType,
		EdsClusterConfig: edsConfig,
		LbStr:            lb,
		HealthChecks:     c.makeHealthChecks(cluster.HealthChecks),
		Endpoints:        endpoints,
	}, nil
}

func (c *CdsManager) makeLoadBalancePolicy(lb string) (model.LbPolicyType, error) {
	if lb == "" {
		return model.LoadBalancerRand, nil
	}
	policy, ok := model.LbPolicyTypeValue[lb]
	if !ok {
		return "", errors.Errorf("unsupported load-balancing policy %q", lb)
	}
	return policy, nil
}

func (c *CdsManager) makeClusterType(cluster *xdsmodel.Cluster) (model.DiscoveryType, error) {
	clusterType, ok := model.DiscoveryTypeValue[cluster.TypeStr]
	if !ok {
		return 0, errors.Errorf("unsupported discovery type %q", cluster.TypeStr)
	}
	return clusterType, nil
}

func (c *CdsManager) makeEndpoints(endpoints []*xdsmodel.Endpoint) ([]*model.Endpoint, error) {
	r := make([]*model.Endpoint, 0, len(endpoints))
	ids := make(map[string]struct{}, len(endpoints))
	addresses := make(map[string]struct{}, len(endpoints))
	for index, endpoint := range endpoints {
		if endpoint == nil {
			return nil, errors.Errorf("endpoint %d is nil", index)
		}
		if endpoint.Id == "" {
			return nil, errors.Errorf("endpoint %d has an empty ID", index)
		}
		if _, duplicate := ids[endpoint.Id]; duplicate {
			return nil, errors.Errorf("duplicate endpoint ID %q", endpoint.Id)
		}
		ids[endpoint.Id] = struct{}{}
		if endpoint.Address == nil || endpoint.Address.Address == "" || endpoint.Address.Port <= 0 || endpoint.Address.Port > 65535 {
			return nil, errors.Errorf("endpoint %q has invalid socket address", endpoint.Id)
		}
		addressKey := fmt.Sprintf("%s:%d", endpoint.Address.Address, endpoint.Address.Port)
		if _, duplicate := addresses[addressKey]; duplicate {
			return nil, errors.Errorf("duplicate endpoint socket address %q", addressKey)
		}
		addresses[addressKey] = struct{}{}
		unhealthy := false
		if value, ok := endpoint.Metadata[endpointHealthMetadataKey]; ok {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.Wrapf(err, "endpoint %q has invalid %s metadata", endpoint.Id, endpointHealthMetadataKey)
			}
			unhealthy = parsed
		}
		r = append(r, &model.Endpoint{
			ID:        endpoint.Id,
			Name:      endpoint.Name,
			Address:   c.makeAddress(endpoint),
			Metadata:  endpoint.Metadata,
			UnHealthy: unhealthy,
		})
	}
	return r, nil
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

func (c *CdsManager) makeEdsClusterConfig(edsConfig *xdsmodel.EdsClusterConfig) (model.EdsClusterConfig, error) {
	if edsConfig == nil {
		return model.EdsClusterConfig{}, nil
	}
	result := model.EdsClusterConfig{ServiceName: edsConfig.ServiceName}
	if configSource := edsConfig.GetEdsConfig(); configSource != nil {
		apiConfig, err := c.makeApiConfigSource(configSource.GetApiConfigSource())
		if err != nil {
			return model.EdsClusterConfig{}, err
		}
		result.EdsConfig = model.ConfigSource{
			Path:            configSource.GetPath(),
			ApiConfigSource: apiConfig,
		}
	}
	return result, nil
}

func (c *CdsManager) makeApiConfigSource(apiConfig *xdsmodel.ApiConfigSource) (result model.ApiConfigSource, err error) {
	if apiConfig == nil {
		return result, nil
	}
	apiType, ok := model.ApiTypeValue[apiConfig.APITypeStr]
	if !ok {
		return result, errors.Errorf("unsupported EDS API type %q", apiConfig.APITypeStr)
	}

	return model.ApiConfigSource{
		APIType:        model.ApiType(apiType),
		APITypeStr:     apiConfig.APITypeStr,
		ClusterName:    apiConfig.ClusterName,
		RefreshDelay:   apiConfig.RefreshDelay,
		RequestTimeout: apiConfig.RequestTimeout,
		GrpcServices:   nil, //todo create node of pb
	}, nil
}

func containsResource(resources []string, target string) bool {
	for _, resourceName := range resources {
		if resourceName == target {
			return true
		}
	}
	return false
}
