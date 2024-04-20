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

// RemoteConfig remote server info which offer server discovery or k/v or distribution function
type RemoteConfig struct {
	Protocol string `yaml:"protocol" json:"protocol" default:"zookeeper"`
	Timeout  string `yaml:"timeout" json:"timeout" default:"10s"`
	Address  string `yaml:"address" json:"address"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Group    string `yaml:"group" json:"group"`
	Root     string `yaml:"root" json:"root" default:"/services"`
}
