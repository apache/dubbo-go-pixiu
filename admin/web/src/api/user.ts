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

import { http } from './request';
import type { ApiResponse, UserListResponse, UserListItem, Role, UpdateUserRequest, Permission } from './types';

// Get user list with pagination
export const getUserList = (page = 1, pageSize = 20): Promise<ApiResponse<UserListResponse>> => {
  return http.get(`/api/users?page=${page}&pageSize=${pageSize}`);
};

// Get user by ID
export const getUser = (id: number): Promise<ApiResponse<UserListItem>> => {
  return http.get(`/api/users/${id}`);
};

// Update user
export const updateUser = (id: number, data: UpdateUserRequest): Promise<ApiResponse<unknown>> => {
  return http.put(`/api/users/${id}`, data as Record<string, unknown>);
};

// Delete user
export const deleteUser = (id: number): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/users/${id}`);
};

// Reset user password (admin operation)
export const resetUserPassword = (id: number, newPassword: string): Promise<ApiResponse<unknown>> => {
  return http.post(`/api/users/${id}/reset-password`, { newPassword });
};

// Assign role to user
export const assignUserRole = (userId: number, roleId: number): Promise<ApiResponse<unknown>> => {
  return http.post(`/api/users/${userId}/assign-role`, { roleId });
};

// Get role list
export const getRoleList = (): Promise<ApiResponse<Role[]>> => {
  return http.get('/api/roles');
};

// Get role by ID
export const getRole = (id: number): Promise<ApiResponse<Role>> => {
  return http.get(`/api/roles/${id}`);
};

// Create role
export const createRole = (data: { roleName: string; description?: string }): Promise<ApiResponse<unknown>> => {
  return http.post('/api/roles', data);
};

// Update role
export const updateRole = (id: number, data: { roleName: string; description?: string }): Promise<ApiResponse<unknown>> => {
  return http.put(`/api/roles/${id}`, data);
};

// Delete role
export const deleteRole = (id: number): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/roles/${id}`);
};

// Get role permissions
export const getRolePermissions = (roleId: number): Promise<ApiResponse<Permission[]>> => {
  return http.get(`/api/roles/${roleId}/permissions`);
};

// Update role permissions
export const updateRolePermissions = (roleId: number, permissionIds: number[]): Promise<ApiResponse<unknown>> => {
  return http.put(`/api/roles/${roleId}/permissions`, { permissionIds });
};

// Get all permissions
export const getPermissionList = (): Promise<ApiResponse<Permission[]>> => {
  return http.get('/api/permissions');
};

// Create user (admin operation)
export const createUser = (data: { username: string; password: string; roleId?: number }): Promise<ApiResponse<UserListItem>> => {
  return http.post('/api/users', data);
};
