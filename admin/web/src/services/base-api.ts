import { pixiuAdminApi } from '../api'
import { request } from './http'
export const baseApi = {
  get: () => request<string>(pixiuAdminApi.base),
  save: (content: string) =>
    request<void>(pixiuAdminApi.base + '/', {
      method: 'POST',
      body: new URLSearchParams({ content }),
    }),
}
