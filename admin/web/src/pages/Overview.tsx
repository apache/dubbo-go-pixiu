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

import { Row, Col, Card, Statistic, Spin, Table, Tag, Button, Space, Tooltip, Progress, Badge } from 'antd';
import {
  CloudServerOutlined,
  ApiOutlined,
  AudioOutlined,
  ClusterOutlined,
  CheckCircleOutlined,
  PlusOutlined,
  SettingOutlined,
  SyncOutlined,
  ThunderboltOutlined,
  SafetyCertificateOutlined,
  GlobalOutlined,
  RocketOutlined,
  HistoryOutlined,
  HeartOutlined,
  CloudOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { SafeECharts } from '../components/SafeECharts';
import { useThemeStore } from '../stores/theme';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import {
  useClusterList,
  useListenerList,
  useResourceList,
  usePluginGroupList,
  useInstanceStats,
} from '../api';
import type { ColumnsType } from 'antd/es/table';
import type { Resource, Cluster, Listener } from '../api';

export function Overview() {
  const { t } = useTranslation();
  const { isDark } = useThemeStore();
  const navigate = useNavigate();
  useDocumentTitle('menu.overview');

  // Fetch real data
  const { data: clusters = [], isLoading: isClustersLoading } = useClusterList();
  const { data: listeners = [], isLoading: isListenersLoading } = useListenerList();
  const { data: resources = [], isLoading: isResourcesLoading } = useResourceList();
  const { data: pluginGroups = [], isLoading: isPluginsLoading } = usePluginGroupList();
  const { data: instanceStats } = useInstanceStats();

  const isLoading = isClustersLoading || isListenersLoading || isResourcesLoading || isPluginsLoading;

  // Resource table columns
  const resourceColumns: ColumnsType<Resource> = [
    {
      title: t('mapping.path'),
      dataIndex: 'path',
      key: 'path',
      align: 'center',
    },
    {
      title: t('mapping.type'),
      dataIndex: 'type',
      key: 'type',
      width: 120,
      align: 'center',
      render: (type) => <Tag color="blue">{type}</Tag>,
    },
    {
      title: t('mapping.description'),
      dataIndex: 'description',
      key: 'description',
      align: 'center',
      ellipsis: true,
    },
  ];

  // Cluster table columns
  const clusterColumns: ColumnsType<Cluster> = [
    {
      title: t('cluster.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('cluster.type'),
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type) => <Tag color="blue">{type || 'N/A'}</Tag>,
    },
    {
      title: t('common.status'),
      key: 'status',
      width: 100,
      render: () => (
        <Tag color="success" icon={<CheckCircleOutlined />}>
          {t('common.running')}
        </Tag>
      ),
    },
  ];

  // Listener table columns
  const listenerColumns: ColumnsType<Listener> = [
    {
      title: t('listener.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('listener.protocol'),
      dataIndex: 'protocol',
      key: 'protocol',
      width: 100,
      render: (protocol) => <Tag color="cyan">{protocol || 'HTTP'}</Tag>,
    },
    {
      title: t('common.status'),
      key: 'status',
      width: 100,
      render: () => (
        <Tag color="success" icon={<CheckCircleOutlined />}>
          {t('common.running')}
        </Tag>
      ),
    },
  ];

  // Theme color palette for charts (using Ant Design blue)
  const chartColors = ['#1677ff', '#4096ff', '#69b1ff', '#91caff', '#bae0ff'];

  const getPieChartOption = () => ({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item' as const,
    },
    legend: {
      orient: 'vertical' as const,
      right: 10,
      top: 'center',
    },
    color: chartColors,
    series: [
      {
        name: t('overview.clusters'),
        type: 'pie' as const,
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 6,
          borderColor: isDark ? '#1e293b' : '#fff',
          borderWidth: 2,
        },
        label: {
          show: false,
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold' as const,
          },
        },
        data: clusters.length > 0
          ? clusters.reduce((acc, cluster) => {
              const type = cluster.type || 'unknown';
              const existing = acc.find((item) => item.name === type);
              if (existing) {
                existing.value += 1;
              } else {
                acc.push({ value: 1, name: type });
              }
              return acc;
            }, [] as { value: number; name: string }[])
          : [{ value: 0, name: t('overview.noData') }],
      },
    ],
  });

  const getBarChartOption = () => ({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis' as const,
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category' as const,
      data: [t('overview.clusters'), t('overview.listeners'), t('overview.resources'), t('menu.plugin')],
    },
    yAxis: {
      type: 'value' as const,
    },
    series: [
      {
        type: 'bar' as const,
        data: [clusters.length, listeners.length, resources.length, pluginGroups.length],
        barWidth: '50%',
        itemStyle: {
          color: '#1677ff',
          borderRadius: [4, 4, 0, 0],
        },
      },
    ],
  });

  // Quick actions
  const quickActions = [
    { key: 'cluster', icon: <ClusterOutlined />, label: t('cluster.shortTitle'), path: '/cluster' },
    { key: 'listener', icon: <AudioOutlined />, label: t('listener.shortTitle'), path: '/listener' },
    { key: 'mapping', icon: <ApiOutlined />, label: t('menu.mapping'), path: '/mapping' },
    { key: 'plugin', icon: <SettingOutlined />, label: t('menu.plugin'), path: '/plugin' },
  ];

  // Calculate health percentage (0% when no instances or data unavailable)
  const healthPercentage = instanceStats
    ? instanceStats.total > 0
      ? Math.round((instanceStats.connected / instanceStats.total) * 100)
      : 0
    : 0;

  // Unified card height styles
  const cardBodyStyle = { height: '100%' };

  return (
    <Spin spinning={isLoading}>
      <div>
        <div className="flex justify-between items-center mb-6">
          <h2 className="m-0 text-xl font-semibold">{t('overview.title')}</h2>
          <Space>
            <Button icon={<SyncOutlined />} onClick={() => window.location.reload()}>
              {t('common.refresh')}
            </Button>
          </Space>
        </div>

        {/* Statistics Cards - Row 1 */}
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} lg={6}>
            <Card hoverable onClick={() => navigate('/cluster')} className="h-full">
              <Statistic
                title={t('cluster.shortTitle')}
                value={clusters.length}
                prefix={<ClusterOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card hoverable onClick={() => navigate('/listener')} className="h-full">
              <Statistic
                title={t('listener.shortTitle')}
                value={listeners.length}
                prefix={<AudioOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card hoverable onClick={() => navigate('/mapping')} className="h-full">
              <Statistic
                title={t('menu.mapping')}
                value={resources.length}
                prefix={<ApiOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card hoverable onClick={() => navigate('/plugin')} className="h-full">
              <Statistic
                title={t('menu.plugin')}
                value={pluginGroups.length}
                prefix={<CloudServerOutlined />}
              />
            </Card>
          </Col>
        </Row>

        {/* Server Status & Quick Actions - Row 2 */}
        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24} lg={12}>
            <Card className="h-full" styles={{ body: cardBodyStyle }}>
              <Row gutter={[16, 16]} align="middle" className="h-full">
                <Col xs={24} md={12}>
                  <div className="flex items-center gap-3">
                    <RocketOutlined className="text-3xl text-slate-600 dark:text-slate-400" />
                    <div>
                      <h3 className="m-0 text-lg font-bold">
                        {t('overview.gatewayName')}
                      </h3>
                      <p className="m-0 text-slate-500 text-sm">
                        {t('overview.gatewayDesc')}
                      </p>
                    </div>
                  </div>
                </Col>
                <Col xs={24} md={12}>
                  <Row gutter={8}>
                    <Col span={8}>
                      <div className="text-center">
                        <Tooltip title={t('overview.httpServer')}>
                          <GlobalOutlined className="text-xl text-slate-500" />
                          <div className="text-xs text-slate-400 mt-1">HTTP</div>
                          <Tag color="blue" className="mt-1">:8081</Tag>
                        </Tooltip>
                      </div>
                    </Col>
                    <Col span={8}>
                      <div className="text-center">
                        <Tooltip title={t('overview.xdsServer')}>
                          <ThunderboltOutlined className="text-xl text-slate-500" />
                          <div className="text-xs text-slate-400 mt-1">xDS</div>
                          <Tag color="purple" className="mt-1">:18000</Tag>
                        </Tooltip>
                      </div>
                    </Col>
                    <Col span={8}>
                      <div className="text-center">
                        <Tooltip title={t('overview.systemStatus')}>
                          <SafetyCertificateOutlined className="text-xl text-slate-500" />
                          <div className="text-xs text-slate-400 mt-1">{t('common.status')}</div>
                          <Tag color="success" className="mt-1">{t('common.running')}</Tag>
                        </Tooltip>
                      </div>
                    </Col>
                  </Row>
                </Col>
              </Row>
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card title={t('overview.quickActions')} className="h-full" styles={{ body: cardBodyStyle }}>
              <Row gutter={[12, 12]}>
                {quickActions.map((action) => (
                  <Col xs={12} sm={6} key={action.key}>
                    <Card
                      hoverable
                      size="small"
                      className="text-center"
                      onClick={() => navigate(action.path)}
                    >
                      <div className="text-xl mb-1">{action.icon}</div>
                      <div className="text-sm font-medium">{action.label}</div>
                      <Button type="link" size="small" icon={<PlusOutlined />} className="p-0">
                        {t('common.add')}
                      </Button>
                    </Card>
                  </Col>
                ))}
              </Row>
            </Card>
          </Col>
        </Row>

        {/* Instance Health & Recent Activity - Row 3 */}
        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24} lg={8}>
            <Card
              title={
                <Space>
                  <HeartOutlined />
                  {t('overview.instanceHealth')}
                </Space>
              }
              className="h-full"
              styles={{ body: { ...cardBodyStyle, minHeight: 280 } }}
            >
              <div className="text-center py-2">
                <Progress
                  type="dashboard"
                  percent={healthPercentage}
                  size={120}
                  strokeColor={healthPercentage === 100 ? '#52c41a' : healthPercentage >= 50 ? '#faad14' : '#ff4d4f'}
                  format={(percent) => (
                    <div>
                      <div className="text-xl font-bold">{percent}%</div>
                      <div className="text-xs text-slate-400">{t('overview.healthy')}</div>
                    </div>
                  )}
                />
                <Row gutter={8} className="mt-3">
                  <Col span={8}>
                    <Statistic
                      title={<span className="text-xs">{t('overview.total')}</span>}
                      value={instanceStats?.total || 0}
                      prefix={<CloudOutlined />}
                      valueStyle={{ fontSize: 16 }}
                    />
                  </Col>
                  <Col span={8}>
                    <Badge status="success" />
                    <Statistic
                      title={<span className="text-xs">{t('overview.connected')}</span>}
                      value={instanceStats?.connected || 0}
                      valueStyle={{ fontSize: 16, color: '#52c41a' }}
                    />
                  </Col>
                  <Col span={8}>
                    <Badge status="error" />
                    <Statistic
                      title={<span className="text-xs">{t('overview.disconnected')}</span>}
                      value={instanceStats?.disconnected || 0}
                      valueStyle={{ fontSize: 16, color: '#ff4d4f' }}
                    />
                  </Col>
                </Row>
              </div>
            </Card>
          </Col>

          <Col xs={24} lg={16}>
            <Card
              title={
                <Space>
                  <HistoryOutlined />
                  {t('overview.recentActivity')}
                </Space>
              }
              className="h-full"
              styles={{ body: { ...cardBodyStyle, minHeight: 280, overflow: 'auto' } }}
            >
              <div className="text-center py-8 text-slate-400">
                <HistoryOutlined className="text-3xl mb-2" />
                <div>{t('overview.noRecentActivity')}</div>
              </div>
            </Card>
          </Col>
        </Row>

        {/* Charts - Row 4 */}
        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24} lg={12}>
            <Card title={t('overview.resourceDistribution')} className="h-full">
              <SafeECharts
                option={getBarChartOption()}
                style={{ height: 260 }}
              />
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card title={t('overview.clusterTypes')} className="h-full">
              <SafeECharts
                option={getPieChartOption()}
                style={{ height: 260 }}
              />
            </Card>
          </Col>
        </Row>

        {/* Data Tables - Row 5 */}
        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24} lg={12}>
            <Card
              title={t('overview.clusters')}
              extra={<Button type="link" size="small" onClick={() => navigate('/cluster')}>{t('common.view')}</Button>}
              className="h-full"
              styles={{ body: { ...cardBodyStyle, minHeight: 200 } }}
            >
              <Table
                columns={clusterColumns}
                dataSource={clusters.slice(0, 4) as Cluster[]}
                rowKey="name"
                pagination={false}
                size="small"
                locale={{ emptyText: t('common.noData') }}
              />
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card
              title={t('overview.listeners')}
              extra={<Button type="link" size="small" onClick={() => navigate('/listener')}>{t('common.view')}</Button>}
              className="h-full"
              styles={{ body: { ...cardBodyStyle, minHeight: 200 } }}
            >
              <Table
                columns={listenerColumns}
                dataSource={listeners.slice(0, 4) as Listener[]}
                rowKey="name"
                pagination={false}
                size="small"
                locale={{ emptyText: t('common.noData') }}
              />
            </Card>
          </Col>
        </Row>

        {/* Recent Resources - Row 6 */}
        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24}>
            <Card
              title={t('overview.recentResources')}
              extra={<Button type="link" size="small" onClick={() => navigate('/mapping')}>{t('common.view')}</Button>}
            >
              <Table
                columns={resourceColumns}
                dataSource={resources.slice(0, 5)}
                rowKey={(record, index) => record.path || String(record.id) || `resource-${index}`}
                pagination={false}
                size="middle"
                locale={{ emptyText: t('common.noData') }}
              />
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  );
}
