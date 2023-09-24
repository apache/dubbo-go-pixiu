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
	"strconv"
	"time"
)

import (
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	dubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	istioextensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
)

type ServiceMetadataServer struct {
	dubbov1alpha1.UnimplementedServiceMetadataServiceServer

	KubeClient kube.Client

	// CommittedUpdates describes the number of configuration updates the discovery server has
	// received, process, and stored in the push context. If this number is less than InboundUpdates,
	// there are updates we have not yet processed.
	// Note: This does not mean that all proxies have received these configurations; it is strictly
	// the push context, which means that the next push to a proxy will receive this configuration.
	CommittedUpdates *atomic.Int64

	ch chan *pushRequest

	debounceOptions debounceOptions
	debounceHelper  DebounceHelper
}

type pushRequest struct {
	// application name + revision => metadata info
	ConfigsUpdated map[model.ConfigKey]metadataConfig
}

type metadataConfig struct {
	applicationName string
	revision        string
	metadataInfo    string
	timestamp       time.Time
}

func (pr *pushRequest) Merge(other *pushRequest) *pushRequest {
	if pr == nil {
		return other
	}
	if other == nil {
		return pr
	}

	// Keep the first (older) start time

	// Do not merge when any one is empty
	if len(pr.ConfigsUpdated) == 0 || len(other.ConfigsUpdated) == 0 {
		pr.ConfigsUpdated = nil
	} else {
		for key, value := range other.ConfigsUpdated {
			if prValue, ok := pr.ConfigsUpdated[key]; ok {
				if value.timestamp.After(prValue.timestamp) {
					pr.ConfigsUpdated[key] = value
				}
			} else {
				pr.ConfigsUpdated[key] = value
			}
		}
	}

	return pr
}

func NewServiceMetadataServer(client kube.Client) *ServiceMetadataServer {
	return &ServiceMetadataServer{
		CommittedUpdates: atomic.NewInt64(0),
		KubeClient:       client,
		ch:               make(chan *pushRequest, 10),
		debounceOptions: debounceOptions{
			debounceAfter:  features.SMDebounceAfter,
			debounceMax:    features.SMDebounceMax,
			enableDebounce: features.SMEnableDebounce,
		},
	}
}

func (s *ServiceMetadataServer) Start(stopCh <-chan struct{}) {
	if s == nil {
		log.Warn("ServiceMetadataServer is not init, skip start")
		return
	}
	go s.handleUpdate(stopCh)
	go s.removeOutdatedCRD(stopCh)
}

func (s *ServiceMetadataServer) Register(rpcs *grpc.Server) {
	// Register v3 server
	dubbov1alpha1.RegisterServiceMetadataServiceServer(rpcs, s)
}

func (s *ServiceMetadataServer) Publish(ctx context.Context, request *dubbov1alpha1.PublishServiceMetadataRequest) (*dubbov1alpha1.PublishServiceMetadataResponse, error) {
	pushReq := &pushRequest{ConfigsUpdated: map[model.ConfigKey]metadataConfig{}}

	key := model.ConfigKey{
		Name:      getConfigKeyName(request.GetApplicationName(), request.GetRevision()),
		Namespace: request.GetNamespace(),
	}

	pushReq.ConfigsUpdated[key] = metadataConfig{
		applicationName: request.GetApplicationName(),
		revision:        request.GetRevision(),
		metadataInfo:    request.GetMetadataInfo(),
		timestamp:       time.Now(),
	}

	s.ch <- pushReq

	return &dubbov1alpha1.PublishServiceMetadataResponse{}, nil
}

func getConfigKeyName(applicationName, revision string) string {
	log.Infof("application name: %s, revision: %s", applicationName, revision)
	return fmt.Sprintf("%s-%s", applicationName, revision)
}

func (s *ServiceMetadataServer) Get(ctx context.Context, request *dubbov1alpha1.GetServiceMetadataRequest) (*dubbov1alpha1.GetServiceMetadataResponse, error) {
	metadata, err := s.KubeClient.Istio().ExtensionsV1alpha1().ServiceMetadatas(request.Namespace).Get(ctx, getConfigKeyName(request.GetApplicationName(), request.GetRevision()), v1.GetOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return &dubbov1alpha1.GetServiceMetadataResponse{MetadataInfo: metadata.Spec.GetMetadataInfo()}, nil
}

func (s *ServiceMetadataServer) handleUpdate(stopCh <-chan struct{}) {
	s.debounce(stopCh)
}

func (s *ServiceMetadataServer) Push(req *pushRequest) {
	ctx := context.TODO()

	for key, config := range req.ConfigsUpdated {
		metadata, err := getOrCreateCRD(s.KubeClient, ctx, key.Namespace, config.applicationName, config.revision, config.metadataInfo, config.timestamp)
		if err != nil {
			log.Errorf("Failed to getOrCreateCRD, %s", err.Error())
			return
		}

		metadata.Spec.MetadataInfo = config.metadataInfo

		_, err = s.KubeClient.Istio().ExtensionsV1alpha1().ServiceMetadatas(metadata.Namespace).Update(ctx, metadata, v1.UpdateOptions{})
		if err != nil {
			log.Errorf("Failed to update metadata, %s", err.Error())
			return
		}
	}
}

func getOrCreateCRD(kubeClient kube.Client, ctx context.Context, namespace string, applicationName, revision, metadataInfo string, createTime time.Time) (*v1alpha1.ServiceMetadata, error) {
	var (
		metadata *v1alpha1.ServiceMetadata
	)

	serviceMetadataInterface := kubeClient.Istio().ExtensionsV1alpha1().ServiceMetadatas(namespace)
	crdName := getConfigKeyName(applicationName, revision)
	metadata, err := serviceMetadataInterface.Get(ctx, crdName, v1.GetOptions{})

	log.Infof("metadata: %+v", metadata)

	if err != nil {
		if apierror.IsNotFound(err) {
			mt := &v1alpha1.ServiceMetadata{
				ObjectMeta: v1.ObjectMeta{
					Name:      crdName,
					Namespace: namespace,
					Annotations: map[string]string{
						"updateTime": strconv.FormatInt(createTime.Unix(), 10),
					},
				},
				Spec: istioextensionsv1alpha1.ServiceMetadata{
					ApplicationName: applicationName,
					Revision:        revision,
					MetadataInfo:    metadataInfo,
				},
			}

			metadata, err = serviceMetadataInterface.Create(ctx, mt, v1.CreateOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "Failed to create metadata")
			}
		} else {
			return nil, errors.Wrap(err, "Failed to check metadata exists or not")
		}
	} else {
		// if metadata exist in crd, then update timestamp
		if metadata.ObjectMeta.Annotations != nil {
			metadata.ObjectMeta.Annotations["updateTime"] = strconv.FormatInt(createTime.Unix(), 10)
		} else {
			metadata.ObjectMeta.Annotations = map[string]string{
				"updateTime": strconv.FormatInt(createTime.Unix(), 10),
			}
		}
	}

	return metadata, nil
}

func (s *ServiceMetadataServer) debounce(stopCh <-chan struct{}) {
	s.debounceHelper.Debounce(s.ch, stopCh, s.debounceOptions, s.Push, s.CommittedUpdates)
}

func (s *ServiceMetadataServer) removeOutdatedCRD(stopChan <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.TODO())
	ticker := time.NewTicker(2 * 60 * time.Hour)

	for {
		select {
		case <-stopChan:
			cancel()
			ticker.Stop()
			return
		case <-ticker.C:
			namespaces, err := s.KubeClient.Kube().CoreV1().Namespaces().List(ctx, v1.ListOptions{})
			if err != nil {
				log.Errorf("Failed to get namespaces, %s", err.Error())
				continue
			}

			for _, namespace := range namespaces.Items {
				metadataList, err := s.KubeClient.Istio().ExtensionsV1alpha1().ServiceMetadatas(namespace.GetNamespace()).List(ctx, v1.ListOptions{})
				if err != nil {
					log.Errorf("Failed to get metadata with namespace {%s}, %s", namespace, err.Error())
					continue
				}

				for _, metadata := range metadataList.Items {
					if metadata.Annotations != nil {
						if ts, ok := metadata.Annotations["updateTime"]; ok {
							tsInt, err := strconv.ParseInt(ts, 10, 64)
							if err != nil {
								log.Errorf("Failed to get metadata with namespace {%s}, %s", namespace, err.Error())
							} else {
								if time.Now().Sub(time.Unix(tsInt, 0)) > 24*time.Hour {
									err := s.KubeClient.Istio().ExtensionsV1alpha1().ServiceMetadatas(namespace.GetNamespace()).Delete(ctx, metadata.GetName(), v1.DeleteOptions{})
									if err != nil {
										log.Errorf("Failed to delete outdated metadata with namespace {%s}, metadata {%s}, err {%s}", namespace, metadata.GetName(), err.Error())
									}
								}
							}
						} else {
							log.Errorf("Failed to get metadata with namespace {%s}, metadata {%s}", namespace, metadata.GetName())
						}
					}
				}
			}

		}
	}
}
