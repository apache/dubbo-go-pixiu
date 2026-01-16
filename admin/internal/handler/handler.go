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

package handler

import (
	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/store"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
)

// Handler contains all HTTP handlers and their dependencies.
type Handler struct {
	mysql     *store.MySQL
	etcd      *store.Etcd
	logger    *zap.Logger
	jwtConfig config.JWTConfig
}

// New creates a new Handler with the given dependencies.
func New(mysql *store.MySQL, etcd *store.Etcd, logger *zap.Logger, jwtConfig config.JWTConfig) *Handler {
	return &Handler{
		mysql:     mysql,
		etcd:      etcd,
		logger:    logger,
		jwtConfig: jwtConfig,
	}
}
