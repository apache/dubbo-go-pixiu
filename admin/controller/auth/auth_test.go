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
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSignKey = "d89b63e84a41f593bb17c5c6c3ee76a1fd314ca89ee47e5d7af6b48b1a2d9c36"

func testJWT(t *testing.T) *JWT {
	t.Helper()
	t.Setenv(jwtSignKeyEnv, testSignKey)
	j, err := NewJWT()
	require.NoError(t, err)
	return j
}

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

	_, err = testJWT(t).ParseToken(tokenString)

	assert.ErrorIs(t, err, TokenInvalid)
}

func TestGetSignKeyUsesEnvironmentOverride(t *testing.T) {
	t.Setenv(jwtSignKeyEnv, testSignKey)

	key, err := GetSignKey()
	require.NoError(t, err)
	assert.Equal(t, testSignKey, key)
}

func TestEmptySigningKeyUsesProcessLocalRandomKey(t *testing.T) {
	t.Setenv(jwtSignKeyEnv, "")
	type result struct {
		key string
		err error
	}
	results := make(chan result, 16)
	for i := 0; i < cap(results); i++ {
		go func() {
			key, err := GetSignKey()
			results <- result{key: key, err: err}
		}()
	}
	var key string
	for i := 0; i < cap(results); i++ {
		got := <-results
		require.NoError(t, got.err)
		if key == "" {
			key = got.key
		}
		assert.Equal(t, key, got.key)
	}
	material, err := hex.DecodeString(key)
	require.NoError(t, err)
	assert.Len(t, material, minSignKeyLength)
	assert.NotEqual(t, legacySignKey, key)

	for _, blank := range []string{"", "   "} {
		t.Setenv(jwtSignKeyEnv, blank)
		got, err := GetSignKey()
		require.NoError(t, err)
		assert.Equal(t, key, got)
	}
	t.Setenv(jwtSignKeyEnv, "temporary")
	require.NoError(t, os.Unsetenv(jwtSignKeyEnv))
	j, err := NewJWT()
	require.NoError(t, err)
	assert.Equal(t, key, string(j.SigningKey))
	claims := CustomClaims{Username: "admin"}
	claims.ExpiresAt = time.Now().Add(time.Hour).Unix()
	token, err := j.CreateToken(claims)
	require.NoError(t, err)
	parsed, err := NewJWT()
	require.NoError(t, err)
	_, err = parsed.ParseToken(token)
	require.NoError(t, err)
	t.Setenv(jwtSignKeyEnv, testSignKey)
	configuredKey, err := GetSignKey()
	require.NoError(t, err)
	assert.Equal(t, testSignKey, configuredKey)
}

func TestConfiguredUnsafeSigningKeyIsRejected(t *testing.T) {
	for _, key := range []string{legacySignKey, "short-key"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv(jwtSignKeyEnv, key)
			_, err := NewJWT()
			require.Error(t, err)
		})
	}
}

func TestLegacyDefaultTokenIsRejected(t *testing.T) {
	j := testJWT(t)
	claims := CustomClaims{Username: "admin"}
	claims.ExpiresAt = time.Now().Add(time.Hour).Unix()
	oldToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(legacySignKey))
	require.NoError(t, err)
	_, err = j.ParseToken(oldToken)
	assert.ErrorIs(t, err, TokenInvalid)
}

func TestParseTokenRequiresHS256(t *testing.T) {
	j := testJWT(t)
	claims := CustomClaims{Username: "admin"}
	claims.ExpiresAt = time.Now().Add(time.Hour).Unix()
	for _, method := range []*jwt.SigningMethodHMAC{jwt.SigningMethodHS384, jwt.SigningMethodHS512} {
		t.Run(method.Alg(), func(t *testing.T) {
			token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(testSignKey))
			require.NoError(t, err)
			_, err = j.ParseToken(token)
			assert.ErrorIs(t, err, TokenInvalid)
		})
	}
}

func TestJWTAuthDoesNotAcceptLegacyDefaultToken(t *testing.T) {
	j := testJWT(t)
	router := gin.New()
	reached := false
	router.GET("/protected", JWTAuth(), func(c *gin.Context) { reached = true; c.Status(http.StatusNoContent) })
	claims := CustomClaims{Username: "admin"}
	claims.ExpiresAt = time.Now().Add(time.Hour).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(legacySignKey))
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("token", token)
	router.ServeHTTP(httptest.NewRecorder(), req)
	assert.False(t, reached)

	validToken, err := j.CreateToken(claims)
	require.NoError(t, err)
	validRequest := httptest.NewRequest(http.MethodGet, "/protected", nil)
	validRequest.Header.Set("token", validToken)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, validRequest)
	assert.True(t, reached)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

// TestRefreshTokenReissuesExpiredToken verifies that an otherwise-valid token
// past its ExpiresAt is refreshed into a new, usable token.
func TestRefreshTokenReissuesExpiredToken(t *testing.T) {
	j := testJWT(t)
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
	j := testJWT(t)

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
	j := testJWT(t)
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
	j := testJWT(t)

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
