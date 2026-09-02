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
