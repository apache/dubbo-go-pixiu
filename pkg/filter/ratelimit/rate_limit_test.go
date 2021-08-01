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
	"testing"
)

import (
	"github.com/ghodss/yaml"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

func TestFilter(t *testing.T) {
	factory := newFactory()
	config, ok := factory.Config().(*Config)
	assert.True(t, ok)

	mockYaml, err := yaml.Marshal(mockConfig())
	assert.Nil(t, err)
	err = yaml.Unmarshal(mockYaml, config)
	assert.Nil(t, err)

	apply, err := factory.Apply()
	assert.Nil(t, err)

	ctx := &http.HttpContext{
		BaseContext: context.NewBaseContext(),
	}
	apply(ctx)
}
