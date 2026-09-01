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
