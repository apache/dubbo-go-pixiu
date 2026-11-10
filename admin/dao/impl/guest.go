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

func (d *GuestDao) EnsureDefaultAdmin(username, password string) error {
	_, err := d.db.Exec(`
		INSERT INTO pixiu_user (username, password, role)
		SELECT ?, ?, 1 FROM DUAL
		WHERE NOT EXISTS (SELECT 1 FROM pixiu_user WHERE username = ?)
	`, username, password, username)
	return err
}

func (d *GuestDao) Register(username, password string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRow("SELECT id FROM pixiu_user WHERE username = ? LIMIT 1", username).Scan(&userID)
	if err == nil {
		return errors.New("user already exists")
	}
	if !errors.Is(err, SQL.ErrNoRows) {
		return err
	}

	result, err := tx.Exec("INSERT INTO pixiu_user (username, password, role) VALUES (?, ?, 1)", username, password)
	if err != nil {
		return errors.New("failed to create user")
	}
	userID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	if _, err = tx.Exec("INSERT INTO pixiu_user_role (user_id, role_id) VALUES (?, 1)", userID); err != nil {
		return errors.New("failed to assign administrator role")
	}
	return tx.Commit()
}

func (d *GuestDao) CheckLogin() {
	panic("implement me")
}
