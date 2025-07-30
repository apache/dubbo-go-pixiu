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

package controller

import (
	"context"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/controller/pkg/proxy"
	"github.com/apache/dubbo-go-pixiu/controller/pkg/utils"
	"github.com/rs/zerolog/log"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"time"
)

type IngressController struct {
	clientSet    *kubernetes.Clientset
	queue        workqueue.RateLimitingInterface
	indexer      cache.Indexer
	informer     cache.Controller
	ingressClass string
}

func NewIngressController(client *kubernetes.Clientset, ingressClass string) *IngressController {
	ListWatcher := cache.NewListWatchFromClient(client.NetworkingV1().RESTClient(), "ingresses", "", fields.Everything())
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	indexer, informer := cache.NewIndexerInformer(ListWatcher, &networkingV1.Ingress{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				log.Info().Str("object", key).Msg("watch.add.Ingress")
				queue.Add(&Event{Type: Add, Obj: key})
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				log.Info().Str("object", key).Msg("watch.update.Ingress")
				queue.Add(&Event{Type: Update, Obj: key})
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				log.Info().Str("object", key).Msg("watch.delete.Ingress")
				queue.Add(&Event{Type: Delete, Obj: key, Tombstone: obj})
			}
		},
	}, cache.Indexers{})

	return &IngressController{
		clientSet:    client,
		queue:        queue,
		indexer:      indexer,
		informer:     informer,
		ingressClass: ingressClass,
	}
}

func (c *IngressController) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	event, ok := item.(*Event)
	if !ok {
		log.Error().Msgf("Expected *Event but got %T", item)
		c.queue.Forget(item)
		return true
	}

	key, ok := event.Obj.(string)
	if !ok {
		log.Error().Msgf("Expected string but got %T for object", event.Obj)
		c.queue.Forget(item)
		return true
	}

	log.Info().Str("eventType", event.Type.String()).Str("key", key).Msg("Processing event")

	err := c.syncToStdout(key)
	c.handleErr(err, item)

	return true
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *IngressController) handleErr(err error, key any) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing ingress %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	klog.Infof("Dropping ingress %q out of the queue: %v", key, err)
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the ingress to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *IngressController) syncToStdout(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return fmt.Errorf("invalid key format: %v", err)
	}

	ingress, err := c.clientSet.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Error().Str("Ingress", key).Err(err).Msg("Failed to fetch ingress from API")
		return err
	}

	if ingress.Spec.IngressClassName == nil {
		log.Info().Str("Ingress", key).Msg("ingressClass no resources in use")
		return nil
	} else if *ingress.Spec.IngressClassName == c.ingressClass {
		log.Info().Str("Ingress", key).Str("IngressClass", *ingress.Spec.IngressClassName).Msg("This resource is required")
		newAddress := networkingV1.IngressLoadBalancerIngress{
			IP: "0.0.0.0",
		}
		ingress.Status.LoadBalancer.Ingress = []networkingV1.IngressLoadBalancerIngress{newAddress}

		updatedIngress, err := c.clientSet.NetworkingV1().Ingresses(namespace).UpdateStatus(context.Background(), ingress, metav1.UpdateOptions{})
		if err != nil {
			log.Error().Str("Ingress", key).Err(err).Msg("Failed to update ingress status")
			return err
		}
		log.Info().Str("Ingress", key).Msg("Successfully updated ingress status")
		log.Debug().Interface("UpdatedIngress", updatedIngress).Msg("Updated ingress status")
	}

	return nil
}

func (c *IngressController) RunWorker() {
	for c.processNextItem() {
	}
}

func (c *IngressController) Run(thread int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	log.Info().Msg("pixiu ingress controller enabled")
	go c.informer.Run(stopCh)

	wg := utils.Group{}
	wg.Go(func() {
		proxy.Render(c.clientSet)
	})

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for ingress caches to sync"))
		return
	}

	for i := 0; i < thread; i++ {
		go wait.Until(c.RunWorker, time.Second, stopCh)
	}
	<-stopCh
	log.Info().Msg("Ingress Controller has been gracefully stopped.")
}
