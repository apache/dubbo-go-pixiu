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
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/xds"
	pixiupb "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"google.golang.org/protobuf/types/known/anypb"
)

func makeListeners2() *pixiupb.PixiuExtensionListeners {
	return &pixiupb.PixiuExtensionListeners{Listeners: []*pixiupb.Listener{
		{
			Name: "net/http889",
			Address: &pixiupb.Address{
				SocketAddress: &pixiupb.SocketAddress{
					Address: "0.0.0.0",
					Port:    8889,
				},
				Name: "http_8889",
			},
			FilterChain: makeHttpFilter(),
		},
	}}
}

func GenerateSnapshotPixiu2() cache.Snapshot {
	ldsResource2, _ := anypb.New(makeListeners2())
	cdsResource, _ := anypb.New(makeClusters())
	snap, _ := cache.NewSnapshot("2",
		map[resource.Type][]types.Resource{
			resource.ExtensionConfigType: {
				&core.TypedExtensionConfig{
					Name:        xds.ClusterType,
					TypedConfig: cdsResource,
				},
				&core.TypedExtensionConfig{
					Name:        xds.ListenerType,
					TypedConfig: ldsResource2,
				},
			},
		},
	)
	return snap
}
