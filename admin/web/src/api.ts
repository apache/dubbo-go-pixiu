/** Pixiu Admin API contract, aligned with admin/web controllers and API.md. */
export const pixiuAdminApi = {
  base: '/config/api/base',
  resources: { list: '/config/api/resource/list', detail: '/config/api/resource/detail', create: '/config/api/resource', update: '/config/api/resource', remove: '/config/api/resource' },
  methods: { list: '/config/api/resource/method/list', detail: '/config/api/resource/method/detail', create: '/config/api/resource/method', update: '/config/api/resource/method', remove: '/config/api/resource/method' },
  clusters: { list: '/config/api/cluster/list', detail: '/config/api/cluster/detail', create: '/config/api/cluster', update: '/config/api/cluster', remove: '/config/api/cluster' },
  listeners: { list: '/config/api/listener/list', detail: '/config/api/listener/detail', create: '/config/api/listener', update: '/config/api/listener', remove: '/config/api/listener' },
  pluginGroups: { list: '/config/api/plugin_group/list', detail: '/config/api/plugin_group/detail', create: '/config/api/plugin_group', update: '/config/api/plugin_group', remove: '/config/api/plugin_group' },
  rateLimit: { detail: '/config/api/plugin/ratelimit', create: '/config/api/plugin/ratelimit/', update: '/config/api/plugin/ratelimit/', remove: '/config/api/plugin/ratelimit/' },
  opa: { policy: '/config/api/opa/policy' }
} as const

export const mockSurfaces = [] as const
