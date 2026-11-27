/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { useEffect, useId, useRef } from 'react'

export function ConfirmDialog({
  open,
  title = '确认操作',
  description = '此操作无法撤销，请确认继续。',
  confirmLabel = '确认',
  cancelLabel = '取消',
  onConfirm,
  onCancel,
}: Readonly<{
  open: boolean
  title?: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  onConfirm: () => void
  onCancel: () => void
}>) {
  const titleId = useId()
  const descriptionId = useId()
  const dialogRef = useRef<HTMLDialogElement>(null)
  const cancelRef = useRef<HTMLButtonElement>(null)
  const onCancelRef = useRef(onCancel)
  onCancelRef.current = onCancel
  useEffect(() => {
    const dialog = dialogRef.current
    if (!open || !dialog) return
    const previous = document.activeElement instanceof HTMLElement ? document.activeElement : null
    dialog.showModal()
    cancelRef.current?.focus()
    return () => {
      dialog.close()
      previous?.focus()
    }
  }, [open])
  return (
    <dialog
      ref={dialogRef}
      className="confirm-dialog"
      aria-labelledby={titleId}
      aria-describedby={descriptionId}
      onCancel={(event) => {
        event.preventDefault()
        onCancelRef.current()
      }}
    >
      <h3 id={titleId}>{title}</h3>
      <p id={descriptionId}>{description}</p>
      <div>
        <button ref={cancelRef} type="button" className="secondary" onClick={onCancel}>
          {cancelLabel}
        </button>
        <button type="button" className="primary" onClick={onConfirm}>
          {confirmLabel}
        </button>
      </div>
    </dialog>
  )
}
