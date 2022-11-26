// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tmpl

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/file"
)

// Evaluate parses the template and then executes it with the given parameters.
func Evaluate(tpl string, data interface{}) (string, error) {
	t, err := Parse(tpl)
	if err != nil {
		return "", err
	}

	return Execute(t, data)
}

func EvaluateFile(filePath string, data interface{}) (string, error) {
	tpl, err := file.AsString(filePath)
	if err != nil {
		return "", err
	}
	return Evaluate(tpl, data)
}

// EvaluateOrFail calls Evaluate and fails tests if it returns error.
func EvaluateOrFail(t test.Failer, tpl string, data interface{}) string {
	t.Helper()
	s, err := Evaluate(tpl, data)
	if err != nil {
		t.Fatalf("tmpl.EvaluateOrFail: %v", err)
	}
	return s
}

func EvaluateFileOrFail(t test.Failer, filePath string, data interface{}) string {
	t.Helper()
	s, err := EvaluateFile(filePath, data)
	if err != nil {
		t.Fatalf("tmpl.EvaluateFileOrFail: %v", err)
	}
	return s
}

// MustEvaluate calls Evaluate and panics if there is an error.
func MustEvaluate(tpl string, data interface{}) string {
	s, err := Evaluate(tpl, data)
	if err != nil {
		panic(fmt.Sprintf("tmpl.MustEvaluate: %v", err))
	}
	return s
}

func MustEvaluateFile(filePath string, data interface{}) string {
	s, err := EvaluateFile(filePath, data)
	if err != nil {
		panic(fmt.Sprintf("tmpl.MustEvaluate: %v", err))
	}
	return s
}

// EvaluateAll calls Evaluate the same data args against each of the given templates.
func EvaluateAll(data interface{}, templates ...string) ([]string, error) {
	out := make([]string, 0, len(templates))
	for _, t := range templates {
		content, err := Evaluate(t, data)
		if err != nil {
			return nil, err
		}
		out = append(out, content)
	}
	return out, nil
}

func EvaluateAllFiles(data interface{}, filePaths ...string) ([]string, error) {
	templates, err := file.AsStringArray(filePaths...)
	if err != nil {
		return nil, err
	}
	return EvaluateAll(data, templates...)
}

func MustEvaluateAll(data interface{}, templates ...string) []string {
	out, err := EvaluateAll(data, templates...)
	if err != nil {
		panic(fmt.Sprintf("tmpl.MustEvaluateAll: %v", err))
	}
	return out
}

// EvaluateAllOrFail calls Evaluate and fails t if an error occurs.
func EvaluateAllOrFail(t test.Failer, data interface{}, templates ...string) []string {
	t.Helper()
	out, err := EvaluateAll(data, templates...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func EvaluateAllFilesOrFail(t test.Failer, data interface{}, filePaths ...string) []string {
	t.Helper()
	out, err := EvaluateAllFiles(data, filePaths...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
