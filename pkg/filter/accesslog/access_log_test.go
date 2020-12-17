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

package accesslog

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessLog_Write_to_file(t *testing.T) {
	msg := "this is test msg"
	filePath := "/home/dubbo-go-proxy/logs/dubbo-go-access"
	//filePath := "C:\\Users\\60125\\Desktop\\dubbo-go\\logs\\dubbo-go-proxy-20201217"
	accessLogWriter := &model.AccessLogWriter{AccessLogDataChan: make(chan model.AccessLogData, constant.LogDataBuffer)}
	accessLogWriter.Write()
	accessLogWriter.Writer(model.AccessLogData{AccessLogMsg: msg, AccessLogConfig: model.AccessLogConfig{OutPutPath: filePath, Enable: true}})
	time.Sleep(3e9)
	assert.FileExists(t, filePath, nil)
}
