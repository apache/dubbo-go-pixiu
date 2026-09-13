import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../../services/auth-api'
export function Login() {
  const [user, setUser] = useState('admin')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [language, setLanguage] = useState<'zh' | 'en'>('zh')
  const nav = useNavigate()
  const english = language === 'en'
  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!user || !password) {
      setError(english ? 'Please enter a username and password' : '请输入用户名和密码')
      return
    }
    setLoading(true)
    setError('')
    login(user, password)
      .then((data) => {
        localStorage.setItem('pixiu_session', JSON.stringify(data))
        nav('/gateway/overview', { replace: true })
      })
      .catch((e) => setError(e.message || (english ? 'Login failed' : '登录失败')))
      .finally(() => setLoading(false))
  }
  return (
    <div className="login-page">
      <div className="login-panel">
        <div className="login-language">
          <button className="lang-btn" onClick={() => setLanguage(english ? 'zh' : 'en')}>
            {english ? '中' : 'EN'}
          </button>
        </div>
        <div className="login-brand">
          <div className="brand-mark">P</div>
          <div>
            <strong>PIXIU</strong>
            <span>Gateway Admin</span>
          </div>
        </div>
        <div className="login-copy">
          <p className="eyebrow">{english ? 'GATEWAY ADMIN' : '网关管理平台'}</p>
          <h1>{english ? 'Sign in to Pixiu Admin' : '登录 Pixiu Admin'}</h1>
          <p>
            {english
              ? 'Manage API routes, clusters, listeners, and traffic policies.'
              : '管理 API 路由、集群、监听器与流量策略。'}
          </p>
        </div>
        <form onSubmit={submit}>
          <label>
            {english ? 'Username' : '用户名'}
            <input value={user} onChange={(e) => setUser(e.target.value)} autoComplete="username" />
          </label>
          <label>
            {english ? 'Password' : '密码'}
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
            />
          </label>
          {error && <div className="login-error">{error}</div>}
          <button className="primary login-submit" disabled={loading}>
            {loading ? (english ? 'Signing in…' : '登录中…') : english ? 'Sign in' : '登录'}
          </button>
        </form>
        <small className="login-hint">
          {english ? 'Use your Pixiu Admin account' : '请使用 Pixiu Admin 账号登录'}
        </small>
      </div>
    </div>
  )
}
