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

package filter

import "github.com/apache/dubbo-go-pixiu/pkg/api/context"

// Filter filter func, filter
type Filter func(context.Context)

// Factory filter Factory, for FilterFunc. filter manager will fill the Configuration from local or config center
type Factory interface {
	// Config return the pointer of config
	Config() interface{}

	// Apply return the filter func, initialize whatever you want before return the func
	// use c.next() to next filter, before is pre logic, after is post logic.
	Apply() (Filter, error)
}

// ErrResponse err response.
type ErrResponse struct {
	Message string `json:"message"`
}
