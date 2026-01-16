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

import { useState, useEffect, useCallback, useRef } from 'react';
import { Segmented, Form, Input, InputNumber, Select, theme } from 'antd';
import { useTranslation } from 'react-i18next';
import Editor, { type Monaco } from '@monaco-editor/react';
import { configureMonacoYaml } from 'monaco-yaml';
import { useThemeStore } from '../../stores/theme';
import type { FormInstance } from 'antd';

export type EditorMode = 'form' | 'yaml';

export interface FieldConfig {
  name: string;
  label: string;
  type: 'input' | 'number' | 'select' | 'textarea';
  required?: boolean;
  options?: { label: string; value: string | number | boolean }[];
  placeholder?: string;
  min?: number;
  max?: number;
  disabled?: boolean;  // Always disabled regardless of readOnly
  hidden?: boolean;    // Hide this field
}

interface DualModeEditorProps<T = unknown> {
  mode: EditorMode;
  onModeChange: (mode: EditorMode) => void;
  form: FormInstance<T>;
  yamlValue: string;
  onYamlChange: (value: string) => void;
  fields: FieldConfig[];
  formToYaml: (data: T, existingYaml: string) => string;
  yamlToForm: (yaml: string) => T | null;
  readOnly?: boolean;
}

export function DualModeEditor<T>({
  mode,
  onModeChange,
  form,
  yamlValue,
  onYamlChange,
  fields,
  formToYaml,
  yamlToForm,
  readOnly = false,
}: DualModeEditorProps<T>) {
  const { t } = useTranslation();
  const { isDark } = useThemeStore();
  const { token } = theme.useToken();
  const [syncError, setSyncError] = useState<string | null>(null);
  const monacoConfigured = useRef(false);

  const handleEditorBeforeMount = useCallback((monaco: Monaco) => {
    if (monacoConfigured.current) return;
    monacoConfigured.current = true;

    configureMonacoYaml(monaco, {
      enableSchemaRequest: false,
      validate: true,
      format: true,
      hover: true,
      completion: true,
    });
  }, []);

  const handleModeChange = useCallback(
    (newMode: EditorMode) => {
      if (newMode === mode) return;

      if (newMode === 'yaml') {
        const formData = form.getFieldsValue() as T;
        const yaml = formToYaml(formData, yamlValue);
        onYamlChange(yaml);
      } else {
        const formData = yamlToForm(yamlValue);
        if (formData) {
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          form.setFieldsValue(formData as any);
          setSyncError(null);
        } else {
          setSyncError(t('common.yamlParseError'));
          return;
        }
      }
      onModeChange(newMode);
    },
    [mode, form, yamlValue, formToYaml, yamlToForm, onYamlChange, onModeChange, t]
  );

  useEffect(() => {
    // Clear sync error when yaml value changes externally
    // eslint-disable-next-line react-hooks/set-state-in-effect -- resetting error state on prop change
    setSyncError(null);
  }, [yamlValue]);

  const renderField = (field: FieldConfig) => {
    // Skip hidden fields
    if (field.hidden) return null;

    const rules = field.required ? [{ required: true, message: t('common.required') }] : [];
    const isDisabled = readOnly || field.disabled;

    switch (field.type) {
      case 'input':
        return (
          <Form.Item key={field.name} name={field.name} label={field.label} rules={rules}>
            <Input placeholder={field.placeholder} disabled={isDisabled} />
          </Form.Item>
        );
      case 'number':
        return (
          <Form.Item key={field.name} name={field.name} label={field.label} rules={rules}>
            <InputNumber
              className="w-full"
              min={field.min}
              max={field.max}
              placeholder={field.placeholder}
              disabled={isDisabled}
            />
          </Form.Item>
        );
      case 'select':
        return (
          <Form.Item key={field.name} name={field.name} label={field.label} rules={rules}>
            <Select
              options={field.options}
              placeholder={field.placeholder}
              disabled={isDisabled}
            />
          </Form.Item>
        );
      case 'textarea':
        return (
          <Form.Item key={field.name} name={field.name} label={field.label} rules={rules}>
            <Input.TextArea
              rows={3}
              placeholder={field.placeholder}
              disabled={isDisabled}
            />
          </Form.Item>
        );
      default:
        return null;
    }
  };

  return (
    <div className="space-y-4">
      {/* Mode Switcher */}
      <div className="flex items-center justify-between">
        <Segmented
          value={mode}
          onChange={(value) => handleModeChange(value as EditorMode)}
          options={[
            { label: t('common.form'), value: 'form' },
            { label: 'YAML', value: 'yaml' },
          ]}
          size="middle"
        />
        {readOnly && (
          <span
            className="text-xs px-2 py-1 rounded"
            style={{
              color: token.colorTextSecondary,
              backgroundColor: token.colorFillSecondary,
            }}
          >
            {t('common.view')}
          </span>
        )}
      </div>

      {/* Error Message */}
      {syncError && (
        <div
          className="flex items-center gap-2 text-sm px-3 py-2 rounded-lg"
          style={{
            color: token.colorError,
            backgroundColor: token.colorErrorBg,
            border: `1px solid ${token.colorErrorBorder}`,
          }}
        >
          <svg className="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
          </svg>
          {syncError}
        </div>
      )}

      {/* Content Area */}
      <div
        className="rounded-xl overflow-hidden"
        style={{
          backgroundColor: token.colorBgContainer,
          border: `1px solid ${token.colorBorderSecondary}`,
        }}
      >
        {mode === 'form' ? (
          <div className="p-5">
            <Form
              form={form}
              layout="vertical"
              disabled={readOnly}
              className="max-w-lg"
            >
              {fields.map(renderField)}
            </Form>
          </div>
        ) : (
          <div className="h-[350px]">
            <Editor
              height="100%"
              defaultLanguage="yaml"
              value={yamlValue}
              onChange={(value) => onYamlChange(value || '')}
              beforeMount={handleEditorBeforeMount}
              theme={isDark ? 'vs-dark' : 'light'}
              options={{
                minimap: { enabled: false },
                fontSize: 13,
                lineNumbers: 'on',
                scrollBeyondLastLine: false,
                automaticLayout: true,
                readOnly,
                tabSize: 2,
                padding: { top: 16, bottom: 16 },
              }}
            />
          </div>
        )}
      </div>
    </div>
  );
}
