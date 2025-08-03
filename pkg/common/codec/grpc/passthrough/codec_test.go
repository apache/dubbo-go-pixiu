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

package passthrough_test

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/codec/grpc/passthrough"
)

func TestCodec(t *testing.T) {
	codec := passthrough.Codec{}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, passthrough.Name, codec.Name())
	})

	t.Run("Marshal", func(t *testing.T) {
		t.Run("proto.Message success", func(t *testing.T) {
			msg := &wrapperspb.StringValue{Value: "test_proto_message"}
			expectedBytes, err := proto.Marshal(msg)
			assert.NoError(t, err)

			actualBytes, err := codec.Marshal(msg)
			assert.NoError(t, err)
			assert.Equal(t, expectedBytes, actualBytes)
		})

		t.Run("[]byte success", func(t *testing.T) {
			data := []byte("test_byte_slice")
			actualBytes, err := codec.Marshal(data)
			assert.NoError(t, err)
			assert.Equal(t, data, actualBytes)
		})

		t.Run("unsupported type failure", func(t *testing.T) {
			_, err := codec.Marshal(12345)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot marshal type int")
		})
	})

	t.Run("Unmarshal", func(t *testing.T) {
		t.Run("proto.Message success", func(t *testing.T) {
			originalMsg := &wrapperspb.StringValue{Value: "test_proto_message"}
			data, err := proto.Marshal(originalMsg)
			assert.NoError(t, err)

			var targetMsg wrapperspb.StringValue
			err = codec.Unmarshal(data, &targetMsg)
			assert.NoError(t, err)
			assert.True(t, proto.Equal(originalMsg, &targetMsg))
		})

		t.Run("[]byte success", func(t *testing.T) {
			data := []byte("test_byte_slice")
			var target []byte
			err := codec.Unmarshal(data, &target)
			assert.NoError(t, err)
			assert.Equal(t, data, target)
		})

		t.Run("unsupported type failure", func(t *testing.T) {
			data := []byte("some data")
			var target int
			err := codec.Unmarshal(data, &target)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot unmarshal into type *int")
		})
	})
}
