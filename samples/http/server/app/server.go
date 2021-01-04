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
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
)

func main() {
	http.HandleFunc("/user/", user)
	log.Println("Starting sample server ...")
	log.Fatal(http.ListenAndServe(":1314", nil))
}

func user(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		byts, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		var user User
		err = json.Unmarshal(byts, &user)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		_, ok := cache.Get(user.Name)
		if ok {
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write([]byte("{\"message\":\"data is exist\"}"))
			return
		}
		user.ID = randSeq(5)
		if cache.Add(&user) {
			b, _ := json.Marshal(&user)
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write(b)
			return
		}
		w.Write(nil)
	case "GET":
		subPath := strings.TrimPrefix(r.URL.Path, "/user/")
		userName := strings.Split(subPath, "/")[0]
		var u *User
		var b bool
		if len(userName) != 0 {
			log.Printf("paths: %v", userName)
			u, b = cache.Get(userName)
		} else {
			q := r.URL.Query()
			u, b = cache.Get(q.Get("name"))
		}
		//w.WriteHeader(200)
		if b {
			b, _ := json.Marshal(u)
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write(b)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
