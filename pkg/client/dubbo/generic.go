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
	"encoding/json"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	cst "github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

// supportedGenericModes lists the generic invocation modes accepted in
// integrationRequest.generic. Values outside this set are rejected instead of
// falling back to the map mode, because a silent fallback would change the
// request contract without telling the caller.
var supportedGenericModes = []string{
	cst.GenericModeMap,
	cst.GenericModeMapAlias,
	cst.GenericModeGson,
	cst.GenericModeProtobufJSON,
	cst.GenericModeBean,
}

// NormalizeGenericMode canonicalizes the configured generic invocation mode.
// An empty value selects the map mode, which is what dubbo-go applies when a
// reference does not request a specific mode. The readable alias "map" resolves
// to the "true" value understood by dubbo-go.
func NormalizeGenericMode(mode string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(mode))
	switch normalized {
	case "", cst.GenericModeMap, cst.GenericModeMapAlias:
		return cst.GenericModeMap, nil
	case cst.GenericModeGson, cst.GenericModeProtobufJSON, cst.GenericModeBean:
		return normalized, nil
	default:
		return "", errors.Errorf("generic mode %q is not supported, expected one of %s",
			mode, strings.Join(supportedGenericModes, ", "))
	}
}

// IsJSONTextGenericMode reports whether the mode carries the request as JSON
// text instead of a list of typed arguments. The provider realizes that text
// into its own request type, which is what lets a caller without generated
// code invoke a protobuf defined service.
func IsJSONTextGenericMode(mode string) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case cst.GenericModeGson, cst.GenericModeProtobufJSON:
		return true
	default:
		return false
	}
}

// EncodeJSONTextArgument renders one mapped value as the JSON text that a JSON
// text mode sends. A string value is taken as an already serialized JSON
// document and has to be valid JSON text, because the provider parses the text
// as the request message; a plain string such as a query parameter is rejected
// here instead of failing inside the provider. Values of any other type are
// marshaled, which keeps numbers decoded as json.Number exact.
func EncodeJSONTextArgument(value any) (string, error) {
	if text, ok := value.(string); ok {
		if !json.Valid([]byte(text)) {
			return "", errors.Errorf("generic JSON argument is not valid JSON text: %q", text)
		}
		return text, nil
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "encode generic JSON argument")
	}
	return string(encoded), nil
}

// UnwrapJSONTextResult turns the JSON text a JSON text mode returns into raw
// bytes, so the HTTP layer writes the provider's JSON document verbatim instead
// of JSON encoding it again as a quoted string. JSON text modes serialize every
// result as a JSON document, so a string result already carries its quotes
// (`"hello"`); passing the text through as a string would quote it twice.
// Values that are not valid JSON text stay untouched.
func UnwrapJSONTextResult(value any) any {
	text, ok := value.(string)
	if !ok || !json.Valid([]byte(text)) {
		return value
	}
	return []byte(text)
}
