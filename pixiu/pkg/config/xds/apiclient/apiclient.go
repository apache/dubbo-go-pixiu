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

package apiclient

import (
	"reflect"
)

import (
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	anypb "github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type (
	ResourceTypeName        = string
	PixiuDynamicConfigModel proto.Message

	ProtoAny struct {
		typeConfig *v3.TypedExtensionConfig
		any        *anypb.Any
	}

	DeltaResources struct {
		NewResources    []*ProtoAny
		RemovedResource []string
	}
)

func (p *ProtoAny) GetName() string {
	return p.typeConfig.Name
}

func (p *ProtoAny) To(configModel PixiuDynamicConfigModel) error {
	if p.any != nil {
		return p.any.UnmarshalTo(configModel)
	}

	err := p.typeConfig.TypedConfig.UnmarshalTo(configModel)
	if err != nil {
		panic(err)
	}
	return errors.Wrapf(err, "can not covert to %v", reflect.TypeOf(configModel))
}

func NewProtoAny(typeConfig *v3.TypedExtensionConfig) *ProtoAny {
	return &ProtoAny{typeConfig: typeConfig}
}
