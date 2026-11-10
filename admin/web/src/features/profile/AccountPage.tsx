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
type Feedback = Readonly<{ kind: 'error' | 'success'; text: string }>
type AccountCopy = Readonly<{
  eyebrow: string
  title: string
  description: string
  accountInformation: string
  signedInUser: string
  username: string
  role: string
  administrator: string
  user: string
  loadError: string
  registerAdministrator: string
  registerDescription: string
  password: string
  confirmPassword: string
  accountFieldsRequired: string
  accountPasswordLength: string
  accountPasswordsMismatch: string
  accountCreated: string
  accountCreateError: string
  creatingAccount: string
  createAccount: string
  changePassword: string
  changePasswordDescription: string
  currentPassword: string
  newPassword: string
  confirmNewPassword: string
  passwordFieldsRequired: string
  newPasswordLength: string
  newPasswordsMismatch: string
  passwordMustChange: string
  passwordUpdated: string
  passwordUpdateError: string
  updatingPassword: string
  updatePassword: string
}>

const accountCopy: Record<Language, AccountCopy> = {
  zh: {
    eyebrow: 'ACCOUNT CENTER',
    title: '个人中心',
    description: '查看账号信息并修改登录密码。',
    accountInformation: '账号信息',
    signedInUser: '当前登录账号',
    username: '用户名',
    role: '角色',
    administrator: '管理员',
    user: '用户',
    loadError: '无法加载账号信息',
    registerAdministrator: '注册管理员账号',
    registerDescription: '仅管理员可以创建账号；新账号也将获得管理员角色。',
    password: '密码',
    confirmPassword: '确认密码',
    accountFieldsRequired: '请填写所有新账号字段',
    accountPasswordLength: '密码长度需要为 8 到 128 个字符',
    accountPasswordsMismatch: '两次输入的密码不一致', // NOSONAR -- localized validation message.
    accountCreated: '管理员账号已创建',
    accountCreateError: '无法创建账号',
    creatingAccount: '创建中…',
    createAccount: '创建账号',
    changePassword: '修改密码',
    changePasswordDescription: '输入当前密码后即可设置新密码。', // NOSONAR -- localized help text.
    currentPassword: '当前密码',
    newPassword: '新密码',
    confirmNewPassword: '确认新密码',
    passwordFieldsRequired: '请填写所有密码字段', // NOSONAR -- localized validation message.
    newPasswordLength: '新密码长度需要为 8 到 128 个字符',
    newPasswordsMismatch: '两次输入的新密码不一致', // NOSONAR -- localized validation message.
    passwordMustChange: '新密码不能与当前密码相同', // NOSONAR -- localized validation message.
    passwordUpdated: '密码已更新',
    passwordUpdateError: '无法修改密码',
    updatingPassword: '修改中…',
    updatePassword: '更新密码',
  },
  en: {
    eyebrow: 'ACCOUNT',
    title: 'Account',
    description: 'View your account and update your password.',
    accountInformation: 'Account information',
    signedInUser: 'Signed-in user',
    username: 'Username',
    role: 'Role',
    administrator: 'Administrator',
    user: 'User',
    loadError: 'Unable to load account information',
    registerAdministrator: 'Register administrator',
    registerDescription:
      'Only administrators can create accounts. New accounts also receive the administrator role.',
    password: 'Password',
    confirmPassword: 'Confirm password',
    accountFieldsRequired: 'Complete all account fields',
    accountPasswordLength: 'Use a password between 8 and 128 characters',
    accountPasswordsMismatch: 'The passwords do not match',
    accountCreated: 'Administrator account created',
    accountCreateError: 'Unable to create account',
    creatingAccount: 'Creating…',
    createAccount: 'Create account',
    changePassword: 'Change password',
    changePasswordDescription: 'Enter your current password to save a new one.',
    currentPassword: 'Current password',
    newPassword: 'New password',
    confirmNewPassword: 'Confirm new password',
    passwordFieldsRequired: 'Complete all password fields',
    newPasswordLength: 'Use a new password between 8 and 128 characters',
    newPasswordsMismatch: 'The new passwords do not match',
    passwordMustChange: 'Choose a password different from the current one',
    passwordUpdated: 'Password updated',
    passwordUpdateError: 'Unable to update password',
    updatingPassword: 'Updating…', // NOSONAR -- localized button label.
    updatePassword: 'Update password',
  },
}

type AccountPageProps = Readonly<{ language: Language }>

export function AccountPage({ language }: AccountPageProps) {
  const copy = accountCopy[language]
  const { user, loadError, loading } = useAccountSession(language)
  const gridClassName = user?.role === 1 ? 'account-grid account-grid-admin' : 'account-grid'

  return (
    <div className="account-page">
      <div className="page-head">
        <div>
          <p className="eyebrow">{copy.eyebrow}</p>
          <h1>{copy.title}</h1>
          <p className="muted">{copy.description}</p>
        </div>
      </div>

      <div className={gridClassName}>
        <div className="account-column account-overview-column">
          <AccountInfoCard copy={copy} user={user} loadError={loadError} loading={loading} />
          {user?.role === 1 && <RegistrationForm copy={copy} />}
        </div>
        <div className="account-column account-password-column">
          <PasswordForm copy={copy} />
        </div>
      </div>
    </div>
  )
}

function useAccountSession(language: Language) {
  const [user, setUser] = useState<UserInfo | null>(null)
  const [loadError, setLoadError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let active = true
    profileApi
      .session()
      .then((info) => {
        if (active) setUser(info)
      })
      .catch((error: unknown) => {
        if (active) setLoadError(errorMessage(error, accountCopy[language].loadError))
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [language])

  return { user, loadError, loading }
}

type AccountInfoCardProps = Readonly<{
  copy: AccountCopy
  user: UserInfo | null
  loadError: string
  loading: boolean
}>

function AccountInfoCard({ copy, user, loadError, loading }: AccountInfoCardProps) {
  const roleLabel = user?.role === 1 ? copy.administrator : copy.user

  return (
    <section className="panel account-card account-info-card">
      <div className="panel-head">
        <div>
          <h2>{copy.accountInformation}</h2>
          <span>{copy.signedInUser}</span>
        </div>
      </div>
      {loading && <div className="loading-skeleton" />}
      {!loading && loadError && <div className="account-feedback error">{loadError}</div>}
      {!loading && !loadError && (
        <dl className="account-details">
          <div>
            <dt>{copy.username}</dt>
            <dd>{user?.username || '—'}</dd>
          </div>
          <div>
            <dt>{copy.role}</dt>
            <dd>{roleLabel}</dd>
          </div>
        </dl>
      )}
    </section>
  )
}

type AccountFormProps = Readonly<{ copy: AccountCopy }>

function RegistrationForm({ copy }: AccountFormProps) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmation, setConfirmation] = useState('')
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [registering, setRegistering] = useState(false)

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    setFeedback(null)
    const validationMessage = validateRegistration(username, password, confirmation, copy)
    if (validationMessage) {
      setFeedback({ kind: 'error', text: validationMessage })
      return
    }

    setRegistering(true)
    try {
      await profileApi.register({ username: username.trim(), password })
      setUsername('')
      setPassword('')
      setConfirmation('')
      setFeedback({ kind: 'success', text: copy.accountCreated })
    } catch (error: unknown) {
      setFeedback({ kind: 'error', text: errorMessage(error, copy.accountCreateError) })
    } finally {
      setRegistering(false)
    }
  }

  return (
    <section className="panel account-card account-registration-card">
      <div className="panel-head">
        <div>
          <h2>{copy.registerAdministrator}</h2>
          <span>{copy.registerDescription}</span>
        </div>
      </div>
      <form className="account-form account-registration-form" onSubmit={submit}>
        <label>
          {copy.username}
          <input
            autoComplete="username"
            maxLength={255}
            value={username}
            onChange={(event) => setUsername(event.target.value)}
          />
        </label>
        <label>
          {copy.password}
          <input
            type="password"
            autoComplete="new-password"
            minLength={8}
            maxLength={128}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </label>
        <label>
          {copy.confirmPassword}
          <input
            type="password"
            autoComplete="new-password"
            minLength={8}
            maxLength={128}
            value={confirmation}
            onChange={(event) => setConfirmation(event.target.value)}
          />
        </label>
        {feedback && <div className={`account-feedback ${feedback.kind}`}>{feedback.text}</div>}
        <button type="submit" className="primary account-submit" disabled={registering}>
          {registering ? copy.creatingAccount : copy.createAccount}
        </button>
      </form>
    </section>
  )
}

function validateRegistration(
  username: string,
  password: string,
  confirmation: string,
  copy: AccountCopy,
) {
  if (!username.trim() || !password || !confirmation) return copy.accountFieldsRequired
  if (!isValidPasswordLength(password)) return copy.accountPasswordLength
  if (password !== confirmation) return copy.accountPasswordsMismatch
  return ''
}

function PasswordForm({ copy }: AccountFormProps) {
  const [currentPassword, setCurrentPassword] = useState('')
  const [password, setPassword] = useState('')
  const [confirmation, setConfirmation] = useState('')
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [saving, setSaving] = useState(false)

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    setFeedback(null)
    const validationMessage = validatePasswordChange(currentPassword, password, confirmation, copy)
    if (validationMessage) {
      setFeedback({ kind: 'error', text: validationMessage })
      return
    }

    setSaving(true)
    try {
      await profileApi.password({ oldPassword: currentPassword, newPassword: password })
      setCurrentPassword('')
      setPassword('')
      setConfirmation('')
      setFeedback({ kind: 'success', text: copy.passwordUpdated })
    } catch (error: unknown) {
      setFeedback({ kind: 'error', text: errorMessage(error, copy.passwordUpdateError) })
    } finally {
      setSaving(false)
    }
  }

  return (
    <section className="panel account-card account-password-card">
      <div className="panel-head">
        <div>
          <h2>{copy.changePassword}</h2>
          <span>{copy.changePasswordDescription}</span>
        </div>
      </div>
      <form className="account-form" onSubmit={submit}>
        <label>
          {copy.currentPassword}
          <input
            type="password"
            autoComplete="current-password"
            value={currentPassword}
            onChange={(event) => setCurrentPassword(event.target.value)}
          />
        </label>
        <label>
          {copy.newPassword}
          <input
            type="password"
            autoComplete="new-password"
            minLength={8}
            maxLength={128}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </label>
        <label>
          {copy.confirmNewPassword}
          <input
            type="password"
            autoComplete="new-password"
            minLength={8}
            maxLength={128}
            value={confirmation}
            onChange={(event) => setConfirmation(event.target.value)}
          />
        </label>
        {feedback && <div className={`account-feedback ${feedback.kind}`}>{feedback.text}</div>}
        <button type="submit" className="primary account-submit" disabled={saving}>
          {saving ? copy.updatingPassword : copy.updatePassword}
        </button>
      </form>
    </section>
  )
}

function validatePasswordChange(
  currentPassword: string,
  password: string,
  confirmation: string,
  copy: AccountCopy,
) {
  if (!currentPassword || !password || !confirmation) return copy.passwordFieldsRequired
  if (!isValidPasswordLength(password)) return copy.newPasswordLength
  if (password !== confirmation) return copy.newPasswordsMismatch
  if (password === currentPassword) return copy.passwordMustChange
  return ''
}

function isValidPasswordLength(password: string) {
  const length = Array.from(password).length
  return length >= 8 && length <= 128
}

function errorMessage(error: unknown, fallback: string) {
  if (error instanceof Error) return error.message
  return fallback
}
