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
import { Form, Input, InputNumber, Row, Col, Card, Drawer, Button, Space, Spin, message, Table, Popconfirm, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import yaml from 'js-yaml';
import { YamlEditor } from '../YamlEditor';
import { useClusterList, usePluginGroupList } from '../../api';
import type { RouteItem, HttpFilter, PluginGroup } from '../../api/types';
import { getPluginConfig } from '../../config/pluginConfigs';
import type { FieldConfig } from '../DualModeEditor';

// Protocol to connection manager mapping
const PROTOCOL_CONNECTION_MANAGER: Record<string, string> = {
  HTTP: 'dgp.filter.httpconnectionmanager',
  HTTPS: 'dgp.filter.httpconnectionmanager',
  HTTP2: 'dgp.filter.httpconnectionmanager',
  GRPC: 'dgp.filter.grpcconnectionmanager',
  TCP: 'dgp.filter.tcpproxy',
  TRIPLE: 'dgp.filter.dubboconnectionmanager',
  DUBBO: 'dgp.filter.dubboconnectionmanager',
};

// Filter options by protocol type
const FILTER_OPTIONS_BY_PROTOCOL: Record<string, string[]> = {
  HTTP: [
    'dgp.filter.http.apiconfig',
    'dgp.filter.http.httpproxy',
    'dgp.filter.http.cors',
    'dgp.filter.http.ratelimit',
    'dgp.filter.http.auth.jwt',
    'dgp.filter.http.timeout',
    'dgp.filter.http.accesslog',
    'dgp.filter.http.metric',
    'dgp.filter.http.recovery',
  ],
  GRPC: [
    'dgp.filter.http.grpcproxy',
    'dgp.filter.http.timeout',
    'dgp.filter.http.accesslog',
    'dgp.filter.http.metric',
    'dgp.filter.http.recovery',
  ],
  DUBBO: [
    'dgp.filter.http.dubboproxy',
    'dgp.filter.http.timeout',
    'dgp.filter.http.accesslog',
    'dgp.filter.http.metric',
    'dgp.filter.http.recovery',
  ],
  TCP: [
    'dgp.filter.network.ratelimit',
    'dgp.filter.network.accesslog',
    'dgp.filter.network.metric',
  ],
};

// Get filter options for a protocol
function getFilterOptionsForProtocol(protocol: string): string[] {
  if (protocol === 'HTTPS' || protocol === 'HTTP2') {
    return FILTER_OPTIONS_BY_PROTOCOL.HTTP;
  }
  if (protocol === 'TRIPLE') {
    return FILTER_OPTIONS_BY_PROTOCOL.DUBBO;
  }
  return FILTER_OPTIONS_BY_PROTOCOL[protocol] || FILTER_OPTIONS_BY_PROTOCOL.HTTP;
}

// Get connection manager for a protocol
function getConnectionManagerForProtocol(protocol: string): string {
  return PROTOCOL_CONNECTION_MANAGER[protocol] || PROTOCOL_CONNECTION_MANAGER.HTTP;
}

// Check if protocol supports routes (HTTP-based protocols)
function protocolSupportsRoutes(protocol: string): boolean {
  return ['HTTP', 'HTTPS', 'HTTP2', 'GRPC'].includes(protocol);
}

// Check if protocol supports filters
function protocolSupportsFilters(): boolean {
  // All protocols support some form of filters
  return true;
}

interface ListenerWizardProps {
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

interface RouteFormData {
  prefix: string;
  cluster: string;
}

// Protocol options for listener
const PROTOCOL_OPTIONS = [
  { value: 'HTTP', label: 'HTTP' },
  { value: 'HTTPS', label: 'HTTPS' },
  { value: 'HTTP2', label: 'HTTP/2' },
  { value: 'GRPC', label: 'gRPC' },
  { value: 'TCP', label: 'TCP' },
  { value: 'TRIPLE', label: 'Triple' },
];

// Parse YAML to extract listener data (matching pkg/model.Listener struct with filter_chains)
function parseListenerYaml(yamlString: string): {
  name?: string;
  protocolType?: string;
  address?: string;
  port?: number;
  routes?: RouteItem[];
  httpFilters?: HttpFilter[];
  pluginGroups?: string[];
} | null {
  try {
    const data = yaml.load(yamlString) as Record<string, unknown>;
    if (!data) return null;

    let address: string | undefined;
    let port: number | undefined;

    // Parse address.socket_address (underscore format per YAML tag)
    const addressObj = data.address as Record<string, unknown> | undefined;
    if (addressObj) {
      const socketAddr = addressObj.socket_address as Record<string, unknown> | undefined;
      if (socketAddr) {
        address = socketAddr.address as string | undefined;
        port = socketAddr.port as number | undefined;
      }
    }

    // Extract routes, http_filters and plugin_groups from filter_chains (Pixiu native format)
    const routes: RouteItem[] = [];
    const httpFilters: HttpFilter[] = [];
    let pluginGroups: string[] = [];
    const filterChains = data.filter_chains as Record<string, unknown> | undefined;
    if (filterChains) {
      const filters = filterChains.filters as Array<Record<string, unknown>> | undefined;
      if (filters) {
        for (const filter of filters) {
          const filterName = filter.name as string;
          // Look for any connection manager filter (http, grpc, dubbo)
          const isConnectionManager = [
            'dgp.filter.httpconnectionmanager',
            'dgp.filter.grpcconnectionmanager',
            'dgp.filter.dubboconnectionmanager',
          ].includes(filterName);

          if (isConnectionManager && filter.config) {
            const config = filter.config as Record<string, unknown>;

            // Extract routes
            const routeConfig = config.route_config as Record<string, unknown> | undefined;
            if (routeConfig && routeConfig.routes) {
              const routesArray = routeConfig.routes as Array<Record<string, unknown>>;
              for (const r of routesArray) {
                const match = r.match as Record<string, unknown> | undefined;
                const route = r.route as Record<string, unknown> | undefined;
                if (match && route) {
                  routes.push({
                    match: {
                      prefix: match.prefix as string || '/',
                    },
                    route: {
                      cluster: route.cluster as string || '',
                      cluster_not_found_response_code: route.cluster_not_found_response_code as number | undefined,
                    },
                  });
                }
              }
            }

            // Extract http_filters
            const httpFiltersArray = config.http_filters as Array<Record<string, unknown>> | undefined;
            if (httpFiltersArray) {
              for (const hf of httpFiltersArray) {
                httpFilters.push({
                  name: hf.name as string,
                  config: hf.config as Record<string, unknown> | undefined,
                });
              }
            }

            // Extract plugin_groups
            const pluginGroupsArray = config.plugin_groups as string[] | undefined;
            if (pluginGroupsArray) {
              pluginGroups = pluginGroupsArray;
            }

            break; // Found the httpconnectionmanager, no need to continue
          }
        }
      }
    }

    return {
      name: data.name as string | undefined,
      protocolType: data.protocol_type as string | undefined,
      address,
      port,
      routes: routes.length > 0 ? routes : undefined,
      httpFilters: httpFilters.length > 0 ? httpFilters : undefined,
      pluginGroups: pluginGroups.length > 0 ? pluginGroups : undefined,
    };
  } catch {
    return null;
  }
}

export function ListenerWizard({ open, onClose, onFinish, loading, initialDataLoading, initialData, mode = 'create' }: ListenerWizardProps) {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [routeForm] = Form.useForm();

  // Routes state
  const [routes, setRoutes] = useState<RouteItem[]>([]);
  const [routeModalOpen, setRouteModalOpen] = useState(false);
  const [editingRouteIndex, setEditingRouteIndex] = useState<number | null>(null);

  // HTTP Filters state
  const [httpFilters, setHttpFilters] = useState<HttpFilter[]>([]);
  const [filterModalOpen, setFilterModalOpen] = useState(false);
  const [filterForm] = Form.useForm();
  const [editingFilterIndex, setEditingFilterIndex] = useState<number | null>(null);
  const [selectedFilterName, setSelectedFilterName] = useState<string>('');

  // Plugin Groups state
  const [selectedPluginGroups, setSelectedPluginGroups] = useState<string[]>([]);

  // Fetch cluster list for route target selection
  const { data: clusters = [] } = useClusterList();

  // Fetch plugin group list
  const { data: pluginGroups = [] } = usePluginGroupList();

  // Initialize form with data when editing
  useEffect(() => {
    if (open && mode === 'edit' && initialData) {
      const parsed = parseListenerYaml(initialData);
      if (parsed) {
        form.setFieldsValue({
          name: parsed.name,
          protocolType: parsed.protocolType || 'HTTP',
          address: parsed.address || '0.0.0.0',
          port: parsed.port || 8080,
        });
        // Routes, httpFilters and pluginGroups are initialized separately to avoid effect cascade
        const routesToSet = parsed.routes || [];
        const filtersToSet = parsed.httpFilters || [];
        const pluginGroupsToSet = parsed.pluginGroups || [];
        // Use requestAnimationFrame to defer state update outside effect body
        requestAnimationFrame(() => {
          setRoutes(routesToSet);
          setHttpFilters(filtersToSet);
          setSelectedPluginGroups(pluginGroupsToSet);
        });
      }
    }
  }, [open, mode, initialData, form]);

  const resetState = useCallback(() => {
    form.resetFields();
    setRoutes([]);
    setHttpFilters([]);
    setSelectedPluginGroups([]);
    setRouteModalOpen(false);
    setEditingRouteIndex(null);
    setFilterModalOpen(false);
    setEditingFilterIndex(null);
    setSelectedFilterName('');
    filterForm.resetFields();
  }, [form, filterForm]);

  // Generate YAML matching pkg/model.Listener struct with filter_chains
  const generateYaml = useCallback(() => {
    const values = form.getFieldsValue();
    const protocol = values.protocolType || 'HTTP';
    const connectionManager = getConnectionManagerForProtocol(protocol);

    // Build the connection manager filter config
    const connManagerConfig: Record<string, unknown> = {};

    // Add route_config if protocol supports routes and routes exist
    if (protocolSupportsRoutes(protocol) && routes.length > 0) {
      connManagerConfig.route_config = {
        routes: routes.map((r) => ({
          match: { prefix: r.match.prefix },
          route: {
            cluster: r.route.cluster,
            ...(r.route.cluster_not_found_response_code
              ? { cluster_not_found_response_code: r.route.cluster_not_found_response_code }
              : {}),
          },
        })),
      };
    }

    // Add http_filters if protocol supports filters
    if (protocolSupportsFilters() && httpFilters.length > 0) {
      connManagerConfig.http_filters = httpFilters.map((f) => ({
        name: f.name,
        config: f.config || {},
      }));
    }

    // Add plugin_groups if selected
    if (selectedPluginGroups.length > 0) {
      connManagerConfig.plugin_groups = selectedPluginGroups;
    }

    // Build the listener object matching pkg/model.Listener (Pixiu native format)
    const listener: Record<string, unknown> = {
      name: values.name || '',
      protocol_type: protocol,
      address: {
        socket_address: {
          address: values.address || '0.0.0.0',
          port: values.port || 8080,
        },
      },
      filter_chains: {
        filters: [
          {
            name: connectionManager,
            config: connManagerConfig,
          },
        ],
      },
    };

    return yaml.dump(listener, { indent: 2, lineWidth: -1 });
  }, [form, routes, httpFilters, selectedPluginGroups]);

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

  // Route management handlers
  const handleAddRoute = () => {
    setEditingRouteIndex(null);
    routeForm.resetFields();
    setRouteModalOpen(true);
  };

  const handleEditRoute = (index: number) => {
    const route = routes[index];
    setEditingRouteIndex(index);
    routeForm.setFieldsValue({
      prefix: route.match.prefix || '/',
      cluster: route.route.cluster,
    });
    setRouteModalOpen(true);
  };

  const handleDeleteRoute = (index: number) => {
    setRoutes(routes.filter((_, i) => i !== index));
  };

  const handleRouteModalOk = async () => {
    try {
      const values = await routeForm.validateFields() as RouteFormData;
      const newRoute: RouteItem = {
        match: { prefix: values.prefix },
        route: { cluster: values.cluster },
      };

      if (editingRouteIndex !== null) {
        const newRoutes = [...routes];
        newRoutes[editingRouteIndex] = newRoute;
        setRoutes(newRoutes);
      } else {
        setRoutes([...routes, newRoute]);
      }
      setRouteModalOpen(false);
      routeForm.resetFields();
    } catch {
      message.warning(t('common.pleaseCompleteForm'));
    }
  };

  const handleRouteModalCancel = () => {
    setRouteModalOpen(false);
    routeForm.resetFields();
  };

  // Filter management handlers
  const handleAddFilter = () => {
    setEditingFilterIndex(null);
    setSelectedFilterName('');
    filterForm.resetFields();
    setFilterModalOpen(true);
  };

  const handleEditFilter = (index: number) => {
    const filter = httpFilters[index];
    setEditingFilterIndex(index);
    setSelectedFilterName(filter.name);
    filterForm.setFieldsValue({
      name: filter.name,
      ...filter.config,
    });
    setFilterModalOpen(true);
  };

  const handleDeleteFilter = (index: number) => {
    setHttpFilters(httpFilters.filter((_, i) => i !== index));
  };

  const handleFilterNameChange = (name: string) => {
    setSelectedFilterName(name);
    // Reset config fields when filter type changes
    filterForm.resetFields();
    filterForm.setFieldsValue({ name });
    // Set default values for new filter
    const config = getPluginConfig(name);
    if (config?.defaultValues) {
      filterForm.setFieldsValue(config.defaultValues);
    }
  };

  const handleFilterModalOk = async () => {
    try {
      const values = await filterForm.validateFields();
      const filterName = values.name as string;

      // Extract config fields (all fields except 'name')
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

      const newFilter: HttpFilter = {
        name: filterName,
        config: Object.keys(configData).length > 0 ? configData : undefined,
      };

      if (editingFilterIndex !== null) {
        const newFilters = [...httpFilters];
        newFilters[editingFilterIndex] = newFilter;
        setHttpFilters(newFilters);
      } else {
        setHttpFilters([...httpFilters, newFilter]);
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

  const drawerTitle = mode === 'edit' ? t('wizard.listener.editTitle') : t('wizard.listener.title');

  // Watch form values to trigger YAML preview update
  const formValues = Form.useWatch([], form);

  // Generate YAML preview
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

  // Get current protocol type for conditional rendering
  const currentProtocol = formValues?.protocolType || 'HTTP';
  const showRoutes = protocolSupportsRoutes(currentProtocol);
  const showFilters = protocolSupportsFilters();
  const availableFilterOptions = getFilterOptionsForProtocol(currentProtocol);

  // Route table columns
  const routeColumns = [
    {
      title: t('listener.route.match'),
      key: 'match',
      render: (_: unknown, record: RouteItem) => record.match.prefix || '-',
    },
    {
      title: t('listener.route.cluster'),
      dataIndex: ['route', 'cluster'],
      key: 'cluster',
    },
    {
      title: t('common.operation'),
      key: 'action',
      width: 100,
      render: (_: unknown, _record: RouteItem, index: number) => (
        <Space>
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => handleEditRoute(index)} />
          <Popconfirm
            title={t('common.deleteConfirm')}
            onConfirm={() => handleDeleteRoute(index)}
            okText={t('common.yes')}
            cancelText={t('common.no')}
          >
            <Button type="text" size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // Filter table columns
  const filterColumns = [
    {
      title: t('listener.filter.name'),
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => t(`plugins.${name}`, name),
    },
    {
      title: t('common.operation'),
      key: 'action',
      width: 100,
      render: (_: unknown, _record: HttpFilter, index: number) => (
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

  // Section 1: Basic Info
  const basicInfoSection = (
    <Card title={t('wizard.listener.step1')} size="small">
      <Row gutter={16}>
        <Col span={16}>
          <Form.Item
            name="name"
            label={t('wizard.listener.listenerName')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Input placeholder={t('wizard.listener.listenerNamePlaceholder')} disabled={mode === 'edit'} />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            name="protocolType"
            label={t('wizard.listener.protocol')}
            initialValue="HTTP"
          >
            <Select options={PROTOCOL_OPTIONS} />
          </Form.Item>
        </Col>
      </Row>
    </Card>
  );

  // Section 2: Network Config
  const networkConfigSection = (
    <Card title={t('wizard.listener.step2')} size="small">
      <Row gutter={16}>
        <Col span={16}>
          <Form.Item
            name="address"
            label={t('wizard.listener.bindAddress')}
            rules={[{ required: true, message: t('common.required') }]}
            initialValue="0.0.0.0"
          >
            <Input placeholder={t('wizard.listener.bindAddressPlaceholder')} />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            name="port"
            label={t('wizard.listener.listenPort')}
            rules={[{ required: true, message: t('common.required') }]}
            initialValue={8080}
          >
            <InputNumber min={1} max={65535} placeholder="8080" style={{ width: '100%' }} />
          </Form.Item>
        </Col>
      </Row>
    </Card>
  );

  // Section 3: Route Configuration
  const routeConfigSection = (
    <Card
      title={t('listener.route.title')}
      size="small"
      extra={
        <Button type="primary" size="small" icon={<PlusOutlined />} onClick={handleAddRoute}>
          {t('listener.route.add')}
        </Button>
      }
    >
      <Table
        columns={routeColumns}
        dataSource={routes}
        rowKey={(_, index) => `route-${index}`}
        size="small"
        pagination={false}
        locale={{ emptyText: t('listener.route.empty') }}
      />
    </Card>
  );

  // Section 4: Filters Configuration (dynamic title based on protocol)
  const filterConfigSection = (
    <Card
      title={t('listener.filter.titleWithProtocol', { protocol: currentProtocol })}
      size="small"
      extra={
        <Button type="primary" size="small" icon={<PlusOutlined />} onClick={handleAddFilter}>
          {t('listener.filter.add')}
        </Button>
      }
    >
      <Table
        columns={filterColumns}
        dataSource={httpFilters}
        rowKey={(_, index) => `filter-${index}`}
        size="small"
        pagination={false}
        locale={{ emptyText: t('listener.filter.empty') }}
      />
    </Card>
  );

  // Section 5: Plugin Groups Configuration
  const pluginGroupsSection = (
    <Card title={t('listener.pluginGroup.title')} size="small">
      <Select
        mode="multiple"
        style={{ width: '100%' }}
        placeholder={t('listener.pluginGroup.selectPlaceholder')}
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
          {t('listener.pluginGroup.empty')}
        </div>
      )}
    </Card>
  );

  // Section 6: YAML Preview
  const yamlPreviewSection = (
    <Card title={t('wizard.listener.step4')} size="small">
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
              {networkConfigSection}
              {showRoutes && routeConfigSection}
              {showFilters && filterConfigSection}
              {pluginGroupsSection}
              {yamlPreviewSection}
            </div>
          </Form>
        </Spin>
      </Drawer>

      {/* Route Edit Drawer */}
      <Drawer
        title={editingRouteIndex !== null ? t('listener.route.edit') : t('listener.route.add')}
        open={routeModalOpen}
        onClose={handleRouteModalCancel}
        width={400}
        destroyOnClose
        extra={
          <Space>
            <Button onClick={handleRouteModalCancel}>{t('common.cancel')}</Button>
            <Button type="primary" onClick={handleRouteModalOk}>{t('common.confirm')}</Button>
          </Space>
        }
      >
        <Form form={routeForm} layout="vertical">
          <Form.Item
            name="prefix"
            label={t('listener.route.matchValue')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Input placeholder="/" />
          </Form.Item>

          <Form.Item
            name="cluster"
            label={t('listener.route.cluster')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Select
              placeholder={t('listener.route.clusterPlaceholder')}
              showSearch
              optionFilterProp="label"
              options={clusters.map((c) => ({ label: c.name, value: c.name }))}
            />
          </Form.Item>
        </Form>
      </Drawer>

      {/* Filter Add/Edit Drawer */}
      <Drawer
        title={editingFilterIndex !== null ? t('listener.filter.edit') : t('listener.filter.add')}
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
            label={t('listener.filter.selectFilter')}
            rules={[{ required: true, message: t('common.required') }]}
          >
            <Select
              placeholder={t('listener.filter.selectFilterPlaceholder')}
              showSearch
              optionFilterProp="label"
              options={availableFilterOptions.map((name) => ({
                label: t(`plugins.${name}`, name),
                value: name,
              }))}
              onChange={handleFilterNameChange}
              disabled={editingFilterIndex !== null}
            />
          </Form.Item>

          {/* Filter configuration fields */}
          {filterConfigFields.length > 0 && (
            <Card title={t('listener.filter.config')} size="small" style={{ marginTop: 16 }}>
              {filterConfigFields.map((field) => renderFilterFormField(field))}
            </Card>
          )}

          {/* Show hint when no config available */}
          {selectedFilterName && filterConfigFields.length === 0 && (
            <div style={{ color: '#999', marginTop: 16 }}>
              {t('listener.filter.noConfigNeeded')}
            </div>
          )}
        </Form>
      </Drawer>
    </>
  );
}

export default ListenerWizard;
