// Shared API contract types, mirrored from the Go admin backend.

/** Unified response wrapper: admin/config.RetData. code is a string. */
export interface RetData<T = unknown> {
  code: string
  data: T
}

export const CODE_OK = '10001'

export interface LoginResult {
  username: string
  token: string
}

export interface UserInfo {
  id: number
  username: string
  role: string
}

/** pkg/config.Resource (legacy gateway resource). */
export interface Resource {
  id?: number
  type: string
  path: string
  timeout?: string | number
  description?: string
  filters?: unknown[]
  methods?: Method[]
  resources?: Resource[]
  headers?: Record<string, string>
}

/** pkg/config.Method (legacy method mapping). */
export interface Method {
  id?: number
  resourcePath?: string
  enable?: boolean
  timeout?: string | number
  httpVerb?: string
  inboundRequest?: Record<string, unknown>
  integrationRequest?: Record<string, unknown>
  [key: string]: unknown
}

export interface Cluster {
  id?: number
  name?: string
  type?: string
  address?: string
  port?: number
}

export interface Listener {
  name?: string
  address?: {
    socket_address?: { address?: string; port?: number }
    name?: string
  }
  route_config?: unknown
  http_filters?: unknown
}

// ---- Route binding (schema-driven POC) ----

export interface AdminObjectMeta {
  name?: string
}

export interface AdminObject {
  kind: string
  metadata: AdminObjectMeta
  spec: Record<string, unknown>
}

export type FieldType = 'string' | 'integer' | 'boolean' | 'object' | 'array' | 'map'

export interface UIHints {
  component?: string
  group?: string
  order?: number
  placeholder?: string
  advanced?: boolean
  options?: Record<string, unknown>
}

export interface FieldSchema {
  type: FieldType
  description?: string
  required?: boolean
  default?: unknown
  enum?: unknown[]
  pattern?: string
  minimum?: number
  properties?: Record<string, FieldSchema>
  items?: FieldSchema
  additionalProperties?: FieldSchema
  allowUnknown?: boolean
  ui?: UIHints
}

export interface ObjectSchema {
  kind: string
  description?: string
  fields: Record<string, FieldSchema>
}

export interface RouteBindingSummary {
  name: string
  status: 'draft' | 'published' | 'modified' | string
  draftRevision?: number
  publishedRevision?: number
  resourceId?: number
  methodId?: number
  updatedAt: string
}

export interface RouteBindingRecord {
  object: AdminObject
  resourceId: number
  methodId: number
  revision: number
  updatedAt: string
  resourceYaml: string
  methodYaml: string
  legacyYaml: string
}

export interface RouteBindingDetail {
  name: string
  draft?: RouteBindingRecord
  published?: RouteBindingRecord
}

export interface RouteBindingPreview {
  object: AdminObject
  resourceId?: number
  methodId?: number
  legacyYaml: string
  publishedLegacyYaml?: string
  diff?: string
  publishedRevision?: number
}
