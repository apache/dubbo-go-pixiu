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
	"fmt"
)

import (
	mysql "github.com/go-sql-driver/mysql"

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

func (d *GuestDao) EnsureUsernameUniqueIndex() error {
	indexed, err := hasUniqueUsernameIndex(d.db)
	if err != nil {
		return err
	}
	if indexed {
		return nil
	}

	indexExists, indexValid, err := fixedUsernameIndexStatus(d.db)
	if err != nil {
		return err
	}
	if indexExists && !indexValid {
		return fmt.Errorf("cannot add unique username index: index %q already exists with an incompatible definition; rename or remove it and restart", usernameIndexName)
	}

	if username, err := findDuplicateUsername(d.db); err == nil {
		return fmt.Errorf("cannot add unique username index: duplicate username %q exists; resolve duplicate accounts and restart", username)
	} else if !errors.Is(err, SQL.ErrNoRows) {
		return err
	}

	_, err = d.db.Exec("ALTER TABLE pixiu_user ADD UNIQUE INDEX uniq_pixiu_user_username (username)")
	if err == nil {
		return nil
	}

	// Another Pixiu instance may have added the index while this instance was starting.
	if indexed, checkErr := hasUniqueUsernameIndex(d.db); checkErr == nil && indexed {
		return nil
	}
	if username, duplicateErr := findDuplicateUsername(d.db); duplicateErr == nil {
		return fmt.Errorf("cannot add unique username index: duplicate username %q exists; resolve duplicate accounts and restart", username)
	}
	return fmt.Errorf("failed to add unique username index: %w", err)
}

func hasUniqueUsernameIndex(db *SQL.DB) (bool, error) {
	var indexed bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT index_name
			FROM information_schema.statistics
			WHERE table_schema = DATABASE()
			  AND table_name = 'pixiu_user'
			  AND non_unique = 0
			GROUP BY index_name
			HAVING COUNT(*) = 1 AND MAX(column_name = 'username' AND sub_part IS NULL) = 1
		)
	`).Scan(&indexed)
	return indexed, err
}

const usernameIndexName = "uniq_pixiu_user_username"

func fixedUsernameIndexStatus(db *SQL.DB) (exists, valid bool, err error) {
	var columnCount, matchingColumnCount int
	err = db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(non_unique = 0 AND column_name = 'username' AND sub_part IS NULL), 0)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		  AND table_name = 'pixiu_user'
		  AND index_name = ?
	`, usernameIndexName).Scan(&columnCount, &matchingColumnCount)
	if err != nil {
		return false, false, err
	}
	if columnCount == 0 {
		return false, false, nil
	}
	return true, columnCount == 1 && matchingColumnCount == 1, nil
}

func findDuplicateUsername(db *SQL.DB) (string, error) {
	var username string
	err := db.QueryRow(`
		SELECT username
		FROM pixiu_user
		GROUP BY username
		HAVING COUNT(*) > 1
		LIMIT 1
	`).Scan(&username)
	return username, err
}

func (d *GuestDao) EnsureDefaultAdmin(username, password string) error {
	if _, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS pixiu_admin_migrations (
			migration_key VARCHAR(128) NOT NULL PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB
	`); err != nil {
		return fmt.Errorf("create admin migration table: %w", err)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		INSERT IGNORE INTO pixiu_role (id, role_name, description)
		VALUES (1, 'admin', 'Administrator')
	`); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE pixiu_role SET role_name = 'admin', description = 'Administrator' WHERE id = 1"); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		INSERT IGNORE INTO pixiu_user (username, password, role)
		VALUES (?, ?, 1)
	`, username, password); err != nil {
		return err
	}

	var userID int64
	if err := tx.QueryRow("SELECT id FROM pixiu_user WHERE username = ? FOR UPDATE", username).Scan(&userID); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE pixiu_user SET role = 1 WHERE id = ?", userID); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		INSERT INTO pixiu_user_role (user_id, role_id)
		SELECT ?, 1 FROM DUAL
		WHERE NOT EXISTS (
			SELECT 1 FROM pixiu_user_role WHERE user_id = ? AND role_id = 1
		)
	`, userID, userID); err != nil {
		return err
	}

	const migrationKey = "legacy-admin-accounts-demoted-v1"
	result, err := tx.Exec("INSERT IGNORE INTO pixiu_admin_migrations (migration_key) VALUES (?)", migrationKey)
	if err != nil {
		return err
	}
	applied, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if applied == 1 {
		var userRoleID int
		err := tx.QueryRow("SELECT id FROM pixiu_role WHERE role_name = 'user' ORDER BY id LIMIT 1").Scan(&userRoleID)
		if errors.Is(err, SQL.ErrNoRows) {
			roleResult, insertErr := tx.Exec("INSERT INTO pixiu_role (role_name, description) VALUES ('user', 'User')")
			if insertErr != nil {
				return insertErr
			}
			roleID, insertErr := roleResult.LastInsertId()
			if insertErr != nil {
				return insertErr
			}
			userRoleID = int(roleID)
		} else if err != nil {
			return err
		}

		if _, err := tx.Exec(`
			INSERT INTO pixiu_user_role (user_id, role_id)
			SELECT id, ? FROM pixiu_user
			WHERE role = 1 AND username <> ?
			  AND NOT EXISTS (
				SELECT 1 FROM pixiu_user_role
				WHERE user_id = pixiu_user.id AND role_id = ?
			  )
		`, userRoleID, username, userRoleID); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE pixiu_user SET role = 0 WHERE role = 1 AND username <> ?", username); err != nil {
			return err
		}
		if _, err := tx.Exec(`
			DELETE ur FROM pixiu_user_role AS ur
			JOIN pixiu_user AS u ON u.id = ur.user_id
			WHERE ur.role_id = 1 AND u.username <> ?
		`, username); err != nil {
			return err
		}
	}

	return tx.Commit()
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
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errors.New("user already exists")
		}
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
