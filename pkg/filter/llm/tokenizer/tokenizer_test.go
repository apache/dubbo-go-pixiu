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

package tokenizer

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
)

// TestUnaryResponseWithEncodings is a table-driven test for unary (non-streaming) responses.
// It covers multiple content encodings like gzip and deflate.
func TestUnaryResponseWithEncodings(t *testing.T) {
	// This is the payload we expect to process after decompression.
	const payload = `{
       "usage": {
          "prompt_tokens": 7
       }
    }`

	// Helper function to compress data with gzip for our test case.
	compressGzipBytes := func(data string) []byte {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write([]byte(data))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)
		return buf.Bytes()
	}

	// Helper function to compress data with flate/deflate for our test case.
	compressFlateBytes := func(data string) []byte {
		var buf bytes.Buffer
		writer, err := flate.NewWriter(&buf, -1)
		assert.NoError(t, err)
		_, err = writer.Write([]byte(data))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)
		return buf.Bytes()
	}

	// Define all test cases in a table.
	testCases := []struct {
		name     string
		encoding string
		getData  func(string) []byte
	}{
		{
			name:     "No Encoding",
			encoding: "",
			getData: func(s string) []byte {
				return []byte(s)
			},
		},
		{
			name:     "Gzip Encoding",
			encoding: "gzip",
			getData:  compressGzipBytes,
		},
		{
			name:     "Flate Encoding",
			encoding: "deflate",
			getData:  compressFlateBytes,
		},
	}

	// Run the tests for each case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &Filter{}

			request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
			assert.NoError(t, err)
			c := mock.GetMockHTTPContext(request)

			// Prepare the (potentially) compressed data
			compressedData := tc.getData(payload)

			c.TargetResp = &client.UnaryResponse{
				Data: compressedData,
			}
			c.AddHeader(constant.HeaderKeyContentEncoding, tc.encoding)

			// Call the filter's Encode method
			filter.Encode(c)
		})
	}
}

// TestStreamResponseWithEncodings is a table-driven test for streaming responses.
// It replaces the old TestStreamResponse.
func TestStreamResponseWithEncodings(t *testing.T) {
	// This is the payload we expect to process after decompression.
	const payload = `data: {
       "usage": {
          "prompt_tokens": 7
       }
    }`

	// Helper function to compress data with gzip for our test case.
	compressGzip := func(data string) io.Reader {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write([]byte(data))
		assert.NoError(t, err)
		err = writer.Close() // IMPORTANT: Close flushes the writer.
		assert.NoError(t, err)
		return &buf
	}

	compressFlate := func(data string) io.Reader {
		var buf bytes.Buffer
		writer, _ := flate.NewWriter(&buf, -1)
		_, err := writer.Write([]byte(data))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)
		return &buf
	}

	// Define all test cases in a table.
	testCases := []struct {
		name      string
		encoding  string
		getStream func(string) io.Reader
	}{
		{
			name:     "No Encoding",
			encoding: "",
			getStream: func(s string) io.Reader {
				return strings.NewReader(s)
			},
		},
		{
			name:      "Gzip Encoding",
			encoding:  "gzip",
			getStream: compressGzip,
		},
		{
			name:      "Flate Encoding",
			encoding:  "deflate",
			getStream: compressFlate,
		},
	}

	// Run the tests for each case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &Filter{}

			req, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
			assert.NoError(t, err)
			ctx := mock.GetMockHTTPContext(req)

			// Prepare the compressed stream and the response header
			compressedStream := tc.getStream(payload)

			// Set up the mock response
			ctx.TargetResp = &client.StreamResponse{
				Stream: io.NopCloser(compressedStream),
			}

			ctx.AddHeader(constant.HeaderKeyContentEncoding, tc.encoding)

			// Call the filter's Encode method
			filter.Encode(ctx)

			// Give the goroutine a moment to process the data
			buf := make([]byte, 1024)
			ctx.TargetResp.(*client.StreamResponse).Stream.Read(buf)
			time.Sleep(5 * time.Millisecond)
			ctx.TargetResp.(*client.StreamResponse).Stream.Close()
		})
	}
}
