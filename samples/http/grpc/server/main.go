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

package main

import (
	"context"
	"net"
)

import (
	"google.golang.org/grpc"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/samples/http/grpc/proto"
)

const (
	MsgUserNotFound          = "user not found"
	MsgUserQuerySuccessfully = "user(s) query successfully"
)

type server struct {
	users map[int32]*proto.User
	proto.UnimplementedUserProviderServer
}

func (s *server) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	us := make([]*proto.User, 0)
	if request.GetUserId() == 0 {
		for _, user := range s.users {
			us = append(us, user)
		}
	} else {
		u, ok := s.users[request.GetUserId()]
		if !ok {
			return &proto.GetUserResponse{Message: MsgUserNotFound}, nil
		}
		us = append(us, u)
	}
	return &proto.GetUserResponse{Message: MsgUserQuerySuccessfully, Users: us}, nil
}

func main() {
	l, err := net.Listen("tcp", ":50001") //nolint:gosec
	if err != nil {
		panic(err)
	}

	s := &server{users: make(map[int32]*proto.User)}
	initUsers(s)

	gs := grpc.NewServer()
	proto.RegisterUserProviderServer(gs, s)
	logger.Info("grpc test server is now running...")
	err = gs.Serve(l)
	if err != nil {
		panic(err)
	}
}

func initUsers(s *server) {
	s.users[1] = &proto.User{UserId: 1, Name: "Kenway"}
	s.users[2] = &proto.User{UserId: 2, Name: "Ken"}
}
