package servicenamemapping

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type Snp struct {
	UnimplementedServiceNameMappingServiceServer

	KubeClient kube.Client
}

func (snp *Snp) RegisterServiceAppMapping(ctx context.Context, req *ServiceMappingRequest) (*ServiceMappingResponse, error) {
	name := req.GetInterfaceName()
	applicationName := req.GetApplicationName()
	create, err := snp.KubeClient.Istio().ExtensionsV1alpha1().ServiceNameMappings(req.Namespace).Create(
		ctx, &v1alpha1.ServiceNameMapping{
			Spec: extensionsv1alpha1.ServiceNameMapping{
				InterfaceName:    name[0],
				ApplicationNames: []string{applicationName},
			},
		},
		v1.CreateOptions{},
	)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(create)
	logger.Errorf("snp RegisterServiceAppMapping %s,%s", string(marshal), err)
	return &ServiceMappingResponse{}, nil
}
