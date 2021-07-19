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
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"github.com/stretchr/testify/assert"
)

type mockRegFacade struct{}

func (m *mockRegFacade) LoadInterfaces() ([]router.API, []error) {
	return nil, nil
}
func (m *mockRegFacade) LoadApplications() ([]router.API, []error) {
	return nil, nil
}
func (m *mockRegFacade) DoSubscribe(*common.URL) error {
	return nil
}
func (m *mockRegFacade) DoUnsubscribe(*common.URL) error {
	return nil
}

func CreateMockAPI(urlPattern string) router.API {
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

func TestDeduplicate(t *testing.T) {
	apis := []router.API{}
	for i := 0; i < 10; i++ {
		apis = append(apis, CreateMockAPI(""))
	}
	var facade *mockRegFacade = &mockRegFacade{}
	baseReg := NewBaseRegistry(facade)
	rst := baseReg.deduplication(apis)
	assert.Equal(t, len(rst), 10)

	apis = []router.API{
		CreateMockAPI("/happy/friday"),
		CreateMockAPI("/happy/friday"),
		CreateMockAPI("/happy/Saturday"),
	}
	rst = baseReg.deduplication(apis)
	assert.Equal(t, len(rst), 2)
}
