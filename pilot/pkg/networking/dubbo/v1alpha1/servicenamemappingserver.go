package v1alpha1

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/gvk"
	dubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	"strings"
)

type Snp struct {
	dubbov1alpha1.UnimplementedServiceNameMappingServiceServer

	RWConfigStore model.ConfigStoreController
}

//TODO: RegisterServiceAppMapping
func (snp *Snp) RegisterServiceAppMapping(ctx context.Context, req *dubbov1alpha1.ServiceMappingRequest) (*dubbov1alpha1.ServiceMappingResponse, error) {
	namespace := req.GetNamespace()
	interfaces := req.GetInterfaceNames()
	applicationName := req.GetApplicationName()
	for _, i := range interfaces {
		lowerCaseName := strings.ToLower(strings.ReplaceAll(i, ".", "-"))
		patch, err := snp.RWConfigStore.Patch(
			config.Config{
				Meta: config.Meta{
					Name:             lowerCaseName,
					Namespace:        namespace,
					GroupVersionKind: gvk.ServiceNameMapping,
					Labels: map[string]string{
						"interface": i,
					},
				},
				Spec: &extensionsv1alpha1.ServiceNameMapping{
					InterfaceName:    i,
					ApplicationNames: []string{applicationName},
				},
			}, func(cfg config.Config) (config.Config, kubetypes.PatchType) {
				mapping := cfg.Spec.(*extensionsv1alpha1.ServiceNameMapping)
				mapping.ApplicationNames = append(mapping.ApplicationNames, applicationName)
				return cfg, kubetypes.MergePatchType
			})
		if err != nil {
			logger.Errorf("snp RegisterServiceAppMapping %s,%s", patch, err)
		}
	}
	logger.Debugf("snp RegisterServiceAppMapping finished, application: %s,interfaces:%v", applicationName, interfaces)
	return &dubbov1alpha1.ServiceMappingResponse{}, nil
}
