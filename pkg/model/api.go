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

package model

import (
	"sync"
)

var (
	CacheApi = sync.Map{}
	EmptyApi = &Api{}
)

// Api is api gateway concept, control request from browser、Mobile APP、third party people
type Api struct {
	Name     string      `json:"name" yaml:"name"`
	ITypeStr string      `json:"itype" yaml:"itype"`
	IType    ApiType     `json:"-" yaml:"-"`
	OTypeStr string      `json:"otype" yaml:"otype"`
	OType    ApiType     `json:"-" yaml:"-"`
	Status   Status      `json:"status" yaml:"status"`
	Metadata interface{} `json:"metadata" yaml:"metadata"`
	Method   string      `json:"method" yaml:"method"`
	RequestMethod
}

// NewApi
func NewApi() *Api {
	return &Api{}
}

// FindApi find a api, if not exist, return false
func (a *Api) FindApi(name string) (*Api, bool) {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api), true
	}

	return nil, false
}

// MatchMethod
func (a *Api) MatchMethod(method string) bool {
	i := RequestMethodValue[method]
	if a.RequestMethod == RequestMethod(i) {
		return true
	}

	return false
}

// IsOk api status equals Up
func (a *Api) IsOk(name string) bool {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*Api).Status == Up
	}

	return false
}

// Offline api offline
func (a *Api) Offline(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*Api).Status = Down
	}
}

// Online api online
func (a *Api) Online(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*Api).Status = Up
	}
}
