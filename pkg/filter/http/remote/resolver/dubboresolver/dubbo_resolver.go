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

package dubboresolver

import (
	apiConf "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

import (
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/remote/resolver"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	resolver.RegisterResolver(model.StandardDubboResolver, StandardDubboResolver{})
}

// StandardDubboResolver handles Dubbo generic calls that include parameter types.
type StandardDubboResolver struct {
	resolver.BaseResolver
}

func (s StandardDubboResolver) Resolve(ctx *contexthttp.HttpContext) (*router.API, error) {
	req := ctx.Request
	if err := s.PreCheck(req); err != nil {
		return nil, err // Not applicable
	}

	// This resolver specifically handles generic calls WITH types.
	mappingParams := []apiConf.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
		{Name: "uri.application", MapTo: "opt.application"},
		{Name: "uri.interface", MapTo: "opt.interface"},
		{Name: "uri.method", MapTo: "opt.method"},
	}

	return s.BuildAPI(req, mappingParams)
}
