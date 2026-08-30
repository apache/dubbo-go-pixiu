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

import "github.com/apache/dubbo-go-pixiu/pkg/common/copyutil"

const KindAdminRouteBinding = "AdminRouteBinding"

// AdminObject is the YAML object edited by Admin. Its stable envelope stays
// small; the registered schema gives the dynamic spec its actual contract.
type AdminObject struct {
	Kind     string         `json:"kind" yaml:"kind"`
	Metadata ObjectMetadata `json:"metadata" yaml:"metadata"`
	Spec     map[string]any `json:"spec" yaml:"spec"`
}

type ObjectMetadata struct {
	Name string `json:"name" yaml:"name"`
}

func (o AdminObject) Clone() AdminObject {
	cloned := o
	cloned.Spec = copyutil.CloneStringAnyMap(o.Spec)
	return cloned
}
