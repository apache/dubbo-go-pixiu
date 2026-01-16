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

// Package doc contains the Swagger/OpenAPI documentation for Pixiu Admin API.
//
// The swagger.json file in this directory contains the OpenAPI 3.0 specification
// for all RESTful API endpoints provided by the Pixiu Admin backend.
//
// API Documentation:
//   - Base URL: http://localhost:8081/api
//   - Authentication: JWT token in "token" header
//   - Response format: JSON with {code, message, data} structure
//
// Main API Groups:
//   - Auth: Login, Register
//   - User: Profile management
//   - Clusters: Backend cluster configuration
//   - Listeners: Protocol listener configuration
//   - Resources: API mapping configuration
//   - Methods: HTTP method configuration
//   - Plugins: Plugin group management
//   - Instances: Instance statistics
//   - Users: User management (Admin)
//   - Roles: Role management (Admin)
//   - Permissions: Permission management (Admin)
package doc
