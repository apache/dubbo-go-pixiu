# Pixiu-Admin Deployment and Configuration Guide

**English** | [中文](README_CN.md)

**Pixiu-Admin** is a management platform based on the **Pixiu** ecosystem, primarily used for configuring, monitoring, and managing Pixiu gateway resources. It provides centralized management functionality via a modern web user interface and RESTful API.

For backend API documentation, please refer to [API.md](API.md).

## Architecture

### Technology Stack

| Component | Technology |
|-----------|------------|
| Frontend | React 18 + TypeScript + Vite + Ant Design 5 |
| Backend | Go + Gin + GORM + Google Wire |
| Config Store | etcd |
| User Store | MySQL |
| Authentication | JWT |

### Backend Structure

```
admin/
├── app/                    # Application entry point
├── internal/               # Internal packages
│   ├── handler/           # HTTP handlers (Gin)
│   ├── model/             # Domain models
│   ├── store/             # Data access layer (etcd + MySQL)
│   ├── server/            # HTTP server setup
│   ├── wire/              # Dependency injection
│   └── xds/               # xDS server integration
└── pkg/                   # Shared packages
    ├── config/            # Configuration utilities
    └── i18n/              # Internationalization
```

### Frontend Structure

```
admin/web/
├── src/
│   ├── api/               # API client (axios)
│   ├── components/        # Reusable components
│   │   ├── ClusterWizard/    # Cluster configuration wizard
│   │   ├── ListenerWizard/   # Listener configuration wizard
│   │   ├── MappingWizard/    # API mapping wizard
│   │   ├── PluginWizard/     # Plugin group wizard
│   │   ├── DualModeEditor/   # Form + YAML dual mode editor
│   │   └── YamlEditor/       # Monaco-based YAML editor
│   ├── pages/             # Page components
│   ├── layouts/           # Layout components
│   ├── stores/            # State management (Zustand)
│   ├── locales/           # i18n translations (en/zh)
│   └── hooks/             # Custom React hooks
└── vite.config.ts         # Vite configuration
```

## Features

- **Gateway Configuration Management**
  - Cluster management with load balancing and health checks
  - Listener management with multi-protocol support (HTTP/HTTPS/Triple/Dubbo/gRPC)
  - API mapping configuration
  - Plugin group management

- **User Management (RBAC)**
  - User authentication with JWT
  - Role-based access control
  - Permission management

- **Modern UI**
  - Responsive design with Ant Design
  - Dark/Light theme support
  - Internationalization (English/Chinese)
  - Wizard-based configuration with form validation
  - Dual-mode editing (Form + YAML)

## Deployment

### Prerequisites

- Go 1.21+
- Node.js 18+ with pnpm
- etcd 3.5+
- MySQL 8.0+

### Deploy etcd

```bash
docker run -d -p 2379:2379 --env ALLOW_NONE_AUTHENTICATION=yes --name etcd bitnami/etcd
```

For Apple Silicon (M1/M2/M3):

```bash
docker run -d -p 2379:2379 --platform linux/amd64 --env ALLOW_NONE_AUTHENTICATION=yes --name etcd bitnami/etcd:3.5.1
```

### Deploy MySQL

```bash
docker run -d -p 3306:3306 --name mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=pixiu_admin \
  mysql:8.0
```

### Configuration

Edit `configs/admin_config.yaml`:

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

### Run Backend

```bash
# From project root
go run ./cmd/admin/admin.go -c ./configs/admin_config.yaml
```

### Run Frontend (Development)

```bash
cd admin/web
pnpm install
pnpm dev
```

### Build Frontend (Production)

```bash
cd admin/web
pnpm build
```

The built files will be in `admin/web/dist/`, which will be served by the backend.

### Access Admin UI

Open browser and navigate to: `http://127.0.0.1:8081`

Default credentials:
- Username: `admin`
- Password: `admin123`

## API Endpoints

All API endpoints are prefixed with `/api/`:

| Endpoint | Description |
|----------|-------------|
| `POST /api/auth/login` | User login |
| `POST /api/auth/register` | User registration |
| `GET /api/clusters` | List clusters |
| `GET /api/listeners` | List listeners |
| `GET /api/resources` | List API mappings |
| `GET /api/plugins` | List plugin groups |
| `GET /api/users` | List users (admin only) |
| `GET /api/roles` | List roles (admin only) |

For detailed API documentation, see [API.md](API.md).

## Run with Pixiu Gateway

Configure Pixiu to use etcd for dynamic configuration:

```yaml
api_meta_config:
  address: "127.0.0.1:2379"
  api_config_path: "/pixiu/config/api"
```

Start Pixiu:

```bash
go run ./cmd/pixiu/pixiu.go gateway start -c ./configs/pixiu_with_admin_config.yaml
```

## License

This project is licensed under the Apache License 2.0.
