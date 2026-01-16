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

import axios from 'axios';
import { http, createFormData } from './request';
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  UserInfo,
  ChangePasswordRequest,
} from './types';

// Login - does not require token, use separate axios instance
export const login = async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  const formData = createFormData({
    username: data.username,
    password: data.password,
  });

  const response = await axios.post<ApiResponse<LoginResponse>>('/api/auth/login', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

// Register - does not require token
export const register = async (data: LoginRequest): Promise<ApiResponse<unknown>> => {
  const formData = createFormData({
    username: data.username,
    password: data.password,
  });

  const response = await axios.post<ApiResponse<unknown>>('/api/auth/register', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

// Logout
export const logout = (): Promise<ApiResponse<unknown>> => {
  return http.post('/api/user/logout');
};

// Get user info
export const getUserInfo = (): Promise<ApiResponse<UserInfo>> => {
  return http.get('/api/user/info');
};

// Get user role
export const getUserRole = (): Promise<ApiResponse<string>> => {
  return http.get('/api/user/role');
};

// Check if user is admin
export const checkIsAdmin = (): Promise<ApiResponse<boolean>> => {
  return http.get('/api/user/is-admin');
};

// Change password
export const changePassword = (data: ChangePasswordRequest): Promise<ApiResponse<unknown>> => {
  const formData = createFormData({
    oldPassword: data.oldPassword,
    newPassword: data.newPassword,
  });
  return http.post('/api/user/password', formData);
};
