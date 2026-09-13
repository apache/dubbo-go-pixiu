export type JsonObject = Record<string, unknown>
export function asJsonObject(value: unknown): JsonObject {
  return value && typeof value === 'object' && !Array.isArray(value) ? (value as JsonObject) : {}
}
export type Resource = JsonObject & {
  id?: string
  name?: string
  path?: string
  description?: string
  type?: string
  timeout?: string
  methods?: Method[]
}
export type Method = JsonObject & {
  id?: string
  httpVerb?: string
  resourcePath?: string
  onAir?: boolean
  timeout?: string
  inboundRequest?: unknown
  integrationRequest?: unknown
  plugins?: unknown
}
export type ApiEnvelope<T> = { code: string; data: T }
