/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowRight, Plus, Trash } from '@phosphor-icons/react'
import { createResource, deleteResource, getBaseInfo, getResourceList, setBaseInfo } from '../api'
import { Resource } from '../types'
import { Button, EmptyState, ErrorState, formatDuration, Loading, Panel, Tag } from '../components/common'
import { ConfirmModal, Modal } from '../components/Modal'
import { YamlEditor } from '../components/YamlEditor'
import { useToast } from '../components/Toast'

const RESOURCE_TEMPLATE = `type: restful
path: /api/v1/example
timeout: 1s
description: ''
filters: []
methods: []
`

export function OverviewPage() {
  const navigate = useNavigate()
  const toast = useToast()

  const [baseYaml, setBaseYaml] = useState('')
  const [baseLoading, setBaseLoading] = useState(true)
  const [baseError, setBaseError] = useState('')
  const [baseSaving, setBaseSaving] = useState(false)

  const [resources, setResources] = useState<Resource[]>([])
  const [listLoading, setListLoading] = useState(true)
  const [listError, setListError] = useState('')

  const [createOpen, setCreateOpen] = useState(false)
  const [createYaml, setCreateYaml] = useState(RESOURCE_TEMPLATE)
  const [creating, setCreating] = useState(false)

  const [deleteTarget, setDeleteTarget] = useState<Resource | null>(null)
  const [deleting, setDeleting] = useState(false)

  const loadBase = useCallback(async () => {
    setBaseLoading(true)
    setBaseError('')
    try {
      setBaseYaml(await getBaseInfo())
    } catch (err) {
      // Key missing on a fresh etcd is normal; keep the editor writable so
      // the config can be created for the first time.
      setBaseYaml('')
      setBaseError(err instanceof Error ? err.message : '加载基础配置失败')
    } finally {
      setBaseLoading(false)
    }
  }, [])

  const loadResources = useCallback(async () => {
    setListLoading(true)
    setListError('')
    try {
      setResources(await getResourceList())
    } catch (err) {
      setListError(err instanceof Error ? err.message : '加载 Resource 列表失败')
    } finally {
      setListLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadBase()
    void loadResources()
  }, [loadBase, loadResources])

  const handleSaveBase = async () => {
    setBaseSaving(true)
    try {
      await setBaseInfo(baseYaml)
      toast.success('基础配置已保存')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存失败')
    } finally {
      setBaseSaving(false)
    }
  }

  const handleCreate = async () => {
    setCreating(true)
    try {
      await createResource(createYaml)
      toast.success('Resource 已创建')
      setCreateOpen(false)
      setCreateYaml(RESOURCE_TEMPLATE)
      void loadResources()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '创建失败')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget?.id) return
    setDeleting(true)
    try {
      await deleteResource(deleteTarget.id)
      toast.success(`Resource ${deleteTarget.id} 已删除`)
      setDeleteTarget(null)
      void loadResources()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '删除失败')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="mx-auto max-w-[1400px] space-y-4">
      <Panel
        title="基础网关配置"
        extra={
          <Button variant="primary" onClick={handleSaveBase} loading={baseSaving} disabled={baseLoading}>
            保存
          </Button>
        }
      >
        {baseLoading ? (
          <Loading />
        ) : (
          <>
            {baseError ? (
              <p className="mb-3 flex items-center justify-between gap-2 rounded-md bg-amber-50 px-3 py-2 text-xs text-amber-700">
                <span>未能读取现有配置（{baseError}）。可直接编辑并保存以创建。</span>
                <button type="button" className="shrink-0 font-medium underline hover:text-amber-800" onClick={() => void loadBase()}>
                  重试
                </button>
              </p>
            ) : null}
            <YamlEditor value={baseYaml} onChange={setBaseYaml} height={220} />
          </>
        )}
      </Panel>

      <Panel
        title="Resource 列表"
        extra={
          <div className="flex items-center gap-2">
            <Button onClick={() => void loadResources()}>刷新</Button>
            <Button variant="primary" onClick={() => setCreateOpen(true)}>
              <Plus size={15} />
              新增 Resource
            </Button>
          </div>
        }
        bodyClassName="p-0"
      >
        {listLoading ? (
          <Loading />
        ) : listError ? (
          <ErrorState text={listError} onRetry={loadResources} />
        ) : resources.length === 0 ? (
          <EmptyState
            text="暂无 Resource，点击右上角新增。"
            action={
              <Button variant="primary" onClick={() => setCreateOpen(true)}>
                <Plus size={15} />
                新增 Resource
              </Button>
            }
          />
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100 text-left text-xs text-slate-500">
                <th className="px-4 py-2.5 font-medium">ID</th>
                <th className="px-4 py-2.5 font-medium">类型</th>
                <th className="px-4 py-2.5 font-medium">路径</th>
                <th className="px-4 py-2.5 font-medium">超时</th>
                <th className="px-4 py-2.5 font-medium">描述</th>
                <th className="px-4 py-2.5 text-right font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              {resources.map((item) => (
                <tr key={item.id} className="border-b border-slate-50 last:border-0 hover:bg-slate-50/60">
                  <td className="px-4 py-2.5 font-mono text-xs text-slate-600">{item.id}</td>
                  <td className="px-4 py-2.5">
                    <Tag tone="blue">{item.type || '-'}</Tag>
                  </td>
                  <td className="max-w-64 truncate px-4 py-2.5 font-mono text-xs text-slate-700">{item.path}</td>
                  <td className="px-4 py-2.5 text-slate-600">{formatDuration(item.timeout)}</td>
                  <td className="max-w-52 truncate px-4 py-2.5 text-slate-500">{item.description || '-'}</td>
                  <td className="px-4 py-2.5">
                    <div className="flex items-center justify-end gap-1">
                      <Button variant="ghost" className="px-2" onClick={() => navigate(`/gateway/mapping/${item.id}`)}>
                        映射详情
                        <ArrowRight size={14} />
                      </Button>
                      <Button variant="ghost" className="px-2 text-rose-600 hover:bg-rose-50" onClick={() => setDeleteTarget(item)}>
                        <Trash size={15} />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Panel>

      <Modal
        open={createOpen}
        title="新增 Resource"
        width="lg"
        onClose={() => setCreateOpen(false)}
        footer={
          <>
            <Button onClick={() => setCreateOpen(false)}>取消</Button>
            <Button variant="primary" loading={creating} onClick={handleCreate}>
              创建
            </Button>
          </>
        }
      >
        <p className="mb-2 text-xs text-slate-500">以 YAML 描述 Resource，创建后可在列表中进入其方法映射详情。</p>
        <YamlEditor value={createYaml} onChange={setCreateYaml} height={320} />
      </Modal>

      <ConfirmModal
        open={deleteTarget !== null}
        title="删除 Resource"
        danger
        loading={deleting}
        confirmText="删除"
        message={
          <>
            确定删除 Resource <span className="font-mono text-slate-800">{deleteTarget?.path}</span>（ID {deleteTarget?.id}）吗？
            其下的 Method 映射会一并失效。
          </>
        }
        onConfirm={handleDelete}
        onClose={() => setDeleteTarget(null)}
      />
    </div>
  )
}
