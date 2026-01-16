/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useEffect } from 'react';
import { Form, Input, Button, Checkbox, Space, Dropdown, Card, message, theme } from 'antd';
import {
  UserOutlined,
  LockOutlined,
  GlobalOutlined,
  SunOutlined,
  MoonOutlined,
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useThemeStore } from '../stores/theme';
import { useAuthStore } from '../stores/auth';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import type { MenuProps } from 'antd';

interface LoginForm {
  username: string;
  password: string;
  remember: boolean;
}

interface LocationState {
  from?: {
    pathname: string;
  };
}

export function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t, i18n } = useTranslation();
  const { isDark, toggleTheme } = useThemeStore();
  const { login, isAuthenticated, isLoading } = useAuthStore();
  const { token } = theme.useToken();
  useDocumentTitle('login.title');

  useEffect(() => {
    if (isAuthenticated) {
      const state = location.state as LocationState;
      const from = state?.from?.pathname || '/overview';
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, location.state]);

  const onFinish = async (values: LoginForm) => {
    const success = await login({
      username: values.username,
      password: values.password,
    });

    if (success) {
      message.success(t('login.success'));
      const state = location.state as LocationState;
      const from = state?.from?.pathname || '/overview';
      navigate(from, { replace: true });
    } else {
      message.error(t('login.failed'));
    }
  };

  const languageItems: MenuProps['items'] = [
    {
      key: 'en',
      label: t('language.en'),
      onClick: () => i18n.changeLanguage('en'),
    },
    {
      key: 'zh',
      label: t('language.zh'),
      onClick: () => i18n.changeLanguage('zh'),
    },
  ];

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: token.colorBgLayout,
      }}
    >
      {/* Top toolbar */}
      <div style={{ position: 'absolute', top: 24, right: 24 }}>
        <Space size={4}>
          <Dropdown menu={{ items: languageItems }} placement="bottom">
            <Button type="text" icon={<GlobalOutlined />} />
          </Dropdown>
          <Button
            type="text"
            icon={isDark ? <SunOutlined /> : <MoonOutlined />}
            onClick={toggleTheme}
          />
        </Space>
      </div>

      {/* Login card */}
      <Card style={{ width: 400 }}>
        {/* Logo and title */}
        <div style={{ marginBottom: 32, textAlign: 'center' }}>
          <div
            style={{
              width: 48,
              height: 48,
              borderRadius: 12,
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              background: token.colorPrimary,
              color: '#fff',
              fontWeight: 600,
              fontSize: 20,
              marginBottom: 16,
            }}
          >
            P
          </div>
          <h1 style={{ fontSize: 24, fontWeight: 600, marginBottom: 8 }}>
            {t('login.title')}
          </h1>
          <p style={{ fontSize: 14, color: token.colorTextSecondary }}>
            {t('login.subtitle')}
          </p>
        </div>

        {/* Login form */}
        <Form
          name="login"
          initialValues={{ remember: true }}
          onFinish={onFinish}
          layout="vertical"
          requiredMark={false}
        >
          <Form.Item
            name="username"
            label={t('login.username')}
            rules={[{ required: true, message: t('login.usernameRequired') }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder={t('login.username')}
              size="large"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label={t('login.password')}
            rules={[{ required: true, message: t('login.passwordRequired') }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder={t('login.password')}
              size="large"
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 24 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Form.Item name="remember" valuePropName="checked" noStyle>
                <Checkbox>{t('login.rememberMe')}</Checkbox>
              </Form.Item>
              <a href="#">{t('login.forgotPassword')}</a>
            </div>
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              block
              size="large"
              loading={isLoading}
            >
              {t('login.signIn')}
            </Button>
          </Form.Item>
        </Form>

        {/* Footer */}
        <div style={{ textAlign: 'center', fontSize: 12, color: token.colorTextSecondary }}>
          {t('login.footer')}
        </div>
      </Card>
    </div>
  );
}
