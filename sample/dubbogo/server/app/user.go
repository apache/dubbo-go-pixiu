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

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(UserProvider))
	// ------for hessian2------
	hessian.RegisterPOJO(&User{})

	cache = &UserDB{
		cacheMap: make(map[string]*User, 16),
		lock:     sync.Mutex{},
	}

	cache.Add(&User{ID: "0001", Name: "tc", Age: 18, Time: time.Now()})
}

var cache *UserDB

// UserDB cache user.
type UserDB struct {
	// key is name, value is user obj
	cacheMap map[string]*User
	lock     sync.Mutex
}

// nolint.
func (db *UserDB) Add(u *User) bool {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.cacheMap[u.Name] = u
	return true
}

// nolint.
func (db *UserDB) Get(n string) (*User, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.cacheMap[n]
	return r, ok
}

// User user obj.
type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}

// UserProvider the dubbo provider.
// like: version: 1.0.0 group: test
type UserProvider struct {
}

// CreateUser new user, Proxy config POST.
func (u *UserProvider) CreateUser(ctx context.Context, user *User) (*User, error) {
	println("Req CreateUser data:%#v", user)
	if user == nil {
		return nil, errors.New("not found")
	}
	_, ok := cache.Get(user.Name)
	if ok {
		return nil, errors.New("data is exist")
	}

	cache.Add(user)
	return user, nil
}

// GetUserByName query by name, single param, Proxy config GET.
func (u *UserProvider) GetUserByName(ctx context.Context, name string) (*User, error) {
	println("Req GetUserByName name:%#v", name)
	r, ok := cache.Get(name)
	if ok {
		println("Req GetUserByName result:%#v", r)
		return r, nil
	}
	return nil, nil
}

// GetUserTimeout query by name, will timeout for proxy.
func (u *UserProvider) GetUserTimeout(ctx context.Context, name string) (*User, error) {
	println("Req GetUserByName name:%#v", name)
	// sleep 10s, proxy config less than 10s.
	time.Sleep(10 * time.Second)
	r, ok := cache.Get(name)
	if ok {
		println("Req GetUserByName result:%#v", r)
		return r, nil
	}
	return nil, nil
}

// GetUserByNameAndAge query by name and age, two params, Proxy config GET.
func (u *UserProvider) GetUserByNameAndAge(ctx context.Context, name string, age int32) (*User, error) {
	println("Req GetUserByNameAndAge name:%s, age:%d", name, age)
	r, ok := cache.Get(name)
	if ok && r.Age == age {
		println("Req GetUserByName result:%#v", r)
		return r, nil
	}
	return r, nil
}

// UpdateUser update by user struct, my be another struct, Proxy config POST or PUT.
func (u *UserProvider) UpdateUser(ctx context.Context, user *User) (bool, error) {
	println("Req UpdateUser data:%#v", user)
	r, ok := cache.Get(user.Name)
	if ok {
		if user.ID != "" {
			r.ID = user.ID
		}
		if user.Age >= 0 {
			r.Age = user.Age
		}
		return true, nil
	}
	return false, errors.New("not found")
}

// UpdateUser update by user struct, my be another struct, Proxy config POST or PUT.
func (u *UserProvider) UpdateUserByName(ctx context.Context, name string, user *User) (bool, error) {
	println("Req UpdateUser data:%#v", user)
	r, ok := cache.Get(name)
	if ok {
		if user.ID != "" {
			r.ID = user.ID
		}
		if user.Age >= 0 {
			r.Age = user.Age
		}
		return true, nil
	}
	return false, errors.New("not found")
}

// nolint.
func (u *UserProvider) Reference() string {
	return "UserProvider"
}

// nolint.
func (u User) JavaClassName() string {
	return "com.ikurento.user.User"
}

// nolint.
func println(format string, args ...interface{}) {
	fmt.Printf("\033[32;40m"+format+"\033[0m\n", args...)
}
