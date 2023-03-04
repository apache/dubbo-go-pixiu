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

package proxy

import (
	"context"
)

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
)

type Reflector struct {
	c *grpcreflect.Client
}

// NewReflector creates a new Reflector from the reflection client
func NewReflector(client *grpcreflect.Client) *Reflector {
	return &Reflector{client}
}

// CreateInvocation creates a MethodInvocation by performing reflection
func (r *Reflector) CreateInvocation(ctx context.Context, serviceName, methodName string, input []byte) (*MethodInvocation, error) {
	serviceDesc, err := r.c.ResolveService(serviceName)
	if err != nil {
		return nil, err
	}
	methodDesc := serviceDesc.FindMethodByName(methodName)
	if methodDesc == nil {
		return nil, errors.New("method not found upstream")
	}
	inputMessage := dynamic.NewMessage(methodDesc.GetInputType())
	err = inputMessage.UnmarshalJSON(input)
	if err != nil {
		return nil, err
	}
	return &MethodInvocation{
		MethodDescriptor: methodDesc,
		Message:          inputMessage,
	}, nil
}

// MethodInvocation contains a method and a message used to invoke an RPC
type MethodInvocation struct {
	*desc.MethodDescriptor
	*dynamic.Message
}
