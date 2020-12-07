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
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	_ "github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

const (
	// nolint
	FileDateFormat = "2006-01-02"
	// nolint
	MessageDateLayout = "2006-01-02 15:04:05"
	// nolint
	LogFileMode = 0600

	// console
	Console = "console"
)

func init() {
	extension.SetFilterFunc(constant.AccessLogFilter, Filter())
}

func Filter() context.FilterFunc {
	return func(c context.Context) {
		alc := config.GetBootstrap().StaticResources.AccessLogConfig
		if alc.Enable {
			start := time.Now()
			c.Next()
			latency := time.Now().Sub(start)
			// build access_log message
			accessLogMsg := buildAccessLogMsg(c, latency)
			if len(alc.OutPutPath) == 0 || alc.OutPutPath == Console {
				logger.Info(accessLogMsg)
			} else {
				go func() {
					err := writeToFile(accessLogMsg, alc.OutPutPath)
					if err != nil {
						logger.Warnf("write access log to file err s%, v%", accessLogMsg, err)
					}
				}()
			}
		}
	}
}

// write message to access log file
func writeToFile(accessLogMsg string, filePath string) error {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, LogFileMode)
	if err != nil {
		logger.Warnf("Can not open the access log file: %s, %v", filePath, err)
		return err
	}
	now := time.Now().Format(FileDateFormat)
	fileInfo, err := logFile.Stat()
	if err != nil {
		logger.Warnf("Can not get the info of access log file: %s, %v", filePath, err)
		return err
	}
	last := fileInfo.ModTime().Format(FileDateFormat)

	// this is confused.
	// for example, if the last = '2020-03-04'
	// and today is '2020-03-05'
	// we will create one new file to log access data
	// By this way, we can split the access log based on days.
	if now != last {
		err = os.Rename(fileInfo.Name(), fileInfo.Name()+"."+now)
		if err != nil {
			logger.Warnf("Can not rename access log file: %s, %v", fileInfo.Name(), err)
			return err
		}
		logFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, LogFileMode)
		if err != nil {
			logger.Warnf("can not open access log file: %s, %v", fileInfo.Name(), err)
			return err
		}
	}
	_, err = logFile.WriteString(accessLogMsg + "\n")
	if err != nil {
		logger.Warnf("can not write to access log file: %s, v%", fileInfo.Name(), err)
		return err
	}
	return nil
}

func buildAccessLogMsg(c context.Context, cost time.Duration) string {
	req := c.(*http.HttpContext).Request
	valueStr := req.URL.Query().Encode()
	if len(valueStr) != 0 {
		valueStr = strings.ReplaceAll(valueStr, "&", ",")
	}

	builder := strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(time.Now().Format(MessageDateLayout))
	builder.WriteString("] ")
	builder.WriteString(req.RemoteAddr)
	builder.WriteString(" -> ")
	builder.WriteString(req.Host)
	builder.WriteString(" - ")
	if len(valueStr) > 0 {
		builder.WriteString("request params: [")
		builder.WriteString(valueStr)
		builder.WriteString("] ")
	}
	builder.WriteString("cost time [ ")
	builder.WriteString(strconv.Itoa(int(cost)) + " ]")
	return builder.String()
}
