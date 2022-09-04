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

package constant

const (
	HeaderKeyContextType = "Content-Type"

	HeaderKeyAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderKeyAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderKeyAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderKeyAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderKeyAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	HeaderValueJsonUtf8        = "application/json;charset=UTF-8"
	HeaderValueTextPlain       = "text/plain"
	HeaderValueApplicationJson = "application/json"

	HeaderValueAll = "*"

	PathSlash           = "/"
	ProtocolSlash       = "://"
	PathParamIdentifier = ":"
)

const (
	Http1HeaderKeyHost = "Host"
	Http2HeaderKeyHost = ":authority"
)

const (
	PprofDefaultAddress = "0.0.0.0"
	PprofDefaultPort    = 7070
)

const (
	Get    = "GET"
	Put    = "PUT"
	Post   = "POST"
	Delete = "DELETE"
)

const (
	DubboHttpDubboVersion   = "x-dubbo-http1.1-dubbo-version"
	DubboServiceProtocol    = "x-dubbo-service-protocol"
	DubboServiceVersion     = "x-dubbo-service-version"
	DubboGroup              = "x-dubbo-service-group"
	DubboServiceMethodTypes = "x-dubbo-service-method-overloading"
)
