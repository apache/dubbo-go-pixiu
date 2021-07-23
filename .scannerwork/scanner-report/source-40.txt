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

package model

// FilterChain filter chain
type FilterChain struct {
	FilterChainMatch FilterChainMatch `yaml:"filter_chain_match" json:"filter_chain_match" mapstructure:"filter_chain_match"`
	Filters          []Filter         `yaml:"filters" json:"filters" mapstructure:"filters"`
}

// Filter core struct, filter is extend by user
type Filter struct {
	Name   string      `yaml:"name" json:"name" mapstructure:"name"`       // Name filter name unique
	Config interface{} `yaml:"config" json:"config" mapstructure:"config"` // Config filter config
}

// FilterChainMatch
type FilterChainMatch struct {
	Domains []string `yaml:"domains" json:"domains" mapstructure:"domains"`
}
