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
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
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
	rt.Put("/", n0)
	_, ok := rt.tree.Get("/")
	assert.True(t, ok)

	err := rt.Put("/", n0)
	assert.Error(t, err, "Method GET already exists in path /")

	n1 := getMockMethod(config.MethodPost)
	err = rt.Put("/mock", n0)
	assert.Nil(t, err)
	err = rt.Put("/mock", n1)
	assert.Nil(t, err)
	mNode, ok := rt.tree.Get("/mock")
	assert.True(t, ok)
	assert.Equal(t, len(mNode.(*Node).methods), 2)

	err = rt.Put("/mock/test", n0)
	assert.Nil(t, err)
	_, ok = rt.tree.Get("/mock")
	assert.True(t, ok)

	rt.Put("/test/:id", n0)
	tNode, ok := rt.tree.Get("/test/:id")
	assert.True(t, ok)
	assert.True(t, tNode.(*Node).wildcard)

	err = rt.Put("/test/:id", n1)
	assert.Nil(t, err)
	err = rt.Put("/test/js", n0)
	assert.Error(t, err, "/test/:id wildcard already exist so that cannot add path /test/js")

	err = rt.Put("/test/:id/mock", n0)
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
	e := rt.Put("/theboys", n0)
	assert.Nil(t, e)
	e = rt.Put("/theboys/:id", n0)
	assert.Nil(t, e)
	e = rt.Put("/vought/:id/supe/:name", n1)
	assert.Nil(t, e)

	m, ok := rt.FindMethod("/theboys", config.MethodGet)
	assert.True(t, ok)
	assert.NotNil(t, m)

	m, ok = rt.FindMethod("/theboys", config.MethodPost)
	assert.False(t, ok)
	assert.Nil(t, m)

	m, ok = rt.FindMethod("/vought/123/supe/startlight", config.MethodPost)
	assert.True(t, ok)
	assert.NotNil(t, m)
}

func TestUpdateMethod(t *testing.T) {
	m0 := getMockMethod(config.MethodGet)
	m1 := getMockMethod(config.MethodGet)
	m0.Version = "1.0.0"
	m1.Version = "2.0.0"

	rt := NewRoute()
	rt.Put("/marvel", m0)
	m, _ := rt.FindMethod("/marvel", config.MethodGet)
	assert.Equal(t, m.Version, "1.0.0")
	rt.UpdateMethod("/marvel", config.MethodGet, m1)
	m, ok := rt.FindMethod("/marvel", config.MethodGet)
	assert.True(t, ok)
	assert.Equal(t, m.Version, "2.0.0")

	rt.Put("/theboys/:id", m0)
	m, _ = rt.FindMethod("/theBoys/12345", config.MethodGet)
	assert.Equal(t, m.Version, "1.0.0")
	rt.UpdateMethod("/theBoys/:id", config.MethodGet, m1)
	m, ok = rt.FindMethod("/theBoys/12345", config.MethodGet)
	assert.True(t, ok)
	assert.Equal(t, m.Version, "2.0.0")
}

func TestSearchWildcard(t *testing.T) {
	rt := &Route{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
	n0 := getMockMethod(config.MethodGet)
	e := rt.Put("/theboys", n0)
	assert.Nil(t, e)
	e = rt.Put("/theboys/:id", n0)
	assert.Nil(t, e)
	e = rt.Put("/vought/:id/supe/:name", n0)
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

func TestContainParam(t *testing.T) {
	assert.True(t, containParam("/test/:id"))
	assert.False(t, containParam("/test"))
	assert.True(t, containParam("/test/:id/mock"))
}

func TestWildcardMatch(t *testing.T) {
	assert.True(t, wildcardMatch("/vought/:id", "/vought/12345"))
	assert.True(t, wildcardMatch("/vought/:id", "/vought/125abc"))
	assert.False(t, wildcardMatch("/vought/:id", "/vought/1234abcd/status"))
	assert.True(t, wildcardMatch("/voughT/:id/:action", "/Vought/1234abcd/attack"))
}
