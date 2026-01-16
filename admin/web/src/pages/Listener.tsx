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
import { Table, Button, Space, Card, Input, Drawer, Popconfirm, Spin, Tooltip, Divider, Tag } from 'antd';
import { PlusOutlined, EyeOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { ColumnsType } from 'antd/es/table';
import {
  useListenerList,
  useListenerYaml,
  useCreateListener,
  useUpdateListener,
  useDeleteListener,
} from '../api';
import type { Listener as ListenerType, ListenerAddress, RouteItem } from '../api';
import { ListenerWizard } from '../components/ListenerWizard';
import { YamlEditor } from '../components/YamlEditor';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

// Extract routes from filter_chains structure (Pixiu native format)
// Routes are inside httpconnectionmanager filter's config.route_config.routes
function extractRoutesFromFilterChains(listener: ListenerType): RouteItem[] {
  const filterChains = listener.filter_chains;
  if (!filterChains || !filterChains.filters) return [];

  for (const filter of filterChains.filters) {
    // Look for httpconnectionmanager filter which contains routes
    if (filter.name === 'dgp.filter.httpconnectionmanager' && filter.config) {
      const config = filter.config as Record<string, unknown>;
      const routeConfig = config.route_config as Record<string, unknown> | undefined;
      if (routeConfig && routeConfig.routes) {
        return routeConfig.routes as RouteItem[];
      }
    }
  }
  return [];
}

export function Listener() {
  const { t } = useTranslation();
  useDocumentTitle('menu.listener');
  const [searchText, setSearchText] = useState('');

  // Wizard state
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [wizardMode, setWizardMode] = useState<'create' | 'edit'>('create');
  const [selectedListenerName, setSelectedListenerName] = useState<string | null>(null);

  // View modal state (separate from wizard)
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [viewListenerName, setViewListenerName] = useState<string | null>(null);

  // API hooks
  const { data: listeners = [], isLoading } = useListenerList();
  // Separate queries for edit and view to avoid conflicts
  const { data: editListenerYaml, isLoading: isEditYamlLoading } = useListenerYaml(
    isWizardOpen && wizardMode === 'edit' ? selectedListenerName : null
  );
  const { data: viewListenerYaml, isLoading: isViewYamlLoading } = useListenerYaml(
    isViewModalOpen ? viewListenerName : null
  );
  const createListener = useCreateListener();
  const updateListener = useUpdateListener();
  const deleteListener = useDeleteListener();

  // Filter listeners by search text
  const filteredListeners = listeners.filter((listener) =>
    listener.name?.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setWizardMode('create');
    setSelectedListenerName(null);
    setIsWizardOpen(true);
  };

  const handleWizardFinish = async (yaml: string) => {
    if (wizardMode === 'edit' && selectedListenerName) {
      await updateListener.mutateAsync({ name: selectedListenerName, content: yaml });
    } else {
      await createListener.mutateAsync(yaml);
    }
    setIsWizardOpen(false);
    setSelectedListenerName(null);
  };

  const handleWizardClose = () => {
    setIsWizardOpen(false);
    setSelectedListenerName(null);
  };

  const handleView = (record: ListenerType) => {
    setViewListenerName(record.name);
    setIsViewModalOpen(true);
  };

  const handleEdit = (record: ListenerType) => {
    setWizardMode('edit');
    setSelectedListenerName(record.name);
    setIsWizardOpen(true);
  };

  const handleDelete = async (record: ListenerType) => {
    await deleteListener.mutateAsync(record.name);
  };

  const handleViewModalClose = () => {
    setIsViewModalOpen(false);
    setViewListenerName(null);
  };

  const columns: ColumnsType<ListenerType> = [
    {
      title: t('listener.name'),
      dataIndex: 'name',
      key: 'name',
      align: 'center',
      sorter: (a, b) => (a.name || '').localeCompare(b.name || ''),
    },
    {
      title: t('listener.protocol'),
      dataIndex: 'protocol_type',
      key: 'protocol_type',
      align: 'center',
      render: (protocol: string | undefined) => {
        if (!protocol) return <Tag>HTTP</Tag>;
        const colorMap: Record<string, string> = {
          HTTP: 'blue',
          HTTPS: 'green',
          HTTP2: 'cyan',
          GRPC: 'purple',
          TCP: 'orange',
          TRIPLE: 'magenta',
        };
        return <Tag color={colorMap[protocol] || 'default'}>{protocol}</Tag>;
      },
    },
    {
      title: t('listener.address'),
      dataIndex: 'address',
      key: 'address',
      align: 'center',
      render: (address: ListenerAddress | undefined) => {
        if (address && typeof address === 'object') {
          const socketAddr = address.socket_address;
          if (socketAddr) {
            return `${socketAddr.address || '-'}:${socketAddr.port || '-'}`;
          }
        }
        return '-';
      },
    },
    {
      title: t('listener.route.title'),
      key: 'routes',
      render: (_, record) => {
        const routes = extractRoutesFromFilterChains(record);
        if (routes.length === 0) return '-';
        return (
          <Space direction="vertical" size={4}>
            {routes.map((route, index) => (
              <Tag key={index} color="blue">
                {route.match?.prefix || '/'} → {route.route?.cluster || '-'}
              </Tag>
            ))}
          </Space>
        );
      },
    },
    {
      title: t('common.operation'),
      key: 'action',
      align: 'center',
      render: (_, record) => (
        <>
          <Tooltip title={t('common.view')}>
            <Button type="text" size="small" icon={<EyeOutlined />} onClick={() => handleView(record)} />
          </Tooltip>
          <Divider orientation="vertical" />
          <Tooltip title={t('common.edit')}>
            <Button type="text" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)} />
          </Tooltip>
          <Divider orientation="vertical" />
          <Popconfirm
            title={t('common.deleteConfirm')}
            onConfirm={() => handleDelete(record)}
            okText={t('common.yes')}
            cancelText={t('common.no')}
          >
            <Tooltip title={t('common.delete')}>
              <Button type="text" size="small" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <div>
      <Card
        title={<span className="text-lg font-semibold">{t('listener.title')}</span>}
        className="shadow-sm"
        variant="borderless"
        extra={
          <Space>
            <Input.Search
              placeholder={t('common.search')}
              style={{ width: 220 }}
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
              onSearch={setSearchText}
              allowClear
            />
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreate}
            >
              {t('common.add')}
            </Button>
          </Space>
        }
      >
        <Table
          className="mt-6"
          columns={columns}
          dataSource={filteredListeners}
          loading={isLoading}
          rowKey="name"
          size="middle"
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
          }}
        />
      </Card>

      {/* Create/Edit Wizard */}
      <ListenerWizard
        open={isWizardOpen}
        onClose={handleWizardClose}
        onFinish={handleWizardFinish}
        loading={createListener.isPending || updateListener.isPending}
        initialDataLoading={wizardMode === 'edit' && isEditYamlLoading}
        mode={wizardMode}
        initialData={wizardMode === 'edit' ? editListenerYaml || undefined : undefined}
      />

      {/* View Drawer */}
      <Drawer
        title={t('common.view')}
        open={isViewModalOpen}
        onClose={handleViewModalClose}
        width={600}
        destroyOnClose
      >
        <Spin spinning={isViewYamlLoading}>
          <YamlEditor value={viewListenerYaml || ''} height="calc(100vh - 150px)" readOnly showCopyButton />
        </Spin>
      </Drawer>
    </div>
  );
}
