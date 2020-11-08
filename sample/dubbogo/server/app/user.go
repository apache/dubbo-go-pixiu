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

	cache.Add(&User{Id: "0001", Name: "tc", Age: 18, Time: time.Now()})
}

var cache *UserDB

type UserDB struct {
	// key is name, value is user obj
	cacheMap map[string]*User
	lock     sync.Mutex
}

func (db *UserDB) Add(u *User) bool {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.cacheMap[u.Name] = u
	return true
}

func (db *UserDB) Get(n string) (*User, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.cacheMap[n]
	return r, ok
}

type User struct {
	Id   string
	Name string
	Age  int32
	Time time.Time
}

// version: 1.0.0 group: test
// methods: all
type UserProvider struct {
}

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

func (u *UserProvider) GetUserByName(ctx context.Context, name string) (*User, error) {
	println("Req GetUserByName data:%#v", name)
	r, ok := cache.Get(name)
	if ok {
		println("Req GetUserByName result:%#v", r)
		return r, nil
	}

	return nil, nil
}

func (u *UserProvider) QueryUser(ctx context.Context, user *User) (*User, error) {
	println("Req QueryUser data:%#v", user)
	rsp := User{user.Id, user.Name, user.Age, time.Now()}
	println("Req QueryUser data:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) TwoSimpleParams(ctx context.Context, name string, age int32) (*User, error) {
	println("Req TwoSimpleParams name:%s, age:%d", name, age)
	rsp := User{"1", name, age, time.Now()}
	println("Req TwoSimpleParams data:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) TwoSimpleParamsWithError(ctx context.Context, name string, age int32) (*User, error) {
	println("Req TwoSimpleParamsWithError name:%s, age:%d", name, age)
	rsp := User{"1", name, age, time.Now()}
	println("Req TwoSimpleParamsWithError data:%#v", rsp)
	return &rsp, errors.New("TwoSimpleParams error")
}

// 方法名称映射，做参考
// GetUserByName can call success.
//func (u *UserProvider) MethodMapper() map[string]string {
//	return map[string]string{
//		"GetUserByName": "queryUserByName",
//	}
//}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func (u User) JavaClassName() string {
	return "com.ikurento.user.User"
}

func println(format string, args ...interface{}) {
	fmt.Printf("\033[32;40m"+format+"\033[0m\n", args...)
}
