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

package initialize

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	xdsController "github.com/apache/dubbo-go-pixiu/admin/controller/xds"
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
)

func TestXDSStatusRouteEndToEnd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousBootstrap := adminconfig.Bootstrap
	previousStore := adminxds.DefaultStatusStore
	t.Cleanup(func() {
		gin.SetMode(gin.DebugMode)
		adminconfig.Bootstrap = previousBootstrap
		adminxds.DefaultStatusStore = previousStore
	})

	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{XDS: adminconfig.XDSConfig{
		ListenPort: 19000,
		NodeID:     "gateway-a",
	}}
	adminxds.DefaultStatusStore = adminxds.NewStatusStore("gateway-a")
	adminxds.DefaultStatusStore.RecordListening()
	adminxds.DefaultStatusStore.RecordSuccess("42", 2, 3)
	adminxds.DefaultStatusStore.RecordError(errors.New("candidate rejected"))

	req := httptest.NewRequest(http.MethodGet, "/config/api/xds/status", nil)
	req.Header.Set("token", signToken(t))
	recorder := httptest.NewRecorder()
	Routers().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	var response struct {
		Code string                       `json:"code"`
		Data xdsController.StatusResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, adminconfig.OK, response.Code)
	require.Equal(t, "gateway-a", response.Data.NodeID)
	require.Equal(t, "42", response.Data.SnapshotVersion)
	require.Equal(t, uint(19000), response.Data.ListenPort)
	require.True(t, response.Data.Listening)
	require.True(t, response.Data.Ready)
	require.True(t, response.Data.Degraded)
	require.Equal(t, "candidate rejected", response.Data.LastError)
	require.Equal(t, "supported", response.Data.ResourceSupport["extension_config_cluster"])
	require.Equal(t, "experimental", response.Data.ResourceSupport["standard_cds"])
	require.Equal(t, "unsupported", response.Data.ResourceSupport["standard_lds"])

	adminxds.DefaultStatusStore.RecordListenError(errors.New("address already in use"))
	recorder = httptest.NewRecorder()
	Routers().ServeHTTP(recorder, req.Clone(req.Context()))
	require.Equal(t, http.StatusOK, recorder.Code)
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.False(t, response.Data.Listening)
	require.False(t, response.Data.Ready)
	require.True(t, response.Data.Degraded)
	require.Equal(t, "address already in use", response.Data.ListenError)
}

func TestSwaggerIncludesXDSStatusRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	recorder := httptest.NewRecorder()
	Routers().ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	var document struct {
		Paths map[string]json.RawMessage `json:"paths"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &document))
	require.Contains(t, document.Paths, "/config/api/xds/status")
}
