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

import type { RouteBindingHistoryEntry } from '../../types/api'
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
  const language = isEnglish ? 'en-US' : 'zh-CN'
  const loadingLabel = isEnglish ? 'Restoring…' : '恢复中…'
  const actionLabel = isEnglish
    ? action === 'restore'
      ? 'Restore this revision'
      : 'Roll back to this revision'
    : action === 'restore'
      ? '恢复此版本'
      : '回滚到此版本'

  return (
    <div className="route-history-list">
      {entries
        .map((entry, index) => ({ entry, previous: entries[index - 1] }))
        .reverse()
        .map(({ entry, previous }) => {
          const object = entry.object
          const spec = object?.spec
          const invocation = spec?.invocation
          const service = spec?.serviceSelection
          const entryConfig = spec?.entry
          const backend = spec?.backendRef
          const parameters = object ? object.spec.invocation.parameters || [] : []
          const details = routeHistoryChanges(entry, previous, isEnglish)
          const deletedLabel = isEnglish ? 'Route deleted' : '路由已删除'
          const firstOrLatest = !previous
            ? isEnglish
              ? 'First publication'
              : '首次发布'
            : isEnglish
              ? 'Published snapshot'
              : '发布快照'
          const snapshotLabel = !object
            ? deletedLabel
            : !previous
              ? firstOrLatest
              : !previous.object
                ? isEnglish
                  ? 'Published after deletion'
                  : '删除后重新发布'
                : isEnglish
                  ? 'Published snapshot'
                  : '发布快照'
          const showAction = Boolean(object || includeDeletedAction)
          const disabled = !canRollback || actionDisabled || busyRevision !== null || !showAction

          return (
            <article className="route-history-card" key={entry.revision}>
              <header className="route-history-head">
                <div className="route-history-heading">
                  <code>rev {entry.revision}</code>
                  <span
                    className={!object ? 'route-history-kind is-deleted' : 'route-history-kind'}
                  >
                    {snapshotLabel}
                  </span>
                </div>
                <time className="route-history-time" dateTime={entry.publishedAt}>
                  {new Date(entry.publishedAt).toLocaleString(language)}
                </time>
              </header>

              {object && spec && invocation && service && entryConfig && backend ? (
                <>
                  <div className="route-history-request">
                    <span>{entryConfig.method || 'HTTP'}</span>
                    <code>{entryConfig.path || '—'}</code>
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
                        {isEnglish ? 'Group' : '分组'} {service.group || '—'} ·{' '}
                        {isEnglish ? 'Version' : '版本'} {service.version || '—'}
                      </span>
                    </div>
                    <div className="route-history-detail">
                      <small>{isEnglish ? 'ROUTE SETTINGS' : '路由设置'}</small>
                      <strong>
                        {spec.enabled === false
                          ? isEnglish
                            ? 'Disabled'
                            : '已停用'
                          : isEnglish
                            ? 'Enabled'
                            : '已启用'}
                      </strong>
                      <span>
                        {isEnglish ? 'Timeout' : '超时'} {invocation.timeout || '—'} ·{' '}
                        {backend.type || '—'} {backend.clusterName || ''}
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
              ) : (
                <div className="route-history-deleted">{deletedLabel}</div>
              )}

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
                {object && (
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
        })}
    </div>
  )
}
