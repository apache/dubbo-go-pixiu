import { pixiuAdminApi } from '../api'
import { isNotFoundError, parseArrayResponse, request } from './http'
import type { JsonObject } from '../types/api'
export const pluginGroupApi = {
  list: async () => {
    try {
      return parseArrayResponse<JsonObject>(await request<unknown>(pixiuAdminApi.pluginGroups.list))
    } catch (e: unknown) {
      if (isNotFoundError(e)) return []
      throw e
    }
  },
  detail: (name: string) =>
    request<unknown>(`${pixiuAdminApi.pluginGroups.detail}?name=${encodeURIComponent(name)}`),
  save: (content: string, method = 'POST') =>
    request<void>(pixiuAdminApi.pluginGroups.create, {
      method,
      body: new URLSearchParams({ content }),
    }),
  remove: (name: string) =>
    request<void>(`${pixiuAdminApi.pluginGroups.remove}?name=${encodeURIComponent(name)}`, {
      method: 'DELETE',
    }),
}
