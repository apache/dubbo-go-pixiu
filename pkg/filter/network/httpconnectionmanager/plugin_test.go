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

package httpconnectionmanager

import (
	"fmt"
	"testing"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const rejectingHTTPFilter = "dgp.test.rejecting-http-filter"

type rejectingHTTPPlugin struct{}
type rejectingHTTPFactory struct{}

func (*rejectingHTTPPlugin) Kind() string { return rejectingHTTPFilter }
func (*rejectingHTTPPlugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &rejectingHTTPFactory{}, nil
}
func (*rejectingHTTPFactory) Config() any  { return &struct{}{} }
func (*rejectingHTTPFactory) Apply() error { return fmt.Errorf("invalid plugin configuration") }
func (*rejectingHTTPFactory) PrepareFilterChain(*contexthttp.HttpContext, filter.FilterChain) error {
	return nil
}

func TestPluginRejectsInvalidNestedHTTPFilter(t *testing.T) {
	plugin := &Plugin{}
	_, err := plugin.CreateFilter(&model.HttpConnectionManagerConfig{
		HTTPFilters: []*model.HTTPFilter{{
			Name:   "missing.http.filter.plugin",
			Config: map[string]any{},
		}},
	})

	require.ErrorContains(t, err, "missing.http.filter.plugin")
}

func TestPluginRejectsNestedHTTPFilterConfiguration(t *testing.T) {
	filter.RegisterHttpFilter(&rejectingHTTPPlugin{})
	plugin := &Plugin{}
	_, err := plugin.CreateFilter(&model.HttpConnectionManagerConfig{
		HTTPFilters: []*model.HTTPFilter{{Name: rejectingHTTPFilter}},
	})

	require.ErrorContains(t, err, "invalid plugin configuration")
}
