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
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
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

func TestAbsolutePath(t *testing.T) {
	rg := &Group{
		basePath: "/",
	}
	rst := rg.absolutePath("abc")
	assert.Equal(t, rst, "/abc")

	rst = rg.absolutePath("")
	assert.Equal(t, rst, "/")

	rg = &Group{
		basePath: "/a",
	}

	rst = rg.absolutePath("")
	assert.Equal(t, rst, "/a")

	rst = rg.absolutePath("b")
	assert.Equal(t, rst, "/a/b")

	rst = rg.absolutePath("/b/")
	assert.Equal(t, rst, "/a/b")

	rst = rg.absolutePath("/:id")
	assert.Equal(t, rst, "/a/:id")
}

func TestGroup(t *testing.T) {
	rg := &Group{
		basePath: "/",
		root:     true,
	}
	rg1, err := rg.Group("")
	assert.Nil(t, rg1)
	assert.Error(t, err, "Cannot group router with empty path")

	rg2, err := rg.Group("/test")
	assert.Equal(t, rg2.basePath, "/test")
	assert.Nil(t, err)

	rg3, err := rg2.Group("mock")
	assert.Nil(t, rg3)
	assert.Error(t, err, "Path must start with '/'")

	rg4, err := rg2.Group("/mock")
	assert.Nil(t, err)
	assert.Equal(t, rg4.basePath, "/test/mock")
}

func TestAdd(t *testing.T) {
	rg := NewRouter()
	rg.Add("/", getMockMethod(config.MethodGet))
	_, ok := rg.routerTree.Get("/")
	assert.True(t, ok)

	// rg.Add("/")
}
