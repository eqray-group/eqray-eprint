<!-- 轮播图：门店自有 + 全平台，会员端按 sort 升序展示 -->
<template>
  <div class="print-banner-page art-full-height">
    <BannerSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('print:banner:create')" type="primary" @click="openDialog(null)">
            新增轮播图
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

      <BannerDialog v-model:visible="dialogVisible" :banner-id="currentId" @success="refreshData" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchDeletePrintBanner, fetchGetPrintBannerList } from '@/api/print/banner'
  import { ElImage, ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import BannerSearch from './modules/banner-search.vue'
  import BannerDialog from './modules/banner-dialog.vue'

  defineOptions({ name: 'PrintBanner' })

  type BannerItem = Api.Print.BannerItem

  const { hasAuth } = useAuth()

  const LINK_TEXT: Record<string, string> = { none: '不跳转', page: '站内页', url: '外部链接' }

  const searchForm = ref<Api.Print.BannerSearchParams>({ printShopId: undefined })
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
      apiFn: fetchGetPrintBannerList,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        {
          prop: 'imageUrl',
          label: '图片',
          width: 140,
          formatter: (row: BannerItem) =>
            h(ElImage, {
              src: row.imageUrl,
              previewSrcList: [row.imageUrl],
              fit: 'cover',
              style: 'width: 96px; height: 54px; border-radius: 4px'
            })
        },
        { prop: 'title', label: '标题', minWidth: 140, formatter: (row: BannerItem) => row.title || '-' },
        {
          prop: 'printShopId',
          label: '范围',
          width: 110,
          formatter: (row: BannerItem) => (row.printShopId === '0' ? '全平台' : `门店 ${row.printShopId}`)
        },
        {
          prop: 'linkType',
          label: '跳转',
          width: 110,
          formatter: (row: BannerItem) => LINK_TEXT[row.linkType] ?? row.linkType
        },
        { prop: 'sort', label: '排序', width: 80 },
        {
          prop: 'status',
          label: '状态',
          width: 90,
          formatter: (row: BannerItem) =>
            h(ElTag, { type: row.status === 1 ? 'success' : 'info', effect: 'plain' }, () =>
              row.status === 1 ? '启用' : '停用'
            )
        },
        {
          prop: 'operation',
          label: '操作',
          width: 160,
          fixed: 'right',
          formatter: (row: BannerItem) =>
            h('div', { class: 'flex items-center' }, [
              ...(hasAuth('print:banner:update')
                ? [h(ArtButtonTable, { type: 'edit', text: '编辑', onClick: () => openDialog(row.id) })]
                : []),
              ...(hasAuth('print:banner:delete')
                ? [h(ArtButtonTable, { type: 'delete', text: '删除', onClick: () => remove(row) })]
                : [])
            ])
        }
      ]
    }
  })

  const handleSearch = (params: Api.Print.BannerSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  const openDialog = (id: string | null): void => {
    currentId.value = id
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  const remove = async (row: BannerItem): Promise<void> => {
    try {
      await ElMessageBox.confirm('删除后会员端立即不再展示该图，确定删除？', '删除轮播图', { type: 'warning' })
    } catch {
      return
    }

    try {
      await fetchDeletePrintBanner(row.id)
      ElMessage.success('已删除')
      refreshData()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
