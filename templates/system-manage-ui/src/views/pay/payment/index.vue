<!-- 支付流水：只读查询，支持按会员发起钱包调账 -->
<template>
  <div class="pay-payment-page art-full-height">
    <PaymentSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <span class="text-sm text-g-600">
            支付单由会员端创建，后台只做查询；金额单位均为分
          </span>
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

      <WalletAdjustDialog v-model:visible="adjustVisible" :member-id="adjustMemberId" @success="refreshData" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchGetPaymentList } from '@/api/pay/payment'
  import { ElTag } from 'element-plus'
  import PaymentSearch from './modules/payment-search.vue'
  // 调账弹窗由「会员钱包」页维护，这里复用同一份，避免两处表单逻辑漂移
  import WalletAdjustDialog from '@/views/pay/wallet/modules/wallet-adjust-dialog.vue'

  defineOptions({ name: 'PayPayment' })

  type PaymentItem = Api.Pay.PaymentItem

  const { hasAuth } = useAuth()

  const CHANNEL_TEXT: Record<string, string> = { wechat: '微信', alipay: '支付宝', balance: '钱包余额' }

  const STATUS_TEXT: Record<number, string> = {
    0: '待支付',
    1: '支付中',
    2: '已支付',
    3: '支付失败',
    4: '已关闭',
    5: '已退款'
  }
  const STATUS_TYPE: Record<number, 'warning' | 'success' | 'danger' | 'info'> = {
    0: 'warning',
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'info',
    5: 'info'
  }

  const BIZ_TEXT: Record<string, string> = {
    print_order: '打印订单',
    wallet_recharge: '钱包充值'
  }

  const searchForm = ref<Api.Pay.PaymentSearchParams>({
    paymentNo: undefined,
    memberId: undefined,
    channel: undefined,
    status: undefined
  })

  const adjustVisible = ref(false)
  const adjustMemberId = ref('')

  // 分转元展示
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
      apiFn: fetchGetPaymentList,
      apiParams: {
        pageNo: 1,
        pageSize: 20,
        ...searchForm.value
      },
      columnsFactory: () => [
        { prop: 'paymentNo', label: '支付单号', width: 210 },
        {
          prop: 'bizType',
          label: '业务类型',
          width: 120,
          formatter: (row: PaymentItem) => BIZ_TEXT[row.bizType] ?? row.bizType
        },
        { prop: 'memberUserId', label: '会员 ID', width: 190 },
        {
          prop: 'channel',
          label: '渠道',
          width: 110,
          formatter: (row: PaymentItem) =>
            CHANNEL_TEXT[row.channel] ?? row.channel
        },
        { prop: 'channelCode', label: '配置编码', width: 130 },
        {
          prop: 'amount',
          label: '应付(元)',
          width: 110,
          formatter: (row: PaymentItem) => money(row.amount)
        },
        {
          prop: 'paidAmount',
          label: '实付(元)',
          width: 110,
          formatter: (row: PaymentItem) => money(row.paidAmount)
        },
        {
          prop: 'refundAmount',
          label: '已退(元)',
          width: 110,
          formatter: (row: PaymentItem) => money(row.refundAmount)
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row: PaymentItem) =>
            h(
              ElTag,
              { type: STATUS_TYPE[row.status] ?? 'info', effect: 'plain' },
              () => STATUS_TEXT[row.status] ?? '未知'
            )
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 180,
          formatter: (row: PaymentItem) => row.createdAt ?? '-'
        },
        {
          prop: 'operation',
          label: '操作',
          width: 120,
          fixed: 'right',
          formatter: (row: PaymentItem) =>
            hasAuth('pay:wallet:adjust')
              ? h('div', { class: 'flex items-center' }, [
                  h(ArtButtonTable, {
                    type: 'edit',
                    text: '调账',
                    onClick: () => showAdjust(row.memberUserId)
                  })
                ])
              : '-'
        }
      ]
    }
  })

  // 搜索：合并条件后重新拉取
  const handleSearch = (params: Api.Pay.PaymentSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  // 打开调账弹窗，会员 ID 取自当前行
  const showAdjust = (memberId: string): void => {
    adjustMemberId.value = memberId
    nextTick(() => {
      adjustVisible.value = true
    })
  }
</script>
