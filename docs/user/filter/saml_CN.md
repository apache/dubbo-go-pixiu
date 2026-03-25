# SAML 认证过滤器 (dgp.filter.http.auth.saml)

[English](saml.md) | 中文

---

## 概述

`dgp.filter.http.auth.saml` 过滤器允许 Pixiu 作为 SAML Service Provider（SP，服务提供者）运行。

启用后，Pixiu 可以：

- 暴露 SP metadata 端点，供 IdP 注册使用
- 将未认证的浏览器请求重定向到 IdP 登录页
- 在 ACS 端点接收并处理 SAML 响应
- 在用户登录成功后创建 session cookie
- 将选定的 SAML assertion 属性转发为上游 HTTP 请求头

典型的 IdP 包括 Microsoft Entra ID、Okta 和 Keycloak。

## 请求流程

1. 浏览器访问受保护的 Pixiu 路径，例如 `/app`。
2. Pixiu 检查当前请求是否携带有效的 SAML session cookie。
3. 如果 session 不存在，Pixiu 会将浏览器重定向到 IdP 登录页。
4. 用户完成登录后，IdP 会向 Pixiu 的 ACS 端点 POST 一个 `SAMLResponse`。
5. Pixiu 校验 assertion，创建 session cookie，并将浏览器重定向回最初访问的路径。
6. 之后的请求会继续转发到上游服务，同时将配置好的 SAML 属性映射到 HTTP 头中。

## 最小配置

```yaml
http_filters:
  - name: "dgp.filter.http.auth.saml"
    config:
      entity_id: "pixiu-saml-sp"
      acs_url: "https://pixiu.example.com/saml/acs"
      metadata_url: "https://pixiu.example.com/saml/metadata"
      idp_metadata_url: "https://idp.example.com/app/metadata"
      cert_file: "/etc/pixiu/saml/sp.crt"
      key_file: "/etc/pixiu/saml/sp.key"
      rules:
        - match:
            prefix: "/app"
      forward_attributes:
        - saml_attribute: "email"
          header: "X-User-Email"
        - saml_attribute: "displayName"
          header: "X-User-Name"
```

如果不想通过 `idp_metadata_url` 拉取，也可以从本地文件加载 IdP metadata：

```yaml
idp_metadata_file: "/etc/pixiu/saml/idp-metadata.xml"
```

## 配置字段说明

- `entity_id`：Pixiu 作为 SP 对外声明的 entity ID
- `acs_url`：Assertion Consumer Service 端点地址，用于接收 `SAMLResponse`
- `metadata_url`：Pixiu 的 SP metadata 地址，提供给 IdP 管理员导入配置
- `idp_metadata_url`：Pixiu 拉取 IdP metadata 的 URL
- `idp_metadata_file`：本地 IdP metadata 文件路径，可替代 `idp_metadata_url`
- `cert_file`：SP 证书文件
- `key_file`：SP 私钥文件
- `rules`：需要启用 SAML 登录保护的路径前缀
- `forward_attributes`：将 SAML assertion 属性映射到上游 HTTP 头
- `allow_idp_initiated`：面向本地 HTTP 开发/测试场景的兼容开关。当浏览器在跨站 ACS POST 时丢失请求跟踪 cookie，可通过该开关跳过对 `InResponseTo` 的严格依赖
- `err_msg`：SAML 认证失败时返回的本地错误消息

## 重要说明

- `acs_url` 和 `metadata_url` 必须使用相同的 scheme 和 host。
- `metadata_url` 表示 SP metadata 地址，供 IdP 导入配置。
- `idp_metadata_url` 或 `idp_metadata_file` 用于让 Pixiu 获取 IdP 的签名证书和 SSO 端点。
- Pixiu 会先移除 `forward_attributes` 中配置的客户端自带请求头，再写入来自 SAML assertion 的值，以避免 header spoofing。
- 强烈建议在生产环境使用 HTTPS。

## Cookie 与 HTTPS 行为

SAML 的 SP-initiated 登录流程通常依赖一个“请求跟踪 cookie”，以便在 IdP 向 ACS 发起跨站 POST 时，Pixiu 仍然能将返回的 `SAMLResponse` 与先前发起的认证请求对应起来。

- 在 HTTPS 场景下，Pixiu 会使用 `SameSite=None`，这样浏览器才会在跨站 ACS POST 时携带该 cookie。
- 在纯 HTTP 开发环境下，浏览器通常会拒绝这一组合要求，因此 Pixiu 会退回到浏览器默认的 cookie 策略。
- 因此，在本地 HTTP 开发环境中，你可能需要设置 `allow_idp_initiated: true`，使 ACS 流程在不依赖 `InResponseTo` 校验的情况下也能继续完成。

在生产环境中，建议使用 HTTPS，并保持 `allow_idp_initiated` 为关闭状态，除非你确实需要支持 IdP-initiated 登录。
