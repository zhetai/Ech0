import { request } from '../request'

// 分页获取Echos
export function fetchGetEchosByPage(searchParams: App.Api.Ech0.ParamsByPagination) {
  return request<App.Api.Ech0.PaginationResult>({
    url: `/messages/page`,
    method: 'POST',
    data: searchParams,
  })
}

// 获取status
export function fetchGetStatus() {
  return request<App.Api.Ech0.Status>({
    url: `/status`,
    method: 'GET',
  })
}

// 上传图片
export function fetchUploadImage(file: File) {
  const formData = new FormData()
  formData.append('image', file)
  return request<string>({
    url: `/images/upload`,
    method: 'POST',
    data: formData,
  })
}

// 添加Echo
export function fetchAddEcho(echoToAdd: App.Api.Ech0.EchoToAdd) {
  return request({
    url: `/messages`,
    method: 'POST',
    data: echoToAdd,
  })
}

// 删除Echo
export function fetchDeleteEcho(echoId: number) {
  return request({
    url: `/messages/${echoId}`,
    method: 'DELETE',
  })
}
