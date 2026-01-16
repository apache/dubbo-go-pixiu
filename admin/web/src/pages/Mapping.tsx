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
import { Table, Button, Space, Tag, Card, Input, Drawer, Popconfirm, Spin, Tooltip, Divider } from 'antd';
import { PlusOutlined, EyeOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { ColumnsType } from 'antd/es/table';
import {
  useResourceList,
  useResourceYaml,
  useCreateResource,
  useUpdateResource,
  useDeleteResource,
} from '../api';
import type { Resource } from '../api';
import { MappingWizard } from '../components/MappingWizard';
import { YamlEditor } from '../components/YamlEditor';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

export function Mapping() {
  const { t } = useTranslation();
  useDocumentTitle('menu.mapping');
  const [searchText, setSearchText] = useState('');

  // Wizard state
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [wizardMode, setWizardMode] = useState<'create' | 'edit'>('create');
  const [selectedResourceId, setSelectedResourceId] = useState<string | null>(null);

  // View modal state (separate from wizard)
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [viewResourceId, setViewResourceId] = useState<string | null>(null);

  // API hooks
  const { data: resources = [], isLoading } = useResourceList();
  // Separate queries for edit and view to avoid conflicts
  const { data: editResourceYaml, isLoading: isEditYamlLoading } = useResourceYaml(
    isWizardOpen && wizardMode === 'edit' ? selectedResourceId : null
  );
  const { data: viewResourceYaml, isLoading: isViewYamlLoading } = useResourceYaml(
    isViewModalOpen ? viewResourceId : null
  );
  const createResource = useCreateResource();
  const updateResource = useUpdateResource();
  const deleteResource = useDeleteResource();

  // Filter resources by search text
  const filteredResources = resources.filter((resource) =>
    resource.path?.toLowerCase().includes(searchText.toLowerCase()) ||
    resource.description?.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setWizardMode('create');
    setSelectedResourceId(null);
    setIsWizardOpen(true);
  };

  const handleWizardFinish = async (yaml: string) => {
    if (wizardMode === 'edit' && selectedResourceId) {
      await updateResource.mutateAsync({ id: selectedResourceId, content: yaml });
    } else {
      await createResource.mutateAsync(yaml);
    }
    setIsWizardOpen(false);
    setSelectedResourceId(null);
    setWizardMode('create');
  };

  const handleWizardClose = () => {
    setIsWizardOpen(false);
    setSelectedResourceId(null);
    setWizardMode('create');
  };

  const handleView = (record: Resource) => {
    setViewResourceId(record.id != null ? String(record.id) : null);
    setIsViewModalOpen(true);
  };

  const handleEdit = (record: Resource) => {
    setWizardMode('edit');
    setSelectedResourceId(record.id != null ? String(record.id) : null);
    setIsWizardOpen(true);
  };

  const handleDelete = async (record: Resource) => {
    if (record.id != null) {
      await deleteResource.mutateAsync(String(record.id));
    }
  };

  const handleViewModalClose = () => {
    setIsViewModalOpen(false);
    setViewResourceId(null);
  };

  const columns: ColumnsType<Resource> = [
    {
      title: t('mapping.path'),
      dataIndex: 'path',
      key: 'path',
      align: 'center',
      sorter: (a, b) => (a.path || '').localeCompare(b.path || ''),
    },
    {
      title: t('mapping.type'),
      dataIndex: 'type',
      key: 'type',
      align: 'center',
      render: (type) => (
        <Tag color="blue">{type ? t(`resourceTypes.${type}`, { defaultValue: type }) : '-'}</Tag>
      ),
    },
    {
      title: t('mapping.description'),
      dataIndex: 'description',
      key: 'description',
      align: 'center',
      ellipsis: true,
      render: (desc) => desc || '-',
    },
    {
      title: t('mapping.timeout'),
      dataIndex: 'timeout',
      key: 'timeout',
      align: 'center',
      render: (timeout) => {
        if (!timeout) return '-';
        const timeoutStr = String(timeout);
        // If timeout already has unit (ms, s, m, h), display as-is
        if (/^\d+(ms|s|m|h)$/.test(timeoutStr)) {
          return timeoutStr;
        }
        // Otherwise, assume it's nanoseconds (Go time.Duration)
        const ns = Number(timeout);
        if (isNaN(ns)) return timeoutStr;
        // Convert nanoseconds to human-readable format
        const seconds = ns / 1e9;
        if (seconds >= 60) {
          const minutes = Math.floor(seconds / 60);
          const remainingSecs = seconds % 60;
          return remainingSecs > 0 ? `${minutes}m${remainingSecs}s` : `${minutes}m`;
        }
        if (seconds >= 1) {
          return `${seconds}s`;
        }
        const ms = ns / 1e6;
        if (ms >= 1) {
          return `${ms}ms`;
        }
        return `${ns}ns`;
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
        title={<span className="text-lg font-semibold">{t('menu.mapping')}</span>}
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
          dataSource={filteredResources}
          loading={isLoading}
          rowKey={(record) => record.id != null ? String(record.id) : record.path}
          size="middle"
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
          }}
        />
      </Card>

      {/* Create/Edit Wizard */}
      <MappingWizard
        open={isWizardOpen}
        onClose={handleWizardClose}
        onFinish={handleWizardFinish}
        loading={createResource.isPending || updateResource.isPending}
        mode={wizardMode}
        initialData={wizardMode === 'edit' ? editResourceYaml || undefined : undefined}
        initialDataLoading={wizardMode === 'edit' && isEditYamlLoading}
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
          <YamlEditor value={viewResourceYaml || ''} height="calc(100vh - 150px)" readOnly showCopyButton />
        </Spin>
      </Drawer>
    </div>
  );
}
