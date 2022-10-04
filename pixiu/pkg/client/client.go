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

package client

// Client represents the interface of http/dubbo clients
type Client interface {
	// Apply to init client
	Apply() error

	// Close close the client
	Close() error

	// Call invoke the downstream service.
	Call(req *Request) (res interface{}, err error)

	// MapParams mapping param, uri, query, body ...
	MapParams(req *Request) (reqData interface{}, err error)
}

/**
 * the following option is designed to support dubbo pixiu model. you can see
 * https://github.com/apache/dubbo-go-pixiu/pixiu/tree/master.
 */

// MapOption option map, key : name, value : option
type MapOption map[string]RequestOption

// RequestOption option interface.
type RequestOption interface {
	Action(target, val interface{}) error
}
