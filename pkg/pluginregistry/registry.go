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

package pluginregistry

import (
	// network filters
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/network/httpconnectionmanager"
	// http filters
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/accesslog"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/authority"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/cors"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/header"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/host"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/grpcproxy"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/httpproxy"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/remote"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/metric"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/response"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/timeout"

	// adapter
	_ "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry"
	_ "github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud"
)
