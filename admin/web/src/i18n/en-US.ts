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

export const enUS = {
  published: 'Published',
  draft: 'Draft',
  paused: 'Paused',
  enabled: 'Enabled',
  disabled: 'Disabled',
  healthy: 'Healthy',
  error: 'Error',
  retry: 'Retry',
  cancel: 'Cancel',
  save: 'Save',
  delete: 'Delete',
  edit: 'Edit',
  create: 'Create',
  mock: 'MOCK',
  connected: 'CONNECTED',
} as const

export const moduleMessages = {
  capability: 'API capability',
  pending: 'Pending backend integration',
  description: 'This page is still a placeholder. List and edit operations are not connected yet.',
  contract: 'API paths',
  paths: 'Backend API reference',
  source: 'API paths for integration',
  createCluster: 'Create cluster',
  createListener: 'Create listener',
  createPlugin: 'Create plugin group',
  createLimit: 'Configure rate limit',
  createOpa: 'Create OPA policy',
  reference: 'API reference: ',
} as const
