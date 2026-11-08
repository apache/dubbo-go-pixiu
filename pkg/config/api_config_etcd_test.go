/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package config

import (
	"testing"
)

func TestEtcdRouteMethodsBelongToResourceID(t *testing.T) {
	resourceKeys := []string{"/pixiu/resources/1", "/pixiu/resources/2"}
	resourceValues := []string{
		"id: 1\npath: /users\ntype: restful\n",
		"id: 2\npath: /users\ntype: restful\n",
	}
	methodKeys := []string{"/pixiu/resources/1/method/1", "/pixiu/resources/2/method/2"}
	methodValues := []string{
		"id: 1\nresourcePath: /users\nhttpVerb: GET\n",
		"id: 2\nresourcePath: /users\nhttpVerb: POST\n",
	}
	loaded := &APIConfig{}
	if err := initAPIConfigServiceFromKvList(loaded, resourceKeys, resourceValues); err != nil {
		t.Fatalf("load resources: %v", err)
	}
	if err := initAPIConfigMethodFromKvList(loaded, methodKeys, methodValues); err != nil {
		t.Fatalf("load methods: %v", err)
	}
	assertSeparateRouteResources(t, loaded)

	previousConfig, previousListener := apiConfig, listener
	t.Cleanup(func() {
		apiConfig, listener = previousConfig, previousListener
	})
	apiConfig = &APIConfig{Resources: []Resource{
		{ID: 1, Path: "/users"},
		{ID: 2, Path: "/users"},
	}}
	listener = noopAPIConfigListener{}
	for index, key := range methodKeys {
		handlePutEvent([]byte(key), []byte(methodValues[index]))
	}
	assertSeparateRouteResources(t, apiConfig)
	handleDeleteEvent([]byte(methodKeys[0]), nil)
	if len(apiConfig.Resources[0].Methods) != 0 || len(apiConfig.Resources[1].Methods) != 1 {
		t.Fatalf("deleting GET affected another resource: %+v", apiConfig.Resources)
	}
}

func assertSeparateRouteResources(t *testing.T, loaded *APIConfig) {
	t.Helper()
	if len(loaded.Resources) != 2 {
		t.Fatalf("resource count: got %d, want 2", len(loaded.Resources))
	}
	for index, verb := range []string{"GET", "POST"} {
		resource := loaded.Resources[index]
		if resource.ID != index+1 || len(resource.Methods) != 1 || resource.Methods[0].HTTPVerb != verb {
			t.Fatalf("resource %d methods: %+v", index+1, resource)
		}
	}
}

type noopAPIConfigListener struct{}

func (noopAPIConfigListener) ResourceChange(Resource, Resource) bool { return true }
func (noopAPIConfigListener) ResourceAdd(Resource) bool              { return true }
func (noopAPIConfigListener) ResourceDelete(Resource) bool           { return true }
func (noopAPIConfigListener) MethodChange(Resource, Method, Method) bool {
	return true
}
func (noopAPIConfigListener) MethodAdd(Resource, Method) bool    { return true }
func (noopAPIConfigListener) MethodDelete(Resource, Method) bool { return true }
