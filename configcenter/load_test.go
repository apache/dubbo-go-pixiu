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

package configcenter

import (
	"encoding/json"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

func getBootstrap() *model.Bootstrap {

	return &model.Bootstrap{
		Config: &model.ConfigCenter{
			Type:   "nacos",
			Enable: "true",
		},
		Nacos: &model.Nacos{
			ServerConfigs: []*model.NacosServerConfig{
				{
					IpAddr:      IpAddr,
					Port:        Port,
					Scheme:      Scheme,
					ContextPath: ContextPath,
				},
			},
			ClientConfig: &model.NacosClientConfig{
				CacheDir:            "./.cache",
				LogDir:              "./.log",
				NotLoadCacheAtStart: true,
				NamespaceId:         Namespace,
			},
			DataId: DataId,
			Group:  Group,
		},
	}
}

func getNacosClient() ConfigClient {
	config, _ := NewNacosConfig(getBootstrap())
	return config
}

func TestDefaultConfigLoad_LoadConfigs(t *testing.T) {
	type fields struct {
		bootConfig   *model.Bootstrap
		configClient ConfigClient
	}
	type args struct {
		boot *model.Bootstrap
		opts []Option
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantV   *model.Bootstrap
		wantErr bool
	}{
		{
			name: "Normal_NacosConfigLoadTest",
			fields: fields{
				bootConfig:   getBootstrap(),
				configClient: getNacosClient(),
			},
			args: args{
				boot: getBootstrap(),
				opts: []Option{func(opt *Options) {
					opt.Remote = false // off
					opt.DataId = DataId
					opt.Group = Group
				}},
			},
			wantV:   nil,
			wantErr: false,
		},
		{
			name: "NilDataId_NacosConfigLoadTest",
			fields: fields{
				bootConfig:   getBootstrap(),
				configClient: getNacosClient(),
			},
			args: args{
				boot: getBootstrap(),
				opts: []Option{func(opt *Options) {
					opt.Remote = false // off
					opt.DataId = ""
					opt.Group = Group
				}},
			},
			wantV:   nil,
			wantErr: false,
		},
		{
			name: "ErrorDataId_NacosConfigLoadTest",
			fields: fields{
				bootConfig:   getBootstrap(),
				configClient: getNacosClient(),
			},
			args: args{
				boot: getBootstrap(),
				opts: []Option{func(opt *Options) {
					opt.Remote = false // off
					opt.DataId = "ErrorDataId"
					opt.Group = Group
				}},
			},
			wantV:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultConfigLoad{
				bootConfig:   tt.fields.bootConfig,
				configClient: tt.fields.configClient,
			}
			gotV, err := d.LoadConfigs(tt.args.boot, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//assert.True(t, gotV.Nacos.DataId == DataId, "load config by nacos config center error!")
			//assert.True(t, len(gotV.StaticResources.Listeners) > 0, "load config by nacos config center error!")
			conf, _ := json.Marshal(gotV)
			logger.Infof("config of Bootstrap load by nacos : %v", string(conf))
		})
	}
}
