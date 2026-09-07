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

// Local session storage. The backend JWT expires after 1h; the frontend keeps
// a sliding local session that successful requests renew (see api/client).

const KEY_USERNAME = 'pixiu.admin.username'
const KEY_TOKEN = 'pixiu.admin.token'
const KEY_ACTIVE_AT = 'pixiu.admin.activeAt'

/** Sliding window, slightly shorter than the 1h JWT lifetime. */
const SESSION_TTL_MS = 55 * 60 * 1000

export interface Session {
  username: string
  token: string
}

export function saveSession(session: Session): void {
  localStorage.setItem(KEY_USERNAME, session.username)
  localStorage.setItem(KEY_TOKEN, session.token)
  localStorage.setItem(KEY_ACTIVE_AT, String(Date.now()))
}

export function loadSession(): Session | null {
  const username = localStorage.getItem(KEY_USERNAME)
  const token = localStorage.getItem(KEY_TOKEN)
  const activeAt = Number(localStorage.getItem(KEY_ACTIVE_AT) || 0)
  if (!username || !token) return null
  if (!activeAt || Date.now() - activeAt > SESSION_TTL_MS) {
    clearSession()
    return null
  }
  return { username, token }
}

/** Renew the sliding session window after a successful request. */
export function touchSession(): void {
  if (localStorage.getItem(KEY_TOKEN)) {
    localStorage.setItem(KEY_ACTIVE_AT, String(Date.now()))
  }
}

export function clearSession(): void {
  localStorage.removeItem(KEY_USERNAME)
  localStorage.removeItem(KEY_TOKEN)
  localStorage.removeItem(KEY_ACTIVE_AT)
}
