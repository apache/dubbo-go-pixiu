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

import { useCallback, useEffect, useState } from 'react';
import { Form, Input, Select, Card, Drawer, Button, Space, Spin, message, Table, Popconfirm, InputNumber } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import yaml from 'js-yaml';
import { YamlEditor } from '../YamlEditor';
import { RESOURCE_TYPES } from '../../types/forms';
import { usePluginGroupList } from '../../api';
import type { Filter, PluginGroup } from '../../api/types';
import { getPluginConfig } from '../../config/pluginConfigs';
import type { FieldConfig } from '../DualModeEditor';

// Available filters for API mapping
const MAPPING_FILTER_OPTIONS = [
  'dgp.filter.http.ratelimit',
  'dgp.filter.http.auth.jwt',
  'dgp.filter.http.timeout',
  'dgp.filter.http.cors',
  'dgp.filter.http.header',
  'dgp.filter.http.accesslog',
  'dgp.filter.http.metric',
];

interface MappingWizardProps {
  open: boolean;
  onClose: () => void;
  onFinish: (yaml: string) => Promise<void>;
  loading?: boolean;
  /** Initial YAML data for edit mode */
  initialData?: string;
  /** Wizard mode: create or edit */
  mode?: 'create' | 'edit';
  /** Whether initial data is still loading */
  initialDataLoading?: boolean;
}

// Parse timeout value to milliseconds for form field
function parseTimeoutToMs(timeout: unknown): string | undefined {
  if (!timeout) return undefined;
  const timeoutStr = String(timeout);

  // If timeout already has unit (ms, s, m, h), convert to ms
  const unitMatch = timeoutStr.match(/^(\d+)(ms|s|m|h)$/);
  if (unitMatch) {
    const value = Number(unitMatch[1]);
    const unit = unitMatch[2];
    switch (unit) {
      case 'ms': return String(value);
      case 's': return String(value * 1000);
      case 'm': return String(value * 60 * 1000);
      case 'h': return String(value * 60 * 60 * 1000);
      default: return String(value);
    }
  }

  // Otherwise, assume it's nanoseconds (Go time.Duration)
  const ns = Number(timeout);
  if (!isNaN(ns)) {
    const ms = ns / 1e6;
    return String(Math.round(ms));
  }

  return timeoutStr;
}

// Parse YAML to extract mapping data
function parseMappingYaml(yamlString: string): {
  path?: string;
  type?: string;
  description?: string;
  timeout?: string;
  filters?: Filter[];
  pluginGroups?: string[];
} | null {
  try {
    const data = yaml.load(yamlString) as Record<string, unknown>;
    if (!data) return null;

    // Parse timeout to milliseconds for form field
    const timeout = parseTimeoutToMs(data.timeout);

    // Extract filters
    const filters: Filter[] = [];
    const filtersArray = data.filters as Array<Record<string, unknown>> | undefined;
    if (filtersArray) {
      for (const f of filtersArray) {
        filters.push({
          name: f.name as string,
          config: f.config as Record<string, unknown> | undefined,
        });
      }
    }

    // Extract plugin_groups
    const pluginGroups = data.plugin_groups as string[] | undefined;

    return {
      path: data.path as string | undefined,
      type: data.type as string | undefined,
      description: data.description as string | undefined,
      timeout,
      filters: filters.length > 0 ? filters : undefined,
      pluginGroups: pluginGroups && pluginGroups.length > 0 ? pluginGroups : undefined,
    };
  } catch {
    return null;
  }
}

export function MappingWizard({ open, onClose, onFinish, loading, initialData, mode = 'create', initialDataLoading }: MappingWizardProps) {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [filterForm] = Form.useForm();

  // Filters state
  const [filters, setFilters] = useState<Filter[]>([]);
  const [filterModalOpen, setFilterModalOpen] = useState(false);
  const [editingFilterIndex, setEditingFilterIndex] = useState<number | null>(null);
  const [selectedFilterName, setSelectedFilterName] = useState<string>('');

  // Plugin Groups state
  const [selectedPluginGroups, setSelectedPluginGroups] = useState<string[]>([]);

  // Fetch plugin group list
  const { data: pluginGroups = [] } = usePluginGroupList();

  // Initialize form with data when editing
  useEffect(() => {
    if (open && mode === 'edit' && initialData) {
      const parsed = parseMappingYaml(initialData);
      if (parsed) {
        form.setFieldsValue({
          path: parsed.path,
          description: parsed.description,
          type: parsed.type || 'restful',
          timeout: parsed.timeout,
        });
        // Set filters and plugin groups
        const filtersToSet = parsed.filters || [];
        const pluginGroupsToSet = parsed.pluginGroups || [];
        requestAnimationFrame(() => {
          setFilters(filtersToSet);
          setSelectedPluginGroups(pluginGroupsToSet);
        });
      }
    }
  }, [open, mode, initialData, form]);

  const resetState = useCallback(() => {
    form.resetFields();
    filterForm.resetFields();
    setFilters([]);
    setSelectedPluginGroups([]);
    setFilterModalOpen(false);
    setEditingFilterIndex(null);
    setSelectedFilterName('');
  }, [form, filterForm]);

  const generateYaml = useCallback(() => {
    const values = form.getFieldsValue();
    const resource: Record<string, unknown> = {
      path: values.path || '',
    };

    if (values.type) resource.type = values.type;
    if (values.description) resource.description = values.description;
    // timeout must include unit suffix (e.g., "300ms", "1s") for Go's time.ParseDuration
    if (values.timeout) resource.timeout = `${values.timeout}ms`;

    // Add filters if any
    if (filters.length > 0) {
      resource.filters = filters.map((f) => ({
        name: f.name,
        ...(f.config && Object.keys(f.config).length > 0 ? { config: f.config } : {}),
      }));
    }

    // Add plugin_groups if selected
    if (selectedPluginGroups.length > 0) {
      resource.plugin_groups = selectedPluginGroups;
    }

    if (values.type === 'restful') {
      resource.methods = [
        { httpVerb: 'GET' },
        { httpVerb: 'POST' },
        { httpVerb: 'PUT' },
        { httpVerb: 'DELETE' },
      ];
    }

    if (values.type === 'dubbo') {
      resource.interface = {
        name: '',
        version: '1.0.0',
        group: '',
      };
    }

    return yaml.dump(resource, { indent: 2, lineWidth: -1 });
  }, [form, filters, selectedPluginGroups]);

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

  // Filter management handlers
  const handleAddFilter = () => {
    setEditingFilterIndex(null);
    setSelectedFilterName('');
    filterForm.resetFields();
    setFilterModalOpen(true);
  };

  const handleEditFilter = (index: number) => {
    const filter = filters[index];
    setEditingFilterIndex(index);
    setSelectedFilterName(filter.name || '');
    filterForm.setFieldsValue({
      name: filter.name,
      ...filter.config,
    });
    setFilterModalOpen(true);
  };

  const handleDeleteFilter = (index: number) => {
    setFilters(filters.filter((_, i) => i !== index));
  };

  const handleFilterNameChange = (name: string) => {
    setSelectedFilterName(name);
    filterForm.resetFields();
    filterForm.setFieldsValue({ name });
    const config = getPluginConfig(name);
    if (config?.defaultValues) {
      filterForm.setFieldsValue(config.defaultValues);
    }
  };

  const handleFilterModalOk = async () => {
    try {
      const values = await filterForm.validateFields();
      const filterName = values.name as string;

      const configData: Record<string, unknown> = {};
      const pluginConfig = getPluginConfig(filterName);
      if (pluginConfig) {
        const allFields = [...pluginConfig.basicFields, ...(pluginConfig.advancedFields || [])];
        allFields.forEach((field) => {
          const value = values[field.name];
          if (value !== undefined && value !== '' && value !== null) {
            configData[field.name] = value;
          }
        });
      }

      const newFilter: Filter = {
        name: filterName,
        config: Object.keys(configData).length > 0 ? configData : undefined,
      };

      if (editingFilterIndex !== null) {
        const newFilters = [...filters];
        newFilters[editingFilterIndex] = newFilter;
        setFilters(newFilters);
      } else {
        setFilters([...filters, newFilter]);
      }

      setFilterModalOpen(false);
      setSelectedFilterName('');
      setEditingFilterIndex(null);
      filterForm.resetFields();
    } catch {
      message.warning(t('common.pleaseCompleteForm'));
    }
  };

  const handleFilterModalCancel = () => {
    setFilterModalOpen(false);
    setSelectedFilterName('');
    setEditingFilterIndex(null);
    filterForm.resetFields();
  };

  // Render form field for filter config
  const renderFilterFormField = useCallback(
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

  // Get current filter config fields
  const currentFilterConfig = selectedFilterName ? getPluginConfig(selectedFilterName) : undefined;
  const filterConfigFields = currentFilterConfig
    ? [...currentFilterConfig.basicFields, ...(currentFilterConfig.advancedFields || [])]
    : [];

  const drawerTitle = mode === 'edit' ? t('wizard.mapping.editTitle') : t('wizard.mapping.title');

  // Watch form values to trigger YAML preview update
  const formValues = Form.useWatch([], form);

  // Generate YAML preview - regenerate when form values change
  const yamlPreview = (() => {
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
    <Card title={t('wizard.mapping.step1')} size="small">
      <Form.Item
        name="path"
        label={t('wizard.mapping.routePath')}
        rules={[{ required: true, message: t('common.required') }]}
      >
        <Input placeholder={t('wizard.mapping.routePathPlaceholder')} />
      </Form.Item>
      <Form.Item name="description" label={t('wizard.mapping.description')}>
        <Input.TextArea rows={3} placeholder={t('wizard.mapping.descriptionPlaceholder')} />
      </Form.Item>
    </Card>
  );

  // Section 2: Route Type
  const routeTypeSection = (
    <Card title={t('wizard.mapping.step2')} size="small">
      <Form.Item
        name="type"
        label={t('wizard.mapping.routeType')}
        rules={[{ required: true, message: t('common.required') }]}
        initialValue="restful"
      >
        <Select
          placeholder={t('wizard.mapping.routeTypePlaceholder')}
          options={RESOURCE_TYPES.map((type) => ({
            label: t(`resourceTypes.${type}`, { defaultValue: type }),
            value: type,
          }))}
        />
      </Form.Item>
    </Card>
  );

  // Section 3: Advanced Config
  const advancedConfigSection = (
    <Card title={t('wizard.mapping.step3')} size="small">
      <Form.Item name="timeout" label={t('wizard.mapping.timeout')}>
        <Input placeholder={t('wizard.mapping.timeoutPlaceholder')} addonAfter="ms" />
      </Form.Item>
    </Card>
  );

  // Filter table columns
  const filterColumns = [
    {
      title: t('mapping.filter.name'),
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => t(`plugins.${name}`, name),
    },
    {
      title: t('common.operation'),
      key: 'action',
      width: 100,
      render: (_: unknown, _record: Filter, index: number) => (
        <Space>
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => handleEditFilter(index)} />
          <Popconfirm
            title={t('common.deleteConfirm')}
            onConfirm={() => handleDeleteFilter(index)}
            okText={t('common.yes')}
            cancelText={t('common.no')}
          >
            <Button type="text" size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // Section 4: Filters Configuration
  const filtersSection = (
    <Card
      title={t('mapping.filter.title')}
      size="small"
      extra={
        <Button type="primary" size="small" icon={<PlusOutlined />} onClick={handleAddFilter}>
          {t('mapping.filter.add')}
        </Button>
      }
    >
      <Table
        columns={filterColumns}
        dataSource={filters}
        rowKey={(_, index) => `filter-${index}`}
        size="small"
        pagination={false}
        locale={{ emptyText: t('mapping.filter.empty') }}
      />
    </Card>
  );

  // Section 5: Plugin Groups Configuration
  const pluginGroupsSection = (
    <Card title={t('mapping.pluginGroup.title')} size="small">
      <Select
        mode="multiple"
        style={{ width: '100%' }}
        placeholder={t('mapping.pluginGroup.selectPlaceholder')}
        value={selectedPluginGroups}
        onChange={setSelectedPluginGroups}
        options={pluginGroups.map((pg: PluginGroup) => ({
          label: pg.groupName,
          value: pg.groupName,
        }))}
        optionFilterProp="label"
        showSearch
      />
      {selectedPluginGroups.length === 0 && (
        <div style={{ color: '#999', marginTop: 8, fontSize: 12 }}>
          {t('mapping.pluginGroup.empty')}
        </div>
      )}
    </Card>
  );

  // Section 6: YAML Preview
  const yamlPreviewSection = (
    <Card title={t('wizard.mapping.step4')} size="small">
      <YamlEditor value={yamlPreview} readOnly height="300px" showCopyButton />
    </Card>
  );

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
        <Spin spinning={!!initialDataLoading}>
          <Form form={form} layout="vertical">
            <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
              {basicInfoSection}
              {routeTypeSection}
              {advancedConfigSection}
              {filtersSection}
              {pluginGroupsSection}
              {yamlPreviewSection}
            </div>
          </Form>
        </Spin>
      </Drawer>

      {/* Filter Add/Edit Drawer */}
      <Drawer
        title={editingFilterIndex !== null ? t('mapping.filter.edit') : t('mapping.filter.add')}
        open={filterModalOpen}
        onClose={handleFilterModalCancel}
        width={500}
        destroyOnClose
        extra={
          <Space>
            <Button onClick={handleFilterModalCancel}>{t('common.cancel')}</Button>
            <Button type="primary" onClick={handleFilterModalOk}>{t('common.confirm')}</Button>
          </Space>
        }
      >
        <Form form={filterForm} layout="vertical">
          <Form.Item
            name="name"
            label={t('mapping.filter.selectFilter')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Select
              placeholder={t('mapping.filter.selectFilterPlaceholder')}
              showSearch
              optionFilterProp="label"
              options={MAPPING_FILTER_OPTIONS.map((name) => ({
                label: t(`plugins.${name}`, name),
                value: name,
              }))}
              onChange={handleFilterNameChange}
              disabled={editingFilterIndex !== null}
            />
          </Form.Item>

          {/* Filter configuration fields */}
          {filterConfigFields.length > 0 && (
            <Card title={t('mapping.filter.config')} size="small" style={{ marginTop: 16 }}>
              {filterConfigFields.map((field) => renderFilterFormField(field))}
            </Card>
          )}

          {/* Show hint when no config available */}
          {selectedFilterName && filterConfigFields.length === 0 && (
            <div style={{ color: '#999', marginTop: 16 }}>
              {t('mapping.filter.noConfigNeeded')}
            </div>
          )}
        </Form>
      </Drawer>
    </>
  );
}

export default MappingWizard;
