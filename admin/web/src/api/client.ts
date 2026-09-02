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

import { CODE_OK, RetData } from '../types'
import { loadSession, touchSession } from '../auth/session'

export class ApiError extends Error {
  code: string

  constructor(code: string, message: string) {
    super(message)
    this.code = code
  }
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  /** Query string parameters. */
  query?: Record<string, string | number | undefined>
  /** application/x-www-form-urlencoded body fields (the Admin form convention). */
  form?: Record<string, string | number | undefined>
  /** Raw body (JSON/YAML), used by the route-binding endpoints. */
  rawBody?: string
  /** Attach token + username headers. Default true. */
  auth?: boolean
}

/**
 * Minimal fetch wrapper for the Pixiu Admin API.
 * Every handler returns HTTP 200; success is signalled by body code '10001'.
 */
export async function request<T = unknown>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', query, form, rawBody, auth = true } = options

  let url = path
  if (query) {
    const params = new URLSearchParams()
    for (const [key, value] of Object.entries(query)) {
      if (value !== undefined && value !== '') params.set(key, String(value))
    }
    const qs = params.toString()
    if (qs) url += (url.includes('?') ? '&' : '?') + qs
  }

  const headers: Record<string, string> = {}
  if (auth) {
    const session = loadSession()
    if (session) {
      // Config requests carry both token and username (backend convention).
      headers['token'] = session.token
      headers['username'] = session.username
    }
  }

  let body: string | undefined
  if (form) {
    const params = new URLSearchParams()
    for (const [key, value] of Object.entries(form)) {
      if (value !== undefined) params.set(key, String(value))
    }
    body = params.toString()
    headers['Content-Type'] = 'application/x-www-form-urlencoded'
  } else if (rawBody !== undefined) {
    body = rawBody
    headers['Content-Type'] = 'application/json'
  }

  let response: Response
  try {
    response = await fetch(url, { method, headers, body })
  } catch {
    throw new ApiError('NETWORK', '无法连接 Admin 服务，请确认后端已启动（默认 127.0.0.1:8081）')
  }

  const text = await response.text()
  let payload: RetData<T> | null = null
  try {
    payload = JSON.parse(text) as RetData<T>
  } catch {
    // Some Admin handlers write two envelopes back-to-back (an error envelope
    // followed by a success one). Honor the first so its message surfaces.
    const boundary = text.indexOf('}{')
    if (boundary > 0) {
      try {
        payload = JSON.parse(text.slice(0, boundary + 1)) as RetData<T>
      } catch {
        payload = null
      }
    }
  }
  if (!payload) {
    throw new ApiError('BAD_RESPONSE', `接口返回了非 JSON 内容（HTTP ${response.status}）`)
  }

  if (payload.code !== CODE_OK) {
    const message = typeof payload.data === 'string' ? payload.data : `请求失败（code ${payload.code}）`
    throw new ApiError(payload.code, message)
  }

  // Successful requests renew the local sliding session.
  if (auth) touchSession()
  return payload.data
}

/**
 * Some legacy endpoints double-encode collections: data is a JSON string that
 * must be parsed again (e.g. resource/method lists).
 */
export function parseJsonData<T>(data: unknown, fallback: T): T {
  if (typeof data !== 'string') return (data as T) ?? fallback
  try {
    return JSON.parse(data) as T
  } catch {
    return fallback
  }
}
