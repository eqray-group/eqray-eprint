<!-- 价格配置：按 SKU 组合配置，门店 × 纸张 × 色彩 × 单双面 = 8 个可售组合 -->
<template>
  <div class="print-price-page art-full-height">
    <PriceSearch v-model="searchForm" @search="handleSearch" @reset="handleReset" />

    <ElCard class="art-table-card">
      <div class="mb-4 flex flex-wrap items-center gap-3">
        <ElButton
          v-if="canUpdate"
          type="primary"
          :disabled="!shopId || !dirty"
          :loading="saving"
          @click="handleSave"
        >
          保存价格
        </ElButton>
        <ElButton :disabled="!shopId || !dirty" @click="loadPrices">还原</ElButton>
        <ElButton :disabled="!shopId" :loading="warningLoading" @click="loadWarnings">
          能力体检
        </ElButton>
        <span class="text-xs text-g-500">单价单位：元/页；留空表示该组合不售卖</span>
      </div>

      <ElAlert
        v-if="!shopId"
        type="info"
        :closable="false"
        show-icon
        title="请先选择门店"
        description="价格按门店独立配置，选好门店后再设置各组合的单价"
      />

      <template v-else>
        <ElAlert
          v-if="warnings.length"
          type="warning"
          :closable="true"
          show-icon
          class="mb-4"
          title="以下组合已配置价格，但当前没有支持该能力的启用打印机"
          @close="warnings = []"
        >
          <div v-for="(item, index) in warnings" :key="index" class="text-xs">
            {{ describe(item.paperType, item.colorType, item.printMode) }} — {{ item.message }}
          </div>
        </ElAlert>

        <div v-loading="loading" class="price-matrix">
          <div
            v-for="row in rows"
            :key="row.key"
            class="price-card"
            :class="{ 'is-unset': !row.priceYuan }"
          >
            <div class="price-card__titles">
              <span class="price-card__title">{{ row.title }}</span>
              <ElTag
                v-if="row.id"
                :type="row.enabled ? 'success' : 'info'"
                effect="plain"
                size="small"
              >
                {{ row.enabled ? '启用' : '停用' }}
              </ElTag>
              <ElTag v-else type="warning" effect="plain" size="small">未配置</ElTag>
            </div>

            <div class="price-card__price">
              <ElInputNumber
                v-model="row.priceYuan"
                :min="0"
                :max="9999"
                :precision="2"
                :step="0.5"
                :disabled="!canUpdate"
                controls-position="right"
                placeholder="不售卖"
                class="price-card__input"
              />
              <span class="text-xs text-g-500">元/页</span>
            </div>

            <div class="price-card__foot">
              <ElSwitch
                v-model="row.enabled"
                :disabled="!canUpdate || !row.priceYuan"
                size="small"
                active-text="启用"
                inactive-text="停用"
              />
              <span class="text-xs text-g-500">
                {{
                  row.priceYuan ? `约 ${(row.priceYuan || 0).toFixed(2)} 元/页` : '填写价格后可启用'
                }}
              </span>
            </div>
          </div>
        </div>
      </template>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { useAuth } from '@/hooks/core/useAuth'
  import {
    fetchGetPrintPriceList,
    fetchGetPrintPriceWarnings,
    fetchSavePrintPriceBatch
  } from '@/api/print/price'
  import PriceSearch from './modules/price-search.vue'

  defineOptions({ name: 'PrintPrice' })

  /** 一个 SKU 组合的矩阵行，priceYuan 为 0/空表示不售卖 */
  interface MatrixRow {
    key: string
    paperType: string
    colorType: string
    printMode: string
    title: string
    /** 已有规则 ID，未配置为 null */
    id: string | null
    priceYuan: number | undefined
    enabled: boolean
    /** 载入时的值，用于判断是否有未保存改动 */
    originPrice: number | undefined
    originEnabled: boolean
  }

  /** 可售组合固定 8 个：A4/A3 × 黑白/彩色 × 单面/双面 */
  const COMBOS: Array<[string, string, string]> = [
    ['A4', 'bw', 'single'],
    ['A4', 'bw', 'double'],
    ['A4', 'color', 'single'],
    ['A4', 'color', 'double'],
    ['A3', 'bw', 'single'],
    ['A3', 'bw', 'double'],
    ['A3', 'color', 'single'],
    ['A3', 'color', 'double']
  ]

  const { hasAuth } = useAuth()
  const canUpdate = computed(() => hasAuth('print:price:update'))

  const searchForm = ref<Api.Print.PriceSearchParams>({ printShopId: undefined })
  const shopId = computed(() => searchForm.value.printShopId ?? '')
  const rows = ref<MatrixRow[]>(buildRows([]))
  const loading = ref(false)
  const saving = ref(false)
  const warnings = ref<Api.Print.CapabilityWarning[]>([])
  const warningLoading = ref(false)

  /** 有未保存改动才允许保存/还原 */
  const dirty = computed(() =>
    rows.value.some((row) => row.priceYuan !== row.originPrice || row.enabled !== row.originEnabled)
  )

  // 展示用：组合转中文
  function describe(paper: string, color: string, mode: string): string {
    return `${paper} ${color === 'color' ? '彩色' : '黑白'} ${mode === 'double' ? '双面' : '单面'}`
  }

  // 组合唯一键，与后端「门店 + 纸张 + 色彩 + 单双面」唯一约束对应
  function ruleKey(paper: string, color: string, mode: string): string {
    return `${paper}|${color}|${mode}`
  }

  // 分转元，保留两位，避免 0.1 + 0.2 这类浮点误差
  function toYuan(fen: number): number {
    return Number((fen / 100).toFixed(2))
  }

  /**
   * 用接口返回的价格规则构建 8 行矩阵
   * @param list 门店已配置的价格规则（未配置的组合自动补空行）
   */
  function buildRows(list: Api.Print.PriceItem[]): MatrixRow[] {
    const configured = new Map(
      list.map((item) => [ruleKey(item.paperType, item.colorType, item.printMode), item])
    )

    return COMBOS.map(([paperType, colorType, printMode]) => {
      const found = configured.get(ruleKey(paperType, colorType, printMode))
      const priceYuan = found ? toYuan(found.price) : undefined
      const enabled = found ? found.status === 1 : true

      return {
        key: ruleKey(paperType, colorType, printMode),
        paperType,
        colorType,
        printMode,
        title: describe(paperType, colorType, printMode),
        id: found?.id ?? null,
        priceYuan,
        enabled,
        originPrice: priceYuan,
        originEnabled: enabled
      }
    })
  }

  /**
   * 加载当前门店的价格规则并重建矩阵
   * 参数：无；未选门店时直接清空矩阵
   */
  async function loadPrices(): Promise<void> {
    if (!shopId.value) {
      rows.value = buildRows([])
      warnings.value = []
      return
    }

    loading.value = true
    try {
      const res = await fetchGetPrintPriceList({
        printShopId: shopId.value,
        pageNo: 1,
        pageSize: 100
      })
      rows.value = buildRows(res.list)
    } catch {
      rows.value = buildRows([])
    } finally {
      loading.value = false
    }
  }

  /**
   * 保存矩阵：单价留空的组合传 0，由后端删除该组合
   * 参数：无；成功后重新拉取一次，拿到最新 ID 与状态
   */
  async function handleSave(): Promise<void> {
    if (!shopId.value) return

    saving.value = true
    try {
      await fetchSavePrintPriceBatch({
        printShopId: shopId.value,
        items: rows.value.map((row) => ({
          paperType: row.paperType,
          colorType: row.colorType,
          printMode: row.printMode,
          price: Math.round((row.priceYuan || 0) * 100),
          status: row.enabled ? 1 : 0
        }))
      })

      ElMessage.success('价格已保存')
      await loadPrices()
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      saving.value = false
    }
  }

  /**
   * 能力体检：列出「有价格但没有打印机支持」的组合
   * 参数：无；结果展示在矩阵上方的提示条
   */
  async function loadWarnings(): Promise<void> {
    if (!shopId.value) return

    warningLoading.value = true
    try {
      warnings.value = await fetchGetPrintPriceWarnings(shopId.value)
    } catch {
      warnings.value = []
    } finally {
      warningLoading.value = false
    }
  }

  function handleSearch(params: Api.Print.PriceSearchParams): void {
    // 切门店要清空体检结果，否则提示会串门店
    warnings.value = []
    searchForm.value = params
    loadPrices()
  }

  function handleReset(): void {
    searchForm.value = { printShopId: undefined }
    warnings.value = []
    loadPrices()
  }
</script>

<style scoped lang="scss">
  .price-matrix {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 12px;
  }

  .price-card {
    padding: 14px 16px;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    transition: border-color 0.2s;

    &:hover {
      border-color: var(--el-color-primary-light-5);
    }

    &.is-unset {
      background: var(--el-fill-color-lighter);
    }

    &__titles {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 10px;
    }

    &__title {
      font-size: 14px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    &__price {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 10px;
    }

    &__input {
      flex: 1;
    }

    &__foot {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
    }
  }
</style>
