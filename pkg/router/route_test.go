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

package router

import (
	"testing"
)

import (
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-proxy-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-proxy-filter/pkg/router"
)

func getMockMethod(verb config.HTTPVerb) config.Method {
	inbound := config.InboundRequest{}
	integration := config.IntegrationRequest{}
	return config.Method{
		OnAir:              true,
		HTTPVerb:           verb,
		InboundRequest:     inbound,
		IntegrationRequest: integration,
	}
}

func TestPut(t *testing.T) {
	rt := &Route{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
	n0 := getMockMethod(config.MethodGet)
	rt.PutAPI(router.API{URLPattern: "/", Method: n0})
	_, ok := rt.tree.Get("/")
	assert.True(t, ok)

	err := rt.PutAPI(router.API{URLPattern: "/", Method: n0})
	assert.Error(t, err, "Method GET already exists in path /")

	n1 := getMockMethod(config.MethodPost)
	err = rt.PutAPI(router.API{URLPattern: "/mock", Method: n0})
	assert.Nil(t, err)
	err = rt.PutAPI(router.API{URLPattern: "/mock", Method: n1})
	assert.Nil(t, err)
	mNode, ok := rt.tree.Get("/mock")
	assert.True(t, ok)
	assert.Equal(t, len(mNode.(*Node).methods), 2)

	err = rt.PutAPI(router.API{URLPattern: "/mock/test", Method: n0})
	assert.Nil(t, err)
	_, ok = rt.tree.Get("/mock/test")
	assert.True(t, ok)

	rt.PutAPI(router.API{URLPattern: "/test/:id", Method: n0})
	tNode, ok := rt.tree.Get("/test/:id")
	assert.True(t, ok)
	assert.True(t, tNode.(*Node).wildcard)

	err = rt.PutAPI(router.API{URLPattern: "/test/:id", Method: n1})
	assert.Nil(t, err)
	err = rt.PutAPI(router.API{URLPattern: "/test/js", Method: n0})
	assert.Error(t, err, "/test/:id wildcard already exist so that cannot add path /test/js")

	err = rt.PutAPI(router.API{URLPattern: "/test/:id/mock", Method: n0})
	tNode, ok = rt.tree.Get("/test/:id/mock")
	assert.True(t, ok)
	assert.True(t, tNode.(*Node).wildcard)
	assert.Nil(t, err)
}

func TestFindMethod(t *testing.T) {
	rt := &Route{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
	n0 := getMockMethod(config.MethodGet)
	n1 := getMockMethod(config.MethodPost)
	e := rt.PutAPI(router.API{URLPattern: "/theboys", Method: n0})
	assert.Nil(t, e)
	e = rt.PutAPI(router.API{URLPattern: "/theboys/:id", Method: n0})
	assert.Nil(t, e)
	e = rt.PutAPI(router.API{URLPattern: "/vought/:id/supe/:name", Method: n1})
	assert.Nil(t, e)

	m, ok := rt.FindAPI("/theboys", config.MethodGet)
	assert.True(t, ok)
	assert.NotNil(t, m)
	assert.Equal(t, m.URLPattern, "/theboys")

	m, ok = rt.FindAPI("/theboys", config.MethodPost)
	assert.False(t, ok)
	assert.Nil(t, m)

	m, ok = rt.FindAPI("/vought/爱国者/supe/startlight", config.MethodPost)
	assert.True(t, ok)
	assert.NotNil(t, m)
	assert.Equal(t, m.URLPattern, "/vought/:id/supe/:name")

	m, ok = rt.FindAPI("/vought/123/supe/startlight", config.MethodPost)
	assert.True(t, ok)
	assert.NotNil(t, m)
	assert.Equal(t, m.URLPattern, "/vought/:id/supe/:name")
}

func TestUpdateMethod(t *testing.T) {
	m0 := getMockMethod(config.MethodGet)
	m1 := getMockMethod(config.MethodGet)
	m0.DubboBackendConfig.Version = "1.0.0"
	m1.DubboBackendConfig.Version = "2.0.0"

	rt := NewRoute()
	rt.PutAPI(router.API{URLPattern: "/marvel", Method: m0})
	m, _ := rt.FindAPI("/marvel", config.MethodGet)
	assert.Equal(t, m.DubboBackendConfig.Version, "1.0.0")
	rt.UpdateAPI(router.API{URLPattern: "/marvel", Method: m1})
	m, ok := rt.FindAPI("/marvel", config.MethodGet)
	assert.True(t, ok)
	assert.Equal(t, m.DubboBackendConfig.Version, "2.0.0")

	rt.PutAPI(router.API{URLPattern: "/theboys/:id", Method: m0})
	m, _ = rt.FindAPI("/theBoys/12345", config.MethodGet)
	assert.Equal(t, m.DubboBackendConfig.Version, "1.0.0")
	rt.UpdateAPI(router.API{URLPattern: "/theBoys/:id", Method: m1})
	m, ok = rt.FindAPI("/theBoys/12345", config.MethodGet)
	assert.True(t, ok)
	assert.Equal(t, m.DubboBackendConfig.Version, "2.0.0")
}

func TestSearchWildcard(t *testing.T) {
	rt := &Route{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
	n0 := getMockMethod(config.MethodGet)
	e := rt.PutAPI(router.API{URLPattern: "/theboys", Method: n0})
	assert.Nil(t, e)
	e = rt.PutAPI(router.API{URLPattern: "/theboys/:id", Method: n0})
	assert.Nil(t, e)
	e = rt.PutAPI(router.API{URLPattern: "/vought/:id/supe/:name", Method: n0})
	assert.Nil(t, e)

	_, ok := rt.searchWildcard("/marvel")
	assert.False(t, ok)
	_, ok = rt.searchWildcard("/theboys/:id/age")
	assert.False(t, ok)
	_, ok = rt.searchWildcard("/theboys/butcher")
	assert.True(t, ok)
	_, ok = rt.searchWildcard("/vought/:id/supe/homelander")
	assert.True(t, ok)
}

func TestWildcardMatch(t *testing.T) {
	vals := wildcardMatch("/vought/:id", "/vought/12345")
	assert.NotNil(t, vals)
	assert.Equal(t, vals.Get("id"), "12345")
	vals = wildcardMatch("/vought/:id", "/vought/125abc")
	assert.NotNil(t, vals)
	assert.Equal(t, vals.Get("id"), "125abc")
	vals = wildcardMatch("/vought/:id", "/vought/1234abcd/status")
	assert.Nil(t, vals)
	vals = wildcardMatch("/voughT/:id/:action", "/Vought/1234abcd/attack")
	assert.NotNil(t, vals)
	assert.Equal(t, vals.Get("id"), "1234abcd")
	assert.Equal(t, vals.Get("action"), "attack")
	vals = wildcardMatch("/voughT/:id/status", "/Vought/1234abcd/status")
	assert.NotNil(t, vals)
	assert.Equal(t, vals.Get("id"), "1234abcd")
}

func TestGetFilters(t *testing.T) {
	rt := NewRoute()
	n0 := getMockMethod(config.MethodGet)
	n1 := getMockMethod(config.MethodPost)
	e := rt.PutAPI(router.API{URLPattern: "/theboys", Method: n0})
	assert.Nil(t, e)
	e = rt.PutAPI(router.API{URLPattern: "/theboys/:id", Method: n0})
	assert.Nil(t, e)
	e = rt.PutAPI(router.API{URLPattern: "/vought/:id/supe/:name", Method: n1})
	assert.Nil(t, e)

}
