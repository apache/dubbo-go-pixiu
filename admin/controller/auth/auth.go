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
	"log"
	"net/http"
	"time"
)

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
)

// Check token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, config.WithError(errors.New("Request does not carry token, no access")))
			c.Abort()
			return
		}
		log.Print("get token: ", token)
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
	return SignKey
}

// CreateToken Generate token (based on user basic information)
// HS256 algorithm
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// Returns the structure pointer of the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// Input: token string, custom Claims structure object, custom function
	// Parse the token string into jwt's Token structure pointer
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	// Parse the claims information in the token and verify the original user data, make the following types of assertions
	//, and convert token.Claims into a specific user-defined Claims structure
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// Update token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// Expiration time verification
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// Set token expiration time
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
