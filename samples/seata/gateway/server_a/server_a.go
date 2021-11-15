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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

type Account struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"Amount"`
}

type Inventory struct {
	ID  int64 `json:"id"`
	Qty int64 `json:"qty"`
}

func main() {
	r := gin.Default()

	r.POST("/service-a/begin", func(context *gin.Context) {
		xid := context.Request.Header.Get("x_seata_xid")
		fmt.Println(fmt.Sprintf("xid %s", xid))
		account := &Account{
			ID:     1000024549,
			Amount: 200,
		}
		inv := &Inventory{
			ID:  1000000005,
			Qty: 2,
		}

		accountReq, err := json.Marshal(account)
		if err != nil {
			context.JSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
		invReq, err := json.Marshal(inv)
		if err != nil {
			context.JSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
		req1, err := http.NewRequest("POST", "http://localhost:2047/service-b/try", bytes.NewBuffer(accountReq))
		if err != nil {
			context.JSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
		req1.Header.Set("x_seata_xid", xid)

		req2, err := http.NewRequest("POST", "http://localhost:2048/service-c/try", bytes.NewBuffer(invReq))
		if err != nil {
			context.JSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
		req2.Header.Set("x_seata_xid", xid)

		client := &http.Client{}
		result1, err := client.Do(req1)
		if err != nil {
			context.JSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
		if result1.StatusCode != http.StatusOK {
			result1.Write(context.Writer)
			return
		}

		result2, err := client.Do(req2)
		if err != nil {
			context.JSON(400, gin.H{"success": false, "message": err.Error()})
			return
		}
		if result2.StatusCode == http.StatusOK {
			result2.Write(context.Writer)
			return
		}
		context.JSON(400, gin.H{"success": false, "message": "there is a error"})
	})

	r.Run(":8080")
}
