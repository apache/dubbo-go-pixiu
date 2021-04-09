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
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	"github.com/urfave/cli"
)

import (
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/accesslog"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/logger"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/recovery"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/remote"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/response"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/timeout"
)

// Version pixiu version
var Version = "0.1.0"

// main pixiu run method
func main() {
	app := newPXApp(&cmdStart)

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Run(os.Args)
}

func newPXApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "dubbogo pixiu"
	app.Version = Version
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo"
	app.Usage = "Dubbogo pixiu is a lightweight gateway."
	app.Flags = cmdStart.Flags

	// commands
	app.Commands = []cli.Command{
		cmdStart,
		cmdStop,
		cmdReload,
	}

	// action
	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}

		return startCmd.Action.(func(c *cli.Context) error)(c)
	}

	return app
}
