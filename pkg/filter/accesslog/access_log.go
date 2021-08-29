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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPAccessLogFilter
)

func init() {
	filter.RegisterHttpFilter(&AccessPlugin{})
}

type (
	// AccessPlugin is http filter plugin.
	AccessPlugin struct {
	}
	// AccessFilter is http filter instance
	AccessFilter struct {
		conf *AccessLogConfig
		alw  *AccessLogWriter
	}
)

// Kind return plugin kind
func (ap *AccessPlugin) Kind() string {
	return Kind
}

// CreateFilter create filter
func (ap *AccessPlugin) CreateFilter() (filter.HttpFilter, error) {
	return &AccessFilter{
		conf: &AccessLogConfig{},
		alw: &AccessLogWriter{
			AccessLogDataChan: make(chan AccessLogData, constant.LogDataBuffer),
		},
	}, nil
}

// PrepareFilterChain prepare chain when http context init
func (af *AccessFilter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(af.Handle)
	return nil
}

// Handle handle http context
func (af *AccessFilter) Handle(c *http.HttpContext) {
	start := time.Now()
	c.Next()
	latency := time.Since(start)
	// build access_log message
	accessLogMsg := buildAccessLogMsg(c, latency)
	if len(accessLogMsg) > 0 {
		af.alw.Writer(AccessLogData{AccessLogConfig: *af.conf, AccessLogMsg: accessLogMsg})
	}
}

// Config return config of filter
func (af *AccessFilter) Config() interface{} {
	return af.conf
}

// Apply init after config set
func (af *AccessFilter) Apply() error {
	// init
	af.alw.Write()
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
	err := c.Err
	if err != nil {
		builder.WriteString(fmt.Sprintf("invoke err [ %v", err))
		builder.WriteString("] ")
	}
	resp := c.TargetResp.Data
	rbs, err := getBytes(resp)
	if err != nil {
		builder.WriteString(fmt.Sprintf(" response can not convert to string"))
		builder.WriteString("] ")
	} else {
		builder.WriteString(fmt.Sprintf(" response [ %+v", string(rbs)))
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
