import request from '@/utils/request'


export function gethost(params) {
    return request({
        url: '/api/v1/hosts',
        method: 'get',
        params: params
    })
}
