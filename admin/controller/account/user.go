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
	"net/http"
	"time"
)

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	"github.com/apache/dubbo-go-pixiu/admin/logic/account"
)

// Logout user logout
func Logout(c *gin.Context) {
	// Invalid setting token
	j := auth.NewJWT()
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

// EditPassword modify account password
func EditPassword(c *gin.Context) {

	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("newPassword")

	//if flag, _ := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$", oldPassword); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal oldPassword",
	//		"data":   nil,
	//	})
	//	return
	//}
	//
	//if flag, _ := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$", newPassword); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal newPassword",
	//		"data":   nil,
	//	})
	//	return
	//}
	username := c.Request.Header.Get("username")
	flag, err := account.EditPassword(oldPassword, newPassword, username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("Successfully modify the password!"))
	// TODO Do I need to update the token?
	//generateToken(c, username)
}

// GetUserInfo get user information
func GetUserInfo(c *gin.Context) {
	username := c.Request.Header.Get("username")
	flag, userInfo, err := account.GetUserInfo(username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet(userInfo))
}

// GetUserRole get user role
func GetUserRole(c *gin.Context) {
	username := c.Request.Header.Get("username")
	flag, result, err := account.GetUserRole(username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet(result))
}

// CheckUserIsAdmin determine whether you are an administrator
func CheckUserIsAdmin(c *gin.Context) {
	username := c.Request.Header.Get("username")
	flag, err := account.CheckUserIsAdmin(username)
	if !flag {
		c.JSON(http.StatusOK, config.WithError(err))
		return
	}
	c.JSON(http.StatusOK, config.WithRet("This user is admin"))
}
