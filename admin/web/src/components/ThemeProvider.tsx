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

import { ConfigProvider, theme } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
import { useTranslation } from 'react-i18next';
import { useThemeStore } from '../stores/theme';
import type { ReactNode } from 'react';

interface ThemeProviderProps {
  children: ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  const { isDark } = useThemeStore();
  const { i18n } = useTranslation();

  const locale = i18n.language === 'zh' ? zhCN : enUS;

  return (
    <ConfigProvider
      locale={locale}
      theme={{
        algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: isDark
          ? {
              colorBgContainer: '#000000',
              colorBgElevated: '#141414',
              colorBgLayout: '#000000',
            }
          : {
              colorBgContainer: '#ffffff',
              colorBgElevated: '#ffffff',
              colorBgLayout: '#ffffff',
            },
        components: {
          Layout: isDark
            ? {
                siderBg: '#000000',
                headerBg: '#000000',
                bodyBg: '#000000',
              }
            : {
                siderBg: '#ffffff',
                headerBg: '#ffffff',
                bodyBg: '#ffffff',
              },
          Menu: isDark
            ? {
                darkItemBg: '#000000',
                darkSubMenuItemBg: '#000000',
              }
            : {},
        },
      }}
    >
      {children}
    </ConfigProvider>
  );
}
