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
	"strconv"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/model"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/i18n"
)

func (h *Handler) Register(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.mysql.Register(req.Username, req.Password); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgUserRegistered)
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		config.Error(c, err)
		return
	}
	if _, err := h.mysql.Login(req.Username, req.Password); err != nil {
		config.Error(c, err)
		return
	}
	token, err := h.createToken(req.Username)
	if err != nil {
		config.Error(c, err)
		return
	}
	h.logger.Info("user logged in: " + req.Username)
	config.Success(c, model.LoginResponse{
		Username: req.Username,
		Token:    token,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	token, err := h.createExpiredToken()
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, token)
}

func (h *Handler) EditPassword(c *gin.Context) {
	var req model.PasswordChangeRequest
	if err := c.ShouldBind(&req); err != nil {
		config.Error(c, err)
		return
	}
	username := h.getUsername(c)
	if err := h.mysql.ChangePassword(username, req.OldPassword, req.NewPassword); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgPasswordChanged)
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	username := h.getUsername(c)
	info, err := h.mysql.GetUserInfo(username)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, info)
}

func (h *Handler) GetUserRole(c *gin.Context) {
	username := h.getUsername(c)
	role, err := h.mysql.GetUserRole(username)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, role)
}

func (h *Handler) CheckUserIsAdmin(c *gin.Context) {
	username := h.getUsername(c)
	isAdmin, err := h.mysql.IsAdmin(username)
	if err != nil {
		config.Error(c, err)
		return
	}
	if !isAdmin {
		config.ErrorMessageI18n(c, i18n.MsgUserNotAdmin)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgUserIsAdmin)
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, pageSize := 1, 20
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("pageSize")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}
	result, err := h.mysql.ListUsers(page, pageSize)
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, result)
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	user, err := h.mysql.GetUserByID(uint(id))
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.mysql.UpdateUser(uint(id), &req); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgUserUpdated)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	if err := h.mysql.DeleteUser(uint(id)); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgUserDeleted)
}

func (h *Handler) ResetUserPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.mysql.ResetUserPassword(uint(id), req.NewPassword); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgPasswordReset)
}

func (h *Handler) AssignUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	var req model.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Error(c, err)
		return
	}
	if err := h.mysql.AssignUserRole(uint(id), req.RoleID); err != nil {
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgRoleAssigned)
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.mysql.ListRoles()
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, roles)
}

func (h *Handler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	role, err := h.mysql.GetRole(uint(id))
	if err != nil {
		config.Error(c, err)
		return
	}
	config.Success(c, role)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.ErrorMessageI18n(c, i18n.MsgInvalidParams)
		return
	}
	role := &model.Role{
		RoleName:    req.RoleName,
		Description: req.Description,
	}
	if err := h.mysql.CreateRole(role); err != nil {
		h.logger.Error("failed to create role: " + err.Error())
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgRoleCreated)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.ErrorMessageI18n(c, i18n.MsgInvalidParams)
		return
	}
	if err := h.mysql.UpdateRole(uint(id), req.RoleName, req.Description); err != nil {
		h.logger.Error("failed to update role: " + err.Error())
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgRoleUpdated)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	if err := h.mysql.DeleteRole(uint(id)); err != nil {
		h.logger.Error("failed to delete role: " + err.Error())
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgRoleDeleted)
}

func (h *Handler) ListPermissions(c *gin.Context) {
	permissions, err := h.mysql.ListPermissions()
	if err != nil {
		h.logger.Error("failed to list permissions: " + err.Error())
		config.Error(c, err)
		return
	}
	config.Success(c, permissions)
}

func (h *Handler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	permissions, err := h.mysql.GetRolePermissions(uint(id))
	if err != nil {
		h.logger.Error("failed to get role permissions: " + err.Error())
		config.Error(c, err)
		return
	}
	config.Success(c, permissions)
}

func (h *Handler) UpdateRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		config.Error(c, err)
		return
	}
	var req model.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.ErrorMessageI18n(c, i18n.MsgInvalidParams)
		return
	}
	if err := h.mysql.UpdateRolePermissions(uint(id), req.PermissionIDs); err != nil {
		h.logger.Error("failed to update role permissions: " + err.Error())
		config.Error(c, err)
		return
	}
	config.SuccessMessageI18n(c, i18n.MsgRolePermissionsUpdated)
}

func (h *Handler) CreateUserByAdmin(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.ErrorMessageI18n(c, i18n.MsgInvalidParams)
		return
	}
	user, err := h.mysql.CreateUserByAdmin(req.Username, req.Password, req.RoleID)
	if err != nil {
		h.logger.Error("failed to create user: " + err.Error())
		config.Error(c, err)
		return
	}
	config.Success(c, model.UserListItem{
		ID:          user.ID,
		Username:    user.Username,
		Role:        user.Role,
		Enabled:     user.Enabled,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	})
}

func (h *Handler) getUsername(c *gin.Context) string {
	if username := c.GetString("username"); username != "" {
		return username
	}
	return c.Request.Header.Get("username")
}
