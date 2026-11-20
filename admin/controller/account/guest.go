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

package account

import (
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	"github.com/apache/dubbo-go-pixiu/admin/logic/account"
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

// @Tags Register
// @Summary register an administrator account
// @Description Create an administrator account. The request requires a valid JWT belonging to an administrator.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param token header string true "JWT token"
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} string "{\"code\":\"\",\"data\":\"\"}"
// @Router /register [post]
func Register(c *gin.Context) {
	registerAccount(c, account.CheckUserIsAdmin, account.Register)
}

func registerAccount(c *gin.Context, checkUserIsAdmin func(string) (bool, error), createAccount func(string, string) error) {
	currentUsername := auth.CurrentUsername(c)
	isAdmin, err := checkUserIsAdmin(currentUsername)
	if !isAdmin {
		if err == nil {
			err = errors.New("only administrators can register accounts")
		}
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}

	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")
	usernameLength := utf8.RuneCountInString(username)
	passwordLength := utf8.RuneCountInString(password)
	if usernameLength == 0 || usernameLength > 255 {
		c.JSON(http.StatusOK, config.WithError(errors.New("username must be between 1 and 255 characters")))
		return
	}
	if passwordLength < 8 || passwordLength > 128 {
		c.JSON(http.StatusOK, config.WithError(errors.New("password must be between 8 and 128 characters")))
		return
	}

	if err := createAccount(username, utils.Md5(password)); err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("Register successfully, please login!"))
}

// login result
type LoginResult struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//username := c.Query("username")
	//password := c.Query("password")

	password = utils.Md5(password)
	flag, _ := account.Login(username, password)
	if flag {
		generateToken(c, username)
	} else {
		c.JSON(http.StatusOK, config.WithError(errors.New("Authentication failed, login information is wrong!")))
	}
}

func generateToken(c *gin.Context, username string) {
	j, err := auth.NewJWT()
	if err != nil {
		log.Printf("failed to initialize admin JWT: %v", err)
		c.JSON(http.StatusInternalServerError, config.WithError(errors.New("authentication is not configured")))
		return
	}
	claims := auth.CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // Signature effective time
			ExpiresAt: int64(time.Now().Unix() + 3600), // Signature expiration time
			Issuer:    "dubbo-go-pixiu",
		},
	}
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}

	data := LoginResult{
		Username: username,
		Token:    token,
	}
	c.JSON(http.StatusOK, config.WithRet(data))
}
