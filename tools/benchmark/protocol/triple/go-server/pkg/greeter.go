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

package pkg

import (
	"context"
)

import (
	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/api"
)

// Compile-time check to ensure BenchmarkProvider implements BenchmarkServiceHandler
var _ api.BenchmarkServiceHandler = (*BenchmarkProvider)(nil)

// BenchmarkProvider implements the BenchmarkServiceHandler for Triple protocol
type BenchmarkProvider struct {
	users map[int32]*api.User
}

func NewBenchmarkProvider() *BenchmarkProvider {
	p := &BenchmarkProvider{
		users: make(map[int32]*api.User),
	}
	p.initUsers()
	return p
}

func (p *BenchmarkProvider) initUsers() {
	p.users[1] = &api.User{Id: 1, Name: "Kenway", Age: 25}
	p.users[2] = &api.User{Id: 2, Name: "Ken", Age: 30}
	p.users[3] = &api.User{Id: 3, Name: "Moorse", Age: 28}
}

func (p *BenchmarkProvider) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUsersResponse, error) {
	logger.Infof("Triple BenchmarkProvider GetUser, userId: %d", req.UserId)
	us := make([]*api.User, 0)
	if req.GetUserId() == 0 {
		for i := int32(1); i <= 2; i++ {
			if u, ok := p.users[i]; ok {
				us = append(us, u)
			}
		}
	} else {
		if u, ok := p.users[req.GetUserId()]; ok {
			us = append(us, u)
		}
	}
	return &api.GetUsersResponse{Users: us}, nil
}

func (p *BenchmarkProvider) GetUsers(ctx context.Context, req *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	logger.Infof("Triple BenchmarkProvider GetUsers, userIds: %v", req.UserIds)
	us := make([]*api.User, 0)
	for _, userId := range req.UserIds {
		if u, ok := p.users[userId]; ok {
			us = append(us, u)
		}
	}
	return &api.GetUsersResponse{Users: us}, nil
}

func (p *BenchmarkProvider) GetUserByName(ctx context.Context, req *api.GetUserByNameRequest) (*api.GetUsersResponse, error) {
	logger.Infof("Triple BenchmarkProvider GetUserByName, name: %s", req.Name)
	for _, user := range p.users {
		if user.Name == req.Name {
			return &api.GetUsersResponse{Users: []*api.User{user}}, nil
		}
	}
	return &api.GetUsersResponse{}, nil
}

func (p *BenchmarkProvider) SayHello(ctx context.Context, req *api.HelloRequest) (*api.HelloResponse, error) {
	logger.Infof("Triple BenchmarkProvider SayHello, name: %s", req.Name)
	return &api.HelloResponse{Message: "Hello " + req.Name}, nil
}
