export type PixiuSession = { username: string; token: string }

export function readSession(): PixiuSession | null {
  const raw = localStorage.getItem('pixiu_session')
  if (!raw) return null
  try {
    const session = JSON.parse(raw) as Partial<PixiuSession>
    if (
      typeof session.username !== 'string' ||
      !session.username.trim() ||
      typeof session.token !== 'string' ||
      !session.token.trim()
    )
      return null
    return { username: session.username, token: session.token }
  } catch {
    return null
  }
}

export function clearSession() {
  localStorage.removeItem('pixiu_session')
}

export function redirectToLogin() {
  if (window.location.pathname !== '/login') window.location.assign('/login')
}
