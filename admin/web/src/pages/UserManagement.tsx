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
  Table,
  Card,
  Space,
  Tag,
  Button,
  Drawer,
  Form,
  Input,
  Select,
  Switch,
  Popconfirm,
  Tooltip,
  Empty,
} from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  KeyOutlined,
  UserOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { ColumnsType } from 'antd/es/table';
import {
  useUserList,
  useCreateUser,
  useUpdateUser,
  useDeleteUser,
  useResetUserPassword,
  useRoleList,
} from '../api';
import type { UserListItem } from '../api';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

type DrawerMode = 'create' | 'edit' | 'resetPassword' | null;

export function UserManagement() {
  const { t } = useTranslation();
  useDocumentTitle('menu.userManagement');

  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // Drawer states
  const [drawerMode, setDrawerMode] = useState<DrawerMode>(null);
  const [selectedUser, setSelectedUser] = useState<UserListItem | null>(null);

  const [form] = Form.useForm();

  // API hooks
  const { data: userData, isLoading } = useUserList(page, pageSize);
  const { data: roles = [] } = useRoleList();
  const createUser = useCreateUser();
  const updateUser = useUpdateUser();
  const deleteUser = useDeleteUser();
  const resetPassword = useResetUserPassword();

  const users = userData?.items || [];
  const total = userData?.total || 0;

  // Role display helpers
  const getRoleName = (roleId: number) => {
    const role = roles.find((r) => r.id === roleId);
    return role?.role_name || `Role ${roleId}`;
  };

  const getRoleColor = (roleId: number) => {
    switch (roleId) {
      case 1:
        return 'red'; // admin
      case 2:
        return 'blue'; // user
      default:
        return 'default';
    }
  };

  // Handlers
  const handleCreate = () => {
    setSelectedUser(null);
    form.resetFields();
    setDrawerMode('create');
  };

  const handleEdit = (user: UserListItem) => {
    setSelectedUser(user);
    form.setFieldsValue({
      role: user.role,
      enabled: user.enabled,
    });
    setDrawerMode('edit');
  };

  const handleResetPassword = (user: UserListItem) => {
    setSelectedUser(user);
    form.resetFields();
    setDrawerMode('resetPassword');
  };

  const handleCloseDrawer = () => {
    setDrawerMode(null);
    setSelectedUser(null);
    form.resetFields();
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      if (drawerMode === 'create') {
        await createUser.mutateAsync({
          username: values.username,
          password: values.password,
          roleId: values.roleId,
        });
      } else if (drawerMode === 'edit' && selectedUser) {
        await updateUser.mutateAsync({
          id: selectedUser.id,
          data: {
            role: values.role,
            enabled: values.enabled,
          },
        });
      } else if (drawerMode === 'resetPassword' && selectedUser) {
        await resetPassword.mutateAsync({
          id: selectedUser.id,
          newPassword: values.newPassword,
        });
      }

      handleCloseDrawer();
    } catch {
      // Form validation failed
    }
  };

  const handleDelete = async (user: UserListItem) => {
    await deleteUser.mutateAsync(user.id);
  };

  const getDrawerTitle = () => {
    switch (drawerMode) {
      case 'create':
        return t('userManagement.createUser');
      case 'edit':
        return t('userManagement.editUser');
      case 'resetPassword':
        return t('userManagement.resetPassword');
      default:
        return '';
    }
  };

  const isSubmitting =
    createUser.isPending || updateUser.isPending || resetPassword.isPending;

  const columns: ColumnsType<UserListItem> = [
    {
      title: t('userManagement.id'),
      dataIndex: 'id',
      key: 'id',
      width: 80,
      align: 'center',
    },
    {
      title: t('userManagement.username'),
      dataIndex: 'username',
      key: 'username',
      align: 'center',
      render: (username) => (
        <Space>
          <UserOutlined />
          <span>{username}</span>
        </Space>
      ),
    },
    {
      title: t('userManagement.role'),
      dataIndex: 'role',
      key: 'role',
      align: 'center',
      render: (role) => (
        <Tag color={getRoleColor(role)}>{getRoleName(role)}</Tag>
      ),
    },
    {
      title: t('userManagement.enabled'),
      dataIndex: 'enabled',
      key: 'enabled',
      align: 'center',
      render: (enabled) => (
        <Switch checked={enabled} disabled size="small" />
      ),
    },
    {
      title: t('userManagement.createdAt'),
      dataIndex: 'dateCreated',
      key: 'dateCreated',
      align: 'center',
      render: (date) => new Date(date).toLocaleString(),
    },
    {
      title: t('userManagement.updatedAt'),
      dataIndex: 'dateUpdated',
      key: 'dateUpdated',
      align: 'center',
      render: (date) => new Date(date).toLocaleString(),
    },
    {
      title: t('common.operation'),
      key: 'action',
      align: 'center',
      width: 180,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title={t('common.edit')}>
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Tooltip title={t('userManagement.resetPassword')}>
            <Button
              type="text"
              size="small"
              icon={<KeyOutlined />}
              onClick={() => handleResetPassword(record)}
            />
          </Tooltip>
          {record.id !== 1 && (
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
        title={<span className="text-lg font-semibold">{t('userManagement.title')}</span>}
        className="shadow-sm"
        variant="borderless"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            {t('userManagement.createUser')}
          </Button>
        }
      >
        <Table
          className="mt-6"
          columns={columns}
          dataSource={users}
          loading={isLoading}
          rowKey="id"
          size="middle"
          locale={{
            emptyText: (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description={t('userManagement.emptyDescription')}
              />
            ),
          }}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `${t('common.total')} ${total} ${t('common.items')}`,
            onChange: (newPage, newPageSize) => {
              setPage(newPage);
              setPageSize(newPageSize);
            },
          }}
        />
      </Card>

      {/* Drawer for Create/Edit/Reset Password */}
      <Drawer
        title={getDrawerTitle()}
        open={drawerMode !== null}
        onClose={handleCloseDrawer}
        width={480}
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
        <Form form={form} layout="vertical">
          {drawerMode === 'create' && (
            <>
              <Form.Item
                name="username"
                label={t('userManagement.username')}
                rules={[
                  { required: true, message: t('common.required') },
                  { min: 3, message: t('userManagement.usernameMinLength') },
                ]}
              >
                <Input
                  prefix={<UserOutlined />}
                  placeholder={t('userManagement.usernamePlaceholder')}
                />
              </Form.Item>
              <Form.Item
                name="password"
                label={t('userManagement.password')}
                rules={[
                  { required: true, message: t('common.required') },
                  { min: 6, message: t('userManagement.passwordMinLength') },
                ]}
              >
                <Input.Password placeholder={t('userManagement.passwordPlaceholder')} />
              </Form.Item>
              <Form.Item
                name="confirmPassword"
                label={t('userManagement.confirmPassword')}
                dependencies={['password']}
                rules={[
                  { required: true, message: t('common.required') },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('password') === value) {
                        return Promise.resolve();
                      }
                      return Promise.reject(new Error(t('userManagement.passwordMismatch')));
                    },
                  }),
                ]}
              >
                <Input.Password placeholder={t('userManagement.confirmPasswordPlaceholder')} />
              </Form.Item>
              <Form.Item
                name="roleId"
                label={t('userManagement.role')}
                initialValue={2}
              >
                <Select
                  options={roles.map((r) => ({
                    label: r.role_name,
                    value: r.id,
                  }))}
                />
              </Form.Item>
            </>
          )}

          {drawerMode === 'edit' && (
            <>
              <Form.Item label={t('userManagement.username')}>
                <Input value={selectedUser?.username} disabled />
              </Form.Item>
              <Form.Item
                name="role"
                label={t('userManagement.role')}
                rules={[{ required: true }]}
              >
                <Select
                  options={roles.map((r) => ({
                    label: r.role_name,
                    value: r.id,
                  }))}
                />
              </Form.Item>
              <Form.Item
                name="enabled"
                label={t('userManagement.enabled')}
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </>
          )}

          {drawerMode === 'resetPassword' && (
            <>
              <Form.Item label={t('userManagement.username')}>
                <Input value={selectedUser?.username} disabled />
              </Form.Item>
              <Form.Item
                name="newPassword"
                label={t('userManagement.newPassword')}
                rules={[
                  { required: true, message: t('common.required') },
                  { min: 6, message: t('userManagement.passwordMinLength') },
                ]}
              >
                <Input.Password placeholder={t('userManagement.newPasswordPlaceholder')} />
              </Form.Item>
              <Form.Item
                name="confirmPassword"
                label={t('userManagement.confirmPassword')}
                dependencies={['newPassword']}
                rules={[
                  { required: true, message: t('common.required') },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('newPassword') === value) {
                        return Promise.resolve();
                      }
                      return Promise.reject(new Error(t('userManagement.passwordMismatch')));
                    },
                  }),
                ]}
              >
                <Input.Password placeholder={t('userManagement.confirmPasswordPlaceholder')} />
              </Form.Item>
            </>
          )}
        </Form>
      </Drawer>
    </div>
  );
}
