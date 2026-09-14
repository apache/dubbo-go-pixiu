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

export type JsonObject = Record<string, unknown>
export function asJsonObject(value: unknown): JsonObject {
  return value && typeof value === 'object' && !Array.isArray(value) ? (value as JsonObject) : {}
}
export type Resource = JsonObject & {
  id?: string
  name?: string
  path?: string
  description?: string
  type?: string
  timeout?: string
  methods?: Method[]
}
export type Method = JsonObject & {
  id?: string
  httpVerb?: string
  resourcePath?: string
  onAir?: boolean
  timeout?: string
  inboundRequest?: unknown
  integrationRequest?: unknown
  plugins?: unknown
}
export type ApiEnvelope<T> = { code: string; data: T }
