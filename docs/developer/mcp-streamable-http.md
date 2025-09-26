# MCP Streamable HTTP 传输增强设计（待删除）

## 背景

目前 Pixiu 的 MCP (Model Context Protocol) 服务器实现仅支持传统的 HTTP 请求-响应模式。随着 MCP 协议的发展，新的 **Streamable HTTP** 传输方式成为标准，它在单一端点上同时支持 POST 和 GET 方法，提供更灵活的通信模式。

## 目标

- 增强现有 MCP 服务器支持 Streamable HTTP 传输
- 保持向后兼容，不破坏现有 HTTP POST 功能
- 无需配置改动，直接增强现有实现
- 支持 Server-Sent Events (SSE) 流式响应

## 现状分析

### 当前实现

```text
客户端 → HTTP POST /mcp → MCPServerFilter → 同步 JSON 响应
```

- 仅支持 POST 方法
- 同步请求-响应模式
- JSON 格式响应

### Streamable HTTP 要求

```text
客户端 → POST /mcp (发送消息) → MCPServerFilter → JSON/SSE 响应
客户端 → GET /mcp (建立流) → MCPServerFilter → SSE 流
```

- 同一端点支持 POST 和 GET
- 基于 Accept 头协商响应格式
- 支持会话管理和流式推送

## 设计思路

### 核心原则

1. **最小改动**: 在现有 MCPServerFilter 基础上增强
2. **向后兼容**: 保持现有配置和 API 不变
3. **协议原生**: 直接支持 Streamable HTTP，无需配置开关
4. **职责分离**: 通过内部模块化保持代码清晰

### 架构增强

```text
原有架构: MCPServerFilter (仅 HTTP POST)
     ↓
增强架构: MCPServerFilter
          ├── HTTP POST 处理 (原有逻辑)
          ├── HTTP GET 处理 (新增 SSE)
          ├── 内容协商 (JSON vs SSE)
          └── 会话管理 (SSE 连接)
```

## 技术方案

### 1. 请求分发增强

在现有 `Decode` 方法中根据 HTTP 方法分发：

- **POST 请求**: 消息处理，支持同步和流式响应
- **GET 请求**: 建立 SSE 流，维持长连接

### 2. 内容协商

基于 `Accept` 头决定响应模式：

- `application/json`: 传统同步响应
- `text/event-stream`: SSE 流式响应
- 工具调用优先使用流式响应

### 3. 会话管理

- 支持 `Mcp-Session-Id` 头部
- 内存存储会话和连接映射
- 自动连接清理和心跳机制

### 4. SSE 流处理

- 标准 SSE 格式: `data: {json}\n\n`
- 连接生命周期管理
- 错误处理和降级机制

## 实现要点

### 代码组织

```text
pkg/filter/mcp/mcpserver/
├── filter.go           # 主 Filter，增强请求分发
├── handlers.go         # MCP 协议处理 (保持不变)
├── registry.go         # 工具注册 (保持不变)
└── transport/          # 新增: 传输层模块
    ├── streamable_http.go    # Streamable HTTP 协议处理
    ├── sse_manager.go        # SSE 连接管理
    ├── session_manager.go    # 会话管理
    └── content_negotiator.go # 内容协商
```

### 关键增强点

1. **方法路由**: 同一端点的 POST/GET 分发
2. **响应协商**: 根据 Accept 头选择响应格式  
3. **流式处理**: SSE 连接建立和维护
4. **会话跟踪**: 会话 ID 生成和映射
5. **错误处理**: 优雅降级和错误响应

### 兼容性保证

- 现有 HTTP POST 行为完全不变
- 配置文件无需修改
- API 接口保持一致
- 现有客户端无感知

## 优势

1. **标准合规**: 完全符合 MCP Streamable HTTP 规范
2. **零配置**: 无需用户配置，自动启用
3. **渐进增强**: 不影响现有功能，纯增强
4. **高性能**: 复用现有架构，最小开销
5. **可扩展**: 为未来协议扩展奠定基础

## 实施路径

1. **第一阶段**: 增强请求分发和内容协商
2. **第二阶段**: 实现 SSE 连接管理
3. **第三阶段**: 集成流式响应处理
4. **第四阶段**: 测试和优化

通过这种设计，Pixiu 的 MCP 服务器将无缝支持最新的 Streamable HTTP 传输，为用户提供更灵活和高效的通信体验。
