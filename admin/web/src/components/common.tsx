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

import { ButtonHTMLAttributes, InputHTMLAttributes, ReactNode, SelectHTMLAttributes, TextareaHTMLAttributes } from 'react'
import { CircleNotch, CloudSlash, Tray } from '@phosphor-icons/react'

/*
 * Shape system (locked): cards/panels rounded-lg, controls rounded-md.
 * Accent (locked): blue-600. Neutrals: slate.
 */

type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'ghost'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  loading?: boolean
}

const BUTTON_STYLE: Record<ButtonVariant, string> = {
  primary: 'bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-60',
  secondary: 'border border-slate-300 bg-white text-slate-700 hover:bg-slate-50 disabled:opacity-60',
  danger: 'bg-rose-600 text-white hover:bg-rose-700 disabled:opacity-60',
  ghost: 'text-slate-600 hover:bg-slate-100 disabled:opacity-60',
}

export function Button({ variant = 'secondary', loading, className = '', children, disabled, ...rest }: ButtonProps) {
  return (
    <button
      type="button"
      disabled={disabled || loading}
      className={`inline-flex items-center justify-center gap-1.5 whitespace-nowrap rounded-md px-3.5 py-1.5 text-sm font-medium transition active:translate-y-px ${BUTTON_STYLE[variant]} ${className}`}
      {...rest}
    >
      {loading ? <CircleNotch size={15} className="animate-spin" /> : null}
      {children}
    </button>
  )
}

export function Panel({
  title,
  extra,
  children,
  className = '',
  bodyClassName = '',
}: {
  title?: ReactNode
  extra?: ReactNode
  children: ReactNode
  className?: string
  bodyClassName?: string
}) {
  return (
    <section className={`rounded-lg bg-white shadow-sm ring-1 ring-slate-900/5 ${className}`}>
      {title !== undefined || extra ? (
        <header className="flex min-h-11 flex-wrap items-center justify-between gap-2 border-b border-slate-100 px-4 py-2">
          <h2 className="text-sm font-semibold text-slate-800">{title}</h2>
          {extra}
        </header>
      ) : null}
      <div className={`p-4 ${bodyClassName}`}>{children}</div>
    </section>
  )
}

type TagTone = 'blue' | 'green' | 'amber' | 'rose' | 'slate'

const TAG_STYLE: Record<TagTone, string> = {
  blue: 'bg-blue-50 text-blue-700 ring-blue-600/20',
  green: 'bg-emerald-50 text-emerald-700 ring-emerald-600/20',
  amber: 'bg-amber-50 text-amber-700 ring-amber-600/25',
  rose: 'bg-rose-50 text-rose-700 ring-rose-600/20',
  slate: 'bg-slate-100 text-slate-600 ring-slate-500/15',
}

export function Tag({ tone = 'slate', children }: { tone?: TagTone; children: ReactNode }) {
  return (
    <span className={`inline-flex items-center whitespace-nowrap rounded px-1.5 py-0.5 text-xs font-medium ring-1 ring-inset ${TAG_STYLE[tone]}`}>
      {children}
    </span>
  )
}

export function Field({ label, hint, error, children }: { label: ReactNode; hint?: string; error?: string; children: ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1 block text-xs font-medium text-slate-600">{label}</span>
      {children}
      {hint && !error ? <span className="mt-1 block text-xs text-slate-400">{hint}</span> : null}
      {error ? <span className="mt-1 block text-xs text-rose-600">{error}</span> : null}
    </label>
  )
}

const CONTROL_CLASS =
  'w-full rounded-md border border-slate-300 bg-white px-2.5 py-1.5 text-sm text-slate-800 placeholder:text-slate-400 transition focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 disabled:bg-slate-50'

export function TextInput(props: InputHTMLAttributes<HTMLInputElement>) {
  const { className = '', ...rest } = props
  return <input className={`${CONTROL_CLASS} ${className}`} {...rest} />
}

export function TextArea(props: TextareaHTMLAttributes<HTMLTextAreaElement>) {
  const { className = '', ...rest } = props
  return <textarea className={`${CONTROL_CLASS} ${className}`} {...rest} />
}

export function Select(props: SelectHTMLAttributes<HTMLSelectElement>) {
  const { className = '', children, ...rest } = props
  return (
    <select className={`${CONTROL_CLASS} ${className}`} {...rest}>
      {children}
    </select>
  )
}

export function Loading({ text = '加载中' }: { text?: string }) {
  return (
    <div className="flex items-center justify-center gap-2 py-12 text-sm text-slate-500">
      <CircleNotch size={18} className="animate-spin text-blue-500" />
      {text}
    </div>
  )
}

export function EmptyState({ text, action }: { text: string; action?: ReactNode }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-12 text-center">
      <Tray size={36} weight="duotone" className="text-slate-300" />
      <p className="text-sm text-slate-500">{text}</p>
      {action}
    </div>
  )
}

export function ErrorState({ text, onRetry }: { text: string; onRetry?: () => void }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-12 text-center">
      <CloudSlash size={36} weight="duotone" className="text-rose-300" />
      <p className="max-w-md break-words text-sm text-slate-500">{text}</p>
      {onRetry ? (
        <Button variant="secondary" onClick={onRetry}>
          重试
        </Button>
      ) : null}
    </div>
  )
}

export function formatTime(value?: string): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

/** Legacy list rows carry Go durations as nanosecond integers. */
export function formatDuration(value?: string | number): string {
  if (value === undefined || value === null || value === '') return '-'
  if (typeof value === 'string') return value
  if (!Number.isFinite(value) || value === 0) return '-'
  if (value % 1e9 === 0) return `${value / 1e9}s`
  if (value >= 1e9) return `${(value / 1e9).toFixed(1)}s`
  if (value >= 1e6) return `${Math.round(value / 1e6)}ms`
  return `${value}ns`
}
