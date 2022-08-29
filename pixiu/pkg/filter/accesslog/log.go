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
	"path/filepath"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

// access log config, enable default value true, outputpath default value console
// AccessLogConfig access log will out put into console
type AccessLogConfig struct {
	OutPutPath string `yaml:"outPutPath" json:"outPutPath" mapstructure:"outPutPath" default:"console"`
}

// AccessLogWriter access log chan
type AccessLogWriter struct {
	AccessLogDataChan chan AccessLogData
}

// AccessLogData access log data
type AccessLogData struct {
	AccessLogMsg    string
	AccessLogConfig AccessLogConfig
}

// Writer writer msg into chan
func (alw *AccessLogWriter) Writer(accessLogData AccessLogData) {
	select {
	case alw.AccessLogDataChan <- accessLogData:
		return
	default:
		logger.Warn("the channel is full and the access logIntoChannel data will be dropped")
		return
	}
}

// Write write log into out put path
func (alw *AccessLogWriter) Write() {
	go func() {
		for accessLogData := range alw.AccessLogDataChan {
			alw.writeLogToFile(accessLogData)
		}
	}()
}

// write log to file or console
func (alw *AccessLogWriter) writeLogToFile(ald AccessLogData) {
	alc := ald.AccessLogConfig
	alm := ald.AccessLogMsg
	if len(alc.OutPutPath) == 0 || alc.OutPutPath == constant.Console {
		logger.Info(alm)
		return
	}
	_ = WriteToFile(alm, alc.OutPutPath)
}

// WriteToFile write message to access log file
func WriteToFile(accessLogMsg string, filePath string) error {
	pd := filepath.Dir(filePath)
	if _, err := os.Stat(pd); err != nil {
		if os.IsExist(err) {
			logger.Warnf("can not open log dir: %s, %v", filePath, err)
		}
		err = os.MkdirAll(pd, os.ModePerm)
		if err != nil {
			logger.Warnf("can not create log dir: %s, %v", filePath, err)
			return err
		}
	}
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, constant.LogFileMode)
	if err != nil {
		logger.Warnf("can not open the access log file: %s, %v", filePath, err)
		return err
	}
	now := time.Now().Format(constant.FileDateFormat)
	fileInfo, err := logFile.Stat()
	if err != nil {
		logger.Warnf("can not get the info of access log file: %s, %v", filePath, err)
		return err
	}
	last := fileInfo.ModTime().Format(constant.FileDateFormat)

	// this is confused.
	// for example, if the last = '2020-03-04'
	// and today is '2020-03-05'
	// we will create one new file to log access data
	// By this way, we can split the access log based on days.
	if now != last {
		err = os.Rename(fileInfo.Name(), fileInfo.Name()+"."+now)
		if err != nil {
			logger.Warnf("can not rename access log file: %s, %v", fileInfo.Name(), err)
			return err
		}
		logFile, err = os.OpenFile(fileInfo.Name(), os.O_CREATE|os.O_APPEND|os.O_RDWR, constant.LogFileMode)
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
