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

package dubbo

import (
	"encoding/json"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

func coerceDirectInvokeValue(parameterType string, value any) (any, error) {
	return CoerceDirectInvokeValue(parameterType, value)
}

func resolveDirectInvokePayload(irequest config.IntegrationRequest, target *dubboTarget) ([]string, []hessian.Object, []byte, error) {
	if irequest.ParameterTypes == nil {
		return nil, nil, nil, errors.New("direct generic invoke requires parameterTypes")
	}
	if len(irequest.ParameterTypes) == 0 {
		if target == nil || len(target.Values) == 0 {
			return []string{}, []hessian.Object{}, []byte("[]"), nil
		}
		return nil, nil, nil, errors.New("direct generic invoke requires zero values for zero-parameter method")
	}
	if target == nil {
		return nil, nil, nil, errors.New("direct generic invoke requires mapped values")
	}
	if len(target.Values) != len(irequest.ParameterTypes) {
		return nil, nil, nil, errors.New("direct generic invoke requires values to match parameterTypes")
	}

	vals := make([]hessian.Object, len(target.Values))
	for i, value := range target.Values {
		mapped, err := CoerceDirectInvokeValue(irequest.ParameterTypes[i], value)
		if err != nil {
			return nil, nil, nil, err
		}
		vals[i] = mapped
	}

	finalValues, err := json.Marshal(vals)
	if err != nil {
		return nil, nil, nil, err
	}
	return append([]string(nil), irequest.ParameterTypes...), vals, finalValues, nil
}
