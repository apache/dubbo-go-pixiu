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

	"github.com/lestrrat-go/jwx/v3/jwk"
)

// JWKSLoader loads a jwk.Set for verification without performing network I/O
// during request validation.
type JWKSLoader interface {
	Load(ctx context.Context) (jwk.Set, error)
}

// LocalLoader loads a jwk.Set parsed from inline string or file path.
type LocalLoader struct {
	set jwk.Set
}

func newLocalLoader(local *LocalJWKS) (JWKSLoader, error) {
	if local == nil {
		return nil, errors.New("local jwks config is nil")
	}

	var data []byte
	var err error
	if local.InlineString != "" {
		data = []byte(local.InlineString)
	} else if local.FilePath != "" {
		data, err = os.ReadFile(local.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read JWKS file %s: %w", local.FilePath, err)
		}
	} else {
		return nil, errors.New("either inline_string or file_path must be specified for local JWKS")
	}

	keySet, err := jwk.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}
	return &LocalLoader{set: keySet}, nil
}

func (l *LocalLoader) Load(_ context.Context) (jwk.Set, error) {
	return l.set, nil
}

// RemoteLoader loads a jwk.Set from a prepared jwk.Cache by lookup only.
type RemoteLoader struct {
	uri   string
	cache *jwk.Cache
}

func newRemoteLoader(cache *jwk.Cache, uri string) JWKSLoader {
	return &RemoteLoader{uri: uri, cache: cache}
}

func (r *RemoteLoader) Load(ctx context.Context) (jwk.Set, error) {
	if r.cache == nil || r.uri == "" {
		return nil, errors.New("remote loader not properly initialized")
	}
	return r.cache.Lookup(ctx, r.uri)
}
