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

package openapi

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/pb33f/libopenapi"
	openapiValidator "github.com/pb33f/libopenapi-validator"
	validatorConfig "github.com/pb33f/libopenapi-validator/config"
	validatorErrors "github.com/pb33f/libopenapi-validator/errors"
	validatorPaths "github.com/pb33f/libopenapi-validator/paths"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind is the kind of OpenAPI validation filter.
	Kind = constant.HTTPOpenAPIFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	FilterFactory struct {
		cfg       *Config
		validator openapiValidator.Validator
		model     *v3.Document
	}

	Filter struct {
		validator openapiValidator.Validator
		model     *v3.Document
	}

	Config struct {
		Path string `yaml:"path" json:"path,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	path, err := cleanOpenAPIPath(factory.cfg.Path)
	if err != nil {
		return err
	}
	factory.cfg.Path = path

	validator, model, err := loadValidatorFromFile(path)
	if err != nil {
		return err
	}
	factory.validator = validator
	factory.model = model
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		validator: factory.validator,
		model:     factory.model,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	if f.validator == nil || f.model == nil {
		return filter.Continue
	}

	req := ctx.Request
	pathItem, foundPath, ok := f.findRequestOperation(req)
	if !ok {
		return filter.Continue
	}

	if valid, validationErrs := f.validator.ValidateHttpRequestSyncWithPathItem(req, pathItem, foundPath); !valid {
		validationDetails := formatValidationErrors(validationErrs)
		errResp := contexthttp.BadRequest.WithError(errors.New("openapi request validation failed"))
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		logger.Debugf("openapi request validation failed: %s", validationDetails)
		return filter.Stop
	}
	return filter.Continue
}

func (f *Filter) findRequestOperation(req *http.Request) (*v3.PathItem, string, bool) {
	pathItem, validationErrs, foundPath := validatorPaths.FindPath(req, f.model, nil)
	if len(validationErrs) > 0 {
		logger.Debugf("openapi path lookup errors for %s %s: %s", req.Method, req.URL.Path, formatValidationErrors(validationErrs))
		return nil, "", false
	}
	if pathItem == nil {
		return nil, "", false
	}
	if !hasRequestOperation(req, pathItem) {
		return nil, "", false
	}
	return pathItem, foundPath, true
}

func hasRequestOperation(req *http.Request, pathItem *v3.PathItem) bool {
	switch req.Method {
	case http.MethodGet:
		return pathItem.Get != nil
	case http.MethodPost:
		return pathItem.Post != nil
	case http.MethodPut:
		return pathItem.Put != nil
	case http.MethodDelete:
		return pathItem.Delete != nil
	case http.MethodOptions:
		return pathItem.Options != nil
	case http.MethodHead:
		return pathItem.Head != nil || pathItem.Get != nil
	case http.MethodPatch:
		return pathItem.Patch != nil
	case http.MethodTrace:
		return pathItem.Trace != nil
	default:
		operations := pathItem.GetOperations()
		return operations != nil && operations.GetOrZero(strings.ToLower(req.Method)) != nil
	}
}

func loadValidatorFromFile(path string) (openapiValidator.Validator, *v3.Document, error) {
	spec, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read openapi file")
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(spec, &datamodel.DocumentConfiguration{
		BasePath: filepath.Dir(path),
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse openapi document")
	}

	model, err := doc.BuildV3Model()
	if err != nil {
		return nil, nil, errors.Wrap(err, "build openapi model")
	}

	validator := openapiValidator.NewValidatorFromV3Model(&model.Model, validatorConfig.WithoutSecurityValidation())
	if validator == nil {
		return nil, nil, errors.New("load openapi validator: validator is nil")
	}
	return validator, &model.Model, nil
}

func cleanOpenAPIPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("openapi path is required")
	}
	if filepath.IsAbs(path) {
		return "", errors.Errorf("openapi path must be relative: %s", path)
	}
	if containsParentDirectory(path) {
		return "", errors.Errorf("openapi path must not contain parent directory: %s", path)
	}

	cleanPath := filepath.Clean(path)
	if cleanPath == "." {
		return "", errors.New("openapi path is required")
	}

	basePath, err := filepath.Abs(filepath.Dir(cleanPath))
	if err != nil {
		return "", errors.Wrap(err, "resolve openapi path base")
	}
	if isSensitiveOpenAPIBasePath(basePath) {
		return "", errors.Errorf("openapi path base directory is not allowed: %s", basePath)
	}
	return cleanPath, nil
}

func containsParentDirectory(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

func isSensitiveOpenAPIBasePath(path string) bool {
	path = filepath.Clean(path)
	if path == filepath.Clean(string(filepath.Separator)) {
		return true
	}

	sensitivePaths := []string{"/etc", "/proc", "/sys", "/dev", "/run", "/var/run"}
	for _, sensitivePath := range sensitivePaths {
		if isPathWithin(path, sensitivePath) {
			return true
		}
	}
	return false
}

func isPathWithin(path string, base string) bool {
	base = filepath.Clean(base)
	if path == base {
		return true
	}

	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func formatValidationErrors(errs []*validatorErrors.ValidationError) string {
	messages := make([]string, 0, len(errs))
	for _, validationErr := range errs {
		if validationErr == nil {
			continue
		}
		switch {
		case validationErr.Message != "" && validationErr.Reason != "":
			messages = append(messages, validationErr.Message+": "+validationErr.Reason)
		case validationErr.Message != "":
			messages = append(messages, validationErr.Message)
		case validationErr.Reason != "":
			messages = append(messages, validationErr.Reason)
		default:
			messages = append(messages, validationErr.Error())
		}
	}
	if len(messages) == 0 {
		return "openapi request validation failed"
	}
	return strings.Join(messages, "; ")
}

var _ filter.HttpFilterFactory = new(FilterFactory)
