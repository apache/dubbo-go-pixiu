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
	"github.com/apache/dubbo-go-pixiu/admin/controller/opa"
	xdsController "github.com/apache/dubbo-go-pixiu/admin/controller/xds"
	_ "github.com/apache/dubbo-go-pixiu/admin/doc"
)

// Routers init router
func Routers() *gin.Engine {
	return routers(auth.AdminAuth())
}

// RoutersWithAdminAuth builds the admin router with a supplied administrator
// middleware. Routers uses the database-backed role check by default.
func RoutersWithAdminAuth(adminAuth gin.HandlerFunc) *gin.Engine {
	return routers(adminAuth)
}

func routers(adminAuth gin.HandlerFunc) *gin.Engine {
	var router = gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Guest router
	router.POST("/login", account.Login)

	// auth router
	taR := router.Group("/", auth.JWTAuth())

	// The following router needs to check the token
	{
		// user router
		taR.POST("/register", account.Register)
		taR.POST("/user/logout", account.Logout)
		taR.POST("/user/password/edit", account.EditPassword)
		taR.POST("/user/getInfo", account.GetUserInfo)
		taR.POST("/user/getUserRole", account.GetUserRole)
		taR.POST("/user/checkIsAdmin", account.CheckUserIsAdmin)

		adminR := taR.Group("/", adminAuth)

		adminR.GET("/config/api/base", configInfo.GetBaseInfo)
		adminR.GET("/config/api/xds/status", xdsController.GetStatus)
		adminR.POST("/config/api/base/", configInfo.SetBaseInfo)
		adminR.PUT("/config/api/base/", configInfo.SetBaseInfo)

		adminR.GET("/config/api/cluster/list", configInfo.GetClusterList)
		adminR.GET("/config/api/cluster/detail", configInfo.DetailCluster)
		adminR.POST("/config/api/cluster", configInfo.UpdateCluster)
		adminR.PUT("/config/api/cluster", configInfo.CreateCluster)
		adminR.DELETE("/config/api/cluster", configInfo.DeleteCluster)

		adminR.GET("/config/api/listener/list", configInfo.GetListenerList)
		adminR.GET("/config/api/listener/detail", configInfo.DetailListener)
		adminR.POST("/config/api/listener", configInfo.UpdateListener)
		adminR.PUT("/config/api/listener", configInfo.CreateListener)
		adminR.DELETE("/config/api/listener", configInfo.DeleteListener)

		// Legacy Resource/Method routes stay available for older Admin clients.
		adminR.GET("/config/api/resource/list", configInfo.GetResourceList)
		adminR.GET("/config/api/resource/detail", configInfo.GetResourceDetail)
		adminR.POST("/config/api/resource", configInfo.CreateResourceInfo)
		adminR.PUT("/config/api/resource", configInfo.ModifyResourceInfo)
		adminR.DELETE("/config/api/resource", configInfo.DeleteResourceInfo)
		adminR.GET("/config/api/resource/method/list", configInfo.GetMethodList)
		adminR.GET("/config/api/resource/method/detail", configInfo.GetMethodDetail)
		adminR.POST("/config/api/resource/method", configInfo.CreateMethodInfo)
		adminR.PUT("/config/api/resource/method", configInfo.ModifyMethodInfo)
		adminR.DELETE("/config/api/resource/method", configInfo.DeleteMethodInfo)
		adminR.PUT("/config/api/resource/publish", configInfo.BatchReleaseResource)
		adminR.PUT("/config/api/resource/method/publish", configInfo.BatchReleaseMethod)

		// AdminRouteBinding API. Each route has
		// an independent draft/publish boundary; publishing one route still
		// updates its generated legacy keys in one etcd transaction.
		adminR.GET("/config/api/route/schema", configInfo.GetRouteBindingSchema)
		adminR.GET("/config/api/route/list", configInfo.GetRouteBindingList)
		adminR.GET("/config/api/route/detail", configInfo.GetRouteBindingDetail)
		adminR.POST("/config/api/route", configInfo.CreateRouteBinding)
		adminR.PUT("/config/api/route", configInfo.ModifyRouteBinding)
		adminR.DELETE("/config/api/route", configInfo.DeleteRouteBinding)
		adminR.POST("/config/api/route/validate", configInfo.ValidateRouteBinding)
		adminR.POST("/config/api/route/preview", configInfo.PreviewRouteBinding)
		adminR.PUT("/config/api/route/publish", configInfo.PublishRouteBinding)
		adminR.GET("/config/api/route/status", configInfo.GetRouteBindingStatus)
		adminR.GET("/config/api/route/diff", configInfo.GetRouteBindingDiff)
		adminR.GET("/config/api/route/history", configInfo.GetRouteBindingHistory)
		adminR.GET("/config/api/route/history/list", configInfo.GetRouteBindingHistoryList)
		adminR.PUT("/config/api/route/rollback", configInfo.RollbackRouteBinding)

		// PluginGroup publication is independent of API route publication.
		adminR.PUT("/config/api/plugin_group/publish", configInfo.BatchReleasePluginGroup)

		adminR.GET("/config/api/opa/policy", opa.GetOPAPolicy)
		adminR.PUT("/config/api/opa/policy", opa.PutOPAPolicy)
		adminR.DELETE("/config/api/opa/policy", opa.DeleteOPAPolicy)

	}

	return router
}
