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

import yaml from 'js-yaml';
import type {
  ClusterFormData,
  ListenerFormData,
  ResourceFormData,
  ResourceType,
} from '../types/forms';

// 解析 YAML 为对象
export function parseYaml(yamlString: string): Record<string, unknown> | null {
  try {
    return yaml.load(yamlString) as Record<string, unknown>;
  } catch {
    return null;
  }
}

// 对象转 YAML
export function toYaml(data: Record<string, unknown>): string {
  return yaml.dump(data, { indent: 2, lineWidth: -1, noRefs: true });
}

// Cluster: 表单数据合并到 YAML
export function mergeClusterForm(formData: ClusterFormData, existingYaml: string): string {
  const existing = parseYaml(existingYaml) || {};
  return toYaml({
    ...existing,
    name: formData.name,
    type: formData.type,
    address: formData.address,
    port: formData.port,
  });
}

// Cluster: 从 YAML 提取表单数据
export function extractClusterForm(yamlString: string): ClusterFormData | null {
  const parsed = parseYaml(yamlString);
  if (!parsed) return null;
  return {
    id: parsed.id as string | undefined,
    name: (parsed.name as string) || '',
    type: (parsed.type as string) || '',
    address: (parsed.address as string) || '',
    port: (parsed.port as number) || 8080,
  };
}

// Listener: 表单数据合并到 YAML
export function mergeListenerForm(formData: ListenerFormData, existingYaml: string): string {
  const existing = parseYaml(existingYaml) || {};
  const existingAddress = (existing.address as Record<string, unknown>) || {};

  return toYaml({
    ...existing,
    name: formData.name,
    protocol_type: formData.protocol,
    address: {
      ...existingAddress,
      socket_address: {
        address: formData.address,
        port: formData.port,
      },
      name: formData.name,
    },
  });
}

// Listener: 从 YAML 提取表单数据
export function extractListenerForm(yamlString: string): ListenerFormData | null {
  const parsed = parseYaml(yamlString);
  if (!parsed) return null;

  const address = parsed.address as Record<string, unknown> | undefined;
  const socketAddress = address?.socket_address as Record<string, unknown> | undefined;

  return {
    id: parsed.id as string | undefined,
    name: (parsed.name as string) || '',
    protocol: (parsed.protocol_type as string) || 'HTTP',
    address: (socketAddress?.address as string) || '',
    port: (socketAddress?.port as number) || 8080,
  };
}

// Resource: 表单数据合并到 YAML
export function mergeResourceForm(formData: ResourceFormData, existingYaml: string): string {
  const existing = parseYaml(existingYaml) || {};

  const merged: Record<string, unknown> = {
    ...existing,
    path: formData.path,
    type: formData.type,
  };

  if (formData.description) {
    merged.description = formData.description;
  } else {
    delete merged.description;
  }

  if (formData.timeout) {
    merged.timeout = formData.timeout;
  } else {
    delete merged.timeout;
  }

  return toYaml(merged);
}

// Resource: 从 YAML 提取表单数据
export function extractResourceForm(yamlString: string): ResourceFormData | null {
  const parsed = parseYaml(yamlString);
  if (!parsed) return null;

  return {
    id: parsed.id as string | undefined,
    path: (parsed.path as string) || '',
    type: (parsed.type as string) || 'Restful',
    description: parsed.description as string | undefined,
    timeout: parsed.timeout as string | undefined,
  };
}

// 兼容旧接口
export function clusterFormToYaml(data: ClusterFormData): string {
  return mergeClusterForm(data, '');
}

export function yamlToClusterForm(yamlString: string): ClusterFormData | null {
  return extractClusterForm(yamlString);
}

export function listenerFormToYaml(data: ListenerFormData): string {
  return mergeListenerForm(data, '');
}

export function yamlToListenerForm(yamlString: string): ListenerFormData | null {
  return extractListenerForm(yamlString);
}

export function resourceFormToYaml(data: ResourceFormData): string {
  return mergeResourceForm(data, '');
}

export function yamlToResourceForm(yamlString: string): ResourceFormData | null {
  return extractResourceForm(yamlString);
}

export function getFormToYamlConverter(type: ResourceType) {
  switch (type) {
    case 'cluster':
      return clusterFormToYaml;
    case 'listener':
      return listenerFormToYaml;
    case 'resource':
      return resourceFormToYaml;
  }
}

export function getYamlToFormConverter(type: ResourceType) {
  switch (type) {
    case 'cluster':
      return yamlToClusterForm;
    case 'listener':
      return yamlToListenerForm;
    case 'resource':
      return yamlToResourceForm;
  }
}
