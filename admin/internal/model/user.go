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

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password    string    `json:"-" gorm:"type:varchar(64);not null"`
	Role        int       `json:"role" gorm:"default:0"`
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	DateCreated time.Time `json:"date_created" gorm:"autoCreateTime"`
	DateUpdated time.Time `json:"date_updated" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for User.
func (User) TableName() string {
	return "pixiu_user"
}

// UserInfo is a safe representation of user data for API responses.
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}

// Role represents a user role.
type Role struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleName    string `json:"role_name" gorm:"type:varchar(50);not null"`
	Description string `json:"description" gorm:"type:varchar(200)"`
}

// TableName specifies the table name for Role.
func (Role) TableName() string {
	return "pixiu_role"
}

// UserRole represents the many-to-many relationship between users and roles.
type UserRole struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	RoleID uint `json:"role_id" gorm:"primaryKey"`
}

// TableName specifies the table name for UserRole.
func (UserRole) TableName() string {
	return "pixiu_user_role"
}

// LoginRequest represents login request payload.
type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// LoginResponse represents login response payload.
type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// PasswordChangeRequest represents password change request payload.
type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" form:"oldPassword" binding:"required"`
	NewPassword string `json:"new_password" form:"newPassword" binding:"required"`
}

// UserListItem is a safe representation of user data for list API responses.
type UserListItem struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Role        int       `json:"role"`
	Enabled     bool      `json:"enabled"`
	DateCreated time.Time `json:"dateCreated"`
	DateUpdated time.Time `json:"dateUpdated"`
}

// UserListResponse represents the response for user list API.
type UserListResponse struct {
	Items    []UserListItem `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// UpdateUserRequest represents the request for updating user.
type UpdateUserRequest struct {
	Role    *int  `json:"role"`
	Enabled *bool `json:"enabled"`
}

// ResetPasswordRequest represents the request for resetting user password.
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required"`
}

// AssignRoleRequest represents the request for assigning role to user.
type AssignRoleRequest struct {
	RoleID uint `json:"roleId" binding:"required"`
}

// CreateRoleRequest represents the request for creating a role.
type CreateRoleRequest struct {
	RoleName    string `json:"roleName" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest represents the request for updating a role.
type UpdateRoleRequest struct {
	RoleName    string `json:"roleName" binding:"required"`
	Description string `json:"description"`
}

// UpdateRolePermissionsRequest represents the request for updating role permissions.
type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permissionIds" binding:"required"`
}

// CreateUserRequest represents the request for creating a user (admin operation).
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"roleId"`
}
