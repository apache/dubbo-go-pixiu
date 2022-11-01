package dubbogen

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/networking/util"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/gvk"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	istioioapiextensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
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

func buildSnp(node *model.Proxy, req *model.PushRequest, watchedResourceNames []string) []*istioioapiextensionsv1alpha1.ServiceNameMapping {
	namespace := node.ConfigNamespace
	res := make([]*istioioapiextensionsv1alpha1.ServiceNameMapping, len(watchedResourceNames))

	watchedResourceNamesMap := map[string]interface{}{}
	for _, name := range watchedResourceNames {
		watchedResourceNamesMap[name] = struct{}{}
	}
	updatedMap := map[string]interface{}{}
	for update, _ := range req.ConfigsUpdated {
		if update.Namespace == namespace {
			updatedMap[update.Name] = struct{}{}
		}
	}
	snps := req.Push.ServiceNameMappingsForProxy(namespace)
	for _, snp := range snps {
		mapping := snp.Spec.(*istioioapiextensionsv1alpha1.ServiceNameMapping)
		// if req.ConfigsUpdated, meaning a full push of watched snp resources.
		if _, exists := updatedMap[mapping.InterfaceName]; exists || len(req.ConfigsUpdated) == 0 {
			// filter watched resource
			if _, exists := watchedResourceNamesMap[mapping.InterfaceName]; exists {
				res = append(res, mapping)
			}
		}
	}
	return res
}
