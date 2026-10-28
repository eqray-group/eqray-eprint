<!-- 打印设备：Windows 客户端实例，客户端登录后会自动登记，这里用于预录与维护 -->
<template>
  <div class="print-device-page art-full-height">
    <DeviceSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('print:device:create')" type="primary" @click="openDialog(null)">
            新增设备
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

      <DeviceDialog v-model:visible="dialogVisible" :device-id="currentId" @success="refreshData" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchDeletePrintDevice, fetchGetPrintDeviceList } from '@/api/print/device'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import DeviceSearch from './modules/device-search.vue'
  import DeviceDialog from './modules/device-dialog.vue'

  defineOptions({ name: 'PrintDevice' })

  type DeviceItem = Api.Print.DeviceItem

  const { hasAuth } = useAuth()

  const searchForm = ref<Api.Print.DeviceSearchParams>({
    keyword: undefined,
    printShopId: undefined,
    status: undefined,
    onlineStatus: undefined
  })

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
      apiFn: fetchGetPrintDeviceList,
      apiParams: { pageNo: 1, pageSize: 20, ...searchForm.value },
      columnsFactory: () => [
        { prop: 'name', label: '设备名称', minWidth: 160 },
        { prop: 'code', label: '设备编码', minWidth: 180, showOverflowTooltip: true },
        { prop: 'printShopId', label: '门店 ID', width: 190 },
        {
          prop: 'clientVersion',
          label: '客户端版本',
          width: 120,
          formatter: (row: DeviceItem) => row.clientVersion || '-'
        },
        {
          prop: 'onlineStatus',
          label: '在线',
          width: 90,
          formatter: (row: DeviceItem) =>
            h(ElTag, { type: row.onlineStatus === 1 ? 'success' : 'info', effect: 'plain' }, () =>
              row.onlineStatus === 1 ? '在线' : '离线'
            )
        },
        {
          prop: 'lastHeartbeatAt',
          label: '最后心跳',
          width: 180,
          formatter: (row: DeviceItem) => row.lastHeartbeatAt ?? '-'
        },
        {
          prop: 'status',
          label: '状态',
          width: 90,
          formatter: (row: DeviceItem) =>
            h(ElTag, { type: row.status === 1 ? 'success' : 'danger', effect: 'plain' }, () =>
              row.status === 1 ? '启用' : '禁用'
            )
        },
        {
          prop: 'operation',
          label: '操作',
          width: 160,
          fixed: 'right',
          formatter: (row: DeviceItem) =>
            h('div', { class: 'flex items-center' }, [
              ...(hasAuth('print:device:update')
                ? [h(ArtButtonTable, { type: 'edit', text: '编辑', onClick: () => openDialog(row.id) })]
                : []),
              ...(hasAuth('print:device:delete')
                ? [h(ArtButtonTable, { type: 'delete', text: '删除', onClick: () => remove(row) })]
                : [])
            ])
        }
      ]
    }
  })

  const handleSearch = (params: Api.Print.DeviceSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  const openDialog = (id: string | null): void => {
    currentId.value = id
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  const remove = async (row: DeviceItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(`确定删除设备「${row.name}」吗？其下打印机需先停用。`, '删除设备', {
        type: 'warning'
      })
    } catch {
      return
    }

    try {
      await fetchDeletePrintDevice(row.id)
      ElMessage.success('已删除')
      refreshData()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
