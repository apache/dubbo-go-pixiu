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
	"fmt"
	"net/http"
)

import (
	"go.uber.org/zap"

	"v.marlon.life/toolkit/util"
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
}

// RunServer start server
func RunServer() {
	// load config
	global.VP = Viper()
	global.LOG = Zap()

	config.InitEtcdClient()

	account.InitUserDao()
	account.InitGuestDao()

	router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	s := initServer(address, router)

	var wg util.WaitGroupWrapper

	wg.AddAndRun(func() {
		global.LOG.Info("server run success on ", zap.String("address", address))
		fmt.Printf(helperInfo, address)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.LOG.Error(err.Error())
		}
	})

	wg.AddAndRun(func() {
		global.LOG.Info("xDS server run success on :18000")
		if err := StartxDsServer(); err != nil {
			global.LOG.Error(err.Error())
		}
	})

	wg.Wait()
}
