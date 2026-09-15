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

import React, { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Boxes,
  Code2,
  Database,
  FileCode2,
  GitBranch,
  Network,
  Plus,
  RefreshCw,
  Save,
  Trash2,
  X,
} from 'lucide-react'
import { clusterApi } from '../../services/cluster-api'
import { listenerApi } from '../../services/listener-api'
import { pluginGroupApi } from '../../services/plugin-group-api'
import { rateLimitApi } from '../../services/rate-limit-api'
import { opaApi } from '../../services/opa-api'
import { baseApi } from '../../services/base-api'
import { Locale, translateText } from '../../i18n'
import { RouteBindingEditor } from './RouteBindingEditor'
import { routeBindingApi } from '../../services/route-binding-api'
import { asJsonObject, JsonObject } from '../../types/api'
import type { RouteBinding, RouteBindingPublishStatus } from '../../types/api'

type Kind = 'cluster' | 'listener' | 'plugin'
const templates = {
  cluster: 'name: local-cluster\ntype: Static\naddress: 127.0.0.1\nport: 20880\n',
  listener:
    'name: http-listener\naddress:\n  socket-address:\n    address: 0.0.0.0\n    port: 8888\nroute_config:\n  routes: []\n',
  plugin:
    'groupName: group1\nplugins:\n  - name: rate limit\n    version: 0.0.1\n    priority: 1000\n    externalLookupName: ExternalPluginRateLimit\n',
}

function Editor({
  title,
  value,
  onChange,
  onClose,
  onSave,
  saving,
  locale,
}: {
  title: string
  value: string
  onChange: (v: string) => void
  onClose: () => void
  onSave: () => void
  saving: boolean
  locale: Locale
}) {
  const tx = (value: string) => translateText(locale, value)
  return (
    <div className="drawer-backdrop" onClick={onClose}>
      <aside className="drawer" onClick={(e) => e.stopPropagation()}>
        <div className="drawer-head">
          <div>
            <span className="eyebrow">YAML CONFIGURATION</span>
            <h2>{title}</h2>
          </div>
          <button className="icon-btn" onClick={onClose}>
            <X size={18} />
          </button>
        </div>
        <div className="code-preview" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          <div className="code-top">
            <span>
              <FileCode2 size={15} /> config.yaml
            </span>
            <span>{tx('由 Pixiu Admin 校验')}</span>
          </div>
          <textarea
            className="yaml-editor"
            value={value}
            onChange={(e) => onChange(e.target.value)}
            spellCheck={false}
          />
        </div>
        <div className="drawer-foot">
          <button className="secondary" onClick={onClose}>
            {tx('取消')}
          </button>
          <button className="primary" disabled={saving} onClick={onSave}>
            <Save size={15} />
            {saving ? tx('保存中…') : tx('保存')}
          </button>
        </div>
      </aside>
    </div>
  )
}

function ErrorState({
  message,
  retry,
  locale,
}: {
  message: string
  retry: () => void
  locale: Locale
}) {
  return (
    <div className="empty">
      <span>{message}</span>
      <button className="secondary" onClick={retry}>
        {translateText(locale, '重试')}
      </button>
    </div>
  )
}

type ConfigFilter = 'All' | 'Connected' | 'Draft'
type ConfigRow = {
  item: JsonObject
  id: string
  title: string
  subtitle: string
  kindLabel: string
  detail: string
  status: ConfigFilter
}
function configIdentity(kind: Kind, item: JsonObject) {
  return kind === 'listener'
    ? String(item.name || '')
    : kind === 'plugin'
      ? String(item.name || item.groupName || '')
      : String(item.id ?? item.name ?? '')
}
function pluginNames(item: JsonObject) {
  const plugins = Array.isArray(item.plugins) ? item.plugins : []
  const names = plugins
    .map((plugin: unknown) => {
      const value = asJsonObject(plugin).name ?? asJsonObject(plugin).externalLookupName
      return typeof value === 'string' ? value : ''
    })
    .filter(Boolean)
  if (names.length) return names
  return [...String(item.content || '').matchAll(/^\s*-\s+name:\s*['"]?([^'"\n]+)['"]?\s*$/gm)]
    .map((match) => match[1].trim())
    .filter(Boolean)
}
function configRow(kind: Kind, item: JsonObject, index: number, locale: Locale): ConfigRow {
  const id = configIdentity(kind, item) || `item-${index}`
  if (kind === 'cluster') {
    const endpoint =
      item.address && item.port
        ? `${item.address}:${item.port}`
        : translateText(locale, '未配置端点')
    return {
      item,
      id,
      title: String(item.name || `cluster.${id}`),
      subtitle: `#${id}`,
      kindLabel: String(item.type || 'Static'),
      detail: endpoint,
      status: item.address && item.port ? 'Connected' : 'Draft',
    }
  }
  if (kind === 'listener') {
    const address = asJsonObject(item.address)
    const socket = asJsonObject(address['socket-address'] ?? address.socket_address)
    const endpoint =
      socket.address && socket.port
        ? `${socket.address}:${socket.port}`
        : translateText(locale, '未配置监听地址')
    const routeConfig = asJsonObject(item.route_config)
    const routeCount = Array.isArray(routeConfig.routes) ? routeConfig.routes.length : 0
    return {
      item,
      id,
      title: String(item.name || `listener.${id}`),
      subtitle: endpoint,
      kindLabel: 'HTTP',
      detail: locale === 'en-US' ? `${routeCount} routes` : `${routeCount} 条路由`,
      status: socket.address && socket.port ? 'Connected' : 'Draft',
    }
  }
  const names = pluginNames(item)
  const preview = names.slice(0, 2).join(' · ') || translateText(locale, '暂无插件')
  return {
    item,
    id,
    title: String(item.name || item.groupName || `group.${id}`),
    subtitle: locale === 'en-US' ? `${names.length} plugins` : `${names.length} 个插件`,
    kindLabel: 'PLUGIN GROUP',
    detail: names.length > 2 ? `${preview} +${names.length - 2}` : preview,
    status: names.length || String(item.content || '').trim() ? 'Connected' : 'Draft',
  }
}

function ConfigListPage({ kind, locale }: { kind: Kind; locale: Locale }) {
  const spec =
    kind === 'cluster'
      ? { title: '集群', desc: '管理 Dubbo 集群、注册中心和服务端点。', icon: Boxes }
      : kind === 'listener'
        ? { title: '监听器', desc: '管理监听地址、端口、协议和路由匹配。', icon: Network }
        : { title: '插件组', desc: '管理网关过滤器插件组及其执行顺序。', icon: Code2 }
  const tx = (value: string) => translateText(locale, value)
  const title = tx(spec.title)
  const noun =
    locale === 'en-US'
      ? kind === 'cluster'
        ? 'cluster'
        : kind === 'listener'
          ? 'listener'
          : 'plugin group'
      : spec.title
  const [items, setItems] = useState<JsonObject[]>([])
  const [query, setQuery] = useState('')
  const [status, setStatus] = useState<ConfigFilter>('All')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [editor, setEditor] = useState<{
    mode: 'create' | 'edit'
    id: string
    value: string
  } | null>(null)
  const [saving, setSaving] = useState(false)
  const Icon = spec.icon
  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data =
        kind === 'cluster'
          ? await clusterApi.list()
          : kind === 'listener'
            ? await listenerApi.list()
            : await pluginGroupApi.list()
      setItems(Array.isArray(data) ? data : [])
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('加载失败'))
    } finally {
      setLoading(false)
    }
  }, [kind, locale])
  useEffect(() => {
    void load()
  }, [load])
  const openCreate = () => setEditor({ mode: 'create', id: '', value: templates[kind] })
  const openEdit = async (item: JsonObject) => {
    const id = configIdentity(kind, item)
    setEditor({ mode: 'edit', id, value: '' })
    try {
      const value =
        kind === 'cluster'
          ? await clusterApi.detail(id)
          : kind === 'listener'
            ? await listenerApi.detail(id)
            : await pluginGroupApi.detail(id)
      setEditor({ mode: 'edit', id, value: typeof value === 'string' ? value : '' })
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('加载详情失败'))
      setEditor(null)
    }
  }
  const save = async () => {
    if (!editor) return
    setSaving(true)
    try {
      if (kind === 'cluster')
        await clusterApi.save(editor.value, editor.mode === 'create' ? 'PUT' : 'POST')
      else if (kind === 'listener')
        await listenerApi.save(editor.value, editor.mode === 'create' ? 'PUT' : 'POST')
      else await pluginGroupApi.save(editor.value, editor.mode === 'create' ? 'POST' : 'PUT')
      setEditor(null)
      await load()
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('保存失败'))
    } finally {
      setSaving(false)
    }
  }
  const remove = async (item: JsonObject) => {
    const id = configIdentity(kind, item)
    if (
      !id ||
      !window.confirm(
        locale === 'en-US' ? `Delete ${title} “${id}”?` : `确认删除${spec.title}“${id}”吗？`,
      )
    )
      return
    try {
      if (kind === 'cluster') await clusterApi.remove(id)
      else if (kind === 'listener') await listenerApi.remove(id)
      else await pluginGroupApi.remove(id)
      await load()
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('删除失败'))
    }
  }
  const rows = useMemo(
    () => items.map((item, index) => configRow(kind, item, index, locale)),
    [items, kind, locale],
  )
  const counts = useMemo(
    () => ({
      All: rows.length,
      Connected: rows.filter((row) => row.status === 'Connected').length,
      Draft: rows.filter((row) => row.status === 'Draft').length,
    }),
    [rows],
  )
  const filtered = useMemo(
    () =>
      rows.filter(
        (row) =>
          (status === 'All' || row.status === status) &&
          `${row.title} ${row.subtitle} ${row.kindLabel} ${row.detail}`
            .toLowerCase()
            .includes(query.toLowerCase()),
      ),
    [query, rows, status],
  )
  const statusLabel = (value: ConfigFilter) =>
    value === 'All' ? tx('全部') : value === 'Connected' ? tx('已连接') : tx('待配置')
  const actionTitle = (mode: 'create' | 'edit') =>
    locale === 'en-US'
      ? `${mode === 'create' ? 'Create' : 'Edit'} ${noun}`
      : `${mode === 'create' ? '新建' : '编辑'}${spec.title}`
  return (
    <div className="page-resource data-prototype">
      <div className="placeholder-head">
        <div>
          <p className="eyebrow">CONNECTED API</p>
          <h1>{title}</h1>
          <p className="muted">{tx(spec.desc)}</p>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="secondary" onClick={() => void load()}>
            <RefreshCw size={14} />
            {tx('刷新')}
          </button>
          <button className="primary" onClick={openCreate}>
            <Plus size={15} />
            {actionTitle('create')}
          </button>
        </div>
      </div>
      <div className="panel route-panel">
        <div className="toolbar">
          <label className="search">
            <Icon size={18} />
            <input
              aria-label={`${locale === 'en-US' ? 'Search ' : '搜索'}${title}`}
              placeholder={`${locale === 'en-US' ? 'Search ' : '搜索'}${title}${locale === 'en-US' ? ', config, or address' : '、配置或地址'}`}
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </label>
          <div className="filters">
            {(['All', 'Connected', 'Draft'] as const).map((value) => (
              <button
                key={value}
                className={`filter ${status === value ? 'active' : ''}`}
                onClick={() => setStatus(value)}
              >
                {statusLabel(value)} <span>{counts[value]}</span>
              </button>
            ))}
          </div>
        </div>
        {loading ? (
          <div className="loading-skeleton" />
        ) : error ? (
          <ErrorState message={error} retry={() => void load()} locale={locale} />
        ) : rows.length === 0 ? (
          <div className="empty">
            <Icon size={27} />
            <b>
              {tx('暂无')}
              {title}
            </b>
            <span>{tx('点击右上角新建配置。')}</span>
            <button className="primary" onClick={openCreate}>
              {actionTitle('create')}
            </button>
          </div>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>
                    {kind === 'cluster'
                      ? 'Cluster'
                      : kind === 'listener'
                        ? 'Listener'
                        : 'Plugin group'}
                  </th>
                  <th>
                    {kind === 'cluster'
                      ? tx('类型')
                      : kind === 'listener'
                        ? tx('协议')
                        : tx('插件')}
                  </th>
                  <th>
                    {kind === 'cluster'
                      ? tx('服务端点')
                      : kind === 'listener'
                        ? tx('路由')
                        : tx('执行配置')}
                  </th>
                  <th>Status</th>
                  <th>Updated</th>
                  <th aria-label={tx('操作')} />
                </tr>
              </thead>
              <tbody>
                {filtered.map((row) => (
                  <tr key={row.id}>
                    <td>
                      <b>{row.title}</b>
                      <small>{row.subtitle}</small>
                    </td>
                    <td>
                      <span className="mono">{row.kindLabel}</span>
                    </td>
                    <td>
                      <span className="target">
                        <Database size={18} />
                        {row.detail}
                      </span>
                    </td>
                    <td>
                      <span className={`badge ${row.status.toLowerCase()}`}>
                        <i />
                        {statusLabel(row.status)}
                      </span>
                    </td>
                    <td>
                      <span className="muted">—</span>
                    </td>
                    <td>
                      <div className="row-actions">
                        <button className="link-btn" onClick={() => void openEdit(row.item)}>
                          {tx('编辑')}
                        </button>
                        <button
                          className="link-btn danger-link"
                          onClick={() => void remove(row.item)}
                        >
                          <Trash2 size={13} />
                          {tx('删除')}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filtered.length === 0 && (
              <div className="empty compact-empty">
                {locale === 'en-US'
                  ? `No matching ${title.toLowerCase()}`
                  : `没有匹配的${spec.title}`}
              </div>
            )}
          </div>
        )}
      </div>
      {editor && (
        <Editor
          locale={locale}
          title={actionTitle(editor.mode)}
          value={editor.value}
          onChange={(value) => setEditor({ ...editor, value })}
          onClose={() => setEditor(null)}
          onSave={() => void save()}
          saving={saving}
        />
      )}
    </div>
  )
}

export function ClusterConfigPage({ locale }: { locale: Locale }) {
  return <ConfigListPage kind="cluster" locale={locale} />
}
export function ListenerConfigPage({ locale }: { locale: Locale }) {
  return <ConfigListPage kind="listener" locale={locale} />
}
export function PluginGroupConfigPage({ locale }: { locale: Locale }) {
  return <ConfigListPage kind="plugin" locale={locale} />
}

type RouteStatus = 'Published' | 'Draft'
type RouteRow = {
  binding: RouteBinding
  id: string
  name: string
  path: string
  verb: string
  target: string
  status: RouteStatus
  publishedRevision: number
}

function routeObjectsMatch(left: RouteBinding, right: RouteBinding) {
  return JSON.stringify(left.object) === JSON.stringify(right.object)
}

function routeRow(
  binding: RouteBinding,
  published: Map<string, RouteBinding>,
  locale: Locale,
): RouteRow {
  const object = binding.object
  const spec = asJsonObject(object.spec)
  const entry = asJsonObject(spec.entry)
  const target = asJsonObject(spec.target)
  const name = String(object.metadata?.name || `route.${binding.resourceId}`)
  const publishedBinding = published.get(name)
  const targetLabel = [target.application, target.interface]
    .filter((value) => typeof value === 'string' && value.trim())
    .join(' / ')

  return {
    binding,
    id: name,
    name,
    path: String(entry.path || '-'),
    verb: String(entry.method || 'GET').toUpperCase(),
    target:
      targetLabel ||
      (typeof target.cluster === 'string' && target.cluster.trim()
        ? target.cluster
        : translateText(locale, '未配置后端目标')),
    status:
      publishedBinding && routeObjectsMatch(binding, publishedBinding) ? 'Published' : 'Draft',
    publishedRevision: publishedBinding?.revision || 0,
  }
}

export function ResourcePage({
  onCountChange,
  locale,
}: {
  onCountChange?: (count: number) => void
  locale: Locale
}) {
  const tx = (value: string) => translateText(locale, value)
  const title = locale === 'en-US' ? 'API routes' : 'API 路由'
  const [rows, setRows] = useState<RouteRow[]>([])
  const [query, setQuery] = useState('')
  const [status, setStatus] = useState<'All' | RouteStatus>('All')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [editor, setEditor] = useState<{
    mode: 'create' | 'edit'
    name: string
    binding: RouteBinding | null
    loading: boolean
    published: boolean
    publishStatus: RouteBindingPublishStatus | null
  } | null>(null)
  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const [draftBindings, publishedBindings] = await Promise.all([
        routeBindingApi.list('draft'),
        routeBindingApi.list('published'),
      ])
      const published = new Map(
        publishedBindings.map((binding) => [binding.object.metadata.name, binding]),
      )
      setRows(draftBindings.map((binding) => routeRow(binding, published, locale)))
      onCountChange?.(draftBindings.length)
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('加载路由失败'))
    } finally {
      setLoading(false)
    }
  }, [onCountChange, locale])
  useEffect(() => {
    void load()
  }, [load])
  const counts = useMemo(
    () => ({
      All: rows.length,
      Published: rows.filter((row) => row.status === 'Published').length,
      Draft: rows.filter((row) => row.status === 'Draft').length,
    }),
    [rows],
  )
  const filtered = useMemo(
    () =>
      rows.filter(
        (row) =>
          (status === 'All' || row.status === status) &&
          `${row.name} ${row.path} ${row.target}`.toLowerCase().includes(query.toLowerCase()),
      ),
    [query, rows, status],
  )
  const openCreate = () =>
    setEditor({
      mode: 'create',
      name: '',
      binding: null,
      loading: false,
      published: false,
      publishStatus: null,
    })
  const openEdit = async (row: RouteRow) => {
    const name = row.binding.object.metadata.name
    setEditor({
      mode: 'edit',
      name,
      binding: null,
      loading: true,
      published: row.status === 'Published',
      publishStatus: {
        name,
        draftRevision: row.binding.revision,
        publishedRevision: row.publishedRevision,
        draftExists: true,
        publishedExists: row.status === 'Published',
        dirty: row.status !== 'Published',
      },
    })
    try {
      const [binding, routeStatus] = await Promise.all([
        routeBindingApi.detail(name, 'draft'),
        routeBindingApi.status(name),
      ])
      setEditor({
        mode: 'edit',
        name,
        binding,
        loading: false,
        published: Boolean(routeStatus.publishedExists && !routeStatus.dirty),
        publishStatus: routeStatus,
      })
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('加载路由详情失败'))
      setEditor(null)
    }
  }
  const statusLabel = (value: RouteStatus) => (value === 'Published' ? tx('已发布') : tx('草稿'))
  if (editor) {
    return (
      <RouteBindingEditor
        locale={locale}
        mode={editor.mode}
        binding={editor.binding}
        loading={editor.loading}
        published={editor.published}
        publishStatus={editor.publishStatus}
        onBack={() => {
          setEditor(null)
          void load()
        }}
        onSaved={() => load()}
      />
    )
  }
  return (
    <div className="page-resource route-prototype">
      <div className="placeholder-head">
        <div>
          <p className="eyebrow">CONNECTED API</p>
          <h1>{title}</h1>
          <p className="muted">
            {locale === 'en-US'
              ? 'Manage resource mappings and methods through Pixiu Admin.'
              : '通过 Pixiu Admin 管理资源映射及其方法。'}
          </p>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="secondary" onClick={() => void load()}>
            <RefreshCw size={14} />
            {tx('刷新')}
          </button>
          <button className="primary" onClick={openCreate}>
            <Plus size={15} />
            {locale === 'en-US' ? 'Create API route' : '新建 API 路由'}
          </button>
        </div>
      </div>
      <div className="panel route-panel">
        <div className="toolbar">
          <label className="search">
            <GitBranch size={18} />
            <input
              aria-label={
                locale === 'en-US' ? 'Search routes, paths, or targets' : '搜索路由、路径或目标'
              }
              placeholder={
                locale === 'en-US' ? 'Search routes, paths, or targets' : '搜索路由、路径或目标'
              }
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </label>
          <div className="filters">
            {(['All', 'Published', 'Draft'] as const).map((value) => (
              <button
                key={value}
                className={`filter ${status === value ? 'active' : ''}`}
                onClick={() => setStatus(value)}
              >
                {value === 'All' ? (locale === 'en-US' ? 'All' : '全部') : statusLabel(value)}{' '}
                <span>{counts[value]}</span>
              </button>
            ))}
          </div>
        </div>
        {loading ? (
          <div className="loading-skeleton" />
        ) : error ? (
          <ErrorState message={error} retry={() => void load()} locale={locale} />
        ) : rows.length === 0 ? (
          <div className="empty">
            <GitBranch size={27} />
            <b>{locale === 'en-US' ? 'No API routes' : '暂无 API 路由'}</b>
            <span>
              {locale === 'en-US'
                ? 'No live resources yet. Use the button above to create the first route.'
                : '真实资源列表为空，点击右上角创建第一条路由。'}
            </span>
            <button className="primary" onClick={openCreate}>
              {locale === 'en-US' ? 'Create API route' : '新建 API 路由'}
            </button>
          </div>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Route</th>
                  <th>{tx('请求方法')}</th>
                  <th>{tx('后端目标')}</th>
                  <th>Status</th>
                  <th>Requests</th>
                  <th>Updated</th>
                  <th aria-label={tx('操作')} />
                </tr>
              </thead>
              <tbody>
                {filtered.map((row, index) => (
                  <tr key={`${row.id}-${row.verb}-${index}`}>
                    <td>
                      <b>{row.name}</b>
                      <small>{row.path}</small>
                    </td>
                    <td>
                      <span className={`method ${row.verb.toLowerCase()}`}>{row.verb}</span>
                    </td>
                    <td>
                      <span className="target">
                        <Database size={18} />
                        {row.target}
                      </span>
                    </td>
                    <td>
                      <span className={`badge ${row.status.toLowerCase()}`}>
                        <i />
                        {statusLabel(row.status)}
                      </span>
                    </td>
                    <td>
                      <span className="mono unavailable-value">-</span>
                    </td>
                    <td>
                      <span className="muted">-</span>
                    </td>
                    <td>
                      <div className="row-actions route-row-actions">
                        <button className="link-btn" onClick={() => void openEdit(row)}>
                          {tx('编辑')}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filtered.length === 0 && (
              <div className="empty compact-empty">
                {locale === 'en-US' ? `No matching routes` : '没有匹配的路由'}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export function RateLimitPage({ locale }: { locale: Locale }) {
  const tx = (value: string) => translateText(locale, value)
  const [content, setContent] = useState('')
  const [exists, setExists] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const template =
    'resources:\n  - name: api\n    items:\n      - pattern: /api/*\nrules:\n  - flowRule:\n      resource: api\n      threshold: 100\n      enable: true\n'
  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const value = await rateLimitApi.get()
      setContent(value || '')
      setExists(Boolean(value))
    } catch (e) {
      setContent(template)
      setExists(false)
      setError(
        e instanceof Error
          ? e.message
          : locale === 'en-US'
            ? 'No rate limit configuration'
            : '暂无限流配置',
      )
    } finally {
      setLoading(false)
    }
  }, [locale])
  useEffect(() => {
    void load()
  }, [load])
  const save = async () => {
    setSaving(true)
    try {
      await rateLimitApi.save(content, exists ? 'PUT' : 'POST')
      setExists(true)
      setError('')
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('保存失败'))
    } finally {
      setSaving(false)
    }
  }
  const remove = async () => {
    if (
      !window.confirm(
        locale === 'en-US'
          ? 'Delete the global rate limit configuration?'
          : '确认删除全局限流配置吗？',
      )
    )
      return
    try {
      await rateLimitApi.remove()
      setExists(false)
      setContent(template)
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('删除失败'))
    }
  }
  return (
    <SingleYamlPage
      locale={locale}
      title={locale === 'en-US' ? 'Rate limiting' : '限流'}
      eyebrow="RATE LIMITING"
      desc={
        locale === 'en-US'
          ? 'Manage Pixiu global rate-limit filter configuration.'
          : '管理 Pixiu 全局限流过滤器配置。'
      }
      content={content}
      loading={loading}
      error={error}
      exists={exists}
      saving={saving}
      onChange={setContent}
      onReload={() => void load()}
      onSave={() => void save()}
      onDelete={() => void remove()}
    />
  )
}

export function OpaPage({ locale }: { locale: Locale }) {
  const tx = (value: string) => translateText(locale, value)
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      setContent((await opaApi.get('policy_id=pixiu-authz')) || '')
    } catch (e) {
      setError(
        e instanceof Error
          ? e.message
          : locale === 'en-US'
            ? 'Unable to connect to OPA (127.0.0.1:8181)'
            : '无法连接 OPA 服务（127.0.0.1:8181）',
      )
    } finally {
      setLoading(false)
    }
  }, [locale])
  useEffect(() => {
    void load()
  }, [load])
  const save = async () => {
    setSaving(true)
    try {
      await opaApi.save({ policy_id: 'pixiu-authz', content })
      setError('')
    } catch (e) {
      setError(
        e instanceof Error
          ? e.message
          : locale === 'en-US'
            ? 'Failed to save OPA policy'
            : '保存 OPA 策略失败',
      )
    } finally {
      setSaving(false)
    }
  }
  const remove = async () => {
    if (
      !window.confirm(
        locale === 'en-US' ? 'Delete the pixiu-authz policy?' : '确认删除 pixiu-authz 策略吗？',
      )
    )
      return
    try {
      await opaApi.remove('pixiu-authz')
      setContent('')
    } catch (e) {
      setError(e instanceof Error ? e.message : tx('删除失败'))
    }
  }
  return (
    <SingleYamlPage
      locale={locale}
      title={locale === 'en-US' ? 'OPA policies' : 'OPA 策略'}
      eyebrow="POLICY CONTROL"
      desc={
        locale === 'en-US'
          ? 'Read, update, or delete the OPA policy used by Pixiu Admin.'
          : '读取、更新或删除 Pixiu Admin 使用的 OPA 策略。'
      }
      content={content}
      loading={loading}
      error={error}
      exists={Boolean(content)}
      saving={saving}
      onChange={setContent}
      onReload={() => void load()}
      onSave={() => void save()}
      onDelete={() => void remove()}
    />
  )
}

export function BaseInfoPage({ locale }: { locale: Locale }) {
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      setContent((await baseApi.get()) || '')
    } catch (e) {
      setError(
        e instanceof Error
          ? e.message
          : locale === 'en-US'
            ? 'Unable to read base configuration'
            : '无法读取基础配置',
      )
    } finally {
      setLoading(false)
    }
  }, [locale])
  useEffect(() => {
    void load()
  }, [load])
  const save = async () => {
    setSaving(true)
    try {
      await baseApi.save(content)
      setError('')
    } catch (e) {
      setError(
        e instanceof Error
          ? e.message
          : locale === 'en-US'
            ? 'Failed to save base configuration'
            : '保存基础配置失败',
      )
    } finally {
      setSaving(false)
    }
  }
  return (
    <SingleYamlPage
      locale={locale}
      title={locale === 'en-US' ? 'Base configuration' : '基础配置'}
      eyebrow="GATEWAY SETTINGS"
      desc={
        locale === 'en-US'
          ? 'Read and update Pixiu gateway base information.'
          : '读取和更新 Pixiu 网关的基础信息。'
      }
      content={content}
      loading={loading}
      error={error}
      exists={Boolean(content)}
      saving={saving}
      onChange={setContent}
      onReload={() => void load()}
      onSave={() => void save()}
      onDelete={() => undefined}
      hideDelete
    />
  )
}

function SingleYamlPage({
  locale,
  title,
  eyebrow,
  desc,
  content,
  loading,
  error,
  exists,
  saving,
  onChange,
  onReload,
  onSave,
  onDelete,
  hideDelete,
}: {
  locale: Locale
  title: string
  eyebrow: string
  desc: string
  content: string
  loading: boolean
  error: string
  exists: boolean
  saving: boolean
  onChange: (v: string) => void
  onReload: () => void
  onSave: () => void
  onDelete: () => void
  hideDelete?: boolean
}) {
  const tx = (value: string) => translateText(locale, value)
  return (
    <div className="page-resource">
      <div className="placeholder-head">
        <div>
          <p className="eyebrow">{eyebrow}</p>
          <h1>{title}</h1>
          <p className="muted">{desc}</p>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="secondary" onClick={onReload}>
            <RefreshCw size={14} />
            {tx('刷新')}
          </button>
          {!hideDelete && (
            <button className="secondary" disabled={!exists} onClick={onDelete}>
              <Trash2 size={14} />
              {tx('删除')}
            </button>
          )}
          <button className="primary" disabled={loading || saving} onClick={onSave}>
            <Save size={14} />
            {exists
              ? locale === 'en-US'
                ? 'Save changes'
                : '保存更新'
              : locale === 'en-US'
                ? 'Create configuration'
                : '创建配置'}
          </button>
        </div>
      </div>
      <div className="panel">
        <div className="panel-head">
          <div>
            <h2>{exists ? tx('已连接后端配置') : tx('尚未配置')}</h2>
            <span>{error || tx('配置将通过 Pixiu Admin API 保存。')}</span>
          </div>
          <Code2 size={18} />
        </div>
        {loading ? (
          <div className="loading-skeleton" />
        ) : (
          <textarea
            className="yaml-editor yaml-page"
            value={content}
            onChange={(e) => onChange(e.target.value)}
            spellCheck={false}
          />
        )}
      </div>
    </div>
  )
}
