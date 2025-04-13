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
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

import (
	"github.com/pkoukk/tiktoken-go"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.LLMTokenizerFilter
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
		cfg *Config
	}
	// Filter is http filter instance
	Filter struct {
		cfg *Config
		tkm *tiktoken.Tiktoken
	}
	// Config describe the config of FilterFactory
	Config struct {
		//OfflineLoader    bool   `yaml:"offline_loader" json:"offline_loader" mapstructure:"offline_loader"`
		Encoding         string `yaml:"encoding" json:"encoding" mapstructure:"encoding"`
		EncodingForModel string `yaml:"encoding_for_model" json:"encoding_for_model" mapstructure:"encoding_for_model"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	c := factory.cfg
	if c == nil {
		return nil
	}

	var tkm *tiktoken.Tiktoken
	var err error

	if c.EncodingForModel != "" {
		tkm, err = tiktoken.EncodingForModel(c.EncodingForModel)
		if err != nil {
			logger.Error("Get tokenizer failed", err)
		}
	} else if c.Encoding != "" {
		tkm, err = tiktoken.GetEncoding(c.Encoding)
		if err != nil {
			logger.Error("Get tokenizer failed", err)
		}
	} else {
		return errors.New("no encoding or model specified")
	}

	f := &Filter{
		cfg: factory.cfg,
		tkm: tkm,
	}
	chain.AppendEncodeFilters(f)
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Encode(hc *http.HttpContext) filter.FilterStatus {
	switch res := hc.TargetResp.(type) {
	case client.StreamResponse:
		go func() {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, res.Stream); err != nil {
				return
			}

			scanner := bufio.NewScanner(&buf)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "data:") {
					line = strings.TrimPrefix(line, "data:")
					token := f.tkm.EncodeOrdinary(line)
					logger.Debugf("[Tokenizer] [DOWNSTREAM] receive response | %d | %s | ", len(token), line)
				}
			}
		}()
	case client.ByteResponse:
		go func() {
			resp := res.Data
			token := f.tkm.EncodeOrdinary(string(resp))
			logger.Debugf("[Tokenizer] [DOWNSTREAM] receive response | %d | %s | ", len(token), resp)
		}()
	default:
		logger.Warnf("Respones type not suitable for token calc")
	}

	return filter.Continue
}

func (f *Filter) Decode(hc *http.HttpContext) filter.FilterStatus {
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, hc.Request.Body); err != nil {
			return
		}

		resp := buf.String()
		token := f.tkm.EncodeOrdinary(resp)
		logger.Debugf("[Tokenizer] [UPSTREAM] receive response | %d | %s | ", len(token), resp)
	}()

	return filter.Continue
}
