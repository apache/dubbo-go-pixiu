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
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	dubbov1alpha1 "istio.io/api/dubbo/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/assert"
)

func TestRegisterServiceAppMapping(t *testing.T) {
	client := kube.NewFakeClient()
	s := NewSnp(client)
	stop := make(chan struct{})

	s.Start(stop)
	w := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		w.Add(1)
		j := i
		go func() {
			_, _ = s.RegisterServiceAppMapping(context.Background(), &dubbov1alpha1.ServiceMappingRequest{
				Namespace:       "default",
				ApplicationName: "app" + strconv.Itoa(j),
				InterfaceNames:  []string{"com.test.TestService"},
			})
			w.Done()
		}()
	}

	w.Wait()
	time.Sleep(s.debounceMax)
	stop <- struct{}{}

	lowerCaseName := strings.ToLower(strings.ReplaceAll("com.test.TestService", ".", "-"))
	mapping, err := client.Istio().ExtensionsV1alpha1().ServiceNameMappings("default").Get(context.Background(), lowerCaseName, v1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(mapping.Spec.ApplicationNames), 100)
}

func TestReqMerge(t *testing.T) {
	tests := []struct {
		name string
		reqs []*RegisterRequest
		want *RegisterRequest
	}{
		{
			name: "different",
			reqs: []*RegisterRequest{
				{
					ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{
						model.ConfigKey{
							Name:      "a-b-c-demoservice",
							Namespace: "default",
						}: {
							"app1": struct{}{},
							"app2": struct{}{},
						},
					},
				},
				{
					ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{
						model.ConfigKey{
							Name:      "a-b-c-demoservice",
							Namespace: "default",
						}: {
							"app3": struct{}{},
						},
					},
				},
			},
			want: &RegisterRequest{
				ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{
					model.ConfigKey{
						Name:      "a-b-c-demoservice",
						Namespace: "default",
					}: {
						"app1": struct{}{},
						"app2": struct{}{},
						"app3": struct{}{},
					},
				},
			},
		},
		{
			name: "same",
			reqs: []*RegisterRequest{
				{
					ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{
						model.ConfigKey{
							Name:      "a-b-c-demoservice",
							Namespace: "default",
						}: {
							"app1": struct{}{},
						},
					},
				},
				{
					ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{
						model.ConfigKey{
							Name:      "a-b-c-demoservice",
							Namespace: "default",
						}: {
							"app1": struct{}{},
						},
					},
				},
			},
			want: &RegisterRequest{
				ConfigsUpdated: map[model.ConfigKey]map[string]struct{}{
					model.ConfigKey{
						Name:      "a-b-c-demoservice",
						Namespace: "default",
					}: {
						"app1": struct{}{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		var got *RegisterRequest
		t.Run(tt.name, func(t *testing.T) {
			for _, r := range tt.reqs {
				got = got.Merge(r)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
