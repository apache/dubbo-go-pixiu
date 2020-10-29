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
package registry

import (
	"github.com/apache/dubbo-go/common"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

// TransferURL2Api transfer url and clusterName to IntegrationRequest
func TransferURL2Api(url common.URL, clusterName string) []config.IntegrationRequest {
	var irs []config.IntegrationRequest
	for _, method := range url.Methods {
		irs = append(irs, config.IntegrationRequest{
			RequestType: url.Protocol,
			DubboMetadata: dubbo.DubboMetadata{
				ApplicationName: url.GetParam(constant.NAME_KEY, ""),
				Group:           url.GetParam(constant.GROUP_KEY, ""),
				Version:         url.GetParam(constant.VERSION_KEY, ""),
				Interface:       url.GetParam(constant.INTERFACE_KEY, ""),
				Method:          method,
				Retries:         url.GetParam(constant.RETRIES_KEY, ""),
				ClusterName:     clusterName,
				ProtocolTypeStr: url.Protocol,
			},
		})
	}
	return irs
}
