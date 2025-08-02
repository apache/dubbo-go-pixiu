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
	"fmt"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Validator is responsible for validating JWTs from a set of trusted providers.
type Validator struct {
	providers map[string]*validatedProvider
	cache     *jwk.Cache
}

// validatedProvider holds a pre-configured jwk.Set for a single provider.
type validatedProvider struct {
	isLocal   bool
	keySet    jwk.Set // for local
	remoteURI string  // for remote
	issuer    string
	audiences []string
}

// NewValidator creates a new Validator instance from a list of provider configs.
// It will initialize the auto-refreshing JWKS sets.
func NewValidator(ctx context.Context, providers []Provider) (*Validator, error) {
	cache := jwk.NewCache(ctx)
	validatedProviders := make(map[string]*validatedProvider)

	for _, p := range providers {
		if p.JwksSource.Remote != nil {
			err := cache.Register(p.JwksSource.Remote.URI, jwk.WithHTTPClient(http.DefaultClient))
			if err != nil {
				return nil, fmt.Errorf("failed to register remote JWKS URL for provider %s: %w", p.Name, err)
			}
			validatedProviders[p.Name] = &validatedProvider{
				isLocal:   false,
				remoteURI: p.JwksSource.Remote.URI,
				issuer:    p.Issuer,
				audiences: p.Audiences,
			}
		} else if p.JwksSource.Local != nil {
			keySet, err := jwk.Parse([]byte(p.JwksSource.Local.Inline))
			if err != nil {
				return nil, fmt.Errorf("failed to create JWK set from local inline string for provider %s: %w", p.Name, err)
			}
			validatedProviders[p.Name] = &validatedProvider{
				isLocal:   true,
				keySet:    keySet,
				issuer:    p.Issuer,
				audiences: p.Audiences,
			}
		} else {
			return nil, fmt.Errorf("provider %s must have either a remote or local JWKS source", p.Name)
		}
	}

	if len(validatedProviders) == 0 {
		return nil, fmt.Errorf("no valid providers found")
	}

	return &Validator{providers: validatedProviders, cache: cache}, nil
}

// Validate parses and validates a token string against a specific provider.
func (v *Validator) Validate(ctx context.Context, providerName, tokenString string) (jwt.Token, error) {
	provider, ok := v.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	var keySet jwk.Set
	var err error
	if provider.isLocal {
		keySet = provider.keySet
	} else {
		keySet, err = v.cache.Get(ctx, provider.remoteURI)
		if err != nil {
			return nil, fmt.Errorf("failed to get JWK set from cache for provider %s: %w", providerName, err)
		}
	}

	// Parse with all validations except audience, which we check manually for "any of" logic.
	token, err := jwt.Parse([]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithIssuer(provider.issuer),
		jwt.WithValidate(true),
	)

	if err != nil {
		return nil, fmt.Errorf("token validation failed for provider %s: %w", providerName, err)
	}

	// Manual "any of" audience validation
	if len(provider.audiences) > 0 {
		isAudienceValid := false
		// Create a quick lookup map of the audiences we accept
		acceptedAudiences := make(map[string]struct{}, len(provider.audiences))
		for _, aud := range provider.audiences {
			acceptedAudiences[aud] = struct{}{}
		}

		// Iterate over the audiences in the token
		for _, tokenAud := range token.Audience() {
			// If we find one that we accept, validation is successful
			if _, ok := acceptedAudiences[tokenAud]; ok {
				isAudienceValid = true
				break
			}
		}

		if !isAudienceValid {
			return nil, fmt.Errorf("token audience claim does not contain any of the required audiences for provider %s", providerName)
		}
	}

	return token, nil
}
