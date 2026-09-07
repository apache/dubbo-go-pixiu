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

import { ReactNode, useEffect, useRef, useState } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import {
  Broadcast,
  CaretDown,
  Database,
  Gauge,
  GitBranch,
  PuzzlePiece,
  ShieldCheck,
  SignOut,
  Speedometer,
  UserCircle,
} from '@phosphor-icons/react'
import { clearSession, loadSession } from '../auth/session'
import { logout } from '../api'
import { useToast } from './Toast'

interface MenuChild {
  label: string
  to: string
  icon: ReactNode
}

interface MenuGroup {
  label: string
  children: MenuChild[]
}

const MENU: MenuGroup[] = [
  {
    label: '网关配置',
    children: [
      { label: '概览', to: '/gateway/overview', icon: <Gauge size={17} /> },
      { label: '路由模型', to: '/gateway/route-binding', icon: <GitBranch size={17} /> },
      { label: '插件配置', to: '/gateway/plugin-group', icon: <PuzzlePiece size={17} /> },
      { label: '集群管理', to: '/gateway/cluster', icon: <Database size={17} /> },
      { label: 'Listener 管理', to: '/gateway/listener', icon: <Broadcast size={17} /> },
    ],
  },
  {
    label: '限流配置',
    children: [{ label: '限流配置', to: '/flow/ratelimit', icon: <Speedometer size={17} /> }],
  },
  {
    label: 'OPA 配置',
    children: [{ label: 'OPA 配置', to: '/opa', icon: <ShieldCheck size={17} /> }],
  },
]

function AccountMenu() {
  const navigate = useNavigate()
  const toast = useToast()
  const [open, setOpen] = useState(false)
  const [username, setUsername] = useState('')
  const rootRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setUsername(loadSession()?.username ?? '')
  }, [])

  useEffect(() => {
    if (!open) return
    const onClick = (event: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(event.target as Node)) setOpen(false)
    }
    window.addEventListener('mousedown', onClick)
    return () => window.removeEventListener('mousedown', onClick)
  }, [open])

  const handleLogout = async () => {
    try {
      await logout()
    } catch {
      // Even when the backend call fails, local state must be cleaned up.
    }
    clearSession()
    toast.info('已退出登录')
    navigate('/login', { replace: true })
  }

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        className="flex items-center gap-2 rounded-md px-2.5 py-1.5 text-sm text-slate-700 transition hover:bg-slate-100"
        onClick={() => setOpen((prev) => !prev)}
      >
        <span className="flex h-7 w-7 items-center justify-center rounded-full bg-blue-600/10 text-blue-700">
          <UserCircle size={20} />
        </span>
        <span className="max-w-32 truncate font-medium">{username || '未登录'}</span>
        <CaretDown size={13} className={`text-slate-400 transition-transform ${open ? 'rotate-180' : ''}`} />
      </button>
      {open ? (
        <div className="absolute right-0 top-full z-40 mt-1 w-44 overflow-hidden rounded-lg bg-white py-1 shadow-lg ring-1 ring-slate-900/10">
          <button
            type="button"
            className="flex w-full items-center gap-2 px-3.5 py-2 text-left text-sm text-slate-700 hover:bg-slate-50"
            onClick={() => {
              setOpen(false)
              navigate('/profile')
            }}
          >
            <UserCircle size={16} className="text-slate-400" />
            个人中心
          </button>
          <button
            type="button"
            className="flex w-full items-center gap-2 px-3.5 py-2 text-left text-sm text-rose-600 hover:bg-rose-50"
            onClick={handleLogout}
          >
            <SignOut size={16} />
            退出登录
          </button>
        </div>
      ) : null}
    </div>
  )
}

export function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="flex h-full">
      <aside className="flex w-56 shrink-0 flex-col bg-slate-900">
        <div className="flex h-14 items-center gap-2.5 px-4">
          <img
            src="/apache-asf-logo.svg"
            alt="Apache Software Foundation"
            className="h-9 w-[69px] shrink-0 object-contain object-left"
          />
          <div className="leading-tight">
            <div className="text-sm font-semibold text-white">Pixiu Admin</div>
            <div className="text-[11px] text-slate-400">Dubbo Go Pixiu</div>
          </div>
        </div>
        <nav className="slim-scroll flex-1 overflow-y-auto px-2.5 pb-4">
          {MENU.map((group) => (
            <div key={group.label} className="mt-4 first:mt-2">
              <div className="px-2 pb-1.5 text-[11px] font-medium uppercase tracking-wider text-slate-500">{group.label}</div>
              <ul className="space-y-0.5">
                {group.children.map((item) => (
                  <li key={item.to}>
                    <NavLink
                      to={item.to}
                      className={({ isActive }) =>
                        `flex items-center gap-2.5 rounded-md px-2.5 py-2 text-[13px] transition ${
                          isActive
                            ? 'bg-blue-600 font-medium text-white'
                            : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                        }`
                      }
                    >
                      {item.icon}
                      {item.label}
                    </NavLink>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </nav>
      </aside>
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-14 shrink-0 items-center justify-between border-b border-slate-200 bg-white px-5">
          <div className="text-sm text-slate-500">API 网关配置管理</div>
          <AccountMenu />
        </header>
        <main className="slim-scroll min-h-0 flex-1 overflow-y-auto p-5">{children}</main>
      </div>
    </div>
  )
}
