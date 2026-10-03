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

package router

import (
	"fmt"
	"sort"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/copyutil"
	"github.com/apache/dubbo-go-pixiu/pkg/common/fingerprint"
)

func CloneClaims(claims map[string]any) map[string]any {
	return copyutil.CloneStringAnyMap(claims)
}

func ClaimsFingerprint(sc SelectionContext) (string, error) {
	b := fingerprint.NewBuilder()
	if sc.UserID != "" {
		b.AddString("sub:" + sc.UserID)
	}
	if sc.Tenant != "" {
		b.AddString("tenant:" + sc.Tenant)
	}
	if len(sc.Claims) == 0 {
		return b.SumShort(16), nil
	}
	keys := make([]string, 0, len(sc.Claims))
	for k := range sc.Claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if (k == "sub" && sc.UserID != "") || (k == "tenant" && sc.Tenant != "") {
			continue
		}
		b.AddString("claim:" + k)
		if err := b.AddJSON(sc.Claims[k]); err != nil {
			return "", fmt.Errorf("claim %q is not JSON-canonicalizable: %w", k, err)
		}
	}
	return b.SumShort(16), nil
}
