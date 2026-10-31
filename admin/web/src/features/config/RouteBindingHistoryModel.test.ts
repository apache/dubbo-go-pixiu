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

import { describe, expect, it } from 'vitest'
import type { AdminRouteBindingObject, RouteBindingHistoryEntry } from '../../types/api'
import { routeHistoryChanges } from './RouteBindingHistoryModel'

function historyEntry(revision: string, path: string, enabled = true): RouteBindingHistoryEntry {
  const object: AdminRouteBindingObject = {
    kind: 'AdminRouteBinding',
    metadata: { name: 'sample-greeter' },
    spec: {
      entry: { protocol: 'http', path, method: 'GET' },
      backendRef: { type: 'registry', clusterName: 'backend' },
      serviceSelection: {
        application: 'runtime-e2e',
        interface: 'pixiu.e2e.Greeter',
        group: '',
        version: '',
      },
      invocation: {
        method: 'SayHello',
        timeout: '1s',
        parameters: [{ from: 'queryStrings.name', to: 0, type: 'java.lang.String' }],
      },
      enabled,
    },
  }
  return { revision, publishedAt: '2026-10-31T07:00:00.000Z', object, resourceId: 1, methodId: 1 }
}

describe('route binding history model', () => {
  it('labels the first published snapshot', () => {
    expect(routeHistoryChanges(historyEntry('1', '/sample'), undefined, false)).toEqual([
      { field: '版本', before: '—', after: '首次发布' },
    ])
  })

  it('shows only the route fields changed since the previous snapshot', () => {
    const previous = historyEntry('1', '/sample')
    const current = historyEntry('2', '/sample-v2', false)

    expect(routeHistoryChanges(current, previous, false)).toEqual([
      { field: '入口路径', before: '/sample', after: '/sample-v2' },
      { field: '路由状态', before: '已启用', after: '已停用' },
    ])
  })

  it('reports changes to the HTTP-to-Dubbo argument mapping', () => {
    const previous = historyEntry('1', '/sample')
    const current = historyEntry('2', '/sample')
    current.object!.spec.invocation.parameters = [
      { from: 'queryStrings.id', to: 0, type: 'java.lang.Long' },
    ]

    expect(routeHistoryChanges(current, previous, false)).toEqual([
      {
        field: '参数映射',
        before: 'queryStrings.name → arg[0] · java.lang.String',
        after: 'queryStrings.id → arg[0] · java.lang.Long',
      },
    ])
  })

  it('describes a deleted snapshot and a later published snapshot', () => {
    const deleted: RouteBindingHistoryEntry = {
      revision: '2',
      publishedAt: '2026-10-31T07:01:00.000Z',
    }
    const restored = historyEntry('3', '/sample')

    expect(routeHistoryChanges(deleted, historyEntry('1', '/sample'), false)).toEqual([
      { field: '路由状态', before: '已发布', after: '已删除' },
    ])
    expect(routeHistoryChanges(restored, deleted, false)).toEqual([
      { field: '路由状态', before: '已删除', after: '重新发布' },
    ])
  })
})
