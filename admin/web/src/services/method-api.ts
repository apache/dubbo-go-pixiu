import { pixiuAdminApi } from '../api'
import { isNotFoundError, parseArrayResponse, request } from './http'
import type { Method } from '../types/api'
export const methodApi = {
  list: async (resourceId: string) => {
    try {
      return parseArrayResponse<Method>(
        await request<unknown>(
          `${pixiuAdminApi.methods.list}?resourceId=${encodeURIComponent(resourceId)}`,
        ),
      )
    } catch (e: unknown) {
      if (isNotFoundError(e)) return []
      throw e
    }
  },
  detail: (resourceId: string, id: string) =>
    request<string>(
      `${pixiuAdminApi.methods.detail}?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,
    ),
  create: (resourceId: string, content: string) =>
    request<void>(`${pixiuAdminApi.methods.create}?resourceId=${encodeURIComponent(resourceId)}`, {
      method: 'POST',
      body: new URLSearchParams({ content }),
    }),
  update: (resourceId: string, id: string, content: string) =>
    request<void>(
      `${pixiuAdminApi.methods.update}?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,
      { method: 'PUT', body: new URLSearchParams({ content }) },
    ),
  remove: (resourceId: string, id: string) =>
    request<void>(
      `${pixiuAdminApi.methods.remove}?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,
      { method: 'DELETE' },
    ),
}
