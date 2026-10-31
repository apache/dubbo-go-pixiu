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
import { routeHistoryChanges } from './RouteBindingHistoryModel'

type RouteBindingHistoryListProps = {
  readonly isEnglish: boolean
  readonly entries: RouteBindingHistoryEntry[]
  readonly canRollback: boolean
  readonly action: 'rollback' | 'restore'
  readonly actionDisabled?: boolean
  readonly busyRevision?: string | null
  readonly includeDeletedAction?: boolean
  readonly onAction: (revision: string) => void
}

export function RouteBindingHistoryList({
  isEnglish,
  entries,
  canRollback,
  action,
  actionDisabled = false,
  busyRevision = null,
  includeDeletedAction = false,
  onAction,
}: RouteBindingHistoryListProps) {
  const revisions = entries
    .map((entry, index) => ({ entry, previous: entries[index - 1] }))
    .reverse()

  return (
    <div className="route-history-list">
      {revisions.map(({ entry, previous }) => (
        <RouteBindingHistoryCard
          key={entry.revision}
          isEnglish={isEnglish}
          entry={entry}
          previous={previous}
          canRollback={canRollback}
          action={action}
          actionDisabled={actionDisabled}
          busyRevision={busyRevision}
          includeDeletedAction={includeDeletedAction}
          onAction={onAction}
        />
      ))}
    </div>
  )
}

type RouteBindingHistoryCardProps = Omit<RouteBindingHistoryListProps, 'entries'> & {
  readonly entry: RouteBindingHistoryEntry
  readonly previous?: RouteBindingHistoryEntry
}

function RouteBindingHistoryCard({
  isEnglish,
  entry,
  previous,
  canRollback,
  action,
  actionDisabled = false,
  busyRevision = null,
  includeDeletedAction = false,
  onAction,
}: RouteBindingHistoryCardProps) {
  const language = isEnglish ? 'en-US' : 'zh-CN'
  const loadingLabel = isEnglish ? 'Restoring…' : '恢复中…'
  const actionLabel = historyActionLabel(action, isEnglish)
  const snapshotLabel = historySnapshotLabel(entry, previous, isEnglish)
  const details = routeHistoryChanges(entry, previous, isEnglish)
  const showAction = Boolean(entry.object || includeDeletedAction)
  const disabled = !canRollback || actionDisabled || busyRevision !== null || !showAction

  return (
    <article className="route-history-card">
      <header className="route-history-head">
        <div className="route-history-heading">
          <code>rev {entry.revision}</code>
          <span className={entry.object ? 'route-history-kind' : 'route-history-kind is-deleted'}>
            {snapshotLabel}
          </span>
        </div>
        <time className="route-history-time" dateTime={entry.publishedAt}>
          {new Date(entry.publishedAt).toLocaleString(language)}
        </time>
      </header>

      <RouteBindingHistorySnapshot object={entry.object} isEnglish={isEnglish} />

      <div className="route-history-changes">
        <small>{isEnglish ? 'CHANGES FROM PREVIOUS REVISION' : '与上一版本相比'}</small>
        {details.length > 0 ? (
          <ul>
            {details.map((change) => (
              <li key={change.field}>
                <strong>{change.field}</strong>
                <code>{change.before}</code>
                <span aria-hidden="true">→</span>
                <code>{change.after}</code>
              </li>
            ))}
          </ul>
        ) : (
          <span>
            {isEnglish
              ? `Same route configuration as revision ${previous?.revision ?? '—'}`
              : `路由配置与 rev ${previous?.revision ?? '—'} 一致`}
          </span>
        )}
      </div>

      <footer className="route-history-actions">
        {entry.object && (
          <span className="route-history-record-ids">
            {isEnglish ? 'Resource' : '资源'} {entry.resourceId ?? '—'} ·{' '}
            {isEnglish ? 'Method' : '方法'} {entry.methodId ?? '—'}
          </span>
        )}
        {showAction && (
          <button
            className="secondary"
            type="button"
            disabled={disabled}
            onClick={() => onAction(entry.revision)}
          >
            {busyRevision === entry.revision ? loadingLabel : actionLabel}
          </button>
        )}
      </footer>
    </article>
  )
}

function historyActionLabel(action: 'rollback' | 'restore', isEnglish: boolean) {
  if (action === 'restore') return isEnglish ? 'Restore this revision' : '恢复此版本'
  return isEnglish ? 'Roll back to this revision' : '回滚到此版本'
}

function historySnapshotLabel(
  entry: RouteBindingHistoryEntry,
  previous: RouteBindingHistoryEntry | undefined,
  isEnglish: boolean,
) {
  if (!entry.object) return isEnglish ? 'Route deleted' : '路由已删除'
  if (!previous) return isEnglish ? 'First publication' : '首次发布'
  if (!previous.object) return isEnglish ? 'Published after deletion' : '删除后重新发布'
  return isEnglish ? 'Published snapshot' : '发布快照'
}

function RouteBindingHistorySnapshot({
  object,
  isEnglish,
}: {
  readonly object: AdminRouteBindingObject | undefined
  readonly isEnglish: boolean
}) {
  if (!object) {
    return <div className="route-history-deleted">{isEnglish ? 'Route deleted' : '路由已删除'}</div>
  }

  const { spec } = object
  const { invocation, serviceSelection: service, entry, backendRef: backend } = spec
  const parameters = invocation.parameters || []

  return (
    <>
      <div className="route-history-request">
        <span>{entry.method || 'HTTP'}</span>
        <code>{entry.path || '—'}</code>
      </div>

      <div className="route-history-details">
        <div className="route-history-detail">
          <small>{isEnglish ? 'DUBBO SERVICE' : 'DUBBO 服务'}</small>
          <strong>{service.application || '—'}</strong>
          <code>{service.interface || '—'}</code>
        </div>
        <div className="route-history-detail">
          <small>{isEnglish ? 'INVOCATION' : '调用'}</small>
          <strong>{invocation.method || '—'}</strong>
          <span>
            {isEnglish ? 'Group' : '分组'} {service.group || '—'} · {isEnglish ? 'Version' : '版本'}{' '}
            {service.version || '—'}
          </span>
        </div>
        <div className="route-history-detail">
          <small>{isEnglish ? 'ROUTE SETTINGS' : '路由设置'}</small>
          <strong>{routeStatusLabel(spec.enabled, isEnglish)}</strong>
          <span>
            {isEnglish ? 'Timeout' : '超时'} {invocation.timeout || '—'} · {backend.type || '—'}{' '}
            {backend.clusterName || ''}
          </span>
        </div>
      </div>

      <div className="route-history-parameters">
        <small>{isEnglish ? 'HTTP → DUBBO ARGUMENTS' : 'HTTP → Dubbo 参数映射'}</small>
        {parameters.length > 0 ? (
          <ul>
            {parameters.map((parameter, index) => (
              <li key={`${parameter.from}-${parameter.to}-${index}`}>
                <code>
                  {parameter.from} → arg[{parameter.to}] · {parameter.type}
                </code>
              </li>
            ))}
          </ul>
        ) : (
          <span>{isEnglish ? 'No argument mappings' : '无参数映射'}</span>
        )}
      </div>
    </>
  )
}

function routeStatusLabel(enabled: boolean | undefined, isEnglish: boolean) {
  if (enabled === false) return isEnglish ? 'Disabled' : '已停用'
  return isEnglish ? 'Enabled' : '已启用'
}
