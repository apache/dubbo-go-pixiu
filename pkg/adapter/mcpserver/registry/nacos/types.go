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

package nacos

// NacosTool is the structure for defining a tool in Nacos
type NacosTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type RequestTemplate struct {
	URL            string              `json:"url"`
	Method         string              `json:"method"`
	Headers        []map[string]string `json:"headers"`
	ArgsToJsonBody bool                `json:"argsToJsonBody"`
	ArgsToUrlParam bool                `json:"argsToUrlParam"`
}

type ResponseTemplate struct {
	PrependBody string `json:"prependBody"`
}

type JsonGoTemplate struct {
	RequestTemplate  RequestTemplate  `json:"requestTemplate"`
	ResponseTemplate ResponseTemplate `json:"responseTemplate"`
}

type ToolMeta struct {
	Enabled   bool           `json:"enabled"`
	Templates map[string]any `json:"templates"`
}

type ToolsSpec struct {
	Tools     []NacosTool         `json:"tools"`
	ToolsMeta map[string]ToolMeta `json:"toolsMeta"`
}
