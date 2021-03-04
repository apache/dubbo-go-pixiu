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
	config.SetProviderService(new(StudentProvider))
	// ------for hessian2------
	hessian.RegisterPOJO(&Student{})

	studentCache = &StudentDB{
		nameIndex: make(map[string]*Student, 16),
		codeIndex: make(map[int64]*Student, 16),
		lock:      sync.Mutex{},
	}

	studentCache.Add(&Student{ID: "0001", Code: 1, Name: "tc-student", Age: 18, Time: time.Now()})
	studentCache.Add(&Student{ID: "0002", Code: 2, Name: "ic-student", Age: 88, Time: time.Now()})
}

var studentCache *StudentDB

// StudentDB cache student.
type StudentDB struct {
	// key is name, value is student obj
	nameIndex map[string]*Student
	// key is code, value is student obj
	codeIndex map[int64]*Student
	lock      sync.Mutex
}

// nolint
func (db *StudentDB) Add(u *Student) bool {
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
func (db *StudentDB) AddForName(u *Student) bool {
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
func (db *StudentDB) AddForCode(u *Student) bool {
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
func (db *StudentDB) GetByName(n string) (*Student, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.nameIndex[n]
	return r, ok
}

// nolint
func (db *StudentDB) GetByCode(n int64) (*Student, bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	r, ok := db.codeIndex[n]
	return r, ok
}

func (db *StudentDB) existName(name string) bool {
	if len(name) <= 0 {
		return false
	}

	_, ok := db.nameIndex[name]
	if ok {
		return true
	}

	return false
}

func (db *StudentDB) existCode(code int64) bool {
	if code <= 0 {
		return false
	}

	_, ok := db.codeIndex[code]
	if ok {
		return true
	}

	return false
}

// Student student obj.
type Student struct {
	ID   string    `json:"id,omitempty"`
	Code int64     `json:"code,omitempty"`
	Name string    `json:"name,omitempty"`
	Age  int32     `json:"age,omitempty"`
	Time time.Time `json:"time,omitempty"`
}

// StudentProvider the dubbo provider.
// like: version: 1.0.0 group: test
type StudentProvider struct {
}

// CreateStudent new Student, PX config POST.
func (s *StudentProvider) CreateStudent(ctx context.Context, student *Student) (*Student, error) {
	outLn("Req CreateStudent data:%#v", student)
	if student == nil {
		return nil, errors.New("not found")
	}
	_, ok := studentCache.GetByName(student.Name)
	if ok {
		return nil, errors.New("data is exist")
	}

	b := studentCache.Add(student)
	if b {
		return student, nil
	}

	return nil, errors.New("add error")
}

// GetStudentByName query by name, single param, PX config GET.
func (s *StudentProvider) GetStudentByName(ctx context.Context, name string) (*Student, error) {
	outLn("Req GetStudentByName name:%#v", name)
	r, ok := studentCache.GetByName(name)
	if !ok {
		return nil, nil
	}
	outLn("Req GetStudentByName result:%#v", r)
	return r, nil
}

// GetStudentByCode query by code, single param, PX config GET.
func (s *StudentProvider) GetStudentByCode(ctx context.Context, code int64) (*Student, error) {
	outLn("Req GetStudentByCode name:%#v", code)
	r, ok := studentCache.GetByCode(code)
	if !ok {
		return nil, nil
	}
	outLn("Req GetStudentByCode result:%#v", r)
	return r, nil
}

// GetStudentTimeout query by name, will timeout for pixiu.
func (s *StudentProvider) GetStudentTimeout(ctx context.Context, name string) (*Student, error) {
	outLn("Req GetStudentByName name:%#v", name)
	// sleep 10s, pixiu config less than 10s.
	time.Sleep(10 * time.Second)
	r, ok := studentCache.GetByName(name)
	if !ok {
		return nil, nil
	}
	outLn("Req GetStudentByName result:%#v", r)
	return r, nil
}

// GetStudentByNameAndAge query by name and age, two params, PX config GET.
func (s *StudentProvider) GetStudentByNameAndAge(ctx context.Context, name string, age int32) (*Student, error) {
	outLn("Req GetStudentByNameAndAge name:%s, age:%d", name, age)
	r, ok := studentCache.GetByName(name)
	if ok && r.Age == age {
		outLn("Req GetStudentByNameAndAge result:%#v", r)
		return r, nil
	}
	return r, nil
}

// UpdateStudent update by Student struct, my be another struct, PX config POST or PUT.
func (s *StudentProvider) UpdateStudent(ctx context.Context, student *Student) (bool, error) {
	outLn("Req UpdateStudent data:%#v", student)
	r, ok := studentCache.GetByName(student.Name)
	if !ok {
		return false, errors.New("not found")
	}
	if len(student.ID) != 0 {
		r.ID = student.ID
	}
	if student.Age >= 0 {
		r.Age = student.Age
	}
	return true, nil
}

// UpdateStudent update by Student struct, my be another struct, PX config POST or PUT.
func (s *StudentProvider) UpdateStudentByName(ctx context.Context, name string, student *Student) (bool, error) {
	outLn("Req UpdateStudentByName data:%#v", student)
	r, ok := studentCache.GetByName(name)
	if !ok {
		return false, errors.New("not found")
	}
	if len(student.ID) != 0 {
		r.ID = student.ID
	}
	if student.Age >= 0 {
		r.Age = student.Age
	}
	return true, nil
}

// nolint
func (s *StudentProvider) Reference() string {
	return "StudentProvider"
}

// nolint
func (s Student) JavaClassName() string {
	return "com.dubbogo.pixiu.StudentService"
}

// nolint
func outLn(format string, args ...interface{}) {
	fmt.Printf("\033[32;40m"+format+"\033[0m\n", args...)
}
