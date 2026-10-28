/**
 * 高德地图 JS API 加载工具
 *
 * Key 不写在前端代码里，而是放在「基础设施 → 参数配置」（infra_config）中，
 * 由 /api/v1/infra/configs/public 下发，换 Key 不用重新打包前端。
 *
 * @module utils/amap
 */
import { fetchGetPublicConfigs } from '@/api/infra/config'

/** 参数配置里存高德 Web 端 Key 的键名 */
export const AMAP_KEY_CONFIG = 'print.amap.web.key'

/** 参数配置里存高德安全密钥的键名（JS API 2.0 启用安全密钥后才需要） */
export const AMAP_SECURITY_CONFIG = 'print.amap.web.security_code'

/** 脚本地址，插件按需加载：地理编码、地点搜索、工具条 */
const AMAP_SCRIPT = 'https://webapi.amap.com/maps?v=2.0'

const AMAP_PLUGINS = 'AMap.Geocoder,AMap.PlaceSearch,AMap.Scale,AMap.ToolBar'

/** 进行中的加载任务，避免并发重复插入 script */
let loading: Promise<void> | null = null

/** 高德接口返回的凭据 */
export interface AmapCredentials {
  /** Web 端 Key，为空表示还没在参数配置里填写 */
  key: string
  /** 安全密钥，为空表示未启用 */
  securityCode: string
}

/**
 * 从参数配置读取高德 Key 与安全密钥
 * @returns key 为空时调用方应提示去「参数配置」里补全
 */
export async function fetchAmapCredentials(): Promise<AmapCredentials> {
  const values = await fetchGetPublicConfigs([AMAP_KEY_CONFIG, AMAP_SECURITY_CONFIG])

  return {
    key: (values?.[AMAP_KEY_CONFIG] ?? '').trim(),
    securityCode: (values?.[AMAP_SECURITY_CONFIG] ?? '').trim()
  }
}

/**
 * 加载高德地图 JS API，已加载完成时直接返回
 * @param credentials 参数配置里读到的 Key 与安全密钥
 */
export function loadAmap(credentials: AmapCredentials): Promise<void> {
  const amapWindow = window as any
  if (amapWindow.AMap) {
    return Promise.resolve()
  }
  if (loading) {
    return loading
  }

  // 2.0 启用安全密钥时必须先注入，再加载脚本，否则接口会报 INVALID_USER_SCODE
  if (credentials.securityCode) {
    amapWindow._AMapSecurityConfig = { securityJsCode: credentials.securityCode }
  }

  loading = new Promise<void>((resolve, reject) => {
    const script = document.createElement('script')
    script.src = `${AMAP_SCRIPT}&key=${encodeURIComponent(credentials.key)}&plugin=${AMAP_PLUGINS}`
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => {
      // 加载失败要把缓存清掉，否则重试会一直拿到失败的 Promise
      loading = null
      script.remove()
      reject(new Error('高德地图脚本加载失败，请检查网络或参数配置里的 Key'))
    }
    document.head.appendChild(script)
  })

  return loading
}
