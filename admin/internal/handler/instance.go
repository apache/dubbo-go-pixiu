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

package handler

import (
	"fmt"
	"time"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/model"
	"github.com/apache/dubbo-go-pixiu/admin/internal/xds"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/i18n"
)

// xdsServerRef holds reference to xDS server for instance queries.
var xdsServerRef *xds.Server

// RegisterXDSServer registers the xDS server globally for instance queries.
func RegisterXDSServer(server *xds.Server) {
	xdsServerRef = server
}

// ListInstances handles GET /instances
func (h *Handler) ListInstances(c *gin.Context) {
	if xdsServerRef == nil {
		config.ErrorMessageI18n(c, i18n.MsgXDSNotInitialized)
		return
	}

	instances := xdsServerRef.GetInstances()

	// Convert to InstanceInfo for API response
	var infos []model.InstanceInfo
	now := time.Now()
	for _, inst := range instances {
		uptime := ""
		if inst.Status == model.InstanceStatusConnected {
			uptime = formatDuration(now.Sub(inst.ConnectedAt))
		}

		infos = append(infos, model.InstanceInfo{
			NodeID:      inst.NodeID,
			Address:     inst.Address,
			Cluster:     inst.Cluster,
			Version:     inst.Version,
			Metadata:    inst.Metadata,
			Status:      inst.Status,
			LastSeen:    inst.LastSeen.Format(time.RFC3339),
			ConnectedAt: inst.ConnectedAt.Format(time.RFC3339),
			Uptime:      uptime,
		})
	}

	config.Success(c, gin.H{
		"total":     len(infos),
		"instances": infos,
	})
}

// GetInstanceStats handles GET /instances/stats
func (h *Handler) GetInstanceStats(c *gin.Context) {
	if xdsServerRef == nil {
		config.ErrorMessageI18n(c, i18n.MsgXDSNotInitialized)
		return
	}

	instances := xdsServerRef.GetInstances()

	connected := 0
	disconnected := 0
	for _, inst := range instances {
		if inst.Status == model.InstanceStatusConnected {
			connected++
		} else {
			disconnected++
		}
	}

	config.Success(c, gin.H{
		"total":        len(instances),
		"connected":    connected,
		"disconnected": disconnected,
	})
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%dd%dh", days, hours)
}
