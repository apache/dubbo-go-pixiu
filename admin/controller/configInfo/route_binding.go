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
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

const maxRouteBindingRequestBytes = 4 << 20

type routeBindingRequest struct {
	Object           *schema.AdminObject `json:"object"`
	ExpectedRevision int64               `json:"expectedRevision"`
}

type routeBindingValidationResponse struct {
	Valid  bool               `json:"valid"`
	Object schema.AdminObject `json:"object"`
}

type routeBindingPreviewResponse struct {
	Object schema.AdminObject `json:"object"`
	YAML   string             `json:"yaml"`
}

type routeBindingErrorResponse struct {
	Message string                   `json:"message"`
	Issues  []schema.ValidationIssue `json:"issues,omitempty"`
}

// GetRouteBindingSchema returns the AdminRouteBinding form and validation
// schema. It has no etcd dependency so the UI can render before a config
// namespace contains any routes.
//
// @Tags Config
// @Summary get API route binding schema
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/route/schema [get]
func GetRouteBindingSchema(c *gin.Context) {
	registry, err := schema.NewBuiltinRegistry()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(registry.List()))
}

// GetRouteBindingList returns high-level route bindings from the draft or
// published namespace. The default scope is draft; pass scope=published or
// unpublished=0 to read the published snapshot.
//
// @Tags Config
// @Summary get API route binding list
// @Produce application/json
// @Param scope query string false "draft or published"
// @Param unpublished query string false "1 for draft, 0 for published"
// @Success 200 {object} string
// @Router /config/api/route/list [get]
func GetRouteBindingList(c *gin.Context) {
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	bindings, err := store.List(routeBindingUnpublished(c))
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(bindings))
}

// GetRouteBindingDetail returns one high-level route binding by metadata.name.
//
// @Tags Config
// @Summary get API route binding detail
// @Produce application/json
// @Param name query string true "Route binding name"
// @Param scope query string false "draft or published"
// @Param unpublished query string false "1 for draft, 0 for published"
// @Success 200 {object} string
// @Router /config/api/route/detail [get]
func GetRouteBindingDetail(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = c.Param("name")
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	binding, err := store.Get(name, routeBindingUnpublished(c))
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(binding))
}

// CreateRouteBinding validates and saves one route binding in the draft
// namespace. The JSON body may be either an AdminObject or
// {"object": AdminObject, "expectedRevision": number}; a form content field
// containing the same object as YAML is also accepted for Admin compatibility.
//
// @Tags Config
// @Summary create API route binding draft
// @Accept application/json
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/route [post]
func CreateRouteBinding(c *gin.Context) {
	object, expectedRevision, err := decodeRouteBindingRequest(c)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	binding, err := store.SaveDraft(object, true, expectedRevision)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(binding))
}

// ModifyRouteBinding validates and replaces one route binding draft. If the
// draft does not exist but a published binding with the same name exists, its
// stable runtime identity is reused.
//
// @Tags Config
// @Summary modify API route binding draft
// @Accept application/json
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/route [put]
func ModifyRouteBinding(c *gin.Context) {
	object, expectedRevision, err := decodeRouteBindingRequest(c)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	binding, err := store.SaveDraft(object, false, expectedRevision)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(binding))
}

// DeleteRouteBinding removes a draft. If the binding is currently published,
// the next atomic publish interprets its absence as a route deletion.
//
// @Tags Config
// @Summary delete API route binding draft
// @Produce application/json
// @Param name query string true "Route binding name"
// @Param expectedRevision query int false "Draft key mod revision"
// @Success 200 {object} string
// @Router /config/api/route [delete]
func DeleteRouteBinding(c *gin.Context) {
	name := c.Query("name")
	expectedRevision, err := routeBindingExpectedRevision(c)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	if err := store.DeleteDraft(name, expectedRevision); err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet("Success"))
}

// ValidateRouteBinding validates an object and returns schema defaults without
// writing anything to etcd.
//
// @Tags Config
// @Summary validate API route binding
// @Accept application/json
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/route/validate [post]
func ValidateRouteBinding(c *gin.Context) {
	object, _, err := decodeRouteBindingRequest(c)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	registry, err := schema.NewBuiltinRegistry()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	normalized, err := registry.Normalize(object)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(routeBindingValidationResponse{Valid: true, Object: normalized}))
}

// PreviewRouteBinding validates and compiles an object to the legacy
// api_config.yaml shape without writing anything to etcd.
//
// @Tags Config
// @Summary preview API route binding
// @Accept application/json
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/route/preview [post]
func PreviewRouteBinding(c *gin.Context) {
	object, _, err := decodeRouteBindingRequest(c)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	registry, err := schema.NewBuiltinRegistry()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	normalized, err := registry.Normalize(object)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	compiled, err := schema.CompileAdminRouteBinding(registry, normalized)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	preview, err := compiled.PreviewYAML()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(routeBindingPreviewResponse{Object: normalized, YAML: string(preview)}))
}

// PublishRouteBinding publishes one route draft in one etcd transaction. Pass
// name and expectedRevision to reject a stale UI snapshot.
//
// @Tags Config
// @Summary atomically publish one API route binding
// @Produce application/json
// @Param name query string true "Route binding name"
// @Param expectedRevision query int false "Draft route key mod revision"
// @Success 200 {object} string
// @Router /config/api/route/publish [put]
func PublishRouteBinding(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = c.Param("name")
	}
	if strings.TrimSpace(name) == "" {
		writeRouteBindingError(c, errors.New("route binding name is required for publish"))
		return
	}
	expectedRevision, err := routeBindingExpectedRevision(c)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	result, err := store.Publish(name, expectedRevision)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(result))
}

// GetRouteBindingStatus returns the independent draft and published
// revisions for one route binding.
//
// @Tags Config
// @Summary get API route binding status
// @Produce application/json
// @Param name query string true "Route binding name"
// @Success 200 {object} string
// @Router /config/api/route/status [get]
func GetRouteBindingStatus(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = c.Param("name")
	}
	if strings.TrimSpace(name) == "" {
		writeRouteBindingError(c, errors.New("route binding name is required for status"))
		return
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	status, err := store.Status(name)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(status))
}

// GetRouteBindingDiff returns the field-level diff between one route's draft
// and published objects.
//
// @Tags Config
// @Summary get API route binding diff
// @Produce application/json
// @Param name query string true "Route binding name"
// @Success 200 {object} string
// @Router /config/api/route/diff [get]
func GetRouteBindingDiff(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = c.Param("name")
	}
	if strings.TrimSpace(name) == "" {
		writeRouteBindingError(c, errors.New("route binding name is required for diff"))
		return
	}
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	diff, err := store.Diff(name)
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(diff))
}

// GetRouteBindingPublishStatus returns the draft/published marker revisions
// used by the legacy full-snapshot publish contract.
//
// @Tags Config
// @Summary get API route binding publish status
// @Produce application/json
// @Success 200 {object} string
// @Router /config/api/route/publish/status [get]
func GetRouteBindingPublishStatus(c *gin.Context) {
	store, err := logic.NewAdminRouteBindingStore()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	status, err := store.PublishStatus()
	if err != nil {
		writeRouteBindingError(c, err)
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithRet(status))
}

func decodeRouteBindingRequest(c *gin.Context) (schema.AdminObject, int64, error) {
	if content := c.PostForm("content"); strings.TrimSpace(content) != "" {
		object, err := schema.DecodeAdminObjectYAML([]byte(content))
		if err != nil {
			return schema.AdminObject{}, 0, err
		}
		expectedRevision, err := routeBindingExpectedRevision(c)
		return object, expectedRevision, err
	}

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxRouteBindingRequestBytes+1))
	if err != nil {
		return schema.AdminObject{}, 0, fmt.Errorf("read route binding request: %w", err)
	}
	if len(body) > maxRouteBindingRequestBytes {
		return schema.AdminObject{}, 0, fmt.Errorf("route binding request exceeds %d bytes", maxRouteBindingRequestBytes)
	}
	body = []byte(strings.TrimSpace(string(body)))
	if len(body) == 0 {
		return schema.AdminObject{}, 0, errors.New("route binding request body is empty")
	}

	var request routeBindingRequest
	if err := json.Unmarshal(body, &request); err == nil && request.Object != nil {
		expectedRevision := request.ExpectedRevision
		if expectedRevision == 0 {
			expectedRevision, err = routeBindingExpectedRevision(c)
		}
		return *request.Object, expectedRevision, err
	}

	var object schema.AdminObject
	if err := json.Unmarshal(body, &object); err == nil && (object.Kind != "" || object.Metadata.Name != "" || object.Spec != nil) {
		expectedRevision, revisionErr := routeBindingExpectedRevision(c)
		return object, expectedRevision, revisionErr
	}

	object, err = schema.DecodeAdminObjectYAML(body)
	if err != nil {
		return schema.AdminObject{}, 0, fmt.Errorf("decode route binding request: %w", err)
	}
	expectedRevision, revisionErr := routeBindingExpectedRevision(c)
	return object, expectedRevision, revisionErr
}

func routeBindingExpectedRevision(c *gin.Context) (int64, error) {
	value := strings.TrimSpace(c.Query("expectedRevision"))
	if value == "" {
		value = strings.Trim(strings.TrimSpace(c.GetHeader("If-Match")), `"`)
	}
	if value == "" {
		return 0, nil
	}
	revision, err := strconv.ParseInt(value, 10, 64)
	if err != nil || revision < 0 {
		return 0, fmt.Errorf("expected revision %q is invalid", value)
	}
	return revision, nil
}

func routeBindingUnpublished(c *gin.Context) bool {
	if strings.EqualFold(strings.TrimSpace(c.Query("scope")), "published") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("scope")), "draft") {
		return true
	}
	return c.Query("unpublished") != "0"
}

func writeRouteBindingError(c *gin.Context, err error) {
	if err == nil {
		err = errors.New("unknown route binding error")
	}
	var validationErrors schema.ValidationErrors
	if errors.As(err, &validationErrors) {
		issues := append([]schema.ValidationIssue(nil), validationErrors...)
		c.JSON(http.StatusOK, adminconfig.RetData{
			Code: adminconfig.ERR,
			Data: routeBindingErrorResponse{Message: err.Error(), Issues: issues},
		})
		return
	}
	c.JSON(http.StatusOK, adminconfig.WithError(err))
}
