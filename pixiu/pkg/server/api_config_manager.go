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

package server

import (
	"fmt"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type (
	ApiConfigListener interface {
		OnAddAPI(r router.API) error
		OnRemoveAPI(r router.API) error
		OnDeleteRouter(r config.Resource) error
	}
	// ApiConfigManager similar to RouterManager
	ApiConfigManager struct {
		als map[string]ApiConfigListener
	}
)

func CreateDefaultApiConfigManager(server *Server, bs *model.Bootstrap) *ApiConfigManager {
	acm := &ApiConfigManager{als: make(map[string]ApiConfigListener)}
	return acm
}

func (acm *ApiConfigManager) AddApiConfigListener(adapterID string, l ApiConfigListener) {
	existedListener, existed := acm.als[adapterID]
	if existed {
		panic(fmt.Errorf("AddApiConfigListener %T and %T listen same adapter: %s", l, existedListener, adapterID))
	}
	acm.als[adapterID] = l
}

func (acm *ApiConfigManager) AddAPI(adapterID string, r router.API) error {
	l, existed := acm.als[adapterID]
	if !existed {
		return errors.Errorf("no listener found")
	}
	return l.OnAddAPI(r)
}

func (acm *ApiConfigManager) RemoveAPI(adapterID string, r router.API) error {
	l, existed := acm.als[adapterID]
	if !existed {
		return errors.Errorf("no listener found")
	}
	return l.OnRemoveAPI(r)
}

func (acm *ApiConfigManager) DeleteRouter(adapterID string, r config.Resource) error {
	l, existed := acm.als[adapterID]
	if !existed {
		return errors.Errorf("no listener found")
	}
	return l.OnDeleteRouter(r)
}
