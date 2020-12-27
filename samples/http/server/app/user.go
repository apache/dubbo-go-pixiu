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
	"sync"
	"time"
)

func init() {
	cache = &UserDB{
		cacheMap: make(map[string]*User, 16),
		lock:     sync.Mutex{},
	}

	cache.Add(&User{ID: "0001", Name: "tc", Age: 18, Time: time.Now()})
	cache.Add(&User{ID: "0002", Name: "ic", Age: 88, Time: time.Now()})
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
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Age  int32     `json:"age"`
	Time time.Time `json:"time"`
}
