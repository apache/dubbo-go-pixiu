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

package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

type Builder struct {
	hash hashWriter
}

type hashWriter interface {
	Write([]byte) (int, error)
	Sum([]byte) []byte
}

func NewBuilder() *Builder {
	return &Builder{hash: sha256.New()}
}

func (b *Builder) AddString(value string) {
	_, _ = b.hash.Write([]byte(value))
	_, _ = b.hash.Write([]byte{0})
}

func (b *Builder) AddJSON(value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	b.AddString(string(data))
	return nil
}

func (b *Builder) Sum() string {
	return hex.EncodeToString(b.hash.Sum(nil))
}

func (b *Builder) SumShort(n int) string {
	sum := b.Sum()
	if n <= 0 || n >= len(sum) {
		return sum
	}
	return sum[:n]
}

func StringsOrdered(items []string) string {
	if len(items) == 0 {
		return ""
	}
	b := NewBuilder()
	for _, item := range items {
		b.AddString(item)
	}
	return b.Sum()
}

func StringsSorted(items []string) string {
	if len(items) == 0 {
		return ""
	}
	cp := make([]string, len(items))
	copy(cp, items)
	sort.Strings(cp)
	return StringsOrdered(cp)
}

func JSONStable(value any) (string, error) {
	b := NewBuilder()
	if err := b.AddJSON(value); err != nil {
		return "", err
	}
	return b.Sum(), nil
}

func JSONStableOrFallback(value any) string {
	sum, err := JSONStable(value)
	if err == nil {
		return sum
	}
	return StringsOrdered([]string{fmt.Sprintf("%#v", value)})
}
