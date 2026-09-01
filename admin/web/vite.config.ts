import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// The Admin API (gin, default :8081) does not ship a CORS middleware, so the
// dev server proxies every backend path. GET requests on /login and /register
// are SPA navigations and must bypass the proxy (the backend only accepts POST).
const backend = 'http://127.0.0.1:8081'

const spaBypass = (req: { method?: string; url?: string }) => {
  if (req.method === 'GET' || req.method === 'HEAD') return req.url
  return undefined
}

export default defineConfig({
  plugins: [react()],
  server: {
    port: 8088,
    proxy: {
      '/login': { target: backend, changeOrigin: true, bypass: spaBypass },
      '/register': { target: backend, changeOrigin: true, bypass: spaBypass },
      '/user': backend,
      '/config': backend,
      '/swagger': backend,
    },
  },
})
