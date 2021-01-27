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

package main

import (
	"log"
)

import (
	"github.com/dubbogo/dubbo-go-proxy-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy-filter/pkg/filter"
)

func main() {}

// Access filter
type Access struct {
}

// ExternalPluginAccess export filter
func ExternalPluginAccess() filter.Filter {
	return &Access{}
}

// Do to export func(c context.Context)
func (r *Access) Do() context.FilterFunc {
	return func(c context.Context) {
		log.Print("DoAccessFilter")
		c.Next()
	}
}
