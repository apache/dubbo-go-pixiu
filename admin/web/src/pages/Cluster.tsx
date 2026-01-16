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
  useClusterList,
  useClusterYaml,
  useCreateCluster,
  useUpdateCluster,
  useDeleteCluster,
} from '../api';
import type { Cluster as ClusterType } from '../api';
import { ClusterWizard } from '../components/ClusterWizard';
import { YamlEditor } from '../components/YamlEditor';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

export function Cluster() {
  const { t } = useTranslation();
  useDocumentTitle('menu.cluster');
  const [searchText, setSearchText] = useState('');

  // Wizard state
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [wizardMode, setWizardMode] = useState<'create' | 'edit'>('create');
  const [selectedClusterName, setSelectedClusterName] = useState<string | null>(null);

  // View modal state
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [viewClusterName, setViewClusterName] = useState<string | null>(null);

  // API hooks
  const { data: clusters = [], isLoading } = useClusterList();
  // Separate queries for edit and view to avoid cache conflicts
  const { data: editClusterYaml, isLoading: isEditYamlLoading } = useClusterYaml(
    isWizardOpen && wizardMode === 'edit' ? selectedClusterName : null,
    'edit'
  );
  const { data: viewClusterYaml, isLoading: isViewYamlLoading } = useClusterYaml(
    isViewModalOpen ? viewClusterName : null,
    'view'
  );
  const createCluster = useCreateCluster();
  const updateCluster = useUpdateCluster();
  const deleteCluster = useDeleteCluster();

  // Filter clusters by search text
  const filteredClusters = clusters.filter((cluster) =>
    cluster.name?.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setWizardMode('create');
    setSelectedClusterName(null);
    setIsWizardOpen(true);
  };

  const handleView = (record: ClusterType) => {
    setViewClusterName(record.name);
    setIsViewModalOpen(true);
  };

  const handleEdit = (record: ClusterType) => {
    setWizardMode('edit');
    setSelectedClusterName(record.name);
    setIsWizardOpen(true);
  };

  const handleDelete = async (record: ClusterType) => {
    await deleteCluster.mutateAsync(record.name);
  };

  const handleWizardFinish = async (yaml: string) => {
    if (wizardMode === 'create') {
      await createCluster.mutateAsync(yaml);
    } else if (selectedClusterName) {
      await updateCluster.mutateAsync({ name: selectedClusterName, content: yaml });
    }
    setIsWizardOpen(false);
    setSelectedClusterName(null);
  };

  const handleWizardClose = () => {
    setIsWizardOpen(false);
    setSelectedClusterName(null);
  };

  const handleViewModalClose = () => {
    setIsViewModalOpen(false);
    setViewClusterName(null);
  };

  const columns: ColumnsType<ClusterType> = [
    {
      title: t('cluster.name'),
      dataIndex: 'name',
      key: 'name',
      align: 'center',
      sorter: (a, b) => (a.name || '').localeCompare(b.name || ''),
    },
    {
      title: t('cluster.type'),
      dataIndex: 'type',
      key: 'type',
      align: 'center',
      render: (type) => (
        <Tag color="blue">{type ? t(`clusterTypes.${type}`, { defaultValue: type }) : '-'}</Tag>
      ),
    },
    {
      title: t('cluster.endpoints'),
      dataIndex: 'endpoints',
      key: 'endpoints',
      align: 'center',
      render: (endpoints: ClusterType['endpoints']) => {
        if (!endpoints || endpoints.length === 0) return '-';
        // Show first endpoint, with count if multiple
        const first = endpoints[0];
        const addr = first?.socket_address;
        const display = addr ? `${addr.address || '-'}:${addr.port || '-'}` : '-';
        if (endpoints.length > 1) {
          return (
            <Tooltip title={endpoints.map(e => {
              const a = e?.socket_address;
              return a ? `${a.address}:${a.port}` : '-';
            }).join(', ')}>
              <span>{display} <Tag color="blue">+{endpoints.length - 1}</Tag></span>
            </Tooltip>
          );
        }
        return display;
      },
    },
    {
      title: t('cluster.loadBalancer'),
      dataIndex: 'lb_policy',
      key: 'lb_policy',
      align: 'center',
      render: (lb_policy) => (
        <Tag color="green">{lb_policy ? t(`lbPolicies.${lb_policy}`, { defaultValue: lb_policy }) : '-'}</Tag>
      ),
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
        title={<span className="text-lg font-semibold">{t('cluster.title')}</span>}
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
          dataSource={filteredClusters}
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
      <ClusterWizard
        open={isWizardOpen}
        onClose={handleWizardClose}
        onFinish={handleWizardFinish}
        loading={createCluster.isPending || updateCluster.isPending}
        initialDataLoading={wizardMode === 'edit' && isEditYamlLoading}
        mode={wizardMode}
        initialData={wizardMode === 'edit' ? editClusterYaml || undefined : undefined}
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
          <YamlEditor
            value={viewClusterYaml || ''}
            height="calc(100vh - 150px)"
            readOnly
            showCopyButton
          />
        </Spin>
      </Drawer>
    </div>
  );
}
