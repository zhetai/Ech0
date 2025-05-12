import { request } from '../request'

// 获取当前登录用户信息
export function fetchGetCurrentUser() {
  return request<App.Api.User.User>({
    url: '/user',
    method: 'GET',
  })
}

// 更新用户信息
export function fetchUpdateUser(user: App.Api.User.UserInfo) {
  return request({
    url: '/user',
    method: 'PUT',
    data: user,
  })
}

