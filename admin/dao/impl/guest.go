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

package impl

import (
	SQL "database/sql"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/dao"
	"github.com/apache/dubbo-go-pixiu/admin/dao/database"
)

type GuestDao struct {
	db *SQL.DB
}

// TODO
func NewGuestDao() *GuestDao {
	return &GuestDao{
		db: database.GetConnection(),
	}
}

func (d *GuestDao) Create(db *SQL.DB) (any, error) {
	d.db = db
	var i dao.GuestDao = d
	return &i, nil
}

func (d *GuestDao) Login(username, password string) (bool, int) {
	db := d.db
	var id int
	err := db.QueryRow("SELECT id FROM pixiu_user WHERE username = ? AND password = ?;", username, password).Scan(&id)
	if err != nil {
		return false, 0
	}
	return true, id
}

func (d *GuestDao) Register(username, password string) error {

	if username == "" {
		return errors.New("void username")
	}

	db := d.db
	var id int
	err := db.QueryRow("SELECT id FROM pixiu_user WHERE username = ?", username).Scan(&id)
	if err == nil {
		return errors.New("User already exists, please login")
	}
	//now := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO pixiu_user (username,password) VALUES (?,?)")
	if err != nil {
		return errors.New("Illegal SQL statement!")
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password)
	if err != nil {
		return errors.New("Failed to create data!")
	}
	// TODO Set transactions and dynamically set user roles
	_ = db.QueryRow("SELECT id FROM pixiu_user WHERE username = ?", username).Scan(&id)
	stmt, _ = db.Prepare("INSERT INTO pixiu_user_role(user_id, role_id) VALUE (?, ?)")
	_, _ = stmt.Exec(id, 1)
	return err
}

func (d *GuestDao) CheckLogin() {
	panic("implement me")
}
