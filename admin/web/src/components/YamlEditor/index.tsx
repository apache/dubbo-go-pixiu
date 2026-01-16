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

import { useCallback } from 'react';
import Editor from '@monaco-editor/react';
import { theme, Button, message, Tooltip } from 'antd';
import { CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useThemeStore } from '../../stores/theme';

interface YamlEditorProps {
  value: string;
  onChange?: (value: string) => void;
  height?: string;
  readOnly?: boolean;
  /** Show copy button in the top-right corner */
  showCopyButton?: boolean;
}

export function YamlEditor({
  value,
  onChange,
  height = '300px',
  readOnly = false,
  showCopyButton = false,
}: YamlEditorProps) {
  const { isDark } = useThemeStore();
  const { token } = theme.useToken();
  const { t } = useTranslation();

  const handleCopy = useCallback(async () => {
    if (!value) {
      message.warning(t('common.noContentToCopy'));
      return;
    }
    try {
      await navigator.clipboard.writeText(value);
      message.success(t('common.copySuccess'));
    } catch {
      message.error(t('common.copyFailed'));
    }
  }, [value, t]);

  return (
    <div
      style={{
        height,
        borderRadius: token.borderRadius,
        border: `1px solid ${token.colorBorderSecondary}`,
        overflow: 'hidden',
        backgroundColor: isDark ? '#1e1e1e' : '#ffffff',
        position: 'relative',
      }}
    >
      {showCopyButton && (
        <Tooltip title={t('common.copy')}>
          <Button
            type="text"
            size="small"
            icon={<CopyOutlined />}
            onClick={handleCopy}
            style={{
              position: 'absolute',
              top: 8,
              right: 8,
              zIndex: 10,
              backgroundColor: isDark ? 'rgba(30, 30, 30, 0.8)' : 'rgba(255, 255, 255, 0.8)',
            }}
          />
        </Tooltip>
      )}
      <Editor
        height="100%"
        language="yaml"
        theme={isDark ? 'vs-dark' : 'light'}
        value={value}
        onChange={(val) => onChange?.(val || '')}
        options={{
          minimap: { enabled: false },
          fontSize: 13,
          lineNumbers: 'on',
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 2,
          wordWrap: 'on',
          folding: true,
          readOnly,
          padding: { top: 12, bottom: 12 },
          lineNumbersMinChars: 3,
          glyphMargin: false,
          renderLineHighlight: 'line',
        }}
      />
    </div>
  );
}

export default YamlEditor;
