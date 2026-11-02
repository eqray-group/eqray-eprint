<!-- 门店成员：绑定后台账号后，该账号登录客户端才能登记设备与核销取件 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="`门店成员${shopName ? ' · ' + shopName : ''}`"
    width="680px"
    align-center
  >
    <ElAlert
      type="info"
      :closable="false"
      show-icon
      title="客户端使用后台账号登录；未绑定门店的账号无法登记设备、无法核销取件"
      class="mb-4"
    />

    <div class="flex items-end gap-3 mb-4">
      <ElForm :inline="true" :model="bindForm" class="flex-1">
        <ElFormItem label="后台账号">
          <ElSelect
            v-model="bindForm.userId"
            filterable
            remote
            reserve-keyword
            clearable
            :remote-method="searchCandidates"
            :loading="searchingUsers"
            placeholder="输入用户名 / 昵称 / 手机号搜索"
            style="width: 260px"
          >
            <ElOption
              v-for="item in candidates"
              :key="item.id"
              :label="candidateLabel(item)"
              :value="String(item.id)"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="角色">
          <ElSelect v-model="bindForm.role" style="width: 120px">
            <ElOption label="店主" value="owner" />
            <ElOption label="店员" value="staff" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <ElButton
        v-if="hasAuth('print:staff:create')"
        type="primary"
        :loading="binding"
        @click="handleBind"
      >
        绑定
      </ElButton>
    </div>

    <ElTable v-loading="loading" :data="staffs" border max-height="360">
      <ElTableColumn label="后台账号" min-width="200">
        <template #default="{ row }">
          <div class="text-sm">{{ row.nickname || row.username || '账号已删除' }}</div>
          <div class="text-xs text-g-500">{{ row.username || `ID ${row.userId}` }}</div>
        </template>
      </ElTableColumn>
      <ElTableColumn prop="role" label="角色" width="100">
        <template #default="{ row }">
          <ElTag :type="row.role === 'owner' ? 'success' : 'info'" effect="plain">
            {{ row.role === 'owner' ? '店主' : '店员' }}
          </ElTag>
        </template>
      </ElTableColumn>
      <ElTableColumn prop="createdAt" label="绑定时间" width="180">
        <template #default="{ row }">{{ row.createdAt ?? '-' }}</template>
      </ElTableColumn>
      <ElTableColumn label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <ElButton
            v-if="hasAuth('print:staff:delete')"
            link
            type="danger"
            @click="handleUnbind(row)"
          >
            解绑
          </ElButton>
        </template>
      </ElTableColumn>
    </ElTable>

    <template #footer>
      <ElButton @click="dialogVisible = false">关闭</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useAuth } from '@/hooks/core/useAuth'
  import {
    fetchBindPrintShopStaff,
    fetchGetPrintShopStaffs,
    fetchSearchPrintShopStaffCandidates,
    fetchUnbindPrintShopStaff
  } from '@/api/print/shop'

  defineOptions({ name: 'PrintShopStaffDialog' })

  interface Props {
    visible: boolean
    shopId: string | null
    shopName?: string
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const { hasAuth } = useAuth()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const loading = ref(false)
  const binding = ref(false)
  const staffs = ref<Api.Print.StaffItem[]>([])
  const bindForm = reactive({ userId: '', role: 'staff' })
  const candidates = ref<Api.Print.StaffCandidate[]>([])
  const searchingUsers = ref(false)

  /**
   * 远程搜索可绑定的后台账号
   * @param keyword 用户名 / 昵称 / 手机号，空串直接清空候选
   */
  const searchCandidates = async (keyword: string): Promise<void> => {
    const text = keyword.trim()
    if (!text) {
      candidates.value = []
      return
    }

    searchingUsers.value = true
    try {
      candidates.value = await fetchSearchPrintShopStaffCandidates(text)
    } catch {
      candidates.value = []
    } finally {
      searchingUsers.value = false
    }
  }

  /**
   * 下拉项文案：优先「昵称（用户名）」，没昵称时只显示用户名
   * @param item 候选账号
   */
  const candidateLabel = (item: Api.Print.StaffCandidate): string =>
    item.nickname ? `${item.nickname}（${item.username}）` : item.username

  const loadStaffs = async (): Promise<void> => {
    if (!props.shopId) return

    loading.value = true
    try {
      staffs.value = await fetchGetPrintShopStaffs(props.shopId)
    } catch {
      staffs.value = []
    } finally {
      loading.value = false
    }
  }

  watch(
    () => props.visible,
    (visible) => {
      if (!visible) return

      bindForm.userId = ''
      bindForm.role = 'staff'
      candidates.value = []
      loadStaffs()
    }
  )

  const handleBind = async (): Promise<void> => {
    if (!props.shopId) return

    const userId = bindForm.userId.trim()
    if (!userId) {
      ElMessage.warning('请选择要绑定的后台账号')
      return
    }

    binding.value = true
    try {
      await fetchBindPrintShopStaff(props.shopId, userId, bindForm.role)
      ElMessage.success('绑定成功')
      bindForm.userId = ''
      candidates.value = []
      loadStaffs()
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      binding.value = false
    }
  }

  const handleUnbind = async (row: Api.Print.StaffItem): Promise<void> => {
    if (!props.shopId) return

    try {
      const name = row.nickname || row.username || row.userId
      await ElMessageBox.confirm(
        `确定解绑账号 ${name} 吗？解绑后该账号立即失去门店权限。`,
        '解绑成员',
        {
          type: 'warning'
        }
      )
    } catch {
      return
    }

    try {
      await fetchUnbindPrintShopStaff(props.shopId, row.id)
      ElMessage.success('已解绑')
      loadStaffs()
    } catch {
      // 错误提示已由请求层统一弹出
    }
  }
</script>
