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
