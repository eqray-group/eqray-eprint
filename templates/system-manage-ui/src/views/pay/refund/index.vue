<!-- 退款管理：查询退款单并发起退款 -->
<template>
  <div class="pay-refund-page art-full-height">
    <RefundSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('pay:refund:create')" type="primary" @click="openDialog">
            发起退款
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

      <RefundDialog v-model:visible="dialogVisible" @success="refreshData" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchGetRefundList } from '@/api/pay/refund'
  import { ElButton, ElTag } from 'element-plus'
  import RefundSearch from './modules/refund-search.vue'
  import RefundDialog from './modules/refund-dialog.vue'

  defineOptions({ name: 'PayRefund' })

  type RefundItem = Api.Pay.RefundItem

  const { hasAuth } = useAuth()

  const CHANNEL_TEXT: Record<string, string> = { wechat: '微信', alipay: '支付宝', balance: '钱包余额' }

  const STATUS_TEXT: Record<number, string> = {
    0: '待退款',
    1: '退款中',
    2: '退款成功',
    3: '退款失败'
  }
  const STATUS_TYPE: Record<number, 'warning' | 'success' | 'danger' | 'info'> = {
    0: 'warning',
    1: 'warning',
    2: 'success',
    3: 'danger'
  }

  const searchForm = ref<Api.Pay.RefundSearchParams>({
    paymentNo: undefined,
    memberId: undefined,
    status: undefined
  })

  const dialogVisible = ref(false)

  const money = (value: number): string => (value / 100).toFixed(2)

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
      apiFn: fetchGetRefundList,
      apiParams: {
        pageNo: 1,
        pageSize: 20,
        ...searchForm.value
      },
      columnsFactory: () => [
        { prop: 'refundNo', label: '退款单号', width: 210 },
        { prop: 'paymentNo', label: '支付单号', width: 210 },
        { prop: 'memberUserId', label: '会员 ID', width: 190 },
        {
          prop: 'channel',
          label: '渠道',
          width: 110,
          formatter: (row: RefundItem) => CHANNEL_TEXT[row.channel] ?? row.channel
        },
        {
          prop: 'amount',
          label: '退款(元)',
          width: 110,
          formatter: (row: RefundItem) => money(row.amount)
        },
        {
          prop: 'status',
          label: '状态',
          width: 110,
          formatter: (row: RefundItem) =>
            h(
              ElTag,
              { type: STATUS_TYPE[row.status] ?? 'info', effect: 'plain' },
              () => STATUS_TEXT[row.status] ?? '未知'
            )
        },
        { prop: 'reason', label: '退款原因', minWidth: 160, showOverflowTooltip: true },
        {
          prop: 'failReason',
          label: '失败原因',
          minWidth: 160,
          showOverflowTooltip: true,
          formatter: (row: RefundItem) => row.failReason || '-'
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 180,
          formatter: (row: RefundItem) => row.createdAt ?? '-'
        }
      ]
    }
  })

  // 搜索：合并条件后重新拉取
  const handleSearch = (params: Api.Pay.RefundSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  // 打开退款弹窗
  const openDialog = (): void => {
    dialogVisible.value = true
  }
</script>
