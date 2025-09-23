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
	pb "dubbo-go-pixiu-benchmark/api"
)

type GreeterProvider struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &pb.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}
