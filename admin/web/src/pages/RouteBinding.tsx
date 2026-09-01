import { useCallback, useEffect, useMemo, useState } from 'react'
import yaml from 'js-yaml'
import {
  ArrowClockwise,
  ArrowCounterClockwise,
  ClockCounterClockwise,
  CloudArrowUp,
  Eye,
  FilePlus,
  FloppyDisk,
  GitDiff,
  Plus,
  Trash,
} from '@phosphor-icons/react'
import {
  deleteRouteBindingDraft,
  getRouteBindingDetail,
  getRouteBindingHistory,
  getRouteBindingSchema,
  listRouteBindings,
  previewRouteBinding,
  publishRouteBinding,
  rollbackRouteBinding,
  saveRouteBindingDraft,
} from '../api'
import {
  AdminObject,
  ObjectSchema,
  RouteBindingPreview,
  RouteBindingRecord,
  RouteBindingSummary,
} from '../types'
import { Button, EmptyState, ErrorState, Field, formatTime, Loading, Panel, Select, Tag, TextInput } from '../components/common'
import { ConfirmModal, Modal } from '../components/Modal'
import { YamlEditor } from '../components/YamlEditor'
import { useToast } from '../components/Toast'

// ---------------------------------------------------------------------------
// Model helpers
// ---------------------------------------------------------------------------

const KIND = 'AdminRouteBinding'
const NAME_PATTERN = /^[A-Za-z0-9][A-Za-z0-9._-]{0,62}$/

const FALLBACK_ENTRY_PROTOCOLS = ['http']
const FALLBACK_HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']
const FALLBACK_TARGET_PROTOCOLS = ['dubbo']
const FALLBACK_PARAM_TYPES = [
  'string',
  'int',
  'long',
  'float',
  'double',
  'boolean',
  'char',
  'short',
  'date',
  'object',
  'java.lang.String',
  'java.lang.Integer',
  'java.lang.Long',
  'java.lang.Double',
  'java.lang.Boolean',
]

function emptyModel(): AdminObject {
  return {
    kind: KIND,
    metadata: { name: '' },
    spec: {
      entry: { protocol: 'http', path: '/api/v1/example/:id', method: 'GET' },
      target: { protocol: 'dubbo', application: '', interface: '', method: '', version: '', group: '', cluster: '' },
      params: [],
      timeout: '1s',
      publish: { mode: 'draft', validate: true },
    },
  }
}

function dumpModel(model: AdminObject): string {
  return yaml.dump(model, { indent: 2, lineWidth: 100, noRefs: true })
}

/** Parse YAML or JSON text into an AdminObject; throws with a readable message. */
function parseModel(text: string): AdminObject {
  const value = yaml.load(text)
  if (value === null || typeof value !== 'object' || Array.isArray(value)) {
    throw new Error('内容必须是 YAML/JSON 对象')
  }
  const object = value as Record<string, unknown>
  if (typeof object.kind !== 'string' || !object.kind) {
    throw new Error('缺少 kind 字段')
  }
  const metadata = (object.metadata ?? {}) as Record<string, unknown>
  if (metadata !== null && typeof metadata !== 'object') {
    throw new Error('metadata 必须是对象')
  }
  const spec = (object.spec ?? {}) as Record<string, unknown>
  if (spec !== null && typeof spec !== 'object') {
    throw new Error('spec 必须是对象')
  }
  return { kind: object.kind, metadata: metadata as AdminObject['metadata'], spec }
}

type SpecValue = Record<string, unknown>

function specOf(model: AdminObject, key: string): SpecValue {
  const value = model.spec[key]
  return value !== null && typeof value === 'object' && !Array.isArray(value) ? (value as SpecValue) : {}
}

/** Patch one spec sub-object, preserving every other (including unknown) field. */
function patchSpec(model: AdminObject, key: string, patch: SpecValue): AdminObject {
  return { ...model, spec: { ...model.spec, [key]: { ...specOf(model, key), ...patch } } }
}

interface ParamRow {
  from: string
  to: number | string
  type: string
}

function paramsOf(model: AdminObject): ParamRow[] {
  const raw = model.spec.params
  if (!Array.isArray(raw)) return []
  return raw
    .filter((item): item is SpecValue => item !== null && typeof item === 'object')
    .map((item) => ({
      from: typeof item.from === 'string' ? item.from : '',
      to: typeof item.to === 'number' ? item.to : String(item.to ?? ''),
      type: typeof item.type === 'string' ? item.type : 'string',
    }))
}

function enumOptions(schema: ObjectSchema | null, path: string[], fallback: string[]): string[] {
  let node = schema?.fields
  let options: unknown[] | undefined
  for (let i = 0; i < path.length; i++) {
    const field = node?.[path[i]]
    if (!field) break
    if (i === path.length - 1) options = field.enum
    node = field.properties
  }
  const values = (options ?? []).filter((item): item is string => typeof item === 'string')
  return values.length > 0 ? values : fallback
}

// ---------------------------------------------------------------------------
// Diff rendering
// ---------------------------------------------------------------------------

function DiffView({ diff }: { diff: string }) {
  const lines = diff.split('\n')
  return (
    <pre className="slim-scroll max-h-72 overflow-auto rounded-md bg-slate-900 p-3 font-mono text-xs leading-5">
      {lines.map((line, index) => {
        let cls = 'text-slate-400'
        if (line.startsWith('+') && !line.startsWith('+++')) cls = 'text-emerald-400'
        else if (line.startsWith('-') && !line.startsWith('---')) cls = 'text-rose-400'
        else if (line.startsWith('---') || line.startsWith('+++')) cls = 'text-slate-500'
        return (
          <div key={index} className={cls}>
            {line || ' '}
          </div>
        )
      })}
    </pre>
  )
}

// ---------------------------------------------------------------------------
// Status tag
// ---------------------------------------------------------------------------

function StatusTag({ status }: { status: string }) {
  if (status === 'published') return <Tag tone="green">已发布</Tag>
  if (status === 'modified') return <Tag tone="amber">已修改</Tag>
  return <Tag tone="blue">草稿</Tag>
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

type EditorView = 'form' | 'yaml'

export function RouteBindingPage() {
  const toast = useToast()

  const [bindings, setBindings] = useState<RouteBindingSummary[]>([])
  const [listLoading, setListLoading] = useState(true)
  const [listError, setListError] = useState('')

  const [schema, setSchema] = useState<ObjectSchema | null>(null)
  const [schemaReady, setSchemaReady] = useState(false)

  const [selectedName, setSelectedName] = useState<string | null>(null)
  const [model, setModel] = useState<AdminObject>(() => emptyModel())
  const [text, setText] = useState(() => dumpModel(emptyModel()))
  const [view, setView] = useState<EditorView>('form')
  const [parseError, setParseError] = useState('')
  const [dirty, setDirty] = useState(false)
  const [recordMeta, setRecordMeta] = useState<{ revision?: number; updatedAt?: string; source: 'draft' | 'published' | null }>({
    source: null,
  })

  const [saving, setSaving] = useState(false)
  const [previewing, setPreviewing] = useState(false)
  const [preview, setPreview] = useState<RouteBindingPreview | null>(null)

  const [publishOpen, setPublishOpen] = useState(false)
  const [publishing, setPublishing] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [deleting, setDeleting] = useState(false)

  const [historyOpen, setHistoryOpen] = useState(false)
  const [historyLoading, setHistoryLoading] = useState(false)
  const [history, setHistory] = useState<RouteBindingRecord[]>([])
  const [rollbackTarget, setRollbackTarget] = useState<RouteBindingRecord | null>(null)
  const [rollingBack, setRollingBack] = useState(false)

  const entryProtocols = useMemo(() => enumOptions(schema, ['entry', 'protocol'], FALLBACK_ENTRY_PROTOCOLS), [schema])
  const httpMethods = useMemo(() => enumOptions(schema, ['entry', 'method'], FALLBACK_HTTP_METHODS), [schema])
  const targetProtocols = useMemo(() => enumOptions(schema, ['target', 'protocol'], FALLBACK_TARGET_PROTOCOLS), [schema])

  // ---- data loading ----

  const loadList = useCallback(async () => {
    setListLoading(true)
    setListError('')
    try {
      setBindings((await listRouteBindings()) ?? [])
    } catch (err) {
      setListError(err instanceof Error ? err.message : '加载路由绑定列表失败')
    } finally {
      setListLoading(false)
    }
  }, [])

  const loadSchema = useCallback(async () => {
    try {
      const schemas = await getRouteBindingSchema()
      setSchema((schemas ?? []).find((item) => item.kind === KIND) ?? null)
      setSchemaReady(true)
    } catch {
      setSchemaReady(false)
    }
  }, [])

  useEffect(() => {
    void loadList()
    void loadSchema()
  }, [loadList, loadSchema])

  // ---- editor state transitions ----

  const applyModel = useCallback((next: AdminObject, markDirty: boolean) => {
    setModel(next)
    setText(dumpModel(next))
    setParseError('')
    if (markDirty) setDirty(true)
  }, [])

  const loadIntoEditor = useCallback(
    (object: AdminObject, name: string | null, source: 'draft' | 'published' | null, record?: RouteBindingRecord) => {
      const cloned = yaml.load(yaml.dump(object)) as AdminObject
      applyModel(cloned, false)
      setSelectedName(name)
      setDirty(false)
      setPreview(null)
      setRecordMeta({ revision: record?.revision, updatedAt: record?.updatedAt, source })
    },
    [applyModel],
  )

  const handleSelect = useCallback(
    async (name: string) => {
      try {
        const detail = await getRouteBindingDetail(name)
        if (detail.draft) {
          loadIntoEditor(detail.draft.object, name, 'draft', detail.draft)
        } else if (detail.published) {
          loadIntoEditor(detail.published.object, name, 'published', detail.published)
        }
      } catch (err) {
        toast.error(err instanceof Error ? err.message : '加载路由绑定失败')
      }
    },
    [loadIntoEditor, toast],
  )

  const handleNew = useCallback(() => {
    loadIntoEditor(emptyModel(), null, null)
    setView('form')
  }, [loadIntoEditor])

  const handleReload = useCallback(async () => {
    if (!selectedName) {
      handleNew()
      return
    }
    await handleSelect(selectedName)
    toast.info('已重新载入服务端配置')
  }, [selectedName, handleNew, handleSelect, toast])

  // ---- form <-> yaml sync ----

  const handleFormChange = useCallback(
    (next: AdminObject) => {
      applyModel(next, true)
    },
    [applyModel],
  )

  const handleTextChange = useCallback((next: string) => {
    setText(next)
    setDirty(true)
    try {
      const parsed = parseModel(next)
      setModel(parsed)
      setParseError('')
    } catch (err) {
      // Keep the previous model; block form view and submission.
      setParseError(err instanceof Error ? err.message : '无法解析内容')
    }
  }, [])

  const switchView = (next: EditorView) => {
    if (next === 'form' && parseError) return
    setView(next)
  }

  // ---- validations shared by submit actions ----

  const currentName = (model.metadata.name ?? '').trim()
  const nameError = !currentName
    ? '请填写模型名称'
    : !NAME_PATTERN.test(currentName)
      ? '名称只能包含字母、数字、点、下划线或中划线，且以字母或数字开头'
      : ''

  const guardSubmittable = (): boolean => {
    if (parseError) {
      toast.error('YAML/JSON 无法解析，请先修正后再提交')
      return false
    }
    if (nameError) {
      toast.error(nameError)
      return false
    }
    return true
  }

  // ---- actions ----

  const handleSaveDraft = async () => {
    if (!guardSubmittable()) return
    setSaving(true)
    try {
      const record = await saveRouteBindingDraft(model)
      toast.success(`草稿已保存（revision ${record.revision}）`)
      setSelectedName(currentName)
      setDirty(false)
      setRecordMeta({ revision: record.revision, updatedAt: record.updatedAt, source: 'draft' })
      void loadList()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存草稿失败')
    } finally {
      setSaving(false)
    }
  }

  const handlePreview = async () => {
    if (!guardSubmittable()) return
    setPreviewing(true)
    try {
      setPreview(await previewRouteBinding(model))
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '预览失败')
    } finally {
      setPreviewing(false)
    }
  }

  const handlePublish = async () => {
    setPublishing(true)
    try {
      const record = await publishRouteBinding(currentName)
      toast.success(`已发布 revision ${record.revision}，运行时已更新`)
      setPublishOpen(false)
      setDirty(false)
      void loadList()
      void handleSelect(currentName)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '发布失败')
    } finally {
      setPublishing(false)
    }
  }

  const handleDeleteDraft = async () => {
    setDeleting(true)
    try {
      await deleteRouteBindingDraft(currentName)
      toast.success('草稿已删除（已发布配置保留）')
      setDeleteOpen(false)
      void loadList()
      handleNew()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '删除草稿失败')
    } finally {
      setDeleting(false)
    }
  }

  const openHistory = async () => {
    if (!currentName) {
      toast.error('请先选择或保存一个路由绑定')
      return
    }
    setHistoryOpen(true)
    setHistoryLoading(true)
    try {
      setHistory((await getRouteBindingHistory(currentName)) ?? [])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '加载发布历史失败')
      setHistoryOpen(false)
    } finally {
      setHistoryLoading(false)
    }
  }

  const handleRollback = async () => {
    if (!rollbackTarget) return
    setRollingBack(true)
    try {
      const record = await rollbackRouteBinding(currentName, rollbackTarget.revision)
      toast.success(`已回滚到 revision ${rollbackTarget.revision}（新 revision ${record.revision}）`)
      setRollbackTarget(null)
      setHistoryOpen(false)
      void loadList()
      void handleSelect(currentName)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '回滚失败')
    } finally {
      setRollingBack(false)
    }
  }

  // ---- form field helpers ----

  const entry = specOf(model, 'entry')
  const target = specOf(model, 'target')
  const publish = specOf(model, 'publish')
  const params = paramsOf(model)
  const timeout = typeof model.spec.timeout === 'string' ? model.spec.timeout : String(model.spec.timeout ?? '')
  const hasExtensions = model.spec.extensions !== undefined && Object.keys(specOf(model, 'extensions')).length > 0

  const setParams = (rows: ParamRow[]) => {
    const next: AdminObject = {
      ...model,
      spec: {
        ...model.spec,
        params: rows.map((row) => ({ from: row.from, to: typeof row.to === 'number' ? row.to : Number(row.to) || 0, type: row.type })),
      },
    }
    handleFormChange(next)
  }

  // -------------------------------------------------------------------------
  // render
  // -------------------------------------------------------------------------

  return (
    <div className="mx-auto max-w-[1400px] space-y-4">
      <Panel bodyClassName="flex flex-wrap items-center justify-between gap-3 py-3">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-base font-semibold text-slate-900">路由模型</h1>
            <Tag tone={schemaReady ? 'green' : 'slate'}>{schemaReady ? 'Schema 已加载' : 'Schema 未加载'}</Tag>
          </div>
          <p className="mt-0.5 text-xs text-slate-500">
            用 AdminRouteBinding 描述 HTTP 到 Dubbo 的调用；草稿不影响运行中的 Pixiu，发布时投影为旧版 Resource/Method。
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => void loadList()}>
            <ArrowClockwise size={15} />
            刷新
          </Button>
          <Button variant="primary" onClick={handleNew}>
            <FilePlus size={15} />
            新建
          </Button>
        </div>
      </Panel>

      <div className="grid grid-cols-1 gap-4 xl:grid-cols-[300px_1fr]">
        <Panel title={`绑定列表（${bindings.length}）`} bodyClassName="p-0" className="xl:sticky xl:top-0">
          {listLoading ? (
            <Loading />
          ) : listError ? (
            <ErrorState text={listError} onRetry={loadList} />
          ) : bindings.length === 0 ? (
            <EmptyState text="暂无路由绑定，点击“新建”开始。" />
          ) : (
            <ul className="slim-scroll max-h-[560px] divide-y divide-slate-50 overflow-y-auto">
              {bindings.map((item) => (
                <li key={item.name}>
                  <button
                    type="button"
                    className={`w-full px-4 py-3 text-left transition hover:bg-slate-50 ${
                      selectedName === item.name ? 'bg-blue-50/70 ring-1 ring-inset ring-blue-600/20' : ''
                    }`}
                    onClick={() => void handleSelect(item.name)}
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span className="truncate font-mono text-[13px] font-medium text-slate-800">{item.name}</span>
                      <StatusTag status={item.status} />
                    </div>
                    <div className="mt-1 flex items-center gap-3 text-[11px] text-slate-400">
                      {item.draftRevision ? <span>草稿 r{item.draftRevision}</span> : null}
                      {item.publishedRevision ? <span>发布 r{item.publishedRevision}</span> : null}
                      <span>{formatTime(item.updatedAt)}</span>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </Panel>

        <div className="min-w-0 space-y-4">
          <Panel
            title={
              <div className="flex items-center gap-2">
                <span>{selectedName ? `编辑：${selectedName}` : '新建路由绑定'}</span>
                {dirty ? <Tag tone="amber">未保存</Tag> : null}
                {recordMeta.source ? (
                  <span className="text-xs font-normal text-slate-400">
                    {recordMeta.source === 'draft' ? '草稿' : '已发布'}
                    {recordMeta.revision ? ` r${recordMeta.revision}` : ''}
                    {recordMeta.updatedAt ? ` · ${formatTime(recordMeta.updatedAt)}` : ''}
                  </span>
                ) : null}
              </div>
            }
            extra={
              <div className="flex flex-wrap items-center gap-2">
                <Button onClick={() => void handleReload()} title="重新载入服务端配置">
                  <ArrowCounterClockwise size={15} />
                </Button>
                <Button variant="primary" loading={saving} disabled={!!parseError} onClick={() => void handleSaveDraft()}>
                  <FloppyDisk size={15} />
                  保存草稿
                </Button>
                <Button loading={previewing} disabled={!!parseError} onClick={() => void handlePreview()}>
                  <Eye size={15} />
                  预览
                </Button>
                <Button
                  variant="primary"
                  disabled={!!parseError || !currentName || !!nameError}
                  onClick={() => setPublishOpen(true)}
                >
                  <CloudArrowUp size={15} />
                  发布
                </Button>
                <Button disabled={!currentName} onClick={() => void openHistory()}>
                  <ClockCounterClockwise size={15} />
                  历史
                </Button>
                <Button variant="danger" disabled={!currentName} onClick={() => setDeleteOpen(true)}>
                  <Trash size={15} />
                  删除草稿
                </Button>
              </div>
            }
          >
            <div className="mb-3 flex items-center justify-between">
              <div className="flex rounded-lg bg-slate-100 p-1">
                <button
                  type="button"
                  className={`rounded-md px-3 py-1 text-xs font-medium transition ${
                    view === 'form' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'
                  }`}
                  onClick={() => switchView('form')}
                  disabled={!!parseError}
                  title={parseError ? '修正 YAML/JSON 后才能切回表单' : undefined}
                >
                  结构化表单
                </button>
                <button
                  type="button"
                  className={`rounded-md px-3 py-1 text-xs font-medium transition ${
                    view === 'yaml' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'
                  }`}
                  onClick={() => switchView('yaml')}
                >
                  模型配置（YAML / JSON）
                </button>
              </div>
              <span className="text-[11px] text-slate-400">表单仅管理已定义字段，其余合法字段会被保留</span>
            </div>

            {parseError ? (
              <p className="mb-3 rounded-md bg-rose-50 px-3 py-2 text-xs text-rose-600">
                配置无法解析：{parseError}。修正前无法切回表单或提交，避免用不完整的表单数据覆盖原始配置。
              </p>
            ) : null}

            {view === 'form' ? (
              <div className="space-y-5">
                <section>
                  <h3 className="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-400">基本信息</h3>
                  <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <Field label="模型名称" error={dirty || currentName ? nameError : ''} hint="字母、数字、点、下划线或中划线">
                      <TextInput
                        value={model.metadata.name ?? ''}
                        placeholder="user-get"
                        onChange={(e) => handleFormChange({ ...model, metadata: { ...model.metadata, name: e.target.value } })}
                      />
                    </Field>
                    <Field label="超时时间" hint="Duration，如 300ms、1s、10s">
                      <TextInput
                        value={timeout}
                        placeholder="1s"
                        onChange={(e) => handleFormChange({ ...model, spec: { ...model.spec, timeout: e.target.value } })}
                      />
                    </Field>
                  </div>
                </section>

                <section>
                  <h3 className="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-400">HTTP 入口</h3>
                  <div className="grid grid-cols-1 gap-3 md:grid-cols-[140px_1fr_140px]">
                    <Field label="协议">
                      <Select
                        value={String(entry.protocol ?? 'http')}
                        onChange={(e) => handleFormChange(patchSpec(model, 'entry', { protocol: e.target.value }))}
                      >
                        {entryProtocols.map((item) => (
                          <option key={item} value={item}>
                            {item}
                          </option>
                        ))}
                      </Select>
                    </Field>
                    <Field label="路径" hint="支持 :param 路径参数，如 /api/v1/users/:id">
                      <TextInput
                        value={String(entry.path ?? '')}
                        placeholder="/api/v1/users/:id"
                        className="font-mono"
                        onChange={(e) => handleFormChange(patchSpec(model, 'entry', { path: e.target.value }))}
                      />
                    </Field>
                    <Field label="HTTP 方法">
                      <Select
                        value={String(entry.method ?? 'GET')}
                        onChange={(e) => handleFormChange(patchSpec(model, 'entry', { method: e.target.value }))}
                      >
                        {httpMethods.map((item) => (
                          <option key={item} value={item}>
                            {item}
                          </option>
                        ))}
                      </Select>
                    </Field>
                  </div>
                </section>

                <section>
                  <h3 className="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-400">Dubbo 目标</h3>
                  <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <Field label="协议">
                      <Select
                        value={String(target.protocol ?? 'dubbo')}
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { protocol: e.target.value }))}
                      >
                        {targetProtocols.map((item) => (
                          <option key={item} value={item}>
                            {item}
                          </option>
                        ))}
                      </Select>
                    </Field>
                    <Field label="应用 Application">
                      <TextInput
                        value={String(target.application ?? '')}
                        placeholder="UserProvider"
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { application: e.target.value }))}
                      />
                    </Field>
                    <Field label="接口 Interface">
                      <TextInput
                        value={String(target.interface ?? '')}
                        placeholder="org.apache.dubbo.UserService"
                        className="font-mono"
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { interface: e.target.value }))}
                      />
                    </Field>
                    <Field label="方法 Method">
                      <TextInput
                        value={String(target.method ?? '')}
                        placeholder="GetUser"
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { method: e.target.value }))}
                      />
                    </Field>
                    <Field label="版本 Version">
                      <TextInput
                        value={String(target.version ?? '')}
                        placeholder="1.0.0"
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { version: e.target.value }))}
                      />
                    </Field>
                    <Field label="分组 Group">
                      <TextInput
                        value={String(target.group ?? '')}
                        placeholder="（可选）"
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { group: e.target.value }))}
                      />
                    </Field>
                    <Field label="集群 Cluster">
                      <TextInput
                        value={String(target.cluster ?? '')}
                        placeholder="user-dubbo"
                        onChange={(e) => handleFormChange(patchSpec(model, 'target', { cluster: e.target.value }))}
                      />
                    </Field>
                  </div>
                </section>

                <section>
                  <div className="mb-2 flex items-center justify-between">
                    <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400">参数映射</h3>
                    <Button className="px-2 py-1 text-xs" onClick={() => setParams([...params, { from: '', to: params.length, type: 'string' }])}>
                      <Plus size={13} />
                      添加映射
                    </Button>
                  </div>
                  {params.length === 0 ? (
                    <p className="rounded-md bg-slate-50 px-3 py-2.5 text-xs text-slate-500">
                      无参数映射。来源支持 uri./queryStrings./headers./requestBody. 前缀，目标为从 0 开始连续的 Dubbo 参数索引。
                    </p>
                  ) : (
                    <div className="overflow-hidden rounded-md ring-1 ring-slate-200">
                      <table className="w-full text-sm">
                        <thead>
                          <tr className="bg-slate-50 text-left text-xs text-slate-500">
                            <th className="px-3 py-2 font-medium">来源 from</th>
                            <th className="w-28 px-3 py-2 font-medium">目标索引 to</th>
                            <th className="w-44 px-3 py-2 font-medium">类型 type</th>
                            <th className="w-16 px-3 py-2 text-right font-medium">操作</th>
                          </tr>
                        </thead>
                        <tbody>
                          {params.map((row, index) => (
                            <tr key={index} className="border-t border-slate-100">
                              <td className="px-2 py-1.5">
                                <TextInput
                                  value={row.from}
                                  placeholder="uri.id"
                                  className="border-transparent font-mono text-xs focus:border-blue-500"
                                  onChange={(e) => setParams(params.map((r, i) => (i === index ? { ...r, from: e.target.value } : r)))}
                                />
                              </td>
                              <td className="px-2 py-1.5">
                                <TextInput
                                  type="number"
                                  min={0}
                                  value={row.to}
                                  className="border-transparent text-xs focus:border-blue-500"
                                  onChange={(e) => setParams(params.map((r, i) => (i === index ? { ...r, to: e.target.value } : r)))}
                                />
                              </td>
                              <td className="px-2 py-1.5">
                                <Select
                                  value={row.type}
                                  className="border-transparent text-xs focus:border-blue-500"
                                  onChange={(e) => setParams(params.map((r, i) => (i === index ? { ...r, type: e.target.value } : r)))}
                                >
                                  {FALLBACK_PARAM_TYPES.map((item) => (
                                    <option key={item} value={item}>
                                      {item}
                                    </option>
                                  ))}
                                  {!FALLBACK_PARAM_TYPES.includes(row.type) ? <option value={row.type}>{row.type}</option> : null}
                                </Select>
                              </td>
                              <td className="px-2 py-1.5 text-right">
                                <Button variant="ghost" className="px-1.5 text-rose-600 hover:bg-rose-50" onClick={() => setParams(params.filter((_, i) => i !== index))}>
                                  <Trash size={14} />
                                </Button>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </section>

                <section>
                  <h3 className="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-400">发布设置</h3>
                  <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <Field label="模式" hint="保存草稿时由服务端置为 draft">
                      <Select
                        value={String(publish.mode ?? 'draft')}
                        onChange={(e) => handleFormChange(patchSpec(model, 'publish', { mode: e.target.value }))}
                      >
                        <option value="draft">draft</option>
                        <option value="published">published</option>
                      </Select>
                    </Field>
                    <Field label="发布前校验">
                      <Select
                        value={String(publish.validate ?? true)}
                        onChange={(e) => handleFormChange(patchSpec(model, 'publish', { validate: e.target.value === 'true' }))}
                      >
                        <option value="true">开启</option>
                        <option value="false">关闭</option>
                      </Select>
                    </Field>
                  </div>
                  {hasExtensions ? (
                    <p className="mt-2 rounded-md bg-slate-50 px-3 py-2 text-xs text-slate-500">
                      检测到 extensions 扩展字段，表单不展示但会原样保留；可在 YAML 视图中编辑。
                    </p>
                  ) : null}
                </section>
              </div>
            ) : (
              <YamlEditor value={text} onChange={handleTextChange} height={520} />
            )}
          </Panel>

          {preview ? (
            <Panel
              title={
                <div className="flex items-center gap-2">
                  <GitDiff size={16} className="text-slate-500" />
                  <span>Projection 预览</span>
                  {preview.resourceId ? <Tag>Resource {preview.resourceId}</Tag> : null}
                  {preview.methodId ? <Tag>Method {preview.methodId}</Tag> : null}
                  {preview.publishedRevision ? <Tag tone="green">当前发布 r{preview.publishedRevision}</Tag> : <Tag tone="blue">首次发布</Tag>}
                </div>
              }
            >
              {preview.diff ? (
                <div className="space-y-3">
                  <div>
                    <h4 className="mb-1.5 text-xs font-medium text-slate-500">与已发布配置的 Diff</h4>
                    <DiffView diff={preview.diff} />
                  </div>
                  <div>
                    <h4 className="mb-1.5 text-xs font-medium text-slate-500">候选 Legacy 配置</h4>
                    <YamlEditor value={preview.legacyYaml} readOnly height={240} />
                  </div>
                </div>
              ) : (
                <div>
                  <h4 className="mb-1.5 text-xs font-medium text-slate-500">
                    {preview.publishedLegacyYaml ? '与已发布配置一致，无 Diff；候选 Legacy 配置如下' : '编译后的 Legacy 配置（发布时写入运行时）'}
                  </h4>
                  <YamlEditor value={preview.legacyYaml} readOnly height={280} />
                </div>
              )}
            </Panel>
          ) : null}
        </div>
      </div>

      <ConfirmModal
        open={publishOpen}
        title="发布路由绑定"
        loading={publishing}
        confirmText="发布"
        message={
          <>
            将把 <span className="font-mono text-slate-800">{currentName}</span> 的已存草稿投影到旧版 Resource/Method
            键并记录一条发布历史。运行中的 Pixiu 会立即读取新配置。
          </>
        }
        onConfirm={() => void handlePublish()}
        onClose={() => setPublishOpen(false)}
      />

      <ConfirmModal
        open={deleteOpen}
        title="删除草稿"
        danger
        loading={deleting}
        confirmText="删除"
        message={
          <>
            确定删除 <span className="font-mono text-slate-800">{currentName}</span> 的草稿吗？已发布的投影会保留，
            未发布的修改将丢失。
          </>
        }
        onConfirm={() => void handleDeleteDraft()}
        onClose={() => setDeleteOpen(false)}
      />

      <Modal open={historyOpen} title={`发布历史：${currentName}`} width="lg" onClose={() => setHistoryOpen(false)}>
        {historyLoading ? (
          <Loading />
        ) : history.length === 0 ? (
          <EmptyState text="暂无发布历史。" />
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100 text-left text-xs text-slate-500">
                <th className="px-2 py-2 font-medium">Revision</th>
                <th className="px-2 py-2 font-medium">Resource / Method</th>
                <th className="px-2 py-2 font-medium">发布时间</th>
                <th className="px-2 py-2 text-right font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              {history.map((record) => (
                <tr key={record.revision} className="border-b border-slate-50 last:border-0">
                  <td className="px-2 py-2.5 font-mono text-xs text-slate-700">r{record.revision}</td>
                  <td className="px-2 py-2.5 font-mono text-xs text-slate-500">
                    {record.resourceId} / {record.methodId}
                  </td>
                  <td className="px-2 py-2.5 text-slate-600">{formatTime(record.updatedAt)}</td>
                  <td className="px-2 py-2.5 text-right">
                    <Button className="px-2 py-1 text-xs" onClick={() => setRollbackTarget(record)}>
                      回滚到此版本
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Modal>

      <ConfirmModal
        open={rollbackTarget !== null}
        title="回滚路由绑定"
        loading={rollingBack}
        confirmText="回滚"
        message={
          <>
            将把 <span className="font-mono text-slate-800">{currentName}</span> 回滚到 revision {rollbackTarget?.revision}
            的运行时快照，并生成一条新的发布记录。
          </>
        }
        onConfirm={() => void handleRollback()}
        onClose={() => setRollbackTarget(null)}
      />
    </div>
  )
}
