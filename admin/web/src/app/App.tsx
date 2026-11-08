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

import React, { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Activity,
  BarChart3,
  Boxes,
  CircleHelp,
  Gauge,
  GitBranch,
  LayoutDashboard,
  LifeBuoy,
  Network,
  PanelLeftClose,
  PanelLeftOpen,
  Settings2,
  ShieldCheck,
  SlidersHorizontal,
} from 'lucide-react'
import '../styles.css'
import { Providers } from './providers'
import { clusterApi } from '../services/cluster-api'
import { listenerApi } from '../services/listener-api'
import { profileApi } from '../services/profile-api'
import { routeBindingApi } from '../services/route-binding-api'
import {
  BaseInfoPage,
  ClusterConfigPage,
  ListenerConfigPage,
  PluginGroupConfigPage,
  OpaPage,
  RateLimitPage,
  ResourcePage,
} from '../features/config/ConfigPages'

const nav = [
  ['Overview', LayoutDashboard],
  ['API routes', GitBranch],
  ['Clusters', Boxes],
  ['Listeners', Network],
  ['Plugin groups', SlidersHorizontal],
  ['Rate limiting', Gauge],
  ['OPA policies', ShieldCheck],
] as const
const zhMap: Record<string, string> = {
  Overview: '总览',
  'API routes': 'API 路由',
  Clusters: '集群',
  Listeners: '监听器',
  'Plugin groups': '插件组',
  'Rate limiting': '限流',
  'OPA policies': 'OPA 策略',
  Settings: '设置',
  'Docs & support': '文档与支持',
  'All systems operational': '系统运行正常',
  'New API route': '新建 API 路由',
}
type Section = (typeof nav)[number][0] | 'Settings'

export function App() {
  return (
    <Providers>
      <AppContent />
    </Providers>
  )
}

function AppContent() {
  const [section, setSection] = useState<Section>('Overview')
  const [dark, setDark] = useState(true)
  const [language, setLanguage] = useState<'zh' | 'en'>('zh')
  const [sidebarCollapsed, setSidebarCollapsed] = useState(
    () => localStorage.getItem('pixiu_sidebar_collapsed') === '1',
  )
  const [sidebarMotion, setSidebarMotion] = useState<'idle' | 'collapsing' | 'expanding'>('idle')
  const [counts, setCounts] = useState({ resources: 0, clusters: 0, listeners: 0 })
  const updateResourceCount = useCallback(
    (count: number) =>
      setCounts((value) => (value.resources === count ? value : { ...value, resources: count })),
    [],
  )
  const locale = language === 'zh' ? 'zh-CN' : 'en-US'
  const tx = (value: string) => (language === 'zh' ? zhMap[value] || value : value)
  useEffect(() => {
    Promise.allSettled([routeBindingApi.list('draft'), clusterApi.list(), listenerApi.list()]).then(
      ([resources, clusters, listeners]) =>
        setCounts({
          resources:
            resources.status === 'fulfilled' && Array.isArray(resources.value)
              ? resources.value.length
              : 0,
          clusters:
            clusters.status === 'fulfilled' && Array.isArray(clusters.value)
              ? clusters.value.length
              : 0,
          listeners:
            listeners.status === 'fulfilled' && Array.isArray(listeners.value)
              ? listeners.value.length
              : 0,
        }),
    )
  }, [section])
  useEffect(() => {
    localStorage.setItem('pixiu_sidebar_collapsed', sidebarCollapsed ? '1' : '0')
  }, [sidebarCollapsed])
  const logout = async () => {
    try {
      await profileApi.logout()
    } catch {
      /* Logout is best effort. */
    }
    localStorage.removeItem('pixiu_session')
    window.location.assign('/login')
  }
  const sidebarToggleLabel = sidebarCollapsed
    ? language === 'zh'
      ? '展开侧边栏'
      : 'Expand sidebar'
    : language === 'zh'
      ? '折叠侧边栏'
      : 'Collapse sidebar'
  const toggleSidebar = () => {
    const next = !sidebarCollapsed
    setSidebarCollapsed(next)
    setSidebarMotion(next ? 'collapsing' : 'expanding')
  }
  const content = useMemo(
    () =>
      section === 'Overview' ? (
        <Overview
          language={language}
          tx={tx}
          counts={counts}
          onRoutes={() => setSection('API routes')}
        />
      ) : section === 'API routes' ? (
        <ResourcePage locale={locale} onCountChange={updateResourceCount} />
      ) : section === 'Clusters' ? (
        <ClusterConfigPage locale={locale} />
      ) : section === 'Listeners' ? (
        <ListenerConfigPage locale={locale} />
      ) : section === 'Plugin groups' ? (
        <PluginGroupConfigPage locale={locale} />
      ) : section === 'Rate limiting' ? (
        <RateLimitPage locale={locale} />
      ) : section === 'Settings' ? (
        <BaseInfoPage locale={locale} />
      ) : (
        <OpaPage locale={locale} />
      ),
    [section, language, locale, counts, updateResourceCount],
  )
  return (
    <div
      className={`${dark ? 'app dark' : 'app'}${sidebarCollapsed ? ' sidebar-collapsed' : ''}${sidebarMotion !== 'idle' ? ` sidebar-${sidebarMotion}` : ''}`}
    >
      <aside className={`sidebar${sidebarCollapsed ? ' collapsed' : ''}`}>
        <div className="brand">
          <div className="brand-mark">P</div>
          <div>
            <strong>PIXIU</strong>
            <span>Gateway Admin · v2</span>
          </div>
        </div>
        <nav>
          {nav.map(([label, Icon]) => (
            <button
              className={section === label ? 'nav-item active' : 'nav-item'}
              onClick={() => setSection(label)}
              key={label}
              title={sidebarCollapsed ? tx(label) : undefined}
              aria-label={tx(label)}
            >
              <Icon size={17} />
              <span>{tx(label)}</span>
              {label === 'API routes' && <em>{counts.resources}</em>}
            </button>
          ))}
        </nav>
        <div className="sidebar-bottom">
          <button
            className={section === 'Settings' ? 'nav-item active' : 'nav-item'}
            onClick={() => setSection('Settings')}
            title={sidebarCollapsed ? tx('Settings') : undefined}
            aria-label={tx('Settings')}
          >
            <Settings2 size={17} />
            <span>{tx('Settings')}</span>
          </button>
          <button
            className="nav-item"
            title={sidebarCollapsed ? tx('Docs & support') : undefined}
            aria-label={tx('Docs & support')}
          >
            <CircleHelp size={17} />
            <span>{tx('Docs & support')}</span>
          </button>
          <button
            className="user user-action"
            onClick={() => void logout()}
            title={
              sidebarCollapsed
                ? language === 'zh'
                  ? '点击退出登录'
                  : 'Click to log out'
                : undefined
            }
            aria-label={language === 'zh' ? '点击退出登录' : 'Click to log out'}
          >
            <div className="avatar">A</div>
            <div>
              <b>admin</b>
              <span>{language === 'zh' ? '点击退出登录' : 'Click to log out'}</span>
            </div>
          </button>
        </div>
      </aside>
      <main
        onAnimationEnd={(event) => {
          if (event.animationName.startsWith('sidebar-main-')) setSidebarMotion('idle')
        }}
      >
        <header className="topbar">
          <div className="topbar-leading">
            <button
              className="icon-btn sidebar-toggle-btn"
              onClick={toggleSidebar}
              aria-label={sidebarToggleLabel}
              title={sidebarToggleLabel}
            >
              {sidebarCollapsed ? <PanelLeftOpen size={17} /> : <PanelLeftClose size={17} />}
            </button>
            <div className="crumb">
              <span>Pixiu</span>
              <span>/</span>
              <b>{section === 'Settings' ? tx('Settings') : tx(section)}</b>
            </div>
          </div>
          <div className="top-actions">
            <div className="status-live">
              <i />
              {tx('All systems operational')}
            </div>
            <button
              className="lang-btn"
              onClick={() => setLanguage(language === 'zh' ? 'en' : 'zh')}
            >
              {language === 'zh' ? 'EN' : '中'}
            </button>
            <button className="icon-btn" onClick={() => setDark(!dark)}>
              {dark ? '☼' : '◐'}
            </button>
            <button className="icon-btn" title="Support">
              <LifeBuoy size={17} />
            </button>
            <div className="avatar small">A</div>
          </div>
        </header>
        <div className="content">{content}</div>
      </main>
    </div>
  )
}

function Overview({
  language,
  tx,
  counts,
  onRoutes,
}: {
  language: 'zh' | 'en'
  tx: (v: string) => string
  counts: { resources: number; clusters: number; listeners: number }
  onRoutes: () => void
}) {
  const metrics = [
    ['API routes', String(counts.resources), GitBranch],
    ['Clusters', String(counts.clusters), Boxes],
    ['Listeners', String(counts.listeners), Network],
    ['P95 latency', '—', Activity],
  ] as const
  const recentChanges =
    language === 'zh'
      ? ([
          ['新增路由', 'order.query', '/api/orders/:id', '8 分钟前', 'created'],
          ['更新集群', 'order-service', '连接配置已更新', '22 分钟前', 'updated'],
          ['暂停路由', 'catalog.search', '/api/catalog/search', '昨天', 'paused'],
        ] as const)
      : ([
          ['New route', 'order.query', '/api/orders/:id', '8 min ago', 'created'],
          [
            'Updated cluster',
            'order-service',
            'Connection config updated',
            '22 min ago',
            'updated',
          ],
          ['Paused route', 'catalog.search', '/api/catalog/search', 'Yesterday', 'paused'],
        ] as const)
  return (
    <>
      <div className="page-head">
        <div>
          <p className="eyebrow">{language === 'zh' ? '运行概览' : 'OPERATIONS OVERVIEW'}</p>
          <h1>{language === 'zh' ? '网关运行概览' : 'Gateway overview'}</h1>
          <p className="muted">
            {language === 'zh'
              ? '数据来自当前 Pixiu Admin backend。'
              : 'Live counts from the Pixiu Admin backend.'}
          </p>
        </div>
        <button className="primary" onClick={onRoutes}>
          <GitBranch size={16} /> {tx('API routes')}
        </button>
      </div>
      <div className="metric-grid">
        {metrics.map(([label, value, Icon]) => (
          <div className="metric" key={label}>
            <div className="metric-top">
              <span>{tx(label)}</span>
              <Icon size={17} />
            </div>
            <strong>{value}</strong>
            <small>
              {label === 'P95 latency'
                ? language === 'zh'
                  ? '后端未提供指标接口'
                  : 'No metrics endpoint'
                : language === 'zh'
                  ? '实时配置数量'
                  : 'Live config count'}
            </small>
          </div>
        ))}
      </div>
      <div className="grid-main">
        <section className="panel traffic">
          <div className="panel-head">
            <div>
              <h2>{language === 'zh' ? '流量趋势' : 'Traffic volume'}</h2>
              <span>
                {language === 'zh'
                  ? '当前后端未提供流量指标接口'
                  : 'The backend does not expose traffic metrics yet'}
              </span>
            </div>
            <span className="badge paused">
              <i />
              {language === 'zh' ? '待接入' : 'Not available'}
            </span>
          </div>
          <div className="unavailable">
            <Activity size={28} />
            <b>{language === 'zh' ? '暂无流量指标' : 'No traffic metrics'}</b>
            <span>
              {language === 'zh'
                ? '配置管理接口已接通；指标接口加入后这里会自动接入。'
                : 'Configuration APIs are connected. This panel is ready for a metrics endpoint.'}
            </span>
          </div>
        </section>
        <section className="panel rank">
          <div className="panel-head">
            <div>
              <h2>{language === 'zh' ? '请求排行' : 'Request ranking'}</h2>
              <span>
                {language === 'zh' ? '按后台请求次数排序' : 'Ranked by backend request count'}
              </span>
            </div>
            <BarChart3 size={17} />
          </div>
          <div className="unavailable request-rank-empty">
            <BarChart3 size={28} />
            <b>{language === 'zh' ? '待接入' : 'Pending integration'}</b>
            <span>
              {language === 'zh'
                ? '后续这里会列出后台请求次数排名'
                : 'Backend request counts will be ranked here once metrics are connected.'}
            </span>
          </div>
        </section>
      </div>
      <section className="panel recent recent-changes">
        <div className="panel-head">
          <div>
            <h2>{language === 'zh' ? '最近配置变更' : 'Recent configuration changes'}</h2>
            <span>
              {language === 'zh'
                ? '网关资源的最近操作记录'
                : 'Latest activity across gateway resources'}
            </span>
          </div>
          <span className="badge draft">
            <i />
            {language === 'zh' ? '演示数据' : 'Mock data'}
          </span>
        </div>
        <div className="change-list">
          {recentChanges.map(([kind, title, detail, time, tone]) => (
            <div className="change-row" key={`${kind}-${title}`}>
              <i className={`change-marker ${tone}`} />
              <div className="change-copy">
                <strong>{title}</strong>
                <small>{detail}</small>
              </div>
              <span className="change-kind">{kind}</span>
              <span className="change-time">{time}</span>
            </div>
          ))}
        </div>
      </section>
    </>
  )
}
