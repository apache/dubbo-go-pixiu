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
import type { ApiResponse, Listener, ListenerDetail } from './types';

// Get listener list
export const getListenerList = (): Promise<ApiResponse<Listener[]>> => {
  return http.get('/api/listeners');
};

// Get listener detail by name (structured object)
export const getListenerDetail = (name: string): Promise<ApiResponse<Listener>> => {
  return http.get(`/api/listeners/${encodeURIComponent(name)}`);
};

// Get listener detail as YAML string (for editing)
export const getListenerYaml = (name: string): Promise<ApiResponse<string>> => {
  return http.get(`/api/listeners/${encodeURIComponent(name)}?format=yaml`);
};

// Create listener (POST for create per RESTful convention)
export const createListener = (content: string): Promise<ApiResponse<ListenerDetail>> => {
  const formData = createFormData({ content });
  return http.post('/api/listeners', formData);
};

// Update listener (PUT for update per RESTful convention)
export const updateListener = (name: string, content: string): Promise<ApiResponse<ListenerDetail>> => {
  const formData = createFormData({ content });
  return http.put(`/api/listeners/${encodeURIComponent(name)}`, formData);
};

// Delete listener by name
export const deleteListener = (name: string): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/listeners/${encodeURIComponent(name)}`);
};
