import { useState } from 'react'
import { ArrowClockwise, CloudArrowDown, FloppyDisk, Trash } from '@phosphor-icons/react'
import { deleteOpaPolicy, getOpaPolicy, putOpaPolicy } from '../api'
import { Button, Field, Panel, Tag, TextInput } from '../components/common'
import { ConfirmModal } from '../components/Modal'
import { YamlEditor } from '../components/YamlEditor'
import { useToast } from '../components/Toast'

const REGO_TEMPLATE = `package pixiu.authz

default allow = false

allow {
  input.headers["x-internal"] == "true"
}
`

export function OpaPage() {
  const toast = useToast()
  const [serverUrl, setServerUrl] = useState('')
  const [policyId, setPolicyId] = useState('')
  const [bearerToken, setBearerToken] = useState('')

  const [content, setContent] = useState(REGO_TEMPLATE)
  const [loaded, setLoaded] = useState(false)
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [syncing, setSyncing] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [deleting, setDeleting] = useState(false)

  const query = {
    server_url: serverUrl || undefined,
    policy_id: policyId || undefined,
    bearer_token: bearerToken || undefined,
  }

  const handleFetch = async () => {
    setLoading(true)
    try {
      const result = await getOpaPolicy(query)
      setContent(result || '')
      setLoaded(true)
      toast.success('已拉取策略')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '拉取策略失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    if (!content.trim()) {
      toast.error('策略内容不能为空')
      return
    }
    setSaving(true)
    try {
      await putOpaPolicy(query, content)
      toast.success('策略已保存')
      setLoaded(true)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleSync = async () => {
    if (!content.trim()) {
      toast.error('策略内容不能为空')
      return
    }
    setSyncing(true)
    try {
      await putOpaPolicy(query, content)
      toast.success('已同步到 OPA 服务')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '同步失败')
    } finally {
      setSyncing(false)
    }
  }

  const handleDelete = async () => {
    setDeleting(true)
    try {
      await deleteOpaPolicy(query)
      toast.success('策略已删除')
      setDeleteOpen(false)
      setLoaded(false)
      setContent(REGO_TEMPLATE)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '删除失败')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="mx-auto max-w-[1400px] space-y-4">
      <Panel title="策略来源">
        <div className="grid grid-cols-1 gap-3 md:grid-cols-[1fr_220px_240px_auto]">
          <Field label="OPA Server URL" hint="留空使用后端默认配置">
            <TextInput value={serverUrl} placeholder="http://127.0.0.1:8181" onChange={(e) => setServerUrl(e.target.value)} />
          </Field>
          <Field label="策略 ID">
            <TextInput value={policyId} placeholder="pixiu/authz" className="font-mono" onChange={(e) => setPolicyId(e.target.value)} />
          </Field>
          <Field label="Bearer Token（可选）">
            <TextInput type="password" value={bearerToken} placeholder="（可选）" onChange={(e) => setBearerToken(e.target.value)} />
          </Field>
          <div className="flex items-end">
            <Button variant="primary" loading={loading} onClick={() => void handleFetch()} className="w-full md:w-auto">
              <CloudArrowDown size={15} />
              拉取策略
            </Button>
          </div>
        </div>
      </Panel>

      <Panel
        title={
          <div className="flex items-center gap-2">
            <span>Rego 策略</span>
            <Tag tone={loaded ? 'green' : 'slate'}>{loaded ? '已拉取' : '未拉取'}</Tag>
          </div>
        }
        extra={
          <div className="flex items-center gap-2">
            <Button variant="danger" disabled={!loaded} onClick={() => setDeleteOpen(true)}>
              <Trash size={15} />
              删除
            </Button>
            <Button loading={syncing} onClick={() => void handleSync()}>
              <ArrowClockwise size={15} />
              手动同步
            </Button>
            <Button variant="primary" loading={saving} onClick={() => void handleSave()}>
              <FloppyDisk size={15} />
              保存
            </Button>
          </div>
        }
      >
        <YamlEditor value={content} onChange={setContent} height={440} language="rego" />
      </Panel>

      <ConfirmModal
        open={deleteOpen}
        title="删除 OPA 策略"
        danger
        loading={deleting}
        confirmText="删除"
        message={
          <>
            确定从 OPA 服务删除策略 <span className="font-mono text-slate-800">{policyId || '（默认策略 ID）'}</span> 吗？
          </>
        }
        onConfirm={() => void handleDelete()}
        onClose={() => setDeleteOpen(false)}
      />
    </div>
  )
}
