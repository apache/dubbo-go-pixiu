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

package main

import (
	"os"
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	config2 "github.com/apache/dubbo-go-pixiu/admin/config"
)

func TestAdminRejectsUnsafeSigningKeyBeforeLoadingConfiguration(t *testing.T) {
	command := getRootCmd()
	originalConfigPath, originalAPIConfigPath := configPath, apiConfigPath
	originalBootstrap := config2.Bootstrap
	t.Cleanup(func() {
		configPath, apiConfigPath = originalConfigPath, originalAPIConfigPath
		config2.Bootstrap = originalBootstrap
		command.SetArgs(nil)
	})
	configurationPath := filepath.Join(t.TempDir(), "admin.yaml")
	require.NoError(t, os.WriteFile(configurationPath, []byte("name: signing-key-regression\n"), 0o600))
	_, err := config2.LoadAPIConfigFromFile(configurationPath)
	require.NoError(t, err)
	config2.Bootstrap = originalBootstrap

	for _, key := range []string{"dubbo-go-pixiu", "short-key"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv("DUBBOGO_PIXIU_JWT_SIGN_KEY", key)
			command.SetArgs([]string{"--config", configurationPath})
			err := command.Execute()
			require.ErrorContains(t, err, "DUBBOGO_PIXIU_JWT_SIGN_KEY")
			require.True(t, config2.Bootstrap == originalBootstrap, "configuration loaded before key validation")

			err = Start()
			require.ErrorContains(t, err, "DUBBOGO_PIXIU_JWT_SIGN_KEY")
		})
	}
}
