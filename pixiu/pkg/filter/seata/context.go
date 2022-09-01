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

package seata

import (
	"bytes"
	"encoding/json"
	"net/http"
)

import (
	"vimagination.zapto.org/byteio"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

var (
	VarHost        = "host"
	VarQueryString = "queryString"
)

type RequestContext struct {
	ActionContext map[string]string
	Headers       http.Header
	Body          []byte
	Trailers      http.Header
}

func (ctx *RequestContext) Encode() ([]byte, error) {
	var (
		err               error
		actionContextData []byte
		headersData       []byte
		trailersData      []byte
		b                 bytes.Buffer
	)

	w := byteio.BigEndianWriter{Writer: &b}

	if ctx.ActionContext == nil || len(ctx.ActionContext) == 0 {
		w.WriteUint32(0)
	} else {
		actionContextData, err = json.Marshal(ctx.ActionContext)
		if err != nil {
			return nil, err
		}

		w.WriteUint32(uint32(len(actionContextData)))
		w.Write(actionContextData)
	}

	if ctx.Headers == nil || len(ctx.Headers) == 0 {
		w.WriteUint32(0)
	} else {
		headersData, err = json.Marshal(ctx.Headers)
		if err != nil {
			return nil, err
		}
		w.WriteUint32(uint32(len(headersData)))
		w.Write(headersData)
	}

	if ctx.Trailers == nil || len(ctx.Trailers) == 0 {
		w.WriteUint32(0)
	} else {
		trailersData, err = json.Marshal(ctx.Trailers)
		if err != nil {
			return nil, err
		}
		w.WriteUint32(uint32(len(trailersData)))
		w.Write(trailersData)
	}

	w.WriteUint32(uint32(len(ctx.Body)))
	b.Write(ctx.Body)

	return b.Bytes(), nil
}

func (ctx *RequestContext) Decode(b []byte) error {
	var (
		actionContextData []byte
		headersData       []byte
		bodyData          []byte
		trailersData      []byte
	)
	r := byteio.BigEndianReader{Reader: bytes.NewReader(b)}

	contextLength, _, err := r.ReadUint32()
	if err != nil {
		return err
	}
	if contextLength > 0 {
		actionContextData = make([]byte, contextLength, contextLength)
		r.Read(actionContextData)
	}

	headerLength, _, err := r.ReadUint32()
	if err != nil {
		return err
	}
	if headerLength > 0 {
		headersData = make([]byte, headerLength, headerLength)
		r.Read(headersData)
	}

	trailersLength, _, err := r.ReadUint32()
	if err != nil {
		return err
	}
	if trailersLength > 0 {
		trailersData = make([]byte, trailersLength, trailersLength)
		r.Read(trailersData)
	}

	bodyLength, _, err := r.ReadUint32()
	if err != nil {
		return err
	}
	if bodyLength > 0 {
		bodyData = make([]byte, bodyLength, bodyLength)
		r.Read(bodyData)
	}

	if actionContextData != nil {
		err = json.Unmarshal(actionContextData, &(ctx.ActionContext))
		if err != nil {
			logger.Errorf("unmarshal action context failed, %v", err)
		}
	}

	if headersData != nil {
		err = json.Unmarshal(headersData, &(ctx.Headers))
		if err != nil {
			logger.Errorf("unmarshal headers failed, %v", err)
		}
	}

	if trailersData != nil {
		err = json.Unmarshal(trailersData, &(ctx.ActionContext))
		if err != nil {
			logger.Errorf("unmarshal trailers failed, %v", err)
		}
	}

	ctx.Body = bodyData
	return nil
}
