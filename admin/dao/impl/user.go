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
	"database/sql"
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/dao"
	"github.com/apache/dubbo-go-pixiu/admin/dao/database"
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao() *UserDao {
	return &UserDao{
		db: database.GetConnection(),
	}
}

func (d *UserDao) Create(db *sql.DB) (any, error) {
	d.db = db
	var i dao.UserDao = d
	return &i, nil
}

func (d *UserDao) EditPassword(oldPassword, newPassword, username string) (bool, error) {
	db := d.db
	var id int
	var password string
	now := time.Now().Format("2006-01-02 15:04:05")
	oldPassword = utils.Md5(oldPassword)
	newPassword = utils.Md5(newPassword)
	err := db.QueryRow("SELECT id, password FROM  pixiu_user WHERE username = ? AND password = ?;", username, oldPassword).Scan(&id, &password)
	if err != nil {
		return false, errors.New("oldPassword error")
	}
	stmt, err := db.Prepare("UPDATE pixiu_user SET password = ?, date_updated = ? WHERE id = ? AND username = ?;")
	if err != nil {
		return false, errors.New("illegal sql statement")
	}
	defer stmt.Close()
	r, err := stmt.Exec(newPassword, now, id, username)
	if err != nil {
		return false, errors.New("fail to update data, Exec fail")
	}
	_, err = r.RowsAffected()
	if err != nil {
		return false, errors.New("fail to update data, RowAffected fail")
	}
	return true, err
}

func (d *UserDao) GetUserInfo(username string) (bool, any, error) {
	db := d.db
	var userId, role int
	err := db.QueryRow("SELECT id, username, `role` FROM pixiu_user WHERE  username = ?", username).Scan(&userId, &username, &role)
	if err != nil {
		return false, "", errors.New("This user does not exist!")
	}
	userInfo := map[string]any{
		"userId":   userId,
		"username": username,
		"role":     role,
	}
	return true, userInfo, err
}

func (d *UserDao) GetUserRole(username string) (bool, any, error) {
	db := d.db
	var roleId int
	var role, description string
	err := db.QueryRow("SELECT role FROM pixiu_user WHERE username = ?", username).Scan(&roleId)
	if err != nil {
		return false, "", err
	}
	err = db.QueryRow("SELECT role_name, description FROM pixiu_role where id = (SELECT role_id FROM pixiu_user_role WHERE user_id = 1)").Scan(&role, &description)
	if err != nil {
		return false, "", err
	}
	result := map[string]any{
		"role":        role,
		"description": description,
	}
	return true, result, err
}

func (d *UserDao) CheckUserIsAdmin(username string) (bool, error) {
	db := d.db
	err := db.QueryRow("SELECT username FROM  pixiu_user WHERE username = ? AND role  = 1", username).Scan(&username)
	if err != nil {
		return false, errors.New("This user is not admin")
	}
	return true, err
}
