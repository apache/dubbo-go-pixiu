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

package metric

type (
	// Rules router match
	Rules struct {
		Match       Match       `yaml:"match" json:"match" mapstructure:"match"`                      // router
		Requirement Requirement `yaml:"require_ment" json:"require_ment" mapstructure:"require_ment"` // router requires jwks check
	}

	Match struct {
		Prefix string `yaml:"prefix" json:"prefix" mapstructure:"prefix"` // url
	}
	Requirement struct {
		ProviderName string `yaml:"provider_name" json:"provider_name" mapstructure:"provider_name"` // jwks providers name
	}
	ApiStatsResponse struct {
		Name     string    `yaml:"name" json:"name" mapstructure:"name"`
		ApiStats []ApiStat `yaml:"api_stats" json:"api_stats" mapstructure:"api_stats"`
	}

	ClusterStatsResponse struct {
		Name         string        `yaml:"name" json:"name" mapstructure:"name"`
		ClusterStats []ClusterStat `yaml:"cluster_stats" json:"cluster_stats" mapstructure:"cluster_stats"`
	}

	CommonStatsResponse struct {
		Name         string       `yaml:"name" json:"name" mapstructure:"name"`
		BackendStats BackendStat  `yaml:"backend_stats" json:"backend_stats" mapstructure:"backend_stats"`
		FrontedStats FrontendStat `yaml:"fronted_stats" json:"fronted_stats" mapstructure:"fronted_stats"`
		ServerStats  ServerStat   `yaml:"server_stats"  json:"server_stats"  mapstructure:"server_stats"`
	}

	ClusterStat struct {
		ClusterName                 string  `yaml:"cluster_name" json:"cluster_name" mapstructure:"cluster_name"`
		Status                      string  `yaml:"status" json:"status" mapstructure:"status"`
		TimedOut                    bool    `yaml:"time_out" json:"time_out" mapstructure:"time_out"`
		NumberOfNodes               int     `yaml:"number_of_nodes" json:"number_of_nodes" mapstructure:"number_of_nodes"`
		NumberOfDataNodes           int     `yaml:"number_of_data_nodes" json:"number_of_data_nodes" mapstructure:"number_of_data_nodes"`
		ActivePrimaryShards         int     `yaml:"active_primary_shards" json:"active_primary_shards" mapstructure:"active_primary_shards"`
		ActiveShards                int     `yaml:"active_shards" json:"active_shards" mapstructure:"active_shards"`
		RelocatingShards            int     `yaml:"relocating_shards" json:"relocating_shards" mapstructure:"relocating_shards"`
		InitializingShards          int     `yaml:"initializing_shards" json:"initializing_shards" mapstructure:"initializing_shards"`
		UnassignedShards            int     `yaml:"unassigned_shards" json:"unassigned_shards" mapstructure:"unassigned_shards"`
		DelayedUnassignedShards     int     `yaml:"delayed_unassigned_shards" json:"delayed_unassigned_shards" mapstructure:"delayed_unassigned_shards"`
		NumberOfPendingTasks        int     `yaml:"number_of_pending_tasks" json:"number_of_pending_tasks" mapstructure:"number_of_pending_tasks"`
		NumberOfInFlightFetch       int     `yaml:"number_of_in_flight_fetch" json:"number_of_in_flight_fetch" mapstructure:"number_of_in_flight_fetch"`
		TaskMaxWaitingInQueueMillis int     `yaml:"task_max_waiting_in_queue_millis" json:"task_max_waiting_in_queue_millis" mapstructure:"task_max_waiting_in_queue_millis"`
		ActiveShardsPercentAsNumber float64 `yaml:"active_shards_percent_as_number" json:"active_shards_percent_as_number" mapstructure:"active_shards_percent_as_number"`
	}
	ApiStat struct {
		ApiName            string `yaml:"api_name" json:"api_name" mapstructure:"api_name"`
		ApiRequests        int64  `yaml:"api_requests" json:"api_requests" mapstructure:"api_requests"`
		ApiRequestsLatency int64  `yaml:"api_requests_latency" json:"api_requests_latency" mapstructure:"api_requests_latency"`
	}

	FrontendStat struct {
		FrontName string `yaml:"front_name" json:"front_name" mapstructure:"front_name"`
	}
	BackendStat struct {
		BackendName string `yaml:"backend_name" json:"backend_name" mapstructure:"backend_name"`
	}

	ServerStat struct {
		Backend string `yaml:"backend" json:"backend" mapstructure:"backend"`
		Server  string `yaml:"server" json:"server" mapstructure:"server"`
	}
)
