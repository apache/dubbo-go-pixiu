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
	"bytes"
	"math"
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/recovery"
)

func TestAuthority(t *testing.T) {
	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPAuthContext(request, true, New().Do(), recovery.New().Do())
	c.Next()
	assert.Equal(t, int8(math.MaxInt8/2+1), c.BaseContext.Index)

	request.Header.Set("X-Forwarded-For", "0.0.0.0")
	c1 := mock.GetMockHTTPAuthContext(request, false, New().Do(), recovery.New().Do())
	c1.Next()
	assert.Equal(t, int8(math.MaxInt8/2+1), c1.BaseContext.Index)

	c2 := mock.GetMockHTTPAuthContext(request, true, New().Do(), recovery.New().Do())
	c2.Next()
	assert.Equal(t, int8(len(c2.Filters)*2), c2.BaseContext.Index)

	request.Header.Set("X-Forwarded-For", "127.0.0.1")
	c3 := mock.GetMockHTTPAuthContext(request, false, New().Do(), recovery.New().Do())
	c3.Next()
	assert.Equal(t, int8(len(c3.Filters)*2), c3.BaseContext.Index)
}
