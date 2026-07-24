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

[English](README.md) | 中文

随 [Apache dubbo-go-pixiu](https://github.com/apache/dubbo-go-pixiu) 一起维护的 AI Agent Skills——让编码助手真正会配 Pixiu 网关。

## 它能给你什么

四个 skill 覆盖 Pixiu 最常见的网关任务：HTTP→Dubbo 路由、自定义 filter、LLM 网关、MCP 网关。它们按真实字段和 SPI 约定生成配置与代码，避开"能编译却没生效"这类坑。

## 安装

不同 AI 工具方式不同：Claude Code / Cursor 内置市场，Codex / OpenCode 需手动安装。

### Claude Code / Cursor

```bash
/plugin marketplace add apache/dubbo-go-pixiu
/plugin install dubbo-go-pixiu@dubbo-go-pixiu-agent-skills
```

### Codex

```bash
codex plugin marketplace add apache/dubbo-go-pixiu
```

然后打开 `/plugins`，选择 `dubbo-go-pixiu` 市场 → Install plugin。

## Skills

### pixiu-filter-author

创建并调试 Pixiu HTTP/Network 过滤器：SPI 实现、`init()` 注册（`RegisterHttpFilter` / `RegisterNetworkFilterPlugin`）、`registry.go` blank import、Kind 对齐 `key.go`、挂载到 HCM 或 filter-chain。

触发词示例：

- "给 Pixiu 写一个自定义 HTTP filter"
- "我的 filter 编译过了但没运行"
- "把这个 Envoy filter 迁移到 Pixiu"

### pixiu-http-to-dubbo

创建并调试 HTTP-to-Dubbo 路由：`api_config.yaml` / `conf.yaml` 路由、注册中心/直连切换、`mappingParams`、`mapType`、`parameterTypes`、`opt.values` / `opt.types`、POJO 处理。

触发词示例：

- "把 /api/v1/users 路由到 Dubbo 的 UserService"
- "consumer 报 no provider found 怎么排查"
- "Dubbo 方法的参数该怎么映射"

### pixiu-llm-gateway

创建并校验 LLM 网关 `conf.yaml`，把 Pixiu 配成多提供方 HTTP 代理：proxy / tokenizer / kvcache 过滤器、`llm_meta`、重试/兜底、vLLM/LMCache、Nacos LLM 服务发现。

触发词示例：

- "配一个用 vLLM 的 LLM 网关"
- "给 LLM 网关加重试和兜底"
- "用 Nacos 做 LLM 服务发现"

### pixiu-mcp-gateway

创建并校验 MCP 网关 `conf.yaml`，把后端 HTTP API 暴露成 MCP 工具：静态/Nacos 动态模式、MCP server 过滤器、tools/resources/prompts、Mcp-Session-Id、OAuth/JWT 保护、Streamable HTTP/SSE。

触发词示例：

- "把我的后端 HTTP API 暴露成 MCP 工具"
- "配一个带 OAuth 保护的 MCP 网关"
- "用 Nacos 动态管理 MCP 工具"

## 贡献

Skill 文件位于 `.agents/skills/<name>/SKILL.md`。

1. 直接修改或新增 skill。
2. 保持 frontmatter 的 `description` 聚焦触发场景——Agent 就是靠这一行决定要不要使用该 skill。
3. 只有当安装路径、skill 名称或支持的 Agent 发生变化时，才修改 `.agents/` 下的元数据。
4. 按 dubbo-go-pixiu 正常贡献流程提交修改。

## License

Apache License 2.0
