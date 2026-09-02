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

import { parseJsonData, request } from './client'
import {
  AdminObject,
  Cluster,
  Listener,
  LoginResult,
  Method,
  ObjectSchema,
  Resource,
  RouteBindingDetail,
  RouteBindingPreview,
  RouteBindingRecord,
  RouteBindingSummary,
  UserInfo,
} from '../types'

// ---- Account ----

export function login(username: string, password: string) {
  return request<LoginResult>('/login', { method: 'POST', form: { username, password }, auth: false })
}

export function register(username: string, password: string) {
  return request<string>('/register', { method: 'POST', form: { username, password }, auth: false })
}

export function logout() {
  return request<string>('/user/logout', { method: 'POST' })
}

export function getUserInfo() {
  return request<UserInfo>('/user/getInfo', { method: 'POST' })
}

export function editPassword(oldPassword: string, newPassword: string) {
  return request<string>('/user/password/edit', { method: 'POST', form: { oldPassword, newPassword } })
}

// ---- Base info & legacy resources / methods ----

export function getBaseInfo() {
  return request<string>('/config/api/base')
}

export function setBaseInfo(content: string) {
  // The backend registers POST/PUT only on the trailing-slash path.
  return request<string>('/config/api/base/', { method: 'POST', form: { content } })
}

export async function getResourceList(): Promise<Resource[]> {
  const data = await request<unknown>('/config/api/resource/list')
  return parseJsonData<Resource[]>(data, [])
}

export function getResourceDetail(resourceId: number | string) {
  return request<string>('/config/api/resource/detail', { query: { resourceId } })
}

export function createResource(content: string) {
  return request<string>('/config/api/resource', { method: 'POST', form: { content } })
}

export function modifyResource(resourceId: number | string, content: string) {
  return request<string>('/config/api/resource', {
    method: 'PUT',
    query: { resourceId },
    form: { content },
  })
}

export function deleteResource(resourceId: number | string) {
  return request<string>('/config/api/resource', { method: 'DELETE', query: { resourceId } })
}

export async function getMethodList(resourceId: number | string): Promise<Method[]> {
  const data = await request<unknown>('/config/api/resource/method/list', { query: { resourceId } })
  return parseJsonData<Method[]>(data, [])
}

export function getMethodDetail(resourceId: number | string, methodId: number | string) {
  return request<string>('/config/api/resource/method/detail', { query: { resourceId, methodId } })
}

export function createMethod(resourceId: number | string, content: string) {
  return request<string>('/config/api/resource/method', {
    method: 'POST',
    query: { resourceId },
    form: { content },
  })
}

export function modifyMethod(resourceId: number | string, methodId: number | string, content: string) {
  return request<string>('/config/api/resource/method', {
    method: 'PUT',
    query: { resourceId, methodId },
    form: { content },
  })
}

export function deleteMethod(resourceId: number | string, methodId: number | string) {
  return request<string>('/config/api/resource/method', {
    method: 'DELETE',
    query: { resourceId, methodId },
  })
}

// ---- Route binding (schema-driven POC) ----

export function getRouteBindingSchema() {
  return request<ObjectSchema[]>('/config/api/route-binding/schema')
}

export function listRouteBindings() {
  return request<RouteBindingSummary[]>('/config/api/route-binding/list')
}

export function getRouteBindingDetail(name: string) {
  return request<RouteBindingDetail>('/config/api/route-binding/detail', { query: { name } })
}

export function saveRouteBindingDraft(object: AdminObject) {
  return request<RouteBindingRecord>('/config/api/route-binding', {
    method: 'POST',
    rawBody: JSON.stringify(object),
  })
}

export function previewRouteBinding(object?: AdminObject, name?: string) {
  return request<RouteBindingPreview>('/config/api/route-binding/preview', {
    method: 'POST',
    query: { name },
    rawBody: object ? JSON.stringify(object) : undefined,
  })
}

export function publishRouteBinding(name: string) {
  return request<RouteBindingRecord>('/config/api/route-binding/publish', {
    method: 'PUT',
    query: { name },
  })
}

export function getRouteBindingHistory(name: string) {
  return request<RouteBindingRecord[]>('/config/api/route-binding/history', { query: { name } })
}

export function rollbackRouteBinding(name: string, revision?: number) {
  return request<RouteBindingRecord>('/config/api/route-binding/rollback', {
    method: 'POST',
    query: { name, revision: revision && revision > 0 ? revision : undefined },
  })
}

export function deleteRouteBindingDraft(name: string) {
  return request<string>('/config/api/route-binding', { method: 'DELETE', query: { name } })
}

// ---- Plugin groups (endpoints may be absent on some backends; page surfaces errors) ----

export function getPluginGroupList() {
  return request<unknown>('/config/api/plugin_group/list')
}

export function getPluginGroupDetail(name: string) {
  return request<string>('/config/api/plugin_group/detail', { query: { name } })
}

export function createPluginGroup(content: string) {
  return request<string>('/config/api/plugin_group', { method: 'POST', form: { content } })
}

export function updatePluginGroup(content: string) {
  return request<string>('/config/api/plugin_group', { method: 'PUT', form: { content } })
}

export function deletePluginGroup(name: string) {
  return request<string>('/config/api/plugin_group', { method: 'DELETE', query: { name } })
}

// ---- Clusters ----

export function getClusterList() {
  return request<Cluster[]>('/config/api/cluster/list')
}

export function getClusterDetail(id: number | string) {
  return request<string>('/config/api/cluster/detail', { query: { clusterId: id } })
}

export function createCluster(content: string) {
  return request<string>('/config/api/cluster', { method: 'PUT', form: { content } })
}

export function updateCluster(content: string) {
  return request<string>('/config/api/cluster', { method: 'POST', form: { content } })
}

export function deleteCluster(id: number | string) {
  return request<string>('/config/api/cluster', { method: 'DELETE', query: { clusterId: id } })
}

// ---- Listeners ----

export function getListenerList() {
  return request<Listener[]>('/config/api/listener/list')
}

export function getListenerDetail(name: string) {
  return request<string>('/config/api/listener/detail', { query: { listener: name } })
}

export function createListener(content: string) {
  return request<string>('/config/api/listener', { method: 'PUT', form: { content } })
}

export function updateListener(content: string) {
  return request<string>('/config/api/listener', { method: 'POST', form: { content } })
}

export function deleteListener(name: string) {
  return request<string>('/config/api/listener', { method: 'DELETE', query: { listener: name } })
}

// ---- Rate limiter (single shared filter config) ----

export function getRateLimiter() {
  return request<string>('/config/api/plugin/ratelimit')
}

export function createRateLimiter(content: string) {
  return request<string>('/config/api/plugin/ratelimit/', { method: 'POST', form: { content } })
}

export function updateRateLimiter(content: string) {
  return request<string>('/config/api/plugin/ratelimit/', { method: 'PUT', form: { content } })
}

export function deleteRateLimiter() {
  return request<string>('/config/api/plugin/ratelimit/', { method: 'DELETE' })
}

// ---- OPA ----

export interface OpaQuery {
  server_url?: string
  policy_id?: string
  bearer_token?: string
}

export function getOpaPolicy(query: OpaQuery) {
  return request<string>('/config/api/opa/policy', { query: { ...query } })
}

export function putOpaPolicy(query: OpaQuery, content: string) {
  return request<string>('/config/api/opa/policy', {
    method: 'PUT',
    form: { ...query, content },
  })
}

export function deleteOpaPolicy(query: OpaQuery) {
  return request<string>('/config/api/opa/policy', { method: 'DELETE', query: { ...query } })
}
