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

package mq

import (
	"strings"
	"testing"
)

func TestClient_Apply(t *testing.T) {
	c := Client{}

	err := c.Apply()

	if err == nil {
		t.Error("Apply should return an error")
	}
	if !strings.Contains(err.Error(), "Apply not implemented") {
		t.Errorf("Apply error should contain 'Apply not implemented', got %s", err.Error())
	}
}