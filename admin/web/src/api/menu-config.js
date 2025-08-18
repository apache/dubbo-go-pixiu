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
export const menuList = [{
  name: '网关配置',
  id: 'Gateway',
  children: [{
    name: '概览',
    id: 'Overview',
    componentName: '/Overview'
  },
  {
    name: '插件配置',
    id: 'Plug',
    componentName: 'Plug'
  },{
    name: '集群管理',
    id: 'Cluster',
    componentName: 'Cluster'
  },{
    name: 'Listener管理',
    id: 'Listener',
    componentName: 'Listener'
  }]
}, {
  name: '限流配置',
  id: 'Flow',
  children: [{
    name: '限流配置',
    id: 'RateLimiter',
    componentName: '/RateLimiter'
  }]
}]
