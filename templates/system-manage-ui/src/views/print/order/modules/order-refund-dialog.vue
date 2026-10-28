<!-- 订单退款：全额退，金额取订单实付，不接受手工输入 -->
<template>
  <ElDialog v-model="dialogVisible" title="订单退款" width="480px" align-center :close-on-click-modal="false">
    <ElAlert
      type="warning"
      :closable="false"
      show-icon
      title="全额退款不可撤销；余额支付立即到账，微信 / 支付宝到账以渠道回调为准"
      class="mb-4"
    />

    <ElForm label-width="90px">
      <ElFormItem label="订单号">
        <ElInput :model-value="order?.orderNo ?? '-'" disabled />
      </ElFormItem>
      <ElFormItem label="退款金额">
        <span class="text-base font-medium">{{ money(order?.payAmount ?? 0) }} 元</span>
      </ElFormItem>
      <ElFormItem label="退款原因" required>
        <ElInput v-model="reason" type="textarea" :rows="2" maxlength="255" placeholder="必填，会记录在订单与退款单上" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">确认退款</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { fetchRefundPrintOrder } from '@/api/print/order'

  defineOptions({ name: 'PrintOrderRefundDialog' })

  interface Props {
    visible: boolean
    order: Api.Print.OrderItem | null
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
  const reason = ref('')

  const money = (value: number): string => (value / 100).toFixed(2)

  watch(
    () => props.visible,
    (visible) => {
      if (visible) reason.value = ''
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!props.order) return
    if (!reason.value.trim()) {
      ElMessage.warning('请填写退款原因')
      return
    }

    submitting.value = true
    try {
      await fetchRefundPrintOrder(props.order.id, reason.value.trim())
      ElMessage.success('退款已提交')
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
