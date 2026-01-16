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
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/model"
)

func (h *Handler) SetupRouter(r *gin.Engine) {
	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}

	user := api.Group("/user")
	user.Use(h.JWTAuth())
	{
		user.POST("/logout", h.Logout)
		user.POST("/password", h.EditPassword)
		user.GET("/info", h.GetUserInfo)
		user.GET("/role", h.GetUserRole)
		user.GET("/is-admin", h.CheckUserIsAdmin)
	}

	clusters := api.Group("/clusters")
	clusters.Use(h.JWTAuth())
	{
		clusters.GET("", h.RequirePermission(model.ResourceClusters, model.ActionRead), h.ListClusters)
		clusters.GET("/:name", h.RequirePermission(model.ResourceClusters, model.ActionRead), h.GetCluster)
		clusters.POST("", h.RequirePermission(model.ResourceClusters, model.ActionCreate), h.CreateCluster)
		clusters.PUT("/:name", h.RequirePermission(model.ResourceClusters, model.ActionUpdate), h.UpdateCluster)
		clusters.DELETE("/:name", h.RequirePermission(model.ResourceClusters, model.ActionDelete), h.DeleteCluster)
	}

	listeners := api.Group("/listeners")
	listeners.Use(h.JWTAuth())
	{
		listeners.GET("", h.RequirePermission(model.ResourceListeners, model.ActionRead), h.ListListeners)
		listeners.GET("/:name", h.RequirePermission(model.ResourceListeners, model.ActionRead), h.GetListener)
		listeners.POST("", h.RequirePermission(model.ResourceListeners, model.ActionCreate), h.CreateListener)
		listeners.PUT("/:name", h.RequirePermission(model.ResourceListeners, model.ActionUpdate), h.UpdateListener)
		listeners.DELETE("/:name", h.RequirePermission(model.ResourceListeners, model.ActionDelete), h.DeleteListener)
	}

	resources := api.Group("/resources")
	resources.Use(h.JWTAuth())
	{
		resources.GET("", h.RequirePermission(model.ResourceResources, model.ActionRead), h.ListResources)
		resources.GET("/:id", h.RequirePermission(model.ResourceResources, model.ActionRead), h.GetResource)
		resources.POST("", h.RequirePermission(model.ResourceResources, model.ActionCreate), h.CreateResource)
		resources.PUT("/:id", h.RequirePermission(model.ResourceResources, model.ActionUpdate), h.UpdateResource)
		resources.DELETE("/:id", h.RequirePermission(model.ResourceResources, model.ActionDelete), h.DeleteResource)
	}

	methods := api.Group("/methods")
	methods.Use(h.JWTAuth())
	{
		methods.GET("", h.RequirePermission(model.ResourceMethods, model.ActionRead), h.ListMethods)
		methods.GET("/:id", h.RequirePermission(model.ResourceMethods, model.ActionRead), h.GetMethod)
		methods.POST("", h.RequirePermission(model.ResourceMethods, model.ActionCreate), h.CreateMethod)
		methods.PUT("/:id", h.RequirePermission(model.ResourceMethods, model.ActionUpdate), h.UpdateMethod)
		methods.DELETE("/:id", h.RequirePermission(model.ResourceMethods, model.ActionDelete), h.DeleteMethod)
	}

	plugins := api.Group("/plugins")
	plugins.Use(h.JWTAuth())
	{
		plugins.GET("", h.RequirePermission(model.ResourcePlugins, model.ActionRead), h.ListPluginGroups)
		plugins.GET("/:name", h.RequirePermission(model.ResourcePlugins, model.ActionRead), h.GetPluginGroup)
		plugins.POST("", h.RequirePermission(model.ResourcePlugins, model.ActionCreate), h.CreatePluginGroup)
		plugins.PUT("/:name", h.RequirePermission(model.ResourcePlugins, model.ActionUpdate), h.UpdatePluginGroup)
		plugins.DELETE("/:name", h.RequirePermission(model.ResourcePlugins, model.ActionDelete), h.DeletePluginGroup)
	}

	instances := api.Group("/instances")
	instances.Use(h.JWTAuth())
	{
		instances.GET("", h.RequirePermission(model.ResourceClusters, model.ActionRead), h.ListInstances)
		instances.GET("/stats", h.RequirePermission(model.ResourceClusters, model.ActionRead), h.GetInstanceStats)
	}

	users := api.Group("/users")
	users.Use(h.JWTAuth())
	{
		users.GET("", h.RequirePermission(model.ResourceUsers, model.ActionRead), h.ListUsers)
		users.GET("/:id", h.RequirePermission(model.ResourceUsers, model.ActionRead), h.GetUser)
		users.POST("", h.RequirePermission(model.ResourceUsers, model.ActionCreate), h.CreateUserByAdmin)
		users.PUT("/:id", h.RequirePermission(model.ResourceUsers, model.ActionUpdate), h.UpdateUser)
		users.DELETE("/:id", h.RequirePermission(model.ResourceUsers, model.ActionDelete), h.DeleteUser)
		users.POST("/:id/reset-password", h.RequirePermission(model.ResourceUsers, model.ActionUpdate), h.ResetUserPassword)
		users.POST("/:id/assign-role", h.RequirePermission(model.ResourceUsers, model.ActionUpdate), h.AssignUserRole)
	}

	roles := api.Group("/roles")
	roles.Use(h.JWTAuth())
	{
		roles.GET("", h.RequirePermission(model.ResourceUsers, model.ActionRead), h.ListRoles)
		roles.GET("/:id", h.RequirePermission(model.ResourceUsers, model.ActionRead), h.GetRoleByID)
		roles.POST("", h.RequirePermission(model.ResourceUsers, model.ActionCreate), h.CreateRole)
		roles.PUT("/:id", h.RequirePermission(model.ResourceUsers, model.ActionUpdate), h.UpdateRole)
		roles.DELETE("/:id", h.RequirePermission(model.ResourceUsers, model.ActionDelete), h.DeleteRole)
		roles.GET("/:id/permissions", h.RequirePermission(model.ResourceUsers, model.ActionRead), h.GetRolePermissions)
		roles.PUT("/:id/permissions", h.RequirePermission(model.ResourceUsers, model.ActionUpdate), h.UpdateRolePermissions)
	}

	permissions := api.Group("/permissions")
	permissions.Use(h.JWTAuth())
	{
		permissions.GET("", h.RequirePermission(model.ResourceUsers, model.ActionRead), h.ListPermissions)
	}
}
