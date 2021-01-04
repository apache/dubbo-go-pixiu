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

// LbPolicy the load balance policy enum
type LbPolicy int32

const (
	RoundRobin LbPolicy = 0
	IPHash     LbPolicy = 1
	WightRobin LbPolicy = 2
	Rand       LbPolicy = 3
)

// LbPolicyName key int32 for LbPolicy, value string
var LbPolicyName = map[int32]string{
	0: "RoundRobin",
	1: "IPHash",
	2: "WightRobin",
	3: "Rand",
}

// LbPolicyValue key string, value int32 for LbPolicy
var LbPolicyValue = map[string]int32{
	"RoundRobin": 0,
	"IPHash":     1,
	"WightRobin": 2,
	"Rand":       3,
}
