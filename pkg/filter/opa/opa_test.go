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

package opa

import (
	"net/http/httptest"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import(
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const testPolicy = `
package test
import future.keywords.if

default allow := false

allow if {
	 input.headers.X[0] == "1"
}
`

func setupFilterWithoutFile(t *testing.T, policy string) *Filter {
	p := Plugin{}

	filterFactory, err := p.CreateFilterFactory()
	assert.Nil(t, err)

	fFactory := filterFactory.(*FilterFactory)
	
	fFactory.cfg = &Config{
		Policy:     policy,
		Entrypoint: "data.test.allow",
	}

	err = fFactory.Apply()
	assert.Nil(t, err)

	return &Filter{
		cfg:           fFactory.cfg,
		preparedQuery: fFactory.preparedQuery,
	}
}

func TestAllowedRule(t *testing.T) {
	f := setupFilterWithoutFile(t, testPolicy)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X", "1")

	rec := httptest.NewRecorder()
	ctx := &http.HttpContext{
		Writer:  rec,
		Request: req,
	}

	result := f.Decode(ctx)
	assert.Equal(t, filter.Continue, result)
}

