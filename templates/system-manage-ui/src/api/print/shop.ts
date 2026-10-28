import request from '@/utils/http'

// 分页查询打印门店
export function fetchGetPrintShopList(params: Api.Print.ShopSearchParams) {
  return request.get<{ list: Api.Print.ShopItem[]; total: number }>({
    url: '/api/v1/print/shops',
    params
  })
}

// 门店详情
export function fetchGetPrintShop(id: string) {
  return request.get<Api.Print.ShopItem>({
    url: `/api/v1/print/shops/${id}`
  })
}

// 新增门店
export function fetchCreatePrintShop(data: Api.Print.ShopSaveParams) {
  return request.post<{ id: string }>({
    url: '/api/v1/print/shops',
    params: data
  })
}

// 编辑门店
export function fetchUpdatePrintShop(id: string, data: Api.Print.ShopSaveParams) {
  return request.put<null>({
    url: `/api/v1/print/shops/${id}`,
    params: data
  })
}

// 删除门店
export function fetchDeletePrintShop(id: string) {
  return request.del<null>({
    url: `/api/v1/print/shops/${id}`
  })
}

// 切换营业状态
export function fetchChangePrintShopStatus(id: string, status: number) {
  return request.put<null>({
    url: `/api/v1/print/shops/${id}/status`,
    params: { status }
  })
}

// 审核门店：1 通过 / 2 拒绝
export function fetchAuditPrintShop(id: string, auditStatus: number, remark = '') {
  return request.put<null>({
    url: `/api/v1/print/shops/${id}/audit`,
    params: { auditStatus, remark }
  })
}

// 门店成员
export function fetchGetPrintShopStaffs(id: string) {
  return request.get<Api.Print.StaffItem[]>({
    url: `/api/v1/print/shops/${id}/staffs`
  })
}

// 绑定成员
export function fetchBindPrintShopStaff(id: string, userId: string, role: string) {
  return request.post<null>({
    url: `/api/v1/print/shops/${id}/staffs`,
    params: { userId, role }
  })
}

// 解绑成员
export function fetchUnbindPrintShopStaff(id: string, staffId: string) {
  return request.del<null>({
    url: `/api/v1/print/shops/${id}/staffs/${staffId}`
  })
}
