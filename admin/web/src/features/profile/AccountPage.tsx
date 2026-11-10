/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useEffect, useState } from 'react'
import { profileApi, type UserInfo } from '../../services/profile-api'

type Language = 'zh' | 'en'

export function AccountPage({ language }: { language: Language }) {
  const english = language === 'en'
  const [user, setUser] = useState<UserInfo | null>(null)
  const [loadError, setLoadError] = useState('')
  const [loading, setLoading] = useState(true)
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [feedback, setFeedback] = useState<{ kind: 'error' | 'success'; text: string } | null>(null)
  const [saving, setSaving] = useState(false)
  const [newUsername, setNewUsername] = useState('')
  const [newAccountPassword, setNewAccountPassword] = useState('')
  const [confirmAccountPassword, setConfirmAccountPassword] = useState('')
  const [registrationFeedback, setRegistrationFeedback] = useState<{
    kind: 'error' | 'success'
    text: string
  } | null>(null)
  const [registering, setRegistering] = useState(false)

  useEffect(() => {
    let active = true
    profileApi
      .session()
      .then((info) => {
        if (active) setUser(info)
      })
      .catch((error: unknown) => {
        if (active)
          setLoadError(
            error instanceof Error
              ? error.message
              : english
                ? 'Unable to load account information'
                : '无法加载账号信息',
          )
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [english])

  const submit = (event: React.FormEvent) => {
    event.preventDefault()
    setFeedback(null)
    if (!oldPassword || !newPassword || !confirmPassword) {
      setFeedback({
        kind: 'error',
        text: english ? 'Complete all password fields' : '请填写所有密码字段',
      })
      return
    }
    const newPasswordLength = Array.from(newPassword).length
    if (newPasswordLength < 8 || newPasswordLength > 128) {
      setFeedback({
        kind: 'error',
        text: english
          ? 'Use a new password between 8 and 128 characters'
          : '新密码长度需要为 8 到 128 个字符',
      })
      return
    }
    if (newPassword !== confirmPassword) {
      setFeedback({
        kind: 'error',
        text: english ? 'The new passwords do not match' : '两次输入的新密码不一致',
      })
      return
    }
    if (newPassword === oldPassword) {
      setFeedback({
        kind: 'error',
        text: english
          ? 'Choose a password different from the current one'
          : '新密码不能与当前密码相同',
      })
      return
    }

    setSaving(true)
    profileApi
      .password({ oldPassword, newPassword })
      .then(() => {
        setOldPassword('')
        setNewPassword('')
        setConfirmPassword('')
        setFeedback({ kind: 'success', text: english ? 'Password updated' : '密码已更新' })
      })
      .catch((error: unknown) =>
        setFeedback({
          kind: 'error',
          text:
            error instanceof Error
              ? error.message
              : english
                ? 'Unable to update password'
                : '无法修改密码',
        }),
      )
      .finally(() => setSaving(false))
  }

  const submitRegistration = (event: React.FormEvent) => {
    event.preventDefault()
    setRegistrationFeedback(null)
    if (!newUsername.trim() || !newAccountPassword || !confirmAccountPassword) {
      setRegistrationFeedback({
        kind: 'error',
        text: english ? 'Complete all account fields' : '请填写所有新账号字段',
      })
      return
    }
    const passwordLength = Array.from(newAccountPassword).length
    if (passwordLength < 8 || passwordLength > 128) {
      setRegistrationFeedback({
        kind: 'error',
        text: english
          ? 'Use a password between 8 and 128 characters'
          : '密码长度需要为 8 到 128 个字符',
      })
      return
    }
    if (newAccountPassword !== confirmAccountPassword) {
      setRegistrationFeedback({
        kind: 'error',
        text: english ? 'The passwords do not match' : '两次输入的密码不一致',
      })
      return
    }

    setRegistering(true)
    profileApi
      .register({ username: newUsername.trim(), password: newAccountPassword })
      .then(() => {
        setNewUsername('')
        setNewAccountPassword('')
        setConfirmAccountPassword('')
        setRegistrationFeedback({
          kind: 'success',
          text: english ? 'Administrator account created' : '管理员账号已创建',
        })
      })
      .catch((error: unknown) =>
        setRegistrationFeedback({
          kind: 'error',
          text:
            error instanceof Error
              ? error.message
              : english
                ? 'Unable to create account'
                : '无法创建账号',
        }),
      )
      .finally(() => setRegistering(false))
  }

  return (
    <div className="account-page">
      <div className="page-head">
        <div>
          <p className="eyebrow">{english ? 'ACCOUNT' : 'ACCOUNT CENTER'}</p>
          <h1>{english ? 'Account' : '个人中心'}</h1>
          <p className="muted">
            {english
              ? 'View your account and update your password.'
              : '查看账号信息并修改登录密码。'}
          </p>
        </div>
      </div>

      <div className={`account-grid${user?.role === 1 ? ' account-grid-admin' : ''}`}>
        <div className="account-column account-overview-column">
          <section className="panel account-card account-info-card">
            <div className="panel-head">
              <div>
                <h2>{english ? 'Account information' : '账号信息'}</h2>
                <span>{english ? 'Signed-in user' : '当前登录账号'}</span>
              </div>
            </div>
            {loading ? (
              <div className="loading-skeleton" />
            ) : loadError ? (
              <div className="account-feedback error">{loadError}</div>
            ) : (
              <dl className="account-details">
                <div>
                  <dt>{english ? 'Username' : '用户名'}</dt>
                  <dd>{user?.username || '—'}</dd>
                </div>
                <div>
                  <dt>{english ? 'Role' : '角色'}</dt>
                  <dd>
                    {user?.role === 1
                      ? english
                        ? 'Administrator'
                        : '管理员'
                      : english
                        ? 'User'
                        : '用户'}
                  </dd>
                </div>
              </dl>
            )}
          </section>

          {user?.role === 1 && (
            <section className="panel account-card account-registration-card">
              <div className="panel-head">
                <div>
                  <h2>{english ? 'Register administrator' : '注册管理员账号'}</h2>
                  <span>
                    {english
                      ? 'Only administrators can create accounts. New accounts also receive the administrator role.'
                      : '仅管理员可以创建账号；新账号也将获得管理员角色。'}
                  </span>
                </div>
              </div>
              <form className="account-form account-registration-form" onSubmit={submitRegistration}>
                <label>
                  {english ? 'Username' : '用户名'}
                  <input
                    autoComplete="username"
                    maxLength={255}
                    value={newUsername}
                    onChange={(event) => setNewUsername(event.target.value)}
                  />
                </label>
                <label>
                  {english ? 'Password' : '密码'}
                  <input
                    type="password"
                    autoComplete="new-password"
                    minLength={8}
                    maxLength={128}
                    value={newAccountPassword}
                    onChange={(event) => setNewAccountPassword(event.target.value)}
                  />
                </label>
                <label>
                  {english ? 'Confirm password' : '确认密码'}
                  <input
                    type="password"
                    autoComplete="new-password"
                    minLength={8}
                    maxLength={128}
                    value={confirmAccountPassword}
                    onChange={(event) => setConfirmAccountPassword(event.target.value)}
                  />
                </label>
                {registrationFeedback && (
                  <div className={`account-feedback ${registrationFeedback.kind}`}>
                    {registrationFeedback.text}
                  </div>
                )}
                <button className="primary account-submit" disabled={registering}>
                  {registering
                    ? english
                      ? 'Creating…'
                      : '创建中…'
                    : english
                      ? 'Create account'
                      : '创建账号'}
                </button>
              </form>
            </section>
          )}
        </div>

        <div className="account-column account-password-column">
          <section className="panel account-card account-password-card">
            <div className="panel-head">
              <div>
                <h2>{english ? 'Change password' : '修改密码'}</h2>
                <span>
                  {english
                    ? 'Enter your current password to save a new one.'
                    : '输入当前密码后即可设置新密码。'}
                </span>
              </div>
            </div>
            <form className="account-form" onSubmit={submit}>
              <label>
                {english ? 'Current password' : '当前密码'}
                <input
                  type="password"
                  autoComplete="current-password"
                  value={oldPassword}
                  onChange={(event) => setOldPassword(event.target.value)}
                />
              </label>
              <label>
                {english ? 'New password' : '新密码'}
                <input
                  type="password"
                  autoComplete="new-password"
                  minLength={8}
                  maxLength={128}
                  value={newPassword}
                  onChange={(event) => setNewPassword(event.target.value)}
                />
              </label>
              <label>
                {english ? 'Confirm new password' : '确认新密码'}
                <input
                  type="password"
                  autoComplete="new-password"
                  minLength={8}
                  maxLength={128}
                  value={confirmPassword}
                  onChange={(event) => setConfirmPassword(event.target.value)}
                />
              </label>
              {feedback && <div className={`account-feedback ${feedback.kind}`}>{feedback.text}</div>}
              <button className="primary account-submit" disabled={saving}>
                {saving
                  ? english
                    ? 'Updating…'
                    : '修改中…'
                  : english
                    ? 'Update password'
                    : '更新密码'}
              </button>
            </form>
          </section>
        </div>
      </div>
    </div>
  )
}
