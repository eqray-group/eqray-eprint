<!-- 门店位置：用高德地图选点，回填经纬度与详细地址；Key 取自参数配置 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    title="地图选点"
    width="860px"
    align-center
    :close-on-click-modal="false"
    @opened="handleOpened"
    @closed="handleClosed"
  >
    <div v-loading="loading" element-loading-text="正在加载地图…" class="shop-location">
      <ElAlert
        v-if="missingKey"
        type="warning"
        :closable="false"
        show-icon
        title="尚未配置高德地图 Key"
        description="请到「基础设施 → 参数配置」填写 print.amap.web.key 后再试"
      />

      <template v-else>
        <div class="shop-location__search">
          <ElInput
            v-model="keyword"
            placeholder="搜索地址或地点，如 广州大学城"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #append>
              <ElButton :loading="searching" @click="handleSearch">搜索</ElButton>
            </template>
          </ElInput>

          <div v-if="pois.length" class="shop-location__pois">
            <div
              v-for="(poi, index) in pois"
              :key="`${index}-${poi.name}`"
              class="shop-location__poi"
              @click="pickPoi(poi)"
            >
              <span class="shop-location__poi-name">{{ poi.name }}</span>
              <span class="shop-location__poi-address">
                {{ poi.pname || '' }}{{ poi.address || '' }}
              </span>
            </div>
          </div>
        </div>

        <div ref="mapRef" class="shop-location__map"></div>

        <div class="shop-location__result">
          <span class="text-xs text-g-500">点击地图或拖动标记选点</span>
          <span class="text-xs">
            {{
              hasPoint
                ? `经纬度：${point.longitude.toFixed(6)}, ${point.latitude.toFixed(6)}`
                : '尚未选点'
            }}
          </span>
          <span class="shop-location__address text-xs text-g-500">
            {{ point.address || '地址待解析' }}
          </span>
        </div>
      </template>
    </div>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :disabled="!hasPoint" @click="handleConfirm">确定</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { fetchAmapCredentials, loadAmap } from '@/utils/amap'

  defineOptions({ name: 'PrintShopLocation' })

  interface Props {
    visible: boolean
    /** 当前门店经度，0 表示还没选点 */
    longitude?: number
    /** 当前门店纬度，0 表示还没选点 */
    latitude?: number
    /** 当前门店详细地址，用于打开时的搜索框回填 */
    address?: string
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    /** 确认选点，回传经纬度与逆地理编码得到的地址 */
    (e: 'confirm', payload: { longitude: number; latitude: number; address: string }): void
  }

  /** 高德返回的 POI，只取用到的字段 */
  interface PoiItem {
    name: string
    address?: string
    pname?: string
    location?: { lng: number; lat: number }
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const mapRef = ref<HTMLDivElement>()
  const keyword = ref('')
  const searching = ref(false)
  const loading = ref(false)
  const missingKey = ref(false)
  const pois = ref<PoiItem[]>([])
  const point = reactive({ longitude: 0, latitude: 0, address: '' })

  const hasPoint = computed(() => point.longitude !== 0 || point.latitude !== 0)

  // 高德实例，类型由脚本动态注入，这里不做强类型约束
  let map: any = null
  let marker: any = null
  let geocoder: any = null
  let placeSearch: any = null

  // 每次打开先回填门店已有坐标，关闭后由 handleClosed 清理实例
  watch(
    () => props.visible,
    (visible) => {
      if (!visible) return

      point.longitude = Number(props.longitude) || 0
      point.latitude = Number(props.latitude) || 0
      point.address = props.address ?? ''
      keyword.value = props.address ?? ''
      pois.value = []
      missingKey.value = false
    }
  )

  /**
   * 弹窗动画结束后再初始化地图：此刻容器才有真实尺寸
   * 参数：无；返回 Promise，内部完成 Key 读取与脚本加载
   */
  const handleOpened = async (): Promise<void> => {
    loading.value = true
    try {
      const credentials = await fetchAmapCredentials()
      if (!credentials.key) {
        missingKey.value = true
        return
      }

      await loadAmap(credentials)
      initMap()
    } catch (error) {
      ElMessage.error((error as Error).message || '地图加载失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 创建地图、标记与地理编码 / 地点搜索实例
   * 参数：无；返回 void，实例保存在组件闭包里
   */
  const initMap = (): void => {
    const AMap = (window as any).AMap
    if (!AMap || !mapRef.value) return

    // 已选点复用门店坐标，未选点先落在广州，用户可直接搜索
    const center: [number, number] = hasPoint.value
      ? [point.longitude, point.latitude]
      : [113.264435, 23.129163]

    map = new AMap.Map(mapRef.value, {
      zoom: hasPoint.value ? 16 : 11,
      center,
      resizeEnable: true
    })
    map.addControl(new AMap.Scale())
    map.addControl(new AMap.ToolBar({ position: 'RB' }))

    marker = new AMap.Marker({ position: center, draggable: true, visible: hasPoint.value })
    marker.on('dragend', (event: any) => {
      moveTo(event.lnglat.getLng(), event.lnglat.getLat(), false)
    })
    map.add(marker)

    geocoder = new AMap.Geocoder({ radius: 200 })
    placeSearch = new AMap.PlaceSearch({ pageSize: 10, pageIndex: 1, extensions: 'base' })

    map.on('click', (event: any) => {
      moveTo(event.lnglat.getLng(), event.lnglat.getLat(), false)
    })
  }

  /**
   * 把标记移到目标坐标并解析地址
   * 参数 lng 经度、lat 纬度、needCenter 是否同步地图中心（搜索命中时为 true）
   */
  const moveTo = (lng: number, lat: number, needCenter: boolean): void => {
    point.longitude = Number(lng.toFixed(6))
    point.latitude = Number(lat.toFixed(6))

    if (marker) {
      marker.setPosition([lng, lat])
      marker.show()
    }
    if (needCenter && map) {
      map.setZoomAndCenter(16, [lng, lat])
    }

    resolveAddress(lng, lat)
  }

  /**
   * 逆地理编码：把坐标换成「省市区 + 门牌」的详细地址
   * 参数 lng 经度、lat 纬度；返回 void，结果异步写回 point.address
   */
  const resolveAddress = (lng: number, lat: number): void => {
    if (!geocoder) return

    geocoder.getAddress([lng, lat], (status: string, result: any) => {
      if (status !== 'complete' || !result?.regeocode) return
      point.address = result.regeocode.formattedAddress ?? ''
    })
  }

  /**
   * 关键词搜索地点
   * 参数：无（读取 keyword）；返回 void，结果写入 pois 供下拉选择
   */
  const handleSearch = (): void => {
    const text = keyword.value.trim()
    if (!text || !placeSearch) return

    searching.value = true
    placeSearch.search(text, (status: string, result: any) => {
      searching.value = false
      if (status !== 'complete' || !result?.poiList?.pois?.length) {
        pois.value = []
        ElMessage.warning('没有搜索到相关地点，换个关键词试试')
        return
      }

      pois.value = result.poiList.pois
    })
  }

  /**
   * 选中搜索结果：定位到该 POI 并回填地址
   * 参数 poi 为高德返回的地点；返回 void
   */
  const pickPoi = (poi: PoiItem): void => {
    const lng = poi.location?.lng
    const lat = poi.location?.lat
    if (lng === undefined || lat === undefined) return

    moveTo(lng, lat, true)
    point.address = `${poi.pname ?? ''}${poi.address ?? ''}${poi.name ?? ''}`
    pois.value = []
  }

  /**
   * 确认选点：把经纬度与地址回传给门店表单
   * 参数：无；返回 void
   */
  const handleConfirm = (): void => {
    if (!hasPoint.value) return

    emit('confirm', {
      longitude: point.longitude,
      latitude: point.latitude,
      address: point.address
    })
    dialogVisible.value = false
  }

  /**
   * 关闭后销毁地图实例，避免重复打开累积监听与内存
   * 参数：无；返回 void
   */
  const handleClosed = (): void => {
    map?.destroy()
    map = null
    marker = null
    geocoder = null
    placeSearch = null
    pois.value = []
    keyword.value = ''
  }
</script>

<style scoped lang="scss">
  .shop-location {
    &__search {
      position: relative;
      margin-bottom: 12px;
    }

    &__pois {
      position: absolute;
      top: 42px;
      right: 0;
      left: 0;
      z-index: 20;
      max-height: 220px;
      overflow-y: auto;
      background: var(--el-bg-color-overlay);
      border: 1px solid var(--el-border-color-lighter);
      border-radius: 6px;
      box-shadow: 0 8px 24px rgb(0 0 0 / 8%);
    }

    &__poi {
      display: flex;
      flex-direction: column;
      gap: 2px;
      padding: 8px 12px;
      cursor: pointer;

      &:hover {
        background: var(--el-fill-color-light);
      }
    }

    &__poi-name {
      font-size: 13px;
      color: var(--el-text-color-primary);
    }

    &__poi-address {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }

    &__map {
      width: 100%;
      height: 380px;
      overflow: hidden;
      border-radius: 6px;
    }

    &__result {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      align-items: center;
      margin-top: 12px;
    }

    &__address {
      flex: 1;
      min-width: 200px;
    }
  }
</style>
