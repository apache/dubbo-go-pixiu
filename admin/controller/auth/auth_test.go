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

package auth

import (
	"sync"
	"testing"
	"time"
)

import (
	"github.com/golang-jwt/jwt/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTokenRejectsUnexpectedSigningMethod(t *testing.T) {
	claims := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = NewJWT().ParseToken(tokenString)

	assert.ErrorIs(t, err, TokenInvalid)
}

func TestGetSignKeyUsesEnvironmentOverride(t *testing.T) {
	t.Setenv(jwtSignKeyEnv, "from-env")

	assert.Equal(t, "from-env", GetSignKey())
}

// TestRefreshTokenReissuesExpiredToken verifies that an otherwise-valid token
// past its ExpiresAt is refreshed into a new, usable token.
func TestRefreshTokenReissuesExpiredToken(t *testing.T) {
	j := NewJWT()
	original := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}
	expired, err := j.CreateToken(original)
	require.NoError(t, err)

	// Sanity: the expired token is rejected by ParseToken.
	_, err = j.ParseToken(expired)
	assert.ErrorIs(t, err, TokenExpired)

	refreshed, err := j.RefreshToken(expired)
	require.NoError(t, err)
	assert.NotEqual(t, expired, refreshed)

	// The refreshed token must parse cleanly and carry the original identity.
	claims, err := j.ParseToken(refreshed)
	require.NoError(t, err)
	assert.Equal(t, "admin", claims.Username)
}

// TestRefreshTokenRejectsMalformedToken ensures a malformed refresh request is
// rejected and — critically — does not corrupt the package-global JWT clock, so
// a subsequent ParseToken of a valid token still works. This is the regression
// for the P0 reported in PR #978.
func TestRefreshTokenRejectsMalformedToken(t *testing.T) {
	j := NewJWT()

	_, err := j.RefreshToken("not-a-token")
	assert.ErrorIs(t, err, TokenMalformed)

	// After the failed refresh, a normal ParseToken must still enforce expiration
	// correctly (i.e. the global clock was not frozen at Unix 0).
	valid := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	validString, err := j.CreateToken(valid)
	require.NoError(t, err)
	if _, err := j.ParseToken(validString); err != nil {
		t.Fatalf("ParseToken of a valid token failed after a malformed refresh: %v", err)
	}

	// And an expired token must still be reported as expired, not accepted.
	expired := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}
	expiredString, err := j.CreateToken(expired)
	require.NoError(t, err)
	if _, err := j.ParseToken(expiredString); err != TokenExpired {
		t.Fatalf("ParseToken of an expired token = %v, want %v", err, TokenExpired)
	}
}

// TestRefreshTokenRejectsBadSignature verifies that a token signed with a
// different key is not refreshable.
func TestRefreshTokenRejectsBadSignature(t *testing.T) {
	j := NewJWT()
	claims := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}
	forged := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	forgedString, err := forged.SignedString([]byte("wrong-key"))
	require.NoError(t, err)

	_, err = j.RefreshToken(forgedString)
	assert.Error(t, err)
	assert.NotEqual(t, TokenExpired, err)
}

// TestConcurrentRefreshAndParse stresses RefreshToken and ParseToken running
// concurrently to confirm there is no data race on a shared global clock and
// that expiration checks stay correct. Run with -race.
func TestConcurrentRefreshAndParse(t *testing.T) {
	j := NewJWT()

	valid := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	validString, err := j.CreateToken(valid)
	require.NoError(t, err)

	expired := CustomClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}
	expiredString, err := j.CreateToken(expired)
	require.NoError(t, err)

	const goroutines = 50
	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for n := 0; n < iterations; n++ {
				if _, err := j.RefreshToken(expiredString); err != nil {
					t.Errorf("RefreshToken of an expired token failed: %v", err)
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			for n := 0; n < iterations; n++ {
				if _, err := j.RefreshToken("not-a-token"); err == nil {
					t.Error("RefreshToken of a malformed token unexpectedly succeeded")
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			for n := 0; n < iterations; n++ {
				if _, err := j.ParseToken(validString); err != nil {
					t.Errorf("ParseToken of a valid token failed concurrently: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()
}
