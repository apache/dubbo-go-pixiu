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
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	manager "github.com/apache/dubbo-go-pixiu/pkg/filter"
)

func Init() {
	manager.RegisterFilterFactory(constant.AccessLogFilter, newAccessLog)
}

type accessLog struct {
	conf *AccessLogConfig
	alw  *AccessLogWriter
}

func newAccessLog() filter.Factory {
	return &accessLog{
		conf: &AccessLogConfig{},
		alw: &AccessLogWriter{
			AccessLogDataChan: make(chan AccessLogData, constant.LogDataBuffer),
		},
	}
}

func (a *accessLog) Config() interface{} {
	return a.conf
}

func (a *accessLog) Apply() (filter.Filter, error) {
	// init
	a.alw.Write()

	return accessLogFunc(a.alw, a.conf), nil
}

// access log filter
func accessLogFunc(alw *AccessLogWriter, alc *AccessLogConfig) filter.Filter {
	return func(c context.Context) {

		start := time.Now()
		c.Next()
		latency := time.Now().Sub(start)
		// build access_log message
		accessLogMsg := buildAccessLogMsg(c, latency)
		if len(accessLogMsg) > 0 {
			alw.Writer(AccessLogData{AccessLogConfig: *alc, AccessLogMsg: accessLogMsg})
		}
	}
}

func buildAccessLogMsg(c context.Context, cost time.Duration) string {
	req := c.(*http.HttpContext).Request
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
	err := c.(*http.HttpContext).Err
	if err != nil {
		builder.WriteString(fmt.Sprintf("invoke err [ %v", err))
		builder.WriteString("] ")
	}
	resp := c.(*http.HttpContext).TargetResp.Data
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
