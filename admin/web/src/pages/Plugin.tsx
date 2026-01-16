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
  Row,
  Col,
  Card,
  Button,
  Empty,
  Spin,
  Space,
  Tag,
  Collapse,
  Popconfirm,
  Descriptions,
  Badge,
  Tooltip,
  Input,
  theme,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  AppstoreOutlined,
  SettingOutlined,
  ThunderboltOutlined,
  CodeOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { PluginWizard } from '../components/PluginWizard';
import {
  usePluginGroupList,
  usePluginGroupYaml,
  useCreatePluginGroup,
  useUpdatePluginGroup,
  useDeletePluginGroup,
} from '../api';
import type { PluginGroupItem } from '../api';

export function Plugin() {
  const { t } = useTranslation();
  const { token } = theme.useToken();
  useDocumentTitle('menu.plugin');
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<string | null>(null);
  const [searchText, setSearchText] = useState('');

  const { data: pluginGroups, isLoading, error, refetch } = usePluginGroupList();

  // Filter plugin groups by search text (group name or plugin name)
  const filteredPluginGroups = pluginGroups?.filter((group) =>
    group.groupName?.toLowerCase().includes(searchText.toLowerCase()) ||
    group.plugins?.some((plugin) => plugin.name?.toLowerCase().includes(searchText.toLowerCase()))
  );
  // Use usePluginGroupYaml to get YAML string for editing (same pattern as Cluster/Listener)
  const { data: editPluginGroupYaml, isLoading: isEditYamlLoading } = usePluginGroupYaml(
    isWizardOpen && editingGroup ? editingGroup : null
  );
  const createMutation = useCreatePluginGroup();
  const updateMutation = useUpdatePluginGroup();
  const deleteMutation = useDeletePluginGroup();

  const handleCreate = () => {
    setEditingGroup(null);
    setIsWizardOpen(true);
  };

  const handleWizardFinish = async (yaml: string) => {
    try {
      if (editingGroup) {
        await updateMutation.mutateAsync({ name: editingGroup, content: yaml });
      } else {
        await createMutation.mutateAsync(yaml);
      }
      setIsWizardOpen(false);
      setEditingGroup(null);
    } catch {
      // Error handled in hook
    }
  };

  const handleEdit = (groupName: string) => {
    setEditingGroup(groupName);
    setIsWizardOpen(true);
  };

  const handleDelete = async (groupName: string) => {
    try {
      await deleteMutation.mutateAsync(groupName);
    } catch {
      // Error handled in hook
    }
  };

  const handleWizardClose = () => {
    setIsWizardOpen(false);
    setEditingGroup(null);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-red-500 mb-4">{t('common.loadError')}</p>
        <Button onClick={() => refetch()}>{t('common.retry')}</Button>
      </div>
    );
  }

  return (
    <div>
      <Card
        title={<span className="text-lg font-semibold">{t('plugin.title')}</span>}
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
            <Button icon={<ReloadOutlined />} onClick={() => refetch()}>
              {t('common.refresh')}
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              {t('common.add')}
            </Button>
          </Space>
        }
      >
        {!filteredPluginGroups || filteredPluginGroups.length === 0 ? (
          <div className="mt-6">
            <Empty
              image={<AppstoreOutlined className="text-6xl text-slate-300" />}
              description={
                <div>
                  <p className="text-slate-500 mb-2">{t('plugin.noPlugins')}</p>
                  <p className="text-slate-400 text-sm">
                    {t('plugin.emptyDescription')}
                  </p>
                </div>
              }
            >
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
                {t('plugin.createFirst')}
              </Button>
            </Empty>
          </div>
        ) : (
          <Row gutter={[16, 16]} className="mt-6">
            {filteredPluginGroups.map((group: PluginGroupItem) => (
            <Col xs={24} lg={12} xl={8} key={group.groupName}>
              <Card
                className="hover:shadow-lg transition-all duration-300 h-full flex flex-col"
                styles={{ body: { flex: 1, display: 'flex', flexDirection: 'column' } }}
                title={
                  <div className="flex items-center gap-2">
                    <SettingOutlined className="text-slate-500" />
                    <span className="font-medium">{group.groupName}</span>
                    <Badge
                      count={group.plugins?.length || 0}
                      style={{ backgroundColor: token.colorTextSecondary }}
                      overflowCount={99}
                    />
                  </div>
                }
                extra={
                  <Space size="small">
                    <Tooltip title={t('common.edit')}>
                      <Button
                        type="text"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(group.groupName)}
                      />
                    </Tooltip>
                    <Popconfirm
                      title={t('common.confirmDelete')}
                      onConfirm={() => handleDelete(group.groupName)}
                      okText={t('common.confirm')}
                      cancelText={t('common.cancel')}
                    >
                      <Tooltip title={t('common.delete')}>
                        <Button type="text" danger icon={<DeleteOutlined />} />
                      </Tooltip>
                    </Popconfirm>
                  </Space>
                }
              >
                {group.plugins && group.plugins.length > 0 ? (
                  <Collapse
                    ghost
                    items={group.plugins.map((plugin, idx) => ({
                      key: idx.toString(),
                      label: (
                        <div className="flex items-center justify-between w-full pr-4">
                          <div className="flex items-center gap-2">
                            <ThunderboltOutlined className="text-slate-500" />
                            <span className="font-medium">{plugin.name}</span>
                          </div>
                          <div className="flex items-center gap-1">
                            <Tag color="blue" className="m-0">v{plugin.version}</Tag>
                            <Tag color="gold" className="m-0">P:{plugin.priority}</Tag>
                          </div>
                        </div>
                      ),
                      children: (
                        <div className="space-y-2">
                          <Descriptions size="small" column={1}>
                            <Descriptions.Item label={t('plugin.version')}>{plugin.version}</Descriptions.Item>
                            <Descriptions.Item label={t('plugin.priority')}>{plugin.priority}</Descriptions.Item>
                            {plugin.externalLookupName && (
                              <Descriptions.Item label={t('plugin.external')}>
                                {plugin.externalLookupName}
                              </Descriptions.Item>
                            )}
                          </Descriptions>
                          {plugin.config && Object.keys(plugin.config).length > 0 && (
                            <div>
                              <div className="text-xs mb-1 flex items-center gap-1" style={{ color: token.colorTextSecondary }}>
                                <CodeOutlined /> {t('plugin.configuration')}
                              </div>
                              <pre
                                className="p-3 rounded text-xs overflow-x-auto"
                                style={{ backgroundColor: token.colorFillSecondary, maxHeight: 200 }}
                              >
                                {JSON.stringify(plugin.config, null, 2)}
                              </pre>
                            </div>
                          )}
                        </div>
                      ),
                    }))}
                  />
                ) : (
                  <div className="text-center py-4">
                    <Empty
                      image={Empty.PRESENTED_IMAGE_SIMPLE}
                      description={
                        <span className="text-slate-400 text-sm">
                          {t('plugin.noPluginsInGroup')}
                        </span>
                      }
                    />
                  </div>
                )}
              </Card>
            </Col>
          ))}
        </Row>
        )}
      </Card>

      {/* Plugin Wizard - for both create and edit */}
      <PluginWizard
        open={isWizardOpen}
        onClose={handleWizardClose}
        onFinish={handleWizardFinish}
        loading={createMutation.isPending || updateMutation.isPending}
        initialDataLoading={editingGroup !== null && isEditYamlLoading}
        mode={editingGroup ? 'edit' : 'create'}
        initialData={editingGroup ? editPluginGroupYaml || undefined : undefined}
      />
    </div>
  );
}
