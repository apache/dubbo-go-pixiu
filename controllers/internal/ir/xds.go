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

package ir

// Xds holds the intermediate representation of a Gateway and is
// used by the Converter to convert it into Pixiu conf.yaml resources.
// This is similar to Envoy Gateway's IR.Xds but adapted for Pixiu.
type Xds struct {
	HTTP []*HTTPListener `json:"http,omitempty" yaml:"http,omitempty"`
	TCP  []*TCPListener  `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	UDP  []*UDPListener  `json:"udp,omitempty" yaml:"udp,omitempty"`
}

type HTTPListener struct {
	Name      string       `json:"name" yaml:"name"`
	Address   string       `json:"address" yaml:"address"`
	Port      int32        `json:"port" yaml:"port"`
	Hostnames []string     `json:"hostnames,omitempty" yaml:"hostnames,omitempty"`
	Routes    []*HTTPRoute `json:"routes,omitempty" yaml:"routes,omitempty"`
}

type TCPListener struct {
	Name    string      `json:"name" yaml:"name"`
	Address string      `json:"address" yaml:"address"`
	Port    int32       `json:"port" yaml:"port"`
	Routes  []*TCPRoute `json:"routes,omitempty" yaml:"routes,omitempty"`
}

type UDPListener struct {
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
	Port    int32  `json:"port" yaml:"port"`
}

type HTTPRoute struct {
	Name          string              `json:"name" yaml:"name"`
	PathMatch     *StringMatch        `json:"pathMatch,omitempty" yaml:"pathMatch,omitempty"`
	HeaderMatches []*StringMatch      `json:"headerMatches,omitempty" yaml:"headerMatches,omitempty"`
	Method        *string             `json:"method,omitempty" yaml:"method,omitempty"`
	Destinations  []*RouteDestination `json:"destinations,omitempty" yaml:"destinations,omitempty"`
	Filters       []*HTTPFilter       `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type TCPRoute struct {
	Name         string              `json:"name" yaml:"name"`
	Destinations []*RouteDestination `json:"destinations,omitempty" yaml:"destinations,omitempty"`
}

type StringMatch struct {
	Exact     *string `json:"exact,omitempty" yaml:"exact,omitempty"`
	Prefix    *string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	SafeRegex *string `json:"safeRegex,omitempty" yaml:"safeRegex,omitempty"`
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
}

type RouteDestination struct {
	Name   string  `json:"name" yaml:"name"`
	Weight *uint32 `json:"weight,omitempty" yaml:"weight,omitempty"`
}

type HTTPFilter struct {
	Name   string                 `json:"name" yaml:"name"`
	Config map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
}

type Cluster struct {
	Name               string              `json:"name" yaml:"name"`
	Endpoints          []*Endpoint         `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	LoadBalancerPolicy *LoadBalancerPolicy `json:"loadBalancerPolicy,omitempty" yaml:"loadBalancerPolicy,omitempty"`
}

type Endpoint struct {
	ID      int    `json:"id" yaml:"id"`
	Address string `json:"address" yaml:"address"`
	Port    int32  `json:"port" yaml:"port"`
}

type LoadBalancerPolicy struct {
	Type string `json:"type" yaml:"type"`
}
