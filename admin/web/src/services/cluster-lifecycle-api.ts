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

import { parseArrayResponse, request } from './http'

export type AdminClusterObject = {
  kind: 'AdminCluster'
  metadata: { name: string }
  spec: {
    type?: string
    endpoint?: { address?: string; port?: number | string }
    [key: string]: unknown
  }
}
export type ClusterDraft = {
  id: string
  object: AdminClusterObject
  sourceId?: number
  revision: number
}
export type ClusterStatus = {
  id: string
  sourceId?: number
  draftRevision: number
  publishedRevision: number
  draftExists: boolean
  publishedExists: boolean
  dirty: boolean
}
export type ClusterChange = { path: string; before: unknown; after: unknown }
export type ClusterDiff = {
  id: string
  draft: ClusterDraft
  published?: AdminClusterObject
  changes: ClusterChange[]
}
export type ClusterHistory = {
  revision: number
  publishedAt: string
  runtimeId: number
  object: AdminClusterObject
}
export type ClusterPublishResult = {
  id: string
  runtimeId: number
  draftRevision: number
  publishedRevision: number
  historyRevision: number
}
export type ClusterFieldSchema = {
  type?: string
  default?: string
  enum?: string[]
  minimum?: number
  maximum?: number
  required?: boolean
  ui?: { placeholder?: string; order?: number }
  properties?: Record<string, ClusterFieldSchema>
}
export type ClusterSchema = {
  kind: string
  fields: Record<string, ClusterFieldSchema>
}

const base = '/config/api/cluster'
const query = (values: Record<string, string | number>) =>
  '?' +
  new URLSearchParams(Object.entries(values).map(([key, value]) => [key, String(value)])).toString()

export const clusterLifecycleApi = {
  schema: async () => {
    const definitions = parseArrayResponse<ClusterSchema>(await request<unknown>(`${base}/schema`))
    const cluster = definitions.find((item) => item.kind === 'AdminCluster')
    if (!cluster) throw new Error('AdminCluster schema is not available')
    return cluster
  },
  listDrafts: async () =>
    parseArrayResponse<ClusterDraft>(await request<unknown>(`${base}/draft/list`)),
  detail: (id: string) => request<ClusterDraft>(`${base}/draft/detail${query({ id })}`),
  create: (object: AdminClusterObject) =>
    request<ClusterDraft>(`${base}/draft`, { method: 'POST', body: JSON.stringify(object) }),
  fromPublished: (sourceId: number) =>
    request<ClusterDraft>(`${base}/draft${query({ sourceId })}`, { method: 'POST' }),
  update: (draft: ClusterDraft, object: AdminClusterObject) =>
    request<ClusterDraft>(`${base}/draft${query({ id: draft.id })}`, {
      method: 'PUT',
      body: JSON.stringify({ object, expectedRevision: draft.revision }),
    }),
  removeDraft: (draft: ClusterDraft) =>
    request<unknown>(`${base}/draft${query({ id: draft.id, expectedRevision: draft.revision })}`, {
      method: 'DELETE',
    }),
  validate: (object: AdminClusterObject, draftId?: string) =>
    request<{ valid: boolean; object: AdminClusterObject }>(
      `${base}/validate${draftId ? query({ id: draftId }) : ''}`,
      {
        method: 'POST',
        body: JSON.stringify(object),
      },
    ),
  preview: (object: AdminClusterObject, draftId?: string) =>
    request<{ object: AdminClusterObject; yaml: string }>(
      `${base}/preview${draftId ? query({ id: draftId }) : ''}`,
      {
        method: 'POST',
        body: JSON.stringify(object),
      },
    ),
  status: (id: string) => request<ClusterStatus>(`${base}/status${query({ id })}`),
  diff: (id: string) => request<ClusterDiff>(`${base}/diff${query({ id })}`),
  publish: (draft: ClusterDraft) =>
    request<ClusterPublishResult>(
      `${base}/publish${query({ id: draft.id, expectedRevision: draft.revision })}`,
      {
        method: 'PUT',
      },
    ),
  history: async (id: string) =>
    parseArrayResponse<ClusterHistory>(await request<unknown>(`${base}/history${query({ id })}`)),
  rollback: (draft: ClusterDraft, revision: number) =>
    request<ClusterPublishResult>(
      `${base}/rollback${query({ id: draft.id, revision, expectedRevision: draft.revision })}`,
      { method: 'PUT' },
    ),
}
