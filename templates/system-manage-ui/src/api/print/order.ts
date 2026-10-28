import request from '@/utils/http'

// 分页查询打印订单
export function fetchGetPrintOrderList(params: Api.Print.OrderSearchParams) {
  return request.get<{ list: Api.Print.OrderItem[]; total: number }>({
    url: '/api/v1/print/orders',
    params
  })
}

// 订单详情（含明细与任务状态）
export function fetchGetPrintOrder(id: string) {
  return request.get<Api.Print.OrderDetail>({
    url: `/api/v1/print/orders/${id}`
  })
}

// 订单全额退款
export function fetchRefundPrintOrder(id: string, reason: string) {
  return request.post<null>({
    url: `/api/v1/print/orders/${id}/refund`,
    params: { reason }
  })
}

// 取件核销
export function fetchPickupPrintOrder(data: Api.Print.PickupParams) {
  return request.post<Api.Print.OrderItem>({
    url: '/api/v1/print/orders/pickup',
    params: data
  })
}
