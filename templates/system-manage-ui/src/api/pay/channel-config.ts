import request from '@/utils/http'

// 分页查询支付渠道配置（列表不返回密钥）
export function fetchGetChannelConfigList(params: Api.Pay.ChannelConfigSearchParams) {
  return request.get<Api.Pay.ChannelConfigList>({
    url: '/api/v1/pay/channel-configs',
    params
  })
}

// 查询渠道配置详情（编辑页回显密钥）
export function fetchGetChannelConfig(id: string) {
  return request.get<Api.Pay.ChannelConfigDetail>({
    url: `/api/v1/pay/channel-configs/${id}`
  })
}

// 新增渠道配置
export function fetchCreateChannelConfig(data: Api.Pay.ChannelConfigSaveParams) {
  return request.post<{ id: string }>({
    url: '/api/v1/pay/channel-configs',
    params: data
  })
}

// 编辑渠道配置（密钥留空表示不修改）
export function fetchUpdateChannelConfig(id: string, data: Api.Pay.ChannelConfigSaveParams) {
  return request.put<null>({
    url: `/api/v1/pay/channel-configs/${id}`,
    params: data
  })
}

// 删除渠道配置
export function fetchDeleteChannelConfig(id: string) {
  return request.del<null>({
    url: `/api/v1/pay/channel-configs/${id}`
  })
}
