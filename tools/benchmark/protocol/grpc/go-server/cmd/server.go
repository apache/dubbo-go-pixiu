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
	"fmt"
	perrors "github.com/pkg/errors"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"google.golang.org/grpc"
)

import (
	"dubbo-go-pixiu-benchmark/protocol/grpc/proto"
)

const (
	MsgUserNotFound          = "user not found"
	MsgUserQuerySuccessfully = "user(s) query successfully"
)

// Test Cases
// curl http://127.0.0.1:8881/api/v1/provider.UserProvider/GetUser
// curl http://127.0.0.1:8881/api/v1/provider.UserProvider/GetUser -X POST -d '{"userId":1}'

type server struct {
	users map[int32]*proto.User
	proto.UnimplementedUserProviderServer
}

func (s *server) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUsersResponse, error) {
	us := make([]*proto.User, 0)
	if request.GetUserId() == 0 {
		for i := 1; i <= 2; i++ {
			us = append(us, s.users[int32(i)])
		}
	} else {
		u, ok := s.users[request.GetUserId()]
		if !ok {
			return &proto.GetUsersResponse{}, perrors.New("Invalid User ID")
		}
		us = append(us, u)
	}
	return &proto.GetUsersResponse{Users: us}, nil
}

func (s *server) GetUsers(ctx context.Context, request *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	us := make([]*proto.User, 0)
	for _, userId := range request.UserId {
		u, ok := s.users[userId]
		if ok {
			us = append(us, u)
		}
	}
	return &proto.GetUsersResponse{Users: us}, nil
}

func (s *server) GetUserByName(ctx context.Context, request *proto.GetUserByNameRequest) (*proto.GetUsersResponse, error) {
	for i, user := range s.users {
		if user.Name == request.Name {
			return &proto.GetUsersResponse{Users: []*proto.User{s.users[i]}}, nil
		}
	}
	return &proto.GetUsersResponse{}, perrors.New("Invalid User Name")
}

func initUsers(s *server) {
	s.users[1] = &proto.User{UserId: 1, Name: "Kenway"}
	s.users[2] = &proto.User{UserId: 2, Name: "Ken"}
}

func main() {
	l, err := net.Listen("tcp", ":50001") //nolint:gosec
	if err != nil {
		panic(err)
	}

	s := &server{users: make(map[int32]*proto.User)}
	initUsers(s)

	keepAliveArgs := keepalive.ServerParameters{
		Time:    60 * time.Second,
		Timeout: 5 * time.Second,
	}
	gs := grpc.NewServer(grpc.KeepaliveParams(keepAliveArgs))

	proto.RegisterUserProviderServer(gs, s)

	// registers the server reflection service on the given gRPC server.
	reflection.Register(gs)

	fmt.Println("grpc test server is now running...")
	go func() {
		err = gs.Serve(l)
		if err != nil {
			panic(err)
		}
	}()

	initSignal()
	gs.GracefulStop() // handle request until all of them is done
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(int(3e9)), func() {
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("provider app exit now...")
			return
		}
	}
}
