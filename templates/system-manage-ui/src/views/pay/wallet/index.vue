<!-- 会员钱包：按会员看余额与流水，可人工调账；金额单位均为分 -->
<template>
  <div class="pay-wallet-page art-full-height">
    <WalletSearch v-model="searchForm" @search="handleSearch" @reset="handleReset" />

    <ElRow v-if="memberId" :gutter="12" class="mb-3">
      <ElCol :span="8">
        <ElCard shadow="never">
          <div class="text-xs text-g-500">当前余额（元）</div>
          <div class="mt-1 text-2xl font-medium">{{ money(wallet?.balance ?? 0) }}</div>
        </ElCard>
      </ElCol>
      <ElCol :span="8">
        <ElCard shadow="never">
          <div class="text-xs text-g-500">累计收入（元）</div>
          <div class="mt-1 text-2xl font-medium text-success">{{ money(wallet?.totalIncome ?? 0) }}</div>
        </ElCard>
      </ElCol>
      <ElCol :span="8">
        <ElCard shadow="never">
          <div class="text-xs text-g-500">累计支出（元）</div>
          <div class="mt-1 text-2xl font-medium text-danger">{{ money(wallet?.totalExpense ?? 0) }}</div>
        </ElCard>
      </ElCol>
    </ElRow>

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton
            v-if="hasAuth('pay:wallet:adjust')"
            type="primary"
            :disabled="!memberId"
            @click="openAdjust"
          >
            钱包调账
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

      <WalletAdjustDialog
        v-model:visible="adjustVisible"
        :member-id="memberId"
        @success="handleAdjustSuccess"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchGetWallet, fetchGetWalletTransactions } from '@/api/pay/wallet'
  import { ElTag } from 'element-plus'
  import WalletSearch from './modules/wallet-search.vue'
  import WalletAdjustDialog from './modules/wallet-adjust-dialog.vue'

  defineOptions({ name: 'PayWallet' })

  type TransactionItem = Api.Pay.WalletTransactionItem

  const { hasAuth } = useAuth()

  const BIZ_TEXT: Record<string, string> = {
    wallet_recharge: '钱包充值',
    print_consume: '打印消费',
    print_refund: '打印退款',
    manual_adjust: '后台调账'
  }

  const searchForm = ref<Api.Pay.WalletTransactionSearchParams>({
    memberId: undefined,
    type: undefined,
    bizType: undefined
  })

  const wallet = ref<Api.Pay.WalletItem | null>(null)
  const adjustVisible = ref(false)

  const memberId = computed(() => searchForm.value.memberId?.trim() ?? '')
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
      apiFn: fetchGetWalletTransactions,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        {
          prop: 'createdAt',
          label: '时间',
          width: 180,
          formatter: (row: TransactionItem) => row.createdAt ?? '-'
        },
        { prop: 'memberUserId', label: '会员 ID', width: 190 },
        {
          prop: 'type',
          label: '收支',
          width: 90,
          formatter: (row: TransactionItem) =>
            h(
              ElTag,
              { type: row.type === 1 ? 'success' : 'warning', effect: 'plain' },
              () => (row.type === 1 ? '收入' : '支出')
            )
        },
        {
          prop: 'bizType',
          label: '业务',
          width: 120,
          formatter: (row: TransactionItem) => BIZ_TEXT[row.bizType] ?? row.bizType
        },
        { prop: 'title', label: '摘要', minWidth: 160, showOverflowTooltip: true },
        {
          prop: 'amount',
          label: '金额(元)',
          width: 120,
          formatter: (row: TransactionItem) =>
            `${row.type === 1 ? '+' : '-'}${money(row.amount)}`
        },
        {
          prop: 'balanceAfter',
          label: '变动后余额(元)',
          width: 150,
          formatter: (row: TransactionItem) => money(row.balanceAfter)
        },
        { prop: 'remark', label: '备注', minWidth: 150, showOverflowTooltip: true }
      ]
    }
  })

  // 查余额：会员 ID 变化时刷新汇总卡片
  const loadWallet = async (): Promise<void> => {
    if (!memberId.value) {
      wallet.value = null
      return
    }

    try {
      wallet.value = await fetchGetWallet(memberId.value)
    } catch {
      wallet.value = null
    }
  }

  const handleSearch = (params: Api.Pay.WalletTransactionSearchParams) => {
    Object.assign(searchParams, params, { pageNo: 1 })
    getData()
    loadWallet()
  }

  const handleReset = () => {
    searchForm.value = { memberId: undefined, type: undefined, bizType: undefined }
    wallet.value = null
    resetSearchParams()
  }

  const openAdjust = (): void => {
    if (!memberId.value) return
    nextTick(() => {
      adjustVisible.value = true
    })
  }

  // 调账后余额与流水都要重查，否则数字对不上
  const handleAdjustSuccess = (): void => {
    refreshData()
    loadWallet()
  }
</script>
