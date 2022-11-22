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

package diag

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

import (
	"istio.io/api/analysis/v1alpha1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/url"
)

// MessageType is a type of diagnostic message
type MessageType struct {
	// The level of the message.
	level Level

	// The error code of the message
	code string

	// TODO: Make this localizable
	template string
}

// Level returns the level of the MessageType
func (m *MessageType) Level() Level { return m.level }

// Code returns the code of the MessageType
func (m *MessageType) Code() string { return m.code }

// Template returns the message template used by the MessageType
func (m *MessageType) Template() string { return m.template }

// Message is a specific diagnostic message
// TODO: Implement using Analysis message API
type Message struct {
	Type *MessageType

	// The Parameters to the message
	Parameters []interface{}

	// Resource is the underlying resource instance associated with the
	// message, or nil if no resource is associated with it.
	Resource *resource.Instance

	// DocRef is an optional reference tracker for the documentation URL
	DocRef string

	// Line is the line number of the error place in the message
	Line int
}

// Unstructured returns this message as a JSON-style unstructured map
func (m *Message) Unstructured(includeOrigin bool) map[string]interface{} {
	result := make(map[string]interface{})

	result["code"] = m.Type.Code()
	result["level"] = m.Type.Level().String()
	if includeOrigin && m.Resource != nil {
		result["origin"] = m.Resource.Origin.FriendlyName()
		if m.Resource.Origin.Reference() != nil {
			loc := m.Resource.Origin.Reference().String()
			if m.Line != 0 {
				loc = m.ReplaceLine(loc)
			}
			result["reference"] = loc
		}
	}
	result["message"] = fmt.Sprintf(m.Type.Template(), m.Parameters...)

	docQueryString := ""
	if m.DocRef != "" {
		docQueryString = fmt.Sprintf("?ref=%s", m.DocRef)
	}
	result["documentationUrl"] = fmt.Sprintf("%s/%s/%s", url.ConfigAnalysis, strings.ToLower(m.Type.Code()), docQueryString)

	return result
}

func (m *Message) AnalysisMessageBase() *v1alpha1.AnalysisMessageBase {
	docQueryString := ""
	if m.DocRef != "" {
		docQueryString = fmt.Sprintf("?ref=%s", m.DocRef)
	}
	docURL := fmt.Sprintf("%s/%s/%s", url.ConfigAnalysis, strings.ToLower(m.Type.Code()), docQueryString)

	return &v1alpha1.AnalysisMessageBase{
		DocumentationUrl: docURL,
		Level:            v1alpha1.AnalysisMessageBase_Level(v1alpha1.AnalysisMessageBase_Level_value[strings.ToUpper(m.Type.Level().String())]),
		Type: &v1alpha1.AnalysisMessageBase_Type{
			Code: m.Type.Code(),
		},
	}
}

// UnstructuredAnalysisMessageBase returns this message as a JSON-style unstructured map in AnalaysisMessageBase
// TODO(jasonwzm): Remove once message implements AnalysisMessageBase
func (m *Message) UnstructuredAnalysisMessageBase() map[string]interface{} {
	mb := m.AnalysisMessageBase()

	var r map[string]interface{}

	j, err := json.Marshal(mb)
	if err != nil {
		return r
	}
	json.Unmarshal(j, &r) // nolint: errcheck

	return r
}

// Origin returns the origin of the message
func (m *Message) Origin() string {
	origin := ""
	if m.Resource != nil {
		loc := ""
		if m.Resource.Origin.Reference() != nil {
			loc = " " + m.Resource.Origin.Reference().String()
			if m.Line != 0 {
				loc = m.ReplaceLine(loc)
			}
		}
		origin = " (" + m.Resource.Origin.FriendlyName() + loc + ")"
	}
	return origin
}

// String implements io.Stringer
func (m *Message) String() string {
	return fmt.Sprintf("%v [%v]%s %s",
		m.Type.Level(), m.Type.Code(), m.Origin(),
		fmt.Sprintf(m.Type.Template(), m.Parameters...))
}

// MarshalJSON satisfies the Marshaler interface
func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Unstructured(true))
}

// NewMessageType returns a new MessageType instance.
func NewMessageType(level Level, code, template string) *MessageType {
	return &MessageType{
		level:    level,
		code:     code,
		template: template,
	}
}

// NewMessage returns a new Message instance from an existing type.
func NewMessage(mt *MessageType, r *resource.Instance, p ...interface{}) Message {
	return Message{
		Type:       mt,
		Resource:   r,
		Parameters: p,
	}
}

// ReplaceLine replaces the line number from the input String method of Reference to the line number from Message
func (m Message) ReplaceLine(l string) string {
	colonSep := strings.Split(l, ":")
	if len(colonSep) < 2 {
		return l
	}
	_, err := strconv.Atoi(strings.TrimSpace(colonSep[len(colonSep)-1]))
	if err == nil {
		colonSep[len(colonSep)-1] = fmt.Sprintf("%d", m.Line)
	}
	return strings.Join(colonSep, ":")
}
