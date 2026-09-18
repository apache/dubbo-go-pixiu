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

export type RouteBindingParam = {
  from: string
  to: number
  type: string
}

export type AdminRouteBindingObject = {
  kind: string
  metadata: {
    name: string
  }
  spec: {
    entry: {
      protocol: string
      path: string
      method: string
    }
    target: {
      protocol: string
      application: string
      interface: string
      method: string
      version: string
      group: string
      cluster: string
    }
    params: RouteBindingParam[]
    enabled: boolean
    publish: {
      mode: string
      validate: boolean
    }
    extensions: JsonObject
  }
}

export type RouteBinding = {
  object: AdminRouteBindingObject
  resourceId: number
  methodId: number
  revision: number
}

export type RouteBindingValidationIssue = {
  path: string
  code: string
  message: string
}

export type RouteBindingPreview = {
  object: AdminRouteBindingObject
  yaml: string
}

export type RouteBindingPublishStatus = {
  name?: string
  draftRevision: number
  publishedRevision: number
  draftExists?: boolean
  publishedExists?: boolean
  dirty?: boolean
}

export type RouteBindingDiffChange = {
  path: string
  before: unknown
  after: unknown
}

export type RouteBindingDiff = {
  name: string
  draft?: RouteBinding
  published?: RouteBinding
  changes: RouteBindingDiffChange[]
}

export type RouteBindingPublishResult = {
  name: string
  revision: number
  draftRevision: number
  publishedRevision: number
  publishedCount: number
  deletedCount: number
}

export type ApiEnvelope<T> = { code: string; data: T }
