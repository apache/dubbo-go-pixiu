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

import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowRight, GitBranch, ShieldCheck, Swap } from '@phosphor-icons/react'
import { login, register } from '../api'
import { saveSession } from '../auth/session'
import { Button, Field, TextInput } from '../components/common'

type Mode = 'login' | 'register'

const HIGHLIGHTS = [
  { icon: <Swap size={18} />, title: 'HTTP 到 Dubbo 泛化调用', desc: '用高层路由模型描述入口与目标，发布后投影为运行时配置。' },
  { icon: <GitBranch size={18} />, title: '草稿与发布生命周期', desc: '草稿不影响运行中的网关，支持预览、Diff、发布历史与回滚。' },
  { icon: <ShieldCheck size={18} />, title: 'Schema 驱动校验', desc: '表单选项与服务端 Schema 同源，错误配置在发布前被拦截。' },
]

export function LoginPage() {
  const navigate = useNavigate()
  const [mode, setMode] = useState<Mode>('login')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const switchMode = (next: Mode) => {
    setMode(next)
    setError('')
    setNotice('')
    setConfirmPassword('')
  }

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    setError('')
    setNotice('')
    if (!username.trim() || !password) {
      setError('请输入用户名和密码')
      return
    }
    if (mode === 'register') {
      if (password !== confirmPassword) {
        setError('两次输入的密码不一致')
        return
      }
    }
    setSubmitting(true)
    try {
      if (mode === 'login') {
        const result = await login(username.trim(), password)
        saveSession({ username: result.username, token: result.token })
        navigate('/gateway/overview', { replace: true })
      } else {
        const message = await register(username.trim(), password)
        setNotice(message || '注册成功，请登录')
        setMode('login')
        setPassword('')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '请求失败')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-[100dvh] bg-white">
      <div className="relative hidden w-[46%] flex-col justify-between overflow-hidden bg-slate-900 p-10 lg:flex">
        <img
          src="/login-bg.png"
          alt=""
          aria-hidden
          className="pointer-events-none absolute inset-0 h-full w-full object-cover"
        />
        <div className="pointer-events-none absolute inset-0 bg-gradient-to-b from-slate-900/70 via-slate-900/35 to-slate-900/75" />
        <div className="relative flex items-center gap-2.5">
          <img
            src="/apache-asf-logo.svg"
            alt="Apache Software Foundation"
            className="h-11 w-[84px] shrink-0 object-contain object-left"
          />
          <span className="text-base font-semibold text-white">Pixiu Admin</span>
        </div>
        <div className="relative">
          <h1 className="max-w-md text-3xl font-semibold leading-snug tracking-tight text-white">
            Dubbo Go Pixiu 网关配置管理
          </h1>
          <p className="mt-3 max-w-md text-sm leading-6 text-slate-400">
            面向运维与研发的控制面：资源、方法映射、集群、Listener 与策略配置，统一在一个面板中维护。
          </p>
          <ul className="mt-8 space-y-5">
            {HIGHLIGHTS.map((item) => (
              <li key={item.title} className="flex gap-3">
                <span className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-slate-800 text-blue-400">
                  {item.icon}
                </span>
                <div>
                  <div className="text-sm font-medium text-slate-100">{item.title}</div>
                  <div className="mt-0.5 text-xs leading-5 text-slate-400">{item.desc}</div>
                </div>
              </li>
            ))}
          </ul>
        </div>
        <div className="relative text-xs text-slate-400">Apache Dubbo Go Pixiu</div>
      </div>

      <div className="flex flex-1 items-center justify-center p-6">
        <div className="w-full max-w-sm">
          <div className="mb-8 lg:hidden">
            <img
              src="/apache-asf-logo.svg"
              alt="Apache Software Foundation"
              className="h-12 w-[92px] object-contain object-left"
            />
            <h1 className="mt-4 text-xl font-semibold text-slate-900">Pixiu Admin</h1>
          </div>

          <div className="mb-6 flex rounded-lg bg-slate-100 p-1">
            {(['login', 'register'] as Mode[]).map((item) => (
              <button
                key={item}
                type="button"
                className={`flex-1 rounded-md py-1.5 text-sm font-medium transition ${
                  mode === item ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'
                }`}
                onClick={() => switchMode(item)}
              >
                {item === 'login' ? '登录' : '注册'}
              </button>
            ))}
          </div>

          <h2 className="text-lg font-semibold text-slate-900">{mode === 'login' ? '登录管理面板' : '注册新账户'}</h2>
          <p className="mt-1 text-sm text-slate-500">
            {mode === 'login' ? '使用管理员分配的账户登录。' : '注册成功后返回登录页进入面板。'}
          </p>

          <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
            <Field label="用户名">
              <TextInput
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="请输入用户名"
                autoComplete="username"
                autoFocus
              />
            </Field>
            <Field label="密码">
              <TextInput
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="请输入密码"
                autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
              />
            </Field>
            {mode === 'register' ? (
              <Field label="确认密码">
                <TextInput
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="请再次输入密码"
                  autoComplete="new-password"
                />
              </Field>
            ) : null}

            {error ? <p className="rounded-md bg-rose-50 px-3 py-2 text-sm text-rose-600">{error}</p> : null}
            {notice ? <p className="rounded-md bg-emerald-50 px-3 py-2 text-sm text-emerald-600">{notice}</p> : null}

            <Button type="submit" variant="primary" loading={submitting} className="w-full py-2">
              {mode === 'login' ? '登录' : '注册'}
              <ArrowRight size={15} />
            </Button>
          </form>
        </div>
      </div>
    </div>
  )
}
