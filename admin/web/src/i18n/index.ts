import { zhCN } from './zh-CN'
import { enUS } from './en-US'

export type Locale = 'zh-CN' | 'en-US'
export const messages = { 'zh-CN': zhCN, 'en-US': enUS }
export const translate = (locale: Locale, key: keyof typeof zhCN) => messages[locale][key]

const textMessages: Record<Locale, Record<string, string>> = {
  'zh-CN': {},
  'en-US': {
    集群: 'Clusters',
    监听器: 'Listeners',
    插件组: 'Plugin groups',
    'API 路由': 'API routes',
    已发布: 'Published',
    草稿: 'Draft',
    已暂停: 'Paused',
    '由 Pixiu Admin 校验': 'Validated by Pixiu Admin',
    取消: 'Cancel',
    '保存中…': 'Saving…',
    保存: 'Save',
    重试: 'Retry',
    未配置端点: 'Endpoint not configured',
    未配置监听地址: 'Listener address not configured',
    条路由: ' routes',
    个插件: ' plugins',
    暂无插件: 'No plugins',
    '管理 Dubbo 集群、注册中心和服务端点。':
      'Manage Dubbo clusters, registries, and service endpoints.',
    '管理监听地址、端口、协议和路由匹配。':
      'Manage listener addresses, ports, protocols, and route matching.',
    '管理网关过滤器插件组及其执行顺序。':
      'Manage gateway filter plugin groups and execution order.',
    加载失败: 'Failed to load',
    加载详情失败: 'Failed to load details',
    保存失败: 'Failed to save',
    删除失败: 'Failed to delete',
    全部: 'All',
    已连接: 'Connected',
    待配置: 'Draft',
    刷新: 'Refresh',
    新建: 'Create',
    编辑: 'Edit',
    删除: 'Delete',
    '配置将通过 Pixiu Admin API 保存。': 'Configuration will be saved through the Pixiu Admin API.',
    尚未配置: 'Not configured',
    已连接后端配置: 'Connected backend configuration',
    类型: 'Type',
    协议: 'Protocol',
    服务端点: 'Service endpoint',
    路由: 'Routes',
    插件: 'Plugins',
    执行配置: 'Execution config',
    操作: 'Actions',
    请求方法: 'Request method',
    后端目标: 'Backend target',
    没有匹配的: 'No matching ',
    没有匹配的路由: 'No matching routes',
    暂无方法: 'No methods',
    '编辑 YAML': 'Edit YAML',
    新建方法: 'Create method',
    暂无: 'No ',
    '点击右上角新建配置。': 'Use the button above to create a configuration.',
    '真实资源列表为空，点击右上角创建第一条路由。':
      'No live resources yet. Use the button above to create the first route.',
    未配置后端目标: 'Backend target not configured',
    查看路由操作: 'View route actions',
    加载路由失败: 'Failed to load routes',
    保存路由失败: 'Failed to save route',
    加载路由详情失败: 'Failed to load route details',
    删除路由失败: 'Failed to delete route',
    加载方法失败: 'Failed to load methods',
    加载方法详情失败: 'Failed to load method details',
    保存方法失败: 'Failed to save method',
  },
}

export function translateText(locale: Locale, value: string) {
  return textMessages[locale][value] ?? value
}
