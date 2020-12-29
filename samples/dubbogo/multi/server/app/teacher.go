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
	"sync"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(TeacherProvider))
	// ------for hessian2------
	hessian.RegisterPOJO(&Teacher{})

	teacherCache = &TeacherDB{
		nameIndex: make(map[string]*Teacher, 16),
		codeIndex: make(map[int64]*Teacher, 16),
		lock:      sync.Mutex{},
	}

	teacherCache.Add(&Teacher{ID: "0001", Code: 1, Name: "tc-teacher", Age: 18, Time: time.Now()})
	teacherCache.Add(&Teacher{ID: "0002", Code: 2, Name: "ic-teacher", Age: 88, Time: time.Now()})
}

var teacherCache *TeacherDB

// TeacherDB cache teacher.
type TeacherDB struct {
	// key is name, value is Teacher obj
	nameIndex map[string]*Teacher
	// key is code, value is Teacher obj
	codeIndex map[int64]*Teacher
	lock      sync.Mutex
}

// nolint
func (db *TeacherDB) Add(u *Teacher) bool {
	db.lock.Lock()
	defer db.lock.Unlock()

	if u.Name == "" || u.Code <= 0 {
		return false
	}

	if !db.existName(u.Name) && !db.existCode(u.Code) {
		return db.AddForName(u) && db.AddForCode(u)
	}

	return false
}

// nolint
func (db *TeacherDB) AddForName(u *Teacher) bool {
	if len(u.Name) == 0 {
		return false
	}

	if _, ok := db.nameIndex[u.Name]; ok {
		return false
	}

	db.nameIndex[u.Name] = u
	return true
}

// nolint
func (db *TeacherDB) AddForCode(u *Teacher) bool {
	if u.Code <= 0 {
		return false
	}

	if _, ok := db.codeIndex[u.Code]; ok {
		return false
	}

	db.codeIndex[u.Code] = u
	return true
}

// nolint
func (db *TeacherDB) GetByName(n string) (*Teacher, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.nameIndex[n]
	return r, ok
}

// nolint
func (db *TeacherDB) GetByCode(n int64) (*Teacher, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.codeIndex[n]
	return r, ok
}

func (db *TeacherDB) existName(name string) bool {
	if len(name) <= 0 {
		return false
	}

	_, ok := db.nameIndex[name]
	if ok {
		return true
	}

	return false
}

func (db *TeacherDB) existCode(code int64) bool {
	if code <= 0 {
		return false
	}

	_, ok := db.codeIndex[code]
	if ok {
		return true
	}

	return false
}

// Teacher teacher obj.
type Teacher struct {
	ID   string    `json:"id,omitempty"`
	Code int64     `json:"code,omitempty"`
	Name string    `json:"name,omitempty"`
	Age  int32     `json:"age,omitempty"`
	Time time.Time `json:"time,omitempty"`
}

// TeacherProvider the dubbo provider.
// like: version: 1.0.0 group: test
type TeacherProvider struct {
}

// CreateTeacher new teacher, Proxy config POST.
func (s *TeacherProvider) CreateTeacher(ctx context.Context, Teacher *Teacher) (*Teacher, error) {
	println("Req CreateTeacher data:%#v", Teacher)
	if Teacher == nil {
		return nil, errors.New("not found")
	}
	_, ok := teacherCache.GetByName(Teacher.Name)
	if ok {
		return nil, errors.New("data is exist")
	}

	b := teacherCache.Add(Teacher)
	if b {
		return Teacher, nil
	}

	return nil, errors.New("add error")
}

// GetTeacherByName query by name, single param, Proxy config GET.
func (s *TeacherProvider) GetTeacherByName(ctx context.Context, name string) (*Teacher, error) {
	println("Req GetTeacherByName name:%#v", name)
	r, ok := teacherCache.GetByName(name)
	if !ok {
		return nil, nil
	}

	println("Req GetTeacherByName result:%#v", r)
	return r, nil
}

// GetTeacherByCode query by code, single param, Proxy config GET.
func (s *TeacherProvider) GetTeacherByCode(ctx context.Context, code int64) (*Teacher, error) {
	println("Req GetTeacherByCode name:%#v", code)
	r, ok := teacherCache.GetByCode(code)
	if !ok {
		return nil, nil
	}
	println("Req GetTeacherByCode result:%#v", r)
	return r, nil
}

// GetTeacherTimeout query by name, will timeout for proxy.
func (s *TeacherProvider) GetTeacherTimeout(ctx context.Context, name string) (*Teacher, error) {
	println("Req GetTeacherByName name:%#v", name)
	// sleep 10s, proxy config less than 10s.
	time.Sleep(10 * time.Second)
	r, ok := teacherCache.GetByName(name)
	if !ok {
		return nil, nil
	}
	println("Req GetTeacherByName result:%#v", r)
	return r, nil
}

// GetTeacherByNameAndAge query by name and age, two params, Proxy config GET.
func (s *TeacherProvider) GetTeacherByNameAndAge(ctx context.Context, name string, age int32) (*Teacher, error) {
	println("Req GetTeacherByNameAndAge name:%s, age:%d", name, age)
	r, ok := teacherCache.GetByName(name)
	if ok && r.Age == age {
		println("Req GetTeacherByNameAndAge result:%#v", r)
		return r, nil
	}
	return r, nil
}

// UpdateTeacher update by teacher struct, my be another struct, Proxy config POST or PUT.
func (s *TeacherProvider) UpdateTeacher(ctx context.Context, teacher *Teacher) (bool, error) {
	println("Req UpdateTeacher data:%#v", teacher)
	r, ok := teacherCache.GetByName(teacher.Name)
	if !ok {
		return false, errors.New("not found")
	}
	if teacher.ID != "" {
		r.ID = teacher.ID
	}
	if teacher.Age >= 0 {
		r.Age = teacher.Age
	}
	return true, nil
}

// UpdateTeacher update by teacher struct, my be another struct, Proxy config POST or PUT.
func (s *TeacherProvider) UpdateTeacherByName(ctx context.Context, name string, teacher *Teacher) (bool, error) {
	println("Req UpdateTeacherByName data:%#v", teacher)
	r, ok := teacherCache.GetByName(name)
	if !ok {
		return false, errors.New("not found")
	}
	if teacher.ID != "" {
		r.ID = teacher.ID
	}
	if teacher.Age >= 0 {
		r.Age = teacher.Age
	}
	return true, nil
}

// nolint
func (s *TeacherProvider) Reference() string {
	return "TeacherProvider"
}

// nolint
func (s Teacher) JavaClassName() string {
	return "com.dubbogo.proxy.TeacherService"
}
