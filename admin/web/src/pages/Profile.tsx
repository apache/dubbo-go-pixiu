import { FormEvent, useCallback, useEffect, useState } from 'react'
import { Key, UserCircle } from '@phosphor-icons/react'
import { editPassword, getUserInfo } from '../api'
import { UserInfo } from '../types'
import { Button, ErrorState, Field, Loading, Panel, Tag, TextInput } from '../components/common'
import { useToast } from '../components/Toast'

export function ProfilePage() {
  const toast = useToast()
  const [info, setInfo] = useState<UserInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [formError, setFormError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      setInfo(await getUserInfo())
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载用户信息失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    setFormError('')
    if (!oldPassword || !newPassword) {
      setFormError('请输入旧密码和新密码')
      return
    }
    if (newPassword !== confirmPassword) {
      setFormError('两次输入的新密码不一致')
      return
    }
    setSubmitting(true)
    try {
      const message = await editPassword(oldPassword, newPassword)
      toast.success(message || '密码修改成功')
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch (err) {
      setFormError(err instanceof Error ? err.message : '修改失败')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="mx-auto max-w-3xl space-y-4">
      <Panel title="个人信息">
        {loading ? (
          <Loading />
        ) : error ? (
          <ErrorState text={error} onRetry={load} />
        ) : info ? (
          <div className="flex items-start gap-4">
            <span className="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-blue-600/10 text-blue-700">
              <UserCircle size={36} />
            </span>
            <dl className="grid flex-1 grid-cols-1 gap-x-8 gap-y-3 sm:grid-cols-3">
              <div>
                <dt className="text-xs text-slate-400">用户名</dt>
                <dd className="mt-0.5 text-sm font-medium text-slate-800">{info.username}</dd>
              </div>
              <div>
                <dt className="text-xs text-slate-400">用户 ID</dt>
                <dd className="mt-0.5 font-mono text-sm text-slate-800">{info.id}</dd>
              </div>
              <div>
                <dt className="text-xs text-slate-400">角色</dt>
                <dd className="mt-0.5">
                  <Tag tone="blue">{String(info.role) === '1' ? '管理员' : info.role || '-'}</Tag>
                </dd>
              </div>
            </dl>
          </div>
        ) : null}
      </Panel>

      <Panel
        title={
          <div className="flex items-center gap-2">
            <Key size={16} className="text-slate-500" />
            <span>修改密码</span>
          </div>
        }
      >
        <form className="max-w-sm space-y-4" onSubmit={handleSubmit}>
          <Field label="旧密码">
            <TextInput
              type="password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              placeholder="请输入当前密码"
              autoComplete="current-password"
            />
          </Field>
          <Field label="新密码">
            <TextInput
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              placeholder="请输入新密码"
              autoComplete="new-password"
            />
          </Field>
          <Field label="确认新密码">
            <TextInput
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="请再次输入新密码"
              autoComplete="new-password"
            />
          </Field>
          {formError ? <p className="rounded-md bg-rose-50 px-3 py-2 text-sm text-rose-600">{formError}</p> : null}
          <Button type="submit" variant="primary" loading={submitting}>
            确认修改
          </Button>
        </form>
      </Panel>
    </div>
  )
}
