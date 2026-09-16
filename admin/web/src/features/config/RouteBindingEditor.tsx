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

import React, { useEffect, useMemo, useRef, useState } from 'react'
import {
  AlertCircle,
  ArrowLeft,
  CheckCircle2,
  Eye,
  FileCode2,
  GitCompareArrows,
  Plus,
  Save,
  Send,
  Trash2,
} from 'lucide-react'
import { parse as parseYaml, stringify as stringifyYaml } from 'yaml'
import { Locale, translateText } from '../../i18n'
import { routeBindingApi } from '../../services/route-binding-api'
import { ApiError } from '../../services/http'
import { asJsonObject } from '../../types/api'
import type {
  AdminRouteBindingObject,
  RouteBinding,
  RouteBindingDiff,
  RouteBindingParam,
  RouteBindingPublishStatus,
  RouteBindingValidationIssue,
} from '../../types/api'

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']
const PARAM_TYPES = [
  'string',
  'char',
  'short',
  'int',
  'long',
  'float',
  'double',
  'boolean',
  'byte',
  'date',
  'object',
  'java.lang.String',
  'java.lang.Character',
  'java.lang.Short',
  'java.lang.Integer',
  'java.lang.Long',
  'java.lang.Float',
  'java.lang.Double',
  'java.lang.Boolean',
  'java.lang.Byte',
  'java.lang.Object',
  'java.util.Date',
]

type EditorTab = 'form' | 'preview' | 'diff' | 'yaml'
type BusyAction = '' | 'save' | 'publish' | 'validate' | 'preview' | 'diff'
type NoticeTone = 'success' | 'error'
type Notice = { tone: NoticeTone; text: string }

type RouteBindingEditorProps = {
  locale: Locale
  mode: 'create' | 'edit'
  binding?: RouteBinding | null
  loading: boolean
  published: boolean
  publishStatus?: RouteBindingPublishStatus | null
  onBack: () => void
  onSaved: () => void | Promise<void>
}

function createDefaultObject(): AdminRouteBindingObject {
  return {
    kind: 'AdminRouteBinding',
    metadata: { name: '' },
    spec: {
      entry: { protocol: 'http', path: '/api/v1/example', method: 'GET' },
      target: {
        protocol: 'dubbo',
        application: '',
        interface: '',
        method: '',
        version: '',
        group: '',
        cluster: '',
      },
      params: [],
      publish: { mode: 'draft', validate: true },
      extensions: {},
    },
  }
}

function stringValue(value: unknown, fallback = '') {
  return typeof value === 'string' ? value : fallback
}

function numberValue(value: unknown, fallback: number) {
  const number = typeof value === 'number' ? value : Number(value)
  return Number.isFinite(number) ? number : fallback
}

function normaliseObject(value: unknown): AdminRouteBindingObject {
  const root = asJsonObject(value)
  const metadata = asJsonObject(root.metadata)
  const spec = asJsonObject(root.spec)
  const rawEntry = asJsonObject(spec.entry)
  const rawTarget = asJsonObject(spec.target)
  const rawPublish = asJsonObject(spec.publish)
  const rawParams = Array.isArray(spec.params) ? spec.params : []

  const params = rawParams.map((rawParam, index): RouteBindingParam => {
    const param = asJsonObject(rawParam)
    return {
      from: stringValue(param.from),
      to: numberValue(param.to, index),
      type: stringValue(param.type, 'string'),
    }
  })

  return {
    kind: stringValue(root.kind, 'AdminRouteBinding'),
    metadata: { name: stringValue(metadata.name) },
    spec: {
      entry: {
        protocol: stringValue(rawEntry.protocol, 'http'),
        path: stringValue(rawEntry.path, '/api/v1/example'),
        method: stringValue(rawEntry.method, 'GET'),
      },
      target: {
        protocol: stringValue(rawTarget.protocol, 'dubbo'),
        application: stringValue(rawTarget.application),
        interface: stringValue(rawTarget.interface),
        method: stringValue(rawTarget.method),
        version: stringValue(rawTarget.version),
        group: stringValue(rawTarget.group),
        cluster: stringValue(rawTarget.cluster),
      },
      params,
      publish: {
        mode: stringValue(rawPublish.mode, 'draft'),
        validate: typeof rawPublish.validate === 'boolean' ? rawPublish.validate : true,
      },
      extensions: asJsonObject(spec.extensions),
    },
  }
}

function objectSignature(value: AdminRouteBindingObject) {
  return JSON.stringify(value)
}

function stringifyRouteBindingYaml(value: AdminRouteBindingObject) {
  return stringifyYaml(value, { indent: 2, lineWidth: 0 })
}

function formatYamlError(error: unknown) {
  if (!(error instanceof Error)) return String(error)
  const linePos = (error as Error & { linePos?: Array<{ line: number; col: number }> }).linePos?.[0]
  return linePos ? `${error.message} (${linePos.line}:${linePos.col})` : error.message
}

function parseRouteBindingYaml(value: string) {
  try {
    const parsed = parseYaml(value) as unknown
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      throw new Error('YAML 顶层必须是对象')
    }
    return { object: normaliseObject(parsed), error: '' }
  } catch (error: unknown) {
    return { object: null, error: formatYamlError(error) }
  }
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback
}

function errorIssues(error: unknown) {
  return error instanceof ApiError ? error.issues : []
}

function issueFor(issues: RouteBindingValidationIssue[], path: string) {
  return issues.find((issue) => issue.path === path)?.message || ''
}

function formatDiffValue(value: unknown) {
  if (value === null || value === undefined) return '∅'
  if (typeof value === 'string') return value
  const formatted = JSON.stringify(value, null, 2)
  return formatted === undefined ? String(value) : formatted
}

export function RouteBindingEditor({
  locale,
  mode: initialMode,
  binding,
  loading,
  published: initialPublished,
  publishStatus: initialPublishStatus,
  onBack,
  onSaved,
}: RouteBindingEditorProps) {
  const isEnglish = locale === 'en-US'
  const tx = (value: string) => translateText(locale, value)
  const [mode, setMode] = useState(initialMode)
  const [object, setObject] = useState<AdminRouteBindingObject>(() =>
    binding ? normaliseObject(binding.object) : createDefaultObject(),
  )
  const [revision, setRevision] = useState(binding?.revision || 0)
  const [published, setPublished] = useState(initialPublished)
  const [publishStatus, setPublishStatus] = useState<RouteBindingPublishStatus | null>(
    initialPublishStatus || null,
  )
  const [activeTab, setActiveTab] = useState<EditorTab>('form')
  const [busy, setBusy] = useState<BusyAction>('')
  const [dirty, setDirty] = useState(initialMode === 'create')
  const [issues, setIssues] = useState<RouteBindingValidationIssue[]>([])
  const [notice, setNotice] = useState<Notice | null>(null)
  const [previewYaml, setPreviewYaml] = useState('')
  const [routeYaml, setRouteYaml] = useState(() =>
    stringifyRouteBindingYaml(binding ? normaliseObject(binding.object) : createDefaultObject()),
  )
  const [yamlError, setYamlError] = useState('')
  const [diffData, setDiffData] = useState<RouteBindingDiff | null>(null)
  const hydratedBindingKey = useRef('')

  useEffect(() => {
    if (loading || !binding) return
    const bindingKey = `${binding.object.metadata.name}:${binding.revision}`
    if (hydratedBindingKey.current === bindingKey) return
    hydratedBindingKey.current = bindingKey
    const normalized = normaliseObject(binding.object)
    setMode(initialMode)
    setObject(normalized)
    setRevision(binding.revision || 0)
    setPublished(initialPublished)
    setPublishStatus(initialPublishStatus || null)
    setDirty(false)
    setIssues([])
    setNotice(null)
    setPreviewYaml('')
    setRouteYaml(stringifyRouteBindingYaml(normalized))
    setYamlError('')
    setDiffData(null)
    setActiveTab('form')
  }, [binding, initialMode, initialPublishStatus, initialPublished, loading])

  const entry = object.spec.entry
  const target = object.spec.target
  const signature = useMemo(() => objectSignature(object), [object])
  const routeLabel = object.metadata.name || (isEnglish ? 'New API route' : '新建 API 路由')
  const currentStatus = dirty
    ? isEnglish
      ? 'Unsaved changes'
      : '有未保存改动'
    : published
      ? tx('已发布')
      : tx('草稿')

  const updateObject = (next: AdminRouteBindingObject) => {
    setObject(next)
    setRouteYaml(stringifyRouteBindingYaml(next))
    setYamlError('')
    setDirty(true)
    setIssues([])
    setNotice(null)
    setPreviewYaml('')
    setDiffData(null)
  }

  const updateYaml = (value: string) => {
    setRouteYaml(value)
    setDirty(true)
    setIssues([])
    setNotice(null)
    setPreviewYaml('')
    setDiffData(null)

    const parsed = parseRouteBindingYaml(value)
    if (!parsed.object) {
      setYamlError(parsed.error)
      return
    }
    if (mode === 'edit' && parsed.object.metadata.name !== object.metadata.name) {
      setYamlError(
        isEnglish ? 'metadata.name cannot change while editing.' : '编辑时不能修改 metadata.name。',
      )
      return
    }
    setYamlError('')
    setObject(parsed.object)
  }

  const applyYamlToForm = () => {
    const parsed = parseRouteBindingYaml(routeYaml)
    if (!parsed.object) {
      setYamlError(parsed.error)
      return
    }
    if (mode === 'edit' && parsed.object.metadata.name !== object.metadata.name) {
      setYamlError(
        isEnglish ? 'metadata.name cannot change while editing.' : '编辑时不能修改 metadata.name。',
      )
      return
    }
    setObject(parsed.object)
    setYamlError('')
    setIssues([])
    setNotice({
      tone: 'success',
      text: isEnglish ? 'YAML is synced to the form.' : 'YAML 已同步到表单。',
    })
  }

  const updateEntry = (key: 'path' | 'method', value: string) =>
    updateObject({
      ...object,
      spec: { ...object.spec, entry: { ...object.spec.entry, [key]: value } },
    })

  const updateTarget = (
    key: 'application' | 'interface' | 'method' | 'version' | 'group' | 'cluster',
    value: string,
  ) =>
    updateObject({
      ...object,
      spec: { ...object.spec, target: { ...object.spec.target, [key]: value } },
    })

  const updateParam = (index: number, patch: Partial<RouteBindingParam>) => {
    const params = object.spec.params.map((param, paramIndex) =>
      paramIndex === index ? { ...param, ...patch } : param,
    )
    updateObject({ ...object, spec: { ...object.spec, params } })
  }

  const addParam = () => {
    const nextParam: RouteBindingParam = {
      from: '',
      to: object.spec.params.length,
      type: 'string',
    }
    updateObject({
      ...object,
      spec: { ...object.spec, params: [...object.spec.params, nextParam] },
    })
  }

  const removeParam = (index: number) =>
    updateObject({
      ...object,
      spec: {
        ...object.spec,
        params: object.spec.params.filter((_, paramIndex) => paramIndex !== index),
      },
    })

  const setPublishValidation = (validate: boolean) =>
    updateObject({
      ...object,
      spec: { ...object.spec, publish: { ...object.spec.publish, validate } },
    })

  const validate = async () => {
    setBusy('validate')
    setNotice(null)
    try {
      const result = await routeBindingApi.validate(object)
      const normalized = normaliseObject(result.object)
      const changed = objectSignature(normalized) !== signature
      setObject(normalized)
      setRouteYaml(stringifyRouteBindingYaml(normalized))
      setYamlError('')
      setDirty((current) => current || changed)
      setIssues([])
      setNotice({
        tone: 'success',
        text: isEnglish
          ? 'Configuration is valid. Schema defaults are applied.'
          : '配置校验通过，默认值已补齐。',
      })
    } catch (error: unknown) {
      setIssues(errorIssues(error))
      setNotice({ tone: 'error', text: errorMessage(error, tx('配置校验失败')) })
    } finally {
      setBusy('')
    }
  }

  const saveDraft = async (): Promise<RouteBinding | null> => {
    if (!object.metadata.name.trim()) {
      const missingName = [
        { path: 'metadata.name', code: 'required', message: tx('请输入路由名称') },
      ]
      setIssues(missingName)
      setNotice({ tone: 'error', text: tx('请先补充必填字段') })
      setActiveTab('form')
      return null
    }
    setBusy('save')
    setNotice(null)
    try {
      const saved =
        mode === 'create'
          ? await routeBindingApi.create(object)
          : await routeBindingApi.update(object, revision)
      const normalized = normaliseObject(saved.object)
      setObject(normalized)
      setRouteYaml(stringifyRouteBindingYaml(normalized))
      setYamlError('')
      setMode('edit')
      setRevision(saved.revision || 0)
      setDirty(false)
      setIssues([])
      setPreviewYaml('')
      setDiffData(null)
      setPublished(false)
      setNotice({ tone: 'success', text: isEnglish ? 'Draft saved.' : '草稿已保存。' })
      try {
        setPublishStatus(await routeBindingApi.status(normalized.metadata.name))
      } catch {
        // The draft was saved successfully. A status refresh is best effort.
      }
      await onSaved()
      return saved
    } catch (error: unknown) {
      setIssues(errorIssues(error))
      setNotice({ tone: 'error', text: errorMessage(error, tx('保存草稿失败')) })
      setActiveTab('form')
      return null
    } finally {
      setBusy('')
    }
  }

  const publish = async () => {
    if (dirty || mode === 'create' || revision === 0) {
      const saved = await saveDraft()
      if (!saved) return
    }
    setBusy('publish')
    setNotice(null)
    try {
      const routeName = object.metadata.name
      const status = await routeBindingApi.status(routeName)
      const result = await routeBindingApi.publish(routeName, status.draftRevision)
      const nextStatus = await routeBindingApi.status(routeName)
      setPublishStatus(nextStatus)
      setPublished(Boolean(nextStatus.publishedExists && !nextStatus.dirty))
      setDirty(false)
      setDiffData(null)
      setNotice({
        tone: 'success',
        text: isEnglish
          ? result.deletedCount
            ? 'The current route was removed atomically.'
            : 'The current route was published atomically.'
          : result.deletedCount
            ? '当前路由已通过 etcd 事务原子删除。'
            : '当前路由已通过 etcd 事务原子发布。',
      })
      await onSaved()
    } catch (error: unknown) {
      setIssues(errorIssues(error))
      setNotice({ tone: 'error', text: errorMessage(error, tx('发布失败，请刷新后重试')) })
    } finally {
      setBusy('')
    }
  }

  const loadDiff = async () => {
    if (mode === 'create' || dirty) {
      setActiveTab('form')
      setNotice({
        tone: 'error',
        text: isEnglish ? 'Save the draft before viewing Diff.' : '请先保存草稿，再查看 Diff。',
      })
      return
    }
    setActiveTab('diff')
    setBusy('diff')
    setNotice(null)
    try {
      setDiffData(await routeBindingApi.diff(object.metadata.name))
      setIssues([])
    } catch (error: unknown) {
      setIssues(errorIssues(error))
      setNotice({ tone: 'error', text: errorMessage(error, tx('加载 Diff 失败')) })
    } finally {
      setBusy('')
    }
  }

  const preview = async (tab: 'preview' | 'yaml' = 'preview') => {
    setActiveTab(tab)
    setBusy('preview')
    setNotice(null)
    try {
      const result = await routeBindingApi.preview(object)
      const normalized = normaliseObject(result.object)
      setObject(normalized)
      setRouteYaml(stringifyRouteBindingYaml(normalized))
      setYamlError('')
      setDirty((current) => current || objectSignature(normalized) !== signature)
      setPreviewYaml(result.yaml)
      setIssues([])
    } catch (error: unknown) {
      setIssues(errorIssues(error))
      setNotice({ tone: 'error', text: errorMessage(error, tx('生成预览失败')) })
    } finally {
      setBusy('')
    }
  }

  const changeTab = (tab: EditorTab) => {
    setActiveTab(tab)
    if (tab === 'diff' && !diffData) {
      void loadDiff()
      return
    }
    if (tab === 'preview' && !previewYaml) void preview()
  }

  const issueMessage = (path: string) => issueFor(issues, path)

  return (
    <div className="page-resource route-editor-page route-binding-editor">
      <div className="route-editor-backbar">
        <button className="route-editor-back" type="button" onClick={onBack}>
          <ArrowLeft size={15} />
          {isEnglish ? 'Back to API routes' : '返回 API 路由'}
        </button>
        <div
          className="route-editor-revisions"
          aria-label={isEnglish ? 'Route revisions' : '路由版本'}
        >
          <span className="route-revision draft">
            {isEnglish ? 'Draft rev' : '草稿 rev'} {publishStatus?.draftRevision || '-'}
          </span>
          <span className="route-revision published">
            {isEnglish ? 'Published rev' : '已发布 rev'} {publishStatus?.publishedRevision || '-'}
          </span>
        </div>
      </div>

      <div className="route-editor-titlebar">
        <div>
          <p className="eyebrow">API ROUTE</p>
          <h1>{routeLabel}</h1>
          <div className="route-editor-summary">
            <span className={`method ${entry.method.toLowerCase()}`}>{entry.method || 'GET'}</span>
            <span className="route-editor-path">
              {entry.path || (isEnglish ? 'Path required' : '需要填写路径')}
            </span>
            <span
              className={`route-editor-state ${dirty ? 'draft' : published ? 'published' : 'draft'}`}
            >
              <i /> {currentStatus}
            </span>
          </div>
        </div>
        <div className="route-editor-actions">
          <button
            className="secondary"
            type="button"
            disabled={busy !== '' || Boolean(yamlError)}
            onClick={() => void validate()}
          >
            <CheckCircle2 size={14} />
            {busy === 'validate' ? (isEnglish ? 'Checking...' : '校验中...') : tx('校验')}
          </button>
          <button
            className="secondary"
            type="button"
            disabled={busy !== '' || Boolean(yamlError)}
            onClick={() => void saveDraft()}
          >
            <Save size={14} />
            {busy === 'save'
              ? isEnglish
                ? 'Saving...'
                : '保存中...'
              : isEnglish
                ? 'Save draft'
                : '保存草稿'}
          </button>
          <button
            className="primary"
            type="button"
            disabled={busy !== '' || Boolean(yamlError)}
            onClick={() => void publish()}
          >
            <Send size={14} />
            {busy === 'publish' ? (isEnglish ? 'Publishing...' : '发布中...') : tx('发布')}
          </button>
        </div>
      </div>

      <div className="route-editor-contract">
        <span className="route-contract-mark">{isEnglish ? 'ATOMIC PUBLISH' : '原子发布'}</span>
        <span>
          {isEnglish
            ? 'Saving updates this route draft. Publishing replaces only this route in one etcd transaction.'
            : '保存只更新当前路由草稿。发布会通过一次 etcd 事务只替换当前路由。'}
        </span>
      </div>

      {notice && (
        <div
          className={`route-editor-notice ${notice.tone}`}
          role={notice.tone === 'error' ? 'alert' : 'status'}
        >
          {notice.tone === 'error' ? <AlertCircle size={16} /> : <CheckCircle2 size={16} />}
          <span>{notice.text}</span>
        </div>
      )}

      <div
        className="route-editor-tabs"
        role="tablist"
        aria-label={isEnglish ? 'Route editor views' : '路由编辑视图'}
      >
        <button
          className={activeTab === 'form' ? 'active' : ''}
          type="button"
          role="tab"
          aria-selected={activeTab === 'form'}
          onClick={() => changeTab('form')}
        >
          {isEnglish ? 'Form' : '表单'}
        </button>
        <button
          className={activeTab === 'preview' ? 'active' : ''}
          type="button"
          role="tab"
          aria-selected={activeTab === 'preview'}
          onClick={() => changeTab('preview')}
        >
          <Eye size={14} /> {isEnglish ? 'Preview' : '预览'}
        </button>
        <button
          className={activeTab === 'diff' ? 'active' : ''}
          type="button"
          role="tab"
          aria-selected={activeTab === 'diff'}
          onClick={() => changeTab('diff')}
        >
          <GitCompareArrows size={14} /> Diff
        </button>
        <button
          className="unavailable-tab"
          type="button"
          disabled
          title={
            isEnglish
              ? 'History is not exposed by the current Admin API.'
              : '当前 Admin API 尚未提供历史接口。'
          }
        >
          {isEnglish ? 'History' : '历史'}
        </button>
        <button
          className={activeTab === 'yaml' ? 'active' : ''}
          type="button"
          role="tab"
          aria-selected={activeTab === 'yaml'}
          onClick={() => changeTab('yaml')}
        >
          <FileCode2 size={14} /> YAML
        </button>
      </div>

      {loading ? (
        <div
          className="route-binding-loading"
          aria-label={isEnglish ? 'Loading route' : '正在加载路由'}
        >
          <div className="loading-skeleton" />
          <div className="loading-skeleton" />
        </div>
      ) : activeTab === 'form' ? (
        <>
          <div className="route-editor-layout">
            <section className="panel route-editor-card">
              <div className="panel-head">
                <div>
                  <h2>{isEnglish ? 'Identity and entry' : '身份与入口'}</h2>
                  <span>
                    {isEnglish
                      ? 'Define the HTTP route exposed by Pixiu.'
                      : '定义 Pixiu 对外暴露的 HTTP 路由。'}
                  </span>
                </div>
              </div>
              <div className="route-entry-fields">
                <label className="route-field">
                  <span>{isEnglish ? 'Route name' : '路由名称'}</span>
                  <input
                    value={object.metadata.name}
                    readOnly={mode === 'edit'}
                    aria-invalid={Boolean(issueMessage('metadata.name'))}
                    aria-describedby="route-name-help route-name-error"
                    placeholder={isEnglish ? 'create-user' : 'create-user'}
                    onChange={(event) =>
                      updateObject({
                        ...object,
                        metadata: { ...object.metadata, name: event.target.value },
                      })
                    }
                  />
                  <small id="route-name-help">
                    {mode === 'edit'
                      ? isEnglish
                        ? 'The key stays stable while editing.'
                        : '编辑时保持路由键稳定。'
                      : isEnglish
                        ? 'Used as the AdminRouteBinding metadata.name.'
                        : '对应 AdminRouteBinding 的 metadata.name。'}
                  </small>
                  {issueMessage('metadata.name') && (
                    <small id="route-name-error" className="route-field-error">
                      {issueMessage('metadata.name')}
                    </small>
                  )}
                </label>
                <label className="route-field">
                  <span>{isEnglish ? 'Entry protocol' : '入口协议'}</span>
                  <select value={entry.protocol} disabled>
                    <option value="http">HTTP</option>
                  </select>
                  <small>
                    {isEnglish
                      ? 'The current schema accepts HTTP only.'
                      : '当前 schema 仅支持 HTTP。'}
                  </small>
                </label>
                <label className="route-field route-field-wide">
                  <span>{isEnglish ? 'Path' : '路径'}</span>
                  <input
                    value={entry.path}
                    aria-invalid={Boolean(issueMessage('spec.entry.path'))}
                    placeholder="/api/v1/users/:id"
                    onChange={(event) => updateEntry('path', event.target.value)}
                  />
                  {issueMessage('spec.entry.path') && (
                    <small className="route-field-error">{issueMessage('spec.entry.path')}</small>
                  )}
                </label>
                <div className="route-field route-field-wide">
                  <span>{isEnglish ? 'HTTP method' : '请求方法'}</span>
                  <div
                    className="route-method-options"
                    role="group"
                    aria-label={isEnglish ? 'HTTP method' : '请求方法'}
                  >
                    {HTTP_METHODS.map((method) => (
                      <button
                        className={entry.method === method ? 'active' : ''}
                        key={method}
                        type="button"
                        aria-pressed={entry.method === method}
                        onClick={() => updateEntry('method', method)}
                      >
                        {method}
                      </button>
                    ))}
                  </div>
                  {issueMessage('spec.entry.method') && (
                    <small className="route-field-error">{issueMessage('spec.entry.method')}</small>
                  )}
                </div>
              </div>
            </section>

            <section className="panel route-editor-card">
              <div className="panel-head">
                <div>
                  <h2>{isEnglish ? 'Backend target' : '后端目标'}</h2>
                  <span>
                    {isEnglish
                      ? 'Map the request to one Dubbo service method.'
                      : '将请求映射到一个 Dubbo 服务方法。'}
                  </span>
                </div>
              </div>
              <div className="route-target-fields">
                <label className="route-field">
                  <span>{isEnglish ? 'Target protocol' : '目标协议'}</span>
                  <select value={target.protocol} disabled>
                    <option value="dubbo">Dubbo</option>
                  </select>
                  <small>
                    {isEnglish ? 'The current compiler targets Dubbo.' : '当前编译器目标为 Dubbo。'}
                  </small>
                </label>
                <label className="route-field">
                  <span>Cluster</span>
                  <input
                    value={target.cluster}
                    placeholder="user-service"
                    onChange={(event) => updateTarget('cluster', event.target.value)}
                  />
                  {issueMessage('spec.target.cluster') && (
                    <small className="route-field-error">
                      {issueMessage('spec.target.cluster')}
                    </small>
                  )}
                </label>
                <label className="route-field">
                  <span>Application</span>
                  <input
                    value={target.application}
                    placeholder="UserProvider"
                    onChange={(event) => updateTarget('application', event.target.value)}
                  />
                  {issueMessage('spec.target.application') && (
                    <small className="route-field-error">
                      {issueMessage('spec.target.application')}
                    </small>
                  )}
                </label>
                <label className="route-field">
                  <span>Interface</span>
                  <input
                    value={target.interface}
                    placeholder="com.example.user.UserService"
                    onChange={(event) => updateTarget('interface', event.target.value)}
                  />
                  {issueMessage('spec.target.interface') && (
                    <small className="route-field-error">
                      {issueMessage('spec.target.interface')}
                    </small>
                  )}
                </label>
                <label className="route-field">
                  <span>Method</span>
                  <input
                    value={target.method}
                    placeholder="GetUser"
                    onChange={(event) => updateTarget('method', event.target.value)}
                  />
                  {issueMessage('spec.target.method') && (
                    <small className="route-field-error">
                      {issueMessage('spec.target.method')}
                    </small>
                  )}
                </label>
                <label className="route-field">
                  <span>Group</span>
                  <input
                    value={target.group}
                    placeholder="stable"
                    onChange={(event) => updateTarget('group', event.target.value)}
                  />
                </label>
                <label className="route-field">
                  <span>Version</span>
                  <input
                    value={target.version}
                    placeholder="1.0.0"
                    onChange={(event) => updateTarget('version', event.target.value)}
                  />
                </label>
              </div>
            </section>
          </div>

          <section className="panel route-param-panel">
            <div className="panel-head route-param-head">
              <div>
                <h2>{isEnglish ? 'Inbound request' : '入站请求'}</h2>
                <span>
                  {isEnglish
                    ? 'Map HTTP values to zero-based Dubbo argument indexes.'
                    : '将 HTTP 参数映射到从 0 开始的 Dubbo 参数位置。'}
                </span>
              </div>
              <button className="secondary route-add-param" type="button" onClick={addParam}>
                <Plus size={14} /> {isEnglish ? 'Add parameter' : '添加参数'}
              </button>
            </div>
            <div className="route-param-table">
              {object.spec.params.length === 0 ? (
                <div className="route-param-empty">
                  <span>{isEnglish ? 'No parameter mappings yet.' : '暂未配置参数映射。'}</span>
                  <small>
                    {isEnglish
                      ? 'Add a mapping when the Dubbo method expects HTTP input.'
                      : '当 Dubbo 方法需要接收 HTTP 参数时，再添加映射。'}
                  </small>
                </div>
              ) : (
                <>
                  <div className="route-param-header" aria-hidden="true">
                    <span>{isEnglish ? 'HTTP source' : 'HTTP 来源'}</span>
                    <span>{isEnglish ? 'Argument index' : '参数位置'}</span>
                    <span>{isEnglish ? 'Dubbo type' : 'Dubbo 类型'}</span>
                    <span />
                  </div>
                  {object.spec.params.map((param, index) => (
                    <div className="route-param-row" key={`${index}-${param.from}`}>
                      <label className="route-param-field">
                        <span>{isEnglish ? 'HTTP source' : 'HTTP 来源'}</span>
                        <input
                          value={param.from}
                          aria-invalid={Boolean(issueMessage(`spec.params[${index}].from`))}
                          placeholder="queryStrings.page"
                          onChange={(event) => updateParam(index, { from: event.target.value })}
                        />
                        {issueMessage(`spec.params[${index}].from`) && (
                          <small className="route-field-error">
                            {issueMessage(`spec.params[${index}].from`)}
                          </small>
                        )}
                      </label>
                      <label className="route-param-field">
                        <span>{isEnglish ? 'Argument index' : '参数位置'}</span>
                        <input
                          type="number"
                          min="0"
                          value={param.to}
                          aria-invalid={Boolean(issueMessage(`spec.params[${index}].to`))}
                          onChange={(event) =>
                            updateParam(index, { to: Number(event.target.value) })
                          }
                        />
                        {issueMessage(`spec.params[${index}].to`) && (
                          <small className="route-field-error">
                            {issueMessage(`spec.params[${index}].to`)}
                          </small>
                        )}
                      </label>
                      <label className="route-param-field">
                        <span>{isEnglish ? 'Dubbo type' : 'Dubbo 类型'}</span>
                        <select
                          value={param.type}
                          onChange={(event) => updateParam(index, { type: event.target.value })}
                        >
                          {PARAM_TYPES.map((type) => (
                            <option key={type} value={type}>
                              {type}
                            </option>
                          ))}
                        </select>
                        {issueMessage(`spec.params[${index}].type`) && (
                          <small className="route-field-error">
                            {issueMessage(`spec.params[${index}].type`)}
                          </small>
                        )}
                      </label>
                      <button
                        className="route-icon-action"
                        type="button"
                        aria-label={
                          isEnglish ? `Remove parameter ${index + 1}` : `删除第 ${index + 1} 个参数`
                        }
                        onClick={() => removeParam(index)}
                      >
                        <Trash2 size={15} />
                      </button>
                    </div>
                  ))}
                </>
              )}
            </div>
          </section>

          <section className="panel route-publish-panel">
            <div>
              <h2>{isEnglish ? 'Publish settings' : '发布设置'}</h2>
              <span>
                {isEnglish
                  ? 'Draft changes are validated before the atomic publish.'
                  : '草稿变更会在原子发布前完成校验。'}
              </span>
            </div>
            <label className="route-switch">
              <input
                type="checkbox"
                checked={object.spec.publish.validate}
                onChange={(event) => setPublishValidation(event.target.checked)}
              />
              <span className="route-switch-track" />
              <span>{isEnglish ? 'Validate before publish' : '发布前校验'}</span>
            </label>
          </section>
        </>
      ) : activeTab === 'diff' ? (
        <section className="panel route-diff-panel">
          <div className="panel-head">
            <div>
              <h2>Diff</h2>
              <span>
                {isEnglish
                  ? 'Compare this route draft with its published version.'
                  : '对比当前路由草稿与已发布版本。'}
              </span>
            </div>
            <button
              className="secondary"
              type="button"
              disabled={busy !== ''}
              onClick={() => void loadDiff()}
            >
              <GitCompareArrows size={14} />
              {busy === 'diff' ? (isEnglish ? 'Loading...' : '加载中...') : tx('刷新 Diff')}
            </button>
          </div>
          {diffData ? (
            diffData.changes.length === 0 ? (
              <div className="route-diff-empty">
                <CheckCircle2 size={22} />
                <b>{isEnglish ? 'No unpublished changes' : '当前路由没有未发布改动'}</b>
                <span>
                  {isEnglish
                    ? 'The draft and published versions are identical.'
                    : '草稿版本与已发布版本内容一致。'}
                </span>
              </div>
            ) : (
              <div className="route-diff-list">
                {diffData.changes.map((change) => (
                  <div className="route-diff-row" key={change.path}>
                    <code>{change.path}</code>
                    <div className="route-diff-value before">
                      <small>{isEnglish ? 'Published' : '已发布'}</small>
                      <pre>{formatDiffValue(change.before)}</pre>
                    </div>
                    <div className="route-diff-value after">
                      <small>{isEnglish ? 'Draft' : '草稿'}</small>
                      <pre>{formatDiffValue(change.after)}</pre>
                    </div>
                  </div>
                ))}
              </div>
            )
          ) : busy === 'diff' ? (
            <div className="route-diff-empty">
              <div className="loading-skeleton" />
            </div>
          ) : (
            <div className="route-diff-empty">
              <GitCompareArrows size={22} />
              <b>{isEnglish ? 'Diff is not loaded' : '尚未加载 Diff'}</b>
            </div>
          )}
        </section>
      ) : (
        <section className="panel route-preview-panel">
          <div className="panel-head">
            <div>
              <h2>
                {activeTab === 'yaml' ? 'YAML' : isEnglish ? 'Runtime preview' : '运行时预览'}
              </h2>
              <span>
                {isEnglish
                  ? activeTab === 'yaml'
                    ? 'Edit the AdminRouteBinding source. Valid changes sync back to the form.'
                    : 'Generated by the current AdminRouteBinding compiler.'
                  : activeTab === 'yaml'
                    ? '编辑 AdminRouteBinding 源配置，语法有效的改动会同步回表单。'
                    : '由当前 AdminRouteBinding 编译器生成。'}
              </span>
            </div>
            {activeTab === 'yaml' ? (
              <button
                className="secondary"
                type="button"
                disabled={busy !== '' || Boolean(yamlError)}
                onClick={applyYamlToForm}
              >
                <CheckCircle2 size={14} />
                {isEnglish ? 'Sync to form' : '同步到表单'}
              </button>
            ) : (
              <button
                className="secondary"
                type="button"
                disabled={busy !== ''}
                onClick={() => void preview()}
              >
                <Eye size={14} />
                {busy === 'preview'
                  ? isEnglish
                    ? 'Generating...'
                    : '生成中...'
                  : isEnglish
                    ? 'Generate preview'
                    : '生成预览'}
              </button>
            )}
          </div>
          <div className="route-preview-meta">
            <FileCode2 size={15} />
            <span>{activeTab === 'yaml' ? 'route-binding.yaml' : 'api_config.yaml'}</span>
            {activeTab === 'yaml' ? (
              <span className={`route-yaml-sync ${yamlError ? 'error' : 'synced'}`}>
                {yamlError
                  ? isEnglish
                    ? 'Syntax error'
                    : '语法错误'
                  : isEnglish
                    ? 'Synced with form'
                    : '已与表单同步'}
              </span>
            ) : (
              <span className="muted">{isEnglish ? 'Read only' : '只读'}</span>
            )}
          </div>
          {activeTab === 'yaml' ? (
            <>
              <textarea
                className="route-preview-code route-yaml-editor"
                value={routeYaml}
                spellCheck={false}
                aria-label={isEnglish ? 'Editable route YAML' : '可编辑路由 YAML'}
                aria-invalid={Boolean(yamlError)}
                aria-describedby={yamlError ? 'route-yaml-error' : undefined}
                onChange={(event) => updateYaml(event.target.value)}
              />
              {yamlError && (
                <div className="route-yaml-error" id="route-yaml-error" role="alert">
                  <AlertCircle size={15} />
                  <span>{isEnglish ? `YAML error: ${yamlError}` : `YAML 错误：${yamlError}`}</span>
                </div>
              )}
            </>
          ) : previewYaml ? (
            <pre className="route-preview-code">{previewYaml}</pre>
          ) : (
            <div className="route-preview-empty">
              <FileCode2 size={22} />
              <b>{isEnglish ? 'Preview is not generated' : '还没有生成预览'}</b>
              <span>
                {isEnglish
                  ? 'Generate a preview to inspect the legacy runtime configuration.'
                  : '生成预览后查看将写入运行时的 legacy 配置。'}
              </span>
            </div>
          )}
        </section>
      )}
    </div>
  )
}
