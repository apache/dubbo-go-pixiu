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

package i18n

// Message key constants
const (
	// Cluster operations
	MsgClusterCreated = "cluster.created"
	MsgClusterUpdated = "cluster.updated"
	MsgClusterDeleted = "cluster.deleted"

	// Listener operations
	MsgListenerCreated = "listener.created"
	MsgListenerUpdated = "listener.updated"
	MsgListenerDeleted = "listener.deleted"

	// Resource operations
	MsgResourceCreated = "resource.created"
	MsgResourceUpdated = "resource.updated"
	MsgResourceDeleted = "resource.deleted"

	// Method operations
	MsgMethodCreated = "method.created"
	MsgMethodUpdated = "method.updated"
	MsgMethodDeleted = "method.deleted"

	// Plugin operations
	MsgPluginCreated = "plugin.created"
	MsgPluginUpdated = "plugin.updated"
	MsgPluginDeleted = "plugin.deleted"

	// User operations
	MsgUserRegistered  = "user.registered"
	MsgUserCreated     = "user.created"
	MsgPasswordChanged = "user.password_changed"
	MsgUserNotAdmin    = "user.not_admin"
	MsgUserIsAdmin     = "user.is_admin"
	MsgLoginSuccess    = "user.login_success"
	MsgLogoutSuccess   = "user.logout_success"
	MsgUserUpdated     = "user.updated"
	MsgUserDeleted     = "user.deleted"
	MsgPasswordReset   = "user.password_reset"
	MsgRoleAssigned    = "user.role_assigned"

	// Role operations
	MsgRoleCreated            = "role.created"
	MsgRoleUpdated            = "role.updated"
	MsgRoleDeleted            = "role.deleted"
	MsgRolePermissionsUpdated = "role.permissions_updated"

	// Authentication
	MsgTokenMissing     = "auth.token_missing"
	MsgTokenExpired     = "auth.token_expired"
	MsgNotAuthenticated = "auth.not_authenticated"
	MsgPermissionDenied = "auth.permission_denied"
	MsgCheckPermFailed  = "auth.check_permission_failed"

	// System errors
	MsgXDSNotInitialized = "error.xds_not_init"

	// Validation errors
	MsgInvalidParams = "error.invalid_params"
)

// messages contains all translations organized by language
var messages = map[string]map[string]string{
	LangEnglish: {
		// Cluster operations
		MsgClusterCreated: "cluster created successfully",
		MsgClusterUpdated: "cluster updated successfully",
		MsgClusterDeleted: "cluster deleted successfully",

		// Listener operations
		MsgListenerCreated: "listener created successfully",
		MsgListenerUpdated: "listener updated successfully",
		MsgListenerDeleted: "listener deleted successfully",

		// Resource operations
		MsgResourceCreated: "resource created successfully",
		MsgResourceUpdated: "resource updated successfully",
		MsgResourceDeleted: "resource deleted successfully",

		// Method operations
		MsgMethodCreated: "method created successfully",
		MsgMethodUpdated: "method updated successfully",
		MsgMethodDeleted: "method deleted successfully",

		// Plugin operations
		MsgPluginCreated: "plugin group created successfully",
		MsgPluginUpdated: "plugin group updated successfully",
		MsgPluginDeleted: "plugin group deleted successfully",

		// User operations
		MsgUserRegistered:  "register successfully, please login",
		MsgUserCreated:     "user created successfully",
		MsgPasswordChanged: "password changed successfully",
		MsgUserNotAdmin:    "user is not an admin",
		MsgUserIsAdmin:     "user is admin",
		MsgLoginSuccess:    "login successful",
		MsgLogoutSuccess:   "logout successful",
		MsgUserUpdated:     "user updated successfully",
		MsgUserDeleted:     "user deleted successfully",
		MsgPasswordReset:   "password reset successfully",
		MsgRoleAssigned:    "role assigned successfully",

		// Role operations
		MsgRoleCreated:            "role created successfully",
		MsgRoleUpdated:            "role updated successfully",
		MsgRoleDeleted:            "role deleted successfully",
		MsgRolePermissionsUpdated: "role permissions updated successfully",

		// Authentication
		MsgTokenMissing:     "request does not carry token",
		MsgTokenExpired:     "token has expired, please login again",
		MsgNotAuthenticated: "user not authenticated",
		MsgPermissionDenied: "permission denied",
		MsgCheckPermFailed:  "failed to check permission",

		// System errors
		MsgXDSNotInitialized: "xDS server not initialized",

		// Validation errors
		MsgInvalidParams: "invalid request parameters",
	},
	LangChinese: {
		// Cluster operations
		MsgClusterCreated: "集群创建成功",
		MsgClusterUpdated: "集群更新成功",
		MsgClusterDeleted: "集群删除成功",

		// Listener operations
		MsgListenerCreated: "监听器创建成功",
		MsgListenerUpdated: "监听器更新成功",
		MsgListenerDeleted: "监听器删除成功",

		// Resource operations
		MsgResourceCreated: "资源创建成功",
		MsgResourceUpdated: "资源更新成功",
		MsgResourceDeleted: "资源删除成功",

		// Method operations
		MsgMethodCreated: "方法创建成功",
		MsgMethodUpdated: "方法更新成功",
		MsgMethodDeleted: "方法删除成功",

		// Plugin operations
		MsgPluginCreated: "插件组创建成功",
		MsgPluginUpdated: "插件组更新成功",
		MsgPluginDeleted: "插件组删除成功",

		// User operations
		MsgUserRegistered:  "注册成功，请登录",
		MsgUserCreated:     "用户创建成功",
		MsgPasswordChanged: "密码修改成功",
		MsgUserNotAdmin:    "用户不是管理员",
		MsgUserIsAdmin:     "用户是管理员",
		MsgLoginSuccess:    "登录成功",
		MsgLogoutSuccess:   "退出登录成功",
		MsgUserUpdated:     "用户更新成功",
		MsgUserDeleted:     "用户删除成功",
		MsgPasswordReset:   "密码重置成功",
		MsgRoleAssigned:    "角色分配成功",

		// Role operations
		MsgRoleCreated:            "角色创建成功",
		MsgRoleUpdated:            "角色更新成功",
		MsgRoleDeleted:            "角色删除成功",
		MsgRolePermissionsUpdated: "角色权限更新成功",

		// Authentication
		MsgTokenMissing:     "请求未携带令牌",
		MsgTokenExpired:     "令牌已过期，请重新登录",
		MsgNotAuthenticated: "用户未认证",
		MsgPermissionDenied: "权限不足",
		MsgCheckPermFailed:  "权限检查失败",

		// System errors
		MsgXDSNotInitialized: "xDS 服务器未初始化",

		// Validation errors
		MsgInvalidParams: "无效的请求参数",
	},
}
