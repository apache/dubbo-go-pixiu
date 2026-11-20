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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

import (
	"github.com/gin-gonic/gin"

	"gopkg.in/yaml.v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

type clusterDraftService interface {
	CreateDraft(context.Context, schema.AdminObject) (logic.ClusterDraft, error)
	CreateFromPublished(context.Context, int) (logic.ClusterDraft, error)
	UpdateDraft(context.Context, string, schema.AdminObject, int64) (logic.ClusterDraft, error)
	GetDraft(context.Context, string) (logic.ClusterDraft, error)
	ListDrafts(context.Context) ([]logic.ClusterDraft, error)
	Status(context.Context, string) (logic.ClusterDraftStatus, error)
	Diff(context.Context, string) (logic.ClusterDraftDiff, error)
	ValidateNameAvailability(context.Context, string, string) error
	Publish(context.Context, string, int64) (logic.ClusterPublishResult, error)
	History(context.Context, string) ([]logic.ClusterHistoryEntry, error)
	Rollback(context.Context, string, int64, int64) (logic.ClusterPublishResult, error)
	DeleteDraft(context.Context, string, int64) error
}

var newClusterDraftService = func() (clusterDraftService, error) { return logic.NewAdminClusterDraftStore() }

type clusterDraftRequest struct {
	Object           *schema.AdminObject `json:"object"`
	ExpectedRevision int64               `json:"expectedRevision"`
}

type clusterValidationResponse struct {
	Valid  bool               `json:"valid"`
	Object schema.AdminObject `json:"object"`
}

type clusterPreviewResponse struct {
	Object schema.AdminObject `json:"object"`
	YAML   string             `json:"yaml"`
}

type clusterErrorResponse struct {
	Message string                   `json:"message"`
	Issues  []schema.ValidationIssue `json:"issues,omitempty"`
}

// GetClusterSchema returns only the Cluster form definition. The legacy route
// schema endpoint remains route-specific.
func GetClusterSchema(c *gin.Context) {
	registry, err := schema.DefaultRegistry()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	definition, ok := registry.Lookup(schema.KindAdminCluster)
	if !ok {
		writeClusterDraftError(c, schema.ErrSchemaNotFound)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet([]schema.ObjectSchema{definition}))
}

func CreateClusterDraft(c *gin.Context) {
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if source := strings.TrimSpace(c.Query("sourceId")); source != "" {
		id, parseErr := strconv.Atoi(source)
		if parseErr != nil || id <= 0 {
			writeClusterDraftError(c, errors.New("sourceId must be positive"))
			return
		}
		body, readErr := io.ReadAll(io.LimitReader(c.Request.Body, maxRouteBindingRequestBytes+1))
		if readErr != nil {
			writeClusterDraftError(c, readErr)
			return
		}
		if len(bytes.TrimSpace(body)) > 0 {
			writeClusterDraftError(c, errors.New("sourceId creation does not accept a request body"))
			return
		}
		created, createErr := store.CreateFromPublished(c.Request.Context(), id)
		if createErr != nil {
			writeClusterDraftError(c, createErr)
			return
		}
		c.JSON(http.StatusOK, adminconfig.WithRet(created))
		return
	}
	object, _, err := decodeClusterDraftRequest(c)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	created, err := store.CreateDraft(c.Request.Context(), object)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(created))
}

func UpdateClusterDraft(c *gin.Context) {
	id := strings.TrimSpace(c.Query("id"))
	object, bodyRevision, err := decodeClusterDraftRequest(c)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	revision, err := clusterExpectedRevision(c)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if bodyRevision > 0 {
		if revision > 0 && revision != bodyRevision {
			writeClusterDraftError(c, errors.New("expectedRevision values disagree"))
			return
		}
		revision = bodyRevision
	}
	if revision <= 0 {
		writeClusterDraftError(c, errors.New("positive expectedRevision is required"))
		return
	}
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	updated, err := store.UpdateDraft(c.Request.Context(), id, object, revision)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(updated))
}

func GetClusterDraft(c *gin.Context) {
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	draft, err := store.GetDraft(c.Request.Context(), strings.TrimSpace(c.Query("id")))
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(draft))
}

func ListClusterDrafts(c *gin.Context) {
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	drafts, err := store.ListDrafts(c.Request.Context())
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(drafts))
}

func ValidateClusterDraft(c *gin.Context) {
	object, _, err := decodeClusterDraftRequest(c)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	registry, err := schema.DefaultRegistry()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if object.Kind == "" {
		object.Kind = schema.KindAdminCluster
	}
	if object.Kind != schema.KindAdminCluster {
		writeClusterDraftError(c, errors.New("kind must be AdminCluster"))
		return
	}
	normalized, err := registry.Normalize(object)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if err := checkClusterNameAvailability(c, normalized.Metadata.Name); err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(clusterValidationResponse{Valid: true, Object: normalized}))
}

func PreviewClusterDraft(c *gin.Context) {
	object, _, err := decodeClusterDraftRequest(c)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	registry, err := schema.DefaultRegistry()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if object.Kind == "" {
		object.Kind = schema.KindAdminCluster
	}
	compiled, err := schema.CompileAdminCluster(registry, object)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if err := checkClusterNameAvailability(c, compiled.Cluster.Name); err != nil {
		writeClusterDraftError(c, err)
		return
	}
	preview, err := compiled.PreviewYAML()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(clusterPreviewResponse{Object: compiled.Source, YAML: string(preview)}))
}

func checkClusterNameAvailability(c *gin.Context, name string) error {
	store, err := newClusterDraftService()
	if err != nil {
		return err
	}
	return store.ValidateNameAvailability(c.Request.Context(), strings.TrimSpace(c.Query("id")), name)
}

func GetClusterDraftStatus(c *gin.Context) {
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	status, err := store.Status(c.Request.Context(), strings.TrimSpace(c.Query("id")))
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(status))
}

func GetClusterDraftDiff(c *gin.Context) {
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	diff, err := store.Diff(c.Request.Context(), strings.TrimSpace(c.Query("id")))
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(diff))
}

func PublishClusterDraft(c *gin.Context) {
	id := strings.TrimSpace(c.Query("id"))
	revision, err := clusterExpectedRevision(c)
	if err != nil || revision <= 0 {
		writeClusterDraftError(c, errors.New("positive expectedRevision is required"))
		return
	}
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	result, err := store.Publish(c.Request.Context(), id, revision)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(result))
}

func GetClusterHistory(c *gin.Context) {
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	entries, err := store.History(c.Request.Context(), strings.TrimSpace(c.Query("id")))
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(entries))
}

func RollbackClusterDraft(c *gin.Context) {
	historyRevision, err := strconv.ParseInt(strings.TrimSpace(c.Query("revision")), 10, 64)
	if err != nil || historyRevision <= 0 {
		writeClusterDraftError(c, errors.New("positive history revision is required"))
		return
	}
	expectedRevision, err := clusterExpectedRevision(c)
	if err != nil || expectedRevision <= 0 {
		writeClusterDraftError(c, errors.New("positive expectedRevision is required"))
		return
	}
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	result, err := store.Rollback(c.Request.Context(), strings.TrimSpace(c.Query("id")), historyRevision, expectedRevision)
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(result))
}

func DeleteClusterDraft(c *gin.Context) {
	revision, err := clusterExpectedRevision(c)
	if err != nil || revision <= 0 {
		writeClusterDraftError(c, errors.New("positive expectedRevision is required"))
		return
	}
	store, err := newClusterDraftService()
	if err != nil {
		writeClusterDraftError(c, err)
		return
	}
	if err := store.DeleteDraft(c.Request.Context(), strings.TrimSpace(c.Query("id")), revision); err != nil {
		writeClusterDraftError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("delete cluster draft success!"))
}

func decodeClusterDraftRequest(c *gin.Context) (schema.AdminObject, int64, error) {
	if content := c.PostForm("content"); strings.TrimSpace(content) != "" {
		object, err := decodeClusterObjectYAML([]byte(content))
		return object, 0, err
	}
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxRouteBindingRequestBytes+1))
	if err != nil {
		return schema.AdminObject{}, 0, fmt.Errorf("read cluster draft request: %w", err)
	}
	if len(body) > maxRouteBindingRequestBytes {
		return schema.AdminObject{}, 0, errors.New("cluster draft request exceeds size limit")
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return schema.AdminObject{}, 0, errors.New("cluster draft request body is empty")
	}
	var jsonObject map[string]json.RawMessage
	if err := json.Unmarshal(body, &jsonObject); err == nil {
		if jsonObject == nil {
			return schema.AdminObject{}, 0, errors.New("cluster draft request must be an object")
		}
		if _, wrapped := jsonObject["object"]; wrapped {
			var request clusterDraftRequest
			if err := decodeStrictClusterJSON(body, &request); err != nil {
				return schema.AdminObject{}, 0, err
			}
			if request.Object == nil {
				return schema.AdminObject{}, 0, errors.New("cluster draft request object is required")
			}
			if request.ExpectedRevision < 0 {
				return schema.AdminObject{}, 0, errors.New("expectedRevision must not be negative")
			}
			return *request.Object, request.ExpectedRevision, nil
		}
		var object schema.AdminObject
		if err := decodeStrictClusterJSON(body, &object); err != nil {
			return schema.AdminObject{}, 0, err
		}
		return object, 0, nil
	}
	object, err := decodeClusterObjectYAML(body)
	if err != nil {
		return schema.AdminObject{}, 0, fmt.Errorf("decode cluster draft request: %w", err)
	}
	return object, 0, nil
}

func decodeStrictClusterJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode cluster draft JSON: %w", err)
	}
	return nil
}

func decodeClusterObjectYAML(data []byte) (schema.AdminObject, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var object schema.AdminObject
	if err := decoder.Decode(&object); err != nil {
		return schema.AdminObject{}, fmt.Errorf("decode cluster draft YAML: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return schema.AdminObject{}, errors.New("cluster draft YAML must contain one document")
		}
		return schema.AdminObject{}, fmt.Errorf("decode cluster draft YAML: %w", err)
	}
	return object, nil
}

func clusterExpectedRevision(c *gin.Context) (int64, error) {
	value := strings.TrimSpace(c.Query("expectedRevision"))
	if value == "" {
		value = strings.Trim(strings.TrimSpace(c.GetHeader("If-Match")), `"`)
	}
	if value == "" {
		return 0, nil
	}
	revision, err := strconv.ParseInt(value, 10, 64)
	if err != nil || revision < 0 {
		return 0, fmt.Errorf("expectedRevision %q is invalid", value)
	}
	return revision, nil
}

func writeClusterDraftError(c *gin.Context, err error) {
	if errors.Is(err, logic.ErrClusterDraftConflict) || errors.Is(err, logic.ErrClusterPublishConflict) {
		c.JSON(http.StatusOK, adminconfig.RetData{Code: adminconfig.RETRY, Data: err.Error()})
		return
	}
	var issues schema.ValidationErrors
	if errors.As(err, &issues) {
		c.JSON(http.StatusOK, adminconfig.RetData{Code: adminconfig.ERR, Data: clusterErrorResponse{Message: err.Error(), Issues: []schema.ValidationIssue(issues)}})
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithError(err))
}
