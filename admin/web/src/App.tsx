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

import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { ToastProvider } from './components/Toast'
import { Layout } from './components/Layout'
import { LoginPage } from './pages/Login'
import { OverviewPage } from './pages/Overview'
import { MappingPage } from './pages/Mapping'
import { RouteBindingPage } from './pages/RouteBinding'
import { PluginGroupPage } from './pages/PluginGroup'
import { ClusterPage } from './pages/ClusterPage'
import { ListenerPage } from './pages/ListenerPage'
import { RateLimiterPage } from './pages/RateLimiter'
import { OpaPage } from './pages/OpaPage'
import { ProfilePage } from './pages/Profile'

// Note: intentionally no global route guard. Whether an unauthenticated visit
// can operate is decided by API auth and each page's own state.
export default function App() {
  return (
    <ToastProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/*"
            element={
              <Layout>
                <Routes>
                  <Route path="/" element={<Navigate to="/gateway/overview" replace />} />
                  <Route path="/gateway/overview" element={<OverviewPage />} />
                  <Route path="/gateway/mapping/:resourceId" element={<MappingPage />} />
                  <Route path="/gateway/route-binding" element={<RouteBindingPage />} />
                  <Route path="/gateway/plugin-group" element={<PluginGroupPage />} />
                  <Route path="/gateway/cluster" element={<ClusterPage />} />
                  <Route path="/gateway/listener" element={<ListenerPage />} />
                  <Route path="/flow/ratelimit" element={<RateLimiterPage />} />
                  <Route path="/opa" element={<OpaPage />} />
                  <Route path="/profile" element={<ProfilePage />} />
                  <Route path="*" element={<Navigate to="/gateway/overview" replace />} />
                </Routes>
              </Layout>
            }
          />
        </Routes>
      </BrowserRouter>
    </ToastProvider>
  )
}
