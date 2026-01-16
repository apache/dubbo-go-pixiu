# Pixiu Admin Frontend

Modern web interface for Dubbo Go Pixiu gateway management.

## Technology Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| React | 18 | UI Framework |
| TypeScript | 5.x | Type Safety |
| Vite | 7.x | Build Tool |
| Ant Design | 5.x | UI Components |
| Zustand | 5.x | State Management |
| React Router | 7.x | Routing |
| Axios | 1.x | HTTP Client |
| i18next | 25.x | Internationalization |
| Monaco Editor | - | YAML Editor |
| ECharts | 5.x | Charts |

## Project Structure

```
src/
├── api/                    # API client layer
│   ├── request.ts         # Axios instance configuration
│   ├── hooks.ts           # React Query hooks
│   ├── auth.ts            # Authentication APIs
│   ├── cluster.ts         # Cluster APIs
│   ├── listener.ts        # Listener APIs
│   ├── resource.ts        # Resource APIs
│   ├── plugin.ts          # Plugin APIs
│   └── user.ts            # User management APIs
├── components/            # Reusable components
│   ├── ClusterWizard/     # Cluster configuration wizard
│   ├── ListenerWizard/    # Listener configuration wizard
│   ├── MappingWizard/     # API mapping wizard
│   ├── PluginWizard/      # Plugin group wizard
│   ├── DualModeEditor/    # Form + YAML dual mode editor
│   ├── YamlEditor/        # Monaco-based YAML editor
│   ├── WizardModal/       # Wizard modal wrapper
│   ├── AuthGuard.tsx      # Route authentication guard
│   ├── ThemeProvider.tsx  # Dark/Light theme provider
│   └── SafeECharts.tsx    # Safe ECharts wrapper
├── pages/                 # Page components
│   ├── Login.tsx          # Login page
│   ├── Overview.tsx       # Dashboard overview
│   ├── Cluster.tsx        # Cluster management
│   ├── Listener.tsx       # Listener management
│   ├── Mapping.tsx        # API mapping management
│   ├── Plugin.tsx         # Plugin group management
│   ├── UserManagement.tsx # User management (admin)
│   ├── RoleManagement.tsx # Role management (admin)
│   ├── Settings.tsx       # System settings
│   └── Profile.tsx        # User profile
├── layouts/               # Layout components
│   └── MainLayout.tsx     # Main application layout
├── stores/                # State management (Zustand)
│   └── authStore.ts       # Authentication state
├── locales/               # i18n translations
│   ├── en/               # English translations
│   └── zh/               # Chinese translations
├── hooks/                 # Custom React hooks
├── types/                 # TypeScript type definitions
├── config/                # Application configuration
├── App.tsx                # Application root
└── main.tsx               # Entry point
```

## Development

### Prerequisites

- Node.js 18+
- pnpm 8+

### Install Dependencies

```bash
pnpm install
```

### Start Development Server

```bash
pnpm dev
```

The development server will start at `http://localhost:5173` with hot module replacement enabled.

### Build for Production

```bash
pnpm build
```

Build artifacts will be generated in the `dist/` directory.

### Linting

```bash
pnpm lint
```

## Features

### Gateway Configuration

- **Cluster Management**: Configure backend clusters with load balancing strategies and health checks
- **Listener Management**: Configure protocol listeners (HTTP, HTTPS, Triple, Dubbo, gRPC)
- **API Mapping**: Configure API routing and request/response transformations
- **Plugin Groups**: Manage plugin chains for request processing

### User Interface

- **Wizard-based Configuration**: Step-by-step wizards for complex configurations
- **Dual-mode Editing**: Switch between form-based and YAML editing
- **Real-time Validation**: Instant feedback on configuration errors
- **Responsive Design**: Works on desktop and tablet devices

### User Management

- **RBAC**: Role-based access control
- **User Management**: Create, update, delete users
- **Role Management**: Define roles with specific permissions
- **Permission Control**: Fine-grained access control

### Internationalization

- **Languages**: English and Chinese
- **Auto-detection**: Automatically detects browser language
- **Runtime Switching**: Switch languages without page reload

### Theming

- **Dark Mode**: Full dark mode support
- **Light Mode**: Default light theme
- **System Preference**: Follows system preference

## Configuration

### Environment Variables

Create a `.env.local` file for local development:

```env
VITE_API_BASE_URL=http://localhost:8081/api
```

### API Proxy

The Vite development server is configured to proxy API requests to the backend. See `vite.config.ts` for configuration.

## License

Apache License 2.0
