<!-- 支付渠道配置：微信三种场景 + 支付宝，密钥仅在编辑页回显 -->
<template>
  <div class="pay-channel-page art-full-height">
    <ChannelConfigSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-if="hasAuth('pay:channel:create')" type="primary" @click="showDialog(null)">
            新增渠道
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

      <ChannelConfigDialog
        v-model:visible="dialogVisible"
        :config-id="currentId"
        @success="refreshUpdate"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { useAuth } from '@/hooks/core/useAuth'
  import { fetchDeleteChannelConfig, fetchGetChannelConfigList } from '@/api/pay/channel-config'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import ChannelConfigSearch from './modules/channel-config-search.vue'
  import ChannelConfigDialog from './modules/channel-config-dialog.vue'

  defineOptions({ name: 'PayChannelConfig' })

  type ConfigItem = Api.Pay.ChannelConfigItem

  const { hasAuth } = useAuth()

  // 渠道与场景的展示文案，与后端 pay_channel_config 取值一致
  const CHANNEL_TEXT: Record<string, string> = { wechat: '微信', alipay: '支付宝' }
  const SCENE_TEXT: Record<string, string> = {
    mini: '小程序',
    jsapi: '公众号',
    h5: '手机浏览器',
    wap: '手机网站'
  }

  const searchForm = ref<Api.Pay.ChannelConfigSearchParams>({
    code: undefined,
    channel: undefined,
    status: undefined
  })

  const dialogVisible = ref(false)
  const currentId = ref<string | null>(null)

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
      apiFn: fetchGetChannelConfigList,
      apiParams: {
        pageNo: 1,
        pageSize: 20,
        ...searchForm.value
      },
      columnsFactory: () => [
        { prop: 'code', label: '配置编码', width: 150 },
        { prop: 'name', label: '名称', minWidth: 160 },
        {
          prop: 'channel',
          label: '渠道',
          width: 100,
          formatter: (row: ConfigItem) =>
            h(
              ElTag,
              { type: row.channel === 'wechat' ? 'success' : 'primary', effect: 'plain' },
              () => CHANNEL_TEXT[row.channel] ?? row.channel
            )
        },
        {
          prop: 'scene',
          label: '场景',
          width: 120,
          formatter: (row: ConfigItem) => SCENE_TEXT[row.scene] ?? row.scene
        },
        { prop: 'appId', label: 'AppID', minWidth: 180, showOverflowTooltip: true },
        {
          prop: 'status',
          label: '状态',
          width: 90,
          formatter: (row: ConfigItem) =>
            h(
              ElTag,
              { type: row.status === 1 ? 'success' : 'info', effect: 'plain' },
              () => (row.status === 1 ? '启用' : '停用')
            )
        },
        { prop: 'sort', label: '排序', width: 80 },
        {
          prop: 'operation',
          label: '操作',
          width: 160,
          fixed: 'right',
          formatter: (row: ConfigItem) =>
            h('div', { class: 'flex items-center' }, [
              ...(hasAuth('pay:channel:update')
                ? [
                    h(ArtButtonTable, {
                      type: 'edit',
                      text: '编辑',
                      onClick: () => showDialog(row.id)
                    })
                  ]
                : []),
              ...(hasAuth('pay:channel:delete')
                ? [
                    h(ArtButtonTable, {
                      type: 'delete',
                      text: '删除',
                      onClick: () => removeConfig(row)
                    })
                  ]
                : [])
            ])
        }
      ]
    }
  })

  // 搜索：合并条件后重新拉取
  const handleSearch = (params: Api.Pay.ChannelConfigSearchParams) => {
    Object.assign(searchParams, params)
    getData()
  }

  // 打开弹窗：id 为空表示新增
  const showDialog = (id: string | null): void => {
    currentId.value = id
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  // 删除渠道配置：已产生的支付单仍保留 channelCode 可追溯
  const removeConfig = async (row: ConfigItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(`确定删除渠道配置「${row.code}」吗？`, '删除支付渠道', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return
    }

    try {
      await fetchDeleteChannelConfig(row.id)
      ElMessage.success('已删除')
      refreshUpdate()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
