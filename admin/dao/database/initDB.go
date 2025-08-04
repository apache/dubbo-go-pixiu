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

package database

import (
	"database/sql"
	"fmt"
	"sync"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
)

// MysqlDriver
const MysqlDriver = "mysql"

var db *sql.DB
var initDb sync.Once

func Init(mysqlConf config.MysqlConfig) {
	var err error

	//username, password, host, port, dbname := config.Bootstrap.GetMysqlConfig()
	//dataSourceName := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8"
	//db, err = sql.Open(MysqlDriver, dataSourceName)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlConf.Username, mysqlConf.Password,
		mysqlConf.Host, mysqlConf.Port, mysqlConf.Dbname)
	db, err = sql.Open(MysqlDriver, connectStr)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
}

func GetConnection() *sql.DB {
	initDb.Do(func() {
		Init(config.Bootstrap.MysqlConfig)
	})
	return db
}
