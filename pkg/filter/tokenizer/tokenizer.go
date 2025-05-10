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
	"fmt"
	"io"
	"strings"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind                = constant.LLMTokenizerFilter
	LoggerFmt           = "[Tokenizer] [DOWNSTREAM] "
	PromptTokensDetails = "prompt_tokens_details"
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

func (factory *FilterFactory) Config() any {
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
		pr, pw := io.Pipe()
		res.Stream = newTeeReadCloser(res.Stream, pw)
		go f.processStreamResponse(pr)
	case *client.UnaryResponse:
		f.processUsageData(res.Data)
	default:
		logger.Warnf(LoggerFmt+"Response type not suitable for token calc: %T", res)
	}

	return filter.Continue
}

func (f *Filter) processStreamResponse(stream io.Reader) {
	scanner := bufio.NewScanner(stream)
	currentLine := make([]byte, 0, 1024)
	// read the stream by line
	// and process the data lines
	// the data line is prefixed with "data:"
	// the data line is a json string
	// the for loop is to read the streamline by line and concat the separate "data:" lines
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			f.processUsageData(currentLine)
			currentLine = make([]byte, 0, 1024)
			line = strings.TrimPrefix(line, "data:")
		}
		currentLine = append(currentLine, line...)
	}
	f.processUsageData(currentLine)
	if err := scanner.Err(); err != nil && err != io.EOF {
		logger.Errorf(LoggerFmt+"Error reading stream: %v", err)
	}
}

func (f *Filter) processUsageData(data []byte) {
	var dataCont map[string]any
	err := json.Unmarshal(data, &dataCont)
	if err != nil {
		return
	}

	usage, ok := dataCont["usage"].(map[string]any)
	if !ok || usage == nil {
		return
	}

	// todo: currently we only log the usage, we should export it to metrics
	f.logUsage(usage)
}

func (f *Filter) logUsage(usage map[string]any) {
	for key, value := range usage {
		if key == PromptTokensDetails {
			details, ok := value.(map[string]any)
			if !ok {
				logger.Warnf(LoggerFmt+PromptTokensDetails+" is not a map, value: %+v", value)
				continue
			}
			for detailKey, detailValue := range details {
				logger.Infof(LoggerFmt+"Usage | %s: %v", detailKey, detailValue)
			}
		} else {
			logger.Infof(LoggerFmt+"Usage | %s: %v", key, value)
		}
	}
}

type teeReadCloser struct {
	reader   io.Reader
	closer   io.Closer
	writer   io.Writer
	once     sync.Once
	closeErr error
}

func newTeeReadCloser(r io.ReadCloser, w io.Writer) *teeReadCloser {
	return &teeReadCloser{
		reader: r,
		closer: r,
		writer: w,
	}
}

func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	if n <= 0 || err != nil {
		return
	}
	nw, err := t.writer.Write(p[:n])
	if err != nil {
		logger.Errorf(LoggerFmt+"Error writing to tee writer: %v", err)
		return
	}
	if nw != n {
		logger.Errorf(LoggerFmt+"Short write to tee writer: %d/%d", nw, n)
		//err = fmt.Errorf("short write to tee writer: %d/%d", nw, n)
	}
	return n, nil
}

func (t *teeReadCloser) Close() (err error) {
	var (
		closerErr error
		writerErr error
	)

	t.once.Do(func() {
		closerErr = t.closer.Close()
		if closerErr != nil {
			logger.Errorf(LoggerFmt+"Error closing closer: %v", closerErr)
		}

		if t.writer != nil {
			writerCloser, ok := t.writer.(io.Closer)
			if ok {
				writerErr = writerCloser.Close()
				if writerErr != nil {
					logger.Errorf(LoggerFmt+"Error closing writer: %v", writerErr)
				}
			}
		}
	})

	if closerErr != nil || writerErr != nil {
		err = fmt.Errorf("closing closer error: %w. closing writer error: %w", closerErr, writerErr)
	}
	return err
}
