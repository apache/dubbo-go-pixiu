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

package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

import (
	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/global"
	"github.com/apache/dubbo-go-pixiu/admin/initialize"
	"github.com/apache/dubbo-go-pixiu/admin/logic/account"
)

var (
	helperInfo = `
	Welcome DUBBOGO-PIXIU-ADMIN
	Default doc address: http://127.0.0.1%s/swagger/index.html
	Default running address: http://127.0.0.1:8080
`
)

type server interface {
	ListenAndServe() error
	// Shutdown gracefully shuts down the server without interrupting any
	// active connections.
	Shutdown(context.Context) error
}

// RunServer start server. The admin HTTP server and the xDS gRPC server are
// started together under runCoordinated: the first one to fail cancels the
// other, and its error is returned so the CLI can surface a startup failure
// instead of leaving the process false-healthy.
func RunServer(configPath string) error {
	// load config
	vp, err := Viper(configPath)
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}
	global.VP = vp
	global.LOG = Zap()

	config.InitEtcdClient()

	account.InitUserDao()
	account.InitGuestDao()

	router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	s := initServer(address, router)

	return runCoordinated(context.Background(),
		func(ctx context.Context) error { return runHTTPServer(ctx, s, address) },
		func(ctx context.Context) error { return StartxDsServer(ctx) },
	)
}

// runCoordinated runs the given workers concurrently under a shared context.
// The first worker to return a non-nil error cancels the others and that error
// is returned; any further errors are dropped. If all workers exit cleanly (for
// example after the context is canceled), nil is returned.
func runCoordinated(ctx context.Context, workers ...func(context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for _, w := range workers {
		wg.Add(1)
		go func(w func(context.Context) error) {
			defer wg.Done()
			if err := w(ctx); err != nil {
				// Cancel peers so they initiate shutdown, then keep only the
				// first error.
				cancel()
				select {
				case errCh <- err:
				default:
				}
			}
		}(w)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

// runHTTPServer serves the admin HTTP API until ctx is canceled. A real
// ListenAndServe failure (e.g. port already in use) is returned so the
// coordinator can shut the xDS server down and propagate it to the CLI.
func runHTTPServer(ctx context.Context, s server, address string) error {
	global.LOG.Info("server run success on ", zap.String("address", address))
	fmt.Printf(helperInfo, address)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	case <-ctx.Done():
		// A peer failed; shut the HTTP server down gracefully. Shutdown
		// unblocks ListenAndServe, which then reports ErrServerClosed.
		_ = s.Shutdown(ctx)
		<-serveErr
		return nil
	}
}
