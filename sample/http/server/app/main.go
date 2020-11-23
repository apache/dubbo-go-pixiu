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
	"net/http"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
)

func main() {
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/user", user)
	log.Println("Starting sample server ...")
	log.Fatal(http.ListenAndServe(":1314", nil))
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, I'm http"))
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
		if cache.Add(&user) {
			b, _ := json.Marshal(&user)
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write(b)
			return
		}
		w.Write(nil)
	case "GET":
		q := r.URL.Query()
		u, b := cache.Get(q.Get("name"))
		//w.WriteHeader(200)
		if b {
			b, _ := json.Marshal(u)
			w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			w.Write(b)
			return
		}
		w.Write(nil)
	}
}
