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
	"strconv"
	"strings"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/common/util"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type endpointSink interface {
	SetEndpoint(clusterName string, endpoint *model.Endpoint)
	DeleteEndpoint(clusterName string, endpointID string)
}

type clusterManagerEndpointSink struct{}

func (clusterManagerEndpointSink) SetEndpoint(clusterName string, endpoint *model.Endpoint) {
	server.GetClusterManager().SetEndpoint(clusterName, endpoint)
}

func (clusterManagerEndpointSink) DeleteEndpoint(clusterName string, endpointID string) {
	server.GetClusterManager().DeleteEndpoint(clusterName, endpointID)
}

type endpointOwner struct {
	runtimeID string
	registry  string
	serverID  string
}

type publishedEndpoint struct {
	ClusterName string
	EndpointID  string
	Address     model.SocketAddress
}

type endpointReconciler struct {
	mu        sync.Mutex
	sink      endpointSink
	published map[endpointOwner]map[string]publishedEndpoint
}

func newEndpointReconciler(sink endpointSink) *endpointReconciler {
	if sink == nil {
		sink = clusterManagerEndpointSink{}
	}
	return &endpointReconciler{
		sink:      sink,
		published: make(map[endpointOwner]map[string]publishedEndpoint),
	}
}

func (r *endpointReconciler) ApplyServerConfig(runtimeID, registryName, serverID string, cfg *model.McpServerConfig) error {
	owner, err := normalizeEndpointOwner(runtimeID, registryName, serverID)
	if err != nil {
		return err
	}
	desired, err := buildDesiredEndpoints(owner, cfg)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	previous := r.published[owner]
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
		delete(r.published, owner)
		return nil
	}
	r.published[owner] = clonePublishedEndpoints(desired)
	return nil
}

func (r *endpointReconciler) ValidateServerConfig(runtimeID, registryName, serverID string, cfg *model.McpServerConfig) error {
	owner, err := normalizeEndpointOwner(runtimeID, registryName, serverID)
	if err != nil {
		return err
	}
	_, err = buildDesiredEndpoints(owner, cfg)
	return err
}

func normalizeEndpointOwner(runtimeID, registryName, serverID string) (endpointOwner, error) {
	runtimeID = strings.TrimSpace(runtimeID)
	if runtimeID == "" {
		return endpointOwner{}, fmt.Errorf("mcp endpoint reconcile runtime id is required")
	}
	registryName = strings.TrimSpace(registryName)
	if registryName == "" {
		registryName = "default"
	}
	serverID = strings.TrimSpace(serverID)
	if serverID == "" {
		serverID = "default"
	}
	return endpointOwner{runtimeID: runtimeID, registry: registryName, serverID: serverID}, nil
}

func buildDesiredEndpoints(owner endpointOwner, cfg *model.McpServerConfig) (map[string]publishedEndpoint, error) {
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
			EndpointID:  stableEndpointID(owner, name),
			Address: model.SocketAddress{
				Address: result.Host,
				Port:    result.Port,
			},
		}
	}
	return desired, nil
}

func stableEndpointID(owner endpointOwner, toolName string) string {
	return "mcp/" + sanitizeEndpointIDPart(owner.runtimeID) +
		"/" + sanitizeEndpointIDPart(owner.registry) +
		"/" + sanitizeEndpointIDPart(owner.serverID) +
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
			b.WriteString("%")
			b.WriteString(strconv.FormatInt(int64(r), 16))
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
