<!-- 会员卡列表：卡自带等级与折扣，与会员等级体系无关 -->
<template>
  <div class="member-card-page art-full-height">
    <CardSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <span class="text-sm text-g-600">
            会员卡由后台发放，折扣在计价时按卡上的折扣百分比生效，停用或过期即恢复原价
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

      <CardDialog v-model:visible="dialogVisible" :card-data="currentCard" @success="refreshUpdate" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchDisableMemberCard, fetchGetMemberCardList } from '@/api/member/card'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import CardSearch from './modules/card-search.vue'
  import CardDialog from './modules/card-dialog.vue'

  defineOptions({ name: 'MemberCard' })

  type CardItem = Api.Member.CardItem

  const { hasAuth } = useAuth()

  // 卡状态：0 停用 / 1 正常 / 2 已过期
  const STATUS_TEXT: Record<number, string> = { 0: '停用', 1: '正常', 2: '已过期' }
  const STATUS_TYPE: Record<number, 'success' | 'info' | 'warning'> = {
    0: 'info',
    1: 'success',
    2: 'warning'
  }

  const searchForm = ref<Api.Member.CardSearchParams>({
    cardNo: undefined,
    status: undefined
  })

  const dialogVisible = ref(false)
  const currentCard = ref<CardItem | null>(null)

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
    refreshData,
    refreshUpdate
  } = useTable({
    core: {
      apiFn: fetchGetMemberCardList,
      apiParams: {
        pageNo: 1,
        pageSize: 20,
        ...searchForm.value
      },
      columnsFactory: () => [
        { prop: 'cardNo', label: '卡号', width: 200 },
        { prop: 'memberUserId', label: '会员 ID', width: 200 },
        {
          prop: 'cardName',
          label: '卡名称',
          minWidth: 140,
          formatter: (row: CardItem) => `${row.cardName}（${row.cardLevel}）`
        },
        {
          prop: 'discountPercent',
          label: '折扣',
          width: 100,
          formatter: (row: CardItem) =>
            row.discountPercent === 100 ? '无折扣' : `${row.discountPercent / 10} 折`
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row: CardItem) =>
            h(
              ElTag,
              { type: STATUS_TYPE[row.status] ?? 'info', effect: 'plain' },
              () => STATUS_TEXT[row.status] ?? '未知'
            )
        },
        {
          prop: 'expireTime',
          label: '有效期',
          width: 180,
          formatter: (row: CardItem) => (row.expireTime ? row.expireTime.slice(0, 10) : '永久有效')
        },
        {
          prop: 'operation',
          label: '操作',
          width: 200,
          fixed: 'right',
          formatter: (row: CardItem) =>
            h('div', { class: 'flex items-center' }, [
              ...(hasAuth('member:card:update')
                ? [
                    h(ArtButtonTable, {
                      type: 'edit',
                      text: '编辑',
                      onClick: () => showDialog(row)
                    })
                  ]
                : []),
              ...(hasAuth('member:card:delete') && row.status !== 0
                ? [
                    h(ArtButtonTable, {
                      type: 'delete',
                      text: '停用',
                      onClick: () => disableCard(row)
                    })
                  ]
                : [])
            ])
        }
      ]
    }
  })

  // 搜索：把条件合并进分页参数后重新拉取
  const handleSearch = (params: Api.Member.CardSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  // 打开编辑弹窗
  const showDialog = (row: CardItem): void => {
    currentCard.value = row
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  // 停用会员卡：停用后折扣立即失效
  const disableCard = async (row: CardItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(`确定要停用卡「${row.cardNo}」吗？停用后不再享受折扣。`, '停用会员卡', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return
    }

    try {
      await fetchDisableMemberCard(row.memberUserId)
      ElMessage.success('已停用')
      refreshUpdate()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
