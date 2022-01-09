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

package test

import (
	"context"
	"time"
)

// User user obj.
type User struct {
	ID   string    `json:"id,omitempty"`
	Code int64     `json:"code,omitempty"`
	Name string    `json:"name,omitempty"`
	Age  int32     `json:"age,omitempty"`
	Time time.Time `json:"time,omitempty"`
}

// UserProvider the dubbo provider.
// like: version: 1.0.0 group: test
type UserProvider struct {
	CreateUser          func(ctx context.Context, user *User) (*User, error)
	GetUserByName       func(ctx context.Context, name string) (*User, error)
	GetUserByCode       func(ctx context.Context, code int64) (*User, error)
	GetUserTimeout      func(ctx context.Context, name string) (*User, error)
	GetUserByNameAndAge func(ctx context.Context, name string, age int32) (*User, error)
	UpdateUser          func(ctx context.Context, user *User) (bool, error)
	UpdateUserByName    func(ctx context.Context, name string, user *User) (bool, error)
}

// nolint
func (u *UserProvider) Reference() string {
	return "UserProvider"
}

// nolint
func (u User) JavaClassName() string {
	return "com.dubbogo.pixiu.User"
}
