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

package handler

import (
	"time"
)

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/i18n"
)

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token is not valid yet")
	ErrTokenMalformed   = errors.New("malformed token")
	ErrTokenInvalid     = errors.New("invalid token")
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (h *Handler) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			config.UnauthorizedI18n(c, i18n.MsgTokenMissing)
			c.Abort()
			return
		}
		claims, err := h.parseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				config.UnauthorizedI18n(c, i18n.MsgTokenExpired)
			} else {
				config.Unauthorized(c, err.Error())
			}
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func (h *Handler) createToken(username string) (string, error) {
	expireHours := h.jwtConfig.GetExpireHours()
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
			Issuer:    h.jwtConfig.GetIssuer(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtConfig.GetSignKey()))
}

func (h *Handler) createExpiredToken() (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix(),
			Issuer:    h.jwtConfig.GetIssuer(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtConfig.GetSignKey()))
}

func (h *Handler) parseToken(tokenString string) (*Claims, error) {
	signKey := h.jwtConfig.GetSignKey()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			}
		}
		return nil, ErrTokenInvalid
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}
