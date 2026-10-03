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

package registry

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

func TestParseDubboStringPreservesSerialization(t *testing.T) {
	backend, methods, location, err := ParseDubboString("tri://127.0.0.1:20001/org.apache.dubbogo.samples.api.Greeter?application=BDTService&interface=org.apache.dubbogo.samples.api.Greeter&methods=SayHello&serialization=hessian2")
	require.NoError(t, err)

	assert.Equal(t, "tri", backend.Protocol)
	assert.Equal(t, "hessian2", backend.Serialization)
	assert.Equal(t, []string{"SayHello"}, methods)
	assert.Equal(t, "127.0.0.1:20001", location)
}

func TestCreateAPIConfigDefaultsGenericSerialization(t *testing.T) {
	api := CreateAPIConfig(
		"/dubbo.io/org.apache.dubbo.sample.UserProvider",
		"127.0.0.1:20000",
		config.DubboBackendConfig{
			ApplicationName: "dubbo.io",
			Protocol:        "dubbo",
			Interface:       "org.apache.dubbo.sample.UserProvider",
		},
		"GetUser",
		nil,
	)

	assert.Equal(t, "hessian2", api.Method.IntegrationRequest.Serialization)
}
