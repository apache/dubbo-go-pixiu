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

import { parse as parseYaml, stringify as stringifyYaml } from 'yaml'
import type {
  AdminRouteBindingObject,
  RouteBindingFieldSchema,
  RouteBindingObjectSchema,
} from '../../types/api'

type RouteEditorOptions = {
  entryProtocols: string[]
  targetProtocols: string[]
  httpMethods: string[]
  paramTypes: string[]
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return value !== null && typeof value === 'object' && !Array.isArray(value)
}

function hasOwn(value: object, key: PropertyKey) {
  return Object.prototype.hasOwnProperty.call(value, key)
}

function cloneValue(value: unknown): unknown {
  if (Array.isArray(value)) return value.map(cloneValue)
  if (isRecord(value)) {
    return Object.fromEntries(Object.entries(value).map(([key, child]) => [key, cloneValue(child)]))
  }
  return value
}

function fail(path: string, message: string): never {
  throw new Error(`${path}: ${message}`)
}

function emptyValue(field: RouteBindingFieldSchema): unknown {
  if (hasOwn(field, 'default') && field.type !== 'object') return cloneValue(field.default)
  switch (field.type) {
    case 'string':
      return ''
    case 'integer':
      return 0
    case 'boolean':
      return false
    case 'array':
      return []
    case 'map':
      return {}
    case 'object': {
      const defaults = isRecord(field.default)
        ? (cloneValue(field.default) as Record<string, unknown>)
        : {}
      const properties = Object.fromEntries(
        Object.entries(field.properties || {}).map(([key, child]) => [
          key,
          hasOwn(defaults, key) && !hasOwn(child, 'default')
            ? cloneValue(defaults[key])
            : emptyValue(child),
        ]),
      )
      return { ...defaults, ...properties }
    }
  }
}

export function createDefaultRouteBinding(
  schema: RouteBindingObjectSchema,
): AdminRouteBindingObject {
  const defaultValues = Object.fromEntries(
    Object.entries(schema.fields).map(([key, field]) => [key, emptyValue(field)]),
  )
  return {
    kind: schema.kind,
    metadata: { name: '' },
    spec: defaultValues as AdminRouteBindingObject['spec'],
  }
}

function validateField(value: unknown, field: RouteBindingFieldSchema, path: string): unknown {
  switch (field.type) {
    case 'string':
      if (typeof value !== 'string') fail(path, 'must be a string')
      if (field.pattern && !new RegExp(field.pattern).test(value)) {
        fail(path, `must match ${field.pattern}`)
      }
      break
    case 'integer':
      if (typeof value !== 'number' || !Number.isInteger(value)) fail(path, 'must be an integer')
      if (field.minimum !== undefined && value < field.minimum) {
        fail(path, `must be at least ${field.minimum}`)
      }
      break
    case 'boolean':
      if (typeof value !== 'boolean') fail(path, 'must be a boolean')
      break
    case 'object': {
      if (!isRecord(value)) fail(path, 'must be an object')
      return validateFields(value, field.properties || {}, path, Boolean(field.allowUnknown))
    }
    case 'array':
      if (!Array.isArray(value)) fail(path, 'must be an array')
      if (!field.items) return value.map(cloneValue)
      return value.map((item, index) => validateField(item, field.items!, `${path}[${index}]`))
    case 'map': {
      if (!isRecord(value)) fail(path, 'must be an object')
      if (!field.additionalProperties) return { ...value }
      return Object.fromEntries(
        Object.entries(value).map(([key, child]) => [
          key,
          validateField(child, field.additionalProperties!, `${path}.${key}`),
        ]),
      )
    }
    default:
      fail(path, `uses unsupported schema type ${String(field.type)}`)
  }

  if (field.enum && !field.enum.some((item) => Object.is(item, value))) {
    fail(path, `must be one of: ${field.enum.map(String).join(', ')}`)
  }
  return value
}

function validateFields(
  source: Record<string, unknown>,
  fields: Record<string, RouteBindingFieldSchema>,
  path: string,
  allowUnknown: boolean,
) {
  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(source)) {
    const field = fields[key]
    if (!field) {
      if (allowUnknown) result[key] = cloneValue(value)
      else fail(`${path}.${key}`, 'is not defined by the route schema')
      continue
    }
    result[key] = validateField(value, field, `${path}.${key}`)
  }

  for (const [key, field] of Object.entries(fields)) {
    if (hasOwn(result, key)) continue
    if (hasOwn(field, 'default')) {
      result[key] = validateField(cloneValue(field.default), field, `${path}.${key}`)
    } else if (field.required) fail(`${path}.${key}`, 'is required')
  }
  return result
}

export function normaliseRouteBindingObject(
  value: unknown,
  schema: RouteBindingObjectSchema,
  stripLegacyPublishPreference = false,
): AdminRouteBindingObject {
  if (!isRecord(value)) fail('YAML', 'top level must be an object')
  for (const key of Object.keys(value)) {
    if (!['kind', 'metadata', 'spec'].includes(key)) fail(key, 'is not part of AdminRouteBinding')
  }
  if (value.kind !== schema.kind) fail('kind', `must be ${schema.kind}`)
  if (!isRecord(value.metadata)) fail('metadata', 'must be an object')
  for (const key of Object.keys(value.metadata)) {
    if (key !== 'name') fail(`metadata.${key}`, 'is not defined by the route schema')
  }
  if (typeof value.metadata.name !== 'string') fail('metadata.name', 'must be a string')
  if (!isRecord(value.spec)) fail('spec', 'must be an object')

  // The old publish preference was a no-op and is deliberately removed by the
  // backend normalizer. Reject it in edited YAML so it cannot appear synced
  // while being omitted from the submitted object.
  const specInput = { ...value.spec }
  const hasLegacyPublishPreference = hasOwn(specInput, 'publish')
  if (hasLegacyPublishPreference && !stripLegacyPublishPreference) {
    fail('spec.publish', 'is no longer supported; remove this field')
  }
  if (hasLegacyPublishPreference) delete specInput.publish
  const spec = validateFields(specInput, schema.fields, 'spec', false)
  return {
    kind: schema.kind,
    metadata: { name: value.metadata.name },
    spec: spec as AdminRouteBindingObject['spec'],
  }
}

export function hasLegacyPublishPreference(value: unknown) {
  return isRecord(value) && isRecord(value.spec) && hasOwn(value.spec, 'publish')
}

function yamlErrorMessage(error: unknown) {
  if (!(error instanceof Error)) return String(error)
  const linePosition = (error as Error & { linePos?: Array<{ line: number; col: number }> })
    .linePos?.[0]
  return linePosition
    ? `${error.message} (${linePosition.line}:${linePosition.col})`
    : error.message
}

export function parseRouteBindingYaml(value: string, schema: RouteBindingObjectSchema) {
  try {
    const parsed = parseYaml(value) as unknown
    const object = normaliseRouteBindingObject(parsed, schema)
    const defaultsApplied = stableStringify(parsed) !== stableStringify(object)
    return { object, error: '', defaultsApplied }
  } catch (error: unknown) {
    return { object: null, error: yamlErrorMessage(error), defaultsApplied: false }
  }
}

function stableStringify(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(stableStringify).join(',')}]`
  if (isRecord(value)) {
    return `{${Object.keys(value)
      .sort()
      .map((key) => `${JSON.stringify(key)}:${stableStringify(value[key])}`)
      .join(',')}}`
  }
  return JSON.stringify(value) ?? String(value)
}

export function stringifyRouteBindingYaml(value: AdminRouteBindingObject) {
  return stringifyYaml(value, { indent: 2, lineWidth: 0 })
}

function stringEnum(field?: RouteBindingFieldSchema) {
  return (field?.enum || []).filter((value): value is string => typeof value === 'string')
}

function property(field: RouteBindingFieldSchema | undefined, key: string) {
  return field?.properties?.[key]
}

export function routeEditorOptions(schema: RouteBindingObjectSchema): RouteEditorOptions {
  const entry = schema.fields.entry
  const target = schema.fields.target
  const params = schema.fields.params?.items
  return {
    entryProtocols: stringEnum(property(entry, 'protocol')),
    targetProtocols: stringEnum(property(target, 'protocol')),
    httpMethods: stringEnum(property(entry, 'method')),
    paramTypes: stringEnum(property(params, 'type')),
  }
}
