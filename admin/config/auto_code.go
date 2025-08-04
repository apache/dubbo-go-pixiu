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

type Autocode struct {
	TransferRestart bool   `mapstructure:"transfer-restart" json:"transferRestart" yaml:"transfer-restart"`
	Root            string `mapstructure:"root" json:"root" yaml:"root"`
	Server          string `mapstructure:"server" json:"server" yaml:"server"`
	SApi            string `mapstructure:"server-api" json:"serverApi" yaml:"server-api"`
	SInitialize     string `mapstructure:"server-initialize" json:"serverInitialize" yaml:"server-initialize"`
	SModel          string `mapstructure:"server-model" json:"serverModel" yaml:"server-model"`
	SRequest        string `mapstructure:"server-request" json:"serverRequest"  yaml:"server-request"`
	SRouter         string `mapstructure:"server-router" json:"serverRouter" yaml:"server-router"`
	SService        string `mapstructure:"server-service" json:"serverService" yaml:"server-service"`
	Web             string `mapstructure:"web" json:"web" yaml:"web"`
	WApi            string `mapstructure:"web-api" json:"webApi" yaml:"web-api"`
	WForm           string `mapstructure:"web-form" json:"webForm" yaml:"web-form"`
	WTable          string `mapstructure:"web-table" json:"webTable" yaml:"web-table"`
	WFlow           string `mapstructure:"web-flow" json:"webFlow" yaml:"web-flow"`
}
