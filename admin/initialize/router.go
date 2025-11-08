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
	"github.com/apache/dubbo-go-pixiu/admin/controller/account"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	"github.com/apache/dubbo-go-pixiu/admin/controller/configInfo"
	_ "github.com/apache/dubbo-go-pixiu/admin/doc"
)

// Routers init router
func Routers() *gin.Engine {
	var router = gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Guest router
	router.POST("/login", account.Login)
	router.POST("/register", account.Register)

	// auth router
	taR := router.Group("/", auth.JWTAuth())

	// The following router needs to check the token
	{
		// user router
		taR.POST("/user/logout", account.Logout)
		taR.POST("/user/password/edit", account.EditPassword)
		taR.POST("/user/getInfo", account.GetUserInfo)
		taR.POST("/user/getUserRole", account.GetUserRole)
		taR.POST("/user/checkIsAdmin", account.CheckUserIsAdmin)

		taR.GET("/config/api/base", configInfo.GetBaseInfo)
		taR.POST("/config/api/base/", configInfo.SetBaseInfo)
		taR.PUT("/config/api/base/", configInfo.SetBaseInfo)

		taR.GET("/config/api/resource/list", configInfo.GetResourceList)
		taR.GET("/config/api/resource/detail", configInfo.GetResourceDetail)
		taR.POST("/config/api/resource", configInfo.CreateResourceInfo)
		taR.PUT("/config/api/resource", configInfo.ModifyResourceInfo)
		taR.DELETE("/config/api/resource", configInfo.DeleteResourceInfo)

		taR.GET("/config/api/cluster/list", configInfo.GetClusterList)
		taR.GET("/config/api/cluster/detail", configInfo.DetailCluster)
		taR.POST("/config/api/cluster", configInfo.UpdateCluster)
		taR.PUT("/config/api/cluster", configInfo.CreateCluster)
		taR.DELETE("/config/api/cluster", configInfo.DeleteCluster)

		taR.GET("/config/api/listener/list", configInfo.GetListenerList)
		taR.GET("/config/api/listener/detail", configInfo.DetailListener)
		taR.POST("/config/api/listener", configInfo.UpdateListener)
		taR.PUT("/config/api/listener", configInfo.CreateListener)
		taR.DELETE("/config/api/listener", configInfo.DeleteListener)

		taR.GET("/config/api/resource/method/list", configInfo.GetMethodList)
		taR.GET("/config/api/resource/method/detail", configInfo.GetMethodDetail)
		taR.POST("/config/api/resource/method", configInfo.CreateMethodInfo)
		taR.PUT("/config/api/resource/method", configInfo.ModifyMethodInfo)
		taR.DELETE("/config/api/resource/method", configInfo.DeleteMethodInfo)

		// Which request method to choose, Temporarily choose put method
		taR.PUT("/config/api/resource/publish", configInfo.BatchReleaseResource)
		taR.PUT("/config/api/resource/method/publish", configInfo.BatchReleaseMethod)
		taR.PUT("/config/api/plugin_group/publish", configInfo.BatchReleasePluginGroup)
	}

	return router
}
