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

import { useEffect, useMemo, useRef, useState } from 'react'
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
import {
  RouteEditorContent,
  RouteEditorContract,
  RouteEditorHeader,
  RouteEditorNotice,
  RouteEditorTabs,
} from './RouteBindingEditorSections'
import type { BusyAction, EditorTab, Notice } from './RouteBindingEditorSections'

type RouteBindingEditorProps = {
  readonly locale: Locale
  readonly mode: 'create' | 'edit'
  readonly binding?: RouteBinding | null
  readonly loading: boolean
  readonly published: boolean
  readonly publishStatus?: RouteBindingPublishStatus | null
  readonly onBack: () => void
  readonly onSaved: () => void | Promise<void>
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
      enabled: true,
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
      enabled: typeof spec.enabled === 'boolean' ? spec.enabled : true,
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

function currentStatusLabel(
  isEnglish: boolean,
  dirty: boolean,
  published: boolean,
  translate: (value: string) => string,
) {
  if (dirty) return isEnglish ? 'Unsaved changes' : '有未保存改动'
  if (published) return translate('已发布')
  return translate('草稿')
}

function lifecycleStatusLabel(isEnglish: boolean, enabled: boolean) {
  if (enabled) return isEnglish ? 'Enabled' : '已启用'
  return isEnglish ? 'Disabled' : '已停用'
}

function publishSuccessMessage(isEnglish: boolean, deleted: boolean) {
  if (isEnglish) {
    return deleted
      ? 'The current route was removed atomically.'
      : 'The current route was published atomically.'
  }
  return deleted ? '当前路由已通过 etcd 事务原子删除。' : '当前路由已通过 etcd 事务原子发布。'
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
  const signature = useMemo(() => objectSignature(object), [object])
  const routeLabel = object.metadata.name || (isEnglish ? 'New API route' : '新建 API 路由')
  const currentStatus = currentStatusLabel(isEnglish, dirty, published, tx)
  const lifecycleStatus = lifecycleStatusLabel(isEnglish, object.spec.enabled)

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

  const setRouteEnabled = (enabled: boolean) =>
    updateObject({ ...object, spec: { ...object.spec, enabled } })

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
        text: publishSuccessMessage(isEnglish, result.deletedCount > 0),
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
      <RouteEditorHeader
        isEnglish={isEnglish}
        routeLabel={routeLabel}
        entry={entry}
        publishStatus={publishStatus}
        dirty={dirty}
        published={published}
        enabled={object.spec.enabled}
        currentStatus={currentStatus}
        lifecycleStatus={lifecycleStatus}
        busy={busy}
        yamlError={yamlError}
        onBack={onBack}
        onValidate={() => void validate()}
        onSave={() => void saveDraft()}
        onPublish={() => void publish()}
      />
      <RouteEditorContract isEnglish={isEnglish} />
      {notice && <RouteEditorNotice notice={notice} />}
      <RouteEditorTabs isEnglish={isEnglish} activeTab={activeTab} onChange={changeTab} />
      <RouteEditorContent
        isEnglish={isEnglish}
        loading={loading}
        activeTab={activeTab}
        mode={mode}
        object={object}
        busy={busy}
        diffData={diffData}
        previewYaml={previewYaml}
        routeYaml={routeYaml}
        yamlError={yamlError}
        issueMessage={issueMessage}
        onNameChange={(name) =>
          updateObject({
            ...object,
            metadata: { ...object.metadata, name },
          })
        }
        onEntryChange={updateEntry}
        onTargetChange={updateTarget}
        onParamAdd={addParam}
        onParamUpdate={updateParam}
        onParamRemove={removeParam}
        onEnabledChange={setRouteEnabled}
        onValidateChange={setPublishValidation}
        onRefreshDiff={() => void loadDiff()}
        onApplyYaml={applyYamlToForm}
        onChangeYaml={updateYaml}
        onPreview={() => void preview()}
      />
    </div>
  )
}
