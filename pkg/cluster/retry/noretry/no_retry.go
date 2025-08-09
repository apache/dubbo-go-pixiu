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

package noretry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	retry.RegisterRetryPolicy(model.RetryerNoRetry, newNoRetryPolicy)
}

type NoRetryPolicy struct {
	firstTime bool
}

func (n *NoRetryPolicy) Attempt() bool {
	if !n.firstTime {
		n.firstTime = true
		return true // Allow the first attempt
	}
	return false
}

func (n *NoRetryPolicy) Reset() {
	n.firstTime = false
}

func newNoRetryPolicy(config map[string]any) (retry.RetryPolicy, error) {
	return &NoRetryPolicy{firstTime: false}, nil
}
