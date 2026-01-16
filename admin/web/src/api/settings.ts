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
import type { ApiResponse } from './types';

export interface SystemSettings {
  language?: string;
  darkMode?: boolean;
  pageSize?: number;
  autoRefresh?: boolean;
  refreshInterval?: number;
  [key: string]: string | boolean | number | undefined;
}

export interface GatewayInfo {
  version?: string;
  status?: string;
  httpPort?: string;
  xdsPort?: string;
  startTime?: string;
  uptime?: string;
  logLevel?: string;
  maxConnections?: number;
  requestTimeout?: number;
}

// Get system settings
export const getSystemSettings = (): Promise<ApiResponse<SystemSettings>> => {
  return http.get('/api/system/settings');
};

// Update system settings
export const updateSystemSettings = (data: SystemSettings): Promise<ApiResponse<unknown>> => {
  return http.put('/api/system/settings', data);
};

// Get gateway info
export const getGatewayInfo = (): Promise<ApiResponse<GatewayInfo>> => {
  return http.get('/api/system/gateway');
};

// Update gateway settings
export const updateGatewaySettings = (data: Partial<GatewayInfo>): Promise<ApiResponse<unknown>> => {
  return http.put('/api/system/gateway', data);
};
