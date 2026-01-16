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

import { useCallback, useEffect } from 'react';
import {
  Form,
  Input,
  InputNumber,
  Select,
  Button,
  Row,
  Col,
  Card,
  Space,
  Drawer,
  Spin,
  message,
} from 'antd';
import { useTranslation } from 'react-i18next';
import yaml from 'js-yaml';
import { YamlEditor } from '../YamlEditor';
import { CLUSTER_TYPES } from '../../types/forms';

interface ClusterWizardProps {
  open: boolean;
  onClose: () => void;
  onFinish: (yaml: string) => Promise<void>;
  loading?: boolean;
  /** Loading state for initial data in edit mode */
  initialDataLoading?: boolean;
  /** Initial YAML data for edit mode */
  initialData?: string;
  /** Wizard mode: create or edit */
  mode?: 'create' | 'edit';
}

// Load balancing policy options (matching pkg/model.LbPolicyType values)
const LB_POLICIES = ['RoundRobin', 'Rand', 'RingHashing', 'MaglevHashing', 'WeightRandom'] as const;

// Parse YAML to extract cluster data (matching pkg/model.ClusterConfig struct)
function parseClusterYaml(yamlString: string): {
  name?: string;
  type?: string;
  lb_policy?: string;
  address?: string;
  port?: number;
} | null {
  try {
    const data = yaml.load(yamlString) as Record<string, unknown>;
    if (!data) return null;

    // Handle new format with endpoints array
    let address: string | undefined;
    let port: number | undefined;
    const endpoints = data.endpoints as Array<Record<string, unknown>> | undefined;
    if (endpoints && endpoints.length > 0) {
      const firstEndpoint = endpoints[0];
      const socketAddress = firstEndpoint?.socket_address as Record<string, unknown> | undefined;
      if (socketAddress) {
        address = socketAddress.address as string | undefined;
        port = socketAddress.port as number | undefined;
      }
    }

    return {
      name: data.name as string | undefined,
      type: data.type as string | undefined,
      lb_policy: data.lb_policy as string | undefined,
      address,
      port,
    };
  } catch {
    return null;
  }
}

export function ClusterWizard({ open, onClose, onFinish, loading, initialDataLoading, initialData, mode = 'create' }: ClusterWizardProps) {
  const { t } = useTranslation();
  const [form] = Form.useForm();

  // Initialize form with data when editing
  useEffect(() => {
    if (open && mode === 'edit' && initialData) {
      const parsed = parseClusterYaml(initialData);
      if (parsed) {
        form.setFieldsValue({
          name: parsed.name,
          type: parsed.type || 'Static',
          lb_policy: parsed.lb_policy || 'RoundRobin',
          address: parsed.address,
          port: parsed.port || 80,
        });
      }
    }
  }, [open, mode, initialData, form]);

  const resetState = useCallback(() => {
    form.resetFields();
  }, [form]);

  // Generate YAML matching pkg/model.ClusterConfig struct
  const generateYaml = useCallback(() => {
    const values = form.getFieldsValue();

    // Build config object for proper YAML serialization with special character escaping
    const clusterConfig: Record<string, unknown> = {
      name: values.name || '',
      type: values.type || 'Static',
      lb_policy: values.lb_policy || 'RoundRobin',
    };

    // Add endpoints array if address or port is provided
    if (values.address || values.port) {
      clusterConfig.endpoints = [
        {
          ID: `${values.name || 'endpoint'}-0`,
          socket_address: {
            address: values.address || '',
            port: values.port || 80,
          },
        },
      ];
    }

    return yaml.dump(clusterConfig, { indent: 2, lineWidth: -1, noRefs: true });
  }, [form]);

  const handleSubmit = useCallback(async () => {
    try {
      await form.validateFields();
      const yamlContent = generateYaml();
      await onFinish(yamlContent);
      resetState();
    } catch {
      message.warning(t('common.pleaseCompleteForm'));
    }
  }, [form, generateYaml, onFinish, resetState, t]);

  const handleClose = useCallback(() => {
    resetState();
    onClose();
  }, [resetState, onClose]);

  const drawerTitle = mode === 'edit' ? t('wizard.cluster.editTitle') : t('wizard.cluster.title');

  // Watch form values to trigger YAML preview update
  const formValues = Form.useWatch([], form);

  // Generate YAML preview - regenerate when form values change
  const yamlPreview = (() => {
    // Access formValues to ensure this updates when form changes
    if (formValues) {
      try {
        return generateYaml();
      } catch {
        return '';
      }
    }
    return '';
  })();

  // Section 1: Basic Info
  const basicInfoSection = (
    <Card title={t('wizard.cluster.step1')} size="small">
      <Row gutter={16}>
        <Col span={8}>
          <Form.Item
            name="name"
            label={t('wizard.cluster.clusterName')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Input placeholder={t('wizard.cluster.clusterNamePlaceholder')} disabled={mode === 'edit'} />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            name="type"
            label={t('wizard.cluster.clusterType')}
            rules={[{ required: true, message: t('common.required') }]}
            initialValue="Static"
          >
            <Select
              placeholder={t('wizard.cluster.clusterTypePlaceholder')}
              options={CLUSTER_TYPES.map((type) => ({
                label: t(`clusterTypes.${type}`, { defaultValue: type }),
                value: type,
              }))}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            name="lb_policy"
            label={t('wizard.cluster.lbPolicy')}
            rules={[{ required: true, message: t('common.required') }]}
            initialValue="RoundRobin"
          >
            <Select
              placeholder={t('wizard.cluster.lbPolicyPlaceholder')}
              options={LB_POLICIES.map((policy) => ({
                label: t(`lbPolicies.${policy}`, { defaultValue: policy }),
                value: policy,
              }))}
            />
          </Form.Item>
        </Col>
      </Row>
    </Card>
  );

  // Section 2: Endpoint
  const endpointSection = (
    <Card title={t('wizard.cluster.step2')} size="small">
      <Row gutter={16}>
        <Col span={16}>
          <Form.Item
            name="address"
            label={t('wizard.cluster.endpointAddress')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Input placeholder={t('wizard.cluster.endpointAddressPlaceholder')} />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            name="port"
            label={t('wizard.cluster.endpointPort')}
            rules={[{ required: true, message: t('common.required') }]}
            initialValue={80}
          >
            <InputNumber
              min={1}
              max={65535}
              style={{ width: '100%' }}
              placeholder={t('wizard.cluster.endpointPortPlaceholder')}
            />
          </Form.Item>
        </Col>
      </Row>
    </Card>
  );

  // Section 3: YAML Preview
  const yamlPreviewSection = (
    <Card title={t('wizard.cluster.step4')} size="small">
      <YamlEditor value={yamlPreview} readOnly height="200px" showCopyButton />
    </Card>
  );

  return (
    <Drawer
      open={open}
      onClose={handleClose}
      title={drawerTitle}
      styles={{ wrapper: { width: 720 } }}
      extra={
        <Space>
          <Button onClick={handleClose}>{t('common.cancel')}</Button>
          <Button type="primary" onClick={handleSubmit} loading={loading}>
            {mode === 'edit' ? t('common.save') : t('common.create')}
          </Button>
        </Space>
      }
    >
      <Spin spinning={!!initialDataLoading}>
        <Form form={form} layout="vertical">
          <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
            {basicInfoSection}
            {endpointSection}
            {yamlPreviewSection}
          </div>
        </Form>
      </Spin>
    </Drawer>
  );
}

export default ClusterWizard;
