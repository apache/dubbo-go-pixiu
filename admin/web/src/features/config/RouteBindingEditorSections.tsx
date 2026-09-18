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
import { translateText } from '../../i18n'
import type {
  AdminRouteBindingObject,
  RouteBindingDiff,
  RouteBindingParam,
  RouteBindingPublishStatus,
} from '../../types/api'

export const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']
export const PARAM_TYPES = [
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

export type EditorTab = 'form' | 'preview' | 'diff' | 'yaml'
export type BusyAction = '' | 'save' | 'publish' | 'validate' | 'preview' | 'diff'
export type Notice = { tone: 'success' | 'error'; text: string }
export type IssueMessage = (path: string) => string
export type RouteEntry = AdminRouteBindingObject['spec']['entry']
export type RouteTarget = AdminRouteBindingObject['spec']['target']
type EditableTargetKey = Exclude<keyof RouteTarget, 'protocol'>

type HeaderProps = {
  readonly isEnglish: boolean
  readonly routeLabel: string
  readonly entry: RouteEntry
  readonly publishStatus: RouteBindingPublishStatus | null
  readonly dirty: boolean
  readonly published: boolean
  readonly enabled: boolean
  readonly currentStatus: string
  readonly lifecycleStatus: string
  readonly busy: BusyAction
  readonly yamlError: string
  readonly onBack: () => void
  readonly onValidate: () => void
  readonly onSave: () => void
  readonly onPublish: () => void
}

type TabsProps = {
  readonly isEnglish: boolean
  readonly activeTab: EditorTab
  readonly onChange: (tab: EditorTab) => void
}

type IdentityCardProps = {
  readonly isEnglish: boolean
  readonly mode: 'create' | 'edit'
  readonly entry: RouteEntry
  readonly name: string
  readonly issueMessage: IssueMessage
  readonly onNameChange: (value: string) => void
  readonly onPathChange: (value: string) => void
  readonly onMethodChange: (value: string) => void
}

type TargetCardProps = {
  readonly isEnglish: boolean
  readonly target: RouteTarget
  readonly issueMessage: IssueMessage
  readonly onChange: (key: EditableTargetKey, value: string) => void
}

type ParamCardProps = {
  readonly isEnglish: boolean
  readonly params: RouteBindingParam[]
  readonly issueMessage: IssueMessage
  readonly onAdd: () => void
  readonly onUpdate: (index: number, patch: Partial<RouteBindingParam>) => void
  readonly onRemove: (index: number) => void
}

type PublishCardProps = {
  readonly isEnglish: boolean
  readonly enabled: boolean
  readonly validate: boolean
  readonly onEnabledChange: (enabled: boolean) => void
  readonly onValidateChange: (validate: boolean) => void
}

type FormViewProps = IdentityCardProps & TargetCardProps & ParamCardProps & PublishCardProps

type DiffViewProps = {
  readonly isEnglish: boolean
  readonly diffData: RouteBindingDiff | null
  readonly busy: BusyAction
  readonly onRefresh: () => void
}

type PreviewViewProps = {
  readonly isEnglish: boolean
  readonly activeTab: 'preview' | 'yaml'
  readonly busy: BusyAction
  readonly previewYaml: string
  readonly routeYaml: string
  readonly yamlError: string
  readonly onApplyYaml: () => void
  readonly onChangeYaml: (value: string) => void
  readonly onPreview: () => void
}

type ContentProps = {
  readonly isEnglish: boolean
  readonly loading: boolean
  readonly activeTab: EditorTab
  readonly mode: 'create' | 'edit'
  readonly object: AdminRouteBindingObject
  readonly busy: BusyAction
  readonly diffData: RouteBindingDiff | null
  readonly previewYaml: string
  readonly routeYaml: string
  readonly yamlError: string
  readonly issueMessage: IssueMessage
  readonly onNameChange: (value: string) => void
  readonly onEntryChange: (key: 'path' | 'method', value: string) => void
  readonly onTargetChange: (key: EditableTargetKey, value: string) => void
  readonly onParamAdd: () => void
  readonly onParamUpdate: (index: number, patch: Partial<RouteBindingParam>) => void
  readonly onParamRemove: (index: number) => void
  readonly onEnabledChange: (enabled: boolean) => void
  readonly onValidateChange: (validate: boolean) => void
  readonly onRefreshDiff: () => void
  readonly onApplyYaml: () => void
  readonly onChangeYaml: (value: string) => void
  readonly onPreview: () => void
}

function actionLabel(
  isEnglish: boolean,
  busy: BusyAction,
  action: Exclude<BusyAction, ''>,
  idleLabel: string,
) {
  if (busy !== action) return idleLabel
  const labels: Record<Exclude<BusyAction, ''>, [string, string]> = {
    save: ['Saving...', '保存中...'],
    publish: ['Publishing...', '发布中...'],
    validate: ['Checking...', '校验中...'],
    preview: ['Generating...', '生成中...'],
    diff: ['Loading...', '加载中...'],
  }
  return isEnglish ? labels[action][0] : labels[action][1]
}

function routeNameHelp(isEnglish: boolean, mode: 'create' | 'edit') {
  if (mode === 'edit')
    return isEnglish ? 'The key stays stable while editing.' : '编辑时保持路由键稳定。'
  return isEnglish
    ? 'Used as the AdminRouteBinding metadata.name.'
    : '对应 AdminRouteBinding 的 metadata.name。'
}

function previewTitle(isEnglish: boolean, activeTab: 'preview' | 'yaml') {
  if (activeTab === 'yaml') return 'YAML'
  return isEnglish ? 'Runtime preview' : '运行时预览'
}

function previewDescription(isEnglish: boolean, activeTab: 'preview' | 'yaml') {
  if (activeTab === 'yaml') {
    return isEnglish
      ? 'Edit the AdminRouteBinding source. Valid changes sync back to the form.'
      : '编辑 AdminRouteBinding 源配置，语法有效的改动会同步回表单。'
  }
  return isEnglish
    ? 'Generated by the current AdminRouteBinding compiler.'
    : '由当前 AdminRouteBinding 编译器生成。'
}

function routeStateClass(dirty: boolean, published: boolean) {
  if (dirty || !published) return 'draft'
  return 'published'
}

function yamlSyncLabel(isEnglish: boolean, yamlError: string) {
  if (yamlError) return isEnglish ? 'Syntax error' : '语法错误'
  return isEnglish ? 'Synced with form' : '已与表单同步'
}

function yamlSyncClass(yamlError: string) {
  return yamlError ? 'error' : 'synced'
}

export function RouteEditorHeader({
  isEnglish,
  routeLabel,
  entry,
  publishStatus,
  dirty,
  published,
  enabled,
  currentStatus,
  lifecycleStatus,
  busy,
  yamlError,
  onBack,
  onValidate,
  onSave,
  onPublish,
}: HeaderProps) {
  return (
    <>
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
            <span className={`route-editor-state ${routeStateClass(dirty, published)}`}>
              <i /> {currentStatus}
            </span>
            <span className={`route-editor-lifecycle ${enabled ? 'enabled' : 'disabled'}`}>
              <i /> {lifecycleStatus}
            </span>
          </div>
        </div>
        <div className="route-editor-actions">
          <button
            className="secondary"
            type="button"
            disabled={busy !== '' || Boolean(yamlError)}
            onClick={onValidate}
          >
            <CheckCircle2 size={14} />
            {actionLabel(
              isEnglish,
              busy,
              'validate',
              translateText(isEnglish ? 'en-US' : 'zh-CN', '校验'),
            )}
          </button>
          <button
            className="secondary"
            type="button"
            disabled={busy !== '' || Boolean(yamlError)}
            onClick={onSave}
          >
            <Save size={14} />
            {actionLabel(isEnglish, busy, 'save', isEnglish ? 'Save draft' : '保存草稿')}
          </button>
          <button
            className="primary"
            type="button"
            disabled={busy !== '' || Boolean(yamlError)}
            onClick={onPublish}
          >
            <Send size={14} />
            {actionLabel(isEnglish, busy, 'publish', isEnglish ? 'Publish' : '发布')}
          </button>
        </div>
      </div>
    </>
  )
}

type ContractProps = {
  readonly isEnglish: boolean
}

export function RouteEditorContract({ isEnglish }: ContractProps) {
  return (
    <div className="route-editor-contract">
      <span className="route-contract-mark">{isEnglish ? 'ATOMIC PUBLISH' : '原子发布'}</span>
      <span>
        {isEnglish
          ? 'Saving updates this route draft. Publishing replaces only this route in one etcd transaction.'
          : '保存只更新当前路由草稿。发布会通过一次 etcd 事务只替换当前路由。'}
      </span>
    </div>
  )
}

type NoticeProps = {
  readonly notice: Notice
}

export function RouteEditorNotice({ notice }: NoticeProps) {
  if (notice.tone === 'error') {
    return (
      <div className="route-editor-notice error" role="alert">
        <AlertCircle size={16} />
        <span>{notice.text}</span>
      </div>
    )
  }
  return (
    <output className="route-editor-notice success">
      <CheckCircle2 size={16} />
      <span>{notice.text}</span>
    </output>
  )
}

export function RouteEditorTabs({ isEnglish, activeTab, onChange }: TabsProps) {
  return (
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
        onClick={() => onChange('form')}
      >
        {isEnglish ? 'Form' : '表单'}
      </button>
      <button
        className={activeTab === 'preview' ? 'active' : ''}
        type="button"
        role="tab"
        aria-selected={activeTab === 'preview'}
        onClick={() => onChange('preview')}
      >
        <Eye size={14} /> {isEnglish ? 'Preview' : '预览'}
      </button>
      <button
        className={activeTab === 'diff' ? 'active' : ''}
        type="button"
        role="tab"
        aria-selected={activeTab === 'diff'}
        onClick={() => onChange('diff')}
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
        onClick={() => onChange('yaml')}
      >
        <FileCode2 size={14} /> YAML
      </button>
    </div>
  )
}

function RouteIdentityCard({
  isEnglish,
  mode,
  entry,
  name,
  issueMessage,
  onNameChange,
  onPathChange,
  onMethodChange,
}: IdentityCardProps) {
  const nameIssue = issueMessage('metadata.name')
  const pathIssue = issueMessage('spec.entry.path')
  const methodIssue = issueMessage('spec.entry.method')
  return (
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
            value={name}
            readOnly={mode === 'edit'}
            aria-invalid={Boolean(nameIssue)}
            aria-describedby="route-name-help route-name-error"
            placeholder="create-user"
            onChange={(event) => onNameChange(event.target.value)}
          />
          <small id="route-name-help">{routeNameHelp(isEnglish, mode)}</small>
          {nameIssue && (
            <small id="route-name-error" className="route-field-error">
              {nameIssue}
            </small>
          )}
        </label>
        <label className="route-field">
          <span>{isEnglish ? 'Entry protocol' : '入口协议'}</span>
          <select value={entry.protocol} disabled>
            <option value="http">HTTP</option>
          </select>
          <small>
            {isEnglish ? 'The current schema accepts HTTP only.' : '当前 schema 仅支持 HTTP。'}
          </small>
        </label>
        <label className="route-field route-field-wide">
          <span>{isEnglish ? 'Path' : '路径'}</span>
          <input
            value={entry.path}
            aria-invalid={Boolean(pathIssue)}
            placeholder="/api/v1/users/:id"
            onChange={(event) => onPathChange(event.target.value)}
          />
          {pathIssue && <small className="route-field-error">{pathIssue}</small>}
        </label>
        <fieldset className="route-field route-field-wide route-method-fieldset">
          <legend>{isEnglish ? 'HTTP method' : '请求方法'}</legend>
          <div className="route-method-options">
            {HTTP_METHODS.map((method) => (
              <button
                className={entry.method === method ? 'active' : ''}
                key={method}
                type="button"
                aria-pressed={entry.method === method}
                onClick={() => onMethodChange(method)}
              >
                {method}
              </button>
            ))}
          </div>
          {methodIssue && <small className="route-field-error">{methodIssue}</small>}
        </fieldset>
      </div>
    </section>
  )
}

type TargetFieldProps = {
  readonly label: string
  readonly value: string
  readonly placeholder: string
  readonly issue: string
  readonly onChange: (value: string) => void
}

function TargetField({ label, value, placeholder, issue, onChange }: TargetFieldProps) {
  return (
    <label className="route-field">
      <span>{label}</span>
      <input
        value={value}
        placeholder={placeholder}
        onChange={(event) => onChange(event.target.value)}
      />
      {issue && <small className="route-field-error">{issue}</small>}
    </label>
  )
}

function RouteTargetCard({ isEnglish, target, issueMessage, onChange }: TargetCardProps) {
  const fields: Array<{
    key: EditableTargetKey
    label: string
    placeholder: string
    issuePath: string
  }> = [
    {
      key: 'cluster',
      label: 'Cluster',
      placeholder: 'user-service',
      issuePath: 'spec.target.cluster',
    },
    {
      key: 'application',
      label: 'Application',
      placeholder: 'UserProvider',
      issuePath: 'spec.target.application',
    },
    {
      key: 'interface',
      label: 'Interface',
      placeholder: 'com.example.user.UserService',
      issuePath: 'spec.target.interface',
    },
    { key: 'method', label: 'Method', placeholder: 'GetUser', issuePath: 'spec.target.method' },
    { key: 'group', label: 'Group', placeholder: 'stable', issuePath: '' },
    { key: 'version', label: 'Version', placeholder: '1.0.0', issuePath: '' },
  ]
  return (
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
        {fields.map((field) => (
          <TargetField
            key={field.key}
            label={field.label}
            value={target[field.key]}
            placeholder={field.placeholder}
            issue={field.issuePath ? issueMessage(field.issuePath) : ''}
            onChange={(value) => onChange(field.key, value)}
          />
        ))}
      </div>
    </section>
  )
}

type ParamRowProps = {
  readonly isEnglish: boolean
  readonly param: RouteBindingParam
  readonly index: number
  readonly issueMessage: IssueMessage
  readonly onUpdate: (index: number, patch: Partial<RouteBindingParam>) => void
  readonly onRemove: (index: number) => void
}

function RouteParamRow({
  isEnglish,
  param,
  index,
  issueMessage,
  onUpdate,
  onRemove,
}: ParamRowProps) {
  const fromIssue = issueMessage(`spec.params[${index}].from`)
  const toIssue = issueMessage(`spec.params[${index}].to`)
  const typeIssue = issueMessage(`spec.params[${index}].type`)
  return (
    <div className="route-param-row" key={`${index}-${param.from}`}>
      <label className="route-param-field">
        <span>{isEnglish ? 'HTTP source' : 'HTTP 来源'}</span>
        <input
          value={param.from}
          aria-invalid={Boolean(fromIssue)}
          placeholder="queryStrings.page"
          onChange={(event) => onUpdate(index, { from: event.target.value })}
        />
        {fromIssue && <small className="route-field-error">{fromIssue}</small>}
      </label>
      <label className="route-param-field">
        <span>{isEnglish ? 'Argument index' : '参数位置'}</span>
        <input
          type="number"
          min="0"
          value={param.to}
          aria-invalid={Boolean(toIssue)}
          onChange={(event) => onUpdate(index, { to: Number(event.target.value) })}
        />
        {toIssue && <small className="route-field-error">{toIssue}</small>}
      </label>
      <label className="route-param-field">
        <span>{isEnglish ? 'Dubbo type' : 'Dubbo 类型'}</span>
        <select
          value={param.type}
          onChange={(event) => onUpdate(index, { type: event.target.value })}
        >
          {PARAM_TYPES.map((type) => (
            <option key={type} value={type}>
              {type}
            </option>
          ))}
        </select>
        {typeIssue && <small className="route-field-error">{typeIssue}</small>}
      </label>
      <button
        className="route-icon-action"
        type="button"
        aria-label={isEnglish ? `Remove parameter ${index + 1}` : `删除第 ${index + 1} 个参数`}
        onClick={() => onRemove(index)}
      >
        <Trash2 size={15} />
      </button>
    </div>
  )
}

function RouteParamsCard({
  isEnglish,
  params,
  issueMessage,
  onAdd,
  onUpdate,
  onRemove,
}: ParamCardProps) {
  return (
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
        <button className="secondary route-add-param" type="button" onClick={onAdd}>
          <Plus size={14} /> {isEnglish ? 'Add parameter' : '添加参数'}
        </button>
      </div>
      <div className="route-param-table">
        {params.length === 0 ? (
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
            {params.map((param, index) => (
              <RouteParamRow
                key={`${index}-${param.from}`}
                isEnglish={isEnglish}
                param={param}
                index={index}
                issueMessage={issueMessage}
                onUpdate={onUpdate}
                onRemove={onRemove}
              />
            ))}
          </>
        )}
      </div>
    </section>
  )
}

function RoutePublishCard({
  isEnglish,
  enabled,
  validate,
  onEnabledChange,
  onValidateChange,
}: PublishCardProps) {
  return (
    <section className="panel route-publish-panel">
      <div>
        <h2>{isEnglish ? 'Release and lifecycle' : '发布与生命周期'}</h2>
        <span>
          {isEnglish
            ? 'Both lifecycle changes and validation preferences are published with this route.'
            : '路由启停和校验偏好都会随当前路由一起发布。'}
        </span>
      </div>
      <div className="route-switches">
        <label className="route-switch">
          <input
            type="checkbox"
            checked={enabled}
            aria-label={isEnglish ? 'Enable route' : '启用路由'}
            onChange={(event) => onEnabledChange(event.target.checked)}
          />
          <span className="route-switch-track" />
          <span>{isEnglish ? 'Route enabled' : '启用路由'}</span>
        </label>
        <label className="route-switch">
          <input
            type="checkbox"
            checked={validate}
            aria-label={isEnglish ? 'Validate before publish' : '发布前校验'}
            onChange={(event) => onValidateChange(event.target.checked)}
          />
          <span className="route-switch-track" />
          <span>{isEnglish ? 'Validate before publish' : '发布前校验'}</span>
        </label>
      </div>
    </section>
  )
}

function RouteFormView(props: FormViewProps) {
  const {
    isEnglish,
    mode,
    entry,
    name,
    target,
    params,
    enabled,
    validate,
    issueMessage,
    onNameChange,
    onPathChange,
    onMethodChange,
    onChange,
    onAdd,
    onUpdate,
    onRemove,
    onEnabledChange,
    onValidateChange,
  } = props
  return (
    <>
      <div className="route-editor-layout">
        <RouteIdentityCard
          isEnglish={isEnglish}
          mode={mode}
          entry={entry}
          name={name}
          issueMessage={issueMessage}
          onNameChange={onNameChange}
          onPathChange={onPathChange}
          onMethodChange={onMethodChange}
        />
        <RouteTargetCard
          isEnglish={isEnglish}
          target={target}
          issueMessage={issueMessage}
          onChange={onChange}
        />
      </div>
      <RouteParamsCard
        isEnglish={isEnglish}
        params={params}
        issueMessage={issueMessage}
        onAdd={onAdd}
        onUpdate={onUpdate}
        onRemove={onRemove}
      />
      <RoutePublishCard
        isEnglish={isEnglish}
        enabled={enabled}
        validate={validate}
        onEnabledChange={onEnabledChange}
        onValidateChange={onValidateChange}
      />
    </>
  )
}

function DiffRow({
  isEnglish,
  change,
}: {
  readonly isEnglish: boolean
  readonly change: RouteBindingDiff['changes'][number]
}) {
  return (
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
  )
}

function RouteDiffContent({
  isEnglish,
  diffData,
  busy,
}: Pick<DiffViewProps, 'isEnglish' | 'diffData' | 'busy'>) {
  if (!diffData) {
    return busy === 'diff' ? (
      <div className="route-diff-empty">
        <div className="loading-skeleton" />
      </div>
    ) : (
      <div className="route-diff-empty">
        <GitCompareArrows size={22} />
        <b>{isEnglish ? 'Diff is not loaded' : '尚未加载 Diff'}</b>
      </div>
    )
  }
  if (diffData.changes.length === 0) {
    return (
      <div className="route-diff-empty">
        <CheckCircle2 size={22} />
        <b>{isEnglish ? 'No unpublished changes' : '当前路由没有未发布改动'}</b>
        <span>
          {isEnglish
            ? 'The draft and published versions are identical.'
            : '草稿版本与已发布版本内容一致。'}
        </span>
      </div>
    )
  }
  return (
    <div className="route-diff-list">
      {diffData.changes.map((change) => (
        <DiffRow key={change.path} isEnglish={isEnglish} change={change} />
      ))}
    </div>
  )
}

function RouteDiffView({ isEnglish, diffData, busy, onRefresh }: DiffViewProps) {
  return (
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
        <button className="secondary" type="button" disabled={busy !== ''} onClick={onRefresh}>
          <GitCompareArrows size={14} />
          {actionLabel(
            isEnglish,
            busy,
            'diff',
            translateText(isEnglish ? 'en-US' : 'zh-CN', '刷新 Diff'),
          )}
        </button>
      </div>
      <RouteDiffContent isEnglish={isEnglish} diffData={diffData} busy={busy} />
    </section>
  )
}

function RoutePreviewBody({
  isEnglish,
  activeTab,
  previewYaml,
  routeYaml,
  yamlError,
  onChangeYaml,
}: Pick<
  PreviewViewProps,
  'isEnglish' | 'activeTab' | 'previewYaml' | 'routeYaml' | 'yamlError' | 'onChangeYaml'
>) {
  if (activeTab === 'yaml') {
    return (
      <>
        <textarea
          className="route-preview-code route-yaml-editor"
          value={routeYaml}
          spellCheck={false}
          aria-label={isEnglish ? 'Editable route YAML' : '可编辑路由 YAML'}
          aria-invalid={Boolean(yamlError)}
          aria-describedby={yamlError ? 'route-yaml-error' : undefined}
          onChange={(event) => onChangeYaml(event.target.value)}
        />
        {yamlError && (
          <div className="route-yaml-error" id="route-yaml-error" role="alert">
            <AlertCircle size={15} />
            <span>{isEnglish ? `YAML error: ${yamlError}` : `YAML 错误：${yamlError}`}</span>
          </div>
        )}
      </>
    )
  }
  if (previewYaml) return <pre className="route-preview-code">{previewYaml}</pre>
  return (
    <div className="route-preview-empty">
      <FileCode2 size={22} />
      <b>{isEnglish ? 'Preview is not generated' : '还没有生成预览'}</b>
      <span>
        {isEnglish
          ? 'Generate a preview to inspect the legacy runtime configuration.'
          : '生成预览后查看将写入运行时的 legacy 配置。'}
      </span>
    </div>
  )
}

type PreviewActionProps = Pick<
  PreviewViewProps,
  'isEnglish' | 'activeTab' | 'busy' | 'yamlError' | 'onApplyYaml' | 'onPreview'
>

function RoutePreviewAction({
  isEnglish,
  activeTab,
  busy,
  yamlError,
  onApplyYaml,
  onPreview,
}: PreviewActionProps) {
  if (activeTab === 'yaml') {
    return (
      <button
        className="secondary"
        type="button"
        disabled={busy !== '' || Boolean(yamlError)}
        onClick={onApplyYaml}
      >
        <CheckCircle2 size={14} />
        {isEnglish ? 'Sync to form' : '同步到表单'}
      </button>
    )
  }
  return (
    <button className="secondary" type="button" disabled={busy !== ''} onClick={onPreview}>
      <Eye size={14} />
      {actionLabel(isEnglish, busy, 'preview', isEnglish ? 'Generate preview' : '生成预览')}
    </button>
  )
}

type PreviewMetaProps = Pick<PreviewViewProps, 'isEnglish' | 'activeTab' | 'yamlError'>

function RoutePreviewMeta({ isEnglish, activeTab, yamlError }: PreviewMetaProps) {
  const isYaml = activeTab === 'yaml'
  return (
    <div className="route-preview-meta">
      <FileCode2 size={15} />
      <span>{isYaml ? 'route-binding.yaml' : 'api_config.yaml'}</span>
      {isYaml ? (
        <span className={`route-yaml-sync ${yamlSyncClass(yamlError)}`}>
          {yamlSyncLabel(isEnglish, yamlError)}
        </span>
      ) : (
        <span className="muted">{isEnglish ? 'Read only' : '只读'}</span>
      )}
    </div>
  )
}

function RoutePreviewView({
  isEnglish,
  activeTab,
  busy,
  previewYaml,
  routeYaml,
  yamlError,
  onApplyYaml,
  onChangeYaml,
  onPreview,
}: PreviewViewProps) {
  return (
    <section className="panel route-preview-panel">
      <div className="panel-head">
        <div>
          <h2>{previewTitle(isEnglish, activeTab)}</h2>
          <span>{previewDescription(isEnglish, activeTab)}</span>
        </div>
        <RoutePreviewAction
          isEnglish={isEnglish}
          activeTab={activeTab}
          busy={busy}
          yamlError={yamlError}
          onApplyYaml={onApplyYaml}
          onPreview={onPreview}
        />
      </div>
      <RoutePreviewMeta isEnglish={isEnglish} activeTab={activeTab} yamlError={yamlError} />
      <RoutePreviewBody
        isEnglish={isEnglish}
        activeTab={activeTab}
        previewYaml={previewYaml}
        routeYaml={routeYaml}
        yamlError={yamlError}
        onChangeYaml={onChangeYaml}
      />
    </section>
  )
}

export function RouteEditorContent({
  isEnglish,
  loading,
  activeTab,
  mode,
  object,
  busy,
  diffData,
  previewYaml,
  routeYaml,
  yamlError,
  issueMessage,
  onNameChange,
  onEntryChange,
  onTargetChange,
  onParamAdd,
  onParamUpdate,
  onParamRemove,
  onEnabledChange,
  onValidateChange,
  onRefreshDiff,
  onApplyYaml,
  onChangeYaml,
  onPreview,
}: ContentProps) {
  if (loading) {
    return (
      <div
        className="route-binding-loading"
        aria-label={isEnglish ? 'Loading route' : '正在加载路由'}
      >
        <div className="loading-skeleton" />
        <div className="loading-skeleton" />
      </div>
    )
  }
  if (activeTab === 'form') {
    return (
      <RouteFormView
        isEnglish={isEnglish}
        mode={mode}
        entry={object.spec.entry}
        name={object.metadata.name}
        target={object.spec.target}
        params={object.spec.params}
        enabled={object.spec.enabled}
        validate={object.spec.publish.validate}
        issueMessage={issueMessage}
        onNameChange={onNameChange}
        onPathChange={(value) => onEntryChange('path', value)}
        onMethodChange={(value) => onEntryChange('method', value)}
        onChange={onTargetChange}
        onAdd={onParamAdd}
        onUpdate={onParamUpdate}
        onRemove={onParamRemove}
        onEnabledChange={onEnabledChange}
        onValidateChange={onValidateChange}
      />
    )
  }
  if (activeTab === 'diff') {
    return (
      <RouteDiffView
        isEnglish={isEnglish}
        diffData={diffData}
        busy={busy}
        onRefresh={onRefreshDiff}
      />
    )
  }
  return (
    <RoutePreviewView
      isEnglish={isEnglish}
      activeTab={activeTab}
      busy={busy}
      previewYaml={previewYaml}
      routeYaml={routeYaml}
      yamlError={yamlError}
      onApplyYaml={onApplyYaml}
      onChangeYaml={onChangeYaml}
      onPreview={() => onPreview()}
    />
  )
}

export function formatDiffValue(value: unknown) {
  if (value === null || value === undefined) return '∅'
  if (typeof value === 'string') return value
  if (typeof value === 'bigint') return `${value}n`
  if (typeof value === 'function') return `[Function ${value.name || 'anonymous'}]`
  if (typeof value === 'symbol') return value.toString()
  try {
    const formatted = JSON.stringify(value, null, 2)
    if (formatted !== undefined) return formatted
  } catch {
    return '<unserializable value>'
  }
  return '<unserializable value>'
}
