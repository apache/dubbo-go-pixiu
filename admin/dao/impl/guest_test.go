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

package impl

import (
	"strings"
	"testing"
)

func TestGuestDao_CheckLogin(t *testing.T) {
	d := &GuestDao{}
	err := d.CheckLogin()
	if err == nil || !strings.Contains(err.Error(), "session validation is not implemented") {
		t.Fatalf("unimplemented session validation must return an error, got %v", err)
	}
}
