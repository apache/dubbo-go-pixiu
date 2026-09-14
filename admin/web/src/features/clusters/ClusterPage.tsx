/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { useQuery } from '@tanstack/react-query'
import { clusterApi } from '../../services/cluster-api'
import type { JsonObject } from '../../types/api'
export function ClusterPage() {
  const q = useQuery({ queryKey: ['clusters'], queryFn: clusterApi.list })
  return (
    <div className="page-resource">
      <div className="placeholder-head">
        <div>
          <p className="eyebrow">CONNECTED API</p>
          <h1>集群</h1>
          <p className="muted">管理 Dubbo 集群、注册中心和服务端点。</p>
        </div>
        <button className="primary">+ 新建集群</button>
      </div>
      <div className="panel route-panel">
        {q.isLoading ? (
          <div className="loading-skeleton" />
        ) : q.isError ? (
          <div className="empty">
            {q.error instanceof Error ? q.error.message : '加载失败'}
            <button className="secondary" onClick={() => q.refetch()}>
              重试
            </button>
          </div>
        ) : !q.data?.length ? (
          <div className="empty">
            暂无集群<button className="primary">新建集群</button>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>名称</th>
                <th>类型</th>
                <th>端点</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {q.data.map((x: JsonObject, i: number) => (
                <tr key={String(x.id || i)}>
                  <td>
                    <b>{String(x.name || x.id || '未命名集群')}</b>
                  </td>
                  <td>{String(x.typeStr || x.type || 'Dubbo')}</td>
                  <td>
                    {String(x.address || (Array.isArray(x.endpoints) ? x.endpoints.length : '-'))}
                  </td>
                  <td>
                    <button className="link-btn">编辑</button>
                    <button className="link-btn">删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
