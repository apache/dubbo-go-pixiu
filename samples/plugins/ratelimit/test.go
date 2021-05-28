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
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{Timeout: time.Second * 2}
	url := "http://localhost:8888/api/v1/test-dubbo/user?name=tc"

	ch := make(chan interface{})
	for i := 0; i < 5; i++ {
		go func() {
			for {
				access(client, url)
			}
		}()
	}
	<-ch
}

func access(c *http.Client, url string) {
	res, err := c.Get(url)
	time.Sleep(time.Millisecond * 100)
	if err != nil {
		fmt.Println(time.Now(), "blocked", err.Error())
		return
	}
	fmt.Println(time.Now(), "passed")
	_ = res.Body.Close()
}
