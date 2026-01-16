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

// Types - export types separately
export { API_CODE } from './types';
export type {
  // API Response
  ApiResponse,
  // Auth Types
  LoginRequest,
  LoginResponse,
  UserInfo,
  ChangePasswordRequest,
  // Cluster Types
  Cluster,
  ClusterDetail,
  // Listener Types
  Listener,
  ListenerDetail,
  ListenerAddress,
  SocketAddress,
  RouteConfig,
  RouteItem,
  RouteMatch,
  RouteTarget,
  HttpFilter,
  // Resource Types (Mapping)
  Resource,
  ResourceDetail,
  Method,
  MethodDetail,
  Filter,
  Param,
  BodyDefinition,
  InboundRequest,
  IntegrationRequest,
  DubboBackendConfig,
  HttpBackendConfig,
  MappingParam,
  // Plugin Types
  PluginGroup,
  Plugin,
  PluginGroupItem,
  PluginConfig,
  PluginGroupDetail,
  // API Config Types
  APIConfig,
  Definition,
  // Helper Types
  BaseConfig,
  // User Management Types
  UserListItem,
  UserListResponse,
  UpdateUserRequest,
  Role,
  ROLE_NAMES,
  Permission,
  // Instance Types
  Instance,
  InstanceStats,
} from './types';

// Request utilities
export {
  http,
  createFormData,
  getToken,
  setToken,
  removeToken,
  getStoredUser,
  setStoredUser,
  removeStoredUser,
} from './request';

// Auth API
export {
  login,
  register,
  logout,
  getUserInfo,
  getUserRole,
  checkIsAdmin,
  changePassword,
} from './auth';

// Cluster API
export {
  getClusterList,
  getClusterDetail,
  getClusterYaml,
  createCluster,
  updateCluster,
  deleteCluster,
} from './cluster';

// Listener API
export {
  getListenerList,
  getListenerDetail,
  getListenerYaml,
  createListener,
  updateListener,
  deleteListener,
} from './listener';

// Resource API
export {
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

// Plugin Group API
export {
  getPluginGroupList,
  getPluginGroupDetail,
  getPluginGroupYaml,
  createPluginGroup,
  updatePluginGroup,
  deletePluginGroup,
} from './plugin';

// System Settings API
export {
  getSystemSettings,
  updateSystemSettings,
  getGatewayInfo,
  updateGatewaySettings,
} from './settings';
export type { SystemSettings, GatewayInfo } from './settings';

// User API
export {
  getUserList,
  getUser,
  createUser,
  updateUser,
  deleteUser,
  resetUserPassword,
  assignUserRole,
  getRoleList,
  getRole,
  createRole,
  updateRole,
  deleteRole,
  getRolePermissions,
  updateRolePermissions,
  getPermissionList,
} from './user';

// Instance API
export { getInstances, getInstanceStats } from './instance';

// React Query Hooks
export {
  useClusterList,
  useClusterDetail,
  useClusterYaml,
  useCreateCluster,
  useUpdateCluster,
  useDeleteCluster,
  useListenerList,
  useListenerDetail,
  useListenerYaml,
  useCreateListener,
  useUpdateListener,
  useDeleteListener,
  useResourceList,
  useResourceDetail,
  useResourceYaml,
  useCreateResource,
  useUpdateResource,
  useDeleteResource,
  useMethodList,
  useMethodDetail,
  useCreateMethod,
  useUpdateMethod,
  useDeleteMethod,
  usePluginGroupList,
  usePluginGroupDetail,
  usePluginGroupYaml,
  useCreatePluginGroup,
  useUpdatePluginGroup,
  useDeletePluginGroup,
  useSystemSettings,
  useUpdateSystemSettings,
  useGatewayInfo,
  useUpdateGatewaySettings,
  // User hooks
  useUserList,
  useUser,
  useCreateUser,
  useUpdateUser,
  useDeleteUser,
  useResetUserPassword,
  useAssignUserRole,
  useRoleList,
  useCreateRole,
  useUpdateRole,
  useDeleteRole,
  usePermissionList,
  useRolePermissions,
  useUpdateRolePermissions,
  // Instance hooks
  useInstances,
  useInstanceStats,
} from './hooks';
