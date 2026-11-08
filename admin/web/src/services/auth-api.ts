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

import { API_BASE } from '../app/constants'
export async function login(username: string, password: string) {
  const body = new URLSearchParams({ username, password })
  const r = await fetch(API_BASE + '/login', { method: 'POST', body })
  const data = await r.json()
  if (data.code !== '10001') throw new Error(data.data || '登录失败')
  return data.data as { username: string; token: string }
}
