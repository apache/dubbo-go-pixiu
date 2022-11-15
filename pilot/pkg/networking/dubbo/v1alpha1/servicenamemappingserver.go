package v1alpha1

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	"github.com/pkg/errors"
	dubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	istiolog "istio.io/pkg/log"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var log = istiolog.RegisterScope("Snp server", "Snp register server for Proxyless dubbo", 0)

type Snp struct {
	dubbov1alpha1.UnimplementedServiceNameMappingServiceServer

	KubeClient kube.Client
}

//TODO: RegisterServiceAppMapping
func (s *Snp) RegisterServiceAppMapping(ctx context.Context, req *dubbov1alpha1.ServiceMappingRequest) (*dubbov1alpha1.ServiceMappingResponse, error) {
	namespace := req.GetNamespace()
	interfaces := req.GetInterfaceNames()
	applicationName := req.GetApplicationName()
	for _, i := range interfaces {
		for j := 0; j < 3; j++ {
			err := tryRegister(ctx, s.KubeClient, namespace, applicationName, i)
			if err == nil {
				break
			}
			log.Errorf("RegisterServiceAppMapping error: %v", err)
		}
	}
	log.Debugf("snp RegisterServiceAppMapping finished, application: %s,interfaces:%v", applicationName, interfaces)
	return &dubbov1alpha1.ServiceMappingResponse{}, nil
}

func tryRegister(ctx context.Context, kubeClient kube.Client, namespace string, applicationName string, interfaceName string) error {
	lowerCaseName := strings.ToLower(strings.ReplaceAll(interfaceName, ".", "-"))
	snpInterface := kubeClient.Istio().ExtensionsV1alpha1().ServiceNameMappings(namespace)
	snp, err := snpInterface.Get(ctx, lowerCaseName, v1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			//try to create snp, this operation may be conflict with other goroutine
			snp, err = snpInterface.Create(ctx, &v1alpha1.ServiceNameMapping{
				ObjectMeta: v1.ObjectMeta{
					Name:      lowerCaseName,
					Namespace: namespace,
					Labels: map[string]string{
						"interface": interfaceName,
					},
				},
				Spec: extensionsv1alpha1.ServiceNameMapping{
					InterfaceName:    interfaceName,
					ApplicationNames: []string{applicationName},
				},
			}, v1.CreateOptions{})
			// create success
			if err == nil {
				log.Infof("create snp %s revision %s", interfaceName, snp.ResourceVersion)
				return nil
			}
			// If the creation fails, meaning it already created by other goroutine, then get it
			if apierror.IsAlreadyExists(err) {
				log.Infof("Snp[%s] has been exists, err: %v", err)
				snp, err = snpInterface.Get(ctx, lowerCaseName, v1.GetOptions{})
				// maybe failed to get snp cause of network issue, just return error
				if err != nil {
					return errors.Wrap(err, "tryRegister retry get snp error")
				}
			}
		} else {
			return errors.Wrap(err, "tryRegister get snp error")
		}
	}
	for _, name := range snp.Spec.ApplicationNames {
		if name == applicationName {
			return nil
		}
	}
	snp.Spec.ApplicationNames = append(snp.Spec.ApplicationNames, applicationName)
	snp, err = snpInterface.Update(ctx, snp, v1.UpdateOptions{})
	if err != nil {
		log.Errorf("s RegisterServiceAppMapping revision:%s", err)
		return errors.Wrap(err, "tryRegister retry update snp error")
	}
	log.Debugf("s RegisterServiceAppMapping revision:%s", snp.ResourceVersion)
	return nil
}
