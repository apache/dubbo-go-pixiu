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

package dubbogen

import (
	"strings"
)

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	istioioapidubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	istioioapiextensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/networking/util"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/gvk"
)

// Map of all configs that do not impact LDS
var snpConfigs = map[model.NodeType]map[config.GroupVersionKind]struct{}{
	model.SidecarProxy: {
		gvk.ServiceNameMapping: {},
	},
}

func snpNeedsPush(proxy *model.Proxy, req *model.PushRequest) bool {
	if req == nil {
		return true
	}
	if !req.Full {
		// SNP only handles full push
		return false
	}
	// If none set, we will always push
	if len(req.ConfigsUpdated) == 0 {
		return true
	}
	for conf := range req.ConfigsUpdated {
		if _, f := snpConfigs[proxy.Type][conf.Kind]; f {
			return true
		}
	}
	return false
}

func (d *DubboConfigGenerator) buildServiceNameMappings(node *model.Proxy, req *model.PushRequest, watchedResourceNames []string) model.Resources {
	log.Debugf("building snp for %s", node.ID)
	if !snpNeedsPush(node, req) {
		return nil
	}
	snps := buildSnp(node, req, watchedResourceNames)
	log.Debugf("snp res for %s:\n%v", node.ID, snps)
	resources := model.Resources{}
	for _, c := range snps {
		resources = append(resources, &discovery.Resource{
			Name:     c.InterfaceName,
			Resource: util.MessageToAny(c),
		})
	}
	return resources
}

func buildSnp(node *model.Proxy, req *model.PushRequest, watchedResourceNames []string) []*istioioapidubbov1alpha1.ServiceMappingXdsResponse {
	defaultNs := node.ConfigNamespace
	res := make([]*istioioapidubbov1alpha1.ServiceMappingXdsResponse, 0, len(watchedResourceNames))

	watchedResourceNamesMap := map[string]interface{}{}
	for _, name := range watchedResourceNames {
		watchedResourceNamesMap[name] = struct{}{}
	}
	updatedMap := map[string]interface{}{}
	for update, _ := range req.ConfigsUpdated {
		updatedMap[update.Name] = struct{}{}
	}
	// if req.ConfigsUpdated is empty, it means that all watched resources need to be pushed
	for _, w := range watchedResourceNames {
		// if configsUpdated empty, meaning a full push of watched snp resources.
		if _, exists := updatedMap[w]; exists || len(req.ConfigsUpdated) == 0 {
			namespace, snpName := extractNameAndNameSpace(w, defaultNs)
			if snp := req.Push.ServiceNameMappingsByNameSpaceAndInterfaceName(namespace, snpName); snp != nil {
				mapping := snp.Spec.(*istioioapiextensionsv1alpha1.ServiceNameMapping)
				res = append(res, &istioioapidubbov1alpha1.ServiceMappingXdsResponse{
					Namespace:        snp.Namespace,
					InterfaceName:    mapping.InterfaceName,
					ApplicationNames: mapping.ApplicationNames,
				})
			}
		}
	}
	return res
}

func extractNameAndNameSpace(watchedResource, defaultNamespace string) (string, string) {
	split := strings.Split(watchedResource, "|")
	if len(split) == 1 || split[1] == "" {
		return defaultNamespace, split[1]
	}
	return split[1], split[0]
}
