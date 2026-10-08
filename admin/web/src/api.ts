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

/** Pixiu Admin API contract, aligned with admin/web controllers and API.md. */
export const pixiuAdminApi = {
  base: '/config/api/base',
  routeBindings: {
    schema: '/config/api/route/schema',
    list: '/config/api/route/list',
    detail: '/config/api/route/detail',
    create: '/config/api/route',
    update: '/config/api/route',
    remove: '/config/api/route',
    validate: '/config/api/route/validate',
    preview: '/config/api/route/preview',
    publish: '/config/api/route/publish',
    status: '/config/api/route/status',
    diff: '/config/api/route/diff',
  },
  clusters: {
    list: '/config/api/cluster/list',
    detail: '/config/api/cluster/detail',
    create: '/config/api/cluster',
    update: '/config/api/cluster',
    remove: '/config/api/cluster',
  },
  listeners: {
    list: '/config/api/listener/list',
    detail: '/config/api/listener/detail',
    create: '/config/api/listener',
    update: '/config/api/listener',
    remove: '/config/api/listener',
  },
  pluginGroups: {
    list: '/config/api/plugin_group/list',
    detail: '/config/api/plugin_group/detail',
    create: '/config/api/plugin_group',
    update: '/config/api/plugin_group',
    remove: '/config/api/plugin_group',
  },
  rateLimit: {
    detail: '/config/api/plugin/ratelimit',
    create: '/config/api/plugin/ratelimit/',
    update: '/config/api/plugin/ratelimit/',
    remove: '/config/api/plugin/ratelimit/',
  },
  opa: { policy: '/config/api/opa/policy' },
} as const
