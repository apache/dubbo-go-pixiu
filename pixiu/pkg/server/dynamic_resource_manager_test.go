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

package server

import (
	"errors"
	"fmt"
	"runtime/debug"
	"testing"
)

import (
	"github.com/cch123/supermonkey"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
)

func Test_createDynamicResourceManger(t *testing.T) {
	bs := model.Bootstrap{
		Node: &model.Node{
			Cluster: "cluster-1",
			Id:      "node-xx",
		},
		DynamicResources: &model.DynamicResources{
			LdsConfig: &model.ApiConfigSource{
				APITypeStr:     "GRPC",
				ClusterName:    []string{"t1"},
				RefreshDelay:   "1s",
				RequestTimeout: "2s",
				GrpcServices: []model.GrpcService{
					{
						Timeout:         "3s",
						InitialMetadata: []model.HeaderValue{},
					},
				},
			},
			CdsConfig: &model.ApiConfigSource{
				APITypeStr:     "GRPC",
				ClusterName:    []string{"t2"},
				RefreshDelay:   "11s",
				RequestTimeout: "12s",
				GrpcServices: []model.GrpcService{
					{
						Timeout:         "13s",
						InitialMetadata: []model.HeaderValue{},
					},
				},
			},
			AdsConfig: &model.ApiConfigSource{
				APITypeStr:     "GRPC",
				ClusterName:    []string{"t3"},
				RefreshDelay:   "21s",
				RequestTimeout: "22s",
				GrpcServices: []model.GrpcService{
					{
						Timeout:         "23s",
						InitialMetadata: []model.HeaderValue{},
					},
				},
			},
		},
	}
	type args struct {
		bs *model.Bootstrap
	}
	tests := []struct {
		name      string
		args      args
		expectNil bool
		wantNode  *model.Node
		wantLds   *model.ApiConfigSource
		wantCds   *model.ApiConfigSource
		wantErr   error
	}{
		{
			name:      "empty",
			args:      args{bs: &model.Bootstrap{}},
			expectNil: true,
		},
		{
			name: "full",
			args: args{bs: &model.Bootstrap{
				DynamicResources: &model.DynamicResources{
					LdsConfig: bs.DynamicResources.LdsConfig,
					CdsConfig: bs.DynamicResources.CdsConfig,
				},
				Node: bs.Node,
			}},
			expectNil: false,
			wantNode:  bs.Node,
			wantCds:   bs.DynamicResources.CdsConfig,
			wantLds:   bs.DynamicResources.LdsConfig,
		},
		{
			name: "validate-failed-ads",
			args: args{bs: &model.Bootstrap{
				DynamicResources: &model.DynamicResources{
					LdsConfig: copyConfig(*bs.DynamicResources.LdsConfig),
					CdsConfig: copyConfig(*bs.DynamicResources.CdsConfig),
					AdsConfig: copyConfig(*bs.DynamicResources.AdsConfig),
				},
				Node: bs.Node,
			}},
			expectNil: true,
			wantErr:   errors.New(""),
		},
		{
			name: "validate-failed-ads",
			args: args{bs: func() *model.Bootstrap {
				bs := model.Bootstrap{
					DynamicResources: &model.DynamicResources{
						LdsConfig: copyConfig(*bs.DynamicResources.LdsConfig),
						CdsConfig: copyConfig(*bs.DynamicResources.CdsConfig),
					},
					Node: bs.Node,
				}
				bs.DynamicResources.LdsConfig.APITypeStr = "HTTP"
				return &bs
			}()},
			expectNil: true,
			wantErr:   errors.New(""),
		},
		{
			name: "validate-failed-apitype",
			args: args{bs: func() *model.Bootstrap {
				bs := model.Bootstrap{
					DynamicResources: &model.DynamicResources{
						LdsConfig: copyConfig(*bs.DynamicResources.LdsConfig),
						CdsConfig: copyConfig(*bs.DynamicResources.CdsConfig),
					},
					Node: bs.Node,
				}
				bs.DynamicResources.LdsConfig.APITypeStr = "REST"
				return &bs
			}()},
			expectNil: true,
			wantErr:   errors.New(""),
		},
		{
			name: "validate-failed-api-type",
			args: args{bs: func() *model.Bootstrap {
				bs := model.Bootstrap{
					DynamicResources: &model.DynamicResources{
						LdsConfig: copyConfig(*bs.DynamicResources.LdsConfig),
						CdsConfig: copyConfig(*bs.DynamicResources.CdsConfig),
					},
					Node: bs.Node,
				}
				bs.DynamicResources.LdsConfig.APITypeStr = "DUBBO2"
				return &bs
			}()},
			expectNil: true,
			wantErr:   errors.New(""),
		},
	}

	supermonkey.Patch(xds.StartXdsClient, func(listenerMg controls.ListenerManager, clusterMg controls.ClusterManager, drm controls.DynamicResourceManager) xds.Client {
		return nil
	})
	supermonkey.Patch((*Server).GetListenerManager, func(_ *Server) *ListenerManager {
		return nil
	})
	supermonkey.Patch((*Server).GetClusterManager, func(_ *Server) *ClusterManager {
		return nil
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resMg, err := func() (result DynamicResourceManager, err error) {
				defer func() {
					panicInfo := recover()
					err, _ = panicInfo.(error)
					if err != nil {
						fmt.Println(err)
						debug.PrintStack()
					}
				}()
				result = createDynamicResourceManger(tt.args.bs)
				return
			}()
			assert := require.New(t)
			if tt.wantErr != nil {
				assert.Error(err)
				return
			}
			if tt.expectNil {
				assert.Nil(resMg)
				return
			}
			assert.NoError(err)
			assert.NotNil(resMg)
			assert.Equal(tt.wantLds, resMg.GetLds())
			assert.Equal(tt.wantCds, resMg.GetCds())
			assert.Equal(tt.wantNode, resMg.GetNode())
		})
	}
}

func copyConfig(target model.ApiConfigSource) *model.ApiConfigSource {
	return &target
}
