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
	"google.golang.org/grpc/credentials/insecure"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
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
	// FilterFactory is http Filter instance
	FilterFactory struct {
		conf              *Seata
		transactionInfos  map[string]*TransactionInfo
		tccResources      map[string]*TCCResource
		transactionClient apis.TransactionManagerServiceClient
		resourceClient    apis.ResourceManagerServiceClient
		branchMessages    chan *apis.BranchMessage
	}
	// Filter is http Filter instance
	Filter struct {
		conf              *Seata
		transactionInfos  map[string]*TransactionInfo
		tccResources      map[string]*TCCResource
		transactionClient apis.TransactionManagerServiceClient
		resourceClient    apis.ResourceManagerServiceClient
		branchMessages    chan *apis.BranchMessage

		globalResult         bool
		branchRegisterResult bool
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		conf:             &Seata{},
		transactionInfos: make(map[string]*TransactionInfo),
		tccResources:     make(map[string]*TCCResource),
		branchMessages:   make(chan *apis.BranchMessage),
	}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.conf
}

func (factory *FilterFactory) Apply() error {
	conn, err := grpc.Dial(factory.conf.ServerAddressing,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(factory.conf.GetClientParameters()))
	if err != nil {
		return err
	}

	factory.transactionClient = apis.NewTransactionManagerServiceClient(conn)
	factory.resourceClient = apis.NewResourceManagerServiceClient(conn)

	runtime.GoWithRecover(func() {
		factory.branchCommunicate()
	}, nil)

	for _, ti := range factory.conf.TransactionInfos {
		factory.transactionInfos[ti.RequestPath] = ti
	}

	for _, r := range factory.conf.TCCResources {
		factory.tccResources[r.PrepareRequestPath] = r
	}
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		conf:              factory.conf,
		transactionInfos:  factory.transactionInfos,
		tccResources:      factory.tccResources,
		transactionClient: factory.transactionClient,
		resourceClient:    factory.resourceClient,
		branchMessages:    factory.branchMessages,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	if method != netHttp.MethodPost {
		return filter.Continue
	}

	transactionInfo, found := f.transactionInfos[strings.ToLower(path)]
	if found {
		f.globalResult = f.handleHttp1GlobalBegin(ctx, transactionInfo)
	}

	tccResource, exists := f.tccResources[strings.ToLower(path)]
	if exists {
		f.branchRegisterResult = f.handleHttp1BranchRegister(ctx, tccResource)
	}
	return filter.Continue
}

func (f *Filter) Encode(ctx *http.HttpContext) filter.FilterStatus {
	if f.globalResult {
		f.handleHttp1GlobalEnd(ctx)
	}

	if f.branchRegisterResult {
		f.handleHttp1BranchEnd(ctx)
	}
	return filter.Continue
}
