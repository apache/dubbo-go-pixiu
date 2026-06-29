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

package mcpserver

import (
	"fmt"
	"strings"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/common/util"
	filtermcp "github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type endpointSink interface {
	SetEndpoint(clusterName string, endpoint *model.Endpoint)
	DeleteEndpoint(clusterName, endpointID string)
}

type clusterManagerEndpointSink struct{}

func (clusterManagerEndpointSink) SetEndpoint(clusterName string, endpoint *model.Endpoint) {
	server.GetClusterManager().SetEndpoint(clusterName, endpoint)
}

func (clusterManagerEndpointSink) DeleteEndpoint(clusterName, endpointID string) {
	server.GetClusterManager().DeleteEndpoint(clusterName, endpointID)
}

type publishedEndpoint struct {
	ClusterName string
	EndpointID  string
	Address     model.SocketAddress
}

type endpointReconciler struct {
	mu        sync.Mutex
	sink      endpointSink
	published map[filtermcp.ServerSource]map[string]publishedEndpoint
}

func newEndpointReconciler(sink endpointSink) *endpointReconciler {
	if sink == nil {
		sink = clusterManagerEndpointSink{}
	}
	return &endpointReconciler{
		sink:      sink,
		published: make(map[filtermcp.ServerSource]map[string]publishedEndpoint),
	}
}

func (r *endpointReconciler) ApplyServerConfig(source filtermcp.ServerSource, cfg *model.McpServerConfig) error {
	source = source.Normalize()
	desired, err := buildDesiredEndpoints(source, cfg)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	previous := r.published[source]
	for toolName, old := range previous {
		next, ok := desired[toolName]
		if !ok || old.ClusterName != next.ClusterName || old.EndpointID != next.EndpointID {
			r.sink.DeleteEndpoint(old.ClusterName, old.EndpointID)
		}
	}
	for toolName, next := range desired {
		old, ok := previous[toolName]
		if ok && endpointPublishedEqual(old, next) {
			continue
		}
		endpoint := &model.Endpoint{
			ID:      next.EndpointID,
			Address: next.Address,
		}
		r.sink.SetEndpoint(next.ClusterName, endpoint)
	}

	if len(desired) == 0 {
		delete(r.published, source)
		return nil
	}
	r.published[source] = clonePublishedEndpoints(desired)
	return nil
}

func (r *endpointReconciler) ValidateServerConfig(source filtermcp.ServerSource, cfg *model.McpServerConfig) error {
	_, err := buildDesiredEndpoints(source.Normalize(), cfg)
	return err
}

func (r *endpointReconciler) RemoveSource(source filtermcp.ServerSource) {
	_ = r.ApplyServerConfig(source, nil)
}

func (r *endpointReconciler) RemoveAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for source, previous := range r.published {
		for _, old := range previous {
			r.sink.DeleteEndpoint(old.ClusterName, old.EndpointID)
		}
		delete(r.published, source)
	}
}

func buildDesiredEndpoints(source filtermcp.ServerSource, cfg *model.McpServerConfig) (map[string]publishedEndpoint, error) {
	source = source.Normalize()
	desired := make(map[string]publishedEndpoint)
	if cfg == nil {
		return desired, nil
	}
	for i, tool := range cfg.Tools {
		if tool.BackendURL == "" {
			continue
		}
		name := strings.TrimSpace(tool.Name)
		if name == "" {
			return nil, fmt.Errorf("mcp endpoint reconcile tool at index %d has empty name", i)
		}
		if _, exists := desired[name]; exists {
			return nil, fmt.Errorf("mcp endpoint reconcile duplicate tool %q", name)
		}
		result, err := util.ParseHostPortFromURL(tool.BackendURL)
		if err != nil {
			return nil, fmt.Errorf("parse BackendURL %q for tool %q: %w", tool.BackendURL, tool.Name, err)
		}
		if result.UsedFallback {
			logger.Warnf("[dubbo-go-pixiu] mcp adapter using fallback for tool '%s' with BackendURL '%s': %s",
				tool.Name, tool.BackendURL, result.FallbackInfo)
		}
		desired[name] = publishedEndpoint{
			ClusterName: tool.Cluster,
			EndpointID:  stableEndpointID(source, name),
			Address: model.SocketAddress{
				Address: result.Host,
				Port:    result.Port,
			},
		}
	}
	return desired, nil
}

func stableEndpointID(source filtermcp.ServerSource, toolName string) string {
	return "mcp/" + source.Normalize().Key() +
		"/" + sanitizeEndpointIDPart(toolName)
}

func sanitizeEndpointIDPart(value string) string {
	if value == "" {
		return "_"
	}
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteString("_")
		}
	}
	return b.String()
}

func endpointPublishedEqual(a, b publishedEndpoint) bool {
	return a.ClusterName == b.ClusterName &&
		a.EndpointID == b.EndpointID &&
		a.Address.Equal(b.Address)
}

func clonePublishedEndpoints(in map[string]publishedEndpoint) map[string]publishedEndpoint {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]publishedEndpoint, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
