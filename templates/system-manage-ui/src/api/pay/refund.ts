import request from '@/utils/http'

// 分页查询退款单
export function fetchGetRefundList(params: Api.Pay.RefundSearchParams) {
  return request.get<Api.Pay.RefundList>({
    url: '/api/v1/pay/refunds',
    params
  })
}

// 发起退款：余额支付立即到账，微信/支付宝按渠道处理
export function fetchCreateRefund(data: Api.Pay.RefundCreateParams) {
  return request.post<Api.Pay.RefundItem>({
    url: '/api/v1/pay/refunds',
    params: data
  })
}
