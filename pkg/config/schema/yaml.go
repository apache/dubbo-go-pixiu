/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func DecodeObjectYAML(data []byte) (ConfigObject, error) {
	var object ConfigObject
	if err := yaml.Unmarshal(data, &object); err != nil {
		return ConfigObject{}, fmt.Errorf("decode config object YAML: %w", err)
	}
	return object, nil
}

func EncodeObjectYAML(object ConfigObject) ([]byte, error) {
	data, err := yaml.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("encode config object YAML: %w", err)
	}
	return data, nil
}

func DecodeConfigSetYAML(data []byte) (ConfigSet, error) {
	var configSet ConfigSet
	if err := yaml.Unmarshal(data, &configSet); err != nil {
		return ConfigSet{}, fmt.Errorf("decode config set YAML: %w", err)
	}
	return configSet, nil
}

func EncodeConfigSetYAML(configSet ConfigSet) ([]byte, error) {
	data, err := yaml.Marshal(configSet)
	if err != nil {
		return nil, fmt.Errorf("encode config set YAML: %w", err)
	}
	return data, nil
}
