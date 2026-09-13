export function ConfirmDialog({
  open,
  title = '确认操作',
  onConfirm,
  onCancel,
}: {
  open: boolean
  title?: string
  onConfirm: () => void
  onCancel: () => void
}) {
  if (!open) return null
  return (
    <div className="drawer-backdrop">
      <div className="confirm-dialog">
        <h3>{title}</h3>
        <p>此操作无法撤销，请确认继续。</p>
        <div>
          <button className="secondary" onClick={onCancel}>
            取消
          </button>
          <button className="primary" onClick={onConfirm}>
            确认
          </button>
        </div>
      </div>
    </div>
  )
}
