import { useQuery } from '@tanstack/react-query'
import { listenerApi } from '../../services/listener-api'
import { asJsonObject } from '../../types/api'
export function ListenerPage() {
  const q = useQuery({ queryKey: ['listeners'], queryFn: listenerApi.list })
  return (
    <div className="page-resource">
      <div className="placeholder-head">
        <div>
          <p className="eyebrow">CONNECTED API</p>
          <h1>监听器</h1>
          <p className="muted">管理监听地址、端口、协议和路由匹配。</p>
        </div>
        <button className="primary">+ 新建监听器</button>
      </div>
      <div className="panel route-panel">
        {q.isLoading ? (
          <div className="loading-skeleton" />
        ) : q.isError ? (
          <div className="empty">
            加载失败
            <button className="secondary" onClick={() => q.refetch()}>
              重试
            </button>
          </div>
        ) : !q.data?.length ? (
          <div className="empty">
            暂无监听器<button className="primary">新建监听器</button>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>名称</th>
                <th>协议</th>
                <th>地址</th>
                <th>Cluster</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {q.data.map((x, i) => (
                <tr key={String(x.name || i)}>
                  <td>
                    <b>{String(x.name || '未命名监听器')}</b>
                  </td>
                  <td>{String(x.protocol || x.protocolStr || 'HTTP')}</td>
                  <td>{String(x.address || x.port || '-')}</td>
                  <td>{String(x.cluster || asJsonObject(x.route).cluster || '-')}</td>
                  <td>
                    <button className="link-btn">编辑</button>
                    <button className="link-btn">删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
