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

import { useState, useMemo, useCallback } from 'react';
import {
  Card,
  Table,
  Tag,
  Empty,
  Button,
  Space,
  Drawer,
  Form,
  Input,
  Popconfirm,
  Tooltip,
  Checkbox,
  Spin,
} from 'antd';
import {
  SafetyCertificateOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { ColumnsType } from 'antd/es/table';
import type { CheckboxChangeEvent } from 'antd/es/checkbox';
import {
  useRoleList,
  useCreateRole,
  useUpdateRole,
  useDeleteRole,
  usePermissionList,
  useRolePermissions,
  useUpdateRolePermissions,
} from '../api';
import type { Role, Permission } from '../api';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

type DrawerMode = 'create' | 'edit' | 'permissions' | null;

// Group permissions by resource
function groupPermissions(permissions: Permission[]): Record<string, Permission[]> {
  const grouped: Record<string, Permission[]> = {};
  permissions.forEach((perm) => {
    const resource = perm.resource;
    if (!grouped[resource]) {
      grouped[resource] = [];
    }
    grouped[resource].push(perm);
  });
  return grouped;
}

export function RoleManagement() {
  const { t } = useTranslation();
  useDocumentTitle('menu.roleManagement');

  // Drawer states
  const [drawerMode, setDrawerMode] = useState<DrawerMode>(null);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<number[]>([]);

  const [form] = Form.useForm();

  // API hooks
  const { data: roles = [], isLoading } = useRoleList();
  const { data: allPermissions = [], isLoading: permissionsLoading } = usePermissionList();
  const { data: rolePermissions = [], isLoading: rolePermissionsLoading } = useRolePermissions(
    selectedRole?.id ?? null
  );

  const createRole = useCreateRole();
  const updateRole = useUpdateRole();
  const deleteRole = useDeleteRole();
  const updateRolePermissions = useUpdateRolePermissions();

  // Compute effective permission IDs: use user selection if any, otherwise use loaded role permissions
  const effectivePermissionIds = useMemo(() => {
    if (drawerMode === 'permissions') {
      // If user has made selections, use them
      if (selectedPermissionIds.length > 0) {
        return selectedPermissionIds;
      }
      // Otherwise use loaded role permissions as initial value
      return rolePermissions.map((p) => p.id);
    }
    return selectedPermissionIds;
  }, [drawerMode, selectedPermissionIds, rolePermissions]);

  // Helper to get current permission IDs for updates
  const getCurrentPermissionIds = useCallback(() => {
    if (selectedPermissionIds.length > 0) {
      return selectedPermissionIds;
    }
    return rolePermissions.map((p) => p.id);
  }, [selectedPermissionIds, rolePermissions]);

  const getRoleColor = (roleName: string) => {
    switch (roleName.toLowerCase()) {
      case 'admin':
        return 'red';
      case 'user':
        return 'blue';
      default:
        return 'default';
    }
  };

  // Handlers
  const handleCreate = () => {
    setSelectedRole(null);
    form.resetFields();
    setDrawerMode('create');
  };

  const handleEdit = (role: Role) => {
    setSelectedRole(role);
    form.setFieldsValue({
      roleName: role.role_name,
      description: role.description,
    });
    setDrawerMode('edit');
  };

  const handleEditPermissions = (role: Role) => {
    setSelectedRole(role);
    setSelectedPermissionIds([]);
    setDrawerMode('permissions');
  };

  const handleCloseDrawer = () => {
    setDrawerMode(null);
    setSelectedRole(null);
    form.resetFields();
    setSelectedPermissionIds([]);
  };

  const handleSubmit = async () => {
    try {
      if (drawerMode === 'create') {
        const values = await form.validateFields();
        await createRole.mutateAsync({
          roleName: values.roleName,
          description: values.description,
        });
      } else if (drawerMode === 'edit' && selectedRole) {
        const values = await form.validateFields();
        await updateRole.mutateAsync({
          id: selectedRole.id,
          data: {
            roleName: values.roleName,
            description: values.description,
          },
        });
      } else if (drawerMode === 'permissions' && selectedRole) {
        await updateRolePermissions.mutateAsync({
          roleId: selectedRole.id,
          permissionIds: effectivePermissionIds,
        });
      }

      handleCloseDrawer();
    } catch {
      // Form validation failed
    }
  };

  const handleDelete = async (role: Role) => {
    await deleteRole.mutateAsync(role.id);
  };

  // Permission checkbox handlers
  const handlePermissionChange = (permId: number, checked: boolean) => {
    const currentIds = getCurrentPermissionIds();
    if (checked) {
      setSelectedPermissionIds([...currentIds, permId]);
    } else {
      setSelectedPermissionIds(currentIds.filter((id) => id !== permId));
    }
  };

  const handleResourceSelectAll = (_resource: string, permissions: Permission[], checked: boolean) => {
    const currentIds = getCurrentPermissionIds();
    const permIds = permissions.map((p) => p.id);
    if (checked) {
      setSelectedPermissionIds([...new Set([...currentIds, ...permIds])]);
    } else {
      setSelectedPermissionIds(currentIds.filter((id) => !permIds.includes(id)));
    }
  };

  const isResourceAllSelected = (permissions: Permission[]) => {
    return permissions.every((p) => effectivePermissionIds.includes(p.id));
  };

  const isResourceIndeterminate = (permissions: Permission[]) => {
    const selectedCount = permissions.filter((p) => effectivePermissionIds.includes(p.id)).length;
    return selectedCount > 0 && selectedCount < permissions.length;
  };

  const getDrawerTitle = () => {
    switch (drawerMode) {
      case 'create':
        return t('roleManagement.createRole');
      case 'edit':
        return t('roleManagement.editRole');
      case 'permissions':
        return `${t('roleManagement.configurePermissions')} - ${selectedRole?.role_name}`;
      default:
        return '';
    }
  };

  const isSubmitting =
    createRole.isPending || updateRole.isPending || updateRolePermissions.isPending;

  const groupedPermissions = groupPermissions(allPermissions);

  const columns: ColumnsType<Role> = [
    {
      title: t('roleManagement.id'),
      dataIndex: 'id',
      key: 'id',
      width: 80,
      align: 'center',
    },
    {
      title: t('roleManagement.roleName'),
      dataIndex: 'role_name',
      key: 'role_name',
      align: 'center',
      render: (name) => (
        <Tag color={getRoleColor(name)} icon={<SafetyCertificateOutlined />}>
          {name}
        </Tag>
      ),
    },
    {
      title: t('roleManagement.description'),
      dataIndex: 'description',
      key: 'description',
      align: 'center',
      render: (desc) => desc || '-',
    },
    {
      title: t('common.operation'),
      key: 'action',
      align: 'center',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title={t('roleManagement.configurePermissions')}>
            <Button
              type="text"
              size="small"
              icon={<SettingOutlined />}
              onClick={() => handleEditPermissions(record)}
            />
          </Tooltip>
          <Tooltip title={t('common.edit')}>
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          {record.id !== 1 && record.id !== 2 && (
            <Popconfirm
              title={t('common.deleteConfirm')}
              onConfirm={() => handleDelete(record)}
              okText={t('common.yes')}
              cancelText={t('common.no')}
            >
              <Tooltip title={t('common.delete')}>
                <Button
                  type="text"
                  size="small"
                  danger
                  icon={<DeleteOutlined />}
                />
              </Tooltip>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card
        title={<span className="text-lg font-semibold">{t('roleManagement.title')}</span>}
        className="shadow-sm"
        variant="borderless"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            {t('roleManagement.createRole')}
          </Button>
        }
      >
        <Table
          className="mt-6"
          columns={columns}
          dataSource={roles}
          loading={isLoading}
          rowKey="id"
          size="middle"
          locale={{
            emptyText: (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description={t('roleManagement.emptyDescription')}
              />
            ),
          }}
          pagination={false}
        />
      </Card>

      {/* Drawer for Create/Edit Role */}
      <Drawer
        title={getDrawerTitle()}
        open={drawerMode !== null}
        onClose={handleCloseDrawer}
        width={drawerMode === 'permissions' ? 640 : 480}
        extra={
          <Space>
            <Button onClick={handleCloseDrawer}>{t('common.cancel')}</Button>
            <Button type="primary" onClick={handleSubmit} loading={isSubmitting}>
              {t('common.confirm')}
            </Button>
          </Space>
        }
        destroyOnClose
      >
        {(drawerMode === 'create' || drawerMode === 'edit') && (
          <Form form={form} layout="vertical">
            <Form.Item
              name="roleName"
              label={t('roleManagement.roleName')}
              rules={[
                { required: true, message: t('common.required') },
                { min: 2, message: t('roleManagement.roleNameMinLength') },
              ]}
            >
              <Input placeholder={t('roleManagement.roleNamePlaceholder')} />
            </Form.Item>
            <Form.Item
              name="description"
              label={t('roleManagement.description')}
            >
              <Input.TextArea
                rows={3}
                placeholder={t('roleManagement.descriptionPlaceholder')}
              />
            </Form.Item>
          </Form>
        )}

        {drawerMode === 'permissions' && (
          <Spin spinning={permissionsLoading || rolePermissionsLoading}>
            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
              {Object.entries(groupedPermissions).map(([resource, permissions]) => (
                <Card
                  key={resource}
                  size="small"
                  title={
                    <Checkbox
                      checked={isResourceAllSelected(permissions)}
                      indeterminate={isResourceIndeterminate(permissions)}
                      onChange={(e: CheckboxChangeEvent) =>
                        handleResourceSelectAll(resource, permissions, e.target.checked)
                      }
                    >
                      <span style={{ fontWeight: 500 }}>
                        {t(`roleManagement.resources.${resource}`, { defaultValue: resource })}
                      </span>
                    </Checkbox>
                  }
                  extra={<Tag color="blue">{permissions.length}</Tag>}
                >
                  <Checkbox.Group
                    value={effectivePermissionIds}
                    style={{ width: '100%' }}
                  >
                    <Space wrap size={[24, 8]}>
                      {permissions.map((perm) => (
                        <Checkbox
                          key={perm.id}
                          value={perm.id}
                          onChange={(e: CheckboxChangeEvent) =>
                            handlePermissionChange(perm.id, e.target.checked)
                          }
                        >
                          {t(`roleManagement.actions.${perm.action}`, { defaultValue: perm.action })}
                        </Checkbox>
                      ))}
                    </Space>
                  </Checkbox.Group>
                </Card>
              ))}
              {Object.keys(groupedPermissions).length === 0 && (
                <Empty description={t('roleManagement.noPermissions')} />
              )}
            </Space>
          </Spin>
        )}
      </Drawer>
    </div>
  );
}
