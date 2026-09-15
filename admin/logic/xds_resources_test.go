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

package logic

import (
	"testing"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
)

func TestDecodeXDSResourcesPropagatesMalformedYAML(t *testing.T) {
	previous := adminconfig.Bootstrap
	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{EtcdConfig: adminconfig.EtcdConfig{Path: "/pixiu/config/api"}}
	t.Cleanup(func() { adminconfig.Bootstrap = previous })

	_, err := decodeClusters(
		[]string{getClusterKey("1")},
		[]string{"name: ["},
	)
	require.ErrorContains(t, err, "decode cluster configuration")

	_, err = decodeListeners(
		[]string{getListenerKey("http")},
		[]string{"name: ["},
	)
	require.ErrorContains(t, err, "decode listener configuration")
}

func TestDecodeXDSResourcesKeepsLegitimateEmptySets(t *testing.T) {
	clusters, err := decodeClusters(nil, nil)
	require.NoError(t, err)
	require.Empty(t, clusters)

	listeners, err := decodeListeners(nil, nil)
	require.NoError(t, err)
	require.Empty(t, listeners)
}
