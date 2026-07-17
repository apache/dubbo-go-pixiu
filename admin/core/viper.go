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

package core

import (
	"fmt"
	"os"
	"path/filepath"
)

import (
	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/global"
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

// Viper loads the admin config from configPath. When configPath is empty, the
// GVA_CONFIG env var is consulted, falling back to the default config file name.
// The path is taken from the caller (the admin CLI's --config flag) rather than
// parsed here, so an explicitly specified config file is never silently ignored.
func Viper(configPath string) (*viper.Viper, error) {
	if configPath == "" {
		if configEnv := os.Getenv(utils.ConfigEnv); configEnv != "" {
			configPath = configEnv
		} else {
			configPath = utils.ConfigFile
		}
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
	}
	global.CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v, nil
}
