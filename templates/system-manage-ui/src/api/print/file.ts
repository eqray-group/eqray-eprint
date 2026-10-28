import request from '@/utils/http'

// 分页查询打印文件（只读）
export function fetchGetPrintFileList(params: Api.Print.FileSearchParams) {
  return request.get<{ list: Api.Print.FileItem[]; total: number }>({
    url: '/api/v1/print/files',
    params
  })
}

// 打印业务概览
export function fetchGetPrintStatistics() {
  return request.get<Api.Print.StatisticsOverview>({
    url: '/api/v1/print/statistics/overview'
  })
}
