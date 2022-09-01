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

package authority

import (
	"net/http"
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

func TestAuth(t *testing.T) {

	rules := AuthorityConfiguration{
		[]AuthorityRule{{Strategy: Blacklist, Items: append([]string{}, "127.0.0.1")}},
	}
	//rules = []AuthorityRule{{Strategy: Whitelist, Items: append([]string{}, "127.0.0.1")}}
	p := Plugin{}
	authFilter, _ := p.CreateFilterFactory()
	config := authFilter.Config()
	mockYaml, err := yaml.MarshalYML(rules)
	assert.Nil(t, err)
	err = yaml.UnmarshalYML(mockYaml, config)
	assert.Nil(t, err)

	err = authFilter.Apply()
	assert.Nil(t, err)
	f := &Filter{cfg: &rules}

	request, _ := http.NewRequest("GET", "/", nil)
	ctx := mock.GetMockHTTPContext(request)
	filterStatus := f.Decode(ctx)
	assert.Equal(t, filterStatus, filter.Continue)
}
