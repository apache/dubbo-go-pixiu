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

import Editor, { OnMount } from '@monaco-editor/react'

interface YamlEditorProps {
  value: string
  onChange?: (value: string) => void
  height?: number | string
  readOnly?: boolean
  language?: 'yaml' | 'json' | 'plaintext' | 'rego'
}

/** Thin Monaco wrapper so every YAML surface shares the same skin. */
export function YamlEditor({ value, onChange, height = 380, readOnly, language = 'yaml' }: YamlEditorProps) {
  const handleMount: OnMount = (editor) => {
    editor.focus()
  }

  return (
    <div className="overflow-hidden rounded-md ring-1 ring-slate-900/10">
      <Editor
        height={height}
        language={language}
        theme="vs-dark"
        value={value}
        onChange={(next) => onChange?.(next ?? '')}
        onMount={handleMount}
        options={{
          readOnly,
          minimap: { enabled: false },
          fontSize: 13,
          lineNumbersMinChars: 3,
          scrollBeyondLastLine: false,
          wordWrap: 'on',
          tabSize: 2,
          automaticLayout: true,
          renderLineHighlight: 'none',
          padding: { top: 10, bottom: 10 },
          scrollbar: { verticalScrollbarSize: 8, horizontalScrollbarSize: 8 },
        }}
      />
    </div>
  )
}
