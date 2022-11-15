package v1alpha1

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/gvk"
	dubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"strings"
)

type Snp struct {
	dubbov1alpha1.UnimplementedServiceNameMappingServiceServer

	RWConfigStore model.ConfigStoreController
}

//TODO: RegisterServiceAppMapping
func (s *Snp) RegisterServiceAppMapping(ctx context.Context, req *dubbov1alpha1.ServiceMappingRequest) (*dubbov1alpha1.ServiceMappingResponse, error) {
	namespace := req.GetNamespace()
	interfaces := req.GetInterfaceNames()
	applicationName := req.GetApplicationName()
	for _, i := range interfaces {
		for j := 0; j < 3; j++ {
			err := tryRegister(s.RWConfigStore, namespace, applicationName, i)
			if err == nil {
				continue
			}
			logger.Errorf("RegisterServiceAppMapping error: %v", err)
		}
	}
	logger.Debugf("s RegisterServiceAppMapping finished, application: %s,interfaces:%v", applicationName, interfaces)
	return &dubbov1alpha1.ServiceMappingResponse{}, nil
}

func tryRegister(configStore model.ConfigStoreController, namespace string, applicationName string, interfaceName string) error {
	lowerCaseName := strings.ToLower(strings.ReplaceAll(interfaceName, ".", "-"))
	snp := configStore.Get(gvk.ServiceNameMapping, lowerCaseName, namespace)
	if snp == nil {
		revision, err := configStore.Create(config.Config{
			Meta: config.Meta{
				Name:             lowerCaseName,
				Namespace:        namespace,
				GroupVersionKind: gvk.ServiceNameMapping,
				Labels: map[string]string{
					"interface": interfaceName,
				},
			},
			Spec: &extensionsv1alpha1.ServiceNameMapping{
				InterfaceName:    interfaceName,
				ApplicationNames: []string{applicationName},
			},
		})
		if err == nil {
			logger.Infof("create snp %s revision %s", interfaceName, revision)
			return nil
		}
		// If the creation fails, meaning it already exists and needs to be updated
		if errors.IsAlreadyExists(err) {
			logger.Infof("Snp[%s] has been exists, err: %v", err)
			snp = configStore.Get(gvk.ServiceNameMapping, lowerCaseName, namespace)
		}
	}
	mapping := snp.Spec.(*extensionsv1alpha1.ServiceNameMapping)
	for _, name := range mapping.ApplicationNames {
		if name == applicationName {
			return nil
		}
	}
	mapping.ApplicationNames = append(mapping.ApplicationNames, applicationName)
	revision, err := configStore.Update(*snp)
	if err != nil {
		logger.Errorf("s RegisterServiceAppMapping revision:%s", err)
		return err
	}
	logger.Debugf("s RegisterServiceAppMapping revision:%s", revision)
	return nil
}
