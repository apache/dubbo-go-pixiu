<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# dubbo-go-pixiu Agent Skills

English | [中文](README_CN.md)

AI agent skills maintained alongside [Apache dubbo-go-pixiu](https://github.com/apache/dubbo-go-pixiu) — so your coding assistant actually knows how to configure the Pixiu gateway.

## What you get

Four skills cover Pixiu's most common gateway tasks: HTTP-to-Dubbo routing, custom filters, the LLM gateway, and the MCP gateway. They generate config and code against the real fields and SPI conventions, avoiding the "compiles but doesn't take effect" pitfalls.

## Installation

The method differs by AI tool: Claude Code / Cursor have a built-in marketplace; Codex installs manually.

### Claude Code / Cursor

```bash
/plugin marketplace add apache/dubbo-go-pixiu
/plugin install dubbo-go-pixiu@dubbo-go-pixiu-agent-skills
```

### Codex

```bash
codex plugin marketplace add apache/dubbo-go-pixiu
```

Then open `/plugins`, select the `dubbo-go-pixiu` marketplace → Install plugin.

## Skills

### pixiu-filter-author

Create and debug Pixiu HTTP/Network filters: SPI implementation, registration in `init()` (`RegisterHttpFilter` / `RegisterNetworkFilterPlugin`), the blank import in `registry.go`, aligning `Kind` with `key.go`, and mounting under HCM or the filter-chain.

Example triggers:

- "Write a custom HTTP filter for Pixiu"
- "My filter compiles but doesn't run"
- "Port this Envoy filter to Pixiu"

### pixiu-http-to-dubbo

Create and debug HTTP-to-Dubbo routes: `api_config.yaml` / `conf.yaml` routes, registry/direct mode switching, `mappingParams`, `mapType`, `parameterTypes`, `opt.values` / `opt.types`, and POJO handling.

Example triggers:

- "Route /api/v1/users to the Dubbo UserService"
- "How do I debug a consumer reporting no provider found"
- "How should I map the parameters of a Dubbo method"

### pixiu-llm-gateway

Create and validate an LLM gateway `conf.yaml` that configures Pixiu as a multi-provider HTTP proxy: proxy / tokenizer / kvcache filters, `llm_meta`, retry/fallback, vLLM/LMCache, and Nacos LLM discovery.

Example triggers:

- "Set up an LLM gateway using vLLM"
- "Add retry and fallback to the LLM gateway"
- "Use Nacos for LLM service discovery"

### pixiu-mcp-gateway

Create and validate an MCP gateway `conf.yaml` that exposes backend HTTP APIs as MCP tools: static / Nacos dynamic modes, the MCP server filter, tools/resources/prompts, Mcp-Session-Id, OAuth/JWT protection, and Streamable HTTP/SSE.

Example triggers:

- "Expose my backend HTTP API as MCP tools"
- "Set up an MCP gateway with OAuth protection"
- "Manage MCP tools dynamically with Nacos"

## Contributing

Skill files live at `.agents/skills/<name>/SKILL.md`.

1. Edit or add skills directly.
2. Keep each skill's frontmatter `description` focused on its trigger scenarios — the agent relies on that single line to decide whether to use the skill.
3. Only touch the metadata under `.agents/` when the install path, skill names, or supported agents change.
4. Submit changes through the normal dubbo-go-pixiu contribution process.

## License

Apache License 2.0
