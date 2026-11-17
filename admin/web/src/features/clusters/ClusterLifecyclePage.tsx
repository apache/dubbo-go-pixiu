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

import { useCallback, useEffect, useMemo, useState } from 'react'
import { clusterApi } from '../../services/cluster-api'
import { clusterLifecycleApi } from '../../services/cluster-lifecycle-api'
import type {
  AdminClusterObject,
  ClusterDiff,
  ClusterDraft,
  ClusterHistory,
  ClusterFieldSchema,
  ClusterSchema,
  ClusterStatus,
} from '../../services/cluster-lifecycle-api'
import { ApiError } from '../../services/http'
import type { JsonObject } from '../../types/api'
import type { Locale } from '../../i18n'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import {
  clusterFieldValue,
  emptyCluster,
  parseClusterYaml,
  stringifyClusterYaml,
  updateClusterField,
} from './ClusterEditorModel'

type EditorState = { draft: ClusterDraft | null; object: AdminClusterObject; dirty: boolean }
type ConfirmAction = { kind: 'back' | 'draft' | 'runtime' | 'rollback'; revision?: number }

function displayError(error: unknown) {
  return error instanceof Error ? error.message : String(error)
}

function clusterID(value: JsonObject) {
  const id = Number(value.id)
  return Number.isInteger(id) && id > 0 ? id : 0
}

export function ClusterLifecyclePage({
  locale,
  onDirtyChange,
}: {
  locale: Locale
  onDirtyChange?: (dirty: boolean) => void
}) {
  const en = locale === 'en-US'
  const label = (zh: string, english: string) => (en ? english : zh)
  const [runtime, setRuntime] = useState<JsonObject[]>([])
  const [drafts, setDrafts] = useState<ClusterDraft[]>([])
  const [schema, setSchema] = useState<ClusterSchema | null>(null)
  const [loading, setLoading] = useState(true)
  const [pageError, setPageError] = useState('')
  const [editor, setEditor] = useState<EditorState | null>(null)
  const [status, setStatus] = useState<ClusterStatus | null>(null)
  const [diff, setDiff] = useState<ClusterDiff | null>(null)
  const [history, setHistory] = useState<ClusterHistory[]>([])
  const [mode, setMode] = useState<'form' | 'yaml'>('form')
  const [yamlText, setYamlText] = useState('')
  const [yamlError, setYamlError] = useState('')
  const [issues, setIssues] = useState<Array<{ path: string; message: string }>>([])
  const [preview, setPreview] = useState('')
  const [notice, setNotice] = useState('')
  const [noticeError, setNoticeError] = useState(false)
  const [busy, setBusy] = useState('')
  const [confirm, setConfirm] = useState<ConfirmAction | null>(null)
  const hasUnsavedChanges = Boolean(editor?.dirty || yamlError)
  useEffect(() => {
    onDirtyChange?.(hasUnsavedChanges)
  }, [hasUnsavedChanges, onDirtyChange])
  useEffect(() => {
    if (!hasUnsavedChanges) return
    const onBeforeUnload = (event: BeforeUnloadEvent) => {
      event.preventDefault()
      event.returnValue = ''
    }
    window.addEventListener('beforeunload', onBeforeUnload)
    return () => window.removeEventListener('beforeunload', onBeforeUnload)
  }, [hasUnsavedChanges])

  const load = useCallback(async () => {
    setLoading(true)
    setPageError('')
    try {
      const [published, saved, definition] = await Promise.all([
        clusterApi.list(),
        clusterLifecycleApi.listDrafts(),
        clusterLifecycleApi.schema(),
      ])
      setRuntime(published)
      setDrafts(saved)
      setSchema(definition)
    } catch (error) {
      setPageError(displayError(error))
    } finally {
      setLoading(false)
    }
  }, [])
  useEffect(() => {
    void load()
  }, [load])

  const publishedIDs = useMemo(() => new Set(runtime.map(clusterID)), [runtime])
  const unlinkedRuntime = useMemo(() => {
    const linked = new Set(drafts.map((draft) => draft.sourceId).filter(Boolean))
    return runtime.filter((item) => !linked.has(clusterID(item)))
  }, [drafts, runtime])

  const openEditor = useCallback((draft: ClusterDraft | null, object: AdminClusterObject) => {
    setEditor({ draft, object, dirty: false })
    setMode('form')
    setYamlText(stringifyClusterYaml(object))
    setYamlError('')
    setIssues([])
    setPreview('')
    setNotice('')
    setNoticeError(false)
    setStatus(null)
    setDiff(null)
    setHistory([])
  }, [])

  const refreshEditor = useCallback(async (draft: ClusterDraft) => {
    const [fresh, publication, changes, revisions] = await Promise.all([
      clusterLifecycleApi.detail(draft.id),
      clusterLifecycleApi.status(draft.id),
      clusterLifecycleApi.diff(draft.id),
      clusterLifecycleApi.history(draft.id),
    ])
    setEditor({ draft: fresh, object: fresh.object, dirty: false })
    setYamlText(stringifyClusterYaml(fresh.object))
    setStatus(publication)
    setDiff(changes)
    setHistory(revisions)
    setIssues([])
    setYamlError('')
    setPreview('')
  }, [])

  const run = async (action: string, work: () => Promise<void>) => {
    if (busy) return
    setBusy(action)
    setNotice('')
    setNoticeError(false)
    try {
      await work()
    } catch (error) {
      if (error instanceof ApiError) setIssues(error.issues)
      setNotice(displayError(error))
      setNoticeError(true)
      window.setTimeout(
        () => document.querySelector<HTMLElement>('.cluster-editor [aria-invalid="true"]')?.focus(),
        0,
      )
    } finally {
      setBusy('')
    }
  }

  const openPublished = (item: JsonObject) =>
    void run('open', async () => {
      const id = clusterID(item)
      if (!id) throw new Error(label('集群缺少有效 ID', 'Cluster has no valid ID'))
      const draft = await clusterLifecycleApi.fromPublished(id)
      openEditor(draft, draft.object)
      await refreshEditor(draft)
      await load()
    })

  const save = () =>
    void run('save', async () => {
      if (!editor || yamlError)
        throw new Error(yamlError || label('没有可保存的配置', 'No configuration to save'))
      const saved = editor.draft
        ? await clusterLifecycleApi.update(editor.draft, editor.object)
        : await clusterLifecycleApi.create(editor.object)
      await refreshEditor(saved)
      await load()
      setNotice(label('草稿已保存', 'Draft saved'))
    })

  const validate = () =>
    void run('validate', async () => {
      if (!editor || yamlError)
        throw new Error(yamlError || label('请先检查 YAML', 'Check YAML first'))
      setIssues([])
      await clusterLifecycleApi.validate(editor.object, editor.draft?.id)
      setNotice(label('校验通过', 'Validation passed'))
    })

  const showPreview = () =>
    void run('preview', async () => {
      if (!editor || yamlError)
        throw new Error(yamlError || label('请先检查 YAML', 'Check YAML first'))
      setIssues([])
      const result = await clusterLifecycleApi.preview(editor.object, editor.draft?.id)
      setPreview(result.yaml)
      setNotice(label('已生成运行时 YAML 预览', 'Runtime YAML preview ready'))
    })

  const publish = () =>
    void run('publish', async () => {
      if (!editor?.draft || editor.dirty || yamlError) {
        throw new Error(label('请先保存有效草稿，再发布', 'Save a valid draft before publishing'))
      }
      setIssues([])
      await clusterLifecycleApi.publish(editor.draft)
      await refreshEditor(editor.draft)
      await load()
      setNotice(
        label('集群已发布，正在通过 xDS 生效', 'Cluster published; xDS propagation is in progress'),
      )
    })

  const confirmed = () => {
    const action = confirm
    setConfirm(null)
    if (!action) return
    if (action.kind === 'back') {
      setEditor(null)
      return
    }
    void run(action.kind, async () => {
      if (!editor?.draft) return
      if (action.kind === 'draft') {
        await clusterLifecycleApi.removeDraft(editor.draft)
        setEditor(null)
        await load()
      } else if (action.kind === 'runtime') {
        if (!editor.draft.sourceId) return
        await clusterApi.remove(String(editor.draft.sourceId))
        await refreshEditor(editor.draft)
        await load()
        setNotice(label('已发布集群已删除，草稿仍保留', 'Published cluster deleted; draft remains'))
      } else if (action.kind === 'rollback' && action.revision) {
        await clusterLifecycleApi.rollback(editor.draft, action.revision)
        await refreshEditor(editor.draft)
        await load()
        setNotice(label('历史版本已恢复并发布', 'Historical revision restored and published'))
      }
    })
  }

  const back = () => {
    if (editor?.dirty || yamlError) setConfirm({ kind: 'back' })
    else setEditor(null)
  }

  const changeField = (path: string, value: string) => {
    if (!editor) return
    const object = updateClusterField(editor.object, path, value)
    setEditor({ ...editor, object, dirty: true })
    setYamlText(stringifyClusterYaml(object))
    setIssues(issues.filter((issue) => issue.path !== path))
    setPreview('')
  }

  const changeYaml = (value: string) => {
    setYamlText(value)
    setPreview('')
    if (!editor) return
    try {
      const object = parseClusterYaml(value, editor.draft?.sourceId)
      setEditor({ ...editor, object, dirty: true })
      setYamlError('')
      setIssues([])
    } catch (error) {
      setYamlError(displayError(error))
      setEditor({ ...editor, dirty: true })
    }
  }

  const field = (
    path: string,
    title: string,
    hint: string,
    kind: 'text' | 'number' | 'select' = 'text',
  ) => {
    if (!editor) return null
    const definition: ClusterFieldSchema | undefined =
      path === 'spec.type'
        ? schema?.fields.type
        : path === 'spec.endpoint.address'
          ? schema?.fields.endpoint?.properties?.address
          : path === 'spec.endpoint.port'
            ? schema?.fields.endpoint?.properties?.port
            : undefined
    const error = issues.find((issue) => issue.path === path)?.message
    const id = `cluster-${path.replace(/\./g, '-')}`
    const value = clusterFieldValue(editor.object, path)
    const choices = definition?.enum ?? []
    return (
      <div className="cluster-field" key={path}>
        <label htmlFor={id}>
          {title}
          {definition?.required ? ' *' : ''}
        </label>
        {kind === 'select' ? (
          <select
            id={id}
            value={value}
            onChange={(event) => changeField(path, event.target.value)}
            aria-invalid={Boolean(error)}
            aria-describedby={error ? `${id}-error` : `${id}-hint`}
          >
            <option value="">{label('请选择', 'Select a type')}</option>
            {value && !choices.includes(value) && <option value={value}>{value}</option>}
            {choices.map((option) => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </select>
        ) : (
          <input
            id={id}
            type={kind}
            value={value}
            placeholder={definition?.ui?.placeholder ?? hint}
            min={kind === 'number' ? definition?.minimum : undefined}
            max={kind === 'number' ? definition?.maximum : undefined}
            step={kind === 'number' ? 1 : undefined}
            onChange={(event) => changeField(path, event.target.value)}
            aria-invalid={Boolean(error)}
            aria-describedby={error ? `${id}-error` : `${id}-hint`}
          />
        )}
        {error ? (
          <small id={`${id}-error`} className="route-field-error">
            {error}
          </small>
        ) : (
          <small id={`${id}-hint`}>{hint}</small>
        )}
      </div>
    )
  }

  return (
    <div className="page-resource cluster-lifecycle">
      {editor ? (
        <div className="cluster-editor">
          <div className="route-editor-backbar">
            <button type="button" className="route-editor-back" onClick={back}>
              {label('← 返回集群', '← Back to clusters')}
            </button>
            <span className="cluster-editor-id">
              {editor.draft
                ? `${label('草稿', 'Draft')} ${editor.draft.id.slice(0, 8)}`
                : label('新草稿', 'New draft')}
            </span>
          </div>
          <div className="cluster-editor-heading">
            <div>
              <p className="eyebrow">CLUSTER CONFIGURATION</p>
              <h1>{editor.object.metadata.name || label('未命名集群', 'Unnamed cluster')}</h1>
              <p className="muted">
                {status?.publishedExists
                  ? label('已发布 · 可继续编辑草稿', 'Published · draft can be edited')
                  : label('草稿尚未发布', 'Draft is not published')}
                {editor.dirty ? ` · ${label('有未保存改动', 'Unsaved changes')}` : ''}
              </p>
            </div>
            <div className="cluster-actions">
              <button
                type="button"
                className="secondary"
                disabled={Boolean(busy)}
                onClick={validate}
              >
                {label('校验', 'Validate')}
              </button>
              <button
                type="button"
                className="secondary"
                disabled={Boolean(busy)}
                onClick={showPreview}
              >
                {label('预览', 'Preview')}
              </button>
              <button type="button" className="secondary" disabled={Boolean(busy)} onClick={save}>
                {busy === 'save' ? label('保存中…', 'Saving…') : label('保存草稿', 'Save draft')}
              </button>
              <button
                type="button"
                className="primary"
                disabled={Boolean(busy) || !editor.draft || editor.dirty || Boolean(yamlError)}
                onClick={publish}
              >
                {busy === 'publish' ? label('发布中…', 'Publishing…') : label('发布', 'Publish')}
              </button>
            </div>
          </div>
          {notice && (
            <div
              className={`cluster-notice ${noticeError ? 'error' : ''}`}
              role={noticeError ? 'alert' : 'status'}
            >
              {notice}
            </div>
          )}
          <div
            className="route-editor-tabs"
            role="tablist"
            aria-label={label('编辑模式', 'Edit mode')}
          >
            <button
              role="tab"
              aria-selected={mode === 'form'}
              className={mode === 'form' ? 'active' : ''}
              onClick={() => {
                if (!yamlError) setMode('form')
              }}
            >
              {label('结构化表单', 'Form')}
            </button>
            <button
              role="tab"
              aria-selected={mode === 'yaml'}
              className={mode === 'yaml' ? 'active' : ''}
              onClick={() => {
                if (mode === 'form') setYamlText(stringifyClusterYaml(editor.object))
                setMode('yaml')
              }}
            >
              {label('高级 YAML', 'Advanced YAML')}
            </button>
          </div>
          {mode === 'form' ? (
            <form
              className="panel cluster-form"
              noValidate
              onSubmit={(event) => {
                event.preventDefault()
                save()
              }}
            >
              {field(
                'metadata.name',
                label('集群名称', 'Cluster name'),
                label('例如 user-provider', 'For example user-provider'),
              )}
              {field(
                'spec.type',
                label('发现类型', 'Discovery type'),
                label('选择运行时支持的类型', 'Choose a supported discovery type'),
                'select',
              )}
              {field('spec.endpoint.address', label('端点地址', 'Endpoint address'), '127.0.0.1')}
              {field('spec.endpoint.port', label('端点端口', 'Endpoint port'), '20880', 'number')}
            </form>
          ) : (
            <div className="panel cluster-yaml-panel">
              <label htmlFor="cluster-yaml">
                {label('旧版 Cluster YAML', 'Legacy Cluster YAML')}
              </label>
              <p className="muted">
                {label(
                  '运行时 ID 由系统管理；可粘贴旧版 YAML。',
                  'Runtime ID is managed by Pixiu; you can paste legacy YAML.',
                )}
              </p>
              <textarea
                id="cluster-yaml"
                className="resize-none"
                spellCheck={false}
                value={yamlText}
                onChange={(event) => changeYaml(event.target.value)}
                aria-invalid={Boolean(yamlError)}
                aria-describedby={yamlError ? 'cluster-yaml-error' : undefined}
              />
              {yamlError && (
                <p id="cluster-yaml-error" className="route-field-error" role="alert">
                  {yamlError}
                </p>
              )}
            </div>
          )}
          <div className="cluster-inspector">
            <section className="panel">
              <h2>{label('运行时预览', 'Runtime preview')}</h2>
              {preview ? (
                <pre>{preview}</pre>
              ) : (
                <p className="muted">
                  {label(
                    '点击预览生成即将发布的 YAML。',
                    'Preview the YAML that will be published.',
                  )}
                </p>
              )}
            </section>
            <section className="panel">
              <h2>{label('草稿与已发布差异', 'Draft and published diff')}</h2>
              {editor.dirty ? (
                <p className="muted">
                  {label('保存草稿后刷新差异。', 'Save the draft to refresh the diff.')}
                </p>
              ) : diff?.changes.length ? (
                <ul className="cluster-diff">
                  {diff.changes.map((change) => (
                    <li key={change.path}>
                      <code>{change.path}</code>
                      <span>
                        {String(change.before ?? '∅')} → {String(change.after ?? '∅')}
                      </span>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="muted">{label('暂无差异', 'No changes')}</p>
              )}
            </section>
          </div>
          {editor.draft && (
            <section className="panel cluster-history">
              <div className="cluster-section-head">
                <h2>{label('发布历史', 'Published history')}</h2>
                <button
                  type="button"
                  className="secondary"
                  disabled={Boolean(busy)}
                  onClick={() =>
                    void run('refresh', async () => {
                      if (editor.draft) await refreshEditor(editor.draft)
                    })
                  }
                >
                  {label('刷新', 'Refresh')}
                </button>
              </div>
              {history.length ? (
                <ul>
                  {history.map((entry) => (
                    <li key={entry.revision}>
                      <span>
                        <strong>#{entry.revision}</strong> ·{' '}
                        {new Date(entry.publishedAt).toLocaleString(locale)} ·{' '}
                        {entry.object.metadata.name}
                      </span>
                      <button
                        type="button"
                        className="link-btn"
                        disabled={Boolean(busy) || editor.dirty}
                        onClick={() => setConfirm({ kind: 'rollback', revision: entry.revision })}
                      >
                        {label('恢复并发布', 'Restore and publish')}
                      </button>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="muted">{label('还没有发布历史。', 'No published revisions yet.')}</p>
              )}
            </section>
          )}
          {editor.draft && (
            <div className="cluster-danger-actions">
              <button
                type="button"
                className="link-btn"
                disabled={Boolean(busy)}
                onClick={() => setConfirm({ kind: 'draft' })}
              >
                {label('删除草稿', 'Delete draft')}
              </button>
              {status?.publishedExists && (
                <button
                  type="button"
                  className="link-btn"
                  disabled={Boolean(busy)}
                  onClick={() => setConfirm({ kind: 'runtime' })}
                >
                  {label('删除已发布集群', 'Delete published cluster')}
                </button>
              )}
            </div>
          )}
        </div>
      ) : (
        <>
          <div className="placeholder-head">
            <div>
              <p className="eyebrow">CLUSTER CONFIGURATION</p>
              <h1>{label('集群', 'Clusters')}</h1>
              <p className="muted">
                {label(
                  '编辑草稿，预览配置，再发布到网关。',
                  'Edit drafts, preview configuration, then publish to the gateway.',
                )}
              </p>
            </div>
            <button
              className="primary"
              onClick={() => openEditor(null, emptyCluster(schema?.fields.type?.default))}
              disabled={loading || !schema}
            >
              {label('新建草稿', 'New draft')}
            </button>
          </div>
          <div className="panel route-panel cluster-list">
            {loading ? (
              <div className="loading-skeleton" />
            ) : pageError ? (
              <div className="empty" role="alert">
                {pageError}
                <button className="secondary" onClick={() => void load()}>
                  {label('重试', 'Retry')}
                </button>
              </div>
            ) : !drafts.length && !unlinkedRuntime.length ? (
              <div className="empty">
                {label(
                  '暂无集群。创建草稿开始配置。',
                  'No clusters yet. Create a draft to get started.',
                )}
              </div>
            ) : (
              <div className="cluster-table-scroll">
                <table>
                  <thead>
                    <tr>
                      <th>{label('名称', 'Name')}</th>
                      <th>{label('类型', 'Type')}</th>
                      <th>{label('端点', 'Endpoint')}</th>
                      <th>{label('状态', 'Status')}</th>
                      <th>{label('操作', 'Actions')}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {drafts.map((draft) => (
                      <tr key={draft.id}>
                        <td>
                          <b>
                            {draft.object.metadata.name || label('未命名草稿', 'Unnamed draft')}
                          </b>
                        </td>
                        <td>{draft.object.spec.type || '—'}</td>
                        <td>
                          {draft.object.spec.endpoint?.address || '—'}:
                          {draft.object.spec.endpoint?.port || '—'}
                        </td>
                        <td>
                          {draft.sourceId && publishedIDs.has(draft.sourceId)
                            ? label('已发布 / 草稿', 'Published / draft')
                            : draft.sourceId
                              ? label('已删除 / 草稿', 'Published missing / draft')
                              : label('草稿', 'Draft')}
                        </td>
                        <td>
                          <button
                            className="link-btn"
                            disabled={Boolean(busy)}
                            onClick={() =>
                              void run('open', async () => {
                                openEditor(draft, draft.object)
                                await refreshEditor(draft)
                              })
                            }
                          >
                            {label('编辑', 'Edit')}
                          </button>
                        </td>
                      </tr>
                    ))}
                    {unlinkedRuntime.map((item) => (
                      <tr key={`runtime-${clusterID(item)}`}>
                        <td>
                          <b>{String(item.name || label('未命名集群', 'Unnamed cluster'))}</b>
                        </td>
                        <td>{String(item.type || '—')}</td>
                        <td>
                          {String(item.address || '—')}:{String(item.port || '—')}
                        </td>
                        <td>{label('已发布', 'Published')}</td>
                        <td>
                          <button
                            className="link-btn"
                            disabled={Boolean(busy)}
                            onClick={() => openPublished(item)}
                          >
                            {label('编辑', 'Edit')}
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </>
      )}
      <ConfirmDialog
        open={Boolean(confirm)}
        title={
          confirm?.kind === 'back'
            ? label('放弃未保存改动？', 'Discard unsaved changes?')
            : confirm?.kind === 'rollback'
              ? label('恢复并发布此版本？', 'Restore and publish this revision?')
              : confirm?.kind === 'runtime'
                ? label('删除已发布集群？', 'Delete published cluster?')
                : label('删除草稿？', 'Delete draft?')
        }
        description={
          confirm?.kind === 'draft'
            ? label('已发布配置会保留。', 'Published configuration will remain.')
            : undefined
        }
        confirmLabel={label('确认', 'Confirm')}
        cancelLabel={label('取消', 'Cancel')}
        onConfirm={confirmed}
        onCancel={() => setConfirm(null)}
      />
    </div>
  )
}
