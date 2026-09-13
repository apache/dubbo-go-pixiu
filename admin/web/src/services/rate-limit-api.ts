import { pixiuAdminApi } from '../api'
import { request } from './http'
export const rateLimitApi = {
  get: () => request<string>(pixiuAdminApi.rateLimit.detail),
  save: (content: string, method = 'PUT') =>
    request<void>(pixiuAdminApi.rateLimit.create, {
      method,
      body: new URLSearchParams({ content }),
    }),
  remove: () => request<void>(pixiuAdminApi.rateLimit.remove, { method: 'DELETE' }),
}
