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
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		_, _ = w.Write([]byte(tokenize(query.Get("secret"), query.Get("key"))))
	})

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"message":"success","status":200}`))
	})
	log.Println("Starting sample server ...")
	log.Fatal(http.ListenAndServe(":1314", nil))
}

func tokenize(secret, salt string) string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", salt, secret)))
}
