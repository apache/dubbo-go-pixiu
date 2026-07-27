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

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeYAML writes content to a .yaml file in a fresh tempdir and returns the path.
func writeYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "admin.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}
	return path
}

// restoreBootstrap resets the global Bootstrap pointer after each test so
// tests don't leak state into one another.
func restoreBootstrap(t *testing.T) {
	t.Helper()
	prev := Bootstrap
	t.Cleanup(func() { Bootstrap = prev })
}

// TestLoadAPIConfigFromFile_OPA exercises the real admin startup config-loading
// path used in cmd/admin/admin.go:55 — same function, same yaml library, same
// global var. Asserts that an `opa:` section is parsed into OPAConfig and that
// `request_timeout: 3s` deserializes to a time.Duration of 3 seconds.
func TestLoadAPIConfigFromFile_OPA(t *testing.T) {
	restoreBootstrap(t)

	path := writeYAML(t, `
server:
  address: 127.0.0.1:18091
etcd:
  address: 127.0.0.1:2379
  path: /pixiu/config/api/test
mysql:
  username: root
  password: x
  host: 127.0.0.1
  port: "3306"
  dbname: pixiu
opa:
  server_url: http://127.0.0.1:18181
  policy_id: e2e-policy
  request_timeout: 3s
`)

	b, err := LoadAPIConfigFromFile(path)
	if err != nil {
		t.Fatalf("LoadAPIConfigFromFile: %v", err)
	}
	if b == nil {
		t.Fatal("returned bootstrap is nil")
	}
	if Bootstrap != b {
		t.Fatal("global Bootstrap pointer was not updated")
	}

	if b.OPA.ServerURL != "http://127.0.0.1:18181" {
		t.Errorf("ServerURL: got %q", b.OPA.ServerURL)
	}
	if b.OPA.PolicyID != "e2e-policy" {
		t.Errorf("PolicyID: got %q", b.OPA.PolicyID)
	}
	if b.OPA.RequestTimeout != 3*time.Second {
		t.Errorf("RequestTimeout: want 3s, got %v (ns=%d)",
			b.OPA.RequestTimeout, b.OPA.RequestTimeout.Nanoseconds())
	}
}

// TestLoadAPIConfigFromFile_OPADurationFormats locks in the duration formats
// the YAML loader (gopkg.in/yaml.v3) accepts via time.Duration.UnmarshalText.
// Catches regressions if the loader is swapped for one that doesn't support
// Go duration strings.
func TestLoadAPIConfigFromFile_OPADurationFormats(t *testing.T) {
	cases := []struct {
		yamlValue string
		want      time.Duration
	}{
		{"3s", 3 * time.Second},
		{"500ms", 500 * time.Millisecond},
		{"1m30s", 90 * time.Second},
		{"2h", 2 * time.Hour},
	}
	for _, c := range cases {
		t.Run(c.yamlValue, func(t *testing.T) {
			restoreBootstrap(t)
			path := writeYAML(t, "opa:\n  request_timeout: "+c.yamlValue+"\n")
			b, err := LoadAPIConfigFromFile(path)
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			if b.OPA.RequestTimeout != c.want {
				t.Fatalf("want %v, got %v", c.want, b.OPA.RequestTimeout)
			}
		})
	}
}

// TestLoadAPIConfigFromFile_NoOPASection ensures admin still loads when the
// operator omits `opa:` entirely (back-compat); fields take Go zero values
// and downstream code falls back to DefaultOPA* constants.
func TestLoadAPIConfigFromFile_NoOPASection(t *testing.T) {
	restoreBootstrap(t)
	path := writeYAML(t, `
server:
  address: 127.0.0.1:18091
`)
	b, err := LoadAPIConfigFromFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if b.OPA.ServerURL != "" || b.OPA.PolicyID != "" || b.OPA.RequestTimeout != 0 {
		t.Errorf("expected zero OPAConfig when section omitted, got %+v", b.OPA)
	}
}

// TestLoadAPIConfigFromFile_MissingPath surfaces the explicit error message —
// guards against a refactor that silently swallows misconfiguration.
func TestLoadAPIConfigFromFile_MissingPath(t *testing.T) {
	if _, err := LoadAPIConfigFromFile(""); err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}
