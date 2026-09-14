/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { pixiuAdminApi } from '../api'
import { isNotFoundError, parseArrayResponse, request } from './http'
import type { JsonObject } from '../types/api'
export const resourceApi = {
  list: async () => {
    try {
      return parseArrayResponse<JsonObject>(await request<unknown>(pixiuAdminApi.resources.list))
    } catch (e: unknown) {
      if (isNotFoundError(e)) return []
      throw e
    }
  },
  detail: (id: string) =>
    request<unknown>(`${pixiuAdminApi.resources.detail}?resourceId=${encodeURIComponent(id)}`),
  create: (content: string) =>
    request<void>(pixiuAdminApi.resources.create, {
      method: 'POST',
      body: new URLSearchParams({ content }),
    }),
  update: (id: string, content: string) =>
    request<void>(`${pixiuAdminApi.resources.update}?resourceId=${encodeURIComponent(id)}`, {
      method: 'PUT',
      body: new URLSearchParams({ content }),
    }),
  remove: (id: string) =>
    request<void>(`${pixiuAdminApi.resources.remove}?resourceId=${encodeURIComponent(id)}`, {
      method: 'DELETE',
    }),
}
