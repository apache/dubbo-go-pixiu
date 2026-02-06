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

package controller

import (
	"testing"
)

func TestSplitClusterName(t *testing.T) {
	tests := []struct {
		name          string
		clusterName   string
		wantNamespace string
		wantService   string
		wantOk        bool
	}{
		{
			name:          "valid cluster name",
			clusterName:   "default-myservice",
			wantNamespace: "default",
			wantService:   "myservice",
			wantOk:        true,
		},
		{
			// Note: Current implementation uses SplitN with n=2, so only first hyphen is used as delimiter
			// This means namespace cannot contain hyphens - "kube-system-my-service" splits to ("kube", "system-my-service")
			name:          "hyphen splits at first occurrence only",
			clusterName:   "kube-system-my-service",
			wantNamespace: "kube",
			wantService:   "system-my-service",
			wantOk:        true,
		},
		{
			name:          "no hyphen",
			clusterName:   "invalidname",
			wantNamespace: "",
			wantService:   "",
			wantOk:        false,
		},
		{
			name:          "empty string",
			clusterName:   "",
			wantNamespace: "",
			wantService:   "",
			wantOk:        false,
		},
		{
			name:          "only hyphen",
			clusterName:   "-",
			wantNamespace: "",
			wantService:   "",
			wantOk:        false,
		},
		{
			name:          "empty namespace",
			clusterName:   "-service",
			wantNamespace: "",
			wantService:   "",
			wantOk:        false,
		},
		{
			name:          "empty service",
			clusterName:   "namespace-",
			wantNamespace: "",
			wantService:   "",
			wantOk:        false,
		},
		{
			name:          "multiple hyphens",
			clusterName:   "ns-svc-name-v1",
			wantNamespace: "ns",
			wantService:   "svc-name-v1",
			wantOk:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace, service, ok := splitClusterName(tt.clusterName)

			if ok != tt.wantOk {
				t.Errorf("splitClusterName(%q) ok = %v, want %v", tt.clusterName, ok, tt.wantOk)
			}

			if namespace != tt.wantNamespace {
				t.Errorf("splitClusterName(%q) namespace = %q, want %q", tt.clusterName, namespace, tt.wantNamespace)
			}

			if service != tt.wantService {
				t.Errorf("splitClusterName(%q) service = %q, want %q", tt.clusterName, service, tt.wantService)
			}
		})
	}
}
