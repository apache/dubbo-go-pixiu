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
import { FloppyDisk, Trash } from '@phosphor-icons/react'
import { createRateLimiter, deleteRateLimiter, getRateLimiter, updateRateLimiter } from '../api'
import { Button, ErrorState, Loading, Panel, Tag } from '../components/common'
import { ConfirmModal } from '../components/Modal'
import { YamlEditor } from '../components/YamlEditor'
import { useToast } from '../components/Toast'

const RATELIMIT_TEMPLATE = `rateLimit:
  enable: true
  limit: 100
  window: 1s
`

export function RateLimiterPage() {
  const toast = useToast()
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [exists, setExists] = useState(false)
  const [dirty, setDirty] = useState(false)
  const [saving, setSaving] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [deleting, setDeleting] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data = await getRateLimiter()
      setContent(data || '')
      setExists(!!data)
      setDirty(false)
    } catch (err) {
      // A missing single config is reported as an error by the backend; the
      // page still offers the create flow with a template.
      setExists(false)
      setContent(RATELIMIT_TEMPLATE)
      setError(err instanceof Error ? err.message : '读取限流配置失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      if (exists) {
        await updateRateLimiter(content)
        toast.success('限流配置已更新')
      } else {
        await createRateLimiter(content)
        toast.success('限流配置已创建')
      }
      setExists(true)
      setDirty(false)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    setDeleting(true)
    try {
      await deleteRateLimiter()
      toast.success('限流配置已删除')
      setDeleteOpen(false)
      setExists(false)
      setContent(RATELIMIT_TEMPLATE)
      setDirty(false)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '删除失败')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="mx-auto max-w-[1400px]">
      <Panel
        title={
          <div className="flex items-center gap-2">
            <span>限流配置</span>
            {!loading && !error ? <Tag tone={exists ? 'green' : 'slate'}>{exists ? '已配置' : '未配置'}</Tag> : null}
            {dirty ? <Tag tone="amber">未保存</Tag> : null}
          </div>
        }
        extra={
          <div className="flex items-center gap-2">
            <Button onClick={() => void load()}>刷新</Button>
            <Button variant="danger" disabled={!exists} onClick={() => setDeleteOpen(true)}>
              <Trash size={15} />
              删除
            </Button>
            <Button variant="primary" loading={saving} disabled={loading} onClick={() => void handleSave()}>
              <FloppyDisk size={15} />
              {exists ? '保存更新' : '创建配置'}
            </Button>
          </div>
        }
      >
        {loading ? (
          <Loading />
        ) : (
          <>
            {error ? (
              <p className="mb-3 rounded-md bg-amber-50 px-3 py-2 text-xs text-amber-700">
                未能读取现有配置（{error}）。可直接编辑下方模板并创建。
              </p>
            ) : (
              <p className="mb-3 text-xs text-slate-500">全局限流插件仅有一份配置，保存后立即生效。</p>
            )}
            <YamlEditor
              value={content}
              onChange={(next) => {
                setContent(next)
                setDirty(true)
              }}
              height={420}
            />
          </>
        )}
      </Panel>

      <ConfirmModal
        open={deleteOpen}
        title="删除限流配置"
        danger
        loading={deleting}
        confirmText="删除"
        message="确定删除当前的限流配置吗？删除后网关将不再应用该限流策略。"
        onConfirm={() => void handleDelete()}
        onClose={() => setDeleteOpen(false)}
      />
    </div>
  )
}
