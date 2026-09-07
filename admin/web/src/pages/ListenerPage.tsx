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
import { PencilSimple, Plus, Trash } from '@phosphor-icons/react'
import { createListener, deleteListener, getListenerDetail, getListenerList, updateListener } from '../api'
import { Listener } from '../types'
import { Button, EmptyState, ErrorState, Loading, Panel, Tag } from '../components/common'
import { ConfirmModal, Modal } from '../components/Modal'
import { YamlEditor } from '../components/YamlEditor'
import { useToast } from '../components/Toast'

const LISTENER_TEMPLATE = `name: default-http
address:
  socket_address:
    address: 0.0.0.0
    port: 8888
route_config:
  routes:
    - match:
        prefix: /
      route:
        cluster: test-dubbo
        cluster_not_found_response_code: 505
http_filters:
  - name: dgp.filter.httpconnectionmanager
    config: {}
`

type EditorMode = { kind: 'create' } | { kind: 'edit'; name: string } | null

function listenerAddress(item: Listener): string {
  const socket = item.address?.socket_address
  if (!socket?.address) return '-'
  return `${socket.address}${socket.port ? `:${socket.port}` : ''}`
}

export function ListenerPage() {
  const toast = useToast()
  const [items, setItems] = useState<Listener[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [editorMode, setEditorMode] = useState<EditorMode>(null)
  const [editorYaml, setEditorYaml] = useState('')
  const [editorLoading, setEditorLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  const [deleteTarget, setDeleteTarget] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      setItems((await getListenerList()) ?? [])
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载 Listener 失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const openEdit = async (name: string) => {
    setEditorMode({ kind: 'edit', name })
    setEditorYaml('')
    setEditorLoading(true)
    try {
      setEditorYaml(await getListenerDetail(name))
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '加载 Listener 详情失败')
      setEditorMode(null)
    } finally {
      setEditorLoading(false)
    }
  }

  const handleSave = async () => {
    if (!editorMode) return
    setSaving(true)
    try {
      if (editorMode.kind === 'create') {
        await createListener(editorYaml)
        toast.success('Listener 已创建')
      } else {
        await updateListener(editorYaml)
        toast.success('Listener 已保存')
      }
      setEditorMode(null)
      void load()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    setDeleting(true)
    try {
      await deleteListener(deleteTarget)
      toast.success(`Listener ${deleteTarget} 已删除`)
      setDeleteTarget(null)
      void load()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '删除失败')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="mx-auto max-w-[1400px]">
      <Panel
        title="Listener"
        extra={
          <div className="flex items-center gap-2">
            <Button onClick={() => void load()}>刷新</Button>
            <Button
              variant="primary"
              onClick={() => {
                setEditorMode({ kind: 'create' })
                setEditorYaml(LISTENER_TEMPLATE)
              }}
            >
              <Plus size={15} />
              新增 Listener
            </Button>
          </div>
        }
        bodyClassName="p-0"
      >
        {loading ? (
          <Loading />
        ) : error ? (
          <ErrorState text={error} onRetry={load} />
        ) : items.length === 0 ? (
          <EmptyState text="暂无 Listener，点击右上角新增。" />
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100 text-left text-xs text-slate-500">
                <th className="px-4 py-2.5 font-medium">名称</th>
                <th className="px-4 py-2.5 font-medium">监听地址</th>
                <th className="px-4 py-2.5 text-right font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item, index) => (
                <tr key={item.name ?? index} className="border-b border-slate-50 last:border-0 hover:bg-slate-50/60">
                  <td className="px-4 py-2.5">
                    <Tag tone="blue">{item.name || '-'}</Tag>
                  </td>
                  <td className="px-4 py-2.5 font-mono text-xs text-slate-600">{listenerAddress(item)}</td>
                  <td className="px-4 py-2.5">
                    <div className="flex items-center justify-end gap-1">
                      <Button variant="ghost" className="px-2" onClick={() => item.name && void openEdit(item.name)}>
                        <PencilSimple size={15} />
                      </Button>
                      <Button
                        variant="ghost"
                        className="px-2 text-rose-600 hover:bg-rose-50"
                        onClick={() => item.name && setDeleteTarget(item.name)}
                      >
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
        open={editorMode !== null}
        title={editorMode?.kind === 'create' ? '新增 Listener' : `编辑 Listener：${editorMode?.kind === 'edit' ? editorMode.name : ''}`}
        width="lg"
        onClose={() => setEditorMode(null)}
        footer={
          <>
            <Button onClick={() => setEditorMode(null)}>取消</Button>
            <Button variant="primary" loading={saving} disabled={editorLoading} onClick={() => void handleSave()}>
              保存
            </Button>
          </>
        }
      >
        <p className="mb-2 text-xs text-slate-500">以 YAML 描述 Listener 的监听地址、路由与过滤器。</p>
        {editorLoading ? <Loading /> : <YamlEditor value={editorYaml} onChange={setEditorYaml} height={380} />}
      </Modal>

      <ConfirmModal
        open={deleteTarget !== null}
        title="删除 Listener"
        danger
        loading={deleting}
        confirmText="删除"
        message={<>确定删除 Listener {deleteTarget} 吗？</>}
        onConfirm={() => void handleDelete()}
        onClose={() => setDeleteTarget(null)}
      />
    </div>
  )
}
