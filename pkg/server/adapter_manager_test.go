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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	DemoAdapterPlugin struct{}

	DemoAdapter struct {
		cfg *Config
	}

	Config struct {
		Name string `yaml:"name" json:"name,omitempty"`
	}
)

func (d *DemoAdapterPlugin) Kind() string {
	return "test"
}

type initializationTestPlugin struct {
	name    string
	err     error
	created int
}

func (p *initializationTestPlugin) Kind() string { return p.name }

func (p *initializationTestPlugin) CreateAdapter(*model.Adapter) (adapter.Adapter, error) {
	p.created++
	return &initializationTestAdapter{DemoAdapter: DemoAdapter{cfg: &Config{}}, err: p.err}, nil
}

type initializationTestAdapter struct {
	DemoAdapter
	err error
}

func (a *initializationTestAdapter) Apply() error { return a.err }

func TestAdapterManagerDisabledDoesNotBlockEnabled(t *testing.T) {
	p := &initializationTestPlugin{name: t.Name()}
	adapter.RegisterAdapterPlugin(p)
	am := &AdapterManager{configs: []*model.Adapter{
		{Name: p.name, Enabled: "true"},
		{Name: p.name, Enabled: "false"},
		{Name: "unregistered-disabled-adapter", Enabled: "false"},
		{Name: p.name, Enabled: "TRUE"},
	}}
	am.initAdapters()
	require.Equal(t, 2, p.created, "disabled entries must not invoke the factory")
	require.Len(t, am.adapters, 2, "keep enabled adapters before and after disabled entries")
}

func TestAdapterManagerApplyErrorSeverity(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "adapter.log")
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logPath}
	logger.InitLogger(&cfg)
	t.Cleanup(func() { logger.InitLogger(nil) })

	for _, tt := range []struct {
		name  string
		err   error
		level string
	}{
		{name: "unsupported", err: fmt.Errorf("optional setup: %w", errors.ErrUnsupported), level: "warn"},
		{name: "failure", err: errors.New("connection refused"), level: "error"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := &initializationTestPlugin{name: t.Name(), err: tt.err}
			adapter.RegisterAdapterPlugin(p)
			am := &AdapterManager{configs: []*model.Adapter{{Name: p.name}}}
			am.initAdapters()
			require.NoError(t, logger.GetLogger().Sync())
			data, err := os.ReadFile(logPath)
			require.NoError(t, err)
			var matches int
			for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
				var entry struct {
					Level   string `json:"level"`
					Message string `json:"msg"`
				}
				require.NoError(t, json.Unmarshal([]byte(line), &entry))
				if strings.Contains(entry.Message, p.name) {
					matches++
					require.Equal(t, tt.level, entry.Level)
					require.Contains(t, entry.Message, tt.err.Error())
				}
			}
			require.Equal(t, 1, matches, "initialization should log the failure once")
		})
	}
}

func (d *DemoAdapterPlugin) CreateAdapter(ad *model.Adapter) (adapter.Adapter, error) {
	return &DemoAdapter{cfg: &Config{}}, nil
}

func (a *DemoAdapter) Start() {

}

// Apply init
func (a *DemoAdapter) Apply() error {
	return nil
}

// Config get config for Adapter
func (a *DemoAdapter) Config() any {
	return a.cfg
}

func (a *DemoAdapter) Stop() {

}

func TestAdapterManager(t *testing.T) {
	adapter.RegisterAdapterPlugin(&DemoAdapterPlugin{})

	bs := &model.Bootstrap{
		StaticResources: model.StaticResources{
			Adapters: []*model.Adapter{
				{
					Name:   "test",
					Config: make(map[string]any),
				},
			},
		},
	}

	am := CreateDefaultAdapterManager(nil, bs)

	assert.Equal(t, len(am.adapters), 1)

	am.Start()

	am.Stop()

}
