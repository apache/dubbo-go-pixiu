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
import type { ApiResponse, Resource, ResourceDetail, Method, MethodDetail } from './types';

// ==================== Resource (Mapping) ====================

// Get resource list
export const getResourceList = (): Promise<ApiResponse<Resource[]>> => {
  return http.get('/api/resources');
};

// Get resource detail by ID (structured object)
export const getResourceDetail = (id: string): Promise<ApiResponse<Resource>> => {
  return http.get(`/api/resources/${encodeURIComponent(id)}`);
};

// Get resource detail as YAML string (for editing)
export const getResourceYaml = (id: string): Promise<ApiResponse<string>> => {
  return http.get(`/api/resources/${encodeURIComponent(id)}?format=yaml`);
};

// Create resource (POST for create per RESTful convention)
export const createResource = (content: string): Promise<ApiResponse<ResourceDetail>> => {
  const formData = createFormData({ content });
  return http.post('/api/resources', formData);
};

// Update resource (PUT for update per RESTful convention)
export const updateResource = (id: string, content: string): Promise<ApiResponse<ResourceDetail>> => {
  const formData = createFormData({ content });
  return http.put(`/api/resources/${encodeURIComponent(id)}`, formData);
};

// Delete resource by ID
export const deleteResource = (id: string): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/resources/${encodeURIComponent(id)}`);
};

// ==================== Method ====================

// Get method list
export const getMethodList = (resourceId: string): Promise<ApiResponse<Method[]>> => {
  return http.get('/api/methods', { resourceId });
};

// Get method detail by ID
export const getMethodDetail = (id: string, resourceId: string): Promise<ApiResponse<string>> => {
  return http.get(`/api/methods/${encodeURIComponent(id)}`, { resourceId });
};

// Create method (POST for create per RESTful convention)
export const createMethod = (content: string, resourceId: string): Promise<ApiResponse<MethodDetail>> => {
  const formData = createFormData({ content });
  return http.post('/api/methods', formData, { resourceId });
};

// Update method (PUT for update per RESTful convention)
export const updateMethod = (id: string, content: string, resourceId: string): Promise<ApiResponse<MethodDetail>> => {
  const formData = createFormData({ content });
  return http.put(`/api/methods/${encodeURIComponent(id)}`, formData, { resourceId });
};

// Delete method by ID
export const deleteMethod = (id: string, resourceId: string): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/methods/${encodeURIComponent(id)}`, { resourceId });
};
