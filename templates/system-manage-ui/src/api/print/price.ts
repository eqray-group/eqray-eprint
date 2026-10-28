import request from '@/utils/http'

// 分页查询价格规则
export function fetchGetPrintPriceList(params: Api.Print.PriceSearchParams) {
  return request.get<{ list: Api.Print.PriceItem[]; total: number }>({
    url: '/api/v1/print/price-configs',
    params
  })
}

// 新增价格规则
export function fetchCreatePrintPrice(data: Api.Print.PriceSaveParams) {
  return request.post<{ id: string }>({
    url: '/api/v1/print/price-configs',
    params: data
  })
}

// 修改价格与状态
export function fetchUpdatePrintPrice(id: string, price: number, status: number) {
  return request.put<null>({
    url: `/api/v1/print/price-configs/${id}`,
    params: { price, status }
  })
}

// 删除价格规则
export function fetchDeletePrintPrice(id: string) {
  return request.del<null>({
    url: `/api/v1/print/price-configs/${id}`
  })
}

// 能力体检：有价格但没有打印机支持的组合
export function fetchGetPrintPriceWarnings(printShopId: string) {
  return request.get<Api.Print.CapabilityWarning[]>({
    url: '/api/v1/print/price-configs/warnings',
    params: { printShopId }
  })
}
