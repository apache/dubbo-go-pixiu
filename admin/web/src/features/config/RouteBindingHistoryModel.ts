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

import type { AdminRouteBindingObject, RouteBindingHistoryEntry } from '../../types/api'

export type RouteHistoryChange = {
  field: string
  before: string
  after: string
}

export function routeHistoryParameterSummary(object: AdminRouteBindingObject) {
  const parameters = object.spec.invocation.parameters || []
  if (parameters.length === 0) return '—'
  return parameters
    .map((parameter) => `${parameter.from} → arg[${parameter.to}] · ${parameter.type}`)
    .join('； ')
}

export function routeHistoryChanges(
  current: RouteBindingHistoryEntry,
  previous: RouteBindingHistoryEntry | undefined,
  isEnglish: boolean,
): RouteHistoryChange[] {
  const revisionLabel = isEnglish ? 'Revision' : '版本'
  const routeStateLabel = isEnglish ? 'Route state' : '路由状态'
  const deletedLabel = isEnglish ? 'Deleted' : '已删除'
  const publishedLabel = isEnglish ? 'Published' : '已发布'

  if (!previous) {
    return [
      {
        field: revisionLabel,
        before: '—',
        after: current.object ? (isEnglish ? 'First publication' : '首次发布') : deletedLabel,
      },
    ]
  }

  if (!current.object) {
    return [
      {
        field: routeStateLabel,
        before: previous.object ? publishedLabel : deletedLabel,
        after: deletedLabel,
      },
    ]
  }

  if (!previous.object) {
    return [
      {
        field: routeStateLabel,
        before: deletedLabel,
        after: isEnglish ? 'Published again' : '重新发布',
      },
    ]
  }

  const before = historyValues(previous.object, isEnglish)
  const after = historyValues(current.object, isEnglish)
  return Object.keys(after)
    .filter((key) => before[key] !== after[key])
    .map((key) => ({
      field: key,
      before: before[key],
      after: after[key],
    }))
}

function historyValues(
  object: AdminRouteBindingObject,
  isEnglish: boolean,
): Record<string, string> {
  const spec = object.spec
  const text = (value: string | undefined) => value?.trim() || '—'
  const yes = isEnglish ? 'Enabled' : '已启用'
  const no = isEnglish ? 'Disabled' : '已停用'
  const cluster = [spec.backendRef.type, spec.backendRef.clusterName]
    .filter((value) => value?.trim())
    .map(text)
    .join(' · ')
  const groupLabel = isEnglish ? 'Group' : '分组'
  const versionLabel = isEnglish ? 'Version' : '版本'
  const parameters = routeHistoryParameterSummary(object)

  return {
    [isEnglish ? 'HTTP method' : 'HTTP 方法']: text(spec.entry.method).toUpperCase(),
    [isEnglish ? 'Route path' : '入口路径']: text(spec.entry.path),
    [isEnglish ? 'Backend cluster' : '后端集群']: cluster,
    [isEnglish ? 'Application' : '应用']: text(spec.serviceSelection.application),
    [isEnglish ? 'Service interface' : '服务接口']: text(spec.serviceSelection.interface),
    [isEnglish ? 'Dubbo method' : 'Dubbo 方法']: text(spec.invocation.method),
    [groupLabel]: text(spec.serviceSelection.group),
    [versionLabel]: text(spec.serviceSelection.version),
    [isEnglish ? 'Timeout' : '超时']: text(spec.invocation.timeout),
    [isEnglish ? 'Argument mapping' : '参数映射']: parameters,
    [isEnglish ? 'Route status' : '路由状态']: spec.enabled === false ? no : yes,
  }
}
