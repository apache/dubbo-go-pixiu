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

export type PixiuSession = { username: string; token: string }

export function readSession(): PixiuSession | null {
  const raw = localStorage.getItem('pixiu_session')
  if (!raw) return null
  try {
    const session = JSON.parse(raw) as Partial<PixiuSession>
    if (
      typeof session.username !== 'string' ||
      !session.username.trim() ||
      typeof session.token !== 'string' ||
      !session.token.trim()
    )
      return null
    return { username: session.username, token: session.token }
  } catch {
    return null
  }
}

export function clearSession() {
  localStorage.removeItem('pixiu_session')
}

export function redirectToLogin() {
  if (window.location.pathname !== '/login') window.location.assign('/login')
}
