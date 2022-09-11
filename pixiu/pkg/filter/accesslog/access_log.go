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

package accesslog

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPAccessLogFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// FilterFactory is http filter instance
	FilterFactory struct {
		conf *AccessLogConfig
		alw  *AccessLogWriter
	}
	Filter struct {
		conf *AccessLogConfig
		alw  *AccessLogWriter

		start time.Time
	}
)

// Kind return plugin kind
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilter create filter
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		conf: &AccessLogConfig{},
		alw: &AccessLogWriter{
			AccessLogDataChan: make(chan AccessLogData, constant.LogDataBuffer),
		},
	}, nil
}

// PrepareFilterChain prepare chain when http context init
func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{alw: factory.alw}
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)
	return nil
}

// Decode process http context
func (f *Filter) Decode(c *http.HttpContext) filter.FilterStatus {
	f.start = time.Now()
	return filter.Continue
}

func (f *Filter) Encode(c *http.HttpContext) filter.FilterStatus {
	latency := time.Since(f.start)
	// build access_log message
	accessLogMsg := buildAccessLogMsg(c, latency)
	if len(accessLogMsg) > 0 {
		f.alw.Writer(AccessLogData{AccessLogConfig: *f.conf, AccessLogMsg: accessLogMsg})
	}
	return filter.Continue
}

// Config return config of filter
func (factory *FilterFactory) Config() interface{} {
	return factory.conf
}

// Apply init after config set
func (factory *FilterFactory) Apply() error {
	// init
	factory.alw.Write()
	return nil
}

func buildAccessLogMsg(c *http.HttpContext, cost time.Duration) string {
	req := c.Request
	valueStr := req.URL.Query().Encode()
	if len(valueStr) != 0 {
		valueStr = strings.ReplaceAll(valueStr, "&", ",")
	}

	builder := strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(time.Now().Format(constant.MessageDateLayout))
	builder.WriteString("] ")
	builder.WriteString(req.RemoteAddr)
	builder.WriteString(" -> ")
	builder.WriteString(req.Host)
	builder.WriteString(" - ")
	if len(valueStr) > 0 {
		builder.WriteString("request params: [")
		builder.WriteString(valueStr)
		builder.WriteString("] ")
	}
	builder.WriteString("cost time [ ")
	builder.WriteString(strconv.Itoa(int(cost)) + " ]")
	err := c.GetLocalReplyBody()
	if err != nil {
		builder.WriteString(fmt.Sprintf("invoke err [ %v", err))
		builder.WriteString("] ")
	}
	resp := c.TargetResp.Data
	if err != nil {
		builder.WriteString(fmt.Sprintf(" response can not convert to string"))
		builder.WriteString("] ")
	} else {
		builder.WriteString(fmt.Sprintf(" response [ %+v", string(resp)))
		builder.WriteString("] ")
	}
	// builder.WriteString("\n")
	return builder.String()
}

// converter interface to byte array
func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
