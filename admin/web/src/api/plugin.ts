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

import { http, createFormData } from './request';
import type { ApiResponse, PluginGroupItem, PluginGroupDetail } from './types';

// Get plugin group list
export const getPluginGroupList = (): Promise<ApiResponse<PluginGroupItem[]>> => {
  return http.get('/api/plugins');
};

// Get plugin group detail by name (structured object)
export const getPluginGroupDetail = (name: string): Promise<ApiResponse<PluginGroupItem>> => {
  return http.get(`/api/plugins/${encodeURIComponent(name)}`);
};

// Get plugin group detail as YAML string (for editing)
export const getPluginGroupYaml = (name: string): Promise<ApiResponse<string>> => {
  return http.get(`/api/plugins/${encodeURIComponent(name)}?format=yaml`);
};

// Create plugin group (POST for create per RESTful convention)
export const createPluginGroup = (content: string): Promise<ApiResponse<PluginGroupDetail>> => {
  const formData = createFormData({ content });
  return http.post('/api/plugins', formData);
};

// Update plugin group (PUT for update per RESTful convention)
export const updatePluginGroup = (name: string, content: string): Promise<ApiResponse<PluginGroupDetail>> => {
  const formData = createFormData({ content });
  return http.put(`/api/plugins/${encodeURIComponent(name)}`, formData);
};

// Delete plugin group by name
export const deletePluginGroup = (name: string): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/plugins/${encodeURIComponent(name)}`);
};
