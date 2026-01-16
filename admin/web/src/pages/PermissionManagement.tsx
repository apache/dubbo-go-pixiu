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
import { Card, Table, Tag, Input, Space } from 'antd';
import { SafetyCertificateOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { ColumnsType } from 'antd/es/table';
import { usePermissionList } from '../api';
import type { Permission } from '../api';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

export function PermissionManagement() {
  const { t } = useTranslation();
  useDocumentTitle('menu.permissionManagement');
  const [searchText, setSearchText] = useState('');

  const { data: permissions = [], isLoading } = usePermissionList();

  // Filter permissions by search text
  const filteredPermissions = permissions.filter(
    (perm) =>
      perm.resource?.toLowerCase().includes(searchText.toLowerCase()) ||
      perm.action?.toLowerCase().includes(searchText.toLowerCase())
  );

  // Group permissions by resource for display
  const getActionColor = (action: string) => {
    switch (action.toLowerCase()) {
      case 'create':
        return 'green';
      case 'read':
        return 'blue';
      case 'update':
        return 'orange';
      case 'delete':
        return 'red';
      default:
        return 'default';
    }
  };

  const columns: ColumnsType<Permission> = [
    {
      title: t('permissionManagement.id'),
      dataIndex: 'id',
      key: 'id',
      width: 80,
      align: 'center',
    },
    {
      title: t('permissionManagement.resource'),
      dataIndex: 'resource',
      key: 'resource',
      align: 'center',
      render: (resource) => (
        <Tag icon={<SafetyCertificateOutlined />}>
          {t(`roleManagement.resources.${resource}`, { defaultValue: resource })}
        </Tag>
      ),
      filters: [...new Set(permissions.map((p) => p.resource))].map((r) => ({
        text: t(`roleManagement.resources.${r}`, { defaultValue: r }),
        value: r,
      })),
      onFilter: (value, record) => record.resource === value,
    },
    {
      title: t('permissionManagement.action'),
      dataIndex: 'action',
      key: 'action',
      align: 'center',
      render: (action) => (
        <Tag color={getActionColor(action)}>
          {t(`roleManagement.actions.${action}`, { defaultValue: action })}
        </Tag>
      ),
      filters: [...new Set(permissions.map((p) => p.action))].map((a) => ({
        text: t(`roleManagement.actions.${a}`, { defaultValue: a }),
        value: a,
      })),
      onFilter: (value, record) => record.action === value,
    },
    {
      title: t('permissionManagement.description'),
      key: 'description',
      align: 'center',
      render: (_, record) =>
        t(`permissionManagement.permDesc.${record.resource}.${record.action}`, {
          defaultValue: `${t(`roleManagement.actions.${record.action}`, { defaultValue: record.action })} ${t(`roleManagement.resources.${record.resource}`, { defaultValue: record.resource })}`,
        }),
    },
  ];

  return (
    <div>
      <Card
        title={<span className="text-lg font-semibold">{t('permissionManagement.title')}</span>}
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
          </Space>
        }
      >
        <Table
          className="mt-6"
          columns={columns}
          dataSource={filteredPermissions}
          loading={isLoading}
          rowKey="id"
          size="middle"
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
          }}
        />
      </Card>
    </div>
  );
}
