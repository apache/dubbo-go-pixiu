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
import type { Method } from '../types/api'
export const methodApi = {
  list: async (resourceId: string) => {
    try {
      return parseArrayResponse<Method>(
        await request<unknown>(
          `${pixiuAdminApi.methods.list}?resourceId=${encodeURIComponent(resourceId)}`,
        ),
      )
    } catch (e: unknown) {
      if (isNotFoundError(e)) return []
      throw e
    }
  },
  detail: (resourceId: string, id: string) =>
    request<string>(
      `${pixiuAdminApi.methods.detail}?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,
    ),
  create: (resourceId: string, content: string) =>
    request<void>(`${pixiuAdminApi.methods.create}?resourceId=${encodeURIComponent(resourceId)}`, {
      method: 'POST',
      body: new URLSearchParams({ content }),
    }),
  update: (resourceId: string, id: string, content: string) =>
    request<void>(
      `${pixiuAdminApi.methods.update}?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,
      { method: 'PUT', body: new URLSearchParams({ content }) },
    ),
  remove: (resourceId: string, id: string) =>
    request<void>(
      `${pixiuAdminApi.methods.remove}?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,
      { method: 'DELETE' },
    ),
}
