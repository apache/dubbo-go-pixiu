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
				CacheDir:    "/tmp/nacos/cache",
				LogDir:      "/tmp/nacos/log",
				NamespaceId: Namespace,
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
					opt.Remote = true
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
					opt.Remote = true
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
					opt.Remote = true
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
