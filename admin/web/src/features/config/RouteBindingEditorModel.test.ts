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

import { describe, expect, it } from 'vitest'
import type { RouteBindingObjectSchema } from '../../types/api'
import {
  createDefaultRouteBinding,
  hasLegacyPublishPreference,
  normaliseRouteBindingObject,
  parseRouteBindingYaml,
  routeEditorOptions,
} from './RouteBindingEditorModel'

const schema: RouteBindingObjectSchema = {
  kind: 'AdminRouteBinding',
  fields: {
    entry: {
      type: 'object',
      required: true,
      properties: {
        protocol: { type: 'string', default: 'http', enum: ['http'] },
        path: { type: 'string', required: true, default: '/api/v1/example', pattern: '^/' },
        method: { type: 'string', required: true, default: 'GET', enum: ['GET', 'POST'] },
      },
    },
    target: {
      type: 'object',
      required: true,
      properties: {
        protocol: { type: 'string', default: 'dubbo', enum: ['dubbo'] },
        application: { type: 'string', required: true },
        interface: { type: 'string', required: true },
        method: { type: 'string', required: true },
        version: { type: 'string', default: '' },
        group: { type: 'string', default: '' },
        cluster: { type: 'string', required: true },
      },
    },
    params: {
      type: 'array',
      default: [],
      items: {
        type: 'object',
        properties: {
          from: { type: 'string', required: true },
          to: { type: 'integer', required: true, minimum: 0 },
          type: { type: 'string', required: true, enum: ['string', 'int'] },
        },
      },
    },
    enabled: { type: 'boolean', default: true },
    extensions: { type: 'object', default: {}, properties: {} },
  },
}

function routeObject(spec: Record<string, unknown> = {}) {
  return {
    kind: 'AdminRouteBinding',
    metadata: { name: 'gender-user' },
    spec: {
      entry: { protocol: 'http', path: '/smoke/gender/:id', method: 'GET' },
      target: {
        protocol: 'dubbo',
        application: 'dubbo.io',
        interface: 'org.apache.dubbo.sample.UserProvider',
        method: 'GetGender',
        version: '',
        group: '',
        cluster: 'failover',
      },
      params: [{ from: 'uri.id', to: 0, type: 'int' }],
      enabled: true,
      extensions: {},
      ...spec,
    },
  }
}

describe('route binding editor model', () => {
  it('takes create defaults and enum options from the backend schema', () => {
    const object = createDefaultRouteBinding(schema)
    expect(object.spec.entry).toMatchObject({
      protocol: 'http',
      path: '/api/v1/example',
      method: 'GET',
    })
    expect(object.spec.target.protocol).toBe('dubbo')
    expect(object.spec.enabled).toBe(true)
    expect(routeEditorOptions(schema)).toEqual({
      entryProtocols: ['http'],
      targetProtocols: ['dubbo'],
      httpMethods: ['GET', 'POST'],
      paramTypes: ['string', 'int'],
    })
  })

  it('rejects scalar coercion instead of silently changing the YAML value', () => {
    expect(() => normaliseRouteBindingObject(routeObject({ enabled: 'false' }), schema)).toThrow(
      'spec.enabled: must be a boolean',
    )
    expect(() =>
      normaliseRouteBindingObject(
        routeObject({ params: [{ from: 'uri.id', to: '0', type: 'int' }] }),
        schema,
      ),
    ).toThrow('spec.params[0].to: must be an integer')
  })

  it('rejects unknown and deprecated fields that would be dropped on save', () => {
    expect(() =>
      normaliseRouteBindingObject(routeObject({ publish: { mode: 'draft' } }), schema),
    ).toThrow('spec.publish: is no longer supported')
    expect(() => normaliseRouteBindingObject(routeObject({ typo: true }), schema)).toThrow(
      'spec.typo: is not defined by the route schema',
    )
  })

  it('removes the obsolete no-op publish preference only when loading a stored draft', () => {
    const storedDraft = routeObject({ publish: { mode: 'draft' } })
    expect(hasLegacyPublishPreference(storedDraft)).toBe(true)
    const loaded = normaliseRouteBindingObject(storedDraft, schema, true)
    expect(loaded.spec).not.toHaveProperty('publish')
    expect(hasLegacyPublishPreference(loaded)).toBe(false)
  })

  it('reports type errors while parsing YAML and identifies applied schema defaults', () => {
    const invalid = parseRouteBindingYaml(
      `kind: AdminRouteBinding\nmetadata:\n  name: gender-user\nspec:\n  enabled: "false"\n`,
      schema,
    )
    expect(invalid.object).toBeNull()
    expect(invalid.error).toContain('spec.enabled: must be a boolean')

    const withDefaults = parseRouteBindingYaml(
      `kind: AdminRouteBinding\nmetadata:\n  name: gender-user\nspec:\n  entry:\n    protocol: http\n    path: /smoke/gender/:id\n    method: GET\n  target:\n    protocol: dubbo\n    application: dubbo.io\n    interface: org.apache.dubbo.sample.UserProvider\n    method: GetGender\n    cluster: failover\n`,
      schema,
    )
    expect(withDefaults.error).toBe('')
    expect(withDefaults.defaultsApplied).toBe(true)
    expect(withDefaults.object?.spec.enabled).toBe(true)
  })
})
