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

package v1alpha1

import (
	"context"
	"fmt"
	"strings"
	"time"
)

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	dubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	istiolog "istio.io/pkg/log"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
)

var log = istiolog.RegisterScope("Snp server", "Snp register server for Proxyless dubbo", 0)

type Snp struct {
	dubbov1alpha1.UnimplementedServiceNameMappingServiceServer

	KubeClient     kube.Client
	queue          chan *RegisterRequest
	debounceAfter  time.Duration
	debounceMax    time.Duration
	enableDebounce bool
}

func NewSnp(kubeClient kube.Client) *Snp {
	return &Snp{
		KubeClient:     kubeClient,
		queue:          make(chan *RegisterRequest, 10),
		debounceAfter:  features.SNPDebounceAfter,
		debounceMax:    features.SNPDebounceMax,
		enableDebounce: features.SNPEnableDebounce,
	}
}

// RegisterServiceAppMapping register service app mapping.
func (s *Snp) RegisterServiceAppMapping(ctx context.Context, req *dubbov1alpha1.ServiceMappingRequest) (*dubbov1alpha1.ServiceMappingResponse, error) {
	namespace := req.GetNamespace()
	interfaces := req.GetInterfaceNames()
	applicationName := req.GetApplicationName()

	registerReq := &RegisterRequest{ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{}}
	for _, interfaceName := range interfaces {
		key := model.ConfigKey{
			Name:      interfaceName,
			Namespace: namespace,
		}
		if _, ok := registerReq.ConfigsUpdated[key]; !ok {
			registerReq.ConfigsUpdated[key] = make(map[string]struct{})
		}
		registerReq.ConfigsUpdated[key][applicationName] = struct{}{}
	}
	s.queue <- registerReq

	return &dubbov1alpha1.ServiceMappingResponse{}, nil
}

func (s *Snp) Register(server *grpc.Server) {
	dubbov1alpha1.RegisterServiceNameMappingServiceServer(server, s)
}

func (s *Snp) Start(stop <-chan struct{}) {
	if s == nil {
		log.Warn("Snp server is not init, skip start")
		return
	}
	go s.debounce(stop, s.push)
}

func (s *Snp) push(req *RegisterRequest) {
	for key, m := range req.ConfigsUpdated {
		var appNames []string
		for app := range m {
			appNames = append(appNames, app)
		}
		for i := 0; i < 3; i++ {
			if err := tryRegister(s.KubeClient, key.Namespace, key.Name, appNames); err != nil {
				log.Errorf(" register [%v] failed: %v, try again later", key, err)
			} else {
				break
			}
		}
	}
}

// Each application registers its services with the Snp server at startup,
// and there are usually multiple copies of a deployment. They make requests
// at the same time. So we need to debounce these requests,
// because only one request is valid for the same deployment!
// So we wait for a while, merge incoming requests, and submit them when
// there is a timeout or a short period of no subsequent requests,
// by default a minimum of 200ms and a maximum of 1s, which is usually acceptable.
// The above can significantly reduce the pressure on the registry.
func (s *Snp) debounce(stopCh <-chan struct{}, pushFn func(req *RegisterRequest)) {
	ch := s.queue
	var timeChan <-chan time.Time
	var startDebounce time.Time
	var lastConfigUpdateTime time.Time

	pushCounter := 0
	debouncedEvents := 0

	var req *RegisterRequest

	free := true
	freeCh := make(chan struct{}, 1)

	push := func(req *RegisterRequest) {
		pushFn(req)
		freeCh <- struct{}{}
	}

	pushWorker := func() {
		eventDelay := time.Since(startDebounce)
		quietTime := time.Since(lastConfigUpdateTime)
		// it has been too long or quiet enough
		if eventDelay >= s.debounceMax || quietTime >= s.debounceAfter {
			if req != nil {
				pushCounter++

				if req.ConfigsUpdated != nil {
					log.Infof(" Push debounce stable[%d] %d for config %s: %v since last change, %v since last push",
						pushCounter, debouncedEvents, configsUpdated(req),
						quietTime, eventDelay)
				}
				free = false
				go push(req)
				req = nil
				debouncedEvents = 0
			}
		} else {
			timeChan = time.After(s.debounceAfter - quietTime)
		}
	}

	for {
		select {
		case <-freeCh:
			free = true
			pushWorker()
		case r := <-ch:
			// if debounce is disabled, push immediately
			if !s.enableDebounce {
				go push(r)
				req = nil
				continue
			}

			lastConfigUpdateTime = time.Now()
			if debouncedEvents == 0 {
				timeChan = time.After(200 * time.Millisecond)
				startDebounce = lastConfigUpdateTime
			}
			debouncedEvents++

			req = req.Merge(r)
		case <-timeChan:
			if free {
				pushWorker()
			}
		case <-stopCh:
			return
		}
	}
}

func tryRegister(kubeClient kube.Client, namespace, interfaceName string, newApps []string) error {
	log.Debugf("try register [%s] in namespace [%s] with [%v] apps", interfaceName, namespace, len(newApps))
	snp, created, err := getOrCreateSnp(kubeClient, namespace, interfaceName, newApps)
	if created {
		log.Debugf("register success, revision:%s", snp.ResourceVersion)
		return nil
	}
	if err != nil {
		return err
	}

	previousLen := len(snp.Spec.ApplicationNames)
	previousAppNames := make(map[string]struct{}, previousLen)
	for _, name := range snp.Spec.ApplicationNames {
		previousAppNames[name] = struct{}{}
	}
	for _, newApp := range newApps {
		previousAppNames[newApp] = struct{}{}
	}
	if len(previousAppNames) == previousLen {
		log.Debugf("[%s] has been registered: %v", interfaceName, newApps)
		return nil
	}

	mergedApps := make([]string, 0, len(previousAppNames))
	for name := range previousAppNames {
		mergedApps = append(mergedApps, name)
	}
	snp.Spec.ApplicationNames = mergedApps
	snpInterface := kubeClient.Istio().ExtensionsV1alpha1().ServiceNameMappings(namespace)
	snp, err = snpInterface.Update(context.Background(), snp, v1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(err, " update failed")
	}
	log.Debugf("register update success, revision:%s", snp.ResourceVersion)
	return nil
}

func getOrCreateSnp(kubeClient kube.Client, namespace string, interfaceName string, newApps []string) (*v1alpha1.ServiceNameMapping, bool, error) {
	ctx := context.TODO()
	// snp name is a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-'
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
					ApplicationNames: newApps,
				},
			}, v1.CreateOptions{})
			// create success
			if err == nil {
				log.Debugf("create snp %s revision %s", interfaceName, snp.ResourceVersion)
				return snp, true, nil
			}
			// If the creation fails, meaning it already created by other goroutine, then get it
			if apierror.IsAlreadyExists(err) {
				log.Debugf("[%s] has been exists, err: %v", err)
				snp, err = snpInterface.Get(ctx, lowerCaseName, v1.GetOptions{})
				// maybe failed to get snp cause of network issue, just return error
				if err != nil {
					return nil, false, errors.Wrap(err, "tryRegister retry get snp error")
				}
			}
		} else {
			return nil, false, errors.Wrap(err, "tryRegister get snp error")
		}
	}
	return snp, false, nil
}

type RegisterRequest struct {
	ConfigsUpdated map[model.ConfigKey]map[string]struct{}
}

func (r *RegisterRequest) Merge(req *RegisterRequest) *RegisterRequest {
	if r == nil {
		return req
	}
	for key, newApps := range req.ConfigsUpdated {
		if _, ok := r.ConfigsUpdated[key]; !ok {
			r.ConfigsUpdated[key] = make(map[string]struct{})
		}
		for app, _ := range newApps {
			r.ConfigsUpdated[key][app] = struct{}{}
		}
	}
	return r
}

func configsUpdated(req *RegisterRequest) string {
	configs := ""
	for key := range req.ConfigsUpdated {
		configs += key.String()
		break
	}
	if len(req.ConfigsUpdated) > 1 {
		more := fmt.Sprintf(" and %d more configs", len(req.ConfigsUpdated)-1)
		configs += more
	}
	return configs
}
