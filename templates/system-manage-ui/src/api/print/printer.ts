import request from '@/utils/http'

// 分页查询打印机
export function fetchGetPrintPrinterList(params: Api.Print.PrinterSearchParams) {
  return request.get<{ list: Api.Print.PrinterItem[]; total: number }>({
    url: '/api/v1/print/printers',
    params
  })
}

// 打印机详情
export function fetchGetPrintPrinter(id: string) {
  return request.get<Api.Print.PrinterItem>({
    url: `/api/v1/print/printers/${id}`
  })
}

// 新增打印机
export function fetchCreatePrintPrinter(data: Api.Print.PrinterSaveParams) {
  return request.post<{ id: string }>({
    url: '/api/v1/print/printers',
    params: data
  })
}

// 编辑打印机
export function fetchUpdatePrintPrinter(id: string, data: Api.Print.PrinterSaveParams) {
  return request.put<null>({
    url: `/api/v1/print/printers/${id}`,
    params: data
  })
}

// 停用打印机
export function fetchDisablePrintPrinter(id: string) {
  return request.put<null>({
    url: `/api/v1/print/printers/${id}/disable`
  })
}
