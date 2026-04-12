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

package remote

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

func TestMatchClientRoutesHTTPToHTTPClient(t *testing.T) {
	filter := &Filter{conf: config{DubboProxyConfig: &dubbo.DubboProxyConfig{}}}

	cli, err := filter.matchClient(constant.HTTPRequest)
	require.NoError(t, err)
	assert.Same(t, clienthttp.SingletonHTTPClient(), cli)
}

func TestMatchClientRoutesDubboAndTripleToDubboClient(t *testing.T) {
	filter := &Filter{conf: config{DubboProxyConfig: &dubbo.DubboProxyConfig{}}}

	for _, requestType := range []string{constant.DubboRequest, "triple"} {
		t.Run(requestType, func(t *testing.T) {
			cli, err := filter.matchClient(requestType)
			require.NoError(t, err)
			assert.Same(t, dubbo.SingletonDubboClient(), cli)
		})
	}
}

func TestMatchClientRejectsUnknownRequestType(t *testing.T) {
	filter := &Filter{conf: config{DubboProxyConfig: &dubbo.DubboProxyConfig{}}}

	cli, err := filter.matchClient("grpc")
	assert.Nil(t, cli)
	assert.EqualError(t, err, "not support")
}

func TestFilterFactoryApplyOnlyInitializesDubboClient(t *testing.T) {
	originalDubboInit := initDubboClient
	t.Cleanup(func() {
		initDubboClient = originalDubboInit
	})

	dubboInitCalls := 0
	initDubboClient = func(conf *dubbo.DubboProxyConfig) {
		dubboInitCalls++
	}

	factory := &FilterFactory{
		conf: &config{
			DubboProxyConfig: &dubbo.DubboProxyConfig{},
		},
	}

	require.NoError(t, factory.Apply())
	assert.Equal(t, 1, dubboInitCalls)
}
