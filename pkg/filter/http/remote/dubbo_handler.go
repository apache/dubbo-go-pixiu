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

package remote

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	clientdubbo "github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

type DubboHandler struct{}

type outboundBuildState struct {
	service       string
	method        string
	group         string
	version       string
	address       string
	protocol      string
	serialization string
	arguments     []any
	paramTypes    []string
	optValues     []any
	optTypes      []string
	hasPositional bool
	body          map[string]any
}

func (h *DubboHandler) BuildOutbound(req *http.Request, api router.API) (*clientdubbo.DubboOutboundRequest, error) {
	state, err := h.newState(req, api)
	if err != nil {
		return nil, err
	}

	for _, mp := range api.IntegrationRequest.MappingParams {
		if err := h.applyMapping(state, req, api, mp); err != nil {
			return nil, err
		}
	}

	if err := h.finalizeDirectAddress(state, api.IntegrationRequest); err != nil {
		return nil, err
	}
	if err := h.finalizeArgumentsAndTypes(state, api.IntegrationRequest); err != nil {
		return nil, err
	}

	return &clientdubbo.DubboOutboundRequest{
		Service:       state.service,
		Method:        state.method,
		Group:         state.group,
		Version:       state.version,
		Address:       state.address,
		Protocol:      state.protocol,
		Serialization: state.serialization,
		Arguments:     append([]any(nil), state.arguments...),
		ParamTypes:    append([]string(nil), state.paramTypes...),
	}, nil
}

func (h *DubboHandler) newState(req *http.Request, api router.API) (*outboundBuildState, error) {
	body := map[string]any{}
	if req != nil && req.Body != nil {
		rawBody, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(rawBody))
		if len(bytes.TrimSpace(rawBody)) > 0 {
			if err := json.Unmarshal(rawBody, &body); err != nil {
				return nil, err
			}
		}
	}

	ir := api.IntegrationRequest
	return &outboundBuildState{
		service:       ir.Interface,
		method:        ir.Method,
		group:         ir.Group,
		version:       ir.Version,
		protocol:      h.resolveDeclaredProtocol(ir),
		serialization: strings.TrimSpace(ir.Serialization),
		body:          body,
	}, nil
}

func (h *DubboHandler) resolveDeclaredProtocol(ir config.IntegrationRequest) string {
	if protocol := clientdubbo.NormalizeReferenceProtocol(ir.Protocol); protocol != "" {
		return protocol
	}
	if protocol := clientdubbo.NormalizeReferenceProtocol(ir.RequestType); protocol != "" {
		return protocol
	}
	return "dubbo"
}

func (h *DubboHandler) applyMapping(state *outboundBuildState, req *http.Request, api router.API, mp config.MappingParam) error {
	value, err := h.readSourceValue(state, req, api, mp.Name)
	if err != nil {
		return err
	}

	if strings.HasPrefix(mp.MapTo, "opt.") {
		return h.applyOptMapping(state, mp.MapTo, value)
	}

	pos, err := strconv.Atoi(strings.TrimSpace(mp.MapTo))
	if err != nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}

	converted, err := clientdubbo.MapTypes(mp.MapType, value)
	if err != nil {
		return err
	}

	if pos >= len(state.arguments) {
		state.arguments = append(state.arguments, make([]any, pos+1-len(state.arguments))...)
		state.paramTypes = append(state.paramTypes, make([]string, pos+1-len(state.paramTypes))...)
	}
	state.arguments[pos] = converted
	state.paramTypes[pos] = strings.TrimSpace(mp.MapType)
	state.hasPositional = true
	return nil
}

func (h *DubboHandler) readSourceValue(state *outboundBuildState, req *http.Request, api router.API, source string) (any, error) {
	from, keys, err := h.parseMapSource(source)
	if err != nil {
		return nil, err
	}

	switch from {
	case constant.QueryStrings:
		return h.readQueryValue(req, keys)
	case constant.Headers:
		return h.readHeaderValue(req, keys)
	case constant.RequestBody:
		return h.readBodyValue(state.body, keys)
	case constant.RequestURI:
		return h.readURIValue(req, api, keys)
	default:
		return nil, errors.Errorf("unsupported mapping source %q", from)
	}
}

func (h *DubboHandler) applyOptMapping(state *outboundBuildState, mapTo string, value any) error {
	optKey := strings.TrimSpace(strings.TrimPrefix(mapTo, "opt."))
	switch optKey {
	case "group":
		v, ok := value.(string)
		if !ok {
			return errors.New("Group value is not string")
		}
		state.group = v
		return nil
	case "version":
		v, ok := value.(string)
		if !ok {
			return errors.New("Version value is not string")
		}
		state.version = v
		return nil
	case "interface":
		v, ok := value.(string)
		if !ok {
			return errors.New("Interface value is not string")
		}
		state.service = v
		return nil
	case "method":
		v, ok := value.(string)
		if !ok {
			return errors.New("Method value is not string")
		}
		state.method = v
		return nil
	case "values":
		values, err := h.normalizeOptValues(value)
		if err != nil {
			return err
		}
		state.optValues = values
		return nil
	case "types":
		types, err := h.normalizeOptTypes(value)
		if err != nil {
			return err
		}
		state.optTypes = types
		return nil
	case "application":
		return errors.Errorf("deprecated opt mapping: %s", mapTo)
	default:
		return errors.Errorf("unknown opt mapping: %s", mapTo)
	}
}

func (h *DubboHandler) parseMapSource(source string) (string, []string, error) {
	return client.ParseMapSource(source)
}

func (h *DubboHandler) readQueryValue(req *http.Request, keys []string) (any, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request url is nil")
	}
	value := req.URL.Query().Get(keys[0])
	if value == "" {
		return nil, errors.Errorf("Query parameter %v does not exist", keys)
	}
	return value, nil
}

func (h *DubboHandler) readHeaderValue(req *http.Request, keys []string) (any, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	value := req.Header.Get(keys[0])
	if value == "" {
		return nil, errors.Errorf("Header %s not found", keys[0])
	}
	return value, nil
}

func (h *DubboHandler) readBodyValue(body map[string]any, keys []string) (any, error) {
	return client.GetMapValue(body, keys)
}

func (h *DubboHandler) readURIValue(req *http.Request, api router.API, keys []string) (any, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request url is nil")
	}
	values := router.GetURIParams(&api, *req.URL)
	if values == nil {
		return nil, errors.Errorf("URI parameter %s not found", keys[0])
	}
	value := values.Get(keys[0])
	if value == "" {
		return nil, errors.Errorf("URI parameter %s not found", keys[0])
	}
	return value, nil
}

func (h *DubboHandler) normalizeOptTypes(value any) ([]string, error) {
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return []string{}, nil
		}
		parts := strings.Split(v, ",")
		types := make([]string, len(parts))
		for i, part := range parts {
			types[i] = strings.TrimSpace(part)
		}
		return types, nil
	case []any:
		types := make([]string, len(v))
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, errors.New("opt.types must be string or string array")
			}
			types[i] = strings.TrimSpace(s)
		}
		return types, nil
	default:
		return nil, errors.New("opt.types must be string or string array")
	}
}

func (h *DubboHandler) normalizeOptValues(value any) ([]any, error) {
	switch v := value.(type) {
	case []any:
		return append([]any(nil), v...), nil
	case string:
		if v == "" {
			return []any{}, nil
		}
		return []any{v}, nil
	default:
		return []any{value}, nil
	}
}

func (h *DubboHandler) finalizeDirectAddress(state *outboundBuildState, ir config.IntegrationRequest) error {
	if strings.TrimSpace(ir.URL) == "" {
		return nil
	}
	if state.serialization == "" {
		return errors.New("direct generic invoke requires serialization")
	}

	u, err := url.Parse(strings.TrimSpace(ir.URL))
	if err != nil {
		return err
	}

	directProtocol, err := clientdubbo.DirectURLProtocol(ir.URL)
	if err != nil {
		return err
	}
	if state.protocol != "" && state.protocol != directProtocol {
		return errors.Errorf("direct protocol mismatch: url=%s protocol=%s", directProtocol, state.protocol)
	}

	state.protocol = directProtocol
	state.address = u.Host
	return nil
}

func (h *DubboHandler) finalizeArgumentsAndTypes(state *outboundBuildState, ir config.IntegrationRequest) error {
	if state.hasPositional && state.optValues != nil {
		return errors.New("positional mappings and opt.values are mutually exclusive")
	}
	if state.optValues != nil {
		state.arguments = append([]any(nil), state.optValues...)
	}

	if ir.ParameterTypes != nil {
		state.paramTypes = append([]string(nil), ir.ParameterTypes...)
		return h.coerceDeclaredArguments(state)
	}
	if state.optTypes != nil {
		state.paramTypes = append([]string(nil), state.optTypes...)
		return h.coerceDeclaredArguments(state)
	}

	inferred := clientdubbo.InferJavaClassNames(state.arguments)
	if len(inferred) < len(state.arguments) {
		inferred = append(inferred, make([]string, len(state.arguments)-len(inferred))...)
	}

	paramTypes := make([]string, len(state.arguments))
	for i := range state.arguments {
		if i < len(state.paramTypes) && strings.TrimSpace(state.paramTypes[i]) != "" {
			paramTypes[i] = strings.TrimSpace(state.paramTypes[i])
			continue
		}
		paramTypes[i] = inferred[i]
	}
	state.paramTypes = paramTypes
	return nil
}

func (h *DubboHandler) coerceDeclaredArguments(state *outboundBuildState) error {
	if len(state.arguments) != len(state.paramTypes) {
		return errors.New("direct generic invoke requires values to match parameterTypes")
	}

	values := make([]any, len(state.arguments))
	for i, value := range state.arguments {
		mapped, err := clientdubbo.CoerceDirectInvokeValue(state.paramTypes[i], value)
		if err != nil {
			return err
		}
		values[i] = mapped
	}
	state.arguments = values
	return nil
}
