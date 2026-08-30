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
	"errors"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
)

const defaultRouteTimeout = time.Second

// CompiledRoute is the preview boundary between the Admin-facing object and
// Pixiu's current APIConfig model. Publish metadata is deliberately omitted.
type CompiledRoute struct {
	Source   AdminObject
	Resource legacyconfig.Resource
	Method   legacyconfig.Method
}

// CompileAdminRouteBinding validates and converts one high-level binding into
// the legacy Resource and Method structures consumed by Pixiu.
func CompileAdminRouteBinding(registry *Registry, object AdminObject) (CompiledRoute, error) {
	if registry == nil {
		return CompiledRoute{}, errors.New("admin object schema registry is nil")
	}
	normalized, err := registry.Normalize(object)
	if err != nil {
		return CompiledRoute{}, err
	}
	if normalized.Kind != KindAdminRouteBinding {
		return CompiledRoute{}, fmt.Errorf("cannot compile kind %q", normalized.Kind)
	}

	entry, err := objectField(normalized.Spec, "entry")
	if err != nil {
		return CompiledRoute{}, err
	}
	target, err := objectField(normalized.Spec, "target")
	if err != nil {
		return CompiledRoute{}, err
	}

	path := stringField(entry, "path")
	entryProtocol := stringField(entry, "protocol")
	targetProtocol := stringField(target, "protocol")
	method := legacyconfig.Method{
		ResourcePath: path,
		Enable:       true,
		Timeout:      defaultRouteTimeout,
		HTTPVerb:     stringField(entry, "method"),
		InboundRequest: legacyconfig.InboundRequest{
			RequestType: entryProtocol,
		},
		IntegrationRequest: legacyconfig.IntegrationRequest{
			RequestType: targetProtocol,
			DubboBackendConfig: legacyconfig.DubboBackendConfig{
				ClusterName:     stringField(target, "cluster"),
				ApplicationName: stringField(target, "application"),
				Protocol:        targetProtocol,
				Group:           stringField(target, "group"),
				Version:         stringField(target, "version"),
				Interface:       stringField(target, "interface"),
				Method:          stringField(target, "method"),
			},
		},
	}

	params, _ := normalized.Spec["params"].([]any)
	method.MappingParams = make([]legacyconfig.MappingParam, 0, len(params))
	method.ParameterTypes = make([]string, len(params))
	for _, rawParam := range params {
		param, _ := rawParam.(map[string]any)
		toNumber, _ := numericValue(param["to"])
		to := int(toNumber)
		paramType := stringField(param, "type")
		method.MappingParams = append(method.MappingParams, legacyconfig.MappingParam{
			Name:    stringField(param, "from"),
			MapTo:   strconv.Itoa(to),
			MapType: paramType,
		})
		method.ParameterTypes[to] = paramType
	}

	resource := legacyconfig.Resource{
		Type:    "restful",
		Path:    path,
		Timeout: defaultRouteTimeout,
	}
	return CompiledRoute{
		Source:   normalized,
		Resource: resource,
		Method:   method,
	}, nil
}

// LegacyAPIConfig assembles the generated Resource and Method into the shape
// expected by the current api_config.yaml loader.
func (r CompiledRoute) LegacyAPIConfig() legacyconfig.APIConfig {
	resource := r.Resource
	resource.Methods = []legacyconfig.Method{r.Method}
	return legacyconfig.APIConfig{
		Name:      r.Source.Metadata.Name,
		Resources: []legacyconfig.Resource{resource},
	}
}

// PreviewYAML renders a human-readable legacy api_config.yaml preview while
// preserving duration strings instead of time.Duration nanoseconds.
func (r CompiledRoute) PreviewYAML() ([]byte, error) {
	preview := legacyAPIConfigPreview{
		Name: r.Source.Metadata.Name,
		Resources: []legacyResourcePreview{{
			Path:    r.Resource.Path,
			Type:    r.Resource.Type,
			Timeout: r.Resource.Timeout.String(),
			Methods: []legacyMethodPreview{{
				HTTPVerb: r.Method.HTTPVerb,
				Enable:   r.Method.Enable,
				Timeout:  r.Method.Timeout.String(),
				InboundRequest: legacyInboundRequestPreview{
					RequestType: r.Method.InboundRequest.RequestType,
				},
				IntegrationRequest: legacyIntegrationRequestPreview{
					RequestType:     r.Method.IntegrationRequest.RequestType,
					ApplicationName: r.Method.ApplicationName,
					Protocol:        r.Method.Protocol,
					Group:           r.Method.Group,
					Version:         r.Method.Version,
					Interface:       r.Method.Interface,
					Method:          r.Method.Method,
					ClusterName:     r.Method.ClusterName,
					ParameterTypes:  append([]string(nil), r.Method.ParameterTypes...),
					MappingParams:   append([]legacyconfig.MappingParam(nil), r.Method.MappingParams...),
				},
			}},
		}},
	}
	data, err := yaml.Marshal(preview)
	if err != nil {
		return nil, fmt.Errorf("encode legacy api_config preview: %w", err)
	}
	return data, nil
}

type legacyAPIConfigPreview struct {
	Name      string                  `yaml:"name"`
	Resources []legacyResourcePreview `yaml:"resources"`
}

type legacyResourcePreview struct {
	Path    string                `yaml:"path"`
	Type    string                `yaml:"type"`
	Timeout string                `yaml:"timeout"`
	Methods []legacyMethodPreview `yaml:"methods"`
}

type legacyMethodPreview struct {
	HTTPVerb           string                          `yaml:"httpVerb"`
	Enable             bool                            `yaml:"enable"`
	Timeout            string                          `yaml:"timeout"`
	InboundRequest     legacyInboundRequestPreview     `yaml:"inboundRequest"`
	IntegrationRequest legacyIntegrationRequestPreview `yaml:"integrationRequest"`
}

type legacyInboundRequestPreview struct {
	RequestType string `yaml:"requestType"`
}

type legacyIntegrationRequestPreview struct {
	RequestType     string                      `yaml:"requestType"`
	ApplicationName string                      `yaml:"applicationName"`
	Protocol        string                      `yaml:"protocol"`
	Group           string                      `yaml:"group,omitempty"`
	Version         string                      `yaml:"version,omitempty"`
	Interface       string                      `yaml:"interface"`
	Method          string                      `yaml:"method"`
	ClusterName     string                      `yaml:"clusterName"`
	ParameterTypes  []string                    `yaml:"parameterTypes,omitempty"`
	MappingParams   []legacyconfig.MappingParam `yaml:"mappingParams,omitempty"`
}

func objectField(values map[string]any, name string) (map[string]any, error) {
	value, ok := values[name].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("field %s is not an object", name)
	}
	return value, nil
}

func stringField(values map[string]any, name string) string {
	value, _ := values[name].(string)
	return value
}
