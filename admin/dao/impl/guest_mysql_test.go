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
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/go-sql-driver/mysql"
)

func TestRegisterConcurrentUsernameIsUnique(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)

	dao := &GuestDao{db: db}
	if err := dao.EnsureUsernameUniqueIndex(); err != nil {
		t.Fatalf("EnsureUsernameUniqueIndex() error = %v", err)
	}

	const attempts = 8
	start := make(chan struct{})
	results := make(chan error, attempts)
	var workers sync.WaitGroup
	workers.Add(attempts)
	for i := 0; i < attempts; i++ {
		go func() {
			defer workers.Done()
			<-start
			results <- dao.Register("same-user", "hashed-password")
		}()
	}
	close(start)
	workers.Wait()
	close(results)

	successes, duplicates := 0, 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case err.Error() == "user already exists":
			duplicates++
		default:
			t.Errorf("Register() error = %v", err)
		}
	}
	if successes != 1 || duplicates != attempts-1 {
		t.Fatalf("registration results: successes = %d, duplicates = %d; want 1 and %d", successes, duplicates, attempts-1)
	}

	var userCount, roleCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM pixiu_user WHERE username = ?", "same-user").Scan(&userCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM pixiu_user_role").Scan(&roleCount); err != nil {
		t.Fatal(err)
	}
	if userCount != 1 || roleCount != 1 {
		t.Fatalf("persisted rows: users = %d, roles = %d; want 1 each", userCount, roleCount)
	}
}

func TestEnsureUsernameUniqueIndexLeavesExistingDuplicatesUntouched(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)
	if _, err := db.Exec("INSERT INTO pixiu_user (username, password) VALUES (?, ?), (?, ?)", "duplicate", "one", "duplicate", "two"); err != nil {
		t.Fatal(err)
	}

	err := (&GuestDao{db: db}).EnsureUsernameUniqueIndex()
	if err == nil || !strings.Contains(err.Error(), `duplicate username "duplicate" exists`) {
		t.Fatalf("EnsureUsernameUniqueIndex() error = %v, want duplicate username error", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM pixiu_user WHERE username = ?", "duplicate").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("duplicate rows after failed migration = %d, want 2 untouched rows", count)
	}
}

func TestEnsureDefaultAdminSeedsRoleAndAssociation(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)
	dao := &GuestDao{db: db}
	if err := dao.EnsureUsernameUniqueIndex(); err != nil {
		t.Fatalf("EnsureUsernameUniqueIndex() error = %v", err)
	}

	if err := dao.EnsureDefaultAdmin("admin", "initial-hash"); err != nil {
		t.Fatalf("EnsureDefaultAdmin() error = %v", err)
	}
	if err := dao.EnsureDefaultAdmin("admin", "replacement-hash"); err != nil {
		t.Fatalf("EnsureDefaultAdmin() second call error = %v", err)
	}

	var userID int64
	var password string
	var role int
	if err := db.QueryRow("SELECT id, password, role FROM pixiu_user WHERE username = ?", "admin").Scan(&userID, &password, &role); err != nil {
		t.Fatal(err)
	}
	if password != "initial-hash" || role != 1 {
		t.Fatalf("default admin = (password %q, role %d), want (initial-hash, 1)", password, role)
	}

	var roleName, description string
	if err := db.QueryRow("SELECT role_name, description FROM pixiu_role WHERE id = 1").Scan(&roleName, &description); err != nil {
		t.Fatal(err)
	}
	if roleName != "admin" || description != "Administrator" {
		t.Fatalf("admin role = (%q, %q), want (admin, Administrator)", roleName, description)
	}

	var associationCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM pixiu_user_role WHERE user_id = ? AND role_id = 1", userID).Scan(&associationCount); err != nil {
		t.Fatal(err)
	}
	if associationCount != 1 {
		t.Fatalf("admin role associations = %d, want 1", associationCount)
	}

	found, result, err := (&UserDao{db: db}).GetUserRole("admin")
	if err != nil {
		t.Fatalf("GetUserRole() error = %v", err)
	}
	if !found {
		t.Fatal("GetUserRole() found = false, want true")
	}
	adminRole, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("GetUserRole() result has type %T, want map[string]any", result)
	}
	if adminRole["role"] != "admin" || adminRole["description"] != "Administrator" {
		t.Fatalf("GetUserRole() = %#v, want admin role", adminRole)
	}
}

func TestEnsureDefaultAdminPreservesExistingNonAdminAccount(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)
	dao := &GuestDao{db: db}
	if err := dao.EnsureUsernameUniqueIndex(); err != nil {
		t.Fatalf("EnsureUsernameUniqueIndex() error = %v", err)
	}
	if _, err := db.Exec("INSERT INTO pixiu_user (username, password, role) VALUES (?, ?, 0)", "admin", "existing-hash"); err != nil {
		t.Fatal(err)
	}

	if err := dao.EnsureDefaultAdmin("admin", "default-hash"); err != nil {
		t.Fatalf("EnsureDefaultAdmin() error = %v", err)
	}

	var password string
	var role int
	if err := db.QueryRow("SELECT password, role FROM pixiu_user WHERE username = ?", "admin").Scan(&password, &role); err != nil {
		t.Fatal(err)
	}
	if password != "existing-hash" || role != 0 {
		t.Fatalf("existing account = (password %q, role %d), want (existing-hash, 0)", password, role)
	}

	var associationCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM pixiu_user_role").Scan(&associationCount); err != nil {
		t.Fatal(err)
	}
	if associationCount != 0 {
		t.Fatalf("non-admin account role associations = %d, want 0", associationCount)
	}
}

func TestGetUserRoleUsesRequestedUsername(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)
	if _, err := db.Exec(`
		INSERT INTO pixiu_role (id, role_name, description)
		VALUES (1, 'admin', 'Administrator'), (2, 'reader', 'Read only')
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`
		INSERT INTO pixiu_user (id, username, password, role)
		VALUES (1, 'admin', 'admin-hash', 1), (2, 'reader', 'reader-hash', 0)
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO pixiu_user_role (user_id, role_id) VALUES (1, 1), (2, 2)"); err != nil {
		t.Fatal(err)
	}

	found, result, err := (&UserDao{db: db}).GetUserRole("reader")
	if err != nil {
		t.Fatalf("GetUserRole() error = %v", err)
	}
	if !found {
		t.Fatal("GetUserRole() found = false, want true")
	}
	role, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("GetUserRole() result has type %T, want map[string]any", result)
	}
	if role["role"] != "reader" || role["description"] != "Read only" {
		t.Fatalf("GetUserRole() = %#v, want reader role", role)
	}
}

func newGuestMySQLTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("PIXIU_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set PIXIU_TEST_MYSQL_DSN to a MySQL account allowed to create and drop databases")
	}

	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("parse PIXIU_TEST_MYSQL_DSN: %v", err)
	}
	config.DBName = ""
	adminDB, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		t.Fatalf("open MySQL connection: %v", err)
	}
	if err := adminDB.Ping(); err != nil {
		adminDB.Close()
		t.Fatalf("connect to MySQL: %v", err)
	}

	databaseName := fmt.Sprintf("pixiu_test_%d_%d", time.Now().UnixNano(), os.Getpid())
	if _, err := adminDB.Exec("CREATE DATABASE `" + databaseName + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		adminDB.Close()
		t.Fatalf("create isolated MySQL test database: %v", err)
	}

	config.DBName = databaseName
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		_, _ = adminDB.Exec("DROP DATABASE `" + databaseName + "`")
		adminDB.Close()
		t.Fatalf("open isolated MySQL test database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		if _, err := adminDB.Exec("DROP DATABASE `" + databaseName + "`"); err != nil {
			t.Errorf("drop isolated MySQL test database: %v", err)
		}
		_ = adminDB.Close()
	})
	if err := db.Ping(); err != nil {
		t.Fatalf("connect to isolated MySQL test database: %v", err)
	}
	return db
}

func createGuestTestTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`
		CREATE TABLE pixiu_user (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			role INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		t.Fatalf("create pixiu_user test table: %v", err)
	}
	if _, err := db.Exec(`
		CREATE TABLE pixiu_user_role (
			user_id BIGINT NOT NULL,
			role_id INT NOT NULL
		) ENGINE=InnoDB
	`); err != nil {
		t.Fatalf("create pixiu_user_role test table: %v", err)
	}
	if _, err := db.Exec(`
		CREATE TABLE pixiu_role (
			id INT NOT NULL PRIMARY KEY,
			role_name VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL
		) ENGINE=InnoDB
	`); err != nil {
		t.Fatalf("create pixiu_role test table: %v", err)
	}
}
