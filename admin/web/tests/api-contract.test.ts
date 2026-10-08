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

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { pixiuAdminApi } from '../src/api'
import { routeBindingApi } from '../src/services/route-binding-api'

beforeEach(() => {
  const storage = new Map<string, string>()
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => storage.get(key) ?? null,
    setItem: (key: string, value: string) => storage.set(key, value),
    removeItem: (key: string) => storage.delete(key),
  })
})

afterEach(() => vi.unstubAllGlobals())

describe('Pixiu Admin API contract', () => {
  it('uses the independent API route lifecycle endpoints', () => {
    expect(pixiuAdminApi.routeBindings).toEqual({
      schema: '/config/api/route/schema',
      list: '/config/api/route/list',
      detail: '/config/api/route/detail',
      create: '/config/api/route',
      update: '/config/api/route',
      remove: '/config/api/route',
      validate: '/config/api/route/validate',
      preview: '/config/api/route/preview',
      publish: '/config/api/route/publish',
      status: '/config/api/route/status',
      diff: '/config/api/route/diff',
    })
    expect(JSON.stringify(pixiuAdminApi)).not.toContain('/config/api/resource')
    expect(JSON.stringify(pixiuAdminApi)).not.toContain('/config/api/resource/method')
  })

  it('uses the shared contract for API route list requests', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify({ code: '10001', data: [] }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await routeBindingApi.list('draft')

    expect(fetchMock).toHaveBeenCalledWith(
      `${pixiuAdminApi.routeBindings.list}?scope=draft`,
      expect.objectContaining({ headers: expect.any(Headers) }),
    )
  })

  it('includes the original route name when updating an API route', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify({ code: '10001', data: {} }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await routeBindingApi.update('user-get', {} as Parameters<typeof routeBindingApi.update>[1])

    expect(fetchMock).toHaveBeenCalledWith(
      `${pixiuAdminApi.routeBindings.update}?name=user-get`,
      expect.objectContaining({ method: 'PUT', headers: expect.any(Headers) }),
    )
  })

  it('surfaces malformed list responses instead of treating them as empty', async () => {
    vi.stubGlobal(
      'fetch',
      vi
        .fn()
        .mockResolvedValue(
          new Response(JSON.stringify({ code: '10001', data: 'not-an-array' }), { status: 200 }),
        ),
    )

    await expect(routeBindingApi.list()).rejects.toMatchObject({ code: 'BAD_RESPONSE' })
  })

  it('does not accept an error HTTP status as a successful API response', async () => {
    vi.stubGlobal(
      'fetch',
      vi
        .fn()
        .mockResolvedValue(
          new Response(JSON.stringify({ message: 'backend unavailable' }), { status: 503 }),
        ),
    )

    await expect(routeBindingApi.list()).rejects.toMatchObject({ code: 'HTTP_503' })
  })
})
