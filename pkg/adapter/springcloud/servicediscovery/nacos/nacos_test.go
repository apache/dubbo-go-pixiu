package nacos

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"reflect"
	"testing"
)

func TestNewNacosServiceDiscovery(t *testing.T) {
	type args struct {
		config *model.RemoteConfig
	}
	tests := []struct {
		name    string
		args    args
		want    servicediscovery.ServiceDiscovery
		wantErr bool
	}{
		{
			name: "nacos",
			args: args{config: &model.RemoteConfig{
				Protocol: "nacos",
				Address:  "127.0.0.1:8848",
				Group:    "DEFAULT_GROUP",
			}},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNacosServiceDiscovery(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNacosServiceDiscovery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got.QueryServices()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNacosServiceDiscovery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
