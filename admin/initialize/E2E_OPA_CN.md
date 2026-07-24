# PR2：Admin OPA → Gateway OPA 全链路 E2E 测试

[English](E2E_OPA.md) | **中文**

## 测试套件验证的内容

该测试框架启动了：

1. 一个**进程内的智能 OPA mock**（`regoMockOPA`）—— 使用真实的
   `github.com/open-policy-agent/opa/rego` 库来编译和评估模块，所以测试中的
   策略语义和真实 OPA 守护进程的行为是一致的。无需 docker、etcd 或真实的
   OPA 二进制。
2. **真实的 admin Gin 路由**，通过 `initialize.Routers()` 启动，并将
   `adminconfig.Bootstrap.OPA.ServerURL` 指向该 mock。
3. **真实的 gateway OPA filter**，来自 `pkg/filter/opa`（通过对外暴露的
   `Plugin.CreateFilterFactory()` API），同样指向该 mock URL。

每个测试用例先通过 admin REST PUT 下发一条策略，然后通过 gateway filter
驱动一个或多个 HTTP 请求，并对决策结果做断言。

| 测试用例 | 场景 | 验证内容 |
|---|---|---|
| `TestE2E_AllowedThroughFullChain` | PUT "allow if GET" → GET 请求 | Admin → OPA → gateway 端到端 allow 返回 `filter.Continue`，无本地应答 |
| `TestE2E_DeniedThroughFullChain` | 同策略，POST 请求 | Deny 返回 `filter.Stop` + 403，在到达上游之前短路 |
| `TestE2E_DefaultDenyForAllRequests` | 仅 `default allow := false` | 5 种 HTTP 方法全部被拒绝（任何形状的规则都无法绕过） |
| `TestE2E_PolicyHotReload` | PUT v1，再 PUT v2（无需重启） | 行为在下一个请求上立刻翻转 —— 这是 OPA server 模式最核心的价值 |
| `TestE2E_DeleteCausesMissingResultFailClosed` | 先 PUT 再 DELETE | Gateway 返回 502（与 `test_opa.md` §6.6 一致） |
| `TestE2E_HeaderBasedAllowDeny` | 基于 `input.headers["X-Role"]` 的策略 | 首字母大写的 header 传递正常（admin / user / 缺失 三种场景） |
| `TestE2E_GatewayTimeoutFailClosed` | 决策延迟 200ms，gateway 超时 50ms | 返回 504，耗时被控制在 180ms 以内 |
| `TestE2E_PolicyIDOverrideRoutesThroughGateway` | PUT 时使用 form 字段级别的 `policy_id` 覆写 | 覆写后的 policy ID 会被写入 OPA，并可通过 gateway 决策路径命中 |

## 运行方式

```bash
# 默认 —— 运行全部 PR2 用例，非 verbose
./admin/initialize/run.sh

# 详细输出
VERBOSE=1 ./admin/initialize/run.sh

# 按名称筛选用例
./admin/initialize/run.sh -run AllowedThroughFullChain

# 或直接运行：
go test -count=1 -run TestE2E_ -v ./admin/initialize/
```

预期输出（verbose 模式）：

```
=== RUN   TestE2E_AllowedThroughFullChain
--- PASS: TestE2E_AllowedThroughFullChain (0.03s)
=== RUN   TestE2E_DeniedThroughFullChain
--- PASS: TestE2E_DeniedThroughFullChain (0.01s)
=== RUN   TestE2E_DefaultDenyForAllRequests
--- PASS: TestE2E_DefaultDenyForAllRequests (0.01s)
=== RUN   TestE2E_PolicyHotReload
--- PASS: TestE2E_PolicyHotReload (0.01s)
=== RUN   TestE2E_DeleteCausesMissingResultFailClosed
--- PASS: TestE2E_DeleteCausesMissingResultFailClosed (0.01s)
=== RUN   TestE2E_HeaderBasedAllowDeny
--- PASS: TestE2E_HeaderBasedAllowDeny (0.01s)
=== RUN   TestE2E_GatewayTimeoutFailClosed
--- PASS: TestE2E_GatewayTimeoutFailClosed (0.21s)
=== RUN   TestE2E_PolicyIDOverrideRoutesThroughGateway
--- PASS: TestE2E_PolicyIDOverrideRoutesThroughGateway (0.01s)
PASS
ok      github.com/apache/dubbo-go-pixiu/admin/initialize       0.342s
```
