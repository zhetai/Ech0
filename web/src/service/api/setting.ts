import { request } from '../request'

// 获取系统设置
export function fetchGetSettings() {
  return request<App.Api.Setting.SystemSetting>({
    url: '/settings',
    method: 'GET',
  })
}

// 获取评论设置
export function fetchGetCommentSettings() {
  return request<App.Api.Setting.CommentSetting>({
    url: '/comment/settings',
    method: 'GET',
  })
}

// 更新系统设置
export function fetchUpdateSettings(systemSetting: App.Api.Setting.SystemSetting) {
  return request({
    url: '/settings',
    method: 'PUT',
    data: systemSetting,
  })
}

// 更新评论设置
export function fetchUpdateCommentSettings(commentSetting: App.Api.Setting.CommentSetting) {
  return request({
    url: '/comment/settings',
    method: 'PUT',
    data: commentSetting,
  })
}

// 获取 S3 存储设置
export function fetchGetS3Settings() {
  return request<App.Api.Setting.S3Setting>({
    url: '/s3/settings',
    method: 'GET',
  })
}

// 更新 S3 存储设置
export function fetchUpdateS3Settings(s3Setting: App.Api.Setting.S3Setting) {
  return request({
    url: '/s3/settings',
    method: 'PUT',
    data: s3Setting,
  })
}

// 获取 OAuth2 设置
export function fetchGetOAuth2Settings() {
  return request<App.Api.Setting.OAuth2Setting>({
    url: '/oauth2/settings',
    method: 'GET',
  })
}

// 更新 OAuth2 设置
export function fetchUpdateOAuth2Settings(oauth2Setting: App.Api.Setting.OAuth2Setting) {
  return request({
    url: '/oauth2/settings',
    method: 'PUT',
    data: oauth2Setting,
  })
}
