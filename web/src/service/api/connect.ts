import { request } from '../request'

// 获取Connect列表
export function fetchGetConnects() {
  return request<App.Api.Connect.Connected[]>({
    url: '/connects',
    method: 'GET',
  })
}

// 获取Connect详情 (直接根据URL获取，不需要request的url)
export function fetchGetConnect(connectUrl: string) {
  return request<App.Api.Connect.Connect>({
    dirrectUrl: `${connectUrl}/api/connect`,
    url: '/',
    method: 'GET',
  })
}

// 添加Connect
export function fetchAddConnect(connectUrl: string) {
  return request<App.Api.Connect.Connected>({
    url: '/addConnect',
    method: 'POST',
    data: {
      connect_url: connectUrl,
    },
  })
}

// 删除Connect
export function fetchDeleteConnect(id: number) {
  return request<App.Api.Connect.Connected>({
    url: `/delConnect/${id}`,
    method: 'DELETE',
  })
}
