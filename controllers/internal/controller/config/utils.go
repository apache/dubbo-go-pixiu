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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

import (
	"gopkg.in/yaml.v3"
)

func GetConfigFromFile(filename string) (*Config, error) {
	cfg := NewDefaultConfig()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	envVarMap := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envVarMap[pair[0]] = pair[1]
	}

	tpl := template.New("text").Option("missingkey=error")
	tpl, err = tpl.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("error parsing configuration template %v", err)
	}
	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, envVarMap)
	if err != nil {
		return nil, fmt.Errorf("error execute configuration template %v", err)
	}

	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		err = yaml.Unmarshal(buf.Bytes(), cfg)
	} else {
		err = json.Unmarshal(buf.Bytes(), cfg)
	}

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetControllerName() string {
	return ControllerConfig.ControllerName
}
