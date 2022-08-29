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

package baseregistry

import (
	"math/rand"
	"time"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
)

var _ registry.Registry = new(mockRegFacade)

type mockRegFacade struct {
	BaseRegistry
}

func (m *mockRegFacade) DoSubscribe() error {
	return nil
}
func (m *mockRegFacade) DoUnsubscribe() error {
	return nil
}

func CreateMockRegisteredAPI(urlPattern string) router.API {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randStringRunes := func(n int) string {
		rand.Seed(time.Now().UnixNano())
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}
	if len(urlPattern) == 0 {
		urlPattern = randStringRunes(4) + "/" + randStringRunes(4) + "/" + randStringRunes(4)
	}
	return router.API{
		URLPattern: urlPattern,
	}
}
