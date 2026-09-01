import { createContext, ReactNode, useCallback, useContext, useMemo, useRef, useState } from 'react'
import { CheckCircle, WarningCircle, XCircle, X } from '@phosphor-icons/react'

type ToastKind = 'success' | 'error' | 'info'

interface ToastItem {
  id: number
  kind: ToastKind
  text: string
}

interface ToastContextValue {
  success: (text: string) => void
  error: (text: string) => void
  info: (text: string) => void
}

const ToastContext = createContext<ToastContextValue | null>(null)

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext)
  if (!ctx) throw new Error('useToast must be used within ToastProvider')
  return ctx
}

const KIND_STYLE: Record<ToastKind, { icon: ReactNode; bar: string }> = {
  success: { icon: <CheckCircle size={18} weight="fill" className="text-emerald-500" />, bar: 'border-emerald-500' },
  error: { icon: <XCircle size={18} weight="fill" className="text-rose-500" />, bar: 'border-rose-500' },
  info: { icon: <WarningCircle size={18} weight="fill" className="text-blue-500" />, bar: 'border-blue-500' },
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<ToastItem[]>([])
  const nextId = useRef(1)

  const push = useCallback((kind: ToastKind, text: string) => {
    const id = nextId.current++
    setItems((prev) => [...prev, { id, kind, text }])
    window.setTimeout(() => {
      setItems((prev) => prev.filter((item) => item.id !== id))
    }, 3600)
  }, [])

  const value = useMemo<ToastContextValue>(
    () => ({
      success: (text) => push('success', text),
      error: (text) => push('error', text),
      info: (text) => push('info', text),
    }),
    [push],
  )

  return (
    <ToastContext.Provider value={value}>
      {children}
      <div className="pointer-events-none fixed right-4 top-4 z-[90] flex w-96 max-w-[calc(100vw-2rem)] flex-col gap-2">
        {items.map((item) => (
          <div
            key={item.id}
            className={`pointer-events-auto flex items-start gap-2 rounded-lg border-l-4 bg-white px-3 py-2.5 shadow-lg shadow-slate-900/10 ring-1 ring-slate-900/5 ${KIND_STYLE[item.kind].bar}`}
          >
            <span className="mt-0.5 shrink-0">{KIND_STYLE[item.kind].icon}</span>
            <p className="flex-1 break-words text-sm leading-5 text-slate-700">{item.text}</p>
            <button
              type="button"
              aria-label="关闭"
              className="shrink-0 rounded p-0.5 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
              onClick={() => setItems((prev) => prev.filter((it) => it.id !== item.id))}
            >
              <X size={14} />
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  )
}
