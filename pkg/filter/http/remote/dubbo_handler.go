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
	"regexp"
	"strconv"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	clientdubbo "github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

var mapSourcePattern = regexp.MustCompile(`^(uri|queryStrings|headers|requestBody)\.([\w\d.-]+)$`)

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

	// First map HTTP sources into invocation state, then validate direct mode and types.
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
		// Put the body back because this handler only inspects it while building outbound args.
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

	// opt.* rewrites invoke metadata; numeric MapTo fills positional arguments.
	if strings.HasPrefix(mp.MapTo, "opt.") {
		return h.applyOptMapping(state, mp.MapTo, value, mp.MapType)
	}

	pos, err := strconv.Atoi(strings.TrimSpace(mp.MapTo))
	if err != nil || pos < 0 {
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

func (h *DubboHandler) applyOptMapping(state *outboundBuildState, mapTo string, value any, mapType string) error {
	optKey := strings.TrimSpace(strings.TrimPrefix(mapTo, "opt."))
	switch optKey {
	case "group":
		return applyStringOptValue(&state.group, "Group", value)
	case "version":
		return applyStringOptValue(&state.version, "Version", value)
	case "interface":
		return applyStringOptValue(&state.service, "Interface", value)
	case "method":
		return applyStringOptValue(&state.method, "Method", value)
	case "values":
		values, err := h.normalizeOptValues(value)
		if err != nil {
			return err
		}
		state.optValues = values
		if state.optTypes == nil && strings.TrimSpace(mapType) != "" {
			types, err := h.normalizeOptTypes(mapType)
			if err != nil {
				return err
			}
			state.optTypes = types
		}
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

func applyStringOptValue(target *string, name string, value any) error {
	v, ok := value.(string)
	if !ok {
		return errors.Errorf("%s value is not string", name)
	}
	*target = v
	return nil
}

func (h *DubboHandler) parseMapSource(source string) (string, []string, error) {
	matches := mapSourcePattern.FindStringSubmatch(source)
	if matches == nil {
		return "", nil, errors.New("Parameter mapping config incorrect. Please fix it")
	}
	return matches[1], strings.Split(matches[2], "."), nil
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
	if len(keys) > 0 && keys[0] == constant.DefaultBodyAll {
		return body, nil
	}
	if len(keys) == 0 {
		return nil, errors.New("request body mapping keys are empty")
	}

	current, ok := body[keys[0]]
	if !ok {
		return nil, errors.Errorf("%s does not exist in request body", keys[0])
	}
	if len(keys) == 1 {
		return current, nil
	}

	next, ok := current.(map[string]any)
	if !ok {
		return nil, errors.Errorf("%s is not a map structure. It contains %v", keys[0], current)
	}
	return h.readBodyValue(next, keys[1:])
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
	case nil:
		return nil, nil
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
	case []string:
		types := make([]string, len(v))
		for i, item := range v {
			types[i] = strings.TrimSpace(item)
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
	case nil:
		return nil, nil
	case []any:
		return append([]any(nil), v...), nil
	case []string:
		values := make([]any, len(v))
		for i, item := range v {
			values[i] = item
		}
		return values, nil
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
	rawURL := strings.TrimSpace(ir.URL)
	if rawURL == "" {
		return nil
	}
	if state.serialization == "" {
		return errors.New("direct generic invoke requires serialization")
	}

	if !strings.Contains(rawURL, "://") {
		if state.protocol == "" {
			return errors.New("direct generic invoke requires protocol")
		}
		state.address = rawURL
		return nil
	}

	// Direct URLs bypass registry lookup and must match the declared protocol.
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	directProtocol, err := clientdubbo.DirectURLProtocol(rawURL)
	if err != nil {
		return err
	}
	declaredProtocol := clientdubbo.NormalizeReferenceProtocol(ir.Protocol)
	if declaredProtocol != "" && declaredProtocol != directProtocol {
		return errors.Errorf("direct protocol mismatch: url=%s protocol=%s", directProtocol, declaredProtocol)
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

	// Generic invoke requires values and Java parameter types to stay aligned.
	if ir.ParameterTypes != nil {
		state.paramTypes = append([]string(nil), ir.ParameterTypes...)
		return h.coerceDeclaredArguments(state)
	}
	if state.optTypes != nil {
		state.paramTypes = append([]string(nil), state.optTypes...)
		return h.coerceDeclaredArguments(state)
	}
	if strings.TrimSpace(ir.URL) != "" {
		return errors.New("direct generic invoke requires parameterTypes")
	}

	// Registry mode keeps the historical behavior of inferring omitted Java types.
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
