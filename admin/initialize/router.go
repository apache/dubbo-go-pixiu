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

package initialize

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

import (
	account2 "github.com/apache/dubbo-go-pixiu/admin/controller/account"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	configInfo2 "github.com/apache/dubbo-go-pixiu/admin/controller/configInfo"
	_ "github.com/apache/dubbo-go-pixiu/docs"
)

// Routers init router
func Routers() *gin.Engine {
	var router = gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Guest router
	router.POST("/login", account2.Login)
	router.POST("/register", account2.Register)

	// auth router
	taR := router.Group("/", auth.JWTAuth())

	// The following router needs to check the token
	{
		// user router
		taR.POST("/user/logout", account2.Logout)
		taR.POST("/user/password/edit", account2.EditPassword)
		taR.POST("/user/getInfo", account2.GetUserInfo)
		taR.POST("/user/getUserRole", account2.GetUserRole)
		taR.POST("/user/checkIsAdmin", account2.CheckUserIsAdmin)

		taR.GET("/config/api/base", configInfo2.GetBaseInfo)
		taR.POST("/config/api/base/", configInfo2.SetBaseInfo)
		taR.PUT("/config/api/base/", configInfo2.SetBaseInfo)

		taR.GET("/config/api/resource/list", configInfo2.GetResourceList)
		taR.GET("/config/api/resource/detail", configInfo2.GetResourceDetail)
		taR.POST("/config/api/resource", configInfo2.CreateResourceInfo)
		taR.PUT("/config/api/resource", configInfo2.ModifyResourceInfo)
		taR.DELETE("/config/api/resource", configInfo2.DeleteResourceInfo)

		taR.GET("/config/api/cluster/list", configInfo2.GetClusterList)
		taR.GET("/config/api/cluster/detail", configInfo2.DetailCluster)
		taR.POST("/config/api/cluster", configInfo2.UpdateCluster)
		taR.PUT("/config/api/cluster", configInfo2.CreateCluster)
		taR.DELETE("/config/api/cluster", configInfo2.DeleteCluster)

		taR.GET("/config/api/listener/list", configInfo2.GetListenerList)
		taR.GET("/config/api/listener/detail", configInfo2.DetailListener)
		taR.POST("/config/api/listener", configInfo2.UpdateListener)
		taR.PUT("/config/api/listener", configInfo2.CreateListener)
		taR.DELETE("/config/api/listener", configInfo2.DeleteListener)

		taR.GET("/config/api/resource/method/list", configInfo2.GetMethodList)
		taR.GET("/config/api/resource/method/detail", configInfo2.GetMethodDetail)
		taR.POST("/config/api/resource/method", configInfo2.CreateMethodInfo)
		taR.PUT("/config/api/resource/method", configInfo2.ModifyMethodInfo)
		taR.DELETE("/config/api/resource/method", configInfo2.DeleteMethodInfo)

		// Which request method to choose, Temporarily choose put method
		taR.PUT("/config/api/resource/publish", configInfo2.BatchReleaseResource)
		taR.PUT("/config/api/resource/method/publish", configInfo2.BatchReleaseMethod)
		taR.PUT("/config/api/plugin_group/publish", configInfo2.BatchReleasePluginGroup)
	}

	return router
}
