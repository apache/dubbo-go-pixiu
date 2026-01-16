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

package model

// Permission defines an action that can be performed on a resource.
type Permission struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Resource string `json:"resource" gorm:"type:varchar(50);not null;index:idx_resource_action"`
	Action   string `json:"action" gorm:"type:varchar(20);not null;index:idx_resource_action"`
}

// TableName specifies the table name for Permission.
func (Permission) TableName() string {
	return "pixiu_permission"
}

// RolePermission represents the many-to-many relationship between roles and permissions.
type RolePermission struct {
	RoleID       uint `json:"role_id" gorm:"primaryKey"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey"`
}

// TableName specifies the table name for RolePermission.
func (RolePermission) TableName() string {
	return "pixiu_role_permission"
}

// Resource constants for permission control.
const (
	ResourceClusters  = "clusters"
	ResourceListeners = "listeners"
	ResourceResources = "resources"
	ResourceMethods   = "methods"
	ResourcePlugins   = "plugins"
	ResourceUsers     = "users"
)

// Action constants for permission control.
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// DefaultPermissions defines the default permission set.
var DefaultPermissions = []Permission{
	// Clusters
	{Resource: ResourceClusters, Action: ActionCreate},
	{Resource: ResourceClusters, Action: ActionRead},
	{Resource: ResourceClusters, Action: ActionUpdate},
	{Resource: ResourceClusters, Action: ActionDelete},
	// Listeners
	{Resource: ResourceListeners, Action: ActionCreate},
	{Resource: ResourceListeners, Action: ActionRead},
	{Resource: ResourceListeners, Action: ActionUpdate},
	{Resource: ResourceListeners, Action: ActionDelete},
	// Resources
	{Resource: ResourceResources, Action: ActionCreate},
	{Resource: ResourceResources, Action: ActionRead},
	{Resource: ResourceResources, Action: ActionUpdate},
	{Resource: ResourceResources, Action: ActionDelete},
	// Methods
	{Resource: ResourceMethods, Action: ActionCreate},
	{Resource: ResourceMethods, Action: ActionRead},
	{Resource: ResourceMethods, Action: ActionUpdate},
	{Resource: ResourceMethods, Action: ActionDelete},
	// Plugins
	{Resource: ResourcePlugins, Action: ActionCreate},
	{Resource: ResourcePlugins, Action: ActionRead},
	{Resource: ResourcePlugins, Action: ActionUpdate},
	{Resource: ResourcePlugins, Action: ActionDelete},
	// Users (admin only)
	{Resource: ResourceUsers, Action: ActionCreate},
	{Resource: ResourceUsers, Action: ActionRead},
	{Resource: ResourceUsers, Action: ActionUpdate},
	{Resource: ResourceUsers, Action: ActionDelete},
}
