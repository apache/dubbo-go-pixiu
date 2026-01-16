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

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import i18n from '../locales';
import {
  getClusterList,
  getClusterDetail,
  getClusterYaml,
  createCluster,
  updateCluster,
  deleteCluster,
} from './cluster';
import {
  getListenerList,
  getListenerDetail,
  getListenerYaml,
  createListener,
  updateListener,
  deleteListener,
} from './listener';
import {
  getResourceList,
  getResourceDetail,
  getResourceYaml,
  createResource,
  updateResource,
  deleteResource,
  getMethodList,
  getMethodDetail,
  createMethod,
  updateMethod,
  deleteMethod,
} from './resource';
import { API_CODE } from './types';
import type { Resource } from './types';

// ==================== Utils ====================

/**
 * Safely parse JSON array response from backend.
 * Handles cases where backend returns JSON string instead of object.
 */
function parseJsonArray<T>(data: unknown): T[] {
  try {
    const parsed = typeof data === 'string' ? JSON.parse(data) : data;
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

// ==================== Cluster Hooks ====================

export const useClusterList = () => {
  return useQuery({
    queryKey: ['clusters'],
    queryFn: async () => {
      const response = await getClusterList();
      if (response.code === API_CODE.SUCCESS) {
        return response.data || [];
      }
      throw new Error(response.msg || 'Failed to fetch clusters');
    },
  });
};

export const useClusterDetail = (clusterId: string | null) => {
  return useQuery({
    queryKey: ['cluster', clusterId],
    queryFn: async () => {
      if (!clusterId) return null;
      const response = await getClusterDetail(clusterId);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch cluster detail');
    },
    enabled: !!clusterId,
  });
};

export const useClusterYaml = (clusterId: string | null, mode: 'edit' | 'view') => {
  return useQuery({
    queryKey: ['clusterYaml', clusterId, mode],
    queryFn: async () => {
      if (!clusterId) return null;
      const response = await getClusterYaml(clusterId);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch cluster YAML');
    },
    enabled: !!clusterId,
  });
};

export const useCreateCluster = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createCluster,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['clusters'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useUpdateCluster = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { name: string; content: string }) =>
      updateCluster(params.name, params.content),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['clusters'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeleteCluster = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteCluster,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['clusters'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== Listener Hooks ====================

export const useListenerList = () => {
  return useQuery({
    queryKey: ['listeners'],
    queryFn: async () => {
      const response = await getListenerList();
      if (response.code === API_CODE.SUCCESS) {
        return response.data || [];
      }
      throw new Error(response.msg || 'Failed to fetch listeners');
    },
  });
};

export const useListenerDetail = (listenerName: string | null) => {
  return useQuery({
    queryKey: ['listener', listenerName],
    queryFn: async () => {
      if (!listenerName) return null;
      const response = await getListenerDetail(listenerName);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch listener detail');
    },
    enabled: !!listenerName,
  });
};

export const useListenerYaml = (listenerName: string | null) => {
  return useQuery({
    queryKey: ['listenerYaml', listenerName],
    queryFn: async () => {
      if (!listenerName) return null;
      const response = await getListenerYaml(listenerName);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch listener YAML');
    },
    enabled: !!listenerName,
    staleTime: 0, // Always refetch when opened
  });
};

export const useCreateListener = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createListener,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['listeners'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useUpdateListener = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { name: string; content: string }) =>
      updateListener(params.name, params.content),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['listeners'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeleteListener = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteListener,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['listeners'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== Resource Hooks ====================

export const useResourceList = () => {
  return useQuery({
    queryKey: ['resources'],
    queryFn: async () => {
      const response = await getResourceList();
      if (response.code === API_CODE.SUCCESS) {
        return parseJsonArray<Resource>(response.data);
      }
      throw new Error(response.msg || 'Failed to fetch resources');
    },
  });
};

export const useResourceDetail = (resourceId: string | null) => {
  return useQuery({
    queryKey: ['resource', resourceId],
    queryFn: async () => {
      if (!resourceId) return null;
      const response = await getResourceDetail(resourceId);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch resource detail');
    },
    enabled: !!resourceId,
  });
};

export const useResourceYaml = (resourceId: string | null) => {
  return useQuery({
    queryKey: ['resourceYaml', resourceId],
    queryFn: async () => {
      if (!resourceId) return null;
      const response = await getResourceYaml(resourceId);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch resource YAML');
    },
    enabled: !!resourceId,
  });
};

export const useCreateResource = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createResource,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['resources'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useUpdateResource = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: string; content: string }) =>
      updateResource(params.id, params.content),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['resources'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeleteResource = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteResource,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['resources'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== Method Hooks ====================

export const useMethodList = (resourceId: string | null) => {
  return useQuery({
    queryKey: ['methods', resourceId],
    queryFn: async () => {
      if (!resourceId) return [];
      const response = await getMethodList(resourceId);
      if (response.code === API_CODE.SUCCESS) {
        return response.data || [];
      }
      throw new Error(response.msg || 'Failed to fetch methods');
    },
    enabled: !!resourceId,
  });
};

export const useMethodDetail = (methodId: string | null, resourceId: string | null) => {
  return useQuery({
    queryKey: ['method', methodId, resourceId],
    queryFn: async () => {
      if (!methodId || !resourceId) return null;
      const response = await getMethodDetail(methodId, resourceId);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch method detail');
    },
    enabled: !!methodId && !!resourceId,
  });
};

export const useCreateMethod = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { content: string; resourceId: string }) =>
      createMethod(params.content, params.resourceId),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['methods'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useUpdateMethod = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: string; content: string; resourceId: string }) =>
      updateMethod(params.id, params.content, params.resourceId),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['methods'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeleteMethod = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: string; resourceId: string }) =>
      deleteMethod(params.id, params.resourceId),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['methods'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== Plugin Group Hooks ====================

import {
  getPluginGroupList,
  getPluginGroupDetail,
  getPluginGroupYaml,
  createPluginGroup,
  updatePluginGroup,
  deletePluginGroup,
} from './plugin';
import type { PluginGroupItem } from './types';

export const usePluginGroupList = () => {
  return useQuery({
    queryKey: ['pluginGroups'],
    queryFn: async () => {
      const response = await getPluginGroupList();
      if (response.code === API_CODE.SUCCESS) {
        return parseJsonArray<PluginGroupItem>(response.data);
      }
      throw new Error(response.msg || 'Failed to fetch plugin groups');
    },
  });
};

export const usePluginGroupDetail = (name: string | null) => {
  return useQuery({
    queryKey: ['pluginGroup', name],
    queryFn: async () => {
      if (!name) return null;
      const response = await getPluginGroupDetail(name);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch plugin group detail');
    },
    enabled: !!name,
  });
};

export const usePluginGroupYaml = (name: string | null) => {
  return useQuery({
    queryKey: ['pluginGroupYaml', name],
    queryFn: async () => {
      if (!name) return null;
      const response = await getPluginGroupYaml(name);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch plugin group YAML');
    },
    enabled: !!name,
  });
};

export const useCreatePluginGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createPluginGroup,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['pluginGroups'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useUpdatePluginGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { name: string; content: string }) =>
      updatePluginGroup(params.name, params.content),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['pluginGroups'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeletePluginGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deletePluginGroup,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['pluginGroups'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== System Settings Hooks ====================

import {
  getSystemSettings,
  updateSystemSettings,
  getGatewayInfo,
  updateGatewaySettings,
} from './settings';
import type { SystemSettings, GatewayInfo } from './settings';
import type { UserListItem, Role } from './types';
import {
  getUserList,
  getUser,
  updateUser,
  deleteUser,
  resetUserPassword,
  assignUserRole,
  getRoleList,
} from './user';

export const useSystemSettings = () => {
  return useQuery({
    queryKey: ['systemSettings'],
    queryFn: async () => {
      try {
        const response = await getSystemSettings();
        if (response.code === API_CODE.SUCCESS) {
          return response.data as SystemSettings;
        }
        // Return default settings if API not available
        return {} as SystemSettings;
      } catch {
        // Return default settings if API fails
        return {} as SystemSettings;
      }
    },
  });
};

export const useUpdateSystemSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateSystemSettings,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['systemSettings'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useGatewayInfo = () => {
  return useQuery({
    queryKey: ['gatewayInfo'],
    queryFn: async () => {
      try {
        const response = await getGatewayInfo();
        if (response.code === API_CODE.SUCCESS) {
          return response.data as GatewayInfo;
        }
        // Return default info if API not available
        return {
          version: '1.0.0',
          status: 'Running',
          httpPort: '8081',
          xdsPort: '18000',
          startTime: new Date().toLocaleString(),
          uptime: '-',
        } as GatewayInfo;
      } catch {
        // Return default info if API fails
        return {
          version: '1.0.0',
          status: 'Running',
          httpPort: '8081',
          xdsPort: '18000',
          startTime: new Date().toLocaleString(),
          uptime: '-',
        } as GatewayInfo;
      }
    },
  });
};

export const useUpdateGatewaySettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateGatewaySettings,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['gatewayInfo'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== User Management Hooks ====================

export const useUserList = (page = 1, pageSize = 20) => {
  return useQuery({
    queryKey: ['users', page, pageSize],
    queryFn: async () => {
      const response = await getUserList(page, pageSize);
      if (response.code === API_CODE.SUCCESS) {
        return response.data;
      }
      throw new Error(response.msg || 'Failed to fetch users');
    },
  });
};

export const useUser = (id: number | null) => {
  return useQuery({
    queryKey: ['user', id],
    queryFn: async () => {
      if (!id) return null;
      const response = await getUser(id);
      if (response.code === API_CODE.SUCCESS) {
        return response.data as UserListItem;
      }
      throw new Error(response.msg || 'Failed to fetch user');
    },
    enabled: !!id,
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: number; data: { role?: number; enabled?: boolean } }) =>
      updateUser(params.id, params.data),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['users'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeleteUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteUser,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['users'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useResetUserPassword = () => {
  return useMutation({
    mutationFn: (params: { id: number; newPassword: string }) =>
      resetUserPassword(params.id, params.newPassword),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.success'));
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useAssignUserRole = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { userId: number; roleId: number }) =>
      assignUserRole(params.userId, params.roleId),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.success'));
        queryClient.invalidateQueries({ queryKey: ['users'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useRoleList = () => {
  return useQuery({
    queryKey: ['roles'],
    queryFn: async () => {
      const response = await getRoleList();
      if (response.code === API_CODE.SUCCESS) {
        return (response.data || []) as Role[];
      }
      throw new Error(response.msg || 'Failed to fetch roles');
    },
  });
};

import {
  createRole,
  updateRole,
  deleteRole,
  getRolePermissions,
  updateRolePermissions,
  getPermissionList,
  createUser,
} from './user';
import type { Permission } from './types';

export const useCreateUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createUser,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['users'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useCreateRole = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createRole,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.createSuccess'));
        queryClient.invalidateQueries({ queryKey: ['roles'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useUpdateRole = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: number; data: { roleName: string; description?: string } }) =>
      updateRole(params.id, params.data),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['roles'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const useDeleteRole = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteRole,
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.deleteSuccess'));
        queryClient.invalidateQueries({ queryKey: ['roles'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

export const usePermissionList = () => {
  return useQuery({
    queryKey: ['permissions'],
    queryFn: async () => {
      const response = await getPermissionList();
      if (response.code === API_CODE.SUCCESS) {
        return (response.data || []) as Permission[];
      }
      throw new Error(response.msg || 'Failed to fetch permissions');
    },
  });
};

export const useRolePermissions = (roleId: number | null) => {
  return useQuery({
    queryKey: ['rolePermissions', roleId],
    queryFn: async () => {
      if (!roleId) return [];
      const response = await getRolePermissions(roleId);
      if (response.code === API_CODE.SUCCESS) {
        return (response.data || []) as Permission[];
      }
      throw new Error(response.msg || 'Failed to fetch role permissions');
    },
    enabled: !!roleId,
  });
};

export const useUpdateRolePermissions = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { roleId: number; permissionIds: number[] }) =>
      updateRolePermissions(params.roleId, params.permissionIds),
    onSuccess: (response) => {
      if (response.code === API_CODE.SUCCESS) {
        message.success(i18n.t('common.updateSuccess'));
        queryClient.invalidateQueries({ queryKey: ['rolePermissions'] });
      } else {
        message.error(response.msg || i18n.t('common.operationFailed'));
      }
    },
    onError: () => {
      // Error message shown by request interceptor
    },
  });
};

// ==================== Instance Hooks ====================

import { getInstances, getInstanceStats } from './instance';
import type { Instance, InstanceStats } from './types';

export const useInstances = () => {
  return useQuery({
    queryKey: ['instances'],
    queryFn: async () => {
      try {
        const response = await getInstances();
        if (response.code === API_CODE.SUCCESS) {
          return (response.data || []) as Instance[];
        }
        return [] as Instance[];
      } catch {
        // Return empty array if API fails (xDS server may not be running)
        return [] as Instance[];
      }
    },
  });
};

export const useInstanceStats = () => {
  return useQuery({
    queryKey: ['instanceStats'],
    queryFn: async () => {
      try {
        const response = await getInstanceStats();
        if (response.code === API_CODE.SUCCESS) {
          return response.data as InstanceStats;
        }
        return { total: 0, connected: 0, disconnected: 0 } as InstanceStats;
      } catch {
        // Return default stats if API fails
        return { total: 0, connected: 0, disconnected: 0 } as InstanceStats;
      }
    },
  });
};
