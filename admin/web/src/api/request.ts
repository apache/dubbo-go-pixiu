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
import type { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { message } from 'antd';
import i18n from '../locales';
import { API_CODE } from './types';
import type { ApiResponse } from './types';

// Storage key - must match zustand persist key in auth store
const AUTH_STORAGE_KEY = 'pixiu-auth';

interface AuthStorageState {
  state: {
    token: string | null;
    username: string | null;
    isAuthenticated: boolean;
  };
}

// Get auth state from localStorage (zustand persist format)
const getAuthState = (): AuthStorageState['state'] | null => {
  const data = localStorage.getItem(AUTH_STORAGE_KEY);
  if (data) {
    try {
      const parsed = JSON.parse(data) as AuthStorageState;
      return parsed.state;
    } catch {
      return null;
    }
  }
  return null;
};

// Get token from localStorage
export const getToken = (): string | null => {
  return getAuthState()?.token || null;
};

// Get stored user info
export const getStoredUser = (): { username: string; token: string } | null => {
  const state = getAuthState();
  if (state?.token && state?.username) {
    return {
      token: state.token,
      username: state.username,
    };
  }
  return null;
};

// These are kept for compatibility but auth store should be used for mutations
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const setToken = (_token: string): void => {
  console.warn('setToken is deprecated, use useAuthStore instead');
};

export const removeToken = (): void => {
  localStorage.removeItem(AUTH_STORAGE_KEY);
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const setStoredUser = (_user: { username: string; token: string }): void => {
  console.warn('setStoredUser is deprecated, use useAuthStore instead');
};

export const removeStoredUser = (): void => {
  localStorage.removeItem(AUTH_STORAGE_KEY);
};

// Create axios instance
const request = axios.create({
  baseURL: '',
  timeout: 10000,
  headers: {
    'Content-Type': 'multipart/form-data',
  },
});

// Request interceptor
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const storedUser = getStoredUser();
    if (storedUser?.token) {
      config.headers.token = storedUser.token;
      config.headers.username = storedUser.username;
    }
    // Add Accept-Language header for i18n support
    config.headers['Accept-Language'] = i18n.language === 'zh' ? 'zh-CN' : 'en';
    return config;
  },
  (error: AxiosError) => {
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor
request.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse;

    // If response doesn't have code, return raw data
    if (!res || typeof res.code === 'undefined') {
      return response.data;
    }

    // Token expired
    if (res.code === API_CODE.TOKEN_EXPIRED) {
      message.error(i18n.t('common.tokenExpired'));
      removeToken();
      removeStoredUser();
      window.location.href = '/login';
      return Promise.reject(new Error('Token expired'));
    }

    return res;
  },
  (error: AxiosError) => {
    console.error('Response error:', error);

    // Get specific error message based on status code or error type
    let errorMsg: string;

    if (!error.response) {
      // Network error (no response)
      errorMsg = i18n.t('common.networkError');
    } else {
      const status = error.response.status;
      // Try to get error message from response data
      const responseData = error.response.data as { msg?: string; message?: string } | undefined;
      const serverMsg = responseData?.msg || responseData?.message;

      if (serverMsg) {
        errorMsg = serverMsg;
      } else {
        // Fallback to status-based messages
        switch (status) {
          case 400:
            errorMsg = i18n.t('common.badRequest');
            break;
          case 401:
            errorMsg = i18n.t('common.unauthorized');
            break;
          case 403:
            errorMsg = i18n.t('common.forbidden');
            break;
          case 404:
            errorMsg = i18n.t('common.notFound');
            break;
          case 500:
            errorMsg = i18n.t('common.serverError');
            break;
          case 502:
          case 503:
          case 504:
            errorMsg = i18n.t('common.serviceUnavailable');
            break;
          default:
            errorMsg = i18n.t('common.operationFailed');
        }
      }
    }

    message.error(errorMsg);
    return Promise.reject(error);
  }
);

// Helper methods for different request types
export const http = {
  get: <T>(url: string, params?: Record<string, unknown>): Promise<ApiResponse<T>> => {
    return request.get(url, { params }) as Promise<ApiResponse<T>>;
  },

  post: <T>(url: string, data?: FormData | Record<string, unknown>, params?: Record<string, unknown>): Promise<ApiResponse<T>> => {
    return request.post(url, data, { params }) as Promise<ApiResponse<T>>;
  },

  put: <T>(url: string, data?: FormData | Record<string, unknown>, params?: Record<string, unknown>): Promise<ApiResponse<T>> => {
    return request.put(url, data, { params }) as Promise<ApiResponse<T>>;
  },

  delete: <T>(url: string, params?: Record<string, unknown>): Promise<ApiResponse<T>> => {
    return request.delete(url, { params }) as Promise<ApiResponse<T>>;
  },
};

// Helper to create FormData from object
export const createFormData = (data: Record<string, string | Blob>): FormData => {
  const formData = new FormData();
  Object.entries(data).forEach(([key, value]) => {
    formData.append(key, value);
  });
  return formData;
};

export default request;
