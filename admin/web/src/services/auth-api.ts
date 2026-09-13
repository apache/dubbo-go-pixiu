import { API_BASE } from '../app/constants'
export async function login(username: string, password: string) {
  const body = new URLSearchParams({ username, password })
  const r = await fetch(API_BASE + '/login', { method: 'POST', body })
  const data = await r.json()
  if (data.code !== '10001') throw new Error(data.data || '登录失败')
  return data.data as { username: string; token: string }
}
