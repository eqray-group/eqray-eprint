import request from '@/utils/http'

// 分页查询轮播图
export function fetchGetPrintBannerList(params: Api.Print.BannerSearchParams) {
  return request.get<{ list: Api.Print.BannerItem[]; total: number }>({
    url: '/api/v1/print/banners',
    params
  })
}

// 轮播图详情
export function fetchGetPrintBanner(id: string) {
  return request.get<Api.Print.BannerItem>({
    url: `/api/v1/print/banners/${id}`
  })
}

// 新增轮播图
export function fetchCreatePrintBanner(data: Api.Print.BannerSaveParams) {
  return request.post<{ id: string }>({
    url: '/api/v1/print/banners',
    params: data
  })
}

// 编辑轮播图
export function fetchUpdatePrintBanner(id: string, data: Api.Print.BannerSaveParams) {
  return request.put<null>({
    url: `/api/v1/print/banners/${id}`,
    params: data
  })
}

// 删除轮播图
export function fetchDeletePrintBanner(id: string) {
  return request.del<null>({
    url: `/api/v1/print/banners/${id}`
  })
}
