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

package account

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

type bootstrapGuestDao struct {
	username string
	password string
}

func (d *bootstrapGuestDao) Login(string, string) (bool, int) { return false, 0 }
func (d *bootstrapGuestDao) EnsureUsernameUniqueIndex() error { return nil }
func (d *bootstrapGuestDao) EnsureDefaultAdmin(username, password string) error {
	d.username, d.password = username, password
	return nil
}
func (d *bootstrapGuestDao) Register(string, string) error { return nil }
func (d *bootstrapGuestDao) CheckLogin()                   {}

func TestEnsureDefaultAdminUsesBootstrapCredentials(t *testing.T) {
	previous := guestDao
	spy := &bootstrapGuestDao{}
	guestDao = spy
	t.Cleanup(func() { guestDao = previous })

	if err := EnsureDefaultAdmin(); err != nil {
		t.Fatalf("EnsureDefaultAdmin() error = %v", err)
	}
	if spy.username != "admin" {
		t.Errorf("bootstrap username = %q, want admin", spy.username)
	}
	if want := utils.Md5("pixiuadmin"); spy.password != want {
		t.Errorf("bootstrap password hash = %q, want %q", spy.password, want)
	}
}
