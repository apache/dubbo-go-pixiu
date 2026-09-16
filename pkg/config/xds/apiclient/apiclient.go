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
	"sync"
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
		NewResources     []*ProtoAny
		RemovedResources []string

		applyResult chan error
		complete    sync.Once
		versions    map[string]string
	}
)

func (p *ProtoAny) GetName() string {
	if p == nil || p.typeConfig == nil {
		return ""
	}
	return p.typeConfig.Name
}

func (p *ProtoAny) To(configModel PixiuDynamicConfigModel) error {
	if p == nil {
		return errors.New("resource is nil")
	}
	if configModel == nil {
		return errors.New("target config model is nil")
	}
	if p.any != nil {
		return errors.Wrapf(p.any.UnmarshalTo(configModel), "can not convert to %v", reflect.TypeOf(configModel))
	}
	if p.typeConfig == nil || p.typeConfig.TypedConfig == nil {
		return errors.New("typed extension config is nil")
	}
	return errors.Wrapf(p.typeConfig.TypedConfig.UnmarshalTo(configModel), "can not convert to %v", reflect.TypeOf(configModel))
}

func NewProtoAny(typeConfig *v3.TypedExtensionConfig) *ProtoAny {
	return &ProtoAny{typeConfig: typeConfig}
}

func newDeltaResources() *DeltaResources {
	return &DeltaResources{
		NewResources:     make([]*ProtoAny, 0, 1),
		RemovedResources: make([]string, 0, 1),
		applyResult:      make(chan error, 1),
		versions:         make(map[string]string, 1),
	}
}

// Complete reports whether the response was applied successfully. Extension
// config streams use the result to send an ACK or NACK to the control plane.
// Other discovery clients may leave the completion channel unset.
func (d *DeltaResources) Complete(err error) {
	if d == nil || d.applyResult == nil {
		return
	}
	d.complete.Do(func() {
		d.applyResult <- err
	})
}
