import request from '@/utils/request'

export function getperimision(params) {
    return request({
        url: '/api/v1/user/permission',
        method: 'get',
        params: params
    })
}

export function set_permission(data) {
    return request({
      url: '/api/v1/user/set_permission',
      method: 'post',
      data: data
    })
}

  export function del_permission(data) {
    return request({
      url: '/api/v1/user/set_permission',
      method: 'delete',
      data: data
    })
  }