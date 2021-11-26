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

package seata

import (
	netHttp "net/http"
	"strings"
)

import (
	"github.com/opentrx/seata-golang/v2/pkg/apis"
	"github.com/opentrx/seata-golang/v2/pkg/util/runtime"

	"google.golang.org/grpc"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	Kind = "dgp.filter.http.seata"

	SEATA    = "seata"
	XID      = "x_seata_xid"
	BranchID = "x_seata_branch_id"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// MetricFilter is http Filter plugin.
	Plugin struct {
	}
	// Filter is http Filter instance
	Filter struct {
		conf              *Seata
		transactionInfos  map[string]*TransactionInfo
		tccResources      map[string]*TCCResource
		transactionClient apis.TransactionManagerServiceClient
		resourceClient    apis.ResourceManagerServiceClient
		branchMessages    chan *apis.BranchMessage
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{
		conf:             &Seata{},
		transactionInfos: make(map[string]*TransactionInfo),
		tccResources:     make(map[string]*TCCResource),
		branchMessages:   make(chan *apis.BranchMessage),
	}, nil
}

func (f *Filter) Config() interface{} {
	return f.conf
}

func (f *Filter) Apply() error {
	conn, err := grpc.Dial(f.conf.ServerAddressing,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(f.conf.GetClientParameters()))
	if err != nil {
		return err
	}

	f.transactionClient = apis.NewTransactionManagerServiceClient(conn)
	f.resourceClient = apis.NewResourceManagerServiceClient(conn)

	runtime.GoWithRecover(func() {
		f.branchCommunicate()
	}, nil)

	for _, ti := range f.conf.TransactionInfos {
		f.transactionInfos[ti.RequestPath] = ti
	}

	for _, r := range f.conf.TCCResources {
		f.tccResources[r.PrepareRequestPath] = r
	}
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	if method != netHttp.MethodPost {
		ctx.Next()
		return
	}

	transactionInfo, found := f.transactionInfos[strings.ToLower(path)]
	if found {
		result := f.handleHttp1GlobalBegin(ctx, transactionInfo)
		ctx.Next()
		if result {
			f.handleHttp1GlobalEnd(ctx)
		}
	} else {
		ctx.Next()
	}

	tccResource, exists := f.tccResources[strings.ToLower(path)]
	if exists {
		result := f.handleHttp1BranchRegister(ctx, tccResource)
		ctx.Next()
		if result {
			f.handleHttp1BranchEnd(ctx)
		}
	} else {
		ctx.Next()
	}
}
