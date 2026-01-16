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

import { useState } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Avatar,
  Row,
  Col,
  Divider,
  message,
  Spin,
  Tag,
  Descriptions,
  Tabs,
  Alert,
  Progress,
} from 'antd';
import {
  UserOutlined,
  MailOutlined,
  LockOutlined,
  SafetyCertificateOutlined,
  CheckCircleOutlined,
  CrownOutlined,
  KeyOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../stores/auth';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { changePassword, API_CODE } from '../api';

interface PasswordForm {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

export function Profile() {
  useDocumentTitle('menu.profile');
  const { t } = useTranslation();
  const { username } = useAuthStore();
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [activeTab, setActiveTab] = useState('profile');

  const handleProfileUpdate = () => {
    message.success(t('common.success'));
  };

  const handlePasswordChange = async (values: PasswordForm) => {
    setIsChangingPassword(true);
    try {
      const response = await changePassword({
        oldPassword: values.currentPassword,
        newPassword: values.newPassword,
      });

      if (response.code === API_CODE.SUCCESS) {
        message.success(t('profile.passwordChanged'));
        passwordForm.resetFields();
      } else {
        message.error(response.message || response.msg || t('profile.passwordChangeFailed'));
      }
    } catch {
      message.error(t('profile.passwordChangeFailed'));
    } finally {
      setIsChangingPassword(false);
    }
  };

  // Calculate password strength
  const calculatePasswordStrength = (password: string): number => {
    let strength = 0;
    if (password.length >= 6) strength += 25;
    if (password.length >= 10) strength += 25;
    if (/[A-Z]/.test(password)) strength += 25;
    if (/[0-9]/.test(password) && /[^A-Za-z0-9]/.test(password)) strength += 25;
    return strength;
  };

  const getPasswordStrengthText = (strength: number) => {
    if (strength < 50) return t('profile.passwordWeak');
    if (strength < 75) return t('profile.passwordMedium');
    return t('profile.passwordStrong');
  };

  const tabItems = [
    {
      key: 'profile',
      label: (
        <span className="flex items-center gap-2">
          <UserOutlined />
          {t('profile.basicInfo')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Form
            form={profileForm}
            layout="vertical"
            initialValues={{
              username: username || 'admin',
              email: `${username || 'admin'}@pixiu.apache.org`,
              displayName: username || 'Admin User',
            }}
            onFinish={handleProfileUpdate}
          >
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="username"
                  label={t('login.username')}
                  rules={[{ required: true }]}
                >
                  <Input
                    prefix={<UserOutlined />}
                    disabled
                    size="large"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="email"
                  label={t('profile.email')}
                  rules={[{ required: true, type: 'email' }]}
                >
                  <Input
                    prefix={<MailOutlined />}
                    size="large"
                  />
                </Form.Item>
              </Col>
            </Row>
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="displayName"
                  label={t('profile.displayName')}
                  rules={[{ required: true }]}
                >
                  <Input
                    prefix={<UserOutlined />}
                    size="large"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="phone"
                  label={t('profile.phone')}
                >
                  <Input
                    prefix={<span>+86</span>}
                    size="large"
                    placeholder={t('profile.phonePlaceholder')}
                  />
                </Form.Item>
              </Col>
            </Row>
            <Form.Item>
              <Button type="primary" htmlType="submit" icon={<CheckCircleOutlined />}>
                {t('common.save')}
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
    {
      key: 'security',
      label: (
        <span className="flex items-center gap-2">
          <SafetyCertificateOutlined />
          {t('profile.security')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Alert
            message={t('profile.passwordSecurityTitle')}
            description={t('profile.passwordSecurityDesc')}
            type="info"
            showIcon
            className="mb-6"
          />

          <Spin spinning={isChangingPassword}>
            <Form
              form={passwordForm}
              layout="vertical"
              onFinish={handlePasswordChange}
            >
              <Row gutter={24}>
                <Col xs={24} md={12}>
                  <Form.Item
                    name="currentPassword"
                    label={t('profile.currentPassword')}
                    rules={[{ required: true, message: t('profile.currentPasswordRequired') }]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      size="large"
                      placeholder={t('profile.currentPasswordPlaceholder')}
                    />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={24}>
                <Col xs={24} md={12}>
                  <Form.Item
                    name="newPassword"
                    label={t('profile.newPassword')}
                    rules={[
                      { required: true, message: t('profile.newPasswordRequired') },
                      { min: 6, message: t('profile.passwordMinLength') },
                    ]}
                  >
                    <Input.Password
                      prefix={<KeyOutlined />}
                      size="large"
                      placeholder={t('profile.newPasswordPlaceholder')}
                    />
                  </Form.Item>
                  <Form.Item noStyle shouldUpdate={(prev, curr) => prev.newPassword !== curr.newPassword}>
                    {({ getFieldValue }) => {
                      const password = getFieldValue('newPassword') || '';
                      const strength = calculatePasswordStrength(password);
                      if (!password) return null;
                      return (
                        <div className="mb-4">
                          <div className="text-xs mb-1">{t('profile.passwordStrength')}</div>
                          <Progress
                            percent={strength}
                            size="small"
                            status={strength < 50 ? 'exception' : strength < 75 ? 'normal' : 'success'}
                            showInfo={false}
                          />
                          <div className="text-xs mt-1">
                            {getPasswordStrengthText(strength)}
                          </div>
                        </div>
                      );
                    }}
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item
                    name="confirmPassword"
                    label={t('profile.confirmPassword')}
                    dependencies={['newPassword']}
                    rules={[
                      { required: true, message: t('profile.confirmPasswordRequired') },
                      ({ getFieldValue }) => ({
                        validator(_, value) {
                          if (!value || getFieldValue('newPassword') === value) {
                            return Promise.resolve();
                          }
                          return Promise.reject(new Error(t('profile.passwordMismatch')));
                        },
                      }),
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      size="large"
                      placeholder={t('profile.confirmPasswordPlaceholder')}
                    />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={isChangingPassword} icon={<SafetyCertificateOutlined />}>
                  {t('profile.updatePassword')}
                </Button>
              </Form.Item>
            </Form>
          </Spin>
        </div>
      ),
    },
    {
      key: 'activity',
      label: (
        <span className="flex items-center gap-2">
          <ClockCircleOutlined />
          {t('profile.activity')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label={t('profile.lastLogin')}>
              {new Date().toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label={t('profile.loginIp')}>
              127.0.0.1
            </Descriptions.Item>
            <Descriptions.Item label={t('profile.accountCreated')}>
              {new Date(Date.now() - 86400000 * 30).toLocaleDateString()}
            </Descriptions.Item>
            <Descriptions.Item label={t('profile.loginMethod')}>
              <Tag color="blue">{t('profile.password')}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label={t('profile.sessionStatus')}>
              <Tag color="success" icon={<CheckCircleOutlined />}>{t('profile.active')}</Tag>
            </Descriptions.Item>
          </Descriptions>
        </div>
      ),
    },
  ];

  return (
    <Row gutter={24}>
      {/* Left Column - User Card */}
      <Col xs={24} lg={8}>
        <Card className="text-center mb-4 lg:mb-0" variant="borderless">
          <div className="py-4">
            <Avatar
              size={120}
              icon={<UserOutlined />}
              className="mb-4"
            />
            <h2 className="text-xl font-semibold mb-1">{username || 'User'}</h2>
            <p className="mb-4">
              {username}@pixiu.apache.org
            </p>
            <Tag color="gold" icon={<CrownOutlined />} className="text-sm px-3 py-1">
              {t('profile.administrator')}
            </Tag>
          </div>

          <Divider />

          <Descriptions column={1} size="small" className="text-left">
            <Descriptions.Item label={t('profile.role')}>
              <Tag color="blue">{t('profile.admin')}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label={t('common.status')}>
              <Tag color="success" icon={<CheckCircleOutlined />}>{t('profile.active')}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label={t('profile.memberSince')}>
              {new Date(Date.now() - 86400000 * 30).toLocaleDateString()}
            </Descriptions.Item>
          </Descriptions>
        </Card>
      </Col>

      {/* Right Column - Tabs */}
      <Col xs={24} lg={16}>
        <Card variant="borderless">
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            items={tabItems}
            size="large"
          />
        </Card>
      </Col>
    </Row>
  );
}
