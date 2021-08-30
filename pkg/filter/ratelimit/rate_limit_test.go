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
	http2 "net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/recovery"
)

func TestFilter(t *testing.T) {
	factory := Plugin{}
	f, _ := factory.CreateFilter()
	config, ok := f.Config().(*Config)
	assert.True(t, ok)

	mockYaml, err := yaml.MarshalYML(mockConfig())
	assert.Nil(t, err)
	err = yaml.UnmarshalYML(mockYaml, config)
	assert.Nil(t, err)

	err = f.Apply()
	assert.Nil(t, err)

	request, _ := http2.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	c := mock.GetMockHTTPContext(request, f, recovery.GetMock())
	c.Next()
}
