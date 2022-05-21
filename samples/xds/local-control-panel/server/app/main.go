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
	"flag"
	"os"
)

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
)

var (
	l      Logger
	port   = uint(18000)
	nodeID = "test-id"
)

func init() {
	l = Logger{}
	l.Debug = true
}

func main() {
	flag.Parse()

	// Create a snaphost
	snaphost := cache.NewSnapshotCache(false, cache.IDHash{}, l)

	go func() {
		// Create the config that we'll serve to Envoy
		config := GenerateSnapshotPixiu()
		if err := config.Consistent(); err != nil {
			l.Errorf("config inconsistency: %+v\n%+v", config, err)
			os.Exit(1)
		}
		l.Debugf("will serve config %+v", config)

		// Add the config to the snaphost
		if err := snaphost.SetSnapshot(context.Background(), nodeID, config); err != nil {
			l.Errorf("config error %q for %+v", err, config)
			os.Exit(1)
		}

		//time.Sleep(30 * time.Second)
		//config2 := GenerateSnapshotPixiu2()
		//if err := config2.Consistent(); err != nil {
		//	l.Errorf("config inconsistency: %+v\n%+v", config, err)
		//	os.Exit(1)
		//}
		//if err := snaphost.SetSnapshot(context.Background(), nodeID, config2); err != nil {
		//	l.Errorf("config error %q for %+v", err, config)
		//	os.Exit(1)
		//}
	}()

	// Run the xDS server
	ctx := context.Background()
	cb := &test.Callbacks{Debug: l.Debug}
	srv := server.NewServer(ctx, snaphost, cb)
	RunServer(ctx, srv, port)
}
