<!-- 打印文件：只读排查转换状态，不提供下载与删除 -->
<template>
  <div class="print-file-page art-full-height">
    <FileSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <span class="text-sm text-g-600">
            文件内容属于会员，后台只做状态排查；过期文件由每日任务自动清理
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
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/hooks/core/useTable'
  import { fetchGetPrintFileList } from '@/api/print/file'
  import { ElTag } from 'element-plus'
  import FileSearch from './modules/file-search.vue'

  defineOptions({ name: 'PrintFile' })

  type FileItem = Api.Print.FileItem

  const STATUS_TEXT: Record<number, string> = { 0: '待处理', 1: '转换中', 2: '可打印', 3: '转换失败' }
  const STATUS_TYPE: Record<number, 'warning' | 'success' | 'danger' | 'info'> = {
    0: 'warning',
    1: 'warning',
    2: 'success',
    3: 'danger'
  }

  const searchForm = ref<Api.Print.FileSearchParams>({
    memberId: undefined,
    fileName: undefined,
    convertStatus: undefined
  })

  const formatSize = (size: number): string => {
    if (!size) return '0 B'
    if (size < 1024) return `${size} B`
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
    return `${(size / 1024 / 1024).toFixed(1)} MB`
  }

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
      apiFn: fetchGetPrintFileList,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        { prop: 'fileName', label: '文件名', minWidth: 220, showOverflowTooltip: true },
        { prop: 'fileExt', label: '格式', width: 90, formatter: (row: FileItem) => row.fileExt.toUpperCase() },
        { prop: 'fileSize', label: '大小', width: 110, formatter: (row: FileItem) => formatSize(row.fileSize) },
        { prop: 'pageCount', label: '页数', width: 80 },
        {
          prop: 'convertStatus',
          label: '转换状态',
          width: 110,
          formatter: (row: FileItem) =>
            h(ElTag, { type: STATUS_TYPE[row.convertStatus] ?? 'info', effect: 'plain' }, () =>
              STATUS_TEXT[row.convertStatus] ?? '未知'
            )
        },
        {
          prop: 'convertError',
          label: '失败原因',
          minWidth: 200,
          showOverflowTooltip: true,
          formatter: (row: FileItem) => row.convertError || '-'
        },
        { prop: 'memberUserId', label: '会员 ID', width: 190 },
        {
          prop: 'expireTime',
          label: '过期时间',
          width: 180,
          formatter: (row: FileItem) => row.expireTime ?? '-'
        },
        {
          prop: 'createdAt',
          label: '上传时间',
          width: 180,
          formatter: (row: FileItem) => row.createdAt ?? '-'
        }
      ]
    }
  })

  const handleSearch = (params: Api.Print.FileSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }
</script>
