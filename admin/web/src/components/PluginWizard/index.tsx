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

import { useState, useCallback, useMemo, useEffect } from 'react';
import yaml from 'js-yaml';
import {
  Form,
  Input,
  InputNumber,
  Select,
  Tabs,
  Card,
  Row,
  Col,
  Button,
  List,
  Tag,
  Space,
  Popconfirm,
  Empty,
  Drawer,
  Spin,
  message,
} from 'antd';
import {
  GlobalOutlined,
  ApiOutlined,
  SettingOutlined,
  ThunderboltOutlined,
  PlusOutlined,
  DeleteOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { YamlEditor } from '../YamlEditor';
import {
  getPluginConfig,
  type PluginCategory,
} from '../../config/pluginConfigs';
import {
  HTTP_FILTER_PLUGINS,
  NETWORK_FILTER_PLUGINS,
  RPC_FILTER_PLUGINS,
  OTHER_PLUGINS,
} from '../../types/forms';
import type { FieldConfig } from '../DualModeEditor';

interface PluginWizardProps {
  open: boolean;
  onClose: () => void;
  onFinish: (yaml: string) => Promise<void>;
  loading?: boolean;
  /** Whether initial data is still loading */
  initialDataLoading?: boolean;
  /** Initial YAML data for edit mode */
  initialData?: string;
  /** Wizard mode: create or edit */
  mode?: 'create' | 'edit';
}

// Parse YAML to extract plugin group data
function parsePluginGroupYaml(yamlString: string): {
  groupName?: string;
  plugins?: PluginItem[];
} | null {
  try {
    const data = yaml.load(yamlString) as Record<string, unknown>;
    if (!data) return null;

    const plugins: PluginItem[] = [];
    const rawPlugins = data.plugins as Array<Record<string, unknown>> | undefined;
    if (rawPlugins && Array.isArray(rawPlugins)) {
      rawPlugins.forEach((p) => {
        plugins.push({
          name: (p.name as string) || '',
          version: (p.version as string) || '1.0.0',
          priority: (p.priority as number) ?? 100,
          externalLookupName: p.externalLookupName as string | undefined,
          config: p.config as Record<string, unknown> | undefined,
        });
      });
    }

    return {
      groupName: data.groupName as string | undefined,
      plugins,
    };
  } catch {
    return null;
  }
}

interface PluginItem {
  name: string;
  version: string;
  priority: number;
  externalLookupName?: string;
  config?: Record<string, unknown>;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type PluginFormData = Record<string, any>;

const ALL_PLUGINS_BY_CATEGORY: Record<PluginCategory, readonly string[]> = {
  http: HTTP_FILTER_PLUGINS,
  network: NETWORK_FILTER_PLUGINS,
  rpc: RPC_FILTER_PLUGINS,
  other: OTHER_PLUGINS,
};

export function PluginWizard({ open, onClose, onFinish, loading, initialDataLoading, initialData, mode = 'create' }: PluginWizardProps) {
  const { t } = useTranslation();
  const [groupForm] = Form.useForm();
  const [pluginForm] = Form.useForm<PluginFormData>();

  // State for group
  const [groupName, setGroupName] = useState('');

  // State for plugins list
  const [plugins, setPlugins] = useState<PluginItem[]>([]);

  // Initialize form with data when editing
  // This is a valid pattern for synchronizing external data to local form state
  useEffect(() => {
    if (open && mode === 'edit' && initialData) {
      const parsed = parsePluginGroupYaml(initialData);
      if (parsed) {
        // eslint-disable-next-line react-hooks/set-state-in-effect -- Valid pattern for loading initial form data
        setGroupName(parsed.groupName || '');
        setPlugins(parsed.plugins || []);
        groupForm.setFieldsValue({
          groupName: parsed.groupName,
        });
      }
    }
  }, [open, mode, initialData, groupForm]);

  // State for plugin editor drawer
  const [isPluginDrawerOpen, setIsPluginDrawerOpen] = useState(false);
  const [editingPluginIndex, setEditingPluginIndex] = useState<number | null>(null);
  const [selectedPlugin, setSelectedPlugin] = useState<string>('');
  const [selectedCategory, setSelectedCategory] = useState<PluginCategory>('http');

  // State for YAML preview (editable)
  const [editedYaml, setEditedYaml] = useState<string | null>(null);

  const pluginConfig = selectedPlugin ? getPluginConfig(selectedPlugin) : undefined;
  const allFields = pluginConfig
    ? [...pluginConfig.basicFields, ...(pluginConfig.advancedFields || [])]
    : [];

  // Generate YAML from current state
  const generateYaml = useCallback(() => {
    const lines: string[] = [];
    lines.push(`groupName: "${groupName}"`);
    lines.push('plugins:');

    plugins.forEach((plugin) => {
      lines.push(`  - name: "${plugin.name}"`);
      lines.push(`    version: "${plugin.version}"`);
      lines.push(`    priority: ${plugin.priority}`);
      if (plugin.externalLookupName) {
        lines.push(`    externalLookupName: "${plugin.externalLookupName}"`);
      }
      if (plugin.config && Object.keys(plugin.config).length > 0) {
        lines.push('    config:');
        Object.entries(plugin.config).forEach(([key, value]) => {
          if (typeof value === 'string') {
            if (value.includes('\n')) {
              lines.push(`      ${key}: |`);
              value.split('\n').forEach((line) => lines.push(`        ${line}`));
            } else {
              lines.push(`      ${key}: "${value}"`);
            }
          } else {
            lines.push(`      ${key}: ${value}`);
          }
        });
      }
    });

    return lines.join('\n');
  }, [groupName, plugins]);

  // Compute YAML content - use edited version if user has edited, otherwise generate
  const yamlContent = useMemo(() => {
    if (editedYaml !== null) {
      return editedYaml;
    }
    return generateYaml();
  }, [editedYaml, generateYaml]);

  // Handler to update edited YAML
  const handleYamlChange = useCallback((value: string) => {
    setEditedYaml(value);
  }, []);

  const handlePluginSelect = useCallback(
    (pluginType: string) => {
      setSelectedPlugin(pluginType);
      pluginForm.resetFields();
      const config = getPluginConfig(pluginType);
      if (config?.defaultValues) {
        pluginForm.setFieldsValue(config.defaultValues as PluginFormData);
      }
    },
    [pluginForm]
  );

  const renderFormField = useCallback(
    (field: FieldConfig) => {
      if (field.hidden) return null;
      const rules = field.required ? [{ required: true, message: t('common.required') }] : [];

      switch (field.type) {
        case 'input':
          return (
            <Form.Item key={field.name} name={field.name} label={t(field.label)} rules={rules}>
              <Input placeholder={field.placeholder ? t(field.placeholder) : undefined} />
            </Form.Item>
          );
        case 'number':
          return (
            <Form.Item key={field.name} name={field.name} label={t(field.label)} rules={rules}>
              <InputNumber
                min={field.min}
                max={field.max}
                placeholder={field.placeholder ? t(field.placeholder) : undefined}
                style={{ width: '100%' }}
              />
            </Form.Item>
          );
        case 'select':
          return (
            <Form.Item key={field.name} name={field.name} label={t(field.label)} rules={rules}>
              <Select
                placeholder={field.placeholder ? t(field.placeholder) : undefined}
                options={field.options?.map((opt) => ({
                  label: t(String(opt.label), { defaultValue: String(opt.label) }),
                  value: opt.value,
                }))}
              />
            </Form.Item>
          );
        case 'textarea':
          return (
            <Form.Item key={field.name} name={field.name} label={t(field.label)} rules={rules}>
              <Input.TextArea rows={3} placeholder={field.placeholder ? t(field.placeholder) : undefined} />
            </Form.Item>
          );
        default:
          return null;
      }
    },
    [t]
  );

  // Open drawer to add new plugin
  const handleAddPlugin = () => {
    setEditingPluginIndex(null);
    setSelectedPlugin('');
    setSelectedCategory('http');
    pluginForm.resetFields();
    setIsPluginDrawerOpen(true);
  };

  // Open drawer to edit existing plugin
  const handleEditPlugin = (index: number) => {
    const plugin = plugins[index];
    setEditingPluginIndex(index);
    setSelectedPlugin(plugin.name);

    // Find category for the plugin
    for (const [cat, pluginList] of Object.entries(ALL_PLUGINS_BY_CATEGORY)) {
      if (pluginList.includes(plugin.name)) {
        setSelectedCategory(cat as PluginCategory);
        break;
      }
    }

    // Set form values
    pluginForm.setFieldsValue({
      version: plugin.version,
      priority: plugin.priority,
      externalLookupName: plugin.externalLookupName,
      ...plugin.config,
    });

    setIsPluginDrawerOpen(true);
  };

  // Delete plugin from list
  const handleDeletePlugin = (index: number) => {
    setPlugins((prev) => prev.filter((_, i) => i !== index));
  };

  // Save plugin to list
  const handleSavePlugin = async () => {
    if (!selectedPlugin) {
      message.warning(t('wizard.plugin.selectPluginFirst'));
      return;
    }

    try {
      await pluginForm.validateFields();
    } catch {
      return;
    }

    const formData = pluginForm.getFieldsValue();
    const config = getPluginConfig(selectedPlugin);

    // Build config object from form fields
    const pluginConfigData: Record<string, unknown> = {};
    if (config) {
      const configFields = [...config.basicFields, ...(config.advancedFields || [])];
      configFields.forEach((field) => {
        const value = formData[field.name];
        if (value !== undefined && value !== '' && value !== null) {
          pluginConfigData[field.name] = value;
        }
      });
    }

    const newPlugin: PluginItem = {
      name: selectedPlugin,
      version: formData.version || '1.0.0',
      priority: formData.priority ?? 100,
      externalLookupName: formData.externalLookupName || undefined,
      config: Object.keys(pluginConfigData).length > 0 ? pluginConfigData : undefined,
    };

    if (editingPluginIndex !== null) {
      // Update existing plugin
      setPlugins((prev) => prev.map((p, i) => (i === editingPluginIndex ? newPlugin : p)));
    } else {
      // Add new plugin
      setPlugins((prev) => [...prev, newPlugin]);
    }

    setIsPluginDrawerOpen(false);
    setSelectedPlugin('');
    pluginForm.resetFields();
  };

  // Validate and submit
  const handleSubmit = async () => {
    // Validate group name
    if (!groupName.trim()) {
      message.warning(t('wizard.plugin.groupNameRequired'));
      return;
    }

    // Validate at least one plugin
    if (plugins.length === 0) {
      message.warning(t('wizard.plugin.atLeastOnePlugin'));
      return;
    }

    await onFinish(yamlContent);

    // Reset state
    setGroupName('');
    setPlugins([]);
    setEditedYaml(null);
    groupForm.resetFields();
  };

  const handleClose = useCallback(() => {
    setGroupName('');
    setPlugins([]);
    setEditedYaml(null);
    setSelectedPlugin('');
    groupForm.resetFields();
    pluginForm.resetFields();
    onClose();
  }, [groupForm, pluginForm, onClose]);

  // Close drawer
  const handleCloseDrawer = () => {
    setIsPluginDrawerOpen(false);
    setSelectedPlugin('');
    pluginForm.resetFields();
  };

  // Plugin editor drawer content
  const pluginDrawerContent = (
    <div>
      <Tabs
        activeKey={selectedCategory}
        onChange={(key) => setSelectedCategory(key as PluginCategory)}
        items={[
          { key: 'http', label: t('wizard.plugin.categoryHttp'), icon: <GlobalOutlined /> },
          { key: 'network', label: t('wizard.plugin.categoryNetwork'), icon: <ApiOutlined /> },
          { key: 'rpc', label: t('wizard.plugin.categoryRpc'), icon: <ThunderboltOutlined /> },
          { key: 'other', label: t('wizard.plugin.categoryOther'), icon: <SettingOutlined /> },
        ]}
      />
      <Row gutter={[8, 8]} style={{ marginTop: 16, marginBottom: 16 }}>
        {ALL_PLUGINS_BY_CATEGORY[selectedCategory].map((pluginType) => {
          const isSelected = selectedPlugin === pluginType;
          return (
            <Col xs={12} sm={8} key={pluginType}>
              <Card
                size="small"
                hoverable
                onClick={() => handlePluginSelect(pluginType)}
                style={{
                  borderColor: isSelected ? '#1890ff' : undefined,
                  backgroundColor: isSelected ? '#e6f7ff' : undefined,
                  textAlign: 'center',
                }}
              >
                <div style={{ fontSize: 13 }}>
                  {t(`plugins.${pluginType}`, { defaultValue: pluginType.split('.').pop() })}
                </div>
              </Card>
            </Col>
          );
        })}
      </Row>

      {selectedPlugin && (
        <Form form={pluginForm} layout="vertical">
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="version" label={t('plugin.version')} initialValue="1.0.0">
                <Input placeholder="1.0.0" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="priority" label={t('plugin.priority')} initialValue={100}>
                <InputNumber min={0} max={10000} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="externalLookupName" label={t('plugin.external')}>
                <Input placeholder={t('wizard.plugin.externalPlaceholder')} />
              </Form.Item>
            </Col>
          </Row>
          {allFields.length > 0 && (
            <Card title={t('wizard.plugin.pluginConfig')} size="small">
              <Row gutter={16}>
                {allFields.map((field) => (
                  <Col span={12} key={field.name}>
                    {renderFormField(field)}
                  </Col>
                ))}
              </Row>
            </Card>
          )}
        </Form>
      )}
    </div>
  );

  // Group name section
  const groupNameSection = (
    <Card title={t('wizard.plugin.step1')} size="small">
      <Form form={groupForm} layout="vertical">
        <Form.Item
          name="groupName"
          label={t('wizard.plugin.groupName')}
          rules={[{ required: true, message: t('common.required') }]}
        >
          <Input
            placeholder={t('wizard.plugin.groupNamePlaceholder')}
            value={groupName}
            onChange={(e) => setGroupName(e.target.value)}
            disabled={mode === 'edit'}
          />
        </Form.Item>
        <div className="text-sm text-slate-500">{t('wizard.plugin.groupNameHint')}</div>
      </Form>
    </Card>
  );

  // Plugins list section
  const pluginsListSection = (
    <Card
      title={t('wizard.plugin.step2')}
      size="small"
      extra={
        plugins.length > 0 ? (
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddPlugin} size="small">
            {t('wizard.plugin.addPlugin')}
          </Button>
        ) : null
      }
    >
      {plugins.length === 0 ? (
        <Empty description={t('wizard.plugin.noPluginsYet')} image={Empty.PRESENTED_IMAGE_SIMPLE}>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddPlugin}>
            {t('wizard.plugin.addFirstPlugin')}
          </Button>
        </Empty>
      ) : (
        <List
          dataSource={plugins}
          renderItem={(plugin, index) => (
            <List.Item
              actions={[
                <Button
                  key="edit"
                  type="text"
                  icon={<EditOutlined />}
                  onClick={() => handleEditPlugin(index)}
                />,
                <Popconfirm
                  key="delete"
                  title={t('common.confirmDelete')}
                  onConfirm={() => handleDeletePlugin(index)}
                  okText={t('common.confirm')}
                  cancelText={t('common.cancel')}
                >
                  <Button type="text" danger icon={<DeleteOutlined />} />
                </Popconfirm>,
              ]}
            >
              <List.Item.Meta
                title={
                  <Space>
                    <span>
                      {t(`plugins.${plugin.name}`, { defaultValue: plugin.name.split('.').pop() })}
                    </span>
                    <Tag color="blue">v{plugin.version}</Tag>
                    <Tag color="gold">P:{plugin.priority}</Tag>
                  </Space>
                }
                description={<span className="text-xs text-slate-400">{plugin.name}</span>}
              />
            </List.Item>
          )}
        />
      )}
    </Card>
  );

  // YAML preview section
  const yamlPreviewSection = (
    <Card title={t('wizard.plugin.stepPreview')} size="small">
      <YamlEditor value={yamlContent} onChange={handleYamlChange} height="300px" showCopyButton />
    </Card>
  );

  const drawerTitle = mode === 'edit' ? t('wizard.plugin.editTitle') : t('wizard.plugin.title');

  return (
    <>
      <Drawer
        open={open}
        onClose={handleClose}
        title={drawerTitle}
        styles={{ wrapper: { width: 720 } }}
        extra={
          <Space>
            <Button onClick={handleClose}>{t('common.cancel')}</Button>
            <Button type="primary" onClick={handleSubmit} loading={loading} disabled={initialDataLoading}>
              {mode === 'edit' ? t('common.save') : t('common.create')}
            </Button>
          </Space>
        }
      >
        <Spin spinning={!!initialDataLoading} tip={t('common.loading')}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
            {groupNameSection}
            {pluginsListSection}
            {yamlPreviewSection}
          </div>
        </Spin>
      </Drawer>

      {/* Plugin Editor Drawer */}
      <Drawer
        title={editingPluginIndex !== null ? t('wizard.plugin.editPlugin') : t('wizard.plugin.addPlugin')}
        open={isPluginDrawerOpen}
        onClose={handleCloseDrawer}
        styles={{ wrapper: { width: 560 } }}
        zIndex={1010}
        extra={
          <Space>
            <Button onClick={handleCloseDrawer}>{t('common.cancel')}</Button>
            <Button type="primary" onClick={handleSavePlugin}>
              {t('common.confirm')}
            </Button>
          </Space>
        }
      >
        {pluginDrawerContent}
      </Drawer>
    </>
  );
}

export default PluginWizard;
