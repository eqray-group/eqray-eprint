import request from '@/utils/http'

// 分页查询打印设备
export function fetchGetPrintDeviceList(params: Api.Print.DeviceSearchParams) {
  return request.get<{ list: Api.Print.DeviceItem[]; total: number }>({
    url: '/api/v1/print/devices',
    params
  })
}

// 设备详情
export function fetchGetPrintDevice(id: string) {
  return request.get<Api.Print.DeviceItem>({
    url: `/api/v1/print/devices/${id}`
  })
}

// 新增设备
export function fetchCreatePrintDevice(data: Api.Print.DeviceSaveParams) {
  return request.post<{ id: string }>({
    url: '/api/v1/print/devices',
    params: data
  })
}

// 编辑设备
export function fetchUpdatePrintDevice(id: string, data: Api.Print.DeviceSaveParams) {
  return request.put<null>({
    url: `/api/v1/print/devices/${id}`,
    params: data
  })
}

// 删除设备
export function fetchDeletePrintDevice(id: string) {
  return request.del<null>({
    url: `/api/v1/print/devices/${id}`
  })
}
