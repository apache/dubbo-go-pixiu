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

package nacos

import (
	"strings"
	"testing"
)

func TestNacosServiceDiscovery_Register(t *testing.T) {
	n := &nacosServiceDiscovery{}

	err := n.Register()

	if err == nil {
		t.Error("Register should return an error")
	}
	if !strings.Contains(err.Error(), "Register not implemented") {
		t.Errorf("Register error should contain 'Register not implemented', got %s", err.Error())
	}
}

func TestNacosServiceDiscovery_UnRegister(t *testing.T) {
	n := &nacosServiceDiscovery{}

	err := n.UnRegister()

	if err == nil {
		t.Error("UnRegister should return an error")
	}
	if !strings.Contains(err.Error(), "UnRegister not implemented") {
		t.Errorf("UnRegister error should contain 'UnRegister not implemented', got %s", err.Error())
	}
}

func TestNacosServiceDiscovery_Get(t *testing.T) {
	n := &nacosServiceDiscovery{}

	result := n.Get("test-service")

	if result != nil {
		t.Errorf("Get should return nil, got %v", result)
	}
}

func TestNacosServiceDiscovery_StartPeriodicalRefresh(t *testing.T) {
	n := &nacosServiceDiscovery{}

	err := n.StartPeriodicalRefresh()

	if err == nil {
		t.Error("StartPeriodicalRefresh should return an error")
	}
	if !strings.Contains(err.Error(), "StartPeriodicalRefresh not implemented") {
		t.Errorf("StartPeriodicalRefresh error should contain 'StartPeriodicalRefresh not implemented', got %s", err.Error())
	}
}
