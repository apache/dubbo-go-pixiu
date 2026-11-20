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

import { API_BASE, API_STATUS } from '../app/constants'
import { clearSession, readSession, redirectToLogin } from './session'

export class ApiError extends Error {
  code: string
  raw: string
  issues: Array<{ path: string; code: string; message: string }>
  constructor(
    message: string,
    code = '',
    raw = '',
    issues: Array<{ path: string; code: string; message: string }> = [],
  ) {
    super(message)
    this.code = code
    this.raw = raw
    this.issues = issues
  }
}

function sessionHeaders() {
  const session = readSession()
  if (!session) return {}
  return { token: session.token, username: session.username }
}

function isAuthFailure(status: number, bodyData: unknown) {
  const body =
    bodyData && typeof bodyData === 'object' ? (bodyData as { code?: unknown; data?: unknown }) : {}
  if (status === 401 || status === 403) return true
  if (body.code !== API_STATUS.NOT_FOUND) return false
  const message = typeof body.data === 'string' ? body.data.toLowerCase() : ''
  return /token|login|access|authoriz|认证|登录|权限/.test(message)
}

function errorDetails(data: unknown) {
  if (!data || typeof data !== 'object' || Array.isArray(data))
    return { message: '', issues: [] as Array<{ path: string; code: string; message: string }> }
  const value = data as { message?: unknown; issues?: unknown }
  const issues = Array.isArray(value.issues)
    ? value.issues.filter(
        (issue): issue is { path: string; code: string; message: string } =>
          Boolean(issue) &&
          typeof issue === 'object' &&
          typeof (issue as { path?: unknown }).path === 'string' &&
          typeof (issue as { code?: unknown }).code === 'string' &&
          typeof (issue as { message?: unknown }).message === 'string',
      )
    : []
  return {
    message: typeof value.message === 'string' ? value.message : '',
    issues,
  }
}

function primitiveString(value: unknown, fallback = '') {
  if (typeof value === 'string') return value
  if (typeof value === 'number' || typeof value === 'boolean' || typeof value === 'bigint') {
    return value.toString()
  }
  return fallback
}

function responseDataString(value: unknown, fallback = '') {
  return primitiveString(value, fallback)
}

function httpErrorMessage(
  data: unknown,
  details: ReturnType<typeof errorDetails>,
  envelopeMessage: unknown,
  fallback: string,
) {
  if (typeof data === 'string') return data
  if (details.message) return details.message
  if (typeof envelopeMessage === 'string') return envelopeMessage
  return fallback
}

function codeErrorMessage(data: unknown, details: ReturnType<typeof errorDetails>, code: unknown) {
  if (typeof data === 'string') return data
  if (details.message) return details.message
  if (code === API_STATUS.NOT_FOUND) return '数据不存在'
  if (code === API_STATUS.CONCURRENT) return '操作冲突，请刷新后重试'
  return '请求失败'
}

export function isNotFoundError(error: unknown) {
  return error instanceof ApiError && error.code === API_STATUS.NOT_FOUND
}

export function parseArrayResponse<T>(data: unknown): T[] {
  let parsed = data
  if (typeof data === 'string') {
    try {
      parsed = JSON.parse(data) as unknown
    } catch (error) {
      throw new ApiError(
        '接口返回的数据不是有效 JSON',
        'BAD_RESPONSE',
        error instanceof Error ? error.message : primitiveString(error, '未知错误'),
      )
    }
  }
  if (!Array.isArray(parsed))
    throw new ApiError('接口返回的数据格式无效', 'BAD_RESPONSE', responseDataString(data))
  return parsed as T[]
}

export async function request<T>(path: string, init: RequestInit = {}) {
  const body = init.body
  const headers = new Headers(init.headers)
  Object.entries(sessionHeaders()).forEach(([key, value]) => headers.set(key, value))
  if (
    body &&
    !(body instanceof URLSearchParams) &&
    !(body instanceof FormData) &&
    !headers.has('Content-Type')
  )
    headers.set('Content-Type', 'application/json')
  const r = await fetch(API_BASE + path, { ...init, headers })
  const raw = await r.text()
  if (r.status === 401 || r.status === 403) {
    clearSession()
    redirectToLogin()
    throw new ApiError('登录状态已失效，请重新登录', 'AUTH_REQUIRED', raw)
  }
  if (r.status === 204) return undefined as T
  let bodyData: unknown
  try {
    bodyData = JSON.parse(raw)
  } catch {
    const boundary = raw.indexOf('}{')
    if (boundary > 0) {
      try {
        bodyData = JSON.parse(raw.slice(0, boundary + 1))
      } catch {
        bodyData = null
      }
    }
  }
  if (bodyData === null || typeof bodyData !== 'object')
    throw new ApiError(`接口返回了非 JSON 内容（HTTP ${r.status}）`, 'BAD_RESPONSE', raw)
  const envelope = bodyData as { code?: unknown; data?: unknown; message?: unknown }
  if (!r.ok) {
    if (isAuthFailure(r.status, envelope)) {
      clearSession()
      redirectToLogin()
      throw new ApiError(
        '登录状态已失效，请重新登录',
        'AUTH_REQUIRED',
        responseDataString(envelope.data, raw),
      )
    }
    const details = errorDetails(envelope.data)
    const message = httpErrorMessage(
      envelope.data,
      details,
      envelope.message,
      `请求失败（HTTP ${r.status}）`,
    )
    throw new ApiError(
      message,
      primitiveString(envelope.code, `HTTP_${r.status}`),
      raw,
      details.issues,
    )
  }
  if (envelope.code && envelope.code !== API_STATUS.SUCCESS) {
    if (isAuthFailure(r.status, envelope)) {
      clearSession()
      redirectToLogin()
      throw new ApiError(
        '登录状态已失效，请重新登录',
        'AUTH_REQUIRED',
        responseDataString(envelope.data),
      )
    }
    const details = errorDetails(envelope.data)
    const message = codeErrorMessage(envelope.data, details, envelope.code)
    throw new ApiError(
      message,
      primitiveString(envelope.code),
      responseDataString(envelope.data),
      details.issues,
    )
  }
  if (!('data' in envelope)) throw new ApiError('接口响应缺少 data 字段', 'BAD_RESPONSE', raw)
  return envelope.data as T
}
