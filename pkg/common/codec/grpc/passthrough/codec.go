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

package passthrough

import (
	"fmt"
)

import (
	"google.golang.org/grpc/encoding"

	"google.golang.org/protobuf/proto"
)

// Name is the name of this codec.
const Name = "pass_through"

// Codec is a gRPC codec that passes through bytes as is.
// This is used for transparent proxying where the message types are unknown at compile time.
type Codec struct{}

func init() {
	encoding.RegisterCodec(Codec{})
}

// Marshal checks if the value is already bytes or a proto.Message and marshals accordingly.
func (c Codec) Marshal(v any) ([]byte, error) {
	if p, ok := v.(proto.Message); ok {
		return proto.Marshal(p)
	}
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	return nil, fmt.Errorf("passthrough codec: cannot marshal type %T, want proto.Message or []byte", v)
}

// Unmarshal stores the raw data into the target, which must be a *[]byte or proto.Message.
func (c Codec) Unmarshal(data []byte, v any) error {
	if vb, ok := v.(*[]byte); ok {
		*vb = data
		return nil
	}
	if p, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, p)
	}
	return fmt.Errorf("passthrough codec: cannot unmarshal into type %T, want *[]byte or proto.Message", v)
}

// Name returns the name of the codec.
func (c Codec) Name() string {
	return Name
}
