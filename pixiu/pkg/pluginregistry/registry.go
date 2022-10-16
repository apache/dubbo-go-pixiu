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
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/cluster/loadbalancer/rand"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/cluster/loadbalancer/roundrobin"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/accesslog"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/auth/jwt"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/authority"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/cors"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/csrf"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/header"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/host"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/apiconfig"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/dubboproxy"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/grpcproxy"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/httpproxy"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/loadbalancer"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/proxyrewrite"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/proxywasm"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/remote"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/metric"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/network/dubboproxy"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/network/dubboproxy/filter/http"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/network/dubboproxy/filter/proxy"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/network/grpcconnectionmanager"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/network/httpconnectionmanager"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/prometheus"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/seata"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/tracing"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/traffic"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/listener/http"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/listener/http2"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/listener/tcp"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/listener/triple"
)
