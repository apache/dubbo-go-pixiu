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
