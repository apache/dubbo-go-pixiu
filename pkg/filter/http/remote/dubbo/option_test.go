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

package dubbo

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestValuesOpt(t *testing.T) {
	opt := &valuesOpt{}
	target := &dubboTarget{
		Values: make([]interface{}, 3),
		Types:  make([]string, 3),
	}
	vals := []interface{}{
		struct{ Name string }{"joe"},
		"abc",
		123,
	}
	err := opt.Action(target, [2]interface{}{
		vals,
		"object, string, int",
	})
	assert.Nil(t, err)
	assert.Equal(t, len(target.Values), 3)
	assert.Equal(t, target.Values, vals)
	assert.Equal(t, target.Types[0], "object")
	assert.Equal(t, target.Types[1], "string")
	assert.Equal(t, target.Types[2], "int")

	err = opt.Action(target, []interface{}{
		vals,
		"object, string, int",
	})
	assert.NotNil(t, err)

	target = &dubboTarget{
		Values: make([]interface{}, 3),
		Types:  make([]string, 3),
	}
	vals = []interface{}{
		struct{ Name string }{"joe"},
		"abc",
		123,
	}
	err = opt.Action(target, [2]interface{}{
		vals,
		"",
	})
	assert.Nil(t, err)
	assert.Equal(t, len(target.Values), 3)
	assert.Equal(t, target.Values, vals)
	assert.Equal(t, len(target.Types), 0)
}

func TestParamTypesOptAction(t *testing.T) {
	opt := &paramTypesOpt{}
	target := &dubboTarget{
		Values: make([]interface{}, 3),
		Types:  make([]string, 3),
	}
	err := opt.Action(target, "object,string")
	assert.Nil(t, err)
	assert.Equal(t, "object", target.Types[0])
	assert.Equal(t, "string", target.Types[1])

	err = opt.Action(target, "object,whatsoever")
	assert.EqualError(t, err, "Types invalid whatsoever")

	err = opt.Action("target", []string{})
	assert.EqualError(t, err, "The val type must be string")
	err = opt.Action(target, "object,")
	assert.Nil(t, err)
	assert.Equal(t, "object", target.Types[0])
	err = opt.Action(target, "object")
	assert.Nil(t, err)
	assert.Equal(t, "object", target.Types[0])
}
