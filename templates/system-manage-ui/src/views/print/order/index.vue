<!-- 打印订单：查询、详情、退款与取件核销，金额单位均为分 -->
<template>
  <div class="print-order-page art-full-height">
    <OrderSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('print:order:pickup')" type="primary" @click="openPickup">
            取件核销
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

      <OrderDetailDialog v-model:visible="detailVisible" :order-id="currentId" />
      <OrderRefundDialog v-model:visible="refundVisible" :order="currentOrder" @success="refreshData" />
      <OrderPickupDialog v-model:visible="pickupVisible" @success="refreshData" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchGetPrintOrderList } from '@/api/print/order'
  import { ElTag } from 'element-plus'
  import OrderSearch from './modules/order-search.vue'
  import OrderDetailDialog from './modules/order-detail-dialog.vue'
  import OrderRefundDialog from './modules/order-refund-dialog.vue'
  import OrderPickupDialog from './modules/order-pickup-dialog.vue'

  defineOptions({ name: 'PrintOrder' })

  type OrderItem = Api.Print.OrderItem

  const { hasAuth } = useAuth()

  // 订单状态与 print_order.order_status 一致
  const STATUS_TEXT: Record<number, string> = {
    0: '待支付',
    1: '待打印',
    2: '打印中',
    3: '待取件',
    4: '已完成',
    5: '已取消',
    6: '退款中',
    7: '已退款',
    8: '超期关闭'
  }
  const STATUS_TYPE: Record<number, 'warning' | 'success' | 'danger' | 'info' | 'primary'> = {
    0: 'warning',
    1: 'primary',
    2: 'warning',
    3: 'success',
    4: 'success',
    5: 'info',
    6: 'warning',
    7: 'info',
    8: 'info'
  }

  const PAY_TEXT: Record<number, string> = { 0: '未支付', 1: '已支付', 2: '退款中', 3: '已退款' }

  const searchForm = ref<Api.Print.OrderSearchParams>({
    orderNo: undefined,
    printShopId: undefined,
    orderStatus: undefined,
    payStatus: undefined
  })

  const currentId = ref<string | null>(null)
  const currentOrder = ref<OrderItem | null>(null)
  const detailVisible = ref(false)
  const refundVisible = ref(false)
  const pickupVisible = ref(false)

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
      apiFn: fetchGetPrintOrderList,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        { prop: 'orderNo', label: '订单号', width: 200 },
        {
          prop: 'pickupCode',
          label: '取件码',
          width: 100,
          formatter: (row: OrderItem) => row.pickupCode || '-'
        },
        { prop: 'payAmount', label: '实付(元)', width: 110, formatter: (row: OrderItem) => money(row.payAmount) },
        { prop: 'itemCount', label: '文件数', width: 90 },
        {
          prop: 'orderStatus',
          label: '订单状态',
          width: 110,
          formatter: (row: OrderItem) =>
            h(ElTag, { type: STATUS_TYPE[row.orderStatus] ?? 'info', effect: 'plain' }, () =>
              STATUS_TEXT[row.orderStatus] ?? '未知'
            )
        },
        {
          prop: 'payStatus',
          label: '支付',
          width: 100,
          formatter: (row: OrderItem) => PAY_TEXT[row.payStatus] ?? '-'
        },
        {
          prop: 'createdAt',
          label: '下单时间',
          width: 180,
          formatter: (row: OrderItem) => row.createdAt ?? '-'
        },
        {
          prop: 'operation',
          label: '操作',
          width: 180,
          fixed: 'right',
          formatter: (row: OrderItem) =>
            h('div', { class: 'flex items-center' }, [
              h(ArtButtonTable, { type: 'edit', text: '详情', onClick: () => openDetail(row) }),
              ...(hasAuth('print:order:refund') && row.payStatus === 1
                ? [h(ArtButtonTable, { type: 'delete', text: '退款', onClick: () => openRefund(row) })]
                : [])
            ])
        }
      ]
    }
  })

  const handleSearch = (params: Api.Print.OrderSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  const openDetail = (row: OrderItem): void => {
    currentId.value = row.id
    nextTick(() => {
      detailVisible.value = true
    })
  }

  const openRefund = (row: OrderItem): void => {
    currentOrder.value = row
    nextTick(() => {
      refundVisible.value = true
    })
  }

  const openPickup = (): void => {
    nextTick(() => {
      pickupVisible.value = true
    })
  }
</script>
