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
	"net/http"
	"os"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
)

// @Tags Auth
// @Summary JWT check midware
// @Description Validate the token field in the request header, parse and verify the JWT. If verification fails, the request will be terminated.
// @Produce application/json
// @Success 200 {object} string
// Note: this is a middleware, not a direct API endpoint
// JWTAuth Check token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, config.WithError(errors.New("Request does not carry token, no access")))
			c.Abort()
			return
		}
		//log.Print("get token: ", token)
		j := NewJWT()
		// Parse the information contained in the token
		claims, err := j.ParseToken(token)
		if err != nil {
			// token authorization expiration
			if err == TokenExpired {
				c.JSON(http.StatusOK, config.WithError(errors.New("The token authorization has expired, please reapply for authorization")))
				c.Abort()
				return
			}
			// Other token error conditions
			c.JSON(http.StatusOK, config.WithError(err))
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}

// JWT Signature structure
type JWT struct {
	SigningKey []byte
}

const jwtSignKeyEnv = "DUBBOGO_PIXIU_JWT_SIGN_KEY"

// Constant
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token is not valid yet")
	TokenMalformed   error  = errors.New("This is not a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token")
	SignKey          string = "dubbo-go-pixiu" // TODO: The signature information is set to be dynamically obtained
)

// Custom Claims
type CustomClaims struct {
	Username string `json:"username"`
	// The StandardClaims structure implements the Claims interface (Valid() function)
	jwt.StandardClaims
}

// New jwt instance
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// get signKey
func GetSignKey() string {
	if key := os.Getenv(jwtSignKeyEnv); key != "" {
		return key
	}
	return SignKey
}

// CreateToken Generate token (based on user basic information)
// HS256 algorithm
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// Returns the structure pointer of the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, TokenInvalid
	}
	return j.SigningKey, nil
}

// ParseToken
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// Input: token string, custom Claims structure object, custom function
	// Parse the token string into jwt's Token structure pointer
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, j.keyFunc)
	if err != nil {
		// jwt returns *ValidationError; errors.As matches both value and pointer
		// forms so the expiry/malformed branches below are actually reached.
		return nil, j.classify(err)
	}
	// Parse the claims information in the token and verify the original user data, make the following types of assertions
	//, and convert token.Claims into a specific user-defined Claims structure
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// Update token
//
// RefreshToken reissues a token whose only problem is expiration. It deliberately
// does not mutate the package-global jwt.TimeFunc (the previous implementation did,
// which leaked the frozen clock to ParseToken on the parse-failure path and raced
// with concurrent ParseToken/RefreshToken calls). Instead it parses normally and,
// when the sole validation failure is expiration, reuses the parsed claims with a
// fresh expiration time. Any other failure (malformed token, bad signature,
// non-HMAC algorithm, not-valid-yet, etc.) is rejected just like ParseToken does.
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, j.keyFunc)
	if err != nil {
		// A token that is otherwise valid but past its ExpiresAt is exactly what
		// refresh exists to handle. Only the expired bit is tolerated; any other
		// validation failure is reported with the same semantics as ParseToken.
		var ve *jwt.ValidationError
		if !errors.As(err, &ve) || ve.Errors != jwt.ValidationErrorExpired {
			return "", j.classify(err)
		}
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", TokenInvalid
	}
	// Set token expiration time
	claims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
	return j.CreateToken(*claims)
}

// classify maps a raw parse error onto the package sentinel errors, mirroring the
// mapping ParseToken performs so that both entry points report consistently.
func (j *JWT) classify(err error) error {
	var ve *jwt.ValidationError
	if !errors.As(err, &ve) {
		return TokenInvalid
	}
	switch {
	case ve.Errors&jwt.ValidationErrorMalformed != 0:
		return TokenMalformed
	case ve.Errors&jwt.ValidationErrorExpired != 0:
		return TokenExpired
	case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
		return TokenNotValidYet
	default:
		return TokenInvalid
	}
}
