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
	"encoding/json"
	"io"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind      = constant.LLMTokenizerFilter
	LoggerFmt = "[Tokenizer] [DOWNSTREAM] "
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
	}
	// Config describe the config of FilterFactory
	Config struct {
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
	f := &Filter{
		cfg: factory.cfg,
	}
	chain.AppendEncodeFilters(f)
	return nil
}

func (f *Filter) Encode(hc *http.HttpContext) filter.FilterStatus {
	switch res := hc.TargetResp.(type) {
	case *client.StreamResponse:
		go f.processStreamResponse(res)
	case *client.UnaryResponse:
		f.processUsageData(res.Data)
	default:
		logger.Infof(LoggerFmt+"Response type not suitable for token calc: %T", res)
	}

	return filter.Continue
}

func (f *Filter) processStreamResponse(res *client.StreamResponse) {
	scanner := bufio.NewScanner(res.Stream)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimPrefix(line, "data:")
			f.processUsageData([]byte(line))
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		logger.Errorf(LoggerFmt+"Error reading stream: %v", err)
	}
}

func (f *Filter) processUsageData(data []byte) {
	var dataCont map[string]interface{}
	err := json.Unmarshal(data, &dataCont)
	if err != nil {
		logger.Infof(LoggerFmt+"Unmarshal response data failed: %v, data: %s", err, string(data))
		return
	}

	usage, ok := dataCont["usage"].(map[string]interface{})
	if !ok || usage == nil {
		logger.Debugf("Usage field not found or is not a map")
		return
	}

	f.logUsage(usage)
}

func (f *Filter) logUsage(usage map[string]interface{}) {
	for key, value := range usage {
		if key == "prompt_tokens_details" {
			promptTokensDetails, ok := value.(map[string]interface{})
			if !ok {
				logger.Warnf(LoggerFmt+"prompt_tokens_details is not a map, value: %+v", value)
				continue
			}
			for detailKey, detailValue := range promptTokensDetails {
				logger.Debugf(LoggerFmt+"Usage | %s: %v", detailKey, detailValue)
			}
		} else {
			logger.Debugf(LoggerFmt+"Usage | %s: %v", key, value)
		}
	}
}
