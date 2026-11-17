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
import {
  emptyCluster,
  parseClusterYaml,
  stringifyClusterYaml,
  updateClusterField,
} from './ClusterEditorModel'

describe('Cluster form and YAML synchronization', () => {
  it('round trips a legacy Cluster while keeping the runtime ID outside the form', () => {
    const parsed = parseClusterYaml(
      'id: 17\nname: backend\ntype: Static\naddress: 127.0.0.1\nport: 20880\n',
      17,
    )
    expect(parsed).toEqual({
      kind: 'AdminCluster',
      metadata: { name: 'backend' },
      spec: {
        type: 'Static',
        endpoint: { address: '127.0.0.1', port: 20880 },
      },
    })
    expect(stringifyClusterYaml(parsed)).not.toContain('id:')
    expect(parseClusterYaml(stringifyClusterYaml(parsed))).toEqual(parsed)
    expect(() => parseClusterYaml('id: 18\nname: backend\n', 17)).toThrow(/Runtime ID/)
  })

  it('keeps incomplete and invalid draft values editable', () => {
    let object = emptyCluster()
    object = updateClusterField(object, 'spec.endpoint.port', 'bad-port')
    expect(object.spec.endpoint?.port).toBe('bad-port')
    expect(parseClusterYaml(stringifyClusterYaml(object))).toEqual(object)
    object = updateClusterField(object, 'spec.endpoint.port', '20880')
    expect(object.spec.endpoint?.port).toBe(20880)
  })

  it('rejects malformed and unknown YAML instead of dropping input', () => {
    expect(() => parseClusterYaml('name: backend\nport: [')).toThrow()
    expect(() => parseClusterYaml('name: backend\nunknown: value\n')).toThrow(/unknown/)
    expect(() =>
      parseClusterYaml('kind: AdminCluster\nmetadata:\n  name: backend\nspec:\n  id: 17\n'),
    ).toThrow(/spec.id/)
  })
})
