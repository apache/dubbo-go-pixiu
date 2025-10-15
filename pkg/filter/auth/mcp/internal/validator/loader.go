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
	"os"
)

import (
	"github.com/lestrrat-go/jwx/v3/jwk"
)

// JWKSLoader loads a jwk.Set for verification without performing network I/O
// during request validation.
type JWKSLoader interface {
	Load(ctx context.Context) (jwk.Set, error)
}

// StaticLoader loads a pre-parsed jwk.Set
type StaticLoader struct{ set jwk.Set }

func newStaticLoaderFromBytes(data []byte) (JWKSLoader, error) {
	keySet, err := jwk.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}
	return &StaticLoader{set: keySet}, nil
}

func newStaticLoaderFromFile(path string) (JWKSLoader, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS file %s: %w", path, err)
	}
	return newStaticLoaderFromBytes(data)
}

func (l *StaticLoader) Load(_ context.Context) (jwk.Set, error) { return l.set, nil }

// HTTPLoader loads a jwk.Set from a prepared jwk.Cache by lookup only.
type HTTPLoader struct {
	uri   string
	cache *jwk.Cache
}

func newHTTPLoader(cache *jwk.Cache, uri string) JWKSLoader {
	return &HTTPLoader{uri: uri, cache: cache}
}

func (r *HTTPLoader) Load(ctx context.Context) (jwk.Set, error) {
	if r.cache == nil || r.uri == "" {
		return nil, errors.New("remote loader not properly initialized")
	}
	return r.cache.Lookup(ctx, r.uri)
}
