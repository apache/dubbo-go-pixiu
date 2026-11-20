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

import { parseDocument, stringify } from 'yaml'
import type { AdminClusterObject } from '../../services/cluster-lifecycle-api'

export function emptyCluster(type = 'Static'): AdminClusterObject {
  return {
    kind: 'AdminCluster',
    metadata: { name: '' },
    spec: { type, endpoint: { address: '', port: '' } },
  }
}

function record(value: unknown, label: string): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    throw new TypeError(`${label} must be a mapping`)
  }
  return value as Record<string, unknown>
}

function onlyKeys(value: Record<string, unknown>, allowed: string[], label: string) {
  const unknown = Object.keys(value).find((key) => !allowed.includes(key))
  if (unknown) throw new Error(`${label}.${unknown} is not supported`)
}

function text(value: unknown, label: string) {
  if (value === undefined || value === null) return ''
  if (typeof value !== 'string') throw new TypeError(`${label} must be text`)
  return value
}

function port(value: unknown) {
  if (value === undefined || value === null) return ''
  if (typeof value !== 'number' && typeof value !== 'string') {
    throw new Error('port must be a number or text')
  }
  return value
}

export function parseClusterYaml(source: string, sourceId?: number): AdminClusterObject {
  const document = parseDocument(source, { uniqueKeys: true })
  if (document.errors.length) throw new Error(document.errors[0].message)
  const value = record(document.toJS(), 'Cluster YAML')
  if ('kind' in value || 'metadata' in value || 'spec' in value) {
    onlyKeys(value, ['kind', 'metadata', 'spec'], 'object')
    if (value.kind !== 'AdminCluster') throw new Error('kind must be AdminCluster')
    const metadata = record(value.metadata, 'metadata')
    const spec = record(value.spec, 'spec')
    onlyKeys(metadata, ['name'], 'metadata')
    if ('id' in spec) throw new Error('spec.id is managed by Pixiu')
    const endpoint = spec.endpoint === undefined ? {} : record(spec.endpoint, 'spec.endpoint')
    return {
      kind: 'AdminCluster',
      metadata: { name: text(metadata.name, 'metadata.name') },
      spec: {
        ...spec,
        type: text(spec.type, 'spec.type'),
        endpoint: {
          ...endpoint,
          address: text(endpoint.address, 'spec.endpoint.address'),
          port: port(endpoint.port),
        },
      },
    }
  }
  onlyKeys(value, ['id', 'name', 'type', 'address', 'port'], 'cluster')
  if ('id' in value && (typeof value.id !== 'number' || value.id !== sourceId)) {
    throw new Error('Runtime ID is managed by Pixiu; remove or keep the current id')
  }
  return {
    kind: 'AdminCluster',
    metadata: { name: text(value.name, 'name') },
    spec: {
      type: text(value.type, 'type'),
      endpoint: { address: text(value.address, 'address'), port: port(value.port) },
    },
  }
}

export function stringifyClusterYaml(object: AdminClusterObject) {
  const endpoint = object.spec.endpoint ?? {}
  if (
    Object.keys(object.spec).some((key) => !['type', 'endpoint'].includes(key)) ||
    Object.keys(endpoint).some((key) => !['address', 'port'].includes(key))
  ) {
    return stringify(object)
  }
  return stringify({
    name: object.metadata.name,
    type: object.spec.type ?? '',
    address: object.spec.endpoint?.address ?? '',
    port: object.spec.endpoint?.port ?? '',
  })
}

export function clusterFieldValue(object: AdminClusterObject, path: string): string {
  if (path === 'metadata.name') return object.metadata.name
  if (path === 'spec.type') return object.spec.type ?? ''
  if (path === 'spec.endpoint.address') return object.spec.endpoint?.address ?? ''
  if (path === 'spec.endpoint.port') return String(object.spec.endpoint?.port ?? '')
  return ''
}

export function updateClusterField(
  object: AdminClusterObject,
  path: string,
  value: string,
): AdminClusterObject {
  if (path === 'metadata.name') return { ...object, metadata: { ...object.metadata, name: value } }
  if (path === 'spec.type') return { ...object, spec: { ...object.spec, type: value } }
  if (path === 'spec.endpoint.address') {
    return {
      ...object,
      spec: { ...object.spec, endpoint: { ...object.spec.endpoint, address: value } },
    }
  }
  if (path === 'spec.endpoint.port') {
    const number = /^\d+$/.test(value) ? Number(value) : value
    return {
      ...object,
      spec: { ...object.spec, endpoint: { ...object.spec.endpoint, port: number } },
    }
  }
  return object
}
