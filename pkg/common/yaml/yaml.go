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

package yaml

import (
	"fmt"
	"os"
	"path"
)

import (
	perrors "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// LoadYMLConfig Load yml config byte from file
func LoadYMLConfig(confProFile string) ([]byte, error) {
	if len(confProFile) == 0 {
		return nil, perrors.Errorf("configure file name is nil")
	}
	if path.Ext(confProFile) != ".yml" && path.Ext(confProFile) != ".yaml" {
		return nil, perrors.Errorf("configure file name{%v} suffix must be .yml or .yaml", confProFile)
	}

	return os.ReadFile(confProFile)
}

// UnmarshalYMLConfig Load yml config byte from file, then unmarshal to object
func UnmarshalYMLConfig(confProFile string, out interface{}) error {
	confFileStream, err := LoadYMLConfig(confProFile)
	if err != nil {
		return perrors.Errorf("os.ReadFile(file:%s) = error:%v", confProFile, perrors.WithStack(err))
	}
	return yaml.Unmarshal(confFileStream, out)
}

// UnmarshalYML unmarshals decodes the first document found within the in byte slice and assigns decoded values into the out value.
func UnmarshalYML(data []byte, out interface{}) error {
	return yaml.Unmarshal(data, out)
}

// MarshalYML serializes the value provided into a YAML document.
func MarshalYML(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}

// ParseConfig get config struct from map[string]interface{}
func ParseConfig(factoryConfStruct interface{}, conf map[string]interface{}) error {
	// conf will be map, convert to yaml
	yamlBytes, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	// Unmarshal yamlStr to factoryConf
	err = yaml.Unmarshal(yamlBytes, factoryConfStruct)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
