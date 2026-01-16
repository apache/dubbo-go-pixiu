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

import { useState, useEffect } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Select,
  Switch,
  Row,
  Col,
  message,
  Spin,
  Tabs,
  Descriptions,
  InputNumber,
  Divider,
} from 'antd';
import {
  SettingOutlined,
  DatabaseOutlined,
  SafetyCertificateOutlined,
  SaveOutlined,
  ReloadOutlined,
  CloudServerOutlined,
  ApiOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useThemeStore } from '../stores/theme';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import {
  useSystemSettings,
  useUpdateSystemSettings,
  useGatewayInfo,
  type SystemSettings,
} from '../api';

export function Settings() {
  const { t, i18n } = useTranslation();
  useDocumentTitle('menu.settings');
  const { isDark, toggleTheme } = useThemeStore();
  const [generalForm] = Form.useForm();
  const [gatewayForm] = Form.useForm();
  const [activeTab, setActiveTab] = useState('general');

  const { data: systemSettings, isLoading: isSettingsLoading, refetch } = useSystemSettings();
  const { data: gatewayInfo, isLoading: isGatewayLoading } = useGatewayInfo();
  const updateSettings = useUpdateSystemSettings();

  // Initialize forms with loaded data
  useEffect(() => {
    if (systemSettings) {
      generalForm.setFieldsValue(systemSettings);
    }
  }, [systemSettings, generalForm]);

  useEffect(() => {
    if (gatewayInfo) {
      gatewayForm.setFieldsValue(gatewayInfo);
    }
  }, [gatewayInfo, gatewayForm]);

  const handleGeneralSubmit = async (values: SystemSettings) => {
    try {
      await updateSettings.mutateAsync(values);
    } catch {
      // Error handled in hook
    }
  };

  const handleLanguageChange = (lang: string) => {
    i18n.changeLanguage(lang);
    message.success(t('common.success'));
  };

  const handleThemeChange = (checked: boolean) => {
    if ((checked && !isDark) || (!checked && isDark)) {
      toggleTheme();
    }
  };

  const tabItems = [
    {
      key: 'general',
      label: (
        <span className="flex items-center gap-2">
          <SettingOutlined />
          {t('settings.general')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Form
            form={generalForm}
            layout="vertical"
            onFinish={handleGeneralSubmit}
            initialValues={{
              language: i18n.language,
              darkMode: isDark,
              pageSize: 10,
              autoRefresh: false,
              refreshInterval: 30,
            }}
          >
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="language"
                  label={t('settings.language')}
                >
                  <Select
                    options={[
                      { value: 'zh', label: '中文' },
                      { value: 'en', label: 'English' },
                    ]}
                    onChange={handleLanguageChange}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="darkMode"
                  label={t('settings.darkMode')}
                  valuePropName="checked"
                >
                  <Switch onChange={handleThemeChange} />
                </Form.Item>
              </Col>
            </Row>

            <Divider />

            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="pageSize"
                  label={t('settings.pageSize')}
                >
                  <Select
                    options={[
                      { value: 10, label: '10' },
                      { value: 20, label: '20' },
                      { value: 50, label: '50' },
                      { value: 100, label: '100' },
                    ]}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="autoRefresh"
                  label={t('settings.autoRefresh')}
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="refreshInterval"
                  label={t('settings.refreshInterval')}
                >
                  <InputNumber
                    min={5}
                    max={300}
                    addonAfter={t('settings.seconds')}
                    style={{ width: '100%' }}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item>
              <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={updateSettings.isPending}>
                {t('common.save')}
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
    {
      key: 'gateway',
      label: (
        <span className="flex items-center gap-2">
          <CloudServerOutlined />
          {t('settings.gateway')}
        </span>
      ),
      children: (
        <Spin spinning={isGatewayLoading}>
          <div className="py-4">
            <Descriptions column={{ xs: 1, sm: 2 }} bordered size="small">
              <Descriptions.Item label={t('settings.gatewayVersion')}>
                {gatewayInfo?.version || '-'}
              </Descriptions.Item>
              <Descriptions.Item label={t('settings.gatewayStatus')}>
                {gatewayInfo?.status || t('common.running')}
              </Descriptions.Item>
              <Descriptions.Item label={t('settings.httpPort')}>
                {gatewayInfo?.httpPort || '8081'}
              </Descriptions.Item>
              <Descriptions.Item label={t('settings.xdsPort')}>
                {gatewayInfo?.xdsPort || '18000'}
              </Descriptions.Item>
              <Descriptions.Item label={t('settings.startTime')}>
                {gatewayInfo?.startTime || new Date().toLocaleString()}
              </Descriptions.Item>
              <Descriptions.Item label={t('settings.uptime')}>
                {gatewayInfo?.uptime || '-'}
              </Descriptions.Item>
            </Descriptions>

            <Divider />

            <Form
              form={gatewayForm}
              layout="vertical"
              initialValues={{
                logLevel: 'info',
                maxConnections: 10000,
                requestTimeout: 30000,
              }}
            >
              <Row gutter={24}>
                <Col xs={24} md={12}>
                  <Form.Item
                    name="logLevel"
                    label={t('settings.logLevel')}
                  >
                    <Select
                      options={[
                        { value: 'debug', label: 'Debug' },
                        { value: 'info', label: 'Info' },
                        { value: 'warn', label: 'Warn' },
                        { value: 'error', label: 'Error' },
                      ]}
                    />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item
                    name="maxConnections"
                    label={t('settings.maxConnections')}
                  >
                    <InputNumber
                      min={100}
                      max={100000}
                      style={{ width: '100%' }}
                    />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={24}>
                <Col xs={24} md={12}>
                  <Form.Item
                    name="requestTimeout"
                    label={t('settings.requestTimeout')}
                  >
                    <InputNumber
                      min={1000}
                      max={300000}
                      addonAfter="ms"
                      style={{ width: '100%' }}
                    />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item>
                <Button type="primary" htmlType="submit" icon={<SaveOutlined />}>
                  {t('common.save')}
                </Button>
              </Form.Item>
            </Form>
          </div>
        </Spin>
      ),
    },
    {
      key: 'api',
      label: (
        <span className="flex items-center gap-2">
          <ApiOutlined />
          {t('settings.apiSettings')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Form
            layout="vertical"
            initialValues={{
              apiBaseUrl: '/api',
              timeout: 30000,
              retryCount: 3,
            }}
          >
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="apiBaseUrl"
                  label={t('settings.apiBaseUrl')}
                >
                  <Input placeholder="/api" />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="timeout"
                  label={t('settings.apiTimeout')}
                >
                  <InputNumber
                    min={1000}
                    max={120000}
                    addonAfter="ms"
                    style={{ width: '100%' }}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="retryCount"
                  label={t('settings.retryCount')}
                >
                  <InputNumber min={0} max={10} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item>
              <Button type="primary" htmlType="submit" icon={<SaveOutlined />}>
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
          {t('settings.security')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Form
            layout="vertical"
            initialValues={{
              sessionTimeout: 7200,
              maxLoginAttempts: 5,
              lockoutDuration: 300,
              enableAuditLog: true,
            }}
          >
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="sessionTimeout"
                  label={t('settings.sessionTimeout')}
                >
                  <InputNumber
                    min={300}
                    max={86400}
                    addonAfter={t('settings.seconds')}
                    style={{ width: '100%' }}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="maxLoginAttempts"
                  label={t('settings.maxLoginAttempts')}
                >
                  <InputNumber min={1} max={20} style={{ width: '100%' }} />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="lockoutDuration"
                  label={t('settings.lockoutDuration')}
                >
                  <InputNumber
                    min={60}
                    max={3600}
                    addonAfter={t('settings.seconds')}
                    style={{ width: '100%' }}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="enableAuditLog"
                  label={t('settings.enableAuditLog')}
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item>
              <Button type="primary" htmlType="submit" icon={<SaveOutlined />}>
                {t('common.save')}
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
    {
      key: 'storage',
      label: (
        <span className="flex items-center gap-2">
          <DatabaseOutlined />
          {t('settings.storage')}
        </span>
      ),
      children: (
        <div className="py-4">
          <Descriptions column={{ xs: 1, sm: 2 }} bordered size="small">
            <Descriptions.Item label={t('settings.storageType')}>
              Memory
            </Descriptions.Item>
            <Descriptions.Item label={t('settings.dataPath')}>
              /var/pixiu/data
            </Descriptions.Item>
            <Descriptions.Item label={t('settings.configPath')}>
              /etc/pixiu/config
            </Descriptions.Item>
            <Descriptions.Item label={t('settings.logPath')}>
              /var/log/pixiu
            </Descriptions.Item>
          </Descriptions>

          <Divider />

          <Form
            layout="vertical"
            initialValues={{
              cacheEnabled: true,
              cacheTTL: 3600,
              maxCacheSize: 1024,
            }}
          >
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="cacheEnabled"
                  label={t('settings.cacheEnabled')}
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="cacheTTL"
                  label={t('settings.cacheTTL')}
                >
                  <InputNumber
                    min={60}
                    max={86400}
                    addonAfter={t('settings.seconds')}
                    style={{ width: '100%' }}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="maxCacheSize"
                  label={t('settings.maxCacheSize')}
                >
                  <InputNumber
                    min={64}
                    max={10240}
                    addonAfter="MB"
                    style={{ width: '100%' }}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item>
              <Button type="primary" htmlType="submit" icon={<SaveOutlined />}>
                {t('common.save')}
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center gap-3">
          <SettingOutlined className="text-2xl" />
          <h2 className="m-0">{t('menu.settings')}</h2>
        </div>
        <Button icon={<ReloadOutlined />} onClick={() => refetch()}>
          {t('common.refresh')}
        </Button>
      </div>

      <Spin spinning={isSettingsLoading}>
        <Card variant="borderless">
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            items={tabItems}
            tabPosition="left"
          />
        </Card>
      </Spin>
    </div>
  );
}
