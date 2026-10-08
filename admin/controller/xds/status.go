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

package xds

import (
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
)

// StatusResponse is the operational state of the Admin xDS service. Ready
// means the server is listening and at least one snapshot was published.
// Degraded means snapshot publication or the xDS listener has failed.
type StatusResponse struct {
	adminxds.SnapshotStatus
	ListenPort      uint              `json:"listen_port"`
	Ready           bool              `json:"ready"`
	Degraded        bool              `json:"degraded"`
	ResourceSupport map[string]string `json:"resource_support"`
}

func resourceSupportStatus() map[string]string {
	return map[string]string{
		"extension_config_listener": "supported",
		"extension_config_cluster":  "supported",
		"standard_cds":              "experimental",
		"standard_eds":              "experimental",
		"standard_lds":              "unsupported",
		"standard_rds":              "unsupported",
		"standard_sds":              "unsupported",
		"standard_rtds":             "unsupported",
		"ads_config":                "unsupported",
	}
}

// GetStatus returns the last-good xDS snapshot and latest publication error.
// @Tags XDS
// @Summary get Admin xDS publication status
// @Description Returns xDS listener availability, last-good snapshot metadata, latest errors, readiness, and resource support matrix.
// @Produce application/json
// @Success 200 {object} adminconfig.RetData
// @Router /config/api/xds/status [get]
func GetStatus(c *gin.Context) {
	status := adminxds.DefaultStatusStore.Snapshot()
	response := StatusResponse{
		SnapshotStatus:  status,
		ListenPort:      adminconfig.Bootstrap.GetXDSConfig().ListenPort,
		Ready:           status.Listening && status.SnapshotVersion != "",
		Degraded:        status.LastError != "" || status.ListenError != "",
		ResourceSupport: resourceSupportStatus(),
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(response))
}
