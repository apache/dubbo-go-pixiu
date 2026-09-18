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
import type {
  AdminRouteBindingObject,
  RouteBinding,
  RouteBindingDiff,
  RouteBindingPreview,
  RouteBindingPublishResult,
  RouteBindingPublishStatus,
} from '../types/api'

type RouteBindingScope = 'draft' | 'published'

function mutationBody(object: AdminRouteBindingObject, expectedRevision?: number) {
  return JSON.stringify(
    expectedRevision && expectedRevision > 0 ? { object, expectedRevision } : object,
  )
}

function query(values: Record<string, string | number | undefined>) {
  const params = new URLSearchParams()
  Object.entries(values).forEach(([key, value]) => {
    if (value !== undefined) params.set(key, String(value))
  })
  const encoded = params.toString()
  return encoded ? `?${encoded}` : ''
}

export const routeBindingApi = {
  schema: async () =>
    parseArrayResponse(await request<unknown>(pixiuAdminApi.routeBindings.schema)),
  list: async (scope: RouteBindingScope = 'draft') => {
    try {
      return parseArrayResponse<RouteBinding>(
        await request<unknown>(`${pixiuAdminApi.routeBindings.list}${query({ scope })}`),
      )
    } catch (error: unknown) {
      if (isNotFoundError(error)) return []
      throw error
    }
  },
  detail: (name: string, scope: RouteBindingScope = 'draft') =>
    request<RouteBinding>(`${pixiuAdminApi.routeBindings.detail}${query({ name, scope })}`),
  create: (object: AdminRouteBindingObject) =>
    request<RouteBinding>(pixiuAdminApi.routeBindings.create, {
      method: 'POST',
      body: mutationBody(object),
    }),
  update: (object: AdminRouteBindingObject, expectedRevision = 0) =>
    request<RouteBinding>(pixiuAdminApi.routeBindings.update, {
      method: 'PUT',
      body: mutationBody(object, expectedRevision),
    }),
  remove: (name: string, expectedRevision = 0) =>
    request<string>(
      `${pixiuAdminApi.routeBindings.remove}${query({ name, expectedRevision: expectedRevision || undefined })}`,
      { method: 'DELETE' },
    ),
  validate: (object: AdminRouteBindingObject) =>
    request<{ valid: boolean; object: AdminRouteBindingObject }>(
      pixiuAdminApi.routeBindings.validate,
      { method: 'POST', body: mutationBody(object) },
    ),
  preview: (object: AdminRouteBindingObject) =>
    request<RouteBindingPreview>(pixiuAdminApi.routeBindings.preview, {
      method: 'POST',
      body: mutationBody(object),
    }),
  publish: (name: string, expectedRevision = 0) =>
    request<RouteBindingPublishResult>(
      `${pixiuAdminApi.routeBindings.publish}${query({
        name,
        expectedRevision: expectedRevision || undefined,
      })}`,
      { method: 'PUT' },
    ),
  status: (name: string) =>
    request<RouteBindingPublishStatus>(`${pixiuAdminApi.routeBindings.status}${query({ name })}`),
  diff: (name: string) =>
    request<RouteBindingDiff>(`${pixiuAdminApi.routeBindings.diff}${query({ name })}`),
  publishStatus: () =>
    request<RouteBindingPublishStatus>(pixiuAdminApi.routeBindings.publishStatus),
}
