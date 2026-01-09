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
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	perrors "github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

import (
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/api/grpcstub"
)

// server implements the BenchmarkService
type server struct {
	users map[int32]*grpcstub.User
	grpcstub.UnimplementedBenchmarkServiceServer
}

func (s *server) GetUser(ctx context.Context, request *grpcstub.GetUserRequest) (*grpcstub.GetUsersResponse, error) {
	us := make([]*grpcstub.User, 0)
	if request.GetUserId() == 0 {
		for i := int32(1); i <= 2; i++ {
			if u, ok := s.users[i]; ok {
				us = append(us, u)
			}
		}
	} else {
		u, ok := s.users[request.GetUserId()]
		if !ok {
			return &grpcstub.GetUsersResponse{}, perrors.New("Invalid User ID")
		}
		us = append(us, u)
	}
	return &grpcstub.GetUsersResponse{Users: us}, nil
}

func (s *server) GetUsers(ctx context.Context, request *grpcstub.GetUsersRequest) (*grpcstub.GetUsersResponse, error) {
	us := make([]*grpcstub.User, 0)
	for _, userId := range request.UserIds {
		u, ok := s.users[userId]
		if ok {
			us = append(us, u)
		}
	}
	return &grpcstub.GetUsersResponse{Users: us}, nil
}

func (s *server) GetUserByName(ctx context.Context, request *grpcstub.GetUserByNameRequest) (*grpcstub.GetUsersResponse, error) {
	for _, user := range s.users {
		if user.Name == request.Name {
			return &grpcstub.GetUsersResponse{Users: []*grpcstub.User{user}}, nil
		}
	}
	return &grpcstub.GetUsersResponse{}, perrors.New("Invalid User Name")
}

func (s *server) SayHello(ctx context.Context, request *grpcstub.HelloRequest) (*grpcstub.HelloResponse, error) {
	return &grpcstub.HelloResponse{Message: "Hello " + request.Name}, nil
}

func initUsers(s *server) {
	s.users[1] = &grpcstub.User{Id: 1, Name: "Kenway", Age: 25}
	s.users[2] = &grpcstub.User{Id: 2, Name: "Ken", Age: 30}
	s.users[3] = &grpcstub.User{Id: 3, Name: "Moorse", Age: 28}
}

func main() {
	l, err := net.Listen("tcp", ":50001") //nolint:gosec
	if err != nil {
		panic(err)
	}

	s := &server{users: make(map[int32]*grpcstub.User)}
	initUsers(s)

	keepAliveArgs := keepalive.ServerParameters{
		Time:    60 * time.Second,
		Timeout: 5 * time.Second,
	}
	gs := grpc.NewServer(grpc.KeepaliveParams(keepAliveArgs))

	grpcstub.RegisterBenchmarkServiceServer(gs, s)

	// registers the server reflection service on the given gRPC server.
	reflection.Register(gs)

	fmt.Println("grpc benchmark server is now running on :50001...")
	go func() {
		err = gs.Serve(l)
		if err != nil {
			panic(err)
		}
	}()

	initSignal()
	gs.GracefulStop()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(3*time.Second, func() {
				os.Exit(1)
			})
			fmt.Println("grpc benchmark server exit now...")
			return
		}
	}
}
