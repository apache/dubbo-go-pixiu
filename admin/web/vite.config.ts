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

import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

const spaBypass = (req: { method?: string; url?: string }) =>
  req.method === 'GET' || req.method === 'HEAD' ? req.url : undefined

export default defineConfig(({ mode }) => {
  const { VITE_BACKEND_URL = 'http://127.0.0.1:8081' } = loadEnv(mode, '.', 'VITE_')
  return {
    plugins: [react()],
    server: {
      host: '0.0.0.0',
      port: 8088,
      proxy: {
        '/login': { target: VITE_BACKEND_URL, changeOrigin: true, bypass: spaBypass },
        '/register': { target: VITE_BACKEND_URL, changeOrigin: true },
        '/user': VITE_BACKEND_URL,
        '/config': VITE_BACKEND_URL,
        '/swagger': VITE_BACKEND_URL,
      },
    },
  }
})
