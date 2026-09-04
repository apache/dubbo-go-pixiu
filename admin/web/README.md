# Pixiu Admin Web

Dubbo Go Pixiu 的管理前端，React + TypeScript + Tailwind CSS + Monaco Editor，Vite 构建。

## 开发

```bash
npm install
npm run dev
```

- 开发服务器：`http://localhost:8088`
- Vite 将 `/login`、`/register`、`/user`、`/config`、`/swagger` 代理到 Admin 后端（默认 `http://127.0.0.1:8081`）。可通过 `VITE_BACKEND_URL` 覆盖后端地址，例如 Docker Compose 中的 `http://backend:8081`。
- 后端无 CORS 中间件，跨域开发必须走该代理；`/login` 与 `/register` 的 GET 请求由前端路由接管，POST 才转发后端。

## 构建

```bash
npm run build
```

产物输出到 `dist/`。

## 功能边界

- 页面与交互以本 README 和当前源码为准。
- 统一响应约定：HTTP 200 + `{"code": "10001", "data": ...}` 为成功；配置类请求携带 `token` 与 `username` 头。
- 插件组（`/config/api/plugin_group/*`）与限流（`/config/api/plugin/ratelimit*`）接口在当前后端 router 中未注册，对应页面会显示接口错误，属预期边界。
