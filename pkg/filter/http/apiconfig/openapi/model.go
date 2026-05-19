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

const ValidationPlanMetadataKey = "openapi.validation.plan"

type ValidationPlan struct {
	RoutePattern     string
	PathParameters   []ParameterValidation
	QueryParameters  []ParameterValidation
	HeaderParameters []ParameterValidation
	RequestBody      *BodyValidation
}

type ParameterValidation struct {
	Name      string
	Required  bool
	Enum      []string
	Type      string
	MinLength *int
	MaxLength *int
	Minimum   *float64
	Maximum   *float64
}

type BodyValidation struct {
	Required bool
	Schema   *SchemaValidation
}

type SchemaValidation struct {
	Type       string
	Required   []string
	Properties map[string]*SchemaValidation
	Enum       []string
	MinLength  *int
	MaxLength  *int
	Minimum    *float64
	Maximum    *float64
}
