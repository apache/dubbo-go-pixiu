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

package kvcache

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestSelectPreferredInstanceID_NilResponse(t *testing.T) {
	result := selectPreferredInstanceID(nil)
	assert.Equal(t, "", result)
}

func TestSelectPreferredInstanceID_EmptyLayoutInfo(t *testing.T) {
	resp := &LookupResponse{
		LayoutInfo: map[string]CacheLayout{},
	}

	result := selectPreferredInstanceID(resp)
	assert.Equal(t, "", result)
}

func TestSelectPreferredInstanceID_SingleInstance(t *testing.T) {
	resp := &LookupResponse{
		LayoutInfo: map[string]CacheLayout{
			"instance1": {TokenCount: 100},
		},
	}

	result := selectPreferredInstanceID(resp)
	assert.Equal(t, "instance1", result)
}

func TestSelectPreferredInstanceID_MultipleInstances(t *testing.T) {
	resp := &LookupResponse{
		LayoutInfo: map[string]CacheLayout{
			"instance1": {TokenCount: 50},
			"instance2": {TokenCount: 150},
			"instance3": {TokenCount: 100},
		},
	}

	result := selectPreferredInstanceID(resp)
	assert.Equal(t, "instance2", result)
}

func TestSelectPreferredInstanceID_TieBreaker(t *testing.T) {
	resp := &LookupResponse{
		LayoutInfo: map[string]CacheLayout{
			"instance1": {TokenCount: 100},
			"instance2": {TokenCount: 100},
		},
	}

	result := selectPreferredInstanceID(resp)
	// Should return one of them (map iteration is non-deterministic)
	assert.Contains(t, []string{"instance1", "instance2"}, result)
}

func TestSelectPreferredInstanceID_ZeroTokenCount(t *testing.T) {
	resp := &LookupResponse{
		LayoutInfo: map[string]CacheLayout{
			"instance1": {TokenCount: 0},
			"instance2": {TokenCount: 50},
		},
	}

	result := selectPreferredInstanceID(resp)
	assert.Equal(t, "instance2", result)
}

func TestSelectPreferredInstanceID_AllZeroTokens(t *testing.T) {
	resp := &LookupResponse{
		LayoutInfo: map[string]CacheLayout{
			"instance1": {TokenCount: 0},
			"instance2": {TokenCount: 0},
		},
	}

	result := selectPreferredInstanceID(resp)
	// Should still return one of them
	assert.Contains(t, []string{"instance1", "instance2"}, result)
}

func TestExtractPromptFromString(t *testing.T) {
	// extractPromptFromMessages expects []any, not a string
	// This test should use coercePrompt instead, but that's not exported
	// So we test with the actual expected input format
	result := extractPromptFromMessages("simple string prompt")
	assert.Equal(t, "", result) // Returns empty for non-array input
}

func TestExtractPromptFromMessages_SingleMessage(t *testing.T) {
	messages := []any{
		map[string]any{
			"role":    "user",
			"content": "Hello, world!",
		},
	}

	result := extractPromptFromMessages(messages)
	assert.Equal(t, "Hello, world!", result) // Only extracts content, not role
}

func TestExtractPromptFromMessages_MultipleMessages(t *testing.T) {
	messages := []any{
		map[string]any{
			"role":    "system",
			"content": "You are a helpful assistant.",
		},
		map[string]any{
			"role":    "user",
			"content": "What is the weather?",
		},
		map[string]any{
			"role":    "assistant",
			"content": "I don't have access to weather data.",
		},
	}

	result := extractPromptFromMessages(messages)
	expected := "You are a helpful assistant.\nWhat is the weather?\nI don't have access to weather data."
	assert.Equal(t, expected, result)
}

func TestExtractPromptFromMessages_MissingRole(t *testing.T) {
	messages := []any{
		map[string]any{
			"content": "Message without role",
		},
	}

	result := extractPromptFromMessages(messages)
	assert.Equal(t, "Message without role", result) // Role is not required
}

func TestExtractPromptFromMessages_MissingContent(t *testing.T) {
	messages := []any{
		map[string]any{
			"role": "user",
		},
	}

	result := extractPromptFromMessages(messages)
	assert.Equal(t, "", result) // No content means empty result
}

func TestExtractPromptFromMessages_NonStringContent(t *testing.T) {
	messages := []any{
		map[string]any{
			"role":    "user",
			"content": 12345,
		},
	}

	result := extractPromptFromMessages(messages)
	assert.Equal(t, "", result) // Non-string content is skipped
}

func TestExtractPromptFromMessages_EmptyArray(t *testing.T) {
	messages := []map[string]any{}

	result := extractPromptFromMessages(messages)
	assert.Equal(t, "", result)
}

func TestExtractPromptFromMessages_InvalidType(t *testing.T) {
	result := extractPromptFromMessages(12345)
	assert.Equal(t, "", result)
}

func TestExtractPromptFromMessages_NilValue(t *testing.T) {
	result := extractPromptFromMessages(nil)
	assert.Equal(t, "", result)
}
