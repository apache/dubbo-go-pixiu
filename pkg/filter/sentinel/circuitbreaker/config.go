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

package circuitbreaker

import "github.com/alibaba/sentinel-golang/core/circuitbreaker"

type (
	// Resources circuit breaker base config info
	Resources struct {
		Resource                     string                  `json:"resource" yaml:"resource"`
		Strategy                     circuitbreaker.Strategy `json:"strategy" yaml:"strategy"`
		RetryTimeoutMs               uint32                  `json:"retry_timeout_ms" yaml:"retry_timeout_ms"`
		MinRequestAmount             uint64                  `json:"min_request_amount" yaml:"min_request_amount"`
		StatIntervalMs               uint32                  `json:"stat_interval_ms" yaml:"stat_interval_ms"`
		MaxAllowedRtMs               uint64                  `json:"max_allowed_rt_ms" yaml:"max_allowed_rt_ms"`
		StatSlidingWindowBucketCount uint32                  `json:"stat_sliding_window_bucket_count" yaml:"stat_sliding_window_bucket_count"`
		Threshold                    float64                 `json:"threshold" yaml:"threshold"`
		ProbeNum                     uint64                  `json:"probe_num" yaml:"probe_num"`
	}
)
