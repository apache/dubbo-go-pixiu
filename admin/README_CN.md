# Pixiu-Admin 部署与配置文档

[English](README.md) | **中文**

**Pixiu-Admin** 是基于 **Pixiu** 生态的管理平台，主要用于配置、监控和管理 Pixiu 网关资源。通过现代化的 Web 用户界面和 RESTful API 提供集中式管理功能。

后端 API 接口文档请参考 [API_CN.md](API_CN.md)。

## 架构

### 技术栈

| 组件 | 技术 |
|------|------|
| 前端 | React 18 + TypeScript + Vite + Ant Design 5 |
| 后端 | Go + Gin + GORM + Google Wire |
| 配置存储 | etcd |
| 用户存储 | MySQL |
| 认证方式 | JWT |

### 后端结构

```
admin/
├── app/                    # 应用入口
├── internal/               # 内部包
│   ├── handler/           # HTTP 处理器 (Gin)
│   ├── model/             # 领域模型
│   ├── store/             # 数据访问层 (etcd + MySQL)
│   ├── server/            # HTTP 服务器配置
│   ├── wire/              # 依赖注入
│   └── xds/               # xDS 服务器集成
└── pkg/                   # 共享包
    ├── config/            # 配置工具
    └── i18n/              # 国际化
```

### 前端结构

```
admin/web/
├── src/
│   ├── api/               # API 客户端 (axios)
│   ├── components/        # 可复用组件
│   │   ├── ClusterWizard/    # 集群配置向导
│   │   ├── ListenerWizard/   # 监听器配置向导
│   │   ├── MappingWizard/    # API 映射向导
│   │   ├── PluginWizard/     # 插件组向导
│   │   ├── DualModeEditor/   # 表单 + YAML 双模式编辑器
│   │   └── YamlEditor/       # 基于 Monaco 的 YAML 编辑器
│   ├── pages/             # 页面组件
│   ├── layouts/           # 布局组件
│   ├── stores/            # 状态管理 (Zustand)
│   ├── locales/           # 国际化翻译 (中/英)
│   └── hooks/             # 自定义 React Hooks
└── vite.config.ts         # Vite 配置
```

## 功能特性

- **网关配置管理**
  - 集群管理，支持负载均衡和健康检查
  - 监听器管理，支持多协议 (HTTP/HTTPS/Triple/Dubbo/gRPC)
  - API 映射配置
  - 插件组管理

- **用户管理 (RBAC)**
  - 基于 JWT 的用户认证
  - 基于角色的访问控制
  - 权限管理

- **现代化 UI**
  - 基于 Ant Design 的响应式设计
  - 深色/浅色主题切换
  - 国际化支持 (中文/英文)
  - 向导式配置，带表单验证
  - 双模式编辑 (表单 + YAML)

## 部署

### 环境要求

- Go 1.21+
- Node.js 18+ 和 pnpm
- etcd 3.5+
- MySQL 8.0+

### 部署 etcd

```bash
docker run -d -p 2379:2379 --env ALLOW_NONE_AUTHENTICATION=yes --name etcd bitnami/etcd
```

Apple Silicon (M1/M2/M3) 用户：

```bash
docker run -d -p 2379:2379 --platform linux/amd64 --env ALLOW_NONE_AUTHENTICATION=yes --name etcd bitnami/etcd:3.5.1
```

### 部署 MySQL

```bash
docker run -d -p 3306:3306 --name mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=pixiu_admin \
  mysql:8.0
```

### 配置

编辑 `configs/admin_config.yaml`：

```yaml
admin:
  address: 0.0.0.0
  port: 8081
  static_resources_dir: ./admin/web/dist

etcd:
  endpoints:
    - 127.0.0.1:2379
  base_path: /pixiu/config

mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: root
  database: pixiu_admin

jwt:
  secret: your-secret-key
  issuer: pixiu-admin
  expire_hours: 24
```

### 运行后端

```bash
# 在项目根目录
go run ./cmd/admin/admin.go -c ./configs/admin_config.yaml
```

### 运行前端 (开发模式)

```bash
cd admin/web
pnpm install
pnpm dev
```

### 构建前端 (生产模式)

```bash
cd admin/web
pnpm build
```

构建产物位于 `admin/web/dist/`，将由后端服务提供静态文件服务。

### 访问管理界面

打开浏览器访问：`http://127.0.0.1:8081`

默认账号：
- 用户名：`admin`
- 密码：`admin123`

## API 接口

所有 API 接口以 `/api/` 为前缀：

| 接口 | 描述 |
|------|------|
| `POST /api/auth/login` | 用户登录 |
| `POST /api/auth/register` | 用户注册 |
| `GET /api/clusters` | 获取集群列表 |
| `GET /api/listeners` | 获取监听器列表 |
| `GET /api/resources` | 获取 API 映射列表 |
| `GET /api/plugins` | 获取插件组列表 |
| `GET /api/users` | 获取用户列表 (仅管理员) |
| `GET /api/roles` | 获取角色列表 (仅管理员) |

详细 API 文档请参考 [API_CN.md](API_CN.md)。

## 与 Pixiu 网关配合使用

配置 Pixiu 使用 etcd 进行动态配置：

```yaml
api_meta_config:
  address: "127.0.0.1:2379"
  api_config_path: "/pixiu/config/api"
```

启动 Pixiu：

```bash
go run ./cmd/pixiu/pixiu.go gateway start -c ./configs/pixiu_with_admin_config.yaml
```

## 许可证

本项目采用 Apache License 2.0 开源许可。
