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
package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

type ApiLoader interface {
	// every apiloader has a Priority, since remote api configurer such as nacos may cover fileapiloader.
	// and the larger priority number indicates it's apis can cover apis of lower priority number priority.
	GetPrior() int
	GetLoadedApiConfigs() ([]model.Api, error)
	// load all apis at first time
	InitLoad() error
	// watch apis when apis configured were changed.
	HotLoad() (chan struct{}, error)
	// clear all apis
	Clear() error
}
