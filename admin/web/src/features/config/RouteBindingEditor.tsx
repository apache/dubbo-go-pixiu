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

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Locale, translateText } from '../../i18n'
import { routeBindingApi } from '../../services/route-binding-api'
import { ApiError } from '../../services/http'
import type {
  AdminRouteBindingObject,
  RouteBinding,
  RouteBindingDiff,
  RouteBindingObjectSchema,
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
import {
  createDefaultRouteBinding,
  hasLegacyPublishPreference,
  normaliseRouteBindingObject,
  parseRouteBindingYaml,
  routeEditorOptions,
  stringifyRouteBindingYaml,
} from './RouteBindingEditorModel'

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

function createEmptyObject(): AdminRouteBindingObject {
  return {
    kind: 'AdminRouteBinding',
    metadata: { name: '' },
    spec: {
      entry: { protocol: '', path: '', method: '' },
      target: {
        protocol: '',
        application: '',
        interface: '',
        method: '',
        version: '',
        group: '',
        cluster: '',
      },
      params: [],
      enabled: false,
      extensions: {},
    },
  }
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback
}

function routeBindingHydrationKey(
  binding: RouteBinding | null | undefined,
  mode: RouteBindingEditorProps['mode'],
) {
  if (binding) return `${binding.object.metadata.name}:${binding.revision}`
  if (mode === 'create') return 'new-route'
  return ''
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
  if (enabled) return isEnglish ? 'Draft enabled' : '草稿启用'
  return isEnglish ? 'Draft disabled' : '草稿停用'
}

function publishSuccessMessage(isEnglish: boolean, deleted: boolean) {
  if (isEnglish) {
    return deleted
      ? 'The route was removed atomically from etcd. Pixiu runtime loading is not acknowledged here.'
      : 'The route was committed atomically to etcd. Pixiu runtime loading is not acknowledged here.'
  }
  return deleted
    ? '当前路由已通过 etcd 事务原子删除；此处不会确认 Pixiu 运行时是否已完成加载。'
    : '当前路由已通过 etcd 事务原子写入；此处不会确认 Pixiu 运行时是否已完成加载。'
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
  const [object, setObject] = useState<AdminRouteBindingObject>(createEmptyObject)
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
  const [routeYaml, setRouteYaml] = useState('')
  const [yamlError, setYamlError] = useState('')
  const [yamlDefaultsApplied, setYamlDefaultsApplied] = useState(false)
  const [schema, setSchema] = useState<RouteBindingObjectSchema | null>(null)
  const [schemaLoading, setSchemaLoading] = useState(true)
  const [schemaError, setSchemaError] = useState('')
  const [diffData, setDiffData] = useState<RouteBindingDiff | null>(null)
  const hydratedBindingKey = useRef('')

  const loadSchema = useCallback(async () => {
    setSchemaLoading(true)
    setSchemaError('')
    try {
      setSchema(await routeBindingApi.schema())
    } catch (error: unknown) {
      setSchema(null)
      setSchemaError(
        errorMessage(
          error,
          isEnglish ? 'Could not load the route schema.' : '无法加载路由 Schema。',
        ),
      )
    } finally {
      setSchemaLoading(false)
    }
  }, [isEnglish])

  useEffect(() => {
    void loadSchema()
  }, [loadSchema])

  useEffect(() => {
    if (loading || schemaLoading || !schema) return
    const bindingKey = routeBindingHydrationKey(binding, initialMode)
    if (!bindingKey) return
    if (hydratedBindingKey.current === bindingKey) return
    let normalized: AdminRouteBindingObject
    const removedLegacyPublish = binding ? hasLegacyPublishPreference(binding.object) : false
    try {
      normalized = binding
        ? normaliseRouteBindingObject(binding.object, schema, true)
        : createDefaultRouteBinding(schema)
    } catch (error: unknown) {
      const message = errorMessage(
        error,
        isEnglish ? 'Route data does not match the schema.' : '路由数据不符合 Schema。',
      )
      if (binding) {
        hydratedBindingKey.current = bindingKey
        setObject(createDefaultRouteBinding(schema))
        setRouteYaml(stringifyRouteBindingYaml(binding.object))
        setYamlError(message)
        setDirty(true)
        setActiveTab('yaml')
        setNotice({ tone: 'error', text: message })
      } else {
        setSchemaError(message)
      }
      return
    }
    hydratedBindingKey.current = bindingKey
    setMode(initialMode)
    setObject(normalized)
    setRevision(binding?.revision || 0)
    setPublished(initialPublished)
    setPublishStatus(initialPublishStatus || null)
    setDirty(initialMode === 'create' || removedLegacyPublish)
    setIssues([])
    setNotice(
      removedLegacyPublish
        ? {
            tone: 'success',
            text: isEnglish
              ? 'An obsolete no-op spec.publish field was removed. Save this draft to persist the cleanup.'
              : '已移除旧的无效 spec.publish 字段；保存草稿后才会持久化这次清理。',
          }
        : null,
    )
    setPreviewYaml('')
    setRouteYaml(stringifyRouteBindingYaml(normalized))
    setYamlError('')
    setYamlDefaultsApplied(false)
    setDiffData(null)
    setActiveTab('form')
  }, [
    binding,
    initialMode,
    initialPublishStatus,
    initialPublished,
    isEnglish,
    loading,
    schema,
    schemaLoading,
  ])

  const entry = object.spec.entry
  const signature = useMemo(() => JSON.stringify(object), [object])
  const routeLabel = object.metadata.name || (isEnglish ? 'New API route' : '新建 API 路由')
  const currentStatus = currentStatusLabel(isEnglish, dirty, published, tx)
  const lifecycleStatus = lifecycleStatusLabel(isEnglish, object.spec.enabled)
  const options = useMemo(() => (schema ? routeEditorOptions(schema) : undefined), [schema])

  const editorBlocked = schemaLoading || !schema || Boolean(schemaError)

  const updateObject = (next: AdminRouteBindingObject) => {
    setObject(next)
    setRouteYaml(stringifyRouteBindingYaml(next))
    setYamlError('')
    setYamlDefaultsApplied(false)
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

    if (!schema) {
      setYamlError(isEnglish ? 'Route schema is not loaded.' : '路由 Schema 尚未加载。')
      return
    }
    const parsed = parseRouteBindingYaml(value, schema)
    if (!parsed.object) {
      setYamlError(parsed.error)
      setYamlDefaultsApplied(false)
      return
    }
    if (mode === 'edit' && parsed.object.metadata.name !== object.metadata.name) {
      setYamlError(
        isEnglish ? 'metadata.name cannot change while editing.' : '编辑时不能修改 metadata.name。',
      )
      return
    }
    setYamlError('')
    setYamlDefaultsApplied(parsed.defaultsApplied)
    setObject(parsed.object)
  }

  const applyYamlToForm = () => {
    if (!schema) {
      setYamlError(isEnglish ? 'Route schema is not loaded.' : '路由 Schema 尚未加载。')
      return
    }
    const parsed = parseRouteBindingYaml(routeYaml, schema)
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
    setRouteYaml(stringifyRouteBindingYaml(parsed.object))
    setYamlError('')
    setYamlDefaultsApplied(false)
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
      type: options?.paramTypes[0] || '',
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

  const setRouteEnabled = (enabled: boolean) =>
    updateObject({ ...object, spec: { ...object.spec, enabled } })

  const validate = async () => {
    if (yamlError || editorBlocked || !schema) return
    setBusy('validate')
    setNotice(null)
    try {
      const result = await routeBindingApi.validate(object)
      const normalized = normaliseRouteBindingObject(result.object, schema)
      const changed = JSON.stringify(normalized) !== signature
      setObject(normalized)
      setRouteYaml(stringifyRouteBindingYaml(normalized))
      setYamlError('')
      setYamlDefaultsApplied(false)
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
    if (yamlError || editorBlocked || !schema) return null
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
          : await routeBindingApi.update(
              binding?.object.metadata.name || object.metadata.name,
              object,
              revision,
            )
      const normalized = normaliseRouteBindingObject(saved.object, schema)
      setObject(normalized)
      setRouteYaml(stringifyRouteBindingYaml(normalized))
      setYamlError('')
      setYamlDefaultsApplied(false)
      setMode('edit')
      setRevision(saved.revision || 0)
      setDirty(false)
      setIssues([])
      setPreviewYaml('')
      setDiffData(null)
      setPublished(false)
      setNotice({ tone: 'success', text: isEnglish ? 'Draft saved.' : '草稿已保存。' })
      try {
        const nextStatus = await routeBindingApi.status(normalized.metadata.name)
        setPublishStatus(nextStatus)
        setPublished(Boolean(nextStatus.publishedExists && !nextStatus.dirty))
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
    if (yamlError || editorBlocked || !schema) return
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
    if (yamlError || editorBlocked || !schema) return
    setActiveTab(tab)
    setBusy('preview')
    setNotice(null)
    try {
      const result = await routeBindingApi.preview(object)
      const normalized = normaliseRouteBindingObject(result.object, schema)
      setObject(normalized)
      setRouteYaml(stringifyRouteBindingYaml(normalized))
      setYamlError('')
      setYamlDefaultsApplied(false)
      setDirty((current) => current || JSON.stringify(normalized) !== signature)
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
        actionsDisabled={editorBlocked}
        onBack={onBack}
        onValidate={() => void validate()}
        onSave={() => void saveDraft()}
        onPublish={() => void publish()}
      />
      <RouteEditorContract isEnglish={isEnglish} />
      {notice && <RouteEditorNotice notice={notice} />}
      {schemaError && (
        <div className="route-editor-notice error" role="alert">
          <span>{schemaError}</span>
          <button className="secondary" type="button" onClick={() => void loadSchema()}>
            {isEnglish ? 'Retry schema' : '重试加载 Schema'}
          </button>
        </div>
      )}
      <RouteEditorTabs isEnglish={isEnglish} activeTab={activeTab} onChange={changeTab} />
      <RouteEditorContent
        isEnglish={isEnglish}
        loading={loading || schemaLoading || Boolean(schemaError) || !schema}
        activeTab={activeTab}
        mode={mode}
        object={object}
        busy={busy}
        diffData={diffData}
        previewYaml={previewYaml}
        routeYaml={routeYaml}
        yamlError={yamlError}
        yamlDefaultsApplied={yamlDefaultsApplied}
        options={options}
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
        onRefreshDiff={() => void loadDiff()}
        onApplyYaml={applyYamlToForm}
        onChangeYaml={updateYaml}
        onPreview={() => void preview()}
      />
    </div>
  )
}
