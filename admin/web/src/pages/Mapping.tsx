import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, PencilSimple, Plus, Trash } from '@phosphor-icons/react'
import {
  createMethod,
  deleteMethod,
  getMethodDetail,
  getMethodList,
  getResourceDetail,
  modifyMethod,
  modifyResource,
} from '../api'
import { Method } from '../types'
import { Button, EmptyState, ErrorState, formatDuration, Loading, Panel, Tag } from '../components/common'
import { ConfirmModal, Modal } from '../components/Modal'
import { YamlEditor } from '../components/YamlEditor'
import { useToast } from '../components/Toast'

const METHOD_TEMPLATE = `resourcePath: ''
enable: true
timeout: 1s
httpVerb: GET
inboundRequest:
  requestType: http
  headers: []
  queryStrings: []
  requestBody: []
integrationRequest:
  requestType: dubbo
  mappingParams: []
`

type EditorMode = { kind: 'create' } | { kind: 'edit'; method: Method } | null

export function MappingPage() {
  const { resourceId = '' } = useParams()
  const navigate = useNavigate()
  const toast = useToast()

  const [resourceYaml, setResourceYaml] = useState('')
  const [resourceLoading, setResourceLoading] = useState(true)
  const [resourceError, setResourceError] = useState('')
  const [resourceSaving, setResourceSaving] = useState(false)

  const [methods, setMethods] = useState<Method[]>([])
  const [listLoading, setListLoading] = useState(true)
  const [listError, setListError] = useState('')

  const [editorMode, setEditorMode] = useState<EditorMode>(null)
  const [editorYaml, setEditorYaml] = useState('')
  const [editorLoading, setEditorLoading] = useState(false)
  const [editorSaving, setEditorSaving] = useState(false)

  const [deleteTarget, setDeleteTarget] = useState<Method | null>(null)
  const [deleting, setDeleting] = useState(false)

  const loadResource = useCallback(async () => {
    setResourceLoading(true)
    setResourceError('')
    try {
      setResourceYaml(await getResourceDetail(resourceId))
    } catch (err) {
      setResourceError(err instanceof Error ? err.message : '加载 Resource 失败')
    } finally {
      setResourceLoading(false)
    }
  }, [resourceId])

  const loadMethods = useCallback(async () => {
    setListLoading(true)
    setListError('')
    try {
      setMethods(await getMethodList(resourceId))
    } catch (err) {
      setListError(err instanceof Error ? err.message : '加载 Method 列表失败')
    } finally {
      setListLoading(false)
    }
  }, [resourceId])

  useEffect(() => {
    void loadResource()
    void loadMethods()
  }, [loadResource, loadMethods])

  const handleSaveResource = async () => {
    setResourceSaving(true)
    try {
      await modifyResource(resourceId, resourceYaml)
      toast.success('Resource 已保存')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存失败')
    } finally {
      setResourceSaving(false)
    }
  }

  const openCreate = () => {
    setEditorMode({ kind: 'create' })
    setEditorYaml(METHOD_TEMPLATE)
  }

  const openEdit = async (method: Method) => {
    setEditorMode({ kind: 'edit', method })
    setEditorYaml('')
    setEditorLoading(true)
    try {
      setEditorYaml(await getMethodDetail(resourceId, method.id ?? ''))
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '加载 Method 详情失败')
      setEditorMode(null)
    } finally {
      setEditorLoading(false)
    }
  }

  const handleSaveMethod = async () => {
    if (!editorMode) return
    setEditorSaving(true)
    try {
      if (editorMode.kind === 'create') {
        await createMethod(resourceId, editorYaml)
        toast.success('Method 已创建')
      } else {
        await modifyMethod(resourceId, editorMode.method.id ?? '', editorYaml)
        toast.success('Method 已保存')
      }
      setEditorMode(null)
      void loadMethods()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '保存失败')
    } finally {
      setEditorSaving(false)
    }
  }

  const handleDeleteMethod = async () => {
    if (!deleteTarget?.id) return
    setDeleting(true)
    try {
      await deleteMethod(resourceId, deleteTarget.id)
      toast.success(`Method ${deleteTarget.id} 已删除`)
      setDeleteTarget(null)
      void loadMethods()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : '删除失败')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="mx-auto max-w-[1400px] space-y-4">
      <div className="flex items-center justify-between">
        <button
          type="button"
          className="inline-flex items-center gap-1.5 text-sm text-slate-500 transition hover:text-slate-800"
          onClick={() => navigate('/gateway/overview')}
        >
          <ArrowLeft size={15} />
          返回概览
        </button>
        <span className="text-xs text-slate-400">
          Resource ID <span className="font-mono">{resourceId}</span>
        </span>
      </div>

      <Panel
        title="Resource 配置"
        extra={
          <Button
            variant="primary"
            onClick={handleSaveResource}
            loading={resourceSaving}
            disabled={resourceLoading || !!resourceError}
          >
            保存
          </Button>
        }
      >
        {resourceLoading ? (
          <Loading />
        ) : resourceError ? (
          <ErrorState text={resourceError} onRetry={loadResource} />
        ) : (
          <YamlEditor value={resourceYaml} onChange={setResourceYaml} height={240} />
        )}
      </Panel>

      <Panel
        title="Method 映射"
        extra={
          <div className="flex items-center gap-2">
            <Button onClick={() => void loadMethods()}>刷新</Button>
            <Button variant="primary" onClick={openCreate}>
              <Plus size={15} />
              新增 Method
            </Button>
          </div>
        }
        bodyClassName="p-0"
      >
        {listLoading ? (
          <Loading />
        ) : listError ? (
          <ErrorState text={listError} onRetry={loadMethods} />
        ) : methods.length === 0 ? (
          <EmptyState
            text="该 Resource 下暂无 Method 映射。"
            action={
              <Button variant="primary" onClick={openCreate}>
                <Plus size={15} />
                新增 Method
              </Button>
            }
          />
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100 text-left text-xs text-slate-500">
                <th className="px-4 py-2.5 font-medium">ID</th>
                <th className="px-4 py-2.5 font-medium">HTTP 方法</th>
                <th className="px-4 py-2.5 font-medium">状态</th>
                <th className="px-4 py-2.5 font-medium">超时</th>
                <th className="px-4 py-2.5 font-medium">Resource 路径</th>
                <th className="px-4 py-2.5 text-right font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              {methods.map((item) => (
                <tr key={item.id} className="border-b border-slate-50 last:border-0 hover:bg-slate-50/60">
                  <td className="px-4 py-2.5 font-mono text-xs text-slate-600">{item.id}</td>
                  <td className="px-4 py-2.5">
                    <Tag tone="blue">{item.httpVerb || '-'}</Tag>
                  </td>
                  <td className="px-4 py-2.5">
                    <Tag tone={item.enable ? 'green' : 'slate'}>{item.enable ? '启用' : '停用'}</Tag>
                  </td>
                  <td className="px-4 py-2.5 text-slate-600">{formatDuration(item.timeout)}</td>
                  <td className="max-w-56 truncate px-4 py-2.5 font-mono text-xs text-slate-500">{item.resourcePath || '-'}</td>
                  <td className="px-4 py-2.5">
                    <div className="flex items-center justify-end gap-1">
                      <Button variant="ghost" className="px-2" onClick={() => void openEdit(item)}>
                        <PencilSimple size={15} />
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
        open={editorMode !== null}
        title={editorMode?.kind === 'create' ? '新增 Method' : `编辑 Method ${editorMode?.kind === 'edit' ? editorMode.method.id : ''}`}
        width="lg"
        onClose={() => setEditorMode(null)}
        footer={
          <>
            <Button onClick={() => setEditorMode(null)}>取消</Button>
            <Button variant="primary" loading={editorSaving} disabled={editorLoading} onClick={handleSaveMethod}>
              保存
            </Button>
          </>
        }
      >
        <p className="mb-2 text-xs text-slate-500">以 YAML 描述 Method 的请求映射与后端集成。</p>
        {editorLoading ? <Loading /> : <YamlEditor value={editorYaml} onChange={setEditorYaml} height={360} />}
      </Modal>

      <ConfirmModal
        open={deleteTarget !== null}
        title="删除 Method"
        danger
        loading={deleting}
        confirmText="删除"
        message={<>确定删除 Method（ID {deleteTarget?.id}）吗？删除后该映射立即不可用。</>}
        onConfirm={handleDeleteMethod}
        onClose={() => setDeleteTarget(null)}
      />
    </div>
  )
}
