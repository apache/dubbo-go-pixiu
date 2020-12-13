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

package initialize

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/api"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/authority"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/recovery"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/remote"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/response"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/timeout"
	sa "github.com/dubbogo/dubbo-go-proxy/pkg/service/api"
)

// Run start init.
func Run() {
	filterInit()
	apiDiscoveryServiceInit()
}

func filterInit() {
	api.Init()
	authority.Init()
	logger.Init()
	recovery.Init()
	remote.Init()
	response.Init()
	timeout.Init()
}

func apiDiscoveryServiceInit() {
	sa.Init()
}
