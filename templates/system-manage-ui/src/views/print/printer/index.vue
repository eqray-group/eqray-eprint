<!-- 打印机：一条记录 = 某台设备上的一个 Windows 打印队列 -->
<template>
  <div class="print-printer-page art-full-height">
    <PrinterSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('print:printer:create')" type="primary" @click="openDialog(null)">
            新增打印机
          </ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
      </ArtTable>

      <PrinterDialog v-model:visible="dialogVisible" :printer-id="currentId" @success="refreshData" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchDisablePrintPrinter, fetchGetPrintPrinterList } from '@/api/print/printer'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import PrinterSearch from './modules/printer-search.vue'
  import PrinterDialog from './modules/printer-dialog.vue'

  defineOptions({ name: 'PrintPrinter' })

  type PrinterItem = Api.Print.PrinterItem

  const { hasAuth } = useAuth()

  const searchForm = ref<Api.Print.PrinterSearchParams>({ printShopId: undefined, deviceId: undefined })
  const currentId = ref<string | null>(null)
  const dialogVisible = ref(false)

  const {
    columns,
    columnChecks,
    data,
    loading,
    pagination,
    getData,
    searchParams,
    resetSearchParams,
    handleSizeChange,
    handleCurrentChange,
    refreshData
  } = useTable({
    core: {
      apiFn: fetchGetPrintPrinterList,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        { prop: 'name', label: '打印机名称', minWidth: 160 },
        { prop: 'systemName', label: 'Windows 队列', minWidth: 180, showOverflowTooltip: true },
        {
          prop: 'paperTypes',
          label: '纸张',
          width: 110,
          formatter: (row: PrinterItem) => row.paperTypes || '-'
        },
        {
          prop: 'colorSupport',
          label: '彩色',
          width: 90,
          formatter: (row: PrinterItem) =>
            h(ElTag, { type: row.colorSupport === 1 ? 'success' : 'info', effect: 'plain' }, () =>
              row.colorSupport === 1 ? '支持' : '不支持'
            )
        },
        {
          prop: 'duplexSupport',
          label: '双面',
          width: 90,
          formatter: (row: PrinterItem) =>
            h(ElTag, { type: row.duplexSupport === 1 ? 'success' : 'info', effect: 'plain' }, () =>
              row.duplexSupport === 1 ? '支持' : '不支持'
            )
        },
        {
          prop: 'status',
          label: '状态',
          width: 90,
          formatter: (row: PrinterItem) =>
            h(ElTag, { type: row.status === 1 ? 'success' : 'info', effect: 'plain' }, () =>
              row.status === 1 ? '启用' : '停用'
            )
        },
        {
          prop: 'operation',
          label: '操作',
          width: 160,
          fixed: 'right',
          formatter: (row: PrinterItem) =>
            h('div', { class: 'flex items-center' }, [
              ...(hasAuth('print:printer:update')
                ? [h(ArtButtonTable, { type: 'edit', text: '编辑', onClick: () => openDialog(row.id) })]
                : []),
              ...(hasAuth('print:printer:delete') && row.status === 1
                ? [h(ArtButtonTable, { type: 'delete', text: '停用', onClick: () => disable(row) })]
                : [])
            ])
        }
      ]
    }
  })

  const handleSearch = (params: Api.Print.PrinterSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  const openDialog = (id: string | null): void => {
    currentId.value = id
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  const disable = async (row: PrinterItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(
        `停用后「${row.name}」不再接单，已下单未打印的任务需后台处理。`,
        '停用打印机',
        { type: 'warning' }
      )
    } catch {
      return
    }

    try {
      await fetchDisablePrintPrinter(row.id)
      ElMessage.success('已停用')
      refreshData()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
