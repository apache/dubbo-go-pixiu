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
import type { ApiResponse, Cluster, ClusterDetail } from './types';

// Get cluster list
export const getClusterList = (): Promise<ApiResponse<Cluster[]>> => {
  return http.get('/api/clusters');
};

// Get cluster detail by name (structured object)
export const getClusterDetail = (name: string): Promise<ApiResponse<Cluster>> => {
  return http.get(`/api/clusters/${encodeURIComponent(name)}`);
};

// Get cluster detail as YAML string (for editing)
export const getClusterYaml = (name: string): Promise<ApiResponse<string>> => {
  return http.get(`/api/clusters/${encodeURIComponent(name)}?format=yaml`);
};

// Create cluster (POST for create per RESTful convention)
export const createCluster = (content: string): Promise<ApiResponse<ClusterDetail>> => {
  const formData = createFormData({ content });
  return http.post('/api/clusters', formData);
};

// Update cluster (PUT for update per RESTful convention)
export const updateCluster = (name: string, content: string): Promise<ApiResponse<ClusterDetail>> => {
  const formData = createFormData({ content });
  return http.put(`/api/clusters/${encodeURIComponent(name)}`, formData);
};

// Delete cluster by name
export const deleteCluster = (name: string): Promise<ApiResponse<unknown>> => {
  return http.delete(`/api/clusters/${encodeURIComponent(name)}`);
};
