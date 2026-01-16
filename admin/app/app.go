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

// Package app provides the public API for starting the admin server.
package app

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/server"
)

// Server wraps the internal server implementation.
type Server struct {
	srv *server.Server
}

// New creates a new Server instance.
func New(configPath string) (*Server, error) {
	srv, err := server.New(configPath)
	if err != nil {
		return nil, err
	}
	return &Server{srv: srv}, nil
}

// Run starts the server and blocks until shutdown.
func (s *Server) Run() error {
	return s.srv.Run()
}
