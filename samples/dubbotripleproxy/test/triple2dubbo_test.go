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

package test

import (
	"context"
	"testing"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	tpconst "github.com/dubbogo/triple/pkg/common/constant"
)

func TestTriple2Dubbo(t *testing.T) {
	tripleRefConf := newTripleRefConf("com.dubbogo.pixiu.DubboUserService", tpconst.TRIPLE)
	resp, err := tripleRefConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUserByName",
		[]string{"java.lang.String"},
		[]hessian.Object{"tc"},
	)

	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser1(userId string) res: %+v", resp)
}

func newTripleRefConf(iface, protocol string) config.ReferenceConfig {

	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		Cluster:       "failover",
		RegistryIDs:   []string{"zk"},
		Protocol:      protocol,
		Generic:       "true",
		Group:         "test",
		Version:       "1.0.0",
		URL:           "tri://127.0.0.1:9999/" + iface + "?" + constant.SerializationKey + "=hessian2",
	}

	rootConfig := config.NewRootConfigBuilder().
		Build()
	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad(appName)

	return refConf
}
