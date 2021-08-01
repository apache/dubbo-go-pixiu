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

package header

//import (
//	"bytes"
//	"math"
//	"net/http"
//	"testing"
//)
//
//import (
//	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
//	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
//
//	"github.com/stretchr/testify/assert"
//)
//
//import (
//	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
//	"github.com/apache/dubbo-go-pixiu/pkg/filter/recovery"
//)
//
//func TestHeader(t *testing.T) {
//	api := router.API{
//		URLPattern: "/mock/:id/:name",
//		Method:     getMockMethod(config.MethodGet),
//		Headers:    map[string]string{},
//	}
//	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
//	assert.NoError(t, err)
//	c := mock.GetMockHTTPContext(request, newFilter().Do(), recovery.newFilter().Do())
//	c.API(api)
//	c.Next()
//	assert.Equal(t, int8(len(c.Filters)*2), c.BaseContext.Index)
//
//	request.Header.Set("filter", "test")
//	api.Headers["filter"] = "test"
//	c1 := mock.GetMockHTTPContext(request, newFilter().Do(), recovery.newFilter().Do())
//	c1.API(api)
//	c1.Next()
//	assert.Equal(t, int8(len(c1.Filters)*2-1), c1.BaseContext.Index)
//
//	request.Header.Set("filter", "test1")
//	c2 := mock.GetMockHTTPContext(request, newFilter().Do(), recovery.newFilter().Do())
//	c2.API(api)
//	c2.Next()
//	assert.Equal(t, int8(math.MaxInt8/2+1), c2.BaseContext.Index)
//}
//
//func getMockMethod(verb config.HTTPVerb) config.Method {
//	inbound := config.InboundRequest{}
//	integration := config.IntegrationRequest{}
//	return config.Method{
//		OnAir:              true,
//		HTTPVerb:           verb,
//		InboundRequest:     inbound,
//		IntegrationRequest: integration,
//	}
//}
