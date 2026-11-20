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
	"errors"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	"github.com/apache/dubbo-go-pixiu/admin/logic/account"
)

// @Tags Account
// @Summary user logout
// @Description log out the current user (invalidate the token)
// @Produce application/json
// @Success 200 {object} string
// @Router /user/logout [post]
// Logout user logout
func Logout(c *gin.Context) {
	// Invalid setting token
	j, err := auth.NewJWT()
	if err != nil {
		log.Printf("failed to initialize admin JWT: %v", err)
		c.JSON(http.StatusInternalServerError, config.WithError(errors.New("authentication is not configured")))
		return
	}
	claims := auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()), // Signature effective time
			ExpiresAt: int64(time.Now().Unix()), // Signature expiration time
			Issuer:    "dubbo-go-pixiu",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	// TODO Optimize the returned json
	c.JSON(http.StatusOK, config.WithRet(token))
}

// @Tags Account
// @Summary modify account password
// @Description modify the password for the user identified by the validated JWT.
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param oldPassword formData string true "old password"
// @Param newPassword formData string true "new password"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /user/password/edit [post]
// EditPassword modify account password
func EditPassword(c *gin.Context) {

	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("newPassword")
	newPasswordLength := utf8.RuneCountInString(newPassword)
	if oldPassword == "" || newPassword == "" {
		c.JSON(http.StatusOK, config.WithError(errors.New("old and new passwords are required")))
		return
	}
	if newPasswordLength < 8 || newPasswordLength > 128 {
		c.JSON(http.StatusOK, config.WithError(errors.New("new password must be between 8 and 128 characters")))
		return
	}
	if oldPassword == newPassword {
		c.JSON(http.StatusOK, config.WithError(errors.New("new password must differ from the old password")))
		return
	}
	username := auth.CurrentUsername(c)
	flag, err := account.EditPassword(oldPassword, newPassword, username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("Successfully modify the password!"))
	// TODO Do I need to update the token?
	//generateToken(c, username)
}

// @Tags Account
// @Summary get user information
// @Description return information for the user identified by the validated JWT.
// @Produce application/json
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /user/getInfo [post]
// GetUserInfo get user information
func GetUserInfo(c *gin.Context) {
	username := auth.CurrentUsername(c)
	flag, userInfo, err := account.GetUserInfo(username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet(userInfo))
}

// @Tags Account
// @Summary get user role
// @Description return role information for the user identified by the validated JWT.
// @Produce application/json
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /user/getUserRole [post]
// GetUserRole get user role
func GetUserRole(c *gin.Context) {
	username := auth.CurrentUsername(c)
	flag, result, err := account.GetUserRole(username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet(result))
}

// @Tags Account
// @Summary determine whether you are an administrator
// @Description determine whether the user identified by the validated JWT is an administrator.
// @Produce application/json
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /user/checkIsAdmin [post]
// CheckUserIsAdmin determine whether you are an administrator
func CheckUserIsAdmin(c *gin.Context) {
	username := auth.CurrentUsername(c)
	flag, err := account.CheckUserIsAdmin(username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("This user is admin"))
}
