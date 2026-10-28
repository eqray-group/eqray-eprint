<!-- 打印门店：营业状态与审核状态分离，审核通过且营业中才可接单 -->
<template>
  <div class="print-shop-page art-full-height">
    <ShopSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('print:shop:create')" type="primary" @click="openDialog(null, false)">
            新增门店
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

      <ShopDialog v-model:visible="dialogVisible" :shop-id="currentShopId" @success="refreshData" />
      <ShopAuditDialog v-model:visible="auditVisible" :shop-id="currentShopId" @success="refreshData" />
      <ShopStaffDialog
        v-model:visible="staffVisible"
        :shop-id="currentShopId"
        :shop-name="currentShopName"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import {
    fetchChangePrintShopStatus,
    fetchDeletePrintShop,
    fetchGetPrintShopList
  } from '@/api/print/shop'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import ShopSearch from './modules/shop-search.vue'
  import ShopDialog from './modules/shop-dialog.vue'
  import ShopAuditDialog from './modules/shop-audit-dialog.vue'
  import ShopStaffDialog from './modules/shop-staff-dialog.vue'

  defineOptions({ name: 'PrintShop' })

  type ShopItem = Api.Print.ShopItem

  const { hasAuth } = useAuth()

  const AUDIT_TEXT: Record<number, string> = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  const AUDIT_TYPE: Record<number, 'warning' | 'success' | 'danger'> = {
    0: 'warning',
    1: 'success',
    2: 'danger'
  }

  const searchForm = ref<Api.Print.ShopSearchParams>({
    name: undefined,
    status: undefined,
    auditStatus: undefined
  })

  const currentShopId = ref<string | null>(null)
  const currentShopName = ref('')
  const dialogVisible = ref(false)
  const auditVisible = ref(false)
  const staffVisible = ref(false)

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
      apiFn: fetchGetPrintShopList,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        { prop: 'name', label: '门店名称', minWidth: 160 },
        { prop: 'address', label: '地址', minWidth: 200, showOverflowTooltip: true },
        {
          prop: 'status',
          label: '营业',
          width: 100,
          formatter: (row: ShopItem) =>
            h(
              ElTag,
              { type: row.status === 1 ? 'success' : 'info', effect: 'plain' },
              () => (row.status === 1 ? '营业' : '停业')
            )
        },
        {
          prop: 'auditStatus',
          label: '审核',
          width: 110,
          formatter: (row: ShopItem) =>
            h(
              ElTag,
              { type: AUDIT_TYPE[row.auditStatus] ?? 'info', effect: 'plain' },
              () => AUDIT_TEXT[row.auditStatus] ?? '未知'
            )
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 180,
          formatter: (row: ShopItem) => row.createdAt ?? '-'
        },
        {
          prop: 'operation',
          label: '操作',
          width: 300,
          fixed: 'right',
          formatter: (row: ShopItem) =>
            h('div', { class: 'flex items-center flex-wrap' }, [
              ...(hasAuth('print:shop:update')
                ? [h(ArtButtonTable, { type: 'edit', text: '编辑', onClick: () => openDialog(row.id, false) })]
                : []),
              ...(hasAuth('print:shop:audit') && row.auditStatus !== 1
                ? [h(ArtButtonTable, { type: 'edit', text: '审核', onClick: () => openDialog(row.id, true) })]
                : []),
              ...(hasAuth('print:staff:list')
                ? [
                    h(ArtButtonTable, {
                      type: 'edit',
                      text: '成员',
                      onClick: () => openStaff(row)
                    })
                  ]
                : []),
              ...(hasAuth('print:shop:update')
                ? [
                    h(ArtButtonTable, {
                      type: 'edit',
                      text: row.status === 1 ? '停业' : '营业',
                      onClick: () => toggleStatus(row)
                    })
                  ]
                : []),
              ...(hasAuth('print:shop:delete')
                ? [h(ArtButtonTable, { type: 'delete', text: '删除', onClick: () => remove(row) })]
                : [])
            ])
        }
      ]
    }
  })

  // 搜索：合并条件后重新拉取
  const handleSearch = (params: Api.Print.ShopSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  // 统一打开弹窗：审核与编辑共用当前行
  const openDialog = (id: string | null, audit: boolean): void => {
    currentShopId.value = id
    nextTick(() => {
      if (audit) {
        auditVisible.value = true
      } else {
        dialogVisible.value = true
      }
    })
  }

  const openStaff = (row: ShopItem): void => {
    currentShopId.value = row.id
    currentShopName.value = row.name
    nextTick(() => {
      staffVisible.value = true
    })
  }

  const toggleStatus = async (row: ShopItem): Promise<void> => {
    const next = row.status === 1 ? 0 : 1
    try {
      await fetchChangePrintShopStatus(row.id, next)
      ElMessage.success(next === 1 ? '已开始营业' : '已停止营业')
      refreshData()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }

  const remove = async (row: ShopItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(`确定删除门店「${row.name}」吗？`, '删除门店', { type: 'warning' })
    } catch {
      return
    }

    try {
      await fetchDeletePrintShop(row.id)
      ElMessage.success('已删除')
      refreshData()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
