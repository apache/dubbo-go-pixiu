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
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/api"
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/protocol/triple/go-server/pkg"
)

func main() {
	// Create server using new API
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	// Register BenchmarkService handler
	if err := api.RegisterBenchmarkServiceHandler(srv, pkg.NewBenchmarkProvider()); err != nil {
		panic(err)
	}

	// Start server in goroutine
	go func() {
		if err := srv.Serve(); err != nil {
			logger.Errorf("server serve error: %v", err)
		}
	}()

	fmt.Println("triple benchmark server is now running on :20000...")
	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(3*time.Second, func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})
			fmt.Println("triple benchmark server exit now...")
			return
		}
	}
}
