<!-- 门店审核：通过后才对会员端可见，拒绝需填原因 -->
<template>
  <ElDialog v-model="dialogVisible" title="门店审核" width="480px" align-center :close-on-click-modal="false">
    <ElAlert
      type="info"
      :closable="false"
      show-icon
      title="审核通过后不可重复审核；拒绝需填写原因，门店可修改后再次提交"
      class="mb-4"
    />

    <ElForm label-width="90px">
      <ElFormItem label="审核结果">
        <ElRadioGroup v-model="form.auditStatus">
          <ElRadio :value="1">通过</ElRadio>
          <ElRadio :value="2">拒绝</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="审核意见">
        <ElInput v-model="form.remark" type="textarea" :rows="2" maxlength="255" placeholder="拒绝时建议填写原因" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">提交审核</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { fetchAuditPrintShop } from '@/api/print/shop'

  defineOptions({ name: 'PrintShopAuditDialog' })

  interface Props {
    visible: boolean
    shopId: string | null
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const submitting = ref(false)
  const form = reactive({ auditStatus: 1, remark: '' })

  watch(
    () => props.visible,
    (visible) => {
      if (!visible) return

      Object.assign(form, { auditStatus: 1, remark: '' })
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!props.shopId) return
    if (form.auditStatus === 2 && !form.remark.trim()) {
      ElMessage.warning('拒绝时请填写审核意见')
      return
    }

    submitting.value = true
    try {
      await fetchAuditPrintShop(props.shopId, form.auditStatus, form.remark.trim())
      ElMessage.success('审核已提交')
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
