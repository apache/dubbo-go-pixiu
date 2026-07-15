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

package triple

import (
	"strings"
	"testing"

	"github.com/apache/dubbo-go-pixiu/pkg/client"
)

func TestClient_Apply(t *testing.T) {
	tc := &Client{}

	err := tc.Apply()

	if err == nil {
		t.Error("Apply should return an error")
	}
	if !strings.Contains(err.Error(), "Apply is deprecated and not implemented") {
		t.Errorf("Apply error should contain 'Apply is deprecated and not implemented', got %s", err.Error())
	}
}

func TestClient_MapParams(t *testing.T) {
	tc := &Client{}
	req := &client.Request{}

	reqData, err := tc.MapParams(req)

	if reqData != nil {
		t.Errorf("MapParams should return nil, got %v", reqData)
	}
	if err == nil {
		t.Error("MapParams should return an error")
	}
	if !strings.Contains(err.Error(), "MapParams is deprecated and not implemented") {
		t.Errorf("MapParams error should contain 'MapParams is deprecated and not implemented', got %s", err.Error())
	}
}