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

import { afterEach, expect, it, vi } from 'vitest'

vi.mock('./session', () => ({
  clearSession: vi.fn(),
  readSession: () => ({ username: '张三', token: 'test-token' }),
  redirectToLogin: vi.fn(),
}))

import { request } from './http'

afterEach(() => {
  vi.unstubAllGlobals()
})

it('sends the JWT without a username header for non-ASCII usernames', async () => {
  const fetchMock = vi
    .fn()
    .mockResolvedValue(new Response(JSON.stringify({ code: '10001', data: {} }), { status: 200 }))
  vi.stubGlobal('fetch', fetchMock)

  await request('/user/getInfo', { method: 'POST' })

  const requestInit = fetchMock.mock.calls[0][1] as RequestInit
  const headers = new Headers(requestInit.headers)
  expect(headers.get('token')).toBe('test-token')
  expect(headers.has('username')).toBe(false)
})
