package server

import (
	"errors"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/stretchr/testify/require"
	"testing"
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
					LdsConfig: copy(*bs.DynamicResources.LdsConfig),
					CdsConfig: copy(*bs.DynamicResources.CdsConfig),
					AdsConfig: copy(*bs.DynamicResources.AdsConfig),
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
						LdsConfig: copy(*bs.DynamicResources.LdsConfig),
						CdsConfig: copy(*bs.DynamicResources.CdsConfig),
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
						LdsConfig: copy(*bs.DynamicResources.LdsConfig),
						CdsConfig: copy(*bs.DynamicResources.CdsConfig),
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
						LdsConfig: copy(*bs.DynamicResources.LdsConfig),
						CdsConfig: copy(*bs.DynamicResources.CdsConfig),
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resMg, err := func() (result DynamicResourceManager, err error) {
				defer func() {
					panicInfo := recover()
					err, _ = panicInfo.(error)
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

func copy(target model.ApiConfigSource) *model.ApiConfigSource {
	return &target
}
