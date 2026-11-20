/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package configInfo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

type fakeClusterDraftService struct {
	draft        logic.ClusterDraft
	seenObject   schema.AdminObject
	seenSource   int
	seenRevision int64
	updateErr    error
	nameErr      error
	seenDraftID  string
	seenName     string
}

func (f *fakeClusterDraftService) CreateDraft(_ context.Context, object schema.AdminObject) (logic.ClusterDraft, error) {
	f.seenObject = object
	f.draft.Object = object
	return f.draft, nil
}
func (f *fakeClusterDraftService) CreateFromPublished(_ context.Context, source int) (logic.ClusterDraft, error) {
	f.seenSource = source
	return f.draft, nil
}
func (f *fakeClusterDraftService) UpdateDraft(_ context.Context, _ string, object schema.AdminObject, revision int64) (logic.ClusterDraft, error) {
	f.seenObject, f.seenRevision = object, revision
	return f.draft, f.updateErr
}
func (f *fakeClusterDraftService) GetDraft(context.Context, string) (logic.ClusterDraft, error) {
	return f.draft, nil
}
func (f *fakeClusterDraftService) ListDrafts(context.Context) ([]logic.ClusterDraft, error) {
	return []logic.ClusterDraft{f.draft}, nil
}
func (f *fakeClusterDraftService) Status(context.Context, string) (logic.ClusterDraftStatus, error) {
	return logic.ClusterDraftStatus{ID: f.draft.ID, DraftExists: true}, nil
}
func (f *fakeClusterDraftService) Diff(context.Context, string) (logic.ClusterDraftDiff, error) {
	return logic.ClusterDraftDiff{ID: f.draft.ID, Draft: f.draft}, nil
}
func (f *fakeClusterDraftService) ValidateNameAvailability(_ context.Context, draftID, name string) error {
	f.seenDraftID, f.seenName = draftID, name
	return f.nameErr
}
func (f *fakeClusterDraftService) Publish(_ context.Context, id string, revision int64) (logic.ClusterPublishResult, error) {
	f.seenRevision = revision
	return logic.ClusterPublishResult{ID: id, RuntimeID: 17}, nil
}
func (f *fakeClusterDraftService) History(context.Context, string) ([]logic.ClusterHistoryEntry, error) {
	return []logic.ClusterHistoryEntry{}, nil
}
func (f *fakeClusterDraftService) Rollback(_ context.Context, id string, _, revision int64) (logic.ClusterPublishResult, error) {
	f.seenRevision = revision
	return logic.ClusterPublishResult{ID: id, RuntimeID: 17}, nil
}
func (f *fakeClusterDraftService) DeleteDraft(_ context.Context, _ string, revision int64) error {
	f.seenRevision = revision
	return nil
}

type clusterTestResponse struct {
	Code string          `json:"code"`
	Data json.RawMessage `json:"data"`
}

func invokeClusterHandler(t *testing.T, method, target string, body []byte, handler gin.HandlerFunc) clusterTestResponse {
	t.Helper()
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest(method, target, bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request
	handler(ctx)
	require.Equal(t, http.StatusOK, response.Code)
	var result clusterTestResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	return result
}

func withFakeClusterDraftService(t *testing.T, fake *fakeClusterDraftService) {
	t.Helper()
	previous := newClusterDraftService
	newClusterDraftService = func() (clusterDraftService, error) { return fake, nil }
	t.Cleanup(func() { newClusterDraftService = previous })
}

func TestClusterPublishHistoryRollbackAndDeleteHandlers(t *testing.T) {
	id := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	fake := &fakeClusterDraftService{draft: logic.ClusterDraft{ID: id, Revision: 9}}
	withFakeClusterDraftService(t, fake)
	result := invokeClusterHandler(t, http.MethodPut, "/config/api/cluster/publish?id="+id+"&expectedRevision=9", nil, PublishClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	require.Equal(t, int64(9), fake.seenRevision)
	result = invokeClusterHandler(t, http.MethodGet, "/config/api/cluster/history?id="+id, nil, GetClusterHistory)
	require.Equal(t, adminconfig.OK, result.Code)
	result = invokeClusterHandler(t, http.MethodPut, "/config/api/cluster/rollback?id="+id+"&revision=2&expectedRevision=9", nil, RollbackClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	result = invokeClusterHandler(t, http.MethodDelete, "/config/api/cluster/draft?id="+id+"&expectedRevision=9", nil, DeleteClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	for _, handler := range []gin.HandlerFunc{PublishClusterDraft, RollbackClusterDraft, DeleteClusterDraft} {
		result = invokeClusterHandler(t, http.MethodPut, "/config/api/cluster/action?id="+id, nil, handler)
		require.NotEqual(t, adminconfig.OK, result.Code)
	}
}

func TestClusterSchemaEndpointIsIndependentOfRouteSchema(t *testing.T) {
	result := invokeClusterHandler(t, http.MethodGet, "/config/api/cluster/schema", nil, GetClusterSchema)
	require.Equal(t, adminconfig.OK, result.Code)
	var definitions []schema.ObjectSchema
	require.NoError(t, json.Unmarshal(result.Data, &definitions))
	require.Len(t, definitions, 1)
	require.Equal(t, schema.KindAdminCluster, definitions[0].Kind)
}

func TestClusterValidateAndPreviewContracts(t *testing.T) {
	fake := &fakeClusterDraftService{}
	withFakeClusterDraftService(t, fake)
	invalid := []byte(`{"kind":"AdminCluster","metadata":{"name":""},"spec":{"type":"Other","endpoint":{"address":"","port":0}}}`)
	result := invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/preview", invalid, PreviewClusterDraft)
	require.Equal(t, adminconfig.ERR, result.Code)
	var failure clusterErrorResponse
	require.NoError(t, json.Unmarshal(result.Data, &failure))
	paths := make(map[string]bool)
	for _, issue := range failure.Issues {
		paths[issue.Path] = true
	}
	for _, path := range []string{"metadata.name", "spec.type", "spec.endpoint.address", "spec.endpoint.port"} {
		require.True(t, paths[path], path)
	}

	valid := []byte(`{"kind":"AdminCluster","metadata":{"name":"backend"},"spec":{"endpoint":{"address":"127.0.0.1","port":20880}}}`)
	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/validate", valid, ValidateClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	var validated clusterValidationResponse
	require.NoError(t, json.Unmarshal(result.Data, &validated))
	require.True(t, validated.Valid)
	require.Equal(t, "Static", validated.Object.Spec["type"])
	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/preview", valid, PreviewClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	var preview clusterPreviewResponse
	require.NoError(t, json.Unmarshal(result.Data, &preview))
	require.Contains(t, preview.YAML, "port: 20880")
	require.NotContains(t, preview.YAML, "id:")
	require.Equal(t, "backend", fake.seenName)

	fake.nameErr = schema.ValidationErrors{{Path: "metadata.name", Code: "duplicate", Message: "runtime cluster name already exists"}}
	for _, handler := range []gin.HandlerFunc{ValidateClusterDraft, PreviewClusterDraft} {
		result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/check?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", valid, handler)
		require.Equal(t, adminconfig.ERR, result.Code)
		require.NoError(t, json.Unmarshal(result.Data, &failure))
		require.Equal(t, "metadata.name", failure.Issues[0].Path)
		require.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", fake.seenDraftID)
	}
}

func TestClusterDraftControllerSavesInvalidAndUsesRevision(t *testing.T) {
	fake := &fakeClusterDraftService{draft: logic.ClusterDraft{ID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Revision: 7}}
	withFakeClusterDraftService(t, fake)
	invalid := []byte(`{"kind":"AdminCluster","metadata":{"name":""},"spec":{"type":"Other"}}`)
	result := invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/draft", invalid, CreateClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	require.Empty(t, fake.seenObject.Metadata.Name)
	require.Equal(t, "Other", fake.seenObject.Spec["type"])
	result = invokeClusterHandler(t, http.MethodPut, "/config/api/cluster/draft?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa&expectedRevision=7", invalid, UpdateClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	require.Equal(t, int64(7), fake.seenRevision)

	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/draft?sourceId=17", nil, CreateClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	require.Equal(t, 17, fake.seenSource)
	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/draft?sourceId=17", invalid, CreateClusterDraft)
	require.Equal(t, adminconfig.ERR, result.Code)
	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/draft", []byte(`{"object":null}`), CreateClusterDraft)
	require.Equal(t, adminconfig.ERR, result.Code)
	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/draft", []byte(`{"kind":"AdminCluster","metadata":{"name":"x","id":9}}`), CreateClusterDraft)
	require.Equal(t, adminconfig.ERR, result.Code)
	result = invokeClusterHandler(t, http.MethodPost, "/config/api/cluster/draft", []byte("kind: AdminCluster\nmetadata:\n  name: x\n  id: 9\n"), CreateClusterDraft)
	require.Equal(t, adminconfig.ERR, result.Code)
	result = invokeClusterHandler(t, http.MethodGet, "/config/api/cluster/draft/detail?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil, GetClusterDraft)
	require.Equal(t, adminconfig.OK, result.Code)
	result = invokeClusterHandler(t, http.MethodGet, "/config/api/cluster/draft/list", nil, ListClusterDrafts)
	require.Equal(t, adminconfig.OK, result.Code)
	result = invokeClusterHandler(t, http.MethodGet, "/config/api/cluster/status?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil, GetClusterDraftStatus)
	require.Equal(t, adminconfig.OK, result.Code)
	result = invokeClusterHandler(t, http.MethodGet, "/config/api/cluster/diff?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil, GetClusterDraftDiff)
	require.Equal(t, adminconfig.OK, result.Code)

	result = invokeClusterHandler(t, http.MethodPut, "/config/api/cluster/draft?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", invalid, UpdateClusterDraft)
	require.Equal(t, adminconfig.ERR, result.Code)
	require.Equal(t, int64(7), fake.seenRevision)
	fake.updateErr = logic.ErrClusterDraftConflict
	result = invokeClusterHandler(t, http.MethodPut, "/config/api/cluster/draft?id=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa&expectedRevision=7", invalid, UpdateClusterDraft)
	require.Equal(t, adminconfig.RETRY, result.Code)
}

func TestClusterDraftFormYAMLCanSaveIncompleteObject(t *testing.T) {
	fake := &fakeClusterDraftService{draft: logic.ClusterDraft{ID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Revision: 1}}
	withFakeClusterDraftService(t, fake)
	form := url.Values{"content": {"kind: AdminCluster\nmetadata:\n  name: ''\nspec:\n  endpoint:\n    address: ''\n"}}
	request := httptest.NewRequest(http.MethodPost, "/config/api/cluster/draft", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request
	CreateClusterDraft(ctx)
	var result clusterTestResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	require.Equal(t, adminconfig.OK, result.Code)
	require.Empty(t, fake.seenObject.Metadata.Name)
	require.Equal(t, "", fake.seenObject.Spec["endpoint"].(map[string]any)["address"])
}
