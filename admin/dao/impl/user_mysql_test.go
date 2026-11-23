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
	"sync"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

func TestEditPasswordAllowsOnlyOneConcurrentUpdate(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)
	if _, err := db.Exec("INSERT INTO pixiu_user (username, password) VALUES (?, ?)", "user", utils.Md5("old-password")); err != nil {
		t.Fatal(err)
	}

	dao := &UserDao{db: db}
	start := make(chan struct{})
	type result struct {
		updated bool
		err     error
	}
	results := make(chan result, 2)
	var workers sync.WaitGroup
	for _, password := range []string{"new-password-one", "new-password-two"} {
		workers.Add(1)
		go func(password string) {
			defer workers.Done()
			<-start
			updated, err := dao.EditPassword("old-password", password, "user")
			results <- result{updated: updated, err: err}
		}(password)
	}
	close(start)
	workers.Wait()
	close(results)

	successes, failures := 0, 0
	for result := range results {
		switch {
		case result.updated && result.err == nil:
			successes++
		case !result.updated && result.err != nil:
			failures++
		default:
			t.Errorf("EditPassword() result = (%t, %v), want one success and one stale-password failure", result.updated, result.err)
		}
	}
	if successes != 1 || failures != 1 {
		t.Fatalf("concurrent password updates: successes = %d, failures = %d; want 1 each", successes, failures)
	}

	var storedPassword string
	if err := db.QueryRow("SELECT password FROM pixiu_user WHERE username = ?", "user").Scan(&storedPassword); err != nil {
		t.Fatal(err)
	}
	if storedPassword != utils.Md5("new-password-one") && storedPassword != utils.Md5("new-password-two") {
		t.Fatalf("stored password hash = %q, want one of the concurrent updates", storedPassword)
	}
}

func TestEditPasswordRejectsStalePassword(t *testing.T) {
	db := newGuestMySQLTestDB(t)
	createGuestTestTables(t, db)
	if _, err := db.Exec("INSERT INTO pixiu_user (username, password) VALUES (?, ?)", "user", utils.Md5("current-password")); err != nil {
		t.Fatal(err)
	}

	updated, err := (&UserDao{db: db}).EditPassword("stale-password", "new-password", "user")
	if updated || err == nil {
		t.Fatalf("EditPassword() = (%t, %v), want false and an error for stale password", updated, err)
	}
	var storedPassword string
	if err := db.QueryRow("SELECT password FROM pixiu_user WHERE username = ?", "user").Scan(&storedPassword); err != nil {
		t.Fatal(err)
	}
	if storedPassword != utils.Md5("current-password") {
		t.Fatalf("stored password changed after stale update: %q", storedPassword)
	}
}
