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

package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

func ValidateRequest(req *http.Request, plan *ValidationPlan) error {
	if req == nil || plan == nil {
		return nil
	}

	if err := validatePathParameters(req, plan); err != nil {
		return err
	}

	if err := validateQueryParameters(req, plan.QueryParameters); err != nil {
		return err
	}

	if err := validateHeaderParameters(req, plan.HeaderParameters); err != nil {
		return err
	}

	if err := validateRequestBody(req, plan.RequestBody); err != nil {
		return err
	}

	return nil
}

func validatePathParameters(req *http.Request, plan *ValidationPlan) error {
	if plan.RoutePattern == "" || len(plan.PathParameters) == 0 {
		return nil
	}

	values := router.GetURIParams(&router.API{URLPattern: plan.RoutePattern}, *req.URL)
	for _, param := range plan.PathParameters {
		value := values.Get(param.Name)
		if param.Required && strings.TrimSpace(value) == "" {
			return fmt.Errorf("path parameter %q is required", param.Name)
		}
		if value == "" {
			continue
		}
		if err := validateParameterValue("path parameter", value, param); err != nil {
			return err
		}
	}

	return nil
}

func validateQueryParameters(req *http.Request, params []ParameterValidation) error {
	query := req.URL.Query()
	for _, param := range params {
		value := query.Get(param.Name)
		if param.Required && strings.TrimSpace(value) == "" {
			return fmt.Errorf("query parameter %q is required", param.Name)
		}
		if value == "" {
			continue
		}
		if err := validateParameterValue("query parameter", value, param); err != nil {
			return err
		}
	}
	return nil
}

func validateHeaderParameters(req *http.Request, params []ParameterValidation) error {
	for _, param := range params {
		value := req.Header.Get(param.Name)
		if param.Required && strings.TrimSpace(value) == "" {
			return fmt.Errorf("header %q is required", param.Name)
		}
		if value == "" {
			continue
		}
		if err := validateParameterValue("header", value, param); err != nil {
			return err
		}
	}
	return nil
}

func validateParameterValue(kind string, value string, param ParameterValidation) error {
	if len(param.Enum) > 0 && !contains(param.Enum, value) {
		return fmt.Errorf("%s %q must be one of %v", kind, param.Name, param.Enum)
	}

	switch param.Type {
	case "", "string":
		if param.MinLength != nil && len(value) < *param.MinLength {
			return fmt.Errorf("%s %q length must be >= %d", kind, param.Name, *param.MinLength)
		}
		if param.MaxLength != nil && len(value) > *param.MaxLength {
			return fmt.Errorf("%s %q length must be <= %d", kind, param.Name, *param.MaxLength)
		}
	case "integer":
		number, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%s %q must be an integer", kind, param.Name)
		}
		if param.Minimum != nil && float64(number) < *param.Minimum {
			return fmt.Errorf("%s %q must be >= %v", kind, param.Name, *param.Minimum)
		}
		if param.Maximum != nil && float64(number) > *param.Maximum {
			return fmt.Errorf("%s %q must be <= %v", kind, param.Name, *param.Maximum)
		}
	case "number":
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%s %q must be a number", kind, param.Name)
		}
		if param.Minimum != nil && number < *param.Minimum {
			return fmt.Errorf("%s %q must be >= %v", kind, param.Name, *param.Minimum)
		}
		if param.Maximum != nil && number > *param.Maximum {
			return fmt.Errorf("%s %q must be <= %v", kind, param.Name, *param.Maximum)
		}
	case "boolean":
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("%s %q must be a boolean", kind, param.Name)
		}
	}

	return nil
}

func validateRequestBody(req *http.Request, body *BodyValidation) error {
	if body == nil {
		return nil
	}

	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(raw))

	if body.Required && len(bytes.TrimSpace(raw)) == 0 {
		return fmt.Errorf("request body is required")
	}
	if len(bytes.TrimSpace(raw)) == 0 || body.Schema == nil {
		return nil
	}

	if !strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "application/json") {
		return fmt.Errorf("request body content-type must be application/json")
	}

	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("request body must be valid json: %w", err)
	}

	return validateSchema("$", payload, body.Schema)
}

func validateSchema(path string, value any, schema *SchemaValidation) error {
	if schema == nil {
		return nil
	}

	switch schema.Type {
	case "object":
		obj, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("%s must be an object", path)
		}
		for _, field := range schema.Required {
			if _, exists := obj[field]; !exists {
				return fmt.Errorf("%s.%s is required", path, field)
			}
		}
		for key, child := range schema.Properties {
			childValue, exists := obj[key]
			if !exists {
				continue
			}
			if err := validateSchema(path+"."+key, childValue, child); err != nil {
				return err
			}
		}
	case "string":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s must be a string", path)
		}
		if schema.MinLength != nil && len(str) < *schema.MinLength {
			return fmt.Errorf("%s length must be >= %d", path, *schema.MinLength)
		}
		if schema.MaxLength != nil && len(str) > *schema.MaxLength {
			return fmt.Errorf("%s length must be <= %d", path, *schema.MaxLength)
		}
		if len(schema.Enum) > 0 && !contains(schema.Enum, str) {
			return fmt.Errorf("%s must be one of %v", path, schema.Enum)
		}
	case "number":
		number, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("%s must be a number", path)
		}
		if schema.Minimum != nil && number < *schema.Minimum {
			return fmt.Errorf("%s must be >= %v", path, *schema.Minimum)
		}
		if schema.Maximum != nil && number > *schema.Maximum {
			return fmt.Errorf("%s must be <= %v", path, *schema.Maximum)
		}
	case "integer":
		number, ok := toFloat64(value)
		if !ok || math.Mod(number, 1) != 0 {
			return fmt.Errorf("%s must be an integer", path)
		}
		if schema.Minimum != nil && number < *schema.Minimum {
			return fmt.Errorf("%s must be >= %v", path, *schema.Minimum)
		}
		if schema.Maximum != nil && number > *schema.Maximum {
			return fmt.Errorf("%s must be <= %v", path, *schema.Maximum)
		}
	}

	return nil
}

func contains(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func toFloat64(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	default:
		return 0, false
	}
}
