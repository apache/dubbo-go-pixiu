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
	Init() error
	Close() error

	// Call invoke the downstream service.
	Call(req *Request) (res interface{}, err error)

	// MapParams mapping param, uri, query, body ...
	MapParams(req *Request) (reqData interface{}, err error)
}

/**
 * follow option is design to support dubbo proxy model. you can see
 * https://github.com/dubbogo/dubbo-go-proxy/tree/master.
 */

// MapOption option map, key : name, value : option
type MapOption map[string]RequestOption

// RequestOption option interface.
type RequestOption interface {
	// Usable if option can use
	Usable() bool
	// SetUsable set usable
	SetUsable(b bool)
	// Action do with val for special
	Action(req *Request, val interface{})
	// VirtualPos virtual position
	VirtualPos() int
}

// CommonOption common opt.
type CommonOption struct {
	usable bool
	RequestOption
}

// Usable get usable.
func (opt *CommonOption) Usable() bool {
	return opt.usable
}

// SetUsable set usable.
func (opt *CommonOption) SetUsable(b bool) {
	opt.usable = b
}

// VirtualPos virtual position, default 0.
func (opt *CommonOption) VirtualPos() int {
	return 0
}
