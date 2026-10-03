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
