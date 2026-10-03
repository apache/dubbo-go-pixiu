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
	"bytes"
	"io"
	"net/http"
	"net/url"
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
	validatorRadix "github.com/pb33f/libopenapi-validator/radix"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"
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

	defaultMaxRequestBodyBytes = 1 << 20
)

var errOpenAPIRequestBodyTooLarge = errors.New("openapi request body too large")

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	FilterFactory struct {
		cfg               *Config
		validator         openapiValidator.Validator
		model             *v3.Document
		validationOptions *validatorConfig.ValidationOptions
		maxRequestBody    int64
	}

	Filter struct {
		validator         openapiValidator.Validator
		model             *v3.Document
		validationOptions *validatorConfig.ValidationOptions
		maxRequestBody    int64
	}

	Config struct {
		Path string `yaml:"path" json:"path,omitempty"`
		// MaxRequestBodyBytes limits how much request body data OpenAPI validation may read.
		// Zero uses the default limit.
		MaxRequestBodyBytes int64 `yaml:"max_request_body_bytes" json:"max_request_body_bytes,omitempty"`
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

	maxRequestBody, err := factory.cfg.effectiveMaxRequestBodyBytes()
	if err != nil {
		return err
	}

	validator, model, validationOptions, err := loadValidatorFromFile(path)
	if err != nil {
		return err
	}
	factory.validator = validator
	factory.model = model
	factory.validationOptions = validationOptions
	factory.maxRequestBody = maxRequestBody
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		validator:         factory.validator,
		model:             factory.model,
		validationOptions: factory.validationOptions,
		maxRequestBody:    factory.maxRequestBody,
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

	if err := f.prepareRequestBodyForValidation(req); err != nil {
		if err == errOpenAPIRequestBodyTooLarge {
			errResp := contexthttp.PayloadTooLarge.WithError(err)
			ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
			logger.Debug(errResp.Error())
			return filter.Stop
		}
		errResp := contexthttp.BadRequest.WithError(errors.New("openapi request body read failed"))
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		logger.Debugf("openapi request body read failed: %v", err)
		return filter.Stop
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
	pathItem, validationErrs, foundPath := validatorPaths.FindPath(req, f.model, f.validationOptions)
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
		return pathItem.Head != nil
	case http.MethodPatch:
		return pathItem.Patch != nil
	case http.MethodTrace:
		return pathItem.Trace != nil
	default:
		operations := pathItem.GetOperations()
		return operations != nil && operations.GetOrZero(strings.ToLower(req.Method)) != nil
	}
}

func (f *Filter) prepareRequestBodyForValidation(req *http.Request) error {
	if req.Body == nil || req.Body == http.NoBody {
		return nil
	}
	if f.maxRequestBody <= 0 {
		return nil
	}
	if req.ContentLength > f.maxRequestBody {
		return errOpenAPIRequestBodyTooLarge
	}

	body, err := readRequestBodyWithinLimit(req.Body, f.maxRequestBody)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	return nil
}

func readRequestBodyWithinLimit(body io.ReadCloser, limit int64) ([]byte, error) {
	// The original body is replaced by the caller after this bounded read.
	defer body.Close()

	data, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, errOpenAPIRequestBodyTooLarge
	}
	return data, nil
}

func loadValidatorFromFile(path string) (openapiValidator.Validator, *v3.Document, *validatorConfig.ValidationOptions, error) {
	if err := validateExternalRefs(path); err != nil {
		return nil, nil, nil, err
	}

	spec, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "read openapi file")
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(spec, &datamodel.DocumentConfiguration{
		BasePath: filepath.Dir(path),
	})
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "parse openapi document")
	}

	model, err := doc.BuildV3Model()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "build openapi model")
	}

	validationOptions := validatorConfig.NewValidationOptions(validatorConfig.WithoutSecurityValidation())
	if validationOptions.PathTree == nil && !validationOptions.IsPathTreeDisabled() {
		validationOptions.PathTree = validatorRadix.BuildPathTree(&model.Model)
	}

	validator := openapiValidator.NewValidatorFromV3Model(&model.Model, validatorConfig.WithExistingOpts(validationOptions))
	if validator == nil {
		return nil, nil, nil, errors.New("load openapi validator: validator is nil")
	}
	return validator, &model.Model, validationOptions, nil
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

	// Resolve symlinks so that a relative path like "specs/current/passwd"
	// cannot bypass the sensitive-directory check when "specs/current" is
	// a symlink pointing to /etc.
	resolved, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		return "", errors.Wrap(err, "resolve openapi path symlinks")
	}

	basePath, err := filepath.Abs(filepath.Dir(resolved))
	if err != nil {
		return "", errors.Wrap(err, "resolve openapi path base")
	}
	if isSensitiveOpenAPIBasePath(basePath) {
		return "", errors.Errorf("openapi path base directory is not allowed: %s", basePath)
	}
	return cleanPath, nil
}

func validateExternalRefs(rootPath string) error {
	rootPath = filepath.Clean(rootPath)
	rootDir := filepath.Dir(rootPath)
	visited := map[string]struct{}{}
	return validateExternalRefsInFile(rootPath, rootDir, visited)
}

func validateExternalRefsInFile(path string, rootDir string, visited map[string]struct{}) error {
	path = filepath.Clean(path)
	if _, ok := visited[path]; ok {
		return nil
	}
	visited[path] = struct{}{}

	spec, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read openapi file for ref validation")
	}

	var raw any
	if err := yaml.Unmarshal(spec, &raw); err != nil {
		return errors.Wrap(err, "parse openapi file for ref validation")
	}

	refs := collectExternalRefs(raw, nil)
	for _, ref := range refs {
		if err := validateExternalRefPath(path, rootDir, ref); err != nil {
			return err
		}

		refPath, _, ok := splitRefTarget(ref)
		if !ok {
			continue
		}
		if filepath.Ext(refPath) == "" && strings.HasSuffix(refPath, "/") {
			continue
		}
		resolvedRefPath, err := resolveSafeReferencePath(filepath.Dir(path), refPath, rootDir)
		if err != nil {
			return err
		}
		if err := validateExternalRefsInFile(resolvedRefPath, rootDir, visited); err != nil {
			return err
		}
	}
	return nil
}

func collectExternalRefs(value any, refs []string) []string {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			if key == "$ref" {
				if ref, ok := child.(string); ok && ref != "" && !strings.HasPrefix(ref, "#") {
					refs = append(refs, ref)
				}
				continue
			}
			refs = collectExternalRefs(child, refs)
		}
	case []any:
		for _, child := range v {
			refs = collectExternalRefs(child, refs)
		}
	}
	return refs
}

func splitRefTarget(ref string) (string, string, bool) {
	if ref == "" {
		return "", "", false
	}
	if strings.HasPrefix(ref, "#") {
		return "", ref, false
	}
	parsed, err := url.Parse(ref)
	if err != nil {
		return "", "", false
	}
	if parsed.Scheme != "" && parsed.Scheme != "file" {
		return "", "", false
	}
	target := parsed.Path
	if target == "" {
		target = parsed.Opaque
	}
	fragment := parsed.Fragment
	if idx := strings.Index(target, "#"); idx >= 0 {
		fragment = target[idx+1:]
		target = target[:idx]
	}
	return target, fragment, true
}

func validateExternalRefPath(currentFile string, rootDir string, ref string) error {
	if ref == "" || strings.HasPrefix(ref, "#") {
		return nil
	}
	parsed, err := url.Parse(ref)
	if err != nil {
		return errors.Wrap(err, "parse openapi external ref")
	}
	if parsed.Scheme != "" || parsed.Host != "" {
		return errors.Errorf("openapi external ref must be relative: %s", ref)
	}

	refPath, _, ok := splitRefTarget(ref)
	if !ok {
		return nil
	}
	if filepath.IsAbs(refPath) {
		return errors.Errorf("openapi external ref must be relative: %s", ref)
	}
	if containsParentDirectory(refPath) {
		return errors.Errorf("openapi external ref must not contain parent directory: %s", ref)
	}
	_, err = resolveSafeReferencePath(filepath.Dir(currentFile), refPath, rootDir)
	return err
}

func resolveSafeReferencePath(baseDir string, refPath string, rootDir string) (string, error) {
	candidate := filepath.Clean(filepath.Join(baseDir, refPath))
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", errors.Wrap(err, "resolve openapi external ref symlinks")
	}
	if !isWithinBaseDir(resolved, rootDir) {
		return "", errors.Errorf("openapi external ref escapes allowed directory: %s", resolved)
	}
	return resolved, nil
}

func isWithinBaseDir(path string, baseDir string) bool {
	path = filepath.Clean(path)
	baseDir = filepath.Clean(baseDir)
	if path == baseDir {
		return true
	}
	rel, err := filepath.Rel(baseDir, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func (c *Config) effectiveMaxRequestBodyBytes() (int64, error) {
	if c.MaxRequestBodyBytes < 0 {
		return 0, errors.New("max_request_body_bytes must not be negative")
	}
	if c.MaxRequestBodyBytes == 0 {
		return defaultMaxRequestBodyBytes, nil
	}
	return c.MaxRequestBodyBytes, nil
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
	// Resolve symlinks on the candidate path so that macOS /private/etc
	// matches against the sensitive entry /etc.
	resolved := resolveExistingAncestor(path)

	if resolved == filepath.Clean(string(filepath.Separator)) {
		return true
	}

	sensitivePaths := []string{"/etc", "/proc", "/sys", "/dev", "/run", "/var/run"}
	for _, sensitivePath := range sensitivePaths {
		// Resolve symlinks on the sensitive path too (macOS: /etc -> /private/etc).
		resolvedSensitive := resolveExistingAncestor(sensitivePath)
		if isPathWithin(resolved, resolvedSensitive) {
			return true
		}
	}
	return false
}

// resolveExistingAncestor resolves symlinks on the longest existing ancestor
// of path, then appends any remaining non-existent components. This avoids
// EvalSymlinks failing on paths whose leaf does not yet exist on disk.
func resolveExistingAncestor(path string) string {
	path = filepath.Clean(path)

	// Walk up until we find an existing component.
	candidate := path
	var tail []string
	for {
		_, err := os.Lstat(candidate)
		if err == nil {
			break
		}
		parent := filepath.Dir(candidate)
		if parent == candidate {
			// Reached root without finding an existing component.
			return path
		}
		tail = append([]string{filepath.Base(candidate)}, tail...)
		candidate = parent
	}

	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return path
	}
	if len(tail) == 0 {
		return resolved
	}
	return filepath.Join(append([]string{resolved}, tail...)...)
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
