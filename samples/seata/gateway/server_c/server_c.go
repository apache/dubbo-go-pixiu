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

type Inventory struct {
	ID  int64 `json:"id"`
	Qty int64 `json:"qty"`
}

func main() {
	r := gin.Default()

	r.POST("/service-c/try", func(context *gin.Context) {
		inv := &Inventory{}
		err := context.BindJSON(inv)
		if err == nil {
			context.JSON(200, gin.H{
				"success": true,
				"message": fmt.Sprintf("inventory %d tried!", inv.ID),
			})
			return
		}
		context.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
	})

	r.POST("/service-c/confirm", func(context *gin.Context) {
		inv := &Inventory{}
		err := context.BindJSON(inv)
		if err == nil {
			fmt.Println(fmt.Sprintf("inventory %d confirmed!", inv.ID))
			context.JSON(200, gin.H{
				"success": true,
				"message": fmt.Sprintf("inventory %d confirmed!", inv.ID),
			})
			return
		}
		context.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
	})

	r.POST("/service-c/cancel", func(context *gin.Context) {
		inv := &Inventory{}
		err := context.BindJSON(inv)
		if err == nil {
			fmt.Println(fmt.Sprintf("inventory %d canceled!", inv.ID))
			context.JSON(200, gin.H{
				"success": true,
				"message": fmt.Sprintf("inventory %d canceled!", inv.ID),
			})
			return
		}
		context.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
	})

	r.Run(":8082")
}
