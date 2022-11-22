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
	"bytes"
	"context"
	"fmt"
	"io"
	netHttp "net/http"
	"strconv"
)

import (
	"github.com/opentrx/seata-golang/v2/pkg/apis"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

// handleHttp1GlobalBegin return bool, represent whether continue
func (f *Filter) handleHttp1GlobalBegin(ctx *http.HttpContext, transactionInfo *TransactionInfo) bool {
	// todo support transaction isolation level
	xid, err := f.globalBegin(ctx.Ctx, transactionInfo.RequestPath, transactionInfo.Timeout)
	if err != nil {
		logger.Errorf("failed to begin global transaction, transaction info: %v, err: %v",
			transactionInfo, err)
		ctx.SendLocalReply(netHttp.StatusInternalServerError, []byte(fmt.Sprintf("failed to begin global transaction, %v", err)))
		return false
	}
	ctx.Params[XID] = xid
	ctx.Request.Header.Add(XID, xid)
	return true
}

func (f *Filter) handleHttp1GlobalEnd(ctx *http.HttpContext) {
	xidParam := ctx.Params[XID]
	xid := xidParam.(string)
	response, ok := ctx.SourceResp.(*netHttp.Response)
	if ok {
		if response.StatusCode == netHttp.StatusOK {
			err := f.globalCommit(ctx, xid)
			if err != nil {
				logger.Error(err)
			}
		} else {
			err := f.globalRollback(ctx, xid)
			if err != nil {
				logger.Error(err)
			}
		}
	} else {
		err := f.globalRollback(ctx, xid)
		if err != nil {
			logger.Error(err)
		}
	}
}

// handleHttp1BranchRegister return bool, represent whether continue
func (f *Filter) handleHttp1BranchRegister(hctx *http.HttpContext, tccResource *TCCResource) bool {
	xid := hctx.Request.Header.Get(XID)
	if xid == "" {
		logger.Error("failed to get xid from request header")
		hctx.SendLocalReply(netHttp.StatusInternalServerError, []byte("failed to get xid from request header"))
		return false
	}

	bodyBytes, err := io.ReadAll(hctx.Request.Body)
	if err != nil {
		logger.Error(err)
		hctx.SendLocalReply(netHttp.StatusInternalServerError, []byte("failed to retrieve request body"))
		return false
	}

	requestContext := &RequestContext{
		ActionContext: make(map[string]string),
		Headers:       hctx.Request.Header.Clone(),
		Body:          bodyBytes,
		Trailers:      hctx.Request.Trailer.Clone(),
	}
	// Once read body, the body sawEOF will be true, then send request will have no body
	hctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	rEntry := hctx.GetRouteEntry()
	if rEntry == nil {
		panic("no route entry")
	}
	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster :%v", rEntry.Cluster)

	clusterName := rEntry.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		hctx.SendLocalReply(netHttp.StatusServiceUnavailable, []byte("cluster not found endpoint"))
		return false
	}

	requestContext.ActionContext[VarHost] = endpoint.Address.GetAddress()
	requestContext.ActionContext[CommitRequestPath] = tccResource.CommitRequestPath
	requestContext.ActionContext[RollbackRequestPath] = tccResource.RollbackRequestPath
	queryString := hctx.Request.URL.RawQuery
	if queryString != "" {
		requestContext.ActionContext[VarQueryString] = queryString
	}

	data, err := requestContext.Encode()
	if err != nil {
		logger.Errorf("encode request context failed, request context: %v, err: %v", requestContext, err)
		hctx.SendLocalReply(netHttp.StatusInternalServerError, []byte(fmt.Sprintf("encode request context failed, %v", err)))
		return false
	}
	ctx, cancel := context.WithTimeout(hctx.Ctx, hctx.Timeout)
	defer cancel()
	branchID, err := f.branchRegister(ctx, xid, tccResource.PrepareRequestPath, apis.TCC, data, "")
	if err != nil {
		logger.Errorf("branch transaction register failed, xid: %s, err: %v", xid, err)
		hctx.SendLocalReply(netHttp.StatusInternalServerError, []byte(fmt.Sprintf("branch transaction register failed, %v", err)))
		return false
	}
	hctx.Params[XID] = xid
	hctx.Params[BranchID] = strconv.FormatInt(branchID, 10)
	return true
}

func (f *Filter) handleHttp1BranchEnd(ctx *http.HttpContext) {
	xidParam := ctx.Params[XID]
	xid := xidParam.(string)
	branchIDParam := ctx.Params[BranchID]
	branchID, err := strconv.ParseInt(branchIDParam.(string), 10, 64)
	if err != nil {
		logger.Error(err)
	}
	response, ok := ctx.SourceResp.(*netHttp.Response)
	if ok {
		if response.StatusCode != netHttp.StatusOK {
			err := f.branchReport(ctx.Ctx, xid, branchID, apis.TCC, apis.PhaseOneFailed, nil)
			if err != nil {
				logger.Error(err)
			}
		}
	} else {
		err := f.branchReport(ctx.Ctx, xid, branchID, apis.TCC, apis.PhaseOneFailed, nil)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (f *Filter) globalBegin(ctx context.Context, name string, timeout int32) (string, error) {
	request := &apis.GlobalBeginRequest{
		Addressing:      f.conf.Addressing,
		Timeout:         timeout,
		TransactionName: name,
	}
	resp, err := f.transactionClient.Begin(ctx, request)
	if err != nil {
		return "", err
	}
	if resp.ResultCode == apis.ResultCodeSuccess {
		return resp.XID, nil
	}
	return "", fmt.Errorf(resp.Message)
}

func (f *Filter) globalCommit(ctx *http.HttpContext, xid string) error {
	var (
		err    error
		status apis.GlobalSession_GlobalStatus
	)
	defer func() {
		delete(ctx.Params, XID)
	}()
	retry := f.conf.CommitRetryCount
	for retry > 0 {
		status, err = f.commit(ctx.Ctx, xid)
		if err != nil {
			logger.Errorf("failed to report global commit [%s],Retry Countdown: %d, reason: %s",
				xid, retry, err.Error())
		} else {
			break
		}
		retry--
		if retry == 0 {
			return errors.New("failed to report global commit")
		}
	}
	logger.Infof("[%s] commit status: %s", xid, status.String())
	return nil
}

func (f *Filter) globalRollback(ctx *http.HttpContext, xid string) error {
	var (
		err    error
		status apis.GlobalSession_GlobalStatus
	)
	defer func() {
		delete(ctx.Params, XID)
	}()
	retry := f.conf.RollbackRetryCount
	for retry > 0 {
		status, err = f.rollback(ctx.Ctx, xid)
		if err != nil {
			logger.Errorf("failed to report global rollback [%s],Retry Countdown: %d, reason: %s",
				xid, retry, err.Error())
		} else {
			break
		}
		retry--
		if retry == 0 {
			return errors.New("failed to report global rollback")
		}
	}
	logger.Infof("[%s] rollback status: %s", xid, status.String())
	return nil
}

func (f *Filter) commit(ctx context.Context, xid string) (apis.GlobalSession_GlobalStatus, error) {
	request := &apis.GlobalCommitRequest{XID: xid}
	resp, err := f.transactionClient.Commit(ctx, request)
	if err != nil {
		return 0, err
	}
	if resp.ResultCode == apis.ResultCodeSuccess {
		return resp.GlobalStatus, nil
	}
	return 0, errors.New(resp.Message)
}

func (f *Filter) rollback(ctx context.Context, xid string) (apis.GlobalSession_GlobalStatus, error) {
	request := &apis.GlobalRollbackRequest{XID: xid}
	resp, err := f.transactionClient.Rollback(ctx, request)
	if err != nil {
		return 0, err
	}
	if resp.ResultCode == apis.ResultCodeSuccess {
		return resp.GlobalStatus, nil
	}
	return 0, errors.New(resp.Message)
}

func (f *Filter) branchRegister(ctx context.Context, xid string, resourceID string,
	branchType apis.BranchSession_BranchType, applicationData []byte, lockKeys string) (int64, error) {
	request := &apis.BranchRegisterRequest{
		Addressing:      f.conf.Addressing,
		XID:             xid,
		ResourceID:      resourceID,
		LockKey:         lockKeys,
		BranchType:      branchType,
		ApplicationData: applicationData,
	}
	resp, err := f.resourceClient.BranchRegister(ctx, request)
	if err != nil {
		return 0, err
	}
	if resp.ResultCode == apis.ResultCodeSuccess {
		return resp.BranchID, nil
	} else {
		return 0, fmt.Errorf(resp.Message)
	}
}

func (f *Filter) branchReport(ctx context.Context, xid string, branchID int64,
	branchType apis.BranchSession_BranchType, status apis.BranchSession_BranchStatus, applicationData []byte) error {
	request := &apis.BranchReportRequest{
		XID:             xid,
		BranchID:        branchID,
		BranchType:      branchType,
		BranchStatus:    status,
		ApplicationData: applicationData,
	}
	resp, err := f.resourceClient.BranchReport(ctx, request)
	if err != nil {
		return err
	}
	if resp.ResultCode == apis.ResultCodeFailed {
		return fmt.Errorf(resp.Message)
	}
	return nil
}
