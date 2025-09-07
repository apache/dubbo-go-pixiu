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

package validator

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"
)

import (
	"github.com/lestrrat-go/httprc/v3"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Error code constants to avoid magic strings in responses.
const (
	ErrCodeInvalidProvider = "invalid_provider"
	ErrCodeJWKS            = "jwks_error"
	ErrCodeInvalidToken    = "invalid_token"
	ErrCodeTokenExpired    = "token_expired"
	ErrCodeTokenNotYet     = "token_not_yet_valid"
)

// TODO(validator): dynamic provider update (Add/Update/Remove) via atomic snapshot or RWMutex

const (
	defaultAcceptableSkew        = 60 * time.Second
	defaultRemoteJWKSHTTPTimeout = 5 * time.Second
)

// allowedSignatureAlgorithms defines the whitelist of acceptable JWS algorithms
// for verifying access tokens. This mitigates algorithm confusion/downgrade.
var allowedSignatureAlgorithms = map[string]struct{}{
	// Asymmetric RSA (recommended)
	"RS256": {},
	"RS384": {},
	"RS512": {},
	// RSASSA-PSS (recommended)
	"PS256": {},
	"PS384": {},
	"PS512": {},
	// ECDSA (recommended)
	"ES256": {},
	"ES384": {},
	"ES512": {},
	// Edwards (modern)
	"EdDSA": {},
	// Optionally allow HMAC for compatible deployments. If your AS never uses HMAC,
	// remove HS* to further tighten security.
	"HS256": {},
	"HS384": {},
	"HS512": {},
}

// filterKeySetByAllowedAlgorithms filters a JWK set to only include keys
func filterKeySetByAllowedAlgorithms(source jwk.Set) (jwk.Set, int) {
	if source == nil {
		return nil, 0
	}
	filtered := jwk.NewSet()
	kept := 0
	for i := 0; i < source.Len(); i++ {
		key, ok := source.Key(i)
		if !ok {
			continue
		}
		var algStr string
		if err := key.Get("alg", &algStr); err != nil || algStr == "" {
			// Try retrieving as jwa.SignatureAlgorithm, then stringify
			var sa jwa.SignatureAlgorithm
			if err2 := key.Get("alg", &sa); err2 == nil {
				algStr = sa.String()
			}
		}
		if algStr == "" {
			continue
		}
		if _, ok := allowedSignatureAlgorithms[algStr]; !ok {
			continue
		}
		if err := filtered.AddKey(key); err == nil {
			kept++
		}
	}
	return filtered, kept
}

// Validator represents a JWT validator instance
// Remote providers use jwk.Cache for JWKS auto-refresh; local providers use a static key set.
type Validator struct {
	providers map[string]*providerInfo
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// providerInfo contains the provider configuration and its JWKS loader
type providerInfo struct {
	config Provider
	loader JWKSLoader
}

// ValidationError represents a JWT validation error
type ValidationError struct {
	Code    string `json:"error"`
	Message string `json:"error_description"`
	Err     error  `json:"-"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap exposes the underlying error for errors.Is / errors.As without leaking to clients
func (e ValidationError) Unwrap() error { return e.Err }

// categorizeJWKSLoadError maps loader errors to standardized error code/message
func categorizeJWKSLoadError(err error) (code, msg string) {
	if err == nil {
		return ErrCodeJWKS, "jwks error"
	}
	return ErrCodeJWKS, err.Error()
}

// categorizeJWTError categorizes JWT validation errors into standard error codes
func categorizeJWTError(err error) (code, msg string) {
	if err == nil {
		return ErrCodeInvalidToken, "invalid token"
	}
	// jwx v3 exposes sentinel errors; Validate wraps with fmt.Errorf("validation failed: %w", err)
	if errors.Is(err, jwt.TokenExpiredError()) {
		return ErrCodeTokenExpired, jwt.TokenExpiredError().Error()
	}
	if errors.Is(err, jwt.TokenNotYetValidError()) {
		return ErrCodeTokenNotYet, jwt.TokenNotYetValidError().Error()
	}
	return ErrCodeInvalidToken, err.Error()
}

// NewValidator creates a new JWT validator instance
func NewValidator(config Config) (*Validator, error) {
	if len(config.Providers) == 0 {
		return nil, errors.New("at least one provider must be configured")
	}

	ctx, cancel := context.WithCancel(context.Background())
	v := &Validator{
		providers: make(map[string]*providerInfo),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Initialize each provider
	for _, provider := range config.Providers {
		if err := v.addProvider(provider); err != nil {
			cancel()
			return nil, fmt.Errorf("failed to add provider %s: %w", provider.Name, err)
		}
	}

	return v, nil
}

// addProvider adds a provider to the validator
func (v *Validator) addProvider(provider Provider) error {
	entry := &providerInfo{config: provider}

	loader, err := v.buildLoaderFromJWKS(provider.JWKS)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] jwt validator build loader failed: provider=%s jwks=%s err=%v", provider.Name, provider.JWKS, err)
		return fmt.Errorf("failed to init JWKS loader: %w", err)
	}
	entry.loader = loader

	v.mu.Lock()
	v.providers[provider.Name] = entry
	v.mu.Unlock()
	return nil
}

// ProviderByTokenIssuer parses token without validation to extract the issuer
// and returns the provider name configured for that issuer.
func (v *Validator) ProviderByTokenIssuer(tokenString string) (string, error) {
	// Parse token without validation to read claims
	tok, err := jwt.Parse([]byte(tokenString), jwt.WithValidate(false), jwt.WithVerify(false))
	if err != nil {
		return "", fmt.Errorf("failed to parse token for issuer extraction: %w", err)
	}

	var iss string
	if err := tok.Get("iss", &iss); err != nil || iss == "" {
		// fallback to Issuer() accessor (returns issuer string and ok bool)
		if iss2, ok := tok.Issuer(); ok {
			iss = iss2
		}
		if iss == "" {
			return "", fmt.Errorf("issuer claim not found in token")
		}
	}

	v.mu.RLock()
	defer v.mu.RUnlock()
	for name, entry := range v.providers {
		if entry.config.Issuer == iss {
			return name, nil
		}
	}
	return "", fmt.Errorf("no provider found for issuer %s", iss)
}

// buildLoaderFromJWKS parses provider.JWKS and constructs an appropriate loader.
func (v *Validator) buildLoaderFromJWKS(jwks string) (JWKSLoader, error) {
	if jwks == "" {
		return nil, errors.New("jwks must be specified")
	}
	u, err := url.Parse(jwks)
	if err != nil {
		return nil, fmt.Errorf("invalid jwks uri: %w", err)
	}
	switch u.Scheme {
	case "http", "https":
		timeout := defaultRemoteJWKSHTTPTimeout
		// Build http client with resolved timeout
		httpClient := &http.Client{Timeout: timeout}
		client := httprc.NewClient(httprc.WithHTTPClient(httpClient))
		c, err := jwk.NewCache(v.ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create jwk cache: %w", err)
		}
		if err := c.Register(v.ctx, jwks); err != nil {
			return nil, fmt.Errorf("failed to register JWKS uri %s: %w", jwks, err)
		}
		return newHTTPLoader(c, jwks), nil
	case "file":
		return newStaticLoaderFromFile(u.Path)
	default:
		return nil, fmt.Errorf("unsupported jwks scheme: %s", u.Scheme)
	}
}

// Validate validates a JWT token using the specified provider
func (v *Validator) Validate(providerName, tokenString string) (jwt.Token, error) {
	v.mu.RLock()
	provider, exists := v.providers[providerName]
	v.mu.RUnlock()

	if !exists {
		logger.Warnf("[dubbo-go-pixiu] jwt validator provider not found: name=%s", providerName)
		return nil, ValidationError{Code: ErrCodeInvalidProvider, Message: fmt.Sprintf("provider '%s' not found", providerName)}
	}

	// Resolve key set via loader (no network IO in validation path)
	var (
		keySet jwk.Set
		err    error
	)
	keySet, err = provider.loader.Load(v.ctx)
	if err != nil {
		code, msg := categorizeJWKSLoadError(err)
		logger.Errorf("[dubbo-go-pixiu] jwt validator jwks load failed: provider=%s code=%s err=%v", providerName, code, err)
		return nil, ValidationError{Code: code, Message: msg, Err: err}
	}
	if keySet == nil {
		logger.Warnf("[dubbo-go-pixiu] jwt validator jwks not available: provider=%s", providerName)
		return nil, ValidationError{Code: ErrCodeJWKS, Message: "no JWKS available for provider"}
	}

	// Enforce algorithm whitelist by filtering the key set. Tokens signed with algorithms
	// outside this list will be rejected because no matching key remains.
	filteredKeySet, kept := filterKeySetByAllowedAlgorithms(keySet)
	if kept == 0 {
		logger.Warnf("[dubbo-go-pixiu] jwt validator no acceptable jwk after alg filter: provider=%s", providerName)
		return nil, ValidationError{Code: ErrCodeJWKS, Message: "no acceptable JWKs with allowed algorithms"}
	}

	// Build parse options
	opts := make([]jwt.ParseOption, 0, 5)
	opts = append(opts,
		jwt.WithKeySet(filteredKeySet),
		jwt.WithIssuer(provider.config.Issuer),
		jwt.WithValidate(true),
		jwt.WithAcceptableSkew(defaultAcceptableSkew),
	)
	if provider.config.Audience != "" {
		opts = append(opts, jwt.WithAudience(provider.config.Audience))
	}

	// Parse and validate the token (iss/exp/nbf etc.)
	token, err := jwt.Parse([]byte(tokenString), opts...)
	if err != nil {
		code, msg := categorizeJWTError(err)
		logger.Warnf("[dubbo-go-pixiu] jwt validator token validate failed: provider=%s iss=%s code=%s err=%v", providerName, provider.config.Issuer, code, err)
		return nil, ValidationError{Code: code, Message: msg, Err: err}
	}

	return token, nil
}

// Provider returns the provider configuration by name
func (v *Validator) Provider(name string) (*Provider, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	p, ok := v.providers[name]
	if !ok {
		return nil, false
	}
	cp := p.config
	return &cp, true
}

// Providers returns the list of provider names
func (v *Validator) Providers() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Return sorted names for consistency, using maps.Keys + slices.Sorted
	names := slices.Sorted(maps.Keys(v.providers))

	return names
}

// Close shuts down background resources
func (v *Validator) Close() error {
	if v.cancel != nil {
		logger.Infof("[dubbo-go-pixiu] jwt validator shutting down")
		v.cancel()
	}
	return nil
}
