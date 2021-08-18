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
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.AccessLogFilter
)

func init() {
	extension.RegisterHttpFilter(&AccessPlugin{})
}

type (
	// AccessPlugin is http filter plugin.
	AccessPlugin struct {
	}
	// AccessFilter is http filter instance
	AccessFilter struct {
		cfg *Config
		alw *model.AccessLogWriter
		alc *model.AccessLogConfig
	}
	// Config describe the config of AccessFilter
	Config struct{}
)

func (ap *AccessPlugin) Kind() string {
	return Kind
}

func (ap *AccessPlugin) CreateFilter(hcm *http2.HttpConnectionManager, config interface{}, bs *model.Bootstrap) (extension.HttpFilter, error) {
	alc := bs.StaticResources.AccessLogConfig
	if !alc.Enable {
		return nil, errors.Errorf("AccessPlugin CreateFilter error the access_log config not enable")
	}

	accessLogWriter := &model.AccessLogWriter{AccessLogDataChan: make(chan model.AccessLogData, constant.LogDataBuffer)}
	specConfig := config.(Config)
	return &AccessFilter{cfg: &specConfig, alw: accessLogWriter, alc: &alc}, nil
}

func (af *AccessFilter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(af.Handle)
	return nil
}

func (af *AccessFilter) Handle(c *http.HttpContext) {
	start := time.Now()

	c.Next()

	latency := time.Now().Sub(start)
	// build access_log message
	accessLogMsg := buildAccessLogMsg(c, latency)
	if len(accessLogMsg) > 0 {
		af.alw.Writer(model.AccessLogData{AccessLogConfig: *af.alc, AccessLogMsg: accessLogMsg})
	}
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
