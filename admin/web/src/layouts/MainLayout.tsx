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
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Button, Dropdown, Space, Avatar, message, theme } from 'antd';
import {
  DashboardOutlined,
  ApiOutlined,
  ClusterOutlined,
  AudioOutlined,
  AppstoreOutlined,
  UserOutlined,
  SettingOutlined,
  SunOutlined,
  MoonOutlined,
  GlobalOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  KeyOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useThemeStore } from '../stores/theme';
import { useAuthStore } from '../stores/auth';
import type { MenuProps } from 'antd';

const { Header, Sider, Content } = Layout;

export function MainLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const { t, i18n } = useTranslation();
  const { isDark, toggleTheme } = useThemeStore();
  const { username, logout } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();
  const { token } = theme.useToken();

  const menuItems: MenuProps['items'] = [
    {
      key: '/overview',
      icon: <DashboardOutlined />,
      label: t('menu.overview'),
    },
    {
      key: '/mapping',
      icon: <ApiOutlined />,
      label: t('menu.mapping'),
    },
    {
      key: '/cluster',
      icon: <ClusterOutlined />,
      label: t('menu.cluster'),
    },
    {
      key: '/listener',
      icon: <AudioOutlined />,
      label: t('menu.listener'),
    },
    {
      key: '/plugin',
      icon: <AppstoreOutlined />,
      label: t('menu.plugin'),
    },
    {
      key: '/user-management',
      icon: <TeamOutlined />,
      label: t('menu.userManagement'),
    },
    {
      key: '/role-management',
      icon: <SafetyCertificateOutlined />,
      label: t('menu.roleManagement'),
    },
    {
      key: '/permission-management',
      icon: <KeyOutlined />,
      label: t('menu.permissionManagement'),
    },
  ];

  const languageItems: MenuProps['items'] = [
    {
      key: 'en',
      label: t('language.en'),
      onClick: () => i18n.changeLanguage('en'),
    },
    {
      key: 'zh',
      label: t('language.zh'),
      onClick: () => i18n.changeLanguage('zh'),
    },
  ];

  const handleLogout = async () => {
    await logout();
    message.success(t('common.logoutSuccess') || 'Logged out successfully');
    navigate('/login');
  };

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: t('menu.profile'),
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: t('menu.settings'),
      onClick: () => navigate('/settings'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: t('common.logout') || 'Logout',
      onClick: handleLogout,
    },
  ];

  const handleMenuClick: MenuProps['onClick'] = ({ key }) => {
    navigate(key);
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={240}
        style={{
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          zIndex: 100,
        }}
      >
        {/* Logo */}
        <div
          style={{
            height: 56,
            display: 'flex',
            alignItems: 'center',
            justifyContent: collapsed ? 'center' : 'flex-start',
            padding: collapsed ? 0 : '0 20px',
          }}
        >
          <div
            style={{
              width: 32,
              height: 32,
              borderRadius: 8,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              background: token.colorPrimary,
              color: '#fff',
              fontWeight: 600,
              fontSize: 16,
            }}
          >
            P
          </div>
          {!collapsed && (
            <span
              style={{
                marginLeft: 12,
                fontSize: 16,
                fontWeight: 600,
              }}
            >
              Pixiu
            </span>
          )}
        </div>

        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{
            borderRight: 0,
          }}
        />
      </Sider>

      <Layout
        style={{
          marginLeft: collapsed ? 80 : 240,
          transition: 'margin-left 0.2s ease',
        }}
      >
        <Header
          style={{
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            height: 56,
            position: 'sticky',
            top: 0,
            zIndex: 99,
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />

          <Space size={8}>
            <Dropdown menu={{ items: languageItems }} placement="bottom">
              <Button
                type="text"
                icon={<GlobalOutlined />}
              />
            </Dropdown>
            <Button
              type="text"
              icon={isDark ? <SunOutlined /> : <MoonOutlined />}
              onClick={toggleTheme}
            />
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer', padding: '4px 8px', borderRadius: 6, marginLeft: 8 }}>
                <Avatar size={32}>
                  {username?.[0]?.toUpperCase() || 'U'}
                </Avatar>
                <span style={{ fontSize: 14 }}>
                  {username || 'User'}
                </span>
              </Space>
            </Dropdown>
          </Space>
        </Header>

        <Content
          style={{
            margin: 20,
            padding: 24,
            minHeight: 'calc(100vh - 96px)',
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
