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
)

import (
	"github.com/gin-gonic/gin"
)

type Account struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"Amount"`
}

func main() {
	r := gin.Default()

	r.POST("/service-b/try", func(context *gin.Context) {
		account := &Account{}
		err := context.ShouldBindJSON(account)
		if err == nil {
			context.JSON(200, gin.H{
				"success": true,
				"message": fmt.Sprintf("account %d tried!", account.ID),
			})
			return
		}
		context.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
	})

	r.POST("/service-b/confirm", func(context *gin.Context) {
		account := &Account{}
		err := context.BindJSON(account)
		if err == nil {
			fmt.Println(fmt.Sprintf("account %d confirmed!", account.ID))
			context.JSON(200, gin.H{
				"success": true,
				"message": fmt.Sprintf("account %d confirmed!", account.ID),
			})
			return
		}
		context.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
	})

	r.POST("/service-b/cancel", func(context *gin.Context) {
		account := &Account{}
		err := context.BindJSON(account)
		if err == nil {
			fmt.Println(fmt.Sprintf("account %d canceled!", account.ID))
			context.JSON(200, gin.H{
				"success": true,
				"message": fmt.Sprintf("account %d canceled!", account.ID),
			})
			return
		}
		context.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
	})

	r.Run(":8081")
}
