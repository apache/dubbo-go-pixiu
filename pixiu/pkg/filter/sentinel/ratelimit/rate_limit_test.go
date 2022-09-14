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

package ratelimit

import (
	"bytes"
	stdHttp "net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/mock"
)

func TestFilter(t *testing.T) {
	f := &FilterFactory{conf: &Config{}}

	mockYaml, err := yaml.MarshalYML(mockConfig())
	assert.Nil(t, err)

	assert.Nil(t, yaml.UnmarshalYML(mockYaml, f.Config()))

	assert.Nil(t, f.Apply())

	decoder := &Filter{conf: f.conf, matcher: f.matcher}

	request, _ := stdHttp.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	c := mock.GetMockHTTPContext(request)
	status := decoder.Decode(c)
	assert.Equal(t, status, filter.Continue)
}
