<!-- 价格配置：门店 × 纸张 × 色彩 × 单双面 唯一确定单价（分/页） -->
<template>
  <div class="print-price-page art-full-height">
    <PriceSearch v-model="searchForm" @search="handleSearch" @reset="handleReset" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <div class="flex items-center gap-3">
            <ElButton v-if="hasAuth('print:price:create')" type="primary" :disabled="!shopId" @click="openDialog(null)">
              新增价格
            </ElButton>
            <ElButton :disabled="!shopId" :loading="warningLoading" @click="loadWarnings">能力体检</ElButton>
            <span class="text-xs text-g-500">有价格但没有打印机支持的组合无法下单</span>
          </div>
        </template>
      </ArtTableHeader>

      <ElAlert
        v-if="warnings.length"
        type="warning"
        :closable="true"
        show-icon
        class="mb-3"
        title="以下组合已配置价格，但当前没有支持该能力的启用打印机"
        @close="warnings = []"
      >
        <div v-for="(item, index) in warnings" :key="index" class="text-xs">
          {{ describe(item.paperType, item.colorType, item.printMode) }} — {{ item.message }}
        </div>
      </ElAlert>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
      </ArtTable>

      <PriceDialog
        v-model:visible="dialogVisible"
        :price-id="currentId"
        :shop-id="shopId"
        :origin="currentRow"
        @success="refreshData"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import {
    fetchDeletePrintPrice,
    fetchGetPrintPriceList,
    fetchGetPrintPriceWarnings
  } from '@/api/print/price'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import PriceSearch from './modules/price-search.vue'
  import PriceDialog from './modules/price-dialog.vue'

  defineOptions({ name: 'PrintPrice' })

  type PriceItem = Api.Print.PriceItem

  const { hasAuth } = useAuth()

  const searchForm = ref<Api.Print.PriceSearchParams>({ printShopId: undefined })
  const shopId = computed(() => searchForm.value.printShopId ?? '')
  const warnings = ref<Api.Print.CapabilityWarning[]>([])
  const warningLoading = ref(false)
  const currentId = ref<string | null>(null)
  const currentRow = ref<PriceItem | null>(null)
  const dialogVisible = ref(false)

  // 展示用：组合转中文
  const describe = (paper: string, color: string, mode: string): string =>
    `${paper} ${color === 'color' ? '彩色' : '黑白'} ${mode === 'double' ? '双面' : '单面'}`

  const {
    columns,
    columnChecks,
    data,
    loading,
    pagination,
    getData,
    searchParams,
    handleSizeChange,
    handleCurrentChange,
    refreshData
  } = useTable({
    core: {
      apiFn: fetchGetPrintPriceList,
      apiParams: { pageNo: 1, pageSize: 50 },
      columnsFactory: () => [
        {
          prop: 'paperType',
          label: '组合',
          minWidth: 180,
          formatter: (row: PriceItem) => describe(row.paperType, row.colorType, row.printMode)
        },
        {
          prop: 'price',
          label: '单价（元/页）',
          width: 140,
          formatter: (row: PriceItem) => (row.price / 100).toFixed(2)
        },
        {
          prop: 'status',
          label: '状态',
          width: 90,
          formatter: (row: PriceItem) =>
            h(ElTag, { type: row.status === 1 ? 'success' : 'info', effect: 'plain' }, () =>
              row.status === 1 ? '启用' : '停用'
            )
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 180,
          formatter: (row: PriceItem) => row.createdAt ?? '-'
        },
        {
          prop: 'operation',
          label: '操作',
          width: 160,
          fixed: 'right',
          formatter: (row: PriceItem) =>
            h('div', { class: 'flex items-center' }, [
              ...(hasAuth('print:price:update')
                ? [h(ArtButtonTable, { type: 'edit', text: '改价', onClick: () => openDialog(row) })]
                : []),
              ...(hasAuth('print:price:delete')
                ? [h(ArtButtonTable, { type: 'delete', text: '删除', onClick: () => remove(row) })]
                : [])
            ])
        }
      ]
    }
  })

  const handleSearch = (params: Api.Print.PriceSearchParams) => {
    // 切门店要清空体检结果，否则提示会串门店
    warnings.value = []
    Object.assign(searchParams, params, { pageNo: 1 })
    getData()
  }

  const handleReset = () => {
    searchForm.value = { printShopId: undefined }
    warnings.value = []
    Object.assign(searchParams, { printShopId: undefined, pageNo: 1 })
    getData()
  }

  const loadWarnings = async (): Promise<void> => {
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

  const openDialog = (row: PriceItem | null): void => {
    currentId.value = row?.id ?? null
    currentRow.value = row
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  const remove = async (row: PriceItem): Promise<void> => {
    try {
      await ElMessageBox.confirm('删除后该组合立即不可下单，确定删除？', '删除价格规则', { type: 'warning' })
    } catch {
      return
    }

    try {
      await fetchDeletePrintPrice(row.id)
      ElMessage.success('已删除')
      refreshData()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
