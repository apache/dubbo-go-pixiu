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

package opa

import (
	"testing"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
)

func setBootstrap(t *testing.T, b *adminconfig.AdminBootstrap) {
	t.Helper()
	prev := adminconfig.Bootstrap
	adminconfig.Bootstrap = b
	t.Cleanup(func() { adminconfig.Bootstrap = prev })
}

func TestResolveOPAServerURL_Precedence(t *testing.T) {
	t.Run("nil bootstrap falls back to default", func(t *testing.T) {
		setBootstrap(t, nil)
		if got := resolveOPAServerURL(""); got != adminconfig.DefaultOPAServerURL {
			t.Fatalf("want default %q, got %q", adminconfig.DefaultOPAServerURL, got)
		}
	})

	t.Run("empty bootstrap field falls back to default", func(t *testing.T) {
		setBootstrap(t, &adminconfig.AdminBootstrap{})
		if got := resolveOPAServerURL(""); got != adminconfig.DefaultOPAServerURL {
			t.Fatalf("want default %q, got %q", adminconfig.DefaultOPAServerURL, got)
		}
	})

	t.Run("bootstrap value used when caller passes empty", func(t *testing.T) {
		setBootstrap(t, &adminconfig.AdminBootstrap{
			OPA: adminconfig.OPAConfig{ServerURL: "http://configured:1234"},
		})
		if got := resolveOPAServerURL(""); got != "http://configured:1234" {
			t.Fatalf("want bootstrap value, got %q", got)
		}
	})

	t.Run("caller value overrides bootstrap and gets trimmed", func(t *testing.T) {
		setBootstrap(t, &adminconfig.AdminBootstrap{
			OPA: adminconfig.OPAConfig{ServerURL: "http://configured:1234"},
		})
		if got := resolveOPAServerURL("  http://override:9999  "); got != "http://override:9999" {
			t.Fatalf("caller should override bootstrap, got %q", got)
		}
	})
}

func TestResolveOPAPolicyID_Precedence(t *testing.T) {
	t.Run("nil bootstrap falls back to default", func(t *testing.T) {
		setBootstrap(t, nil)
		if got := resolveOPAPolicyID(""); got != adminconfig.DefaultOPAPolicyID {
			t.Fatalf("want default %q, got %q", adminconfig.DefaultOPAPolicyID, got)
		}
	})

	t.Run("bootstrap value used when caller passes empty", func(t *testing.T) {
		setBootstrap(t, &adminconfig.AdminBootstrap{
			OPA: adminconfig.OPAConfig{PolicyID: "from-config"},
		})
		if got := resolveOPAPolicyID(""); got != "from-config" {
			t.Fatalf("want bootstrap value, got %q", got)
		}
	})

	t.Run("caller value overrides bootstrap", func(t *testing.T) {
		setBootstrap(t, &adminconfig.AdminBootstrap{
			OPA: adminconfig.OPAConfig{PolicyID: "from-config"},
		})
		if got := resolveOPAPolicyID("override-id"); got != "override-id" {
			t.Fatalf("caller should override bootstrap, got %q", got)
		}
	})
}
