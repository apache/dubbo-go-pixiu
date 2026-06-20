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
	"context"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func authorizeCall(cs *CompositeSelector, ctx context.Context, sc SelectionContext, tools []model.ToolConfig) error {
	_, err := cs.AuthorizeCall(ctx, sc, tools)
	return err
}

func recordCallSuccessForTest(t *testing.T, cs *CompositeSelector, sc SelectionContext) (CallSuccessResult, error) {
	t.Helper()
	plan, ok := cs.store.Get(cs.planKey(sc.SessionID))
	if !ok {
		return cs.RecordCallSuccess(context.Background(), AuthorizationReceipt{})
	}
	receipt, err := cs.store.IssueReceipt(cs.planKey(sc.SessionID), sc.Requested, plan, cs.routerID)
	if err != nil {
		return cs.RecordCallSuccess(context.Background(), AuthorizationReceipt{})
	}
	return cs.RecordCallSuccess(context.Background(), receipt)
}
