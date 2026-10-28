<!-- 门店成员：绑定后台账号后，该账号登录客户端才能登记设备与核销取件 -->
<template>
  <ElDialog v-model="dialogVisible" :title="`门店成员${shopName ? ' · ' + shopName : ''}`" width="680px" align-center>
    <ElAlert
      type="info"
      :closable="false"
      show-icon
      title="客户端使用后台账号登录；未绑定门店的账号无法登记设备、无法核销取件"
      class="mb-4"
    />

    <div class="flex items-end gap-3 mb-4">
      <ElForm :inline="true" :model="bindForm" class="flex-1">
        <ElFormItem label="后台账号 ID">
          <ElInput v-model="bindForm.userId" placeholder="system_user.id" style="width: 180px" />
        </ElFormItem>
        <ElFormItem label="角色">
          <ElSelect v-model="bindForm.role" style="width: 120px">
            <ElOption label="店主" value="owner" />
            <ElOption label="店员" value="staff" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <ElButton v-if="hasAuth('print:staff:create')" type="primary" :loading="binding" @click="handleBind">
        绑定
      </ElButton>
    </div>

    <ElTable v-loading="loading" :data="staffs" border max-height="360">
      <ElTableColumn prop="userId" label="后台账号 ID" min-width="180" />
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
      loadStaffs()
    }
  )

  const handleBind = async (): Promise<void> => {
    if (!props.shopId) return

    const userId = bindForm.userId.trim()
    if (!userId) {
      ElMessage.warning('请输入后台账号 ID')
      return
    }

    binding.value = true
    try {
      await fetchBindPrintShopStaff(props.shopId, userId, bindForm.role)
      ElMessage.success('绑定成功')
      bindForm.userId = ''
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
      await ElMessageBox.confirm(`确定解绑账号 ${row.userId} 吗？解绑后该账号立即失去门店权限。`, '解绑成员', {
        type: 'warning'
      })
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
