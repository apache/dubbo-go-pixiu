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

import { ReactNode, useEffect } from 'react'
import { X } from '@phosphor-icons/react'

interface ModalProps {
  open: boolean
  title: string
  onClose: () => void
  children: ReactNode
  footer?: ReactNode
  width?: 'md' | 'lg' | 'xl'
}

const WIDTH_CLASS = {
  md: 'max-w-lg',
  lg: 'max-w-3xl',
  xl: 'max-w-5xl',
}

export function Modal({ open, title, onClose, children, footer, width = 'md' }: ModalProps) {
  useEffect(() => {
    if (!open) return
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-[80] flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-slate-900/45" onClick={onClose} />
      <div
        role="dialog"
        aria-modal="true"
        className={`relative flex max-h-[88vh] w-full ${WIDTH_CLASS[width]} flex-col rounded-lg bg-white shadow-2xl`}
      >
        <div className="flex h-12 shrink-0 items-center justify-between border-b border-slate-200 px-4">
          <h3 className="text-sm font-semibold text-slate-800">{title}</h3>
          <button
            type="button"
            aria-label="关闭"
            className="rounded p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
            onClick={onClose}
          >
            <X size={16} />
          </button>
        </div>
        <div className="slim-scroll min-h-0 flex-1 overflow-y-auto p-4">{children}</div>
        {footer ? (
          <div className="flex shrink-0 items-center justify-end gap-2 border-t border-slate-200 px-4 py-3">{footer}</div>
        ) : null}
      </div>
    </div>
  )
}

interface ConfirmModalProps {
  open: boolean
  title: string
  message: ReactNode
  confirmText?: string
  danger?: boolean
  loading?: boolean
  onConfirm: () => void
  onClose: () => void
}

export function ConfirmModal({
  open,
  title,
  message,
  confirmText = '确认',
  danger,
  loading,
  onConfirm,
  onClose,
}: ConfirmModalProps) {
  return (
    <Modal
      open={open}
      title={title}
      onClose={onClose}
      footer={
        <>
          <button
            type="button"
            className="rounded-md border border-slate-300 bg-white px-3.5 py-1.5 text-sm text-slate-700 transition hover:bg-slate-50 active:translate-y-px"
            onClick={onClose}
          >
            取消
          </button>
          <button
            type="button"
            disabled={loading}
            className={`rounded-md px-3.5 py-1.5 text-sm font-medium text-white transition active:translate-y-px disabled:opacity-60 ${
              danger ? 'bg-rose-600 hover:bg-rose-700' : 'bg-blue-600 hover:bg-blue-700'
            }`}
            onClick={onConfirm}
          >
            {loading ? '处理中' : confirmText}
          </button>
        </>
      }
    >
      <div className="text-sm leading-6 text-slate-600">{message}</div>
    </Modal>
  )
}
