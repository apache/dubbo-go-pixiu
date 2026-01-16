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

import { useState, useCallback, useMemo } from 'react';
import yaml from 'js-yaml';
import { message } from 'antd';
import { useTranslation } from 'react-i18next';

export interface UseYamlPreviewOptions<T> {
  formData: T;
  transformer: (data: T) => object;
  enabled?: boolean;
}

export interface UseYamlPreviewReturn {
  yaml: string;
  error: string | null;
  regenerate: () => void;
  copyToClipboard: () => Promise<void>;
}

export function useYamlPreview<T>(
  options: UseYamlPreviewOptions<T>
): UseYamlPreviewReturn {
  const { formData, transformer, enabled = true } = options;
  const { t } = useTranslation();

  // Use useMemo to compute yaml synchronously when formData changes
  const { yamlString, error } = useMemo(() => {
    if (!enabled) {
      return { yamlString: '', error: null };
    }

    try {
      const transformed = transformer(formData);
      const result = yaml.dump(transformed, {
        indent: 2,
        lineWidth: -1,
        noRefs: true,
        sortKeys: false,
      });
      return { yamlString: result, error: null };
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      return { yamlString: '', error: errorMessage };
    }
  }, [formData, transformer, enabled]);

  // For manual regeneration (useful if user edits yaml directly)
  const [manualYaml, setManualYaml] = useState<string | null>(null);

  const regenerate = useCallback(() => {
    // Clear manual override, fall back to computed value
    setManualYaml(null);
  }, []);

  const copyToClipboard = useCallback(async () => {
    const content = manualYaml ?? yamlString;
    if (!content) {
      message.warning(t('common.noContentToCopy'));
      return;
    }

    try {
      await navigator.clipboard.writeText(content);
      message.success(t('common.copySuccess'));
    } catch {
      message.error(t('common.copyFailed'));
    }
  }, [manualYaml, yamlString, t]);

  return useMemo(
    () => ({
      yaml: manualYaml ?? yamlString,
      error,
      regenerate,
      copyToClipboard,
    }),
    [manualYaml, yamlString, error, regenerate, copyToClipboard]
  );
}
