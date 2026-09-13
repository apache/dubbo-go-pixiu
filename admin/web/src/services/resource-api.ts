import { pixiuAdminApi } from '../api'
import { isNotFoundError, parseArrayResponse, request } from './http'
import type { JsonObject } from '../types/api'
export const resourceApi = {
  list: async () => {
    try {
      return parseArrayResponse<JsonObject>(await request<unknown>(pixiuAdminApi.resources.list))
    } catch (e: unknown) {
      if (isNotFoundError(e)) return []
      throw e
    }
  },
  detail: (id: string) =>
    request<unknown>(`${pixiuAdminApi.resources.detail}?resourceId=${encodeURIComponent(id)}`),
  create: (content: string) =>
    request<void>(pixiuAdminApi.resources.create, {
      method: 'POST',
      body: new URLSearchParams({ content }),
    }),
  update: (id: string, content: string) =>
    request<void>(`${pixiuAdminApi.resources.update}?resourceId=${encodeURIComponent(id)}`, {
      method: 'PUT',
      body: new URLSearchParams({ content }),
    }),
  remove: (id: string) =>
    request<void>(`${pixiuAdminApi.resources.remove}?resourceId=${encodeURIComponent(id)}`, {
      method: 'DELETE',
    }),
}
