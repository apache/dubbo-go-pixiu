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

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  login as loginApi,
  logout as logoutApi,
  getUserInfo as getUserInfoApi,
  API_CODE,
} from '../api';
import type { UserInfo, LoginRequest } from '../api';

interface AuthState {
  // State
  token: string | null;
  username: string | null;
  userInfo: UserInfo | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  login: (data: LoginRequest) => Promise<boolean>;
  logout: () => Promise<void>;
  fetchUserInfo: () => Promise<void>;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // Initial state
      token: null,
      username: null,
      userInfo: null,
      isAuthenticated: false,
      isLoading: false,

      // Login action
      login: async (data: LoginRequest) => {
        set({ isLoading: true });
        try {
          const response = await loginApi(data);
          if (response.code === API_CODE.SUCCESS && response.data) {
            const { token, username } = response.data;
            set({
              token,
              username,
              isAuthenticated: true,
              isLoading: false,
            });
            return true;
          }
          set({ isLoading: false });
          return false;
        } catch (error) {
          console.error('Login failed:', error);
          set({ isLoading: false });
          return false;
        }
      },

      // Logout action
      logout: async () => {
        try {
          await logoutApi();
        } catch (error) {
          console.error('Logout error:', error);
        } finally {
          get().clearAuth();
        }
      },

      // Fetch user info
      fetchUserInfo: async () => {
        if (!get().isAuthenticated) return;

        try {
          const response = await getUserInfoApi();
          if (response.code === API_CODE.SUCCESS && response.data) {
            set({ userInfo: response.data });
          }
        } catch (error) {
          console.error('Failed to fetch user info:', error);
        }
      },

      // Clear auth state
      clearAuth: () => {
        set({
          token: null,
          username: null,
          userInfo: null,
          isAuthenticated: false,
        });
      },
    }),
    {
      name: 'pixiu-auth',
      partialize: (state) => ({
        token: state.token,
        username: state.username,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
